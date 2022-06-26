# asterisk-exporter

Asterisk metric exporter for prometheus.
Tested with asterisk-18.

# Run

```
Usage of ./asterisk-exporter-0.0.1-linux-amd64:
  -asterisk_metric_interval int
        Interval sec for metric getting (default 5)
  -web_listen_address string
        Address to listen on for web interface and telemetry. (default ":9495")
  -web_listen_path string
        Path under which to expose metrics. (default "/metrics")
```

# Metrics
* asterisk_health_fail: Counter of asterisk health check failure.

## channel
* asterisk_crruent_channel_context: Shows current number of channels with context.
* asterisk_crruent_channel_tech: Shows current number of channels with tech.
* asterisk_channel_duration_bucket: Bucket for channel's duration.

## bridge
* asterisk_current_bridge_count: Shows current number of bridges with type/tech.
* asterisk_bridge_duration_bucket: Bucket for bridge's duration.

# Example
```
$ curl localhost:9495/metrics
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
# HELP asterisk_crruent_bridge_count Current number of bridges in the asterisk.
# TYPE asterisk_crruent_bridge_count gauge
asterisk_crruent_bridge_count{tech="simple_bridge",type="stasis"} 120
# HELP asterisk_crruent_channel_context Current number of channels(context) in the asterisk.
# TYPE asterisk_crruent_channel_context gauge
asterisk_crruent_channel_context{context="call-in"} 1
# HELP asterisk_crruent_channel_tech Current number of channels(tech) in the asterisk.
# TYPE asterisk_crruent_channel_tech gauge
asterisk_crruent_channel_tech{tech="PJSIP"} 1
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.002e-05
go_gc_duration_seconds{quantile="0.25"} 9.0926e-05
go_gc_duration_seconds{quantile="0.5"} 0.000123781
go_gc_duration_seconds{quantile="0.75"} 0.000213209
go_gc_duration_seconds{quantile="1"} 0.096603235
go_gc_duration_seconds_sum 0.119995573
go_gc_duration_seconds_count 44
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
...
```
