// Package types 定义了项目中通用的数据结构。
// Package types defines common data structures used in the project.
package types

import (
	"time"

	"github.com/turtacn/mha4rdb/internal/core/types/enum"
)

// NodeInfo 存储单个节点的信息。
// NodeInfo stores information about a single node.
type NodeInfo struct {
	// ID 是节点的唯一标识符。
	// ID is the unique identifier for the node.
	ID string `json:"id"`
	// Name 是节点的名称，可以是主机名或自定义名称。
	// Name is the name of the node, can be hostname or a custom name.
	Name string `json:"name"`
	// Address 是节点的网络地址（例如 IP:Port）。
	// Address is the network address of the node (e.g., IP:Port).
	Address string `json:"address"`
	// Status 是节点的当前状态。
	// Status is the current status of the node.
	Status enum.NodeStatus `json:"status"`
	// DBRole 是节点上数据库实例的角色。
	// DBRole is the role of the database instance on this node.
	DBRole enum.DBRole `json:"db_role"`
	// RaftState 是节点的Raft状态（如果适用）。
	// RaftState is the Raft state of the node (if applicable).
	RaftState enum.RaftState `json:"raft_state,omitempty"`
	// IsLeader 表示此节点是否为Raft集群的Leader。
	// IsLeader indicates if this node is the leader of the Raft cluster.
	IsLeader bool `json:"is_leader"`
	// LastHeartbeat 是最后一次心跳时间。
	// LastHeartbeat is the time of the last heartbeat received from/by this node.
	LastHeartbeat time.Time `json:"last_heartbeat,omitempty"`
	// Tags 是节点的附加标签信息。
	// Tags are additional tags for the node.
	Tags map[string]string `json:"tags,omitempty"`
}

// ClusterStatus 存储整个集群的状态信息。
// ClusterStatus stores status information for the entire cluster.
type ClusterStatus struct {
	// ClusterID 是集群的唯一标识符。
	// ClusterID is the unique identifier for the cluster.
	ClusterID string `json:"cluster_id"`
	// LeaderID 是当前Leader节点的ID。
	// LeaderID is the ID of the current leader node.
	LeaderID string `json:"leader_id,omitempty"`
	// LeaderAddress 是当前Leader节点的地址。
	// LeaderAddress is the address of the current leader node.
	LeaderAddress string `json:"leader_address,omitempty"`
	// Nodes 是集群中所有节点的列表。
	// Nodes is a list of all nodes in the cluster.
	Nodes []NodeInfo `json:"nodes"`
	// Status 是集群的整体状态。
	// Status is the overall status of the cluster.
	OverallStatus enum.NodeStatus `json:"overall_status"`
	// Message 包含关于集群状态的附加信息。
	// Message contains additional information about the cluster status.
	Message string `json:"message,omitempty"`
	// Mode 表示集群的操作模式。
	// Mode indicates the cluster's operation mode.
	Mode enum.OperationMode `json:"mode"`
}

// RaftInfo 存储Raft相关的特定信息。
// RaftInfo stores Raft-specific information.
type RaftInfo struct {
	// Term 是当前的Raft任期。
	// Term is the current Raft term.
	Term uint64 `json:"term"`
	// CommitIndex 是已提交的日志索引。
	// CommitIndex is the commit index of the Raft log.
	CommitIndex uint64 `json:"commit_index"`
	// LastAppliedIndex 是最后应用的日志索引。
	// LastAppliedIndex is the last applied index of the Raft log.
	LastAppliedIndex uint64 `json:"last_applied_index"`
	// LeaderID 是当前Leader的ID (如果已知)。
	// LeaderID is the ID of the current leader (if known).
	LeaderID string `json:"leader_id,omitempty"` // Might be different from Raft's internal uint64 ID
	// State 是当前节点的Raft状态。
	// State is the current Raft state of the node.
	State enum.RaftState `json:"state"`
	// Peers 是集群中的其他对等节点。
	// Peers is a list of other peers in the cluster.
	Peers map[string]string `json:"peers,omitempty"` // map[NodeID]Address
}

// HealthCheckResult 表示健康检查的结果。
// HealthCheckResult represents the result of a health check.
type HealthCheckResult struct {
	// Healthy 指示目标是否健康。
	// Healthy indicates if the target is healthy.
	Healthy bool `json:"healthy"`
	// Message 包含健康检查的详细信息。
	// Message contains details from the health check.
	Message string `json:"message,omitempty"`
	// Error 是健康检查过程中发生的错误（如果有）。
	// Error is any error that occurred during the health check.
	Error error `json:"-"` // Should not be marshalled to JSON directly
	// Timestamp 是健康检查执行的时间。
	// Timestamp is when the health check was performed.
	Timestamp time.Time `json:"timestamp"`
	// Details 包含特定检查项的详细信息。
	// Details contains specific check item details.
	Details map[string]string `json:"details,omitempty"`
}

// VIPInfo 存储虚拟IP的相关信息。
// VIPInfo stores information related to the Virtual IP.
type VIPInfo struct {
	// Address 是VIP地址。
	// Address is the VIP address.
	Address string `json:"address"`
	// Interface 是VIP绑定的网络接口。
	// Interface is the network interface the VIP is bound to.
	Interface string `json:"interface"`
	// CurrentHolderNodeID 是当前持有VIP的节点ID。
	// CurrentHolderNodeID is the ID of the node currently holding the VIP.
	CurrentHolderNodeID string `json:"current_holder_node_id"`
}

// ActionResult 表示一个操作执行的结果。
// ActionResult represents the result of an action execution.
type ActionResult struct {
	// Success 表示操作是否成功。
	// Success indicates if the operation was successful.
	Success bool `json:"success"`
	// Message 包含操作结果的消息。
	// Message contains a message regarding the operation's result.
	Message string `json:"message,omitempty"`
	// Error 是操作过程中发生的错误（如果有）。
	// Error is any error that occurred during the operation.
	Error error `json:"-"` // Should not be marshalled to JSON directly
	// Data 包含操作返回的额外数据。
	// Data contains additional data returned by the operation.
	Data interface{} `json:"data,omitempty"`
}