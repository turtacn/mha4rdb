// Package vastbase 包含了与Vastbase数据库交互的具体实现。
// Package vastbase contains the specific implementation for interacting with a Vastbase database.
package vastbase

import (
	"fmt"
	"time"

	"github.com/turtacn/mha4rdb/internal/errors"
)

const (
	// DefaultPort 是Vastbase的默认端口。
	// DefaultPort is the default port for Vastbase.
	DefaultPort = 5432
	// DefaultConnectTimeout 是默认连接超时时间。
	// DefaultConnectTimeout is the default connection timeout.
	DefaultConnectTimeout = 5 * time.Second
	// DefaultSSLMode 是默认的SSL模式。
	// DefaultSSLMode is the default SSL mode.
	DefaultSSLMode = "disable" // "allow", "prefer", "require", "verify-ca", "verify-full"
)

// Config 结构体定义了连接Vastbase数据库所需的配置参数。
// Config struct defines the configuration parameters required to connect to a Vastbase database.
type Config struct {
	// Host 是数据库服务器的主机名或IP地址。
	// Host is the hostname or IP address of the database server.
	Host string `yaml:"host" json:"host"`
	// Port 是数据库服务器的端口号。
	// Port is the port number of the database server.
	Port int `yaml:"port" json:"port"`
	// User 是连接数据库的用户名。
	// User is the username for connecting to the database.
	User string `yaml:"user" json:"user"`
	// Password 是连接数据库的密码。
	// Password is the password for connecting to the database.
	Password string `yaml:"password" json:"password"`
	// DBName 是要连接的数据库名称。
	// DBName is the name of the database to connect to.
	DBName string `yaml:"dbname" json:"dbname"`
	// SSLMode 指定SSL连接模式。
	// SSLMode specifies the SSL connection mode.
	// 可选值: "disable", "allow", "prefer", "require", "verify-ca", "verify-full".
	// Possible values: "disable", "allow", "prefer", "require", "verify-ca", "verify-full".
	SSLMode string `yaml:"sslmode" json:"sslmode"`
	// SSLCertPath 是客户端SSL证书路径。
	// SSLCertPath is the path to the client SSL certificate.
	SSLCertPath string `yaml:"ssl_cert_path,omitempty" json:"ssl_cert_path,omitempty"`
	// SSLKeyPath 是客户端SSL私钥路径。
	// SSLKeyPath is the path to the client SSL private key.
	SSLKeyPath string `yaml:"ssl_key_path,omitempty" json:"ssl_key_path,omitempty"`
	// SSLRootCertPath 是SSL根证书路径。
	// SSLRootCertPath is the path to the SSL root certificate.
	SSLRootCertPath string `yaml:"ssl_root_cert_path,omitempty" json:"ssl_root_cert_path,omitempty"`
	// ConnectTimeout 是连接超时时间。
	// ConnectTimeout is the connection timeout duration.
	ConnectTimeout time.Duration `yaml:"connect_timeout" json:"connect_timeout"`
	// MaxOpenConns 是最大打开连接数。
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `yaml:"max_open_conns,omitempty" json:"max_open_conns,omitempty"`
	// MaxIdleConns 是最大空闲连接数。
	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	MaxIdleConns int `yaml:"max_idle_conns,omitempty" json:"max_idle_conns,omitempty"`
	// ConnMaxLifetime 是连接可被重用的最大时间量。
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime,omitempty" json:"conn_max_lifetime,omitempty"`
	// ConnMaxIdleTime 是连接在被关闭之前可以空闲的最大时间量。
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle before being closed.
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time,omitempty" json:"conn_max_idle_time,omitempty"`
	// ApplicationName 是在数据库中标识此连接的名称。
	// ApplicationName is the name used to identify this connection in the database.
	ApplicationName string `yaml:"application_name,omitempty" json:"application_name,omitempty"`
}

// Validate 检查配置是否有效。
// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}
	if c.Port <= 0 {
		c.Port = DefaultPort
	}
	if c.User == "" {
		return errors.New("database user is required")
	}
	if c.DBName == "" {
		return errors.New("database dbname is required")
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = DefaultConnectTimeout
	}
	if c.SSLMode == "" {
		c.SSLMode = DefaultSSLMode
	}
	// Add more validation for SSL paths if SSL mode requires them
	return nil
}

// DSN 生成Vastbase (PostgreSQL兼容) 的数据源名称 (Data Source Name)。
// DSN generates the Data Source Name for Vastbase (PostgreSQL compatible).
func (c *Config) DSN() string {
	// Example DSN: "postgres://user:password@host:port/dbname?sslmode=disable&connect_timeout=5"
	// Or "host=myhost port=myport user=gorm password=gorm dbname=gorm sslmode=disable"
	// We will use the latter format for lib/pq
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Host, c.Port, c.User, c.DBName, c.SSLMode, int(c.ConnectTimeout.Seconds()))
	if c.Password != "" {
		dsn += fmt.Sprintf(" password=%s", c.Password)
	}
	if c.SSLCertPath != "" {
		dsn += fmt.Sprintf(" sslcert=%s", c.SSLCertPath)
	}
	if c.SSLKeyPath != "" {
		dsn += fmt.Sprintf(" sslkey=%s", c.SSLKeyPath)
	}
	if c.SSLRootCertPath != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", c.SSLRootCertPath)
	}
	if c.ApplicationName != "" {
		dsn += fmt.Sprintf(" application_name=%s", c.ApplicationName)
	}
	return dsn
}

// GetDefaultConfig 返回一个包含默认值的Vastbase配置。
// GetDefaultConfig returns a Vastbase configuration with default values.
func GetDefaultConfig() Config {
	return Config{
		Port:            DefaultPort,
		SSLMode:         DefaultSSLMode,
		ConnectTimeout:  DefaultConnectTimeout,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		ApplicationName: "mha4rdb",
	}
}