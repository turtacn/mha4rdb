// Package vastbase 包含了与Vastbase数据库交互的具体实现。
// Package vastbase contains the specific implementation for interacting with a Vastbase database.
package vastbase

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	// 注册Vastbase (PostgreSQL) 驱动
	// Register Vastbase (PostgreSQL) driver
	_ "github.com/lib/pq"

	"github.com/turtacn/mha4rdb/internal/core/types"
	"github.com/turtacn/mha4rdb/internal/core/types/enum"
	"github.com/turtacn/mha4rdb/internal/database/iface"
	"github.com/turtacn/mha4rdb/internal/errors"
	logiface "github.com/turtacn/mha4rdb/internal/logger/iface"
)

// Compile-time check to ensure vastbaseDB implements iface.DB
var _ iface.DB = (*vastbaseDB)(nil)

// vastbaseDB 结构体实现了DB接口，用于与Vastbase数据库交互。
// vastbaseDB struct implements the DB interface for interacting with a Vastbase database.
type vastbaseDB struct {
	config Config
	db     *sql.DB
	logger logiface.Logger
	mu     sync.RWMutex // 用于保护db连接的并发访问 (For protecting concurrent access to db connection)
}

// NewVastbaseDB 创建一个新的vastbaseDB实例。
// NewVastbaseDB creates a new vastbaseDB instance.
func NewVastbaseDB(cfg Config, logger logiface.Logger) (iface.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.New("invalid Vastbase configuration: " + err.Error())
	}
	return &vastbaseDB{
		config: cfg,
		logger: logger.WithFields(logiface.Fields{"module": "vastbase", "host": cfg.Host, "port": cfg.Port}),
	}, nil
}

// Connect 实现 iface.DB 接口。
// Connect implements the iface.DB interface.
func (v *vastbaseDB) Connect(ctx context.Context) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.db != nil {
		// 尝试Ping现有连接
		// Try to ping existing connection
		if err := v.db.PingContext(ctx); err == nil {
			v.logger.Debugf("Already connected to Vastbase and connection is alive.")
			return nil
		}
		v.logger.Warnf("Existing connection to Vastbase is dead, attempting to reconnect.")
		// 关闭死连接，忽略错误，因为我们要重新连接
		// Close dead connection, ignore error as we are reconnecting
		_ = v.db.Close()
		v.db = nil
	}

	dsn := v.config.DSN()
	v.logger.Infof("Connecting to Vastbase with DSN: %s (password omitted)", v.filterPasswordFromDSN(dsn))

	var err error
	v.db, err = sql.Open("postgres", dsn) // "postgres" 驱动通常兼容Vastbase ("postgres" driver is usually compatible with Vastbase)
	if err != nil {
		v.logger.Errorf("Failed to open Vastbase connection: %v", err)
		return errors.ErrDBConnectionFailed
	}

	v.db.SetMaxOpenConns(v.config.MaxOpenConns)
	v.db.SetMaxIdleConns(v.config.MaxIdleConns)
	v.db.SetConnMaxLifetime(v.config.ConnMaxLifetime)
	v.db.SetConnMaxIdleTime(v.config.ConnMaxIdleTime)

	ctxTimeout, cancel := context.WithTimeout(ctx, v.config.ConnectTimeout)
	defer cancel()

	if err = v.db.PingContext(ctxTimeout); err != nil {
		v.logger.Errorf("Failed to ping Vastbase after opening connection: %v", err)
		// 关闭失败的连接
		// Close failed connection
		_ = v.db.Close()
		v.db = nil
		return errors.ErrDBConnectionFailed
	}

	v.logger.Infof("Successfully connected to Vastbase at %s:%d", v.config.Host, v.config.Port)
	return nil
}

// filterPasswordFromDSN 从DSN字符串中移除密码信息，用于日志记录。
// filterPasswordFromDSN removes password information from a DSN string for logging purposes.
func (v *vastbaseDB) filterPasswordFromDSN(dsn string) string {
	if strings.Contains(dsn, "password=") {
		parts := strings.Split(dsn, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, "password=") {
				parts[i] = "password=********"
			}
		}
		return strings.Join(parts, " ")
	}
	return dsn
}

// Close 实现 iface.DB 接口。
// Close implements the iface.DB interface.
func (v *vastbaseDB) Close(ctx context.Context) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.db != nil {
		v.logger.Infof("Closing Vastbase connection to %s:%d", v.config.Host, v.config.Port)
		err := v.db.Close()
		v.db = nil // 确保db被置为nil，即使Close返回错误 (Ensure db is set to nil even if Close returns an error)
		if err != nil {
			v.logger.Errorf("Error closing Vastbase connection: %v", err)
			return err
		}
		return nil
	}
	v.logger.Debugf("Vastbase connection to %s:%d was already closed.", v.config.Host, v.config.Port)
	return nil
}

// Ping 实现 iface.DB 接口。
// Ping implements the iface.DB interface.
func (v *vastbaseDB) Ping(ctx context.Context) error {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		v.logger.Warnf("Ping failed: Vastbase connection is not established.")
		return errors.ErrDBConnectionFailed
	}
	err := db.PingContext(ctx)
	if err != nil {
		v.logger.Warnf("Ping to Vastbase failed: %v", err)
	}
	return err
}

// GetRole 实现 iface.DB 接口。
// GetRole implements the iface.DB interface.
func (v *vastbaseDB) GetRole(ctx context.Context) (enum.DBRole, error) {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		return enum.RoleUnknownDB, errors.ErrDBConnectionFailed
	}

	// Vastbase特定的查询来确定角色
	// Vastbase-specific query to determine the role
	// 这通常涉及到查询 pg_is_in_recovery() 或类似的函数/视图
	// This usually involves querying pg_is_in_recovery() or similar functions/views
	query := "SELECT pg_is_in_recovery();"
	var isInRecovery bool

	err := db.QueryRowContext(ctx, query).Scan(&isInRecovery)
	if err != nil {
		v.logger.Errorf("Failed to query database role: %v", err)
		return enum.RoleUnknownDB, errors.ErrDBQueryFailed
	}

	if isInRecovery {
		return enum.RoleSecondary, nil
	}
	return enum.RolePrimary, nil
}

// IsHealthy 实现 iface.DB 接口。
// IsHealthy implements the iface.DB interface.
func (v *vastbaseDB) IsHealthy(ctx context.Context) (bool, error) {
	// 首先Ping数据库
	// First, ping the database
	if err := v.Ping(ctx); err != nil {
		v.logger.Warnf("Health check failed: Ping error: %v", err)
		return false, err
	}

	// 检查角色是否可确定
	// Check if role can be determined
	_, err := v.GetRole(ctx)
	if err != nil {
		v.logger.Warnf("Health check failed: Could not determine DB role: %v", err)
		return false, err
	}

	// TODO: 添加更多Vastbase特定的健康检查，例如：
	// - 检查是否有过多的死锁 (check for excessive deadlocks)
	// - 检查复制延迟 (check replication lag)
	// - 检查连接数是否接近上限 (check if connection count is near limit)
	// - 检查是否有长时间运行的查询 (check for long-running queries)

	v.logger.Debugf("Database instance at %s:%d is healthy.", v.config.Host, v.config.Port)
	return true, nil
}

// PromoteToLeader 实现 iface.DB 接口。
// PromoteToLeader implements the iface.DB interface.
func (v *vastbaseDB) PromoteToLeader(ctx context.Context) error {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		return errors.ErrDBConnectionFailed
	}

	// Vastbase提升为Primary的命令
	// Command to promote Vastbase to Primary
	// 这通常是 `pg_promote()` 或类似的特定命令
	// This is typically `pg_promote()` or a similar specific command
	// 例如: SELECT pg_promote(true, 300); (等待300秒)
	// Example: SELECT pg_promote(true, 300); (wait for 300 seconds)
	// 或者使用特定于Vastbase集群管理工具的命令
	// Or use commands specific to Vastbase cluster management tools

	// 示例（需要根据Vastbase实际API调整）
	// Example (needs adjustment based on actual Vastbase API)
	// query := "SELECT pg_promote(true, 60);" // Wait for 60 seconds
	// _, err := v.ExecuteCommand(ctx, query)
	// if err != nil {
	// 	v.logger.Errorf("Failed to promote to leader: %v", err)
	// 	return err
	// }
	// v.logger.Infof("Successfully promoted to leader: %s:%d", v.config.Host, v.config.Port)
	// return nil

	v.logger.Warnf("PromoteToLeader is not fully implemented for Vastbase. Placeholder used.")
	return errors.ErrNotImplemented // 返回未实现错误，直到有确切的Vastbase命令
}

// DemoteToFollower 实现 iface.DB 接口。
// DemoteToFollower implements the iface.DB interface.
func (v *vastbaseDB) DemoteToFollower(ctx context.Context, newLeaderAddress string) error {
	// 在Vastbase/PostgreSQL中，降级通常是通过重新配置recovery.conf（或等效的配置）
	// 并重启节点或重新加载配置来实现的。
	// In Vastbase/PostgreSQL, demotion is typically achieved by reconfiguring recovery.conf (or equivalent)
	// and restarting the node or reloading configuration.
	// 这通常不由SQL命令直接完成，而是通过配置文件和可能的集群管理工具。
	// This is usually not done directly by SQL commands but through configuration files and possibly cluster management tools.
	v.logger.Warnf("DemoteToFollower for Vastbase usually involves configuration changes and restart/reload, not direct SQL. Placeholder used.")
	// 可能需要停止当前实例，修改配置文件以指向新的主节点，然后重启。
	// Might need to stop the current instance, modify config file to point to new primary, then restart.
	return errors.ErrNotImplemented
}

// GetRaftStatus 实现 iface.DB 接口。
// GetRaftStatus implements the iface.DB interface.
func (v *vastbaseDB) GetRaftStatus(ctx context.Context) (*types.RaftInfo, error) {
	// 如果Vastbase原生支持Raft并通过SQL接口暴露状态，则在此处查询。
	// If Vastbase natively supports Raft and exposes status via SQL interface, query it here.
	// 否则，此信息可能由MHA Agent的Raft模块管理。
	// Otherwise, this information might be managed by the MHA Agent's Raft module.
	v.logger.Debugf("GetRaftStatus called. If Vastbase has Raft status views, query them here.")
	// 示例查询（纯占位）
	// Example query (pure placeholder)
	// query := "SELECT term, commit_index, last_applied, leader_id, state FROM vastbase_raft_status;"
	// var info types.RaftInfo
	// err := v.db.QueryRowContext(ctx, query).Scan(&info.Term, &info.CommitIndex, &info.LastAppliedIndex, &info.LeaderID, &info.State)
	// if err != nil {
	//  if err == sql.ErrNoRows {
	//      return nil, errors.ErrNotFound // Or some other appropriate error
	//  }
	// 	v.logger.Errorf("Failed to query Raft status: %v", err)
	// 	return nil, errors.ErrDBQueryFailed
	// }
	// return &info, nil
	return nil, errors.ErrNotImplemented // 假设Vastbase不直接通过此接口暴露Raft状态
}

// GetReplicationStatus 实现 iface.DB 接口。
// GetReplicationStatus implements the iface.DB interface.
func (v *vastbaseDB) GetReplicationStatus(ctx context.Context) (interface{}, error) {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		return nil, errors.ErrDBConnectionFailed
	}
	// 查询Vastbase的复制状态视图，例如 pg_stat_replication 或特定于Vastbase的视图
	// Query Vastbase's replication status views, e.g., pg_stat_replication or Vastbase-specific views
	// query := "SELECT usename, application_name, client_addr, state, sync_state, write_lag, flush_lag, replay_lag FROM pg_stat_replication;"
	// 此处返回一个map或自定义结构体
	// Return a map or a custom struct here
	v.logger.Warnf("GetReplicationStatus for Vastbase is a placeholder.")
	return nil, errors.ErrNotImplemented
}

// ExecuteQuery 实现 iface.DB 接口。
// ExecuteQuery implements the iface.DB interface.
func (v *vastbaseDB) ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		return nil, errors.ErrDBConnectionFailed
	}
	v.logger.Debugf("Executing query: %s with args: %v", query, args)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		v.logger.Errorf("Failed to execute query '%s': %v", query, err)
		return nil, errors.ErrDBQueryFailed
	}
	return rows, nil
}

// ExecuteCommand 实现 iface.DB 接口。
// ExecuteCommand implements the iface.DB interface.
func (v *vastbaseDB) ExecuteCommand(ctx context.Context, command string, args ...interface{}) (sql.Result, error) {
	v.mu.RLock()
	db := v.db
	v.mu.RUnlock()

	if db == nil {
		return nil, errors.ErrDBConnectionFailed
	}
	v.logger.Debugf("Executing command: %s with args: %v", command, args)
	result, err := db.ExecContext(ctx, command, args...)
	if err != nil {
		v.logger.Errorf("Failed to execute command '%s': %v", command, err)
		return nil, errors.ErrDBQueryFailed
	}
	return result, nil
}

// GetPrimaryGTIDSet 实现 iface.DB 接口。
// GetPrimaryGTIDSet implements the iface.DB interface.
func (v *vastbaseDB) GetPrimaryGTIDSet(ctx context.Context) (string, error) {
	// Vastbase/PostgreSQL 使用LSN (Log Sequence Number) 而不是GTID。
	// Vastbase/PostgreSQL uses LSN (Log Sequence Number) instead of GTID.
	// 可以查询当前WAL写入位置。
	// Can query the current WAL write location.
	// 例如: SELECT pg_current_wal_lsn();
	// Example: SELECT pg_current_wal_lsn();
	query := "SELECT pg_current_wal_lsn();" // For primary
	// query := "SELECT pg_last_wal_replay_lsn();" // For standby to know how much it has replayed

	var lsn string
	rows, err := v.ExecuteQuery(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&lsn); err != nil {
			v.logger.Errorf("Failed to scan LSN: %v", err)
			return "", errors.ErrDBQueryFailed
		}
		return lsn, nil
	}
	if err := rows.Err(); err != nil {
		v.logger.Errorf("Error during LSN query rows iteration: %v", err)
		return "", err
	}
	return "", errors.New("could not retrieve LSN")
}

// GetNodeID 实现 iface.DB 接口。
// GetNodeID implements the iface.DB interface.
func (v *vastbaseDB) GetNodeID(ctx context.Context) (string, error) {
	// Vastbase/PostgreSQL没有内建的全局唯一节点ID。
	// Vastbase/PostgreSQL does not have a built-in globally unique node ID.
	// 通常使用主机名或`system_identifier`（来自`pg_control`，但不容易通过SQL访问）。
	// Typically hostname or `system_identifier` (from `pg_control`, but not easily accessible via SQL) is used.
	// 我们可以返回配置的主机名和端口作为一种标识。
	// We can return the configured hostname and port as a form of identification.
	// 或者，如果集群管理工具设置了 `application_name` 或类似标识，可以使用它。
	// Alternatively, if a cluster management tool sets `application_name` or similar, that could be used.
	// 另一种方法是查询 `pg_settings` 中的 `server_uuid` (如果Vastbase版本支持)。
	// Another approach is to query `server_uuid` from `pg_settings` (if supported by Vastbase version).

	// 作为一个简单的实现，返回 host:port
	// As a simple implementation, return host:port
	nodeID := fmt.Sprintf("%s:%d", v.config.Host, v.config.Port)
	v.logger.Debugf("Using host:port as NodeID: %s", nodeID)
	return nodeID, nil
	// return "", errors.ErrNotImplemented
}

// StartFailoverPromotion 实现 iface.DB 接口。
// StartFailoverPromotion implements the iface.DB interface.
func (v *vastbaseDB) StartFailoverPromotion(ctx context.Context) error {
	// 这与 PromoteToLeader 类似，但可能包含更激进的步骤或不同的参数。
	// This is similar to PromoteToLeader, but might involve more aggressive steps or different parameters.
	// 例如，pg_promote(false) 可能用于不等待所有副本同步的场景。
	// For example, pg_promote(false) might be used for scenarios not waiting for all replicas to sync.
	v.logger.Warnf("StartFailoverPromotion is not fully implemented for Vastbase. Placeholder used.")
	return v.PromoteToLeader(ctx) // 暂时调用 PromoteToLeader
}

// StopReplication 实现 iface.DB 接口。
// StopReplication implements the iface.DB interface.
func (v *vastbaseDB) StopReplication(ctx context.Context) error {
	// 在备库上，停止复制通常意味着从 `recovery.conf` (或等效的PG12+配置) 中移除或注释掉 `primary_conninfo`，
	// 然后重启或 `pg_ctl reload`。
	// On a standby, stopping replication usually means removing or commenting out `primary_conninfo`
	// from `recovery.conf` (or equivalent PG12+ configuration) and then restarting or `pg_ctl reload`.
	// 如果是逻辑复制，则可能是 `ALTER SUBSCRIPTION ... DISABLE`。
	// If it's logical replication, it might be `ALTER SUBSCRIPTION ... DISABLE`.
	v.logger.Warnf("StopReplication for Vastbase usually involves config changes and restart/reload. Placeholder used.")
	return errors.ErrNotImplemented
}

// StartReplicationWith 实现 iface.DB 接口。
// StartReplicationWith implements the iface.DB interface.
func (v *vastbaseDB) StartReplicationWith(ctx context.Context, primaryHost string, primaryPort int, user, password string) error {
	// 配置备库以从新的主库复制。
	// Configure the standby to replicate from the new primary.
	// 这通常涉及到更新 `recovery.conf` (或等效的PG12+配置如 `standby.signal` 和 `postgresql.auto.conf`)
	// 中的 `primary_conninfo`，然后重启或 `pg_ctl reload`。
	// This typically involves updating `primary_conninfo` in `recovery.conf` (or equivalent PG12+ config
	// like `standby.signal` and `postgresql.auto.conf`) and then restarting or `pg_ctl reload`.
	// `primary_conninfo` 示例:
	// `primary_conninfo = 'host=new_primary_host port=5432 user=replication_user password=secret'`
	v.logger.Warnf("StartReplicationWith for Vastbase usually involves config changes and restart/reload. Placeholder used.")
	v.logger.Infof("Would configure replication from %s:%d with user %s", primaryHost, primaryPort, user)
	return errors.ErrNotImplemented
}

// SetReadOnly 实现 iface.DB 接口。
// SetReadOnly implements the iface.DB interface.
func (v *vastbaseDB) SetReadOnly(ctx context.Context, readOnly bool) error {
	var query string
	if readOnly {
		query = "SET default_transaction_read_only = true;"
		// 或者 ALTER SYSTEM SET default_transaction_read_only = on; 然后 pg_reload_conf();
		// Or ALTER SYSTEM SET default_transaction_read_only = on; then pg_reload_conf();
	} else {
		query = "SET default_transaction_read_only = false;"
		// 或者 ALTER SYSTEM SET default_transaction_read_only = off; 然后 pg_reload_conf();
		// Or ALTER SYSTEM SET default_transaction_read_only = off; then pg_reload_conf();
	}
	_, err := v.ExecuteCommand(ctx, query)
	if err != nil {
		v.logger.Errorf("Failed to set read-only mode to %v: %v", readOnly, err)
		return err
	}
	v.logger.Infof("Successfully set read-only mode to %v for current session/transaction on %s:%d", readOnly, v.config.Host, v.config.Port)
	// 注意: SET命令通常只对当前会话有效，或当前事务。
	// Note: SET command is usually only effective for the current session or transaction.
	// 全局设置需要 ALTER SYSTEM 和 pg_reload_conf() 或重启。
	// Global setting requires ALTER SYSTEM and pg_reload_conf() or restart.
	// 对于MHA场景，我们可能期望的是全局设置。
	// For MHA scenarios, we might expect a global setting.
	v.logger.Warnf("SetReadOnly currently affects only the current session/transaction. Global change requires ALTER SYSTEM and reload/restart.")
	return nil // 或返回 ErrNotImplemented 如果坚持要求全局更改
}

// IsReadOnly 实现 iface.DB 接口。
// IsReadOnly implements the iface.DB interface.
func (v *vastbaseDB) IsReadOnly(ctx context.Context) (bool, error) {
	query := "SHOW default_transaction_read_only;"
	var readOnlySetting string
	rows, err := v.ExecuteQuery(ctx, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&readOnlySetting); err != nil {
			v.logger.Errorf("Failed to scan read-only setting: %v", err)
			return false, errors.ErrDBQueryFailed
		}
		return strings.ToLower(readOnlySetting) == "on", nil
	}
	if err := rows.Err(); err != nil {
		v.logger.Errorf("Error during read-only setting query rows iteration: %v", err)
		return false, err
	}
	return false, errors.New("could not retrieve read-only setting")
}

// GetConnection 实现 iface.DB 接口。
// GetConnection implements the iface.DB interface.
func (v *vastbaseDB) GetConnection() *sql.DB {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.db
}