# README Rewrite Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rewrite the README to follow standard Prometheus exporter conventions for new users.

**Architecture:** Single file replacement. No code changes, no new files beyond README.md itself.

**Tech Stack:** Markdown

---

### Task 1: Rewrite README.md

**Files:**
- Modify: `README.md`

**Step 1: Replace README.md with the new content**

Write the following complete content to `README.md`:

```markdown
# asterisk-exporter

Prometheus metrics exporter for Asterisk PBX. Polls Asterisk via CLI commands (`asterisk -rx`) and exposes channel and bridge metrics on an HTTP endpoint.

Tested with Asterisk 18.

## Prerequisites

- Go 1.24+
- Asterisk PBX (tested with Asterisk 18)

The exporter must run on the same host as Asterisk, since it executes CLI commands directly via `asterisk -rx`.

## Build

```
make build
```

This cross-compiles binaries for linux/amd64, linux/arm64, windows/amd64, windows/386, and darwin/amd64. Output binaries are named `asterisk-exporter-<version>-<os>-<arch>`.

## Usage

```
./asterisk-exporter-0.0.5-linux-amd64 [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-web_listen_address` | `127.0.0.1:9495` | Address to listen on for web interface and telemetry |
| `-web_listen_path` | `/metrics` | Path under which to expose metrics |
| `-asterisk_metric_interval` | `5` | Interval in seconds between metric collection |

## Metrics

### Health

| Metric | Type | Description |
|--------|------|-------------|
| `asterisk_health_fail` | counter | Count of Asterisk health check failures |

### Channel

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `asterisk_crruent_channel_context` | gauge | `context` | Current number of channels by dialplan context |
| `asterisk_crruent_channel_tech` | gauge | `tech` | Current number of channels by channel technology |
| `asterisk_channel_duration` | histogram | `tech`, `context` | Duration of channels in seconds |

### Bridge

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `asterisk_crruent_bridge_count` | gauge | `type`, `tech` | Current number of bridges by type and technology |
| `asterisk_bridge_duration` | histogram | `type`, `tech` | Duration of bridges in seconds |

## Example

```
$ curl localhost:9495/metrics
# HELP asterisk_crruent_bridge_count Current number of bridges in the asterisk.
# TYPE asterisk_crruent_bridge_count gauge
asterisk_crruent_bridge_count{tech="simple_bridge",type="stasis"} 120
# HELP asterisk_crruent_channel_context Current number of channels(context) in the asterisk.
# TYPE asterisk_crruent_channel_context gauge
asterisk_crruent_channel_context{context="call-in"} 1
# HELP asterisk_crruent_channel_tech Current number of channels(tech) in the asterisk.
# TYPE asterisk_crruent_channel_tech gauge
asterisk_crruent_channel_tech{tech="PJSIP"} 1
# HELP asterisk_channel_duration A duration time of the channel
# TYPE asterisk_channel_duration histogram
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="5"} 1
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="10"} 2
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="30"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="60"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="120"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="300"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="600"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="1800"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="3600"} 4
asterisk_channel_duration_bucket{context="call-in",tech="PJSIP",le="+Inf"} 4
asterisk_channel_duration_sum{context="call-in",tech="PJSIP"} 49
asterisk_channel_duration_count{context="call-in",tech="PJSIP"} 4
# HELP asterisk_bridge_duration A duration time of the bridge
# TYPE asterisk_bridge_duration histogram
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="5"} 1
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="10"} 2
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="30"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="60"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="120"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="300"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="600"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="1800"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="3600"} 4
asterisk_bridge_duration_bucket{tech="simple_bridge",type="stasis",le="+Inf"} 22852
asterisk_bridge_duration_sum{tech="simple_bridge",type="stasis"} 1.8440478212e+10
asterisk_bridge_duration_count{tech="simple_bridge",type="stasis"} 22852
# HELP asterisk_health_fail Asterisk health check count
# TYPE asterisk_health_fail counter
asterisk_health_fail{} 0
```

## License

MIT License. See [LICENSE](LICENSE).
```

**Step 2: Verify the README renders correctly**

Run: `cat README.md | head -5`
Expected: Title and description lines appear correctly.

**Step 3: Commit**

```bash
git add README.md
git commit -m "Rewrite README for new users

- asterisk-exporter: Restructure README to follow Prometheus exporter conventions
- asterisk-exporter: Fix version from 0.0.1 to 0.0.5
- asterisk-exporter: Fix default listen address to 127.0.0.1:9495
- asterisk-exporter: Add prerequisites and build sections
- asterisk-exporter: Convert metrics list to table format with labels
- asterisk-exporter: Trim example output to asterisk-specific metrics only
- asterisk-exporter: Add license reference"
```
