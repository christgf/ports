# Ports

This repo contains a sample implementation of the port domain service, as an [HTTP API](./cmd/ports/main.go), and as
described in [technical instructions](INSTRUCTIONS.md).

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

You can then run integration tests by setting the MongoDB database connection URI in the appropriate environment
variable:

```shell
PORTS_MONGODB_CONN_URI=mongodb://localhost:27017/papaya make test  
```

To destroy and clean up the service environment, use:

```shell
make down   # stop the service environment, clean up containers, networks, volumes
```

The [Makefile](./Makefile) includes a variety of other useful targets. Simply typing `make` should list all possible
options.

### Flags and Environment variables

The following command-line flags or environment variables can be used to configure the ports HTTP API.

| Flag                 | Description            | Environment variable      | Default value                   |
|----------------------|------------------------|---------------------------|---------------------------------|
| `-http-listen-addr`  | HTTP listener address  | `PORTS_HTTP_LISTEN_ADDR`  | `:80`                           |
| `-mongodb-conn-uri`  | MongoDB connection URI | `PORTS_MONGODB_CONN_URI`  | mongodb://localhost:27017/ports |

---

## File loader

Apart from the HTTP API, the project contains a separate command-line utility, called [portload](./cmd/portload/main.go)
, which can be used to read and parse ports information from JSON files, as described
in [technical instructions](INSTRUCTIONS.md).
You can build `portload` into an executable using the following command:

```shell
go build -o portload ./cmd/portload
```

### Running the file loader

To run the file loader, provide the path to a JSON file using a command-line flag.
For example, to import the ports JSON file `ports.json` residing in a folder named `testdata`, use:  

```shell
portload -f testdata/ports.json
```

The file loader will begin printing the records parsed, prefixed with a record counter. E.g.
```shell

...
main 1: Port: {AEAJM Ajman 52000 Ajman United Arab Emirates Asia/Dubai}
main 2: Port: {AEAUH Abu Dhabi 52001 Abu Dhabi United Arab Emirates Asia/Dubai}
main 3: Port: {AEDXB Dubai 52005 Dubai United Arab Emirates Asia/Dubai}
main 4: Port: {AEFJR Al Fujayrah 52005 Al Fujayrah United Arab Emirates Asia/Dubai}
...
main 1629: Port: {ZMLUN Lusaka 79145 Lusaka Lusaka Zambia  Africa/Lusaka}
main 1630: Port: {ZWBUQ Bulawayo 79145 Bulawayo Bulawayo Zimbabwe Africa/Harare}
main 1631: Port: {ZWHRE Harare 79145 Harare Harare Zimbabwe Africa/Harare}
main 1632: Port: {ZWUTA Mutare 79145 Mutare Manicaland Zimbabwe Africa/Harare}
```

### Malformed records
Note that the file loader will immediately stop processing the file on the first error it encounters.

