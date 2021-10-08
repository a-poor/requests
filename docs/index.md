# requests

[![Go Reference](https://pkg.go.dev/badge/github.com/a-poor/requests.svg)](https://pkg.go.dev/github.com/a-poor/requests)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-poor/requests?style=flat-square)](https://pkg.go.dev/github.com/a-poor/requests)
[![Go Test](https://github.com/a-poor/requests/actions/workflows/go.yml/badge.svg)](https://github.com/a-poor/requests/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-poor/requests)](https://goreportcard.com/report/github.com/a-poor/requests)
[![GitHub](https://img.shields.io/github/license/a-poor/requests?style=flat-square)](https://github.com/a-poor/requests/blob/main/LICENSE)
![GitHub last commit](https://img.shields.io/github/last-commit/a-poor/requests?style=flat-square)
[![Sourcegraph](https://sourcegraph.com/github.com/a-poor/requests/-/badge.svg)](https://sourcegraph.com/github.com/a-poor/requests?badge)
[![CodeFactor](https://www.codefactor.io/repository/github/a-poor/requests/badge/main)](https://www.codefactor.io/repository/github/a-poor/requests/overview/main)

_created by Austin Poor_

Welcome to the documentation for the Go package [github.com/a-poor/request].

`requests` is quick and easy HTTP request library written in Go. 

This library is inspired by the Python Requests library. I wrote it for myself in order to make the HTTP client process a little more _ergonomic_ when writing Go code.

## Table of Contents

* [Installation](#installation)
* [Quick Start](#quick-start)
* [Dependencies](#dependencies)
* [Contributing](#contributing)
* [License](#license)

## Installation

Installation is quick and easy!

```bash
go get github.com/a-poor/requests
```

## Quick Start

Here's a quick example of `requests` in action.

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

## Dependencies

Only the standard library!

## Contributing

Pull requests are super welcome! For major changes, please open an issue first to discuss what you would like to change. And please make sure to update tests as appropriate.

_Or_... feel free to just open an issue with some thoughts or suggestions or even just to say _Hi_ and tell me if this library has been helpful!

## License

[MIT](./LICENSE)


