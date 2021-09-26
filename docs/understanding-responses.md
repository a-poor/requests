# Understanding Responses

Sending an HTTP request with a `requests.Request` object will return a `requests.Response`.

## Anatomy of a Response

Here's what the `requests.Response` struct looks like:

```go
type Response struct {
	Ok         bool              // Was the request successful? (Status codes: 100-399)
	StatusCode int               // HTTP response status code
	Headers    map[string]string // HTTP Response headers
	Body       []byte            // HTTP Response body
}
```

`StatusCode` stores the HTTP Response status-code, and helps the user determine if their request was successful. You can read more about HTTP status codes [here](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status).

The value of `Ok` is based on the status code (anything under `400`) and gives the user a more-simplified idea of the success of the request.

`Headers` stores the HTTP response headers. Like the `Request` struct, `Response` has a `GetHeader` function to perform a case-insensitive check for a header.

Finally, `Body` stores the HTTP response body as a slice of bytes.

