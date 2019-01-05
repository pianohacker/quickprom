# quickprom
## a quick analysis tool for Prometheus-compatible time-series databases

`quickprom` fills the gulf between full-blown visualization tools like Grafana and manual analysis
using the Prometheus API.

## Features

- Automatically summarizes data:
	- Collapses labels that are shared between samples/series
	- Only shows date once if it's the same between all series
- Supports basic authentication, or automatically using authorization from your CloudFoundry CLI session

## Usage
```
quickprom - run queries against Prometheus-compatible databases

Usage:
  quickprom [options] QUERY [--time TIME]
  quickprom [options] range QUERY --start START [--end END] --step STEP

Options:
  -t, --target TARGET     URL of Prometheus-compatible target (QUICKPROM_TARGET)
  -k, --skip-tls-verify   Don't verify remote certificate (QUICKPROM_SKIP_TLS_VERIFY)
  --basic-auth USER:PASS  Use basic authentication (QUICKPROM_BASIC_AUTH)
  --cf-auth               Automatically use current oAuth token from `cf` (QUICKPROM_CF_AUTH)
  --time TIME             Evaluate instant query at `TIME` (defaults to now)
  --start START           Start time of range query
  --end END               End time of range query (inclusive, defaults to now)
  --step STEP             Step of range query
```

## Examples

```console
$ export QUICKPROM_TARGET=http://promserver.example
$ quickprom 'prometheus_engine_query_duration_seconds'
Instant vector:
  At: 2019-01-04 22:37:22.944 MST
  All have labels: __name__: prometheus_engine_query_duration_seconds, instance: demo.robustperception.io:9090, job: prometheus

  quantile  slice                   
       0.5  inner_eval    0.000002  
       0.5  prepare_time  0.000004  
       0.5  queue_time    0.000001  
       0.5  result_sort   0.000001  
       0.9  inner_eval    0.000003  
       0.9  prepare_time  0.000007  
       0.9  queue_time    0.000002  
       0.9  result_sort   0.000001  
      0.99  inner_eval    0.000021  
      0.99  prepare_time  0.000017  
      0.99  queue_time    0.000003  
      0.99  result_sort   0.000001  
$ quickprom 'prometheus_engine_query_duration_seconds[30s]'
Range vector:
  All have labels: __name__: prometheus_engine_query_duration_seconds, instance: demo.robustperception.io:9090, job: prometheus
  All on date: 2019-01-04

quantile: 0.5, slice: inner_eval:
    22:37:12.707: 0.000002
    22:37:22.707: 0.000002
    22:37:32.707: 0.000002
quantile: 0.5, slice: prepare_time:
    22:37:12.707: 0.000004
    22:37:22.707: 0.000004
    22:37:32.707: 0.000004
quantile: 0.5, slice: queue_time:
    22:37:12.707: 0.000001
    22:37:22.707: 0.000001
    22:37:32.707: 0.000001
...
$ quickprom 'avg_over_time(prometheus_engine_query_duration_seconds[30s])'
Instant vector:
  At: 2019-01-04 22:37:57.916 MST
  All have labels: instance: demo.robustperception.io:9090, job: prometheus

  quantile  slice                   
       0.5  inner_eval    0.000002  
       0.5  prepare_time  0.000004  
       0.5  queue_time    0.000001  
       0.5  result_sort   0.000001  
       0.9  inner_eval    0.000003  
       0.9  prepare_time  0.000007  
       0.9  queue_time    0.000002  
       0.9  result_sort   0.000001  
      0.99  inner_eval    0.000021  
      0.99  prepare_time  0.000017  
      0.99  queue_time    0.000003  
      0.99  result_sort   0.000001  
```

## TODO

- [ ] JSON output
- [ ] Acceptance tests of binary
- [ ] Custom sorting
- [ ] Sparklines
- [ ] Scalar support
