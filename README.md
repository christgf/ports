# Ports
This repo contains a sample implementation of the port domain service, as an [HTTP API](./cmd/ports/main.go), and as described in [technical instructions](INSTRUCTIONS.md).

## Build and run tests
To build the project, use the following command:

```shell
make build    # build the service.
```

You can run project tests using the following command:
```shell
make test      # go test with caching disabled and code coverage.
```

Note that, by default, only unit and short tests are executed. 
To run integration tests, you must first set up a service environment.
You can do so with the following command:
```shell
make up   # set up an environment with service and its dependencies
```

You can then run integration tests by setting the MongoDB database connection URI in the appropriate environment variable:
```shell
PORTS_MONGODB_CONN_URI=mongodb://localhost:27017/papaya make test  
```

To destroy and clean up the service environment, use:
```shell
make down   # stop the service environment, clean up containers, networks, volumes
```

The [Makefile](./Makefile) includes a variety of other useful targets. Simply typing `make` should list all possible options.

## File loader
Apart from the HTTP API, the project contains a separate command-line utility, called [portload](./cmd/portload/main.go), which can be used to read and parse ports information from JSON files, as described in [technical instructions](INSTRUCTIONS.md).
You can build `portload` into an executable using the following command:
```shell
go build -ldflags "-w -s" -o portload ./cmd/portload
```
