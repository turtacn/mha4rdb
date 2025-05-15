// Package signal provides utilities for handling OS signals, typically for graceful shutdown.
// Package signal 提供了处理操作系统信号的实用程序，通常用于优雅关闭。
package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/turtacn/mha4rdb/internal/logger/iface"
)

// SetupSignalHandler 注册一个信号处理器，监听SIGINT和SIGTERM信号。
// 当接收到这些信号时，会调用提供的回调函数。
// SetupSignalHandler registers a signal handler for SIGINT and SIGTERM.
// It calls the provided callback function when these signals are received.
// It returns a channel that can be used to wait for the signal to be processed by the callback.
//
// Deprecated: Use SetupSignalChannel which is more idiomatic for Go.
// 废弃: 请使用 SetupSignalChannel，它更符合 Go 的习惯用法。
func SetupSignalHandler(log iface.Logger, shutdownCallback func()) chan struct{} {
	stop := make(chan struct{})
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Infof("Received signal: %s. Shutting down...", sig.String())
		if shutdownCallback != nil {
			shutdownCallback()
		}
		log.Infof("Shutdown callback completed.")
		close(stop) // Signal that shutdown process has been triggered
	}()
	return stop
}

// SetupSignalChannel 注册一个信号处理器，监听SIGINT和SIGTERM信号。
// 返回一个channel，当接收到退出信号时，该channel会接收到信号。
// SetupSignalChannel registers a signal handler for SIGINT and SIGTERM.
// It returns a channel that will receive the signal when an exit signal is caught.
func SetupSignalChannel() <-chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	return sigCh
}

// HandleSignals 监听指定的信号channel，并在接收到信号时调用回调函数。
// 它会阻塞直到信号被处理或上下文被取消。
// HandleSignals listens on the provided signal channel and calls the callback when a signal is received.
// This function is blocking until a signal is received and callback is executed, or the context is cancelled.
// ctx can be used to unblock the function if necessary.
//
// Example usage:
// sigCh := SetupSignalChannel()
// HandleSignals(context.Background(), logger, sigCh, func(s os.Signal) {
//    logger.Infof("Shutdown initiated by signal: %s", s.String())
//    // Perform cleanup
// })
func HandleSignals(
	log iface.Logger,
	sigCh <-chan os.Signal,
	shutdownCallback func(os.Signal),
) {
	select {
	case sig := <-sigCh:
		log.Infof("Received signal: %v. Initiating shutdown sequence.", sig)
		if shutdownCallback != nil {
			shutdownCallback(sig)
		}
		log.Infof("Shutdown sequence completed for signal: %v.", sig)
	}
}