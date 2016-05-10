# BIG-IP exporter
Prometheus exporter for BIG-IP statistics. Uses iControl REST API.

## Usage
The bigip_exporter is easy to use. Example: 
```
./bigip_exporter -bigip.host 127.0.0.1 -bigip.port 443 -bigip.username admin -bigip.password admin
```

### Flags
Flag | Description | Default
-----|-------------|---------
-bigip.host | BIG-IP host | localhost
-bigip.port | BIG-IP port | 443
-bigip.username | BIG-IP username | user
-bigip.password | BIG-IP password | pass
-exporter.bind_address | The address the exporter should bind to | All interfaces
-exporter.bind_port | Which port the exporter should listen on | 9142
-exporter.partitions | A comma separated list containing the partitions that should be exported | All partitions
-exporter.namespace | The namespace used in prometheus labels | bigip

## Limitations
iControl REST API have been known to be pretty slow. Therefore requests to the `/metrics` endpoint may take a while to complete, around 15 seconds. Depending on the amount of configuration you have in your BIG-IP this may differ. It may therefore be necessary to tweak `scrape_interval` and `scrape_timeout` in your prometheus configuration to match this.

## Eventual improvements
### Gather data in the background
Currently the data is gathered when the `/metrics` endpoint is called. This causes the request to take a few seconds before completing. This could be fixed by having a go thread that gathers data at regular intervals and that is returned upon a call to the `/metrics` endpoint. This would however go against the [guidelines](https://prometheus.io/docs/instrumenting/writing_exporters/#scheduling).
