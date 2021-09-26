# Headers and Query Parameters

There are _get_, _set_, and _del_ functions to help with request [Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers) and [Query Parameters](https://en.wikipedia.org/wiki/Query_string).

_Note_: `requests` doesn't support multiple values for headers or query parameters. _Sorry :\\_

## HTTP Headers

Headers are stored as a `map[string]string`.

When forming a request, headers can be set directly, like so:

```go
req := Requests{
    Headers: map[string]string{
        "Content-Type":  "application/json",
        "Accept":        "application/json",
        "Authorization": "Bearer abc123...",
    }
}
```

Headers also have the following helper functions: 

```go
func (req *Request) GetHeader(name string) (string, bool)

func (req *Request) SetHeader(name, value string)

func (req *Request) DelHeader(name string)
```

These functions also come in handy since HTTP headers are [case-insensitive](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers).

Each of the functions will convert the `name` parameter to lowercase before setting / getting / deleting.

Here's an example of headers in action:

```go
// Create a request
req := requests.Request{}

// Set a header
req.SetHeader("CoNtEnT-tYpE", "application/json")

// Get the header (with different case)
v, ok := req.GetHeader("Content-Type")

// Confirm the header was found
fmt.Println(ok) // Output: true
fmt.Println(v)  // Output: application/json
```

## Query Parameters

Similar to Headers, there are the following methods for setting / getting / deleting Query parameters.

Query strings are encoded as part of a URL. URLs have the following structure:

```
scheme://host:port/path?query#fragment
```

The following example...

```
https://example.com/about?name=gopher&language=en
```

...breaks down into the following components...

Component | Example
----------|--------
Scheme | `HTTPS`
Host | `example.com`
Port | none
Path | `/about`
Query | `map[string]string{"name": "gopher", "lang": "en"}`
Fragment | none

_Note_: Unlike HTTP Headers, Query parameters _are_ case sensitive.

The `requests.Request` Query parameters can also be set directly. For example:

```go
req := Requests{
    Query: map[string]string{
        "name": "gopher",
        "lang": "en",
    }
}
```

Or via the following helper functions:

```go
func (req *Request) GetQuery(name string) (string, bool)

func (req *Request) SetQuery(name, value string)

func (req *Request) DelQuery(name string)
```

_Note_: When a `requests.Request` is being sent, the Query parameters are encoded using Go's `net/url` package. This provides safety for the user when passing in the URL rather than naively concatenating the query-string to the URL (_ex_ `fmt.Sprintf("%s?%s", req.URL, encodedQuery)`).

