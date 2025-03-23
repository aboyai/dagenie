# DAGenie DB

## Purpose-Built for DAGs. Engineered for Speed.

DAGenie is a high-performance, distributed database purpose-built to manage Directed Acyclic Graphs (DAGs). With disk persistence via BadgerDB, clustered storage with sharding, replication, and a powerful query language (DQL), DAGenie empowers developers and system architects to build, manage, and introspect complex DAG workflows at scale.

---

## Features Roadmap

- ✨ **DAG-First Design**: Optimized for storing and querying DAG structures.
- 🔢 **Disk Persistence**: Leveraging BadgerDB for high-speed, embedded storage.
- 🛡️ **Raft Consensus**: Strong consistency, replication, and leader election.
- 📊 **Sharded Storage**: Scalable horizontal partitioning via consistent hashing.
- 📊 **Replication**: Configurable replication factor with quorum writes/reads.
- ⚙️ **gRPC & TCP Interfaces**: Dual interfaces for custom app integration.
- 🔍 **DQL Query Language**: Custom query language with WHERE, UPDATE, DELETE, COMMIT, ROLLBACK.
- 🔎 **Full-Text Indexing**: Fast indexed DAG and task lookup.
- 🔄 **Gossip Protocol**: Cluster discovery, failure detection, auto-rebalancing.
- 🤑 **User Privileges**: Admin/user roles with scoped permissions.

---

## Getting Started

### 🔧 Installation

```bash
git clone https://github.com/aboyai/dagenie.git
cd dagenie
./build.bat
```

### 🔍 Running CLI

```bash
dagenine serve [port]
dagenie connect --host localhost --port [port]
```

### 🔍 Execute DQL via CLI

```sql
INSERT INTO dag (id, name, status, payload, dependencies, dagid, duration, retries) VALUES ('1', 'AWS', 'pending', '{}', '[]', 'abc234', 200, 10);
```

```sql
SELECT * FROM dag;
```

```sql
SELECT name, SUM(duration) FROM dag GROUP BY name ORDER BY SUM(duration) DESC LIMIT 5;
```

## Language Clients

- [Go Client](./clients/go/README.md)
- [Java Client](./clients/java/README.md)
- [Python Client](./clients/python/README.md)
- [Node.js Client](./clients/nodejs/README.md)
- [C++ Client](./clients/cpp/README.md)
- [Rust Client](./clients/rust/README.md)
- [.NET Client](./clients/dotnet/README.md)
- [Ruby Client](./clients/ruby/README.md)

---

## Contributing

We welcome contributions! Please submit issues and pull requests.

- [Contribution Guide](./CONTRIBUTING.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)

---

## License

Apache 2.0 License © 2025 Aboyai LLC

---

## Contact & Support

For enterprise support, integration help, or queries:

- Email: [support@aboyai.com](mailto:support@aboyai.com)
- Website: [https://aboyai.com](https://aboyayi.com)

---

*DAGenie — Purpose-Built for DAGs. Engineered for Speed.*

