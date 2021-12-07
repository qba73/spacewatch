![Go](https://github.com/qba73/spacewatch/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/spacewatch)](https://goreportcard.com/report/github.com/qba73/spacewatch)
![GitHub](https://img.shields.io/github/license/qba73/meteo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/meteo)

# spacewatch

The [International Space Station](https://en.wikipedia.org/wiki/International_Space_Station) is a modular space station in low Earth orbit. It is a multinational collaborative project involving five participating space agencies: NASA, Roscosmos, JAXA, ESA, and CSA. The ownership and use of the space station is established by intergovernmental treaties and agreements.

The **Spacewatch service** is a REST API that provides functionality for tracking the International Space Station. The service's primary function is to provide information if the ISS is visible in the sky at the given moment.
We assume the ISS is visible if both conditions are satisfied: the sky cloud coverage is less than 30%, and it is nighttime.

**Note:** The service is under development, and it is not intended for production use.


# development

To use the service at the moment, you need to register at [Weatherbit.io](https://www.weatherbit.io/api) to get an API KEY. Once you get the key, you need to export it as an Env Var ```SPACEWATCH_WEATHER_API_KEY```


### Make targets useful for developemnt and testing:
```
$ make
vet                  Run go vet and shadow
check                Run static check analyzer
cover                Run unit tests and generate test coverage report
test                 Run unit tests locally
tidy                 Run go mod tidy and vendor
run                  Run service locally
```

### Run service locally:
```
$ make run
go run cmd/spacewatch-api/main.go
SPACEWATCH : 2021/12/07 16:56:28.966230 main.go:62: main : Started
SPACEWATCH : 2021/12/07 16:56:28.966388 main.go:69: main : Config :
--web-address=localhost:9000
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--web-cache-ttl=10s

SPACEWATCH : 2021/12/07 16:56:29.024535 main.go:96: main : Spacewatch API listening on localhost:9000
```

### Test the service:
```
$ curl localhost:9000
{"lat":-43.93,"long":-36.43,"timezone":"Atlantic/South_Georgia","cloud_coverage":100,"day_part":"night","is_visible":false}
```
```
$ curl localhost:9000
{"lat":-21.71,"long":-7.31,"timezone":"Atlantic/St_Helena","cloud_coverage":9,"day_part":"day","is_visible":true}
```

Received payload:
```json
{"lat":-21.71,"long":-7.31,"timezone":"Atlantic/St_Helena","cloud_coverage":9,"day_part":"day","is_visible":true}
```

- lat/long - ISS coordinates at the time of the request
- timezone - where is ISS
- clud_coverage in %
- day_part - is it day or night in the timezone
- is_visible - indicates if the ISS can be seen on the sky (cloud coverage <= 30% and night)

# roadmap

- [x] implement server-side caching (minimize load on third party services)
- [ ] replace weather provider (currently there are limitations with no of requests)
- [ ] add [Caddy server](https://caddyserver.com) as a proxy in front of the service (TLS, request rate limiting)
- [ ] add a healthcheck endpoint
- [ ] add request tracing and metrics ([Prometheus](https://prometheus.io/docs/instrumenting/clientlibs/))
- [ ] use context for requests cancellations
- [ ] implement middleware (and move logging functionality)
- [ ] configure autodeployment (GitHub Actions)
- [ ] register the service in the [RapidAPI](https://rapidapi.com/products/enterprise-hub/) for public use

# tbd
- add gRPC interface
