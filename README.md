# asterisk-exporter

Asterisk metric exporter for prometheus

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
* asterisk_crruent_channel_context: Shows current number of channels with context.
* asterisk_crruent_channel_tech: Shows current number of channels with tech.
* asterisk_channel_duration_bucket: Bucket for channel's duration.
