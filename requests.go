package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

// URLEncode will encode `data` as a URL-safe string
// Uses fmt.Sprint and url.PathEscape to do the encoding.
func URLEncode(data interface{}) string {
	return url.PathEscape(fmt.Sprint(data))
}

// makeTemplate creates a new Go template instance
// pre-loaded with the URLEncode function.
func makeTemplate() *template.Template {
	return template.New("url").Funcs(template.FuncMap{
		"URLEncode": URLEncode,
	})
}

// JSONMust marshals a map into a JSON byte slice
// using json.Marshal and panics if there is an error.
func JSONMust(data map[string]interface{}) []byte {
	res, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return res
}

// HTTPMethod is a type that represents an
// HTTP request method.
// Read more here: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
type HTTPMethod int

// Enums representing HTTP methods
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

// Convert an HTTPMethod to it's string format
func (m HTTPMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case OPTIONS:
		return "OPTIONS"
	case HEAD:
		return "HEAD"
	case CONNECT:
		return "CONNECT"
	case TRACE:
		return "TRACE"
	case PATCH:
		return "PATCH"

	}
	return ""
}

// Request is a type that represents an HTTP request
//
// Notes:
// ======
// * URL is assumed not to include query parameters or fragments
// * Headers and Query don't support multiple values
// * Timeout of 0 means no timeout
type Request struct {
	URL     string            // URL to send the request to
	Method  HTTPMethod        // HTTP method to use
	Headers map[string]string // Headers to send with the request
	Query   map[string]string // Query parameters to send with the request
	Body    []byte            // Body to send with the request
	Timeout time.Duration     // Timeout for the request
}

// NewGetRequest creates a new Request object
// with the supplied URL and sets the HTTP method
// to GET.
func NewGetRequest(url string) *Request {
	return &Request{
		URL:    url,
		Method: GET,
	}
}

// SendGetRequest creates a new HTTP GET request
// and sends it to the specified URL.
// Internally, calls `NewGetRequest(url).Send()`
func SendGetRequest(url string) (*Response, error) {
	return NewGetRequest(url).Send()
}

// NewPostRequest creates a new Request object
// with the supplied URL, content-type header, and body sets the HTTP method
// to POST.
func NewPostRequest(url string, contentType string, body []byte) *Request {
	return &Request{
		URL:     url,
		Method:  POST,
		Headers: map[string]string{"content-type": contentType},
		Body:    body,
	}
}

// SendPostRequest creates a new HTTP POST request
// and sends it to the specified URL.
// Internally, calls `NewPostRequest(url, contentType, body).Send()`
func SendPostRequest(url string, contentType string, body []byte) (*Response, error) {
	return NewPostRequest(url, contentType, body).Send()
}

// Copy will create a copy of the Request object
func (req *Request) Copy() *Request {
	r := Request{
		URL:     req.URL,
		Method:  req.Method,
		Timeout: req.Timeout,
	}
	if req.Headers != nil {
		r.Headers = make(map[string]string)
		for k, v := range req.Headers {
			r.Headers[k] = v
		}
	}
	if req.Query != nil {
		r.Query = make(map[string]string)
		for k, v := range req.Query {
			r.Query[k] = v
		}
	}
	if req.Body != nil {
		r.Body = make([]byte, len(req.Body))
		copy(r.Body, req.Body)
	}
	return &r
}

// getReqBody returns the request body as a buffer that can be
// passed to the http.NewRequest function
func (req *Request) getReqBody() *bytes.Buffer {
	return bytes.NewBuffer(req.Body)
}

// getURL returns the string formatted URL with
// the query parameters
func (req *Request) getURL() (string, error) {
	// Make sure there's a URL
	if req.URL == "" {
		return "", fmt.Errorf("URL is required")
	}

	// Encode the query parameters (if any)
	vals := url.Values{}
	for k, v := range req.Query {
		vals.Set(k, v)
	}
	q := vals.Encode()

	// Format the URL with the query parameters (if any)
	u := req.URL
	if q != "" {
		u = fmt.Sprintf("%s?%s", u, q)
	}

	return u, nil
}

// formatPath will use Go templates to format the path
// using the data parameter.
func (req *Request) formatPath(data interface{}) (string, error) {
	tmpl, err := makeTemplate().Parse(req.URL)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParsePathParams will create a copy of the Request object
// and replace URL parameters with the supplied data.
//
// Note: The URL template has access to the `URLEncode` function
// which can be used to safely encode a string. (ex `{{ "Hello world" | URLEncode }}`
// will return `Hello%20world`)
func (req *Request) ParsePathParams(data interface{}) (*Request, error) {
	u, err := req.formatPath(data)
	if err != nil {
		return nil, err
	}
	r := req.Copy()
	r.URL = u
	return r, nil
}

// MustParsePathParams is the same as ParsePathParams except it panics
// if there is an error.
func (req *Request) MustParsePathParams(data interface{}) *Request {
	r, err := req.ParsePathParams(data)
	if err != nil {
		panic(err)
	}
	return r
}

// GetHeader gets a header value from the request. Normalizes the key
// to lowercase before checking. Returns the value associated with the
// key and whether it exists.
func (req *Request) GetHeader(name string) (string, bool) {
	// Create the map if it doesn't exist
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	// Do a case-insensitive check for the key
	for k, v := range req.Headers {
		if strings.ToLower(k) == key {
			return v, true
		}
	}

	// Key not found
	return "", false
}

// SetHeader sets a header value in the request. Normalizes the key
// to lowercase before setting.
func (req *Request) SetHeader(name, value string) {
	// Create the map if it doesn't exist
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	// Add the header to the Headers map
	req.Headers[key] = value
}

// DelHeader deletes a header value from the request headers
// if it exists. Normalizes the key to lowercase
// before deleting.
func (req *Request) DelHeader(name string) {
	// Create the map if it doesn't exist
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	for k := range req.Headers {
		if strings.ToLower(k) == key {
			// Delete the header if it exists
			delete(req.Headers, k)
		}
	}
}

// GetQuery gets a query value from the request. Returns the
// value associated with the key and whether it exists.
func (req *Request) GetQuery(name string) (string, bool) {
	// Create the map if it doesn't exist
	if req.Query == nil {
		req.Query = make(map[string]string)
	}

	val, ok := req.Query[name]
	return val, ok
}

// SetQuery sets a header value in the request.
func (req *Request) SetQuery(name, value string) {
	// Create the map if it doesn't exist
	if req.Query == nil {
		req.Query = make(map[string]string)
	}

	// Set the Query param if it exists
	req.Query[name] = value
}

// DelQuery deletes a query value from the request headers
// if it exists.
func (req *Request) DelQuery(name string) {
	// Create the map if it doesn't exist
	if req.Query == nil {
		req.Query = make(map[string]string)
	}

	// Delete the query param if it exists
	delete(req.Query, name)
}

// Send sends the HTTP request with the supplied parameters
func (req *Request) Send() (*Response, error) {
	// Create an http client (with optional timeout)
	client := http.Client{
		Timeout: req.Timeout,
	}

	// Format the URL with the query parameters (if any)
	u, err := req.getURL()
	if err != nil {
		return nil, err
	}

	// Create the underlying request
	httpRequest, err := http.NewRequest(req.Method.String(), u, req.getReqBody())
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers in the underlying request
	for k, v := range req.Headers {
		httpRequest.Header.Set(k, v)
	}

	// Make the reuquest
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	// Add return headers
	rHeaders := make(map[string]string)
	for k, v := range httpResponse.Header {
		if len(v) > 0 {
			lowerKey := strings.ToLower(k)
			rHeaders[lowerKey] = v[0]
		}
	}

	// Load the request body
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	// Format the response & return
	res := Response{
		Ok:         httpResponse.StatusCode < 400,
		StatusCode: httpResponse.StatusCode,
		Headers:    rHeaders,
		Body:       body,
	}

	return &res, nil
}

// MustSend sends the HTTP request and panic if an error
// is returned. (Calls Send() internally)
func (req *Request) MustSend() *Response {
	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	return res
}

// Response is a type that represents an HTTP response
// returned from an HTTP request
type Response struct {
	Ok         bool              // Was the request successful? (Status codes: 200-399)
	StatusCode int               // HTTP response status code
	Headers    map[string]string // HTTP Response headers
	Body       []byte            // HTTP Response body
}

// GetHeader gets a header value from the response if it exists.
// Normalizes the key to lowercase before checking.
// Returns the value of the header and whether it exists.
func (res *Response) GetHeader(name string) (string, bool) {
	// Create the map if it doesn't exist
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	// Do a case-insensitive check for the key
	for k, v := range res.Headers {
		if strings.ToLower(k) == key {
			return v, true
		}
	}

	// Return not found
	return "", false
}

// JSON unmarshalls the response body into a map
func (res *Response) JSON() (map[string]interface{}, error) {
	// Create a new map to store the JSON data
	data := make(map[string]interface{})

	// Unmarshal the JSON data
	err := json.Unmarshal(res.Body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
