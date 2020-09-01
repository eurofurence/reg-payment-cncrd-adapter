# reg-payment-cncrd-adapter

A backend service that provides payment functionality via the Concardis payment provider.

Implemented in Go.

Command line arguments: `-config [config_file]`

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, `go build main.go` and `go test ./...` will download all required dependencies by default.

To build and run the service with default settings, just use `run.sh`.
Alternatively, you can build the service using `go build main.go` and subsequently launch the executable using the command line arguments as listed above.

Go 1.12 or later is required.
