# Quick Start

Want to get up and running quickly? Here are a few examples...

## Basic GET Requests

The `requests.SendGetRequest` function makes a GET request to a URL.

```go
res, err := requests.SendGetRequest("https://example.com/greet")
if err != nil {
    log.Panic(err)
}

fmt.Printf("Status Code: %d\n", res.StatusCode) // Status Code: 200
fmt.Printf("Body: %q\n", string(res.Body))      // Body: "Hello, World!"
```

## Basic POST Requests

The `requests.SendPostRequest` function makes a POST request to a URL, with the specified `body` and `contentType`.

```go
url := "https://example.com/greet"
contentType := "application/json"
body := []byte(`{"name": "gopher"}`)

res, err := requests.SendPostRequest(url, contentType, body)
if err != nil {
    log.Panic(err)
}

fmt.Printf("Status Code: %d\n", res.StatusCode) // Status Code: 200
fmt.Printf("Body: %q\n", string(res.Body))      // Body: "Hello, Gopher!"
```

---

But that's just the start! Keep reading to learn about more advanced ways to build HTTP requests that meet your needs.
