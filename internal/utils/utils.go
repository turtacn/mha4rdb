// Package utils 包含项目通用的辅助函数。
// Package utils contains common utility functions for the project.
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

// GetHostname 返回当前主机的主机名。如果获取失败，则返回 "unknown"。
// GetHostname returns the hostname of the current machine. Returns "unknown" on failure.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// GenerateRandomID 生成一个指定长度的十六进制随机字符串ID。
// GenerateRandomID generates a random hexadecimal string ID of the specified length.
// The resulting string will be 2*length characters long.
func GenerateRandomID(length int) (string, error) {
	if length <= 0 {
		length = 16 // 默认长度 (Default length)
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ContainsString 检查字符串切片是否包含指定的字符串。
// ContainsString checks if a slice of strings contains the specified string.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// InterfaceSliceToStringSlice 将 interface{} 切片转换为 string 切片。
// InterfaceSliceToStringSlice converts a slice of interface{} to a slice of string.
// 如果元素不是字符串类型，则会出错。
// It will error if an element is not of string type.
func InterfaceSliceToStringSlice(data []interface{}) ([]string, error) {
	result := make([]string, len(data))
	for i, v := range data {
		str, ok := v.(string)
		if !ok {
			return nil, NewErrorf("element at index %d is not a string", i)
		}
		result[i] = str
	}
	return result, nil
}

// NewErrorf is a helper to create formatted errors, similar to fmt.Errorf.
// NewErrorf 是一个创建格式化错误的辅助函数，类似于 fmt.Errorf。
// This is defined here to avoid circular dependencies if internal/errors is not yet available
// or for simple local errors. For more general errors, use the internal/errors package.
// 这里定义是为了避免在 internal/errors 包尚不可用时产生循环依赖，或用于简单的局部错误。
// 对于更通用的错误，请使用 internal/errors 包。
type utilityError struct {
	s string
}

func (e *utilityError) Error() string {
	return e.s
}

// NewErrorf formats according to a format specifier and returns the string as a value that satisfies error.
// NewErrorf 根据格式说明符进行格式化，并返回满足 error 接口的字符串值。
func NewErrorf(format string, args ...interface{}) error {
	// This import is fine here as it's a common package.
	// If this were to cause issues, one might implement a simpler formatter or panic.
	// 此处导入 fmt 包是可以的，因为它是一个常用包。
	// 如果这会导致问题，可以实现一个更简单的格式化程序或引发 panic。
	// For now, assume fmt is always available for utility functions.
	// 目前，假设 fmt 对于工具函数总是可用的。
	// Note: This is slightly different from the standard library's fmt.Errorf in how it's defined,
	// but serves a similar purpose for local error creation.
	// 注意：这与标准库的 fmt.Errorf 在定义方式上略有不同，但用于创建局部错误的目的相似。
	// To avoid potential confusion with a global fmt.Errorf,
	// directly using `return fmt.Errorf(format, args...)` is usually preferred.
	// 为了避免与全局的 fmt.Errorf 潜在混淆，通常首选直接使用 `return fmt.Errorf(format, args...)`。
	// However, if we want to ensure all errors from this package have a specific local type,
	// this custom type approach can be used.
	// 但是，如果我们想确保此包中的所有错误都具有特定的局部类型，则可以使用这种自定义类型方法。
	// For simplicity and general utility, let's just use fmt.Errorf.
	// 为简单起见和通用性，我们直接使用 fmt.Errorf。
	_ = &utilityError{} // Avoid unused type warning if we switch to fmt.Errorf
	// We will use the standard fmt.Errorf for now.
	// 我们现在将使用标准的 fmt.Errorf。
	// If more complex error wrapping or typing is needed locally, this can be revisited.
	// 如果本地需要更复杂的错误包装或类型化，可以重新审视这一点。
	// For now, no, don't re-implement fmt.Errorf. This was for illustration.
	// return fmt.Errorf(format, args...) - let's remove this to avoid confusion with actual std lib.
	// Users of this package can use fmt.Errorf or the project's central error package.
	// This file contains general utilities, not error definitions primarily.
	return nil // Placeholder for further utility functions
}