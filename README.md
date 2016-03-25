
# etcd Golang Integration Test Harness

[![Travis Build](https://travis-ci.org/mwitkow/go-etcd-harness.svg)](https://travis-ci.org/mwitkow/go-etcd-harness)
[![Go Report Card](http://goreportcard.com/badge/mwitkow/go-etcd-harness)](http://goreportcard.com/report/mwitkow/go-etcd-harness)
[![GoDoc](http://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/mwitkow/go-etcd-harness)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Harness code that spins up and manages a local-only [etcd](https://github.com/coreos/etcd) server for Go (golang) 
integration tests.

It is meant to be used with suite-tests, such as [testify Suites](https://godoc.org/github.com/stretchr/testify/suite) or [gocheck fixtures](https://labix.org/gocheck),
to leverage single startup of a test suite.


## Features:

 * dynamic port allocation - multiple Harnesses can be run in parallel
 * dynamic data directory - each Harness has an independent data dirctory
 * full `etcd` server lifecycle management - start, stop and data cleanup
 * uses `etcd` from `$PATH` - makes it easy to test against an `etcd` version you run in production
 * typical boostrap of 400-500ms
 
## Usage

Harness exports the [official Golang client bindings](https://godoc.org/github.com/coreos/etcd/client) under 
 `harness.Client`. For raw access purposes `harness.Endpoint` returns the client-side endpoint of the etcd server.

For an example of usage that utilises [testify Suites](https://godoc.org/github.com/stretchr/testify/suite), please see
[harness_test.go](harness_test.go).

## Travis CI Configuration

If you're creating an open source project, chances are you're using Travis as your CI. Travis by default places your
`$GOPATH/bin` on the build's `$PATH`. Since etcd Harness looks for the `etcd` binary in `$PATH`, using it is as easy as:

```yaml
install:
  - go get github.com/coreos/etcd
  - go get github.com/mwitkow/go-etcd-harness
  ...
```

###License

etcd-harness is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
