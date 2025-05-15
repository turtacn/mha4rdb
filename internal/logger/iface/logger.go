// Package iface 定义了日志记录器的接口。
// Package iface defines the interface for a logger.
package iface

// Fields 类型用于传递结构化的日志字段。
// Fields type is used for passing structured logging fields.
type Fields map[string]interface{}

// Logger 是mha4rdb项目统一的日志接口。
// Logger is the unified logging interface for the mha4rdb project.
type Logger interface {
	// Debugf 记录Debug级别的格式化日志。
	// Debugf logs a formatted debug level log.
	Debugf(format string, args ...interface{})
	// Infof 记录Info级别的格式化日志。
	// Infof logs a formatted info level log.
	Infof(format string, args ...interface{})
	// Warnf 记录Warn级别的格式化日志。
	// Warnf logs a formatted warn level log.
	Warnf(format string, args ...interface{})
	// Errorf 记录Error级别的格式化日志。
	// Errorf logs a formatted error level log.
	Errorf(format string, args ...interface{})
	// Fatalf 记录Fatal级别的格式化日志，并随后调用os.Exit(1)。
	// Fatalf logs a formatted fatal level log and then calls os.Exit(1).
	Fatalf(format string, args ...interface{})
	// Panicf 记录Panic级别的格式化日志，并随后调用panic()。
	// Panicf logs a formatted panic level log and then calls panic().
	Panicf(format string, args ...interface{})

	// Debug 记录Debug级别的日志消息。
	// Debug logs a debug level log message.
	Debug(args ...interface{})
	// Info 记录Info级别的日志消息。
	// Info logs an info level log message.
	Info(args ...interface{})
	// Warn 记录Warn级别的日志消息。
	// Warn logs a warn level log message.
	Warn(args ...interface{})
	// Error 记录Error级别的日志消息。
	// Error logs an error level log message.
	Error(args ...interface{})
	// Fatal 记录Fatal级别的日志消息，并随后调用os.Exit(1)。
	// Fatal logs a fatal level log message and then calls os.Exit(1).
	Fatal(args ...interface{})
	// Panic 记录Panic级别的日志消息，并随后调用panic()。
	// Panic logs a panic level log message and then calls panic().
	Panic(args ...interface{})

	// WithFields 返回一个带有附加结构化字段的新日志记录器实例。
	// WithFields returns a new logger instance with additional structured fields.
	WithFields(fields Fields) Logger
	// WithError 返回一个带有错误信息字段的新日志记录器实例。
	// WithError returns a new logger instance with an error field.
	WithError(err error) Logger
}