// Package iface 定义了与数据库交互的接口。
// Package iface defines the interface for interacting with a database.
package iface

import (
	"context"
	"database/sql"

	"github.com/turtacn/mha4rdb/internal/core/types"
	"github.com/turtacn/mha4rdb/internal/core/types/enum"
)

// DB 定义了数据库操作的接口，用于MHA Agent与Vastbase等数据库进行交互。
// DB defines the interface for database operations, used by MHA Agent to interact with databases like Vastbase.
type DB interface {
	// Connect 连接到数据库。
	// Connect connects to the database.
	Connect(ctx context.Context) error

	// Close 关闭数据库连接。
	// Close closes the database connection.
	Close(ctx context.Context) error

	// Ping 检查数据库连接是否仍然存活。
	// Ping checks if the database connection is still alive.
	Ping(ctx context.Context) error

	// GetRole 获取当前数据库实例的角色（例如Primary, Secondary）。
	// GetRole retrieves the current role of the database instance (e.g., Primary, Secondary).
	GetRole(ctx context.Context) (enum.DBRole, error)

	// IsHealthy 检查数据库实例是否健康。
	// IsHealthy checks if the database instance is healthy.
	// 这可能包括检查只读状态、复制延迟等。
	// This might include checking read-only status, replication lag, etc.
	IsHealthy(ctx context.Context) (bool, error)

	// PromoteToLeader 将当前数据库实例提升为Leader/Primary。
	// PromoteToLeader promotes the current database instance to Leader/Primary.
	PromoteToLeader(ctx context.Context) error

	// DemoteToFollower 将当前数据库实例降级为Follower/Secondary。
	// DemoteToFollower demotes the current database instance to Follower/Secondary.
	DemoteToFollower(ctx context.Context, newLeaderAddress string) error

	// GetRaftStatus 获取与Raft相关的状态信息（如果数据库原生支持Raft或有相关视图）。
	// GetRaftStatus retrieves Raft-related status information (if the DB natively supports Raft or has relevant views).
	// 返回的信息结构可能包含任期、提交索引、Leader ID等。
	// The returned info structure might include term, commit index, leader ID, etc.
	GetRaftStatus(ctx context.Context) (*types.RaftInfo, error)

	// GetReplicationStatus 获取复制状态信息。
	// GetReplicationStatus retrieves replication status information.
	// 例如，对于异步/半同步复制，可能包括上游节点、延迟、应用位点等。
	// For example, for async/semi-sync replication, this might include upstream node, lag, apply LSN, etc.
	GetReplicationStatus(ctx context.Context) (interface{}, error) // interface{} can be a specific struct later

	// ExecuteQuery 执行一个只读查询并返回结果。
	// ExecuteQuery executes a read-only query and returns the results.
	// 适用于简单的状态检查或数据检索。
	// Suitable for simple status checks or data retrieval.
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// ExecuteCommand 执行一个写命令（如DDL/DML，但不包括事务控制）。
	// ExecuteCommand executes a write command (like DDL/DML, but not transaction control).
	ExecuteCommand(ctx context.Context, command string, args ...interface{}) (sql.Result, error)

	// GetPrimaryGTIDSet 获取主节点的GTID集合信息（如果适用，如MySQL）。
	// GetPrimaryGTIDSet retrieves the primary node's GTID set information (if applicable, e.g., for MySQL).
	// 对于Vastbase，可能是等效的LSN或事务ID信息。
	// For Vastbase, this might be equivalent LSN or transaction ID information.
	GetPrimaryGTIDSet(ctx context.Context) (string, error)

	// GetNodeID 获取数据库节点的唯一标识符。
	// GetNodeID retrieves the unique identifier of the database node.
	GetNodeID(ctx context.Context) (string, error)

	// StartFailoverPromotion 在故障转移场景下提升节点为主节点。
	// StartFailoverPromotion promotes the node to primary in a failover scenario.
	// 这可能是一个更激进的提升操作。
	// This might be a more aggressive promotion operation.
	StartFailoverPromotion(ctx context.Context) error

	// StopReplication 停止与主节点的复制。
	// StopReplication stops replication with the primary node.
	StopReplication(ctx context.Context) error

	// StartReplicationWith 启动与指定主节点的复制。
	// StartReplicationWith starts replication with the specified primary node.
	StartReplicationWith(ctx context.Context, primaryHost string, primaryPort int, user, password string) error

	// SetReadOnly 设置数据库为只读模式。
	// SetReadOnly sets the database to read-only mode.
	SetReadOnly(ctx context.Context, readOnly bool) error

	// IsReadOnly 检查数据库是否处于只读模式。
	// IsReadOnly checks if the database is in read-only mode.
	IsReadOnly(ctx context.Context) (bool, error)

	// GetConnection 获取底层的 *sql.DB 连接对象，用于特殊操作。
	// GetConnection returns the underlying *sql.DB connection object for special operations.
	// 使用时需谨慎，因为它绕过了接口的抽象。
	// Use with caution as it bypasses the interface abstraction.
	GetConnection() *sql.DB
}