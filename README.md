# Geolocation Resolver Service

## Overview

Geolocation Resolver Service is a service that imports CSV files with raw geolocation data, parses and stores them in a database, and exposes them via an HTTP API.

Sample content of CSV file:
```
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
```

The service is ready to work with CSV files from unreliable sources. It eliminates duplicates and validates entries, but does not check data correctness.

### Libraries

All libraries are located in `pkg` folder.

- `geoloc` - describes the geolocation entries that the service works with;
- `importer` - contains an importer for CSV files with geolocation data;
- `storage/scylla` - contains a connector for ScyllaDB that can be used as a storage for the service (keyspace and table definitions are located in `databases/scylla`);
-  `storage/hashmap` - contains a simple in-memory storage implementation that can be used as a storage for the service.

### Service geoloc_resolver

Service `geoloc_resolver` uses the above libraries and performs 2 things:
1. imports geolocation data from a CSV file,
2. configures an HTTP server to process API requests to resolve IPv4 and IPv6 addresses.

There is a single API endpoint `resolve_ip` with a one URI param `ip`:
  ```bash
  $ curl -v "127.0.0.1:8080/resolve_ip?ip=160.103.7.140"
  < HTTP/1.1 200 OK
  < Content-Type: application/json
  < Date: Thu, 09 Mar 2023 15:39:55 GMT
  < Content-Length: 41
  <
  {"country":"Nicaragua","city":"New Neva"}
  ```

In case of an error, the response is:
  ```bash
  $ curl -v "127.0.0.1:8080/resolve_ip?ip=invalid_ip"
  < HTTP/1.1 400 Bad Request
  < Content-Type: application/json
  < Date: Thu, 09 Mar 2023 15:39:24 GMT
  < Content-Length: 24
  <
  {"message":"invalid ip"}
  ```

If there is no geolocation entry for the requested IP address, the response is:
  ```bash
  $ curl -v "127.0.0.1:8080/resolve_ip?ip=160.103.7.141"
  < HTTP/1.1 204 No Content
  < Date: Thu, 09 Mar 2023 15:41:03 GMT
  < 
  ```

## Getting started

### Build

To build an app, use `Makefile`:

  ```bash
  $ make
  go build -mod vendor -o bin/geoloc_resolver geolocation-resolver/cmd/geoloc_resolver
  ```

### Configuration and execution

Service configuration can be set using a `yaml` config file by passing its path using the `-config` flag:

  ```bash
  ./bin/geoloc_resolver -config ./configs/geoloc_resolver.hashmap.local.yml
  ```

See the `configs` directory for examples of local configurations.

### Development Environment (Local)

To set up a local ScyllaDB cluster using docker-compose, use `deploy_dev.sh` script (located in `scripts/dev`):

  ```bash
  $ ./scripts/dev/deploy_dev.sh start
  Creating network "dev_default" with the default driver
  Creating scylla-node1 ... done
  Creating scylla-node2 ... done
  Creating scylla-node3 ... done
  Initialising scylla... (may take some time)
  Successfully initialised scylla
  ```

Also you can run unit tests and linter checks locally:
  ```bash
  make test
  make lint
  ```

And don't forget about postman tests (see `tests/postman` dir)!

## Tradeoffs and TODO list
1. Use custom logging library for easier grep, log-levels;
2. Add metrics using prometheus/graphite/etc;
3. (0, 0) latitude and longitude coordinates are considered invalid;
4. Add read and write retries inside the `storage/scylla/pool` package;
5. Add API versioning;
6. Last Write Wins - geolocation entries for the same IP address override themselves;
7. Integration tests are the only way to test `storage/scylla` and `internal/geoloc_resolver/http` packages.
