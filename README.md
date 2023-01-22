# reg-payment-cncrd-adapter

<img src="https://github.com/eurofurence/reg-payment-cncrd-adapter/actions/workflows/go.yml/badge.svg" alt="test status"/>

## Overview

A backend service that provides payment functionality via the Concardis payment provider.

Implemented in Go.

Command line arguments
```
-config <path-to-config-file> [-ecs-json-logging]
```

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository outside your GOPATH, build and test runs will download all required
dependencies by default.

## Running on localhost

Copy the configuration template from `docs/config-template.yaml` to `./config.yaml`. This will set you up
for operation with an in-memory database and sensible defaults.

Build using `go build cmd/main.go`.

Then run `./main -config config.yaml`.

## Installation on the server

See `install.sh`. This assumes a current build, and a valid configuration template in specific filenames.

## Test Coverage

In order to collect full test coverage, set go tool arguments to `-covermode=atomic -coverpkg=./internal/...`,
or manually run
```
go test -covermode=atomic -coverpkg=./internal/... ./...
```

## Contract Testing

This microservice uses [pact-go](https://github.com/pact-foundation/pact-go#installation) for contract tests.

Before you can run the contract tests in this repository, you need to run the consumer side contract tests
in the [reg-payment-service](https://github.com/eurofurence/reg-payment-service) to generate
the contract specifications. It is sufficient to just run what's under `test/contract/consumer`.

You are expected to clone that repository into a directory called `reg-payment-service`
right next to this repository. If you wish to place your contract specs somewhere else, simply change the
path or URL in `test/contract/producer/setup_ctr_test.go`.
