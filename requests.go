package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
type Request struct {
	URL     string            // URL to send the request to
	Method  HTTPMethod        // HTTP method to use
	Headers map[string]string // Headers to send with the request
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

// getReqBody returns the request body as a buffer that can be
// passed to the http.NewRequest function
func (req *Request) getReqBody() *bytes.Buffer {
	return bytes.NewBuffer(req.Body)
}

// GetHeader gets a header value from the request. Normalizes the key
// to lowercase before checking. Returns the value of the
// header and whether it exists.
func (req *Request) GetHeader(name string) (string, bool) {
	// Create the map if it doesn't exist
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	// Return the header and whether it exists
	value, ok := req.Headers[key]
	return value, ok
}

// SetHeader sets a header value in the request. Normalizes the key
// before setting (converts to lowercase).
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

	// Delete the header if it exists
	delete(req.Headers, key)
}

// Send sends the HTTP request with the supplied parameters
func (req *Request) Send() (*Response, error) {
	// Create an http client (with optional timeout)
	client := http.Client{
		Timeout: req.Timeout,
	}

	// Create the underlying request
	httpRequest, err := http.NewRequest(req.Method.String(), req.URL, req.getReqBody())
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
		lowerKey := strings.ToLower(k)
		rHeaders[lowerKey] = v[0]
	}

	// Load the request body
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	// Format the response & return
	resp := Response{
		Ok:         httpResponse.StatusCode < 400,
		StatusCode: httpResponse.StatusCode,
		Headers:    rHeaders,
		Body:       body,
	}

	return &resp, nil
}

// MustSend sends the HTTP request and panic if an error is returned.
// (Calls Send() internally)
func (req *Request) MustSend() *Response {
	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	return resp
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
func (resp *Response) GetHeader(name string) (string, bool) {
	// Create the map if it doesn't exist
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}

	// Normalize the key (convert to lowercase)
	key := strings.ToLower(name)

	// Return the header and whether it exists
	value, ok := resp.Headers[key]
	return value, ok
}

// JSON unmarshalls the response body into a map
func (resp *Response) JSON() (map[string]interface{}, error) {
	// Create a new map to store the JSON data
	data := make(map[string]interface{})

	// Unmarshal the JSON data
	err := json.Unmarshal(resp.Body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
