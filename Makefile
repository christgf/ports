.DEFAULT_GOAL := help
DOCKER_COMPOSE_YAML := docker-compose.yaml
REPO_NAME := ports

# `make help` generates a help message for each target that
# has a comment starting with ##
help:
	@echo "Use 'make <target>' where <target> is one of the following:"
	@echo "  check   -  run lint and static checks."
	@echo "  fmt     -  format the solution."
	@echo "  build   -  build the service."
	@echo "  test    -  run unit tests, with coverage."
	@echo "  up      -  set up and run a service environment."
	@echo "  down    -  destroy and clean up a service environment."

check:
	go vet ./...
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w `find . -name '*.go' | grep -v "vendor"`
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

fmt:
	go fmt ./...

build:
	go build -ldflags "-w -s" -o $(REPO_NAME) ./cmd/$(REPO_NAME)

# Run `make up` to set up the environment, and set PORTS_MONGODB_CONN_URI to include integration tests.
test:
	go test -count=1 -race -coverprofile=cover.out ./... && go tool cover -func=cover.out && rm cover.out

up:
	MONGO_HOSTNAME=localhost docker-compose -f $(DOCKER_COMPOSE_YAML) up

down:
	docker-compose -f $(DOCKER_COMPOSE_YAML) down -v

.PHONY: help check fmt build test up down