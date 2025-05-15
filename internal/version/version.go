// Package version 提供了mha4rdb应用的版本信息。
// Package version provides version information for the mha4rdb application.
package version

import "fmt"

var (
	// Version 是应用程序的语义化版本号。
	// Version is the semantic version number of the application.
	Version = "0.0.1-dev"

	// GitCommit 是构建此应用的Git提交哈希。
	// GitCommit is the Git commit hash from which this application was built.
	GitCommit = "unknown"

	// BuildDate 是应用的构建日期。
	// BuildDate is the build date of the application.
	BuildDate = "unknown"

	// GoVersion 是用于构建此应用的Go版本。
	// GoVersion is the Go version used to build this application.
	GoVersion = "unknown" // This can be set by goreleaser or build scripts
)

// Info 结构体包含所有版本相关信息。
// Info struct contains all version-related information.
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// GetInfo 返回包含所有版本信息的Info结构体。
// GetInfo returns an Info struct containing all version information.
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	}
}

// String 返回格式化的版本信息字符串。
// String returns a formatted version information string.
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s",
		i.Version, i.GitCommit, i.BuildDate, i.GoVersion)
}

// PrintVersion 打印版本信息到标准输出。
// PrintVersion prints version information to standard output.
func PrintVersion() {
	fmt.Println(GetInfo().String())
}