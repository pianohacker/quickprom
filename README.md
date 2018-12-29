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

- [x] Table output
- [ ] JSON output
- [ ] Custom sorting
- [ ] Range query support
- [ ] Sparklines
- [ ] Scalar support
