# quickprom
## a quick analysis tool for Prometheus-compatible time-series databases

`quickprom` fills the gulf between full-blown visualization tools like Grafana and manual analysis
using the Prometheus API.

## Features

- Automatically summarizes data:
	- Collapses labels that are shared between samples/series
	- Only shows date once if it's the same between all series
	- Truncates seconds and milliseconds if they're zero for all samples
	- Tries to format all values identically, using the minimum number of digits
- Supports basic authentication, or automatically using authorization from your CloudFoundry CLI session

## Installation
Go 1.11 is required.

```console
$ git clone https://github.com/pianohacker/quickprom
$ cd quickprom
$ GO111MODULE=on go install ./cmd/quickprom
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
| `--json` | Output JSON result (`QUICKPROM_JSON`) |
| `-b, --range-table` | Output range vectors as tables (`QUICKPROM_RANGE_TABLE`) |

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
$ quickprom 'prometheus_http_request_duration_seconds_bucket{le="1"}'
Instant vector:
  At: 2019-01-06 18:45:50.132 MST
  All samples have labels:
    __name__: prometheus_http_request_duration_seconds_bucket
    instance: promserver.example
    job: prometheus
    le: 1

 handler              value
 /                      128 
 /alerts                 58 
 /config                 16 
 /consoles/*filepath    877 
 /flags                  16 
 /graph                 246 
 /label/:name/values    220 
 /metrics             32197 
 /query               19216 
 /query_range         28551 
 /rules                   5 
 /series                725 
 /service-discovery       6 
 /static/*filepath     3900 
 /status                 13 
 /targets                12
$ quickprom range 'prometheus_engine_query_duration_seconds' --start '1:00' --end '2:00' --step '30m' --range-table
Range vector:
  All on date: 2019-01-04
  All timestamps end with: 00.000
  All series have labels:
    __name__: prometheus_engine_query_duration_seconds
    instance: promserver.example
    job: prometheus

 quantile  slice              01:00       01:30       02:00 
 0.5       inner_eval    2.0050e-06  1.9370e-06  1.9890e-06 
 0.5       prepare_time  4.4220e-06  4.2300e-06  4.3320e-06 
 0.5       queue_time    1.4100e-06  1.3750e-06  1.3830e-06 
 0.5       result_sort   7.5900e-07  7.5800e-07  8.7700e-07 
 0.9       inner_eval    3.2960e-06  3.1980e-06  3.3870e-06 
 0.9       prepare_time  6.8050e-06  6.7350e-06  6.9200e-06 
 0.9       queue_time    2.0240e-06  1.9610e-06  2.0500e-06 
 0.9       result_sort   1.2450e-06  1.1230e-06  1.1610e-06 
 0.99      inner_eval    2.4959e-05  1.7474e-05  2.6921e-05 
 0.99      prepare_time  1.8459e-05  1.8900e-05  1.8408e-05 
 0.99      queue_time    3.2850e-06  3.1800e-06  4.2550e-06 
 0.99      result_sort   1.2450e-06  1.2810e-06  1.1610e-06
$ quickprom 'node_timex_status'
Instant vector:
  At: 2019-01-06 18:10:05.628 MST
  All series have labels:
    __name__: node_timex_status
    instance: promserver
    job: node

 value
  8193
```

## TODO

- [ ] Automatically enable range tables, disable when terminal too narrow (needs a decent heuristic)
- [ ] Custom sorting
- [ ] Sparklines
- [ ] Scalar support
