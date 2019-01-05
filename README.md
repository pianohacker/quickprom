# quickprom
## a quick analysis tool for Prometheus-compatible time-series databases

`quickprom` fills the gulf between full-blown visualization tools like Grafana and manual analysis
using the Prometheus API.

## Features

- Automatically summarizes data:
	- Collapses labels that are shared between samples/series
	- Only shows date once if it's the same between all series
- Supports basic authentication, or automatically using authorization from your CloudFoundry CLI session

## TODO

- [ ] JSON output
- [ ] Acceptance tests of binary
- [ ] Custom sorting
- [ ] Sparklines
- [ ] Scalar support
