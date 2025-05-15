# mha4rdb: Relational Database Master High Availability 

`mha4rdb` is kind of Master High Availability for Relational Databases(e.g. Vastbase、Mysql、Postgresql、...)

Detailed design and architecture could be found at [架构设计 / Architecture Design Document](docs/architecture.md)

## 1. Project Introduction

`mha4rdb` is a Master High Availability (MHA) solution for Vastbase databases specifically designed for Hyper-Converged Infrastructure (HCI) scenarios. It leverages the Raft consistency protocol, either built into or integrated with the Vastbase database, to ensure strong consistency and durability of configuration data. Through a lightweight Agent and Virtual IP (VIP) mechanism, it provides automated fault detection and failover for the database service, offering a stable, reliable, and highly available configuration data backend for HCI infrastructure software.

Compared to general OLTP databases, HCI scenarios place a particularly high emphasis on data consistency, while the requirement for extreme read/write performance is relatively secondary. `mha4rdb` focuses on meeting these core needs of the HCI environment, providing a simple, easy-to-use, highly automated database HA capability that integrates well with the HCI platform.

## 2. Key Features

* **Raft-based Strong Consistency**: Relies on the Raft protocol implemented by the Vastbase database itself to ensure data is not lost and remains consistent during cluster state changes and node failures.
* **Virtual IP (VIP) Management**: By automatically managing the floating of the Virtual IP, it provides a fixed database access point for HCI services, eliminating the need to be aware of backend Leader node switches.
* **Lightweight MHA Agent**: An Agent deployed on each database node is responsible for monitoring the status of the local Vastbase instance, participating in cluster state synchronization, and executing VIP management and decision-making.
* **Automated Failover**: Automatically detects Leader node failures and, coordinated by a majority of Agents, automatically switches the VIP to the new Leader node.
* **Multiple Deployment Modes**: Supports single-node, dual-node + arbiter, and multi-node cluster deployment modes to meet different scale and reliability requirements.
* **Easy Integration**: Provides an MHA Client library for easy integration by HCI services to query database cluster status.

## 3. Architecture Overview

The core components of the `mha4rdb` system include:

* **Vastbase Database Instances**: Form the Raft cluster, responsible for data storage and replication.
* **MHA Agents**: Deployed on each Vastbase node, monitoring database status, participating in cluster management, and managing the VIP.
* **Optional Arbiter Service**: Provides an additional vote in dual-node mode to prevent split-brain scenarios.
* **Virtual IP**: Provides a unified access point.
* **HCI Service**: The control plane of the HCI infrastructure software, accessing the database via the VIP and interacting with the Agents via the MHA Client library to obtain status.

For a more detailed architecture design, please refer to the [Architecture Design Document](https://www.google.com/search?q=docs/architecture.md).

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
git clone [https://github.com/turtacn/mha4rdb.git](https://github.com/turtacn/mha4rdb.git)
cd mha4rdb

# Build the MHA Agent
go build -o bin/mha-agent cmd/agent/main.go

# Build the MHA Client library (to be used as a Go module dependency, no separate executable needed)
# go get [github.com/turtacn/mha4rdb/pkg/client](https://github.com/turtacn/mha4rdb/pkg/client)
````

### 5.3 Configuration

Agent configuration is primarily done via a YAML file. Refer to the `configs/agent.yaml` example.

### 5.4 Running

Run the MHA Agent on each Vastbase database node:

```bash
./bin/mha-agent --config configs/agent.yaml
```

For specific deployment and running steps, as well as configuration details for different modes, please refer to the detailed [Deployment Guide](https://www.google.com/search?q=docs/deployment_guide.md) (To Be Completed).

## 6\. Integration with HCI

HCI infrastructure software can integrate with MHA Agents by including the `pkg/client` package. The MHA Client provides interfaces to query cluster status, get the Leader address, etc. HCI services can use this information for database connection management and fault handling.

HCI services are responsible for connecting to the VIP for normal database read/write operations. Upon detecting a Leader switch event (obtained via the MHA Client), they should update their internal connection state or perform other necessary coordination logic.

## 7\. Contributing

We welcome contributions from the community\! If you are interested in the `mha4rdb` project and would like to report bugs, suggest improvements, or submit code, please refer to the [Contributing Guide](https://www.google.com/search?q=CONTRIBUTING.md) (To Be Completed).

## 8\. License

This project is licensed under the [Apache 2.0 License](https://www.google.com/search?q=LICENSE).

## 9\. Project Tree (WIP)

```text
mha4rdb/
├── api/
│   └── proto/
│       └── v1/
│           ├── mha.proto
│           ├── mha\_grpc.pb.go
│           └── mha.pb.go
├── cmd/
│   └── mha4rdb-agent/
│       └── main.go
├── configs/
│   └── agent.example.yaml
├── docs/
│   └── architecture.md
├── internal/
│   ├── agent/
│   │   ├── agent.go
│   │   ├── config.go
│   │   ├── election/
│   │   │   └── manager.go
│   │   ├── health/
│   │   │   └── checker.go
│   │   ├── monitor/
│   │   │   └── monitor.go
│   │   ├── raft/
│   │   │   ├── raft\_ FSM.go
│   │   │   ├── raft\_node.go
│   │   │   └── transport.go
│   │   ├── rpc/
│   │   │   └── server.go
│   │   └── service/  // Implements gRPC services defined in api/proto
│   │       └── mha\_service.go
│   ├── client/
│   │   ├── client.go
│   │   └── options.go
│   ├── core/
│   │   ├── cluster/
│   │   │   └── state.go
│   │   └── types/
│   │       ├── common.go
│   │       └── enum.go
│   ├── database/
│   │   ├── iface/
│   │   │   └── db.go
│   │   └── vastbase/
│   │       ├── vastbase.go
│   │       └── config.go
│   ├── errors/
│   │   └── errors.go
│   ├── logger/
│   │   ├── iface/
│   │   │   └── logger.go
│   │   └── zerolog\_adapter/ // Example implementation
│   │       └── logger.go
│   ├── network/
│   │   ├── vip/
│   │   │   └── manager.go
│   │   └── iface/
│   │       └── vip.go
│   ├── utils/
│   │   └── utils.go
│   └── version/
│       └── version.go
├── pkg/ // Potentially for public libraries if any part of client becomes one
│   └── signal/
│       └── signal.go
├── scripts/
│   ├── run\_tests.sh
│   └── build.sh
├── test/
│   ├── integration/
│   └── e2e/
├── .gitignore
├── go.mod
├── go.sum
└── README.md

```

## 10\. References
  * [MySQL Group Replication](https://dev.mysql.com/doc/refman/8.0/en/group-replication.html)
  * [PostgreSQL Patroni](https://patroni.readthedocs.io/en/latest/)
  * [MySQL MHA (Master High Availability)](https://github.com/yoshinorim/mha4mysql)
