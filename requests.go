package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	PUT
	DELETE
	OPTIONS
	HEAD
	CONNECT
	TRACE
	PATCH
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

type Request struct {
	URL     string            // URL to send the request to
	Method  HTTPMethod        // HTTP method to use
	Headers map[string]string // Headers to send with the request
	Body    []byte            // Body to send with the request
	Timeout time.Duration     // Timeout for the request
}

// Return the request body as a buffer that can be
// passed to the http.NewRequest function
func (req *Request) getReqBody() *bytes.Buffer {
	return bytes.NewBuffer(req.Body)
}

// Get a header value from the request. Normalizes the key
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

// Sets a header value in the request. Normalizes the key
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

// Deletes a header value from the request headers
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

// Send the HTTP request with the supplied parameters
func (req *Request) Send() (*Response, error) {
	// Create an http client (with optional timeout)
	client := http.Client{
		Timeout: req.Timeout,
	}

	// Create the underlying request
	httpRequest, err := http.NewRequest(req.Method.String(), req.URL, req.getReqBody())
	if err != nil {
		return nil, err
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

// Send the HTTP request and panic if an error is returned.
// (Calls Send() internally)
func (req *Request) MustSend() *Response {
	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	return resp
}

type Response struct {
	Ok         bool              // Was the request successful? (Status codes: 200-399)
	StatusCode int               // HTTP response status code
	Headers    map[string]string // HTTP Response headers
	Body       []byte            // HTTP Response body
}

// Get a header value from the response if it exists.
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

// Unmarshall the response body into a map
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
