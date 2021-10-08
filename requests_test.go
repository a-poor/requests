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

func TestURLEncode(t *testing.T) {
	res := requests.URLEncode("foo")
	expect := "foo"
	if res != expect {
		t.Errorf("URLEncode expected %q not %q", expect, res)
	}

	res = requests.URLEncode("Hello, World!")
	expect = "Hello%2C%20World%21"
	if res != expect {
		t.Errorf("URLEncode expected %q not %q", expect, res)
	}

	res = requests.URLEncode(123)
	expect = "123"
	if res != expect {
		t.Errorf("URLEncode expected %q not %q", expect, res)
	}

	res = requests.URLEncode("1/2")
	expect = "1%2F2"
	if res != expect {
		t.Errorf("URLEncode expected %q not %q", expect, res)
	}
}

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

func TestQueryParams(t *testing.T) {
	params := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for k, v := range params {
			if q.Get(k) != v {
				t.Errorf("Query param %s is \"%s\" not \"%s\"", k, q.Get(k), v)
			}
		}
	}))
	defer ts.Close()

	req := requests.Request{
		Method: requests.GET,
		URL:    ts.URL,
	}
	for k, v := range params {
		req.SetQuery(k, v)
	}
	_, err := req.Send()
	if err != nil {
		t.Error(err)
	}
}

func TestRequestCopy(t *testing.T) {
	r1 := requests.Request{
		Method: requests.GET,
		URL:    "http://example.com",
		Body:   []byte("Hello, World!"),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	r2 := r1.Copy()

	// Are the pointers the same?
	if &r1 == r2 {
		t.Error("request pointer doesn't change after copy")
	}

	// Make sure the header maps are pointing to different maps
	r1.Headers["Content-Type"] = "text/plain"
	if r2.Headers["Content-Type"] == "text/plain" {
		t.Error("coppied request header map not coppied")
	}
}

func TestRequestPathParse(t *testing.T) {
	templatePath := `http://example.com/{{ .UserID }}/{{ .Text | URLEncode }}`
	data := struct {
		UserID int
		Text   string
	}{
		UserID: 123,
		Text:   "Hello, World!",
	}
	resultPath := `http://example.com/123/Hello%2C%20World%21`

	req := &requests.Request{
		URL: templatePath,
	}
	req, err := req.ParsePathParams(data)
	if err != nil {
		t.Error(err)
	}
	if req.URL != resultPath {
		t.Errorf("unexpected path formatting %q not %q", req.URL, resultPath)
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
	resp := requests.Response{
		Ok:         true,
		StatusCode: 200,
		Body:       []byte(`{"message":"pong"}`),
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := resp.JSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSend(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			b.Errorf("Request method is \"%s\" not GET", r.Method)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := requests.SendGetRequest(ts.URL)
		if err != nil {
			b.Fatal(err)
		}
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

func TestSendGetRequestWith500InternalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method is \"%s\" not GET", r.Method)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	res, err := requests.SendGetRequest(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code is not %d", http.StatusInternalServerError)
	}
}

func TestSendGetRequestWith500InternalWithJsonBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method is \"%s\" not GET", r.Method)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": true}`))
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	res, err := requests.SendGetRequest(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code is not %d", http.StatusInternalServerError)
	}

	jresp, err := res.JSON()
	if err != nil {
		t.Error(err)
	}

	if jresp["error"] != true {
		t.Error("expected `error` field to be set in json response")
	}
}

func TestSendGetRequestWithJsonResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method is \"%s\" not GET", r.Method)
			return
		}

		_, err := w.Write([]byte(`
{
	"a": "test",
	"b": true, 
	"c": 3
}
`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	res, err := requests.SendGetRequest(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("status code is not 200")
	}

	jresp, err := res.JSON()
	if err != nil {
		t.Error(err)
	}

	if jresp["a"] != "test" {
		t.Errorf("expected 'test' but was %v", jresp["a"])
	}

	if jresp["b"] != true {
		t.Errorf("expected 'true' but was %v", jresp["b"])
	}

	if jresp["c"] != 3.0 {
		t.Errorf("expected '3' but was %v", jresp["c"])
	}

	if jresp["_"] != nil {
		t.Errorf("expected 'nil' but was %v", jresp["_"])
	}

}

func TestSendGetRequestWithEmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method is \"%s\" not GET", r.Method)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	res, err := requests.SendGetRequest(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("status code is not 200")
	}

	jresp, err := res.JSON()
	if err != nil {
		t.Error(err)
	}

	if len(jresp) != 0 {
		t.Fatal("expected empty json on empty body")
	}
}
