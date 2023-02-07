# ficsit-agent

Prometheus Exporter for [Satisfactory](https://www.satisfactorygame.com/).

We depend on the [Ficsit Remote Monitoring](https://ficsit.app/mod/FicsitRemoteMonitoring) mod in order to fetch data from the game itself.

## Dependencies

1. [Satisfactory](https://www.satisfactorygame.com/)
2. Install the [Ficsit Remote Monitoring](https://ficsit.app/mod/FicsitRemoteMonitoring) mod via the Satisfactory Mod Manager.
3. In game, use `/frmweb start`.

## Usage

### External Prometheus

We use [Grafana Agent](https://grafana.com/docs/agent/latest/) to forward data to an external Prometheus instance.

1. Copy `agent.sample.yml` to `agent.yml` and configure your Prometheus server's remote write URL and auth info.
2. `docker compose up`