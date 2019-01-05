# quickprom
## a quick analysis tool for Prometheus-compatible time-series databases

`quickprom` fills the gulf between full-blown visualization tools like Grafana and manual analysis
using the Prometheus API.

## Features

- Shows labels that are shared between samples/series, and hides date if it's the same between all series
- Can automatically use authorization from your CloudFoundry CLI session
- If outputting to a terminal:
	- Labels are organized into a table

## TODO

- [ ] Basic auth support
- [ ] JSON output
- [ ] Acceptance tests of binary
- [ ] Custom sorting
- [ ] Sparklines
- [ ] Scalar support
