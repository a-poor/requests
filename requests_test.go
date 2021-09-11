package requests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-poor/requests"
)

func TestHTTPMethods(t *testing.T) {
	if requests.GET.String() != "GET" {
		t.Error("requests.GET is not GET")
	}
	if requests.POST.String() != "POST" {
		t.Error("requests.POST is not POST")
	}
	if requests.PUT.String() != "PUT" {
		t.Error("requests.PUT is not GET")
	}
	if requests.DELETE.String() != "DELETE" {
		t.Error("requests.DELETE is not DELETE")
	}
	if requests.HEAD.String() != "HEAD" {
		t.Error("requests.HEAD is not HEAD")
	}
	if requests.OPTIONS.String() != "OPTIONS" {
		t.Error("requests.OPTIONS is not OPTIONS")
	}
	if requests.PATCH.String() != "PATCH" {
		t.Error("requests.PATCH is not PATCH")
	}
	if requests.CONNECT.String() != "CONNECT" {
		t.Error("requests.CONNECT is not CONNECT")
	}
	if requests.TRACE.String() != "TRACE" {
		t.Error("requests.TRACE is not TRACE")
	}
}

func TestRequestHeaders(t *testing.T) {
	r := &requests.Request{}
	if r == nil {
		t.Error("Request is nil")
	}

	r.SetHeader("Content-Type", "application/json")
	ct, ok := r.GetHeader("content-type")
	if !ok {
		t.Error("req Content-Type header not set")
	}
	if ct != "application/json" {
		t.Error("req.SetHeader is not working")
	}

	r.DelHeader("Content-Type")
	_, ok = r.GetHeader("content-type")
	if ok {
		t.Error("req.DelHeader is not working")
	}
}

func TestSendGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")

		if r.Method != "GET" {
			t.Errorf("Request method is \"%s\" not GET", r.Method)
		}

	}))
	defer ts.Close()

	res, err := requests.SendGetRequest(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Error("status code is not 200")
	}

	bod := string(res.Body)
	if bod != "Hello, World!\n" {
		t.Error(fmt.Sprintf("response body is \"%s\" not Hello, World!", bod))
	}

}

func ExampleRequest_Send() {
	r := requests.Request{
		Method: requests.GET,
		URL:    "http://example.com",
	}
	res, err := r.Send()
	if err != nil {
		// handle error
	}
	fmt.Println(res.StatusCode)
	// Output: 200
}

func ExampleSendGetRequest() {
	res, err := requests.SendGetRequest("http://example.com")
	if err != nil {
		// handle error
	}
	fmt.Println(res.StatusCode)
	// Output: 200
}
