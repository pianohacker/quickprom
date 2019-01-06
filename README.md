# quickprom
## a quick analysis tool for Prometheus-compatible time-series databases

`quickprom` fills the gulf between full-blown visualization tools like Grafana and manual analysis
using the Prometheus API.

## Features

- Automatically summarizes data:
	- Collapses labels that are shared between samples/series
	- Only shows date once if it's the same between all series
- Supports basic authentication, or automatically using authorization from your CloudFoundry CLI session

## Installation
Go 1.11 is required.

```console
$ GO111MODULE=on go get github.com/pianohacker/quickprom/cmd/quickprom
```

## Usage
```
  quickprom [options] QUERY [--time TIME]
  quickprom [options] range QUERY --start START [--end END] --step STEP
```

### Global options
| Option | Description |
| ------ | ----------- |
| `-t, --target TARGET` | URL of Prometheus-compatible target (`QUICKPROM_TARGET`) |
| `-k, --skip-tls-verify` | Don't verify remote certificate (`QUICKPROM_SKIP_TLS_VERIFY`)  |
| `--basic-auth USER:PASS` | Use basic authentication (`QUICKPROM_BASIC_AUTH`) |
| `--cf-auth` | Automatically use current oAuth token from `cf` (`QUICKPROM_CF_AUTH`)  |
| --json | Output JSON result (`QUICKPROM_JSON`) |

### Instant query options
| Option | Description |
| ------ | ----------- |
| `-i, --time TIME` | Evaluate instant query at `TIME` (defaults to now) |

### Range query options
| Option | Description |
| ------ | ----------- |
| `-s, --start START` | Start time of range query |
| `-e, --end END` | End time of range query (inclusive, defaults to now) |
| `-p, --step STEP` | Step of range query |

### Timestamp format
quickprom uses the excellent fuzzytime library, and thus supports a number of
formats for the --time, --start, --end and --step options. Each takes a date
and/or time, separated by a space. If you leave out the date, today is
assumed, and if you leave out the time, local midnight is assumed.

Some examples:
  - 2010-04-02
  - 11/02/2008 4:48PM GMT
  - 11.02.10 13:21:36+00:00
  - 14:21:01
  - 14:21
  - 2019-01-01T00:12:34Z

## Examples

```console
$ export QUICKPROM_TARGET=http://promserver.example
$ quickprom 'prometheus_engine_query_duration_seconds'
Instant vector:
  At: 2019-01-04 22:37:22.944 MST
  All have labels: __name__: prometheus_engine_query_duration_seconds, instance: promserver.example, job: prometheus

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
$ quickprom range 'prometheus_engine_query_duration_seconds' --start '1:00' --end '2:00' --step '30m' --range-table
Range vector:
  All have labels: __name__: prometheus_engine_query_duration_seconds, instance: promserver.example, job: prometheus
  All on date: 2019-01-04
  All timestamps end with: 00.000

 quantile  slice         01:00     01:30     02:00
 0.5       inner_eval    0.000002  0.000002  0.000002
 0.5       prepare_time  0.000004  0.000004  0.000004
 0.5       queue_time    0.000001  0.000001  0.000001
 0.5       result_sort   0.000001  0.000001  0.000001
 0.9       inner_eval    0.000003  0.000003  0.000003
 0.9       prepare_time  0.000007  0.000007  0.000007
 0.9       queue_time    0.000002  0.000002  0.000002
 0.9       result_sort   0.000001  0.000001  0.000002
 0.99      inner_eval    0.000017  0.000023  0.000016
 0.99      prepare_time  0.000020  0.000019  0.000017
 0.99      queue_time    0.000003  0.000003  0.000003
 0.99      result_sort   0.000003  0.000003  0.000003
$ quickprom 'avg_over_time(prometheus_engine_query_duration_seconds[30s])' --time '4:00'
Instant vector:
  At: 2019-01-04 04:00:00.000 MST
  All have labels: instance: promserver.example, job: prometheus

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

- [ ] Better automatic value formatting
- [ ] Automatically enable range tables, disable when terminal too narrow (needs a decent heuristic)
- [ ] Custom sorting
- [ ] Sparklines
- [ ] Scalar support
