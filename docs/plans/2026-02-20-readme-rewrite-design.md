# README Rewrite Design

## Goal

Rewrite the README to follow standard Prometheus exporter conventions, making it useful for new users discovering the project.

## Approach

Standard Prometheus exporter README (matching conventions of node_exporter, mysqld_exporter, etc.): reference-oriented, single file, concise.

## Sections

1. **Title + Description** — One-liner: Prometheus exporter for Asterisk PBX. Polls via CLI, exposes channel/bridge metrics. Tested with Asterisk 18.

2. **Prerequisites** — Go 1.24+, Asterisk (tested with 18). Note the exporter must run on the same host as Asterisk since it calls `asterisk -rx` directly.

3. **Build** — `make build` cross-compiles for linux/windows/darwin. Note output binary naming convention.

4. **Usage** — Updated flag reference with version 0.0.5, correct default address `127.0.0.1:9495`. Basic run command example.

5. **Metrics** — Table format: metric name, type (gauge/histogram/counter), description. Grouped by health, channel, bridge.

6. **Example Output** — Trimmed `curl` example showing only asterisk-specific metrics (drop go runtime metrics).

7. **License** — One line pointing to LICENSE file.

## What Changes

- Fix version from 0.0.1 to 0.0.5
- Fix default listen address from `:9495` to `127.0.0.1:9495`
- Add prerequisites section
- Add build instructions
- Convert metrics list to table format
- Trim example output to relevant metrics only
- Add license reference

## What Stays the Same

- Metric names (including intentional `crruent` typo) are documented as-is
- Same example metrics output content (just trimmed)
