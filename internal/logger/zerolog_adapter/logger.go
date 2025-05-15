// Package zerolog_adapter 提供了基于zerolog库的Logger接口实现。
// Package zerolog_adapter provides an implementation of the Logger interface using the zerolog library.
package zerolog_adapter

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/turtacn/mha4rdb/internal/logger/iface"
)

// Compile-time check to ensure zerologLogger implements iface.Logger
var _ iface.Logger = (*zerologLogger)(nil)

// zerologLogger 是 iface.Logger 的zerolog实现。
// zerologLogger is the zerolog implementation of iface.Logger.
type zerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger 创建一个新的zerologLogger实例。
// NewZerologLogger creates a new zerologLogger instance.
// levelStr: "debug", "info", "warn", "error", "fatal", "panic"
// isConsole: true for human-readable console output, false for JSON output
func NewZerologLogger(levelStr string, isConsole bool) iface.Logger {
	var level zerolog.Level
	switch levelStr {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel
	default:
		level = zerolog.InfoLevel
	}

	var zl zerolog.Logger
	if isConsole {
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(level).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		zl = zerolog.New(os.Stderr).
			Level(level).
			With().
			Timestamp().
			Caller(). // Consider removing caller for performance in production JSON logs
			Logger()
	}

	return &zerologLogger{logger: zl}
}

// Debugf 实现 iface.Logger 接口。
// Debugf implements the iface.Logger interface.
func (l *zerologLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Infof 实现 iface.Logger 接口。
// Infof implements the iface.Logger interface.
func (l *zerologLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warnf 实现 iface.Logger 接口。
// Warnf implements the iface.Logger interface.
func (l *zerologLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Errorf 实现 iface.Logger 接口。
// Errorf implements the iface.Logger interface.
func (l *zerologLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatalf 实现 iface.Logger 接口。
// Fatalf implements the iface.Logger interface.
func (l *zerologLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// Panicf 实现 iface.Logger 接口。
// Panicf implements the iface.Logger interface.
func (l *zerologLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}

// Debug 实现 iface.Logger 接口。
// Debug implements the iface.Logger interface.
func (l *zerologLogger) Debug(args ...interface{}) {
	l.logger.Debug().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		// Fallback or handle multiple args formatting if needed
		return zerolog.Arr().Interface(args).ToString()
	})
}

// Info 实现 iface.Logger 接口。
// Info implements the iface.Logger interface.
func (l *zerologLogger) Info(args ...interface{}) {
	l.logger.Info().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		return zerolog.Arr().Interface(args).ToString()
	})
}

// Warn 实现 iface.Logger 接口。
// Warn implements the iface.Logger interface.
func (l *zerologLogger) Warn(args ...interface{}) {
	l.logger.Warn().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		return zerolog.Arr().Interface(args).ToString()
	})
}

// Error 实现 iface.Logger 接口。
// Error implements the iface.Logger interface.
func (l *zerologLogger) Error(args ...interface{}) {
	l.logger.Error().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		return zerolog.Arr().Interface(args).ToString()
	})
}

// Fatal 实现 iface.Logger 接口。
// Fatal implements the iface.Logger interface.
func (l *zerologLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		return zerolog.Arr().Interface(args).ToString()
	})
}

// Panic 实现 iface.Logger 接口。
// Panic implements the iface.Logger interface.
func (l *zerologLogger) Panic(args ...interface{}) {
	l.logger.Panic().MsgFunc(func() string {
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				return str
			}
		}
		return zerolog.Arr().Interface(args).ToString()
	})
}

// WithFields 实现 iface.Logger 接口。
// WithFields implements the iface.Logger interface.
func (l *zerologLogger) WithFields(fields iface.Fields) iface.Logger {
	return &zerologLogger{logger: l.logger.With().Fields(fields).Logger()}
}

// WithError 实现 iface.Logger 接口。
// WithError implements the iface.Logger interface.
func (l *zerologLogger) WithError(err error) iface.Logger {
	return &zerologLogger{logger: l.logger.With().Err(err).Logger()}
}