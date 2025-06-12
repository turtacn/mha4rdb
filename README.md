# mha4rdb: Relational Database Master High Availability 

`mha4rdb` is kind of Master High Availability for Relational Databases(e.g. Vastbase、Mysql、Postgresql、...)

* Here RDB is Vastbase for example
* Here common service is IaaS, to take [5]基础设施 for example

## 1. Project Introduction

`mha4rdb` is a Master High Availability (MHA) solution for Vastbase databases specifically designed for Hyper-Converged Infrastructure (HCI[5][5]基础设施) scenarios. It leverages the Raft consistency protocol, either built into or integrated with the Vastbase database, to ensure strong consistency and durability of configuration data. Through a lightweight Agent and Virtual IP (VIP) mechanism, it provides automated fault detection and failover for the database service, offering a stable, reliable, and highly available configuration data backend for [5]基础设施 infrastructure software.

Compared to general OLTP databases, [5]基础设施 scenarios place a particularly high emphasis on data consistency, while the requirement for extreme read/write performance is relatively secondary. `mha4rdb` focuses on meeting these core needs of the [5]基础设施 environment, providing a simple, easy-to-use, highly automated database HA capability that integrates well with the [5]基础设施 platform.

## 2. Key Features

* **Raft-based Strong Consistency**: Relies on the Raft protocol implemented by the Vastbase database itself to ensure data is not lost and remains consistent during cluster state changes and node failures.
* **Virtual IP (VIP) Management**: By automatically managing the floating of the Virtual IP, it provides a fixed database access point for [5]基础设施 services, eliminating the need to be aware of backend Leader node switches.
* **Lightweight MHA Agent**: An Agent deployed on each database node is responsible for monitoring the status of the local Vastbase instance, participating in cluster state synchronization, and executing VIP management and decision-making.
* **Automated Failover**: Automatically detects Leader node failures and, coordinated by a majority of Agents, automatically switches the VIP to the new Leader node.
* **Multiple Deployment Modes**: Supports single-node, dual-node + arbiter, and multi-node cluster deployment modes to meet different scale and reliability requirements.
* **Easy Integration**: Provides an MHA Client library for easy integration by [5]基础设施 services to query database cluster status.

## 3. Architecture Overview

The core components of the `mha4rdb` system include:

* **Vastbase Database Instances**: Form the Raft cluster, responsible for data storage and replication.
* **MHA Agents**: Deployed on each Vastbase node, monitoring database status, participating in cluster management, and managing the VIP.
* **Optional Arbiter Service**: Provides an additional vote in dual-node mode to prevent split-brain scenarios.
* **Virtual IP**: Provides a unified access point.
* **[5]基础设施 Service**: The control plane of the [5]基础设施 infrastructure software, accessing the database via the VIP and interacting with the Agents via the MHA Client library to obtain status.

For a detailed architecture design, please refer to the [架构设计 / Architecture Design Document](docs/architecture.md).

For a more detailed architecture and more designs, please refer to the [详细设计](docs/design.md).


## 4. Deployment Modes

`mha4rdb` supports the following deployment modes:

* **Single Node**: Suitable for development or test environments, provides basic health check capabilities.
* **Dual Node + Arbiter**: Introduces an independent arbiter node alongside two Vastbase instances to form a minimum majority, providing basic high availability.
* **Multi Node Cluster**: Three or more Vastbase instances form a Raft cluster, providing higher availability and fault tolerance.

## 5. Getting Started

### 5.1 Dependencies

* Go language environment (version >= 1.20.2)
* Vastbase database (needs to support or integrate the Raft protocol and expose MHA-related status query and management interfaces)
* (Optional) Operating system tools or libraries for managing VIP (e.g., `ip` command)
* (Optional) Arbiter service (for dual-node + arbiter mode)

### 5.2 Building

```bash
# Clone the repository
git clone https://github.com/turtacn/mha4rdb.git
cd mha4rdb

# Build the MHA Agent
go build -o bin/mha-agent cmd/agent/main.go

# Build the MHA Client library (to be used as a Go module dependency, no separate executable needed)
# go get https://github.com/turtacn/mha4rdb/pkg/client
````

### 5.3 Configuration

Agent configuration is primarily done via a YAML file. Refer to the `configs/agent.yaml` example.

### 5.4 Running

Run the MHA Agent on each Vastbase database node:

```bash
./bin/mha-agent --config configs/agent.yaml
```

For specific deployment and running steps, as well as configuration details for different modes, please refer to the detailed [Deployment Guide](https://www.google.com/search?q=docs/deployment_guide.md) (To Be Completed).

## 6\. Integration with [5]基础设施

[5]基础设施 infrastructure software can integrate with MHA Agents by including the `pkg/client` package. The MHA Client provides interfaces to query cluster status, get the Leader address, etc. [5]基础设施 services can use this information for database connection management and fault handling.

[5]基础设施 services are responsible for connecting to the VIP for normal database read/write operations. Upon detecting a Leader switch event (obtained via the MHA Client), they should update their internal connection state or perform other necessary coordination logic.

## 7\. Contributing

We welcome contributions from the community\! If you are interested in the `mha4rdb` project and would like to report bugs, suggest improvements, or submit code, please refer to the [Contributing Guide](https://www.google.com/search?q=CONTRIBUTING.md) (To Be Completed).

## 8\. License

This project is licensed under the [Apache 2.0 License](https://www.google.com/search?q=LICENSE).

## 9\. Project Tree (WIP)

```text
mha4rdb/
├── cmd/                                # 命令行入口
│   ├── mha-agent/                      # MHA Agent命令行
│   ├── mha-manager/                    # MHA Manager命令行
│   └── mha-cli/                        # MHA CLI工具命令行
├── docs/                               # 文档
│   ├── architecture.md                 # 架构文档
│   ├── deployment.md                   # 部署文档
│   └── api.md                          # API文档
├── examples/                           # 示例
│   ├── client/                         # 客户端示例
│   ├── config/                         # 配置示例
│   └── scripts/                        # 脚本示例
├── internal/                           # 内部包
│   ├── common/                         # 通用代码
│   │   ├── constants/                  # 常量定义
│   │   ├── errors/                     # 错误定义
│   │   ├── logging/                    # 日志
│   │   ├── metrics/                    # 指标
│   │   └── utils/                      # 工具函数
│   ├── agent/                          # Agent实现
│   │   ├── server/                     # Agent服务
│   │   ├── monitor/                    # 监控实现
│   │   ├── vip/                        # VIP管理
│   │   └── handlers/                   # 请求处理
│   ├── manager/                        # Manager实现
│   │   ├── server/                     # Manager服务
│   │   ├── cluster/                    # 集群管理
│   │   ├── election/                   # 选举管理
│   │   └── handlers/                   # 请求处理
│   └── db/                             # 数据库适配
│       ├── mysql/                      # MySQL适配
│       ├── postgresql/                 # PostgreSQL适配
│       └── vastbase/                   # Vastbase适配
├── pkg/                                # 公共包
│   ├── api/                            # API定义
│   │   ├── client/                     # 客户端API
│   │   ├── agent/                      # Agent API
│   │   └── manager/                    # Manager API
│   ├── config/                         # 配置
│   ├── client/                         # 客户端实现
│   ├── raft/                           # Raft实现
│   ├── db/                             # 数据库接口
│   └── models/                         # 数据模型
├── test/                               # 测试
│   ├── integration/                    # 集成测试
│   ├── benchmark/                      # 基准测试
│   └── utils/                          # 测试工具
├── scripts/                            # 脚本工具
│   ├── install/                        # 安装脚本
│   ├── build/                          # 构建脚本
│   └── ci/                             # CI脚本
├── Makefile                            # 构建规则
├── go.mod                              # Go模块定义

```

## 10\. References
- [1] [MySQL Group Replication](https://dev.mysql.com/doc/refman/8.0/en/group-replication.html)
- [2] [PostgreSQL Patroni](https://patroni.readthedocs.io/en/latest/)
- [3] [Development tree of Master High Availability Manager and tools for MySQL (MHA), Manager part](https://github.com/yoshinorim/mha4mysql-manager)
- [4] [Development tree of Master High Availability Manager and tools for MySQL (MHA), Node (MySQL Server) part](https://github.com/yoshinorim/mha4mysql-node)
- [5] [full-stack hyperconverged infrastructure](https://www.gartner.com/reviews/market/full-stack-hyperconverged-infrastructure-software)
