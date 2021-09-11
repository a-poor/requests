# requests

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/requests.svg)](https://pkg.go.dev/github.com/a-poor/requests)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/requests?style=flat-square)](https://pkg.go.dev/github.com/a-poor/requests)
[![Go Test](https://github.com/a-poor/requests/actions/workflows/go.yml/badge.svg)](https://github.com/a-poor/requests/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/requests)](https://goreportcard.com/report/github.com/a-poor/requests)
[![GitHub](https://img.shields.io/github/license/a-poor/requests?style=flat-square)](https://github.com/a-poor/requests/blob/main/LICENSE)
![GitHub last commit](https://img.shields.io/github/last-commit/a-poor/requests?style=flat-square)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/requests/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/requests?badge)

_created by Austin Poor_

A quick and easy HTTP request library written in Go.

## Installation

```bash
go get github.com/a-poor/requests
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/a-poor/requests"
)

func main() {
    // Send the request
    res, err := requests.SendGetRequest("https://google.com")

    // If there was an error, print and return
    if err != nil {
        fmt.Printf("Error: %e\n", err)
        return
    }

    // Print the response's status code
    fmt.Printf("Status Code: %d\n", res.StatusCode)

}
```
