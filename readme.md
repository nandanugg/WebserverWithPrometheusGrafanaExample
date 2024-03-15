# Go Webserver with Prometheus + Grafana example

This repository will showcase on how I implement metrics "*safely*" by preventing high cardinality for more performant webserver

>About prometheus high cardinality, [read more in here](https://grafana.com/blog/2022/10/20/how-to-manage-high-cardinality-metrics-in-prometheus-and-kubernetes/):

## Known upsides and downsides
- This way, we can prevent high cardinality by not sending the path values of an route (it will send `/:userId` instead of `/isfvDGFs` )
- Not found routes and panics are not handled, but you can make one easily by using each web framework global error handling
- As I forsee, it requires more acrobatic unit test if you want to make one (my skill issue)

## Prerequisites
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/)
- A linux environment (WSL and Macs should be fine)

## How to run
#### Start Prometheus server
```bash
# Prometheus will available at localhost:9090
make startProm
```


#### Start Grafana server
```bash
# Grafana will available at localhost:3000
# the default username & password are both `admin`
make startGrafana
```


#### Start Echo WebServer
```bash
# Echo webserver will available at localhost:8080
go run echoExample/main.go
```
#### Start Fiber WebServer
```bash
# Fiber webserver will available at localhost:8080
go run fiberExample/main.go
```