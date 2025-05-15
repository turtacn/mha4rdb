// Package errors 定义了项目中集中的错误类型和常量。
// Package errors defines centralized error types and constants for the project.
package errors

import "errors"

var (
	// ErrNotFound 表示请求的资源未找到。
	// ErrNotFound indicates that the requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists 表示尝试创建已存在的资源。
	// ErrAlreadyExists indicates an attempt to create a resource that already exists.
	ErrAlreadyExists = errors.New("already exists")

	// ErrInvalidInput 表示提供的输入无效。
	// ErrInvalidInput indicates that the provided input is invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrPermissionDenied 表示操作权限不足。
	// ErrPermissionDenied indicates insufficient permission for the operation.
	ErrPermissionDenied = errors.New("permission denied")

	// ErrTimeout 表示操作超时。
	// ErrTimeout indicates that the operation timed out.
	ErrTimeout = errors.New("operation timed out")

	// ErrNotImplemented 表示功能尚未实现。
	// ErrNotImplemented indicates that a feature or function is not yet implemented.
	ErrNotImplemented = errors.New("not implemented")

	// ErrClusterNotHealthy 表示集群不健康，无法执行操作。
	// ErrClusterNotHealthy indicates the cluster is not healthy and cannot perform the operation.
	ErrClusterNotHealthy = errors.New("cluster is not healthy")

	// ErrLeaderNotFound 表示未能找到集群的Leader。
	// ErrLeaderNotFound indicates that the leader of the cluster could not be found.
	ErrLeaderNotFound = errors.New("leader not found")

	// ErrNotLeader 表示当前节点不是Leader，无法执行Leader专属操作。
	// ErrNotLeader indicates the current node is not the leader and cannot perform leader-specific operations.
	ErrNotLeader = errors.New("current node is not the leader")

	// ErrRaftProposeFailed 表示向Raft集群提交提案失败。
	// ErrRaftProposeFailed indicates failure to propose an entry to the Raft cluster.
	ErrRaftProposeFailed = errors.New("failed to propose to raft")

	// ErrDBConnectionFailed 表示连接数据库失败。
	// ErrDBConnectionFailed indicates failure to connect to the database.
	ErrDBConnectionFailed = errors.New("database connection failed")

	// ErrDBQueryFailed 表示执行数据库查询失败。
	// ErrDBQueryFailed indicates failure to execute a database query.
	ErrDBQueryFailed = errors.New("database query failed")

	// ErrDBRoleAssertionFailed 表示数据库角色断言失败（例如，期望是Primary但实际是Secondary）。
	// ErrDBRoleAssertionFailed indicates a database role assertion failure (e.g., expected Primary but was Secondary).
	ErrDBRoleAssertionFailed = errors.New("database role assertion failed")

	// ErrVIPOperationFailed 表示虚拟IP操作失败。
	// ErrVIPOperationFailed indicates that a Virtual IP operation failed.
	ErrVIPOperationFailed = errors.New("vip operation failed")

	// ErrConfigurationInvalid 表示配置无效。
	// ErrConfigurationInvalid indicates that the configuration is invalid.
	ErrConfigurationInvalid = errors.New("invalid configuration")

	// ErrInternalServer 表示发生了未知的内部服务器错误。
	// ErrInternalServer indicates an unknown internal server error occurred.
	ErrInternalServer = errors.New("internal server error")
)

// Is convenience function to check if an error is of a specific type using errors.Is.
// Is 是一个方便的函数，用于使用 errors.Is 检查错误是否为特定类型。
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As convenience function to check if an error is of a specific type and retrieve it using errors.As.
// As 是一个方便的函数，用于使用 errors.As 检查错误是否为特定类型并检索它。
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// New creates a new error with the given message.
// New 使用给定的消息创建一个新错误。
func New(text string) error {
	return errors.New(text)
}

// Errorf formats according to a format specifier and returns the string as a value that satisfies error.
// Errorf 根据格式说明符进行格式化，并返回满足 error 接口的字符串值。
// Note: This is a simple wrapper around fmt.Errorf for consistency if needed,
// but often direct use of fmt.Errorf is fine.
// 注意：如果需要，这是对 fmt.Errorf 的简单封装以保持一致性，但通常直接使用 fmt.Errorf 也可以。
// For now, we'll just recommend using fmt.Errorf directly for formatted errors
// or wrapping existing errors.
// 目前，我们建议直接使用 fmt.Errorf 处理格式化错误或包装现有错误。

// TODO: Consider more structured errors if needed, e.g.:
// type MHAError struct {
//    Code int
//    Message string
//    Cause error
// }
// func (e *MHAError) Error() string { return fmt.Sprintf("code %d: %s (cause: %v)", e.Code, e.Message, e.Cause) }
// func (e *MHAError) Unwrap() error { return e.Cause }