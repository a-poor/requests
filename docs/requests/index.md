# Building Requests

Most of the time you're making requests, you'll likely want to configure the request parameters.

## Anatomy of a Request

HTTP requests made with the `requests` library are all based around the `Request` struct. 

Here's what the `requests.Request` struct looks like:

```go
type Request struct {
	URL     string            // URL to send the request to
	Method  HTTPMethod        // HTTP method to use
	Headers map[string]string // Headers to send with the request
	Query   map[string]string // Query parameters to send with the request
	Body    []byte            // Body to send with the request
	Timeout time.Duration     // Timeout for the request
}
```

The request fields can be set directly or using helper functions.

## Request Parameters

Most of the parameters are probably pretty self explanitory but here's a breakdown of what each parameter means:

Name | Function | Example
-----|----------|-------
URL  | The HTTP endpoint being accessed  | `http://google.com`
Method | HTTP Method being used | `GET`,`POST`,`PUT`
Headers | HTTP Headers | `Content-Type`,`Authorization`
Query | Query parameters added to the end of the URL | `<url>?name=Gopher`
Body | HTTP Request body | _ie_ JSON POST data
Timeout | Time before the request will timeout (Note: `Timeout: 0` will never timeout) | `0`, `5 * time.Second`


For more information, consult the [docs](https://pkg.go.dev/github.com/a-poor/requests).


## Request Methods

The HTTP methods are represented as enums/constants of type `requests.HTTPMethod` like so:

```go
const (
	GET     HTTPMethod = iota // An HTTP GET method
	POST                      // An HTTP POST method
	PUT                       // An HTTP PUT method
	DELETE                    // An HTTP DELETE method
	OPTIONS                   // An HTTP OPTIONS method
	HEAD                      // An HTTP HEAD method
	CONNECT                   // An HTTP CONNECT method
	TRACE                     // An HTTP TRACE method
	PATCH                     // An HTTP PATCH method
)
```

They have one method, `String` which converts them to their text representation (_ex_ `requests.GET.String == "GET"`).

## Sending Requests

Once you have a request that you're ready to send, you can call the `Send` method.

Send has the following signature:

```go
func (req *Request) Send() (*Response, error)
```

And here's an example of it in action

```go
// Create your Request
req := requests.Request{
	URL: "https://example.com",
	Method: requests.GET
}

// Send your request
res, err := req.Send()
if err != nil {
	log.Panic(err)
}

// Do something with the response...
log.Println(string(res.Body))
// Output: Hello, world!
```

Alternatively, if you're _sure_ no error will be raised, you can use the `req.MustSend` method which has the following signature:

```go
func (req *Request) MustSend() *Response
```

__WARNING__: This is probably a bad idea, most of the time.

Under the hood, this will call `req.Send` and, if an error is encountered, it will cause a panic.

