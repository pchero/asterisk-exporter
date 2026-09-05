# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Prometheus metrics exporter for Asterisk PBX. Polls Asterisk via CLI commands (`asterisk -rx`) and exposes channel/bridge metrics on an HTTP endpoint. Tested with Asterisk 18.

## Build & Test Commands

```bash
# Build for all platforms (linux/windows/macos)
make build

# Build output: asterisk-exporter-0.0.6-{os}-{arch}
# Cross-compile targets defined in build.sh

# Run tests
go test ./pkg/collector/...

# Run tests verbose
go test -v ./pkg/collector/...

# Run a single test
go test -v -run TestChannelParser ./pkg/collector/...

# Run the binary
./asterisk-exporter-0.0.6-linux-amd64 -web_listen_address=":9495" -asterisk_metric_interval=5
```

## Architecture

```
cmd/asterisk-exporter/main.go  → Entry point: flags, signal handling, HTTP server (:9495/metrics)
pkg/collector/main.go          → Collector interface, Prometheus metric definitions, polling loop
pkg/collector/channel.go       → Parses `asterisk -rx "core show channels concise"` (! delimited)
pkg/collector/bridge.go        → Parses `asterisk -rx "bridge show all"`
models/channel/channel.go      → Channel data struct (14 fields)
models/bridge/bridge.go        → Bridge data struct (5 fields)
```

The collector runs in a goroutine on a configurable interval, executes Asterisk CLI commands via `exec.Command`, parses the text output, and updates Prometheus gauge/histogram metrics.

**Data flow:** `main.go` starts HTTP server + collector goroutine → `collector.Run()` loops → calls `channelCollects()` and `bridgeCollects()` → each executes an Asterisk CLI command, parses output, updates Prometheus metrics.

## Key Details

- Go module path: `gitlab.com/voipbin/voip/asterisk-exporter.git` (go 1.24)
- Metrics namespace: `asterisk` (defined in `pkg/collector/main.go`)
- Metric names contain intentional typo `crruent` (not `current`) — do not "fix" this as it would break existing Prometheus queries
- Tests cover parsing logic only (no mocked Asterisk commands)
- No CI/CD pipeline configured; builds are manual via `make build`
