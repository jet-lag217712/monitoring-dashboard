# V2 cutover

## Decision

Production and deployment profiles publish and consume **MQTT v2 only**.
No production workload depended on v1 routes, so the Phase 4 dual-publish
window is closed.

Deployment collector YAML sets:

```yaml
publisher:
  telemetry_version: v2
```

Empty `telemetry_version` in code also defaults to `v2`.

Deployment ingestion YAML subscribes to:

```yaml
topics:
  - "site/+/device/+/telemetry/v2/#"
  - "site/+/collector/+/telemetry/v2/heartbeat"
```

## Smoke

Use `make lab-smoke`.
