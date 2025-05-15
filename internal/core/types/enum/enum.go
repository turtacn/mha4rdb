// Package enum 定义了项目中使用的核心枚举类型。
// Package enum defines core enumerated types used in the project.
package enum

// NodeStatus 表示节点的健康状态或工作状态。
// NodeStatus represents the health or working status of a node.
type NodeStatus int

const (
	// StatusUnknown 表示未知状态。
	// StatusUnknown indicates an unknown status.
	StatusUnknown NodeStatus = iota
	// StatusInitializing 表示节点正在初始化。
	// StatusInitializing indicates the node is initializing.
	StatusInitializing
	// StatusStarting 表示节点正在启动。
	// StatusStarting indicates the node is starting.
	StatusStarting
	// StatusRunning 表示节点正在正常运行。
	// StatusRunning indicates the node is running normally.
	StatusRunning
	// StatusStopping 表示节点正在停止。
	// StatusStopping indicates the node is stopping.
	StatusStopping
	// StatusStopped 表示节点已停止。
	// StatusStopped indicates the node is stopped.
	StatusStopped
	// StatusDegraded 表示节点处于降级状态（例如，部分功能不可用）。
	// StatusDegraded indicates the node is in a degraded state (e.g., some functionalities are unavailable).
	StatusDegraded
	// StatusUnreachable 表示节点不可达。
	// StatusUnreachable indicates the node is unreachable.
	StatusUnreachable
	// StatusError 表示节点遇到错误。
	// StatusError indicates the node has encountered an error.
	StatusError
)

// String 方法返回 NodeStatus 的字符串表示。
// String method returns the string representation of NodeStatus.
func (s NodeStatus) String() string {
	switch s {
	case StatusInitializing:
		return "Initializing"
	case StatusStarting:
		return "Starting"
	case StatusRunning:
		return "Running"
	case StatusStopping:
		return "Stopping"
	case StatusStopped:
		return "Stopped"
	case StatusDegraded:
		return "Degraded"
	case StatusUnreachable:
		return "Unreachable"
	case StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// DBRole 表示数据库实例的角色。
// DBRole represents the role of a database instance.
type DBRole int

const (
	// RoleUnknownDB 表示未知的数据库角色。
	// RoleUnknownDB indicates an unknown database role.
	RoleUnknownDB DBRole = iota
	// RolePrimary 表示数据库主节点。
	// RolePrimary indicates a primary database node.
	RolePrimary
	// RoleSecondary 表示数据库备节点（或Follower）。
	// RoleSecondary indicates a secondary/standby/follower database node.
	RoleSecondary
	// RoleCandidate 表示数据库节点是选举候选者。
	// RoleCandidate indicates a database node is an election candidate.
	RoleCandidate
	// RoleArbiter 表示仲裁节点。
	// RoleArbiter indicates an arbiter node.
	RoleArbiter
)

// String 方法返回 DBRole 的字符串表示。
// String method returns the string representation of DBRole.
func (r DBRole) String() string {
	switch r {
	case RolePrimary:
		return "Primary"
	case RoleSecondary:
		return "Secondary"
	case RoleCandidate:
		return "Candidate"
	case RoleArbiter:
		return "Arbiter"
	default:
		return "UnknownDBRole"
	}
}

// RaftState 表示Raft协议中节点的状态。
// RaftState represents the state of a node in the Raft protocol.
type RaftState string

const (
	// RaftStateFollower 表示Raft跟随者状态。
	// RaftStateFollower indicates the Raft Follower state.
	RaftStateFollower RaftState = "Follower"
	// RaftStateCandidate 表示Raft候选者状态。
	// RaftStateCandidate indicates the Raft Candidate state.
	RaftStateCandidate RaftState = "Candidate"
	// RaftStateLeader 表示Raft领导者状态。
	// RaftStateLeader indicates the Raft Leader state.
	RaftStateLeader RaftState = "Leader"
	// RaftStateShutdown 表示Raft节点关闭状态。
	// RaftStateShutdown indicates the Raft Shutdown state.
	RaftStateShutdown RaftState = "Shutdown"
	// RaftStateUnknown 表示未知的Raft状态。
	// RaftStateUnknown indicates an unknown Raft state.
	RaftStateUnknown RaftState = "Unknown"
)

// OperationMode 表示MHA Agent的操作模式。
// OperationMode represents the MHA Agent's operation mode.
type OperationMode string

const (
	// ModeSingle 表示单机模式。
	// ModeSingle represents single node mode.
	ModeSingle OperationMode = "Single"
	// ModeDualArbiter 表示双机+仲裁模式。
	// ModeDualArbiter represents dual node + arbiter mode.
	ModeDualArbiter OperationMode = "DualArbiter"
	// ModeCluster 表示多机集群模式。
	// ModeCluster represents multi-node cluster mode.
	ModeCluster OperationMode = "Cluster"
)