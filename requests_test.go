package requests_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-poor/requests"
)

func TestJSONMust(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	data := map[string]interface{}{"msg": "ping", "nested": map[string]interface{}{"msg": "pong"}}
	_ = requests.JSONMust(data)
}

func TestJSONMustPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("JSONMust should panic but doesn't")
		}
	}()

	data := map[string]interface{}{"msg": make(chan int)}
	_ = requests.JSONMust(data)
}

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

func TestSendPostRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure the request is a POST
		if r.Method != "POST" {
			t.Errorf("Request method is \"%s\" not POST", r.Method)
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		defer r.Body.Close()

		// Unmarshal the request body
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			t.Error(err)
		}

		// Make sure the request body is correct
		msg, ok := data["message"]
		if !ok {
			t.Error("message not found in request body")
		}
		if msg != "ping" {
			t.Errorf("message is \"%s\" not \"ping\"", msg)
		}

		// Write the response message
		respBody := map[string]string{"message": "pong"}
		rdata, _ := json.Marshal(respBody)
		fmt.Fprintln(w, string(rdata))
		r.Header.Set("Content-Type", "application/json")

	}))
	defer ts.Close()

	// Create a POST request
	data, _ := json.Marshal(map[string]string{"message": "ping"})
	res, err := requests.SendPostRequest(ts.URL, "application/json", data)
	if err != nil {
		t.Error(err)
	}

	// Check the return status code
	if res.StatusCode != 200 {
		t.Error("status code is not 200")
	}

	// Check the response body
	respData := make(map[string]string)
	json.Unmarshal(res.Body, &respData)
	msg, ok := respData["message"]
	if !ok {
		t.Error("message not found in response body")
	}
	if msg != "pong" {
		t.Error(fmt.Sprintf("response body is \"%s\" not Hello, World!", msg))
	}

}

func TestResponseJSON(t *testing.T) {
	resp := requests.Response{
		Ok:         true,
		StatusCode: 200,
		Body:       []byte(`{"message":"pong"}`),
	}
	dat, err := resp.JSON()
	if err != nil {
		t.Error(err)
	}
	msg, ok := dat["message"]
	if !ok {
		t.Error("message not present in response")
	}
	if smsg := msg.(string); smsg != "pong" {
		t.Errorf("response message equals \"%s\", not \"pong\"", smsg)
	}
}

func BenchmarkResponseJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		resp := requests.Response{
			Ok:         true,
			StatusCode: 200,
			Body:       []byte(`{"message":"pong"}`),
		}
		resp.JSON()
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
