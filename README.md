# BIG-IP exporter
Prometheus exporter for BIG-IP statistics. Uses iControl REST API.

## Eventual improvements
### Gather data in the background
Currently the data is gathered when the /metrics endpoint is called. This causes the request to take a few seconds before completing. This could be fixed by having a go thread that gathers data at regular intervals and that is returned upon a call to the /metrcis endpoint
