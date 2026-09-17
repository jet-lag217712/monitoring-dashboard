# System Design

Canonical architecture is documented in
[`.ai/project-context/architecture.md`](../../.ai/project-context/architecture.md).

```mermaid
flowchart LR
    subgraph appliance["Local Equate Appliance"]
        devices["SNMPv2c devices"] -->|"independent polls"| collector["SNMP Collector"]
        inventory["Static + managed inventory"] --> collector
        tui["Local Bubble Tea TUI"] <-->|"Unix socket"| collector
        collector --> outbox["SQLite outbox"]
        collector --> local["metrics · healthz · readyz"]
        outbox -->|"local v2 telemetry"| broker["MQTT/TLS QoS 1"]
        broker --> ingestion["Ingestion"]
        ingestion --> postgres[("PostgreSQL\nsystem of record")]
        postgres --> api["Backend API"]
        api --> dashboard["nginx → React dashboard"]
    end
```

The collector sends device, interface, health, and heartbeat telemetry. Dependency edges are local reachability evidence; they are not a complete physical-network graph. The TUI and status/control surface remain local-only.
