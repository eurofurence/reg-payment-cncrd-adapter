# reg-payment-cncrd-adapter

A backend service that provides payment functionality via the Concardis payment provider.

Implemented in Go.

Command line arguments: `-config [config_file]`

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, go build and go test will download
all required dependencies by default.
