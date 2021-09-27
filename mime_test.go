package requests_test

import (
	"testing"

	"github.com/a-poor/requests"
)

func TestMIMETypes(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{".txt", "text/plain"},
		{".json", "application/json"},
		{".html", "text/html"},
		{".pdf", "application/pdf"},
	}

	for _, tc := range testCases {
		m, ok := requests.MIMETypes[tc.in]
		if !ok {
			t.Errorf("MIME type of %q not found", tc.in)
		}
		if m != tc.out {
			t.Errorf("MIME type of %q expected %q but got %q", tc.in, tc.out, m)
		}
	}

	m, ok := requests.MIMETypes[".foo"]
	if ok || m != "" {
		t.Errorf("MIME type of .foo should not be found, got %q", m)
	}

}

func TestGuessMIME(t *testing.T) {
	testCases := []struct {
		in  string
		ok  bool
		out string
	}{
		{"foo.txt", true, "text/plain"},
		{"my-file.pdf", true, "application/pdf"},
		{"path/to/my/file.json", true, "application/json"},
		{"foo.foo", false, ""},
		{"", false, ""},
		{"a.file.html", true, "text/html"},
	}

	for _, tc := range testCases {
		// Get the guessed MIME type
		m, ok := requests.GuessMIME(tc.in)

		// Check that the MIME type exists (if it should)
		if tc.ok && !ok {
			t.Errorf("MIME type of %q not found", tc.in)
		}

		// Or, check that it doesn't exist (if it shouldn't)
		if !tc.ok && ok {
			t.Errorf("MIME type of %q not found", tc.in)
		}

		// And check that the MIME type is correct
		if tc.ok && ok && m != tc.out {
			t.Errorf("MIME type of %q expected to be %q but got %q", tc.in, tc.out, m)
		}
	}
}

func TestGuessMIMEWithDefault(t *testing.T) {
	in := "foo.foo"
	expect := "application/octet-stream"
	m := requests.GuessMIMEWithDefault(in, requests.MIMEDefaultBinary)
	if m != expect {
		t.Errorf("MIME type of %q should be %q but got %q", in, expect, m)
	}

	in = "foo.bar"
	expect = "text/plain"
	m = requests.GuessMIMEWithDefault(in, requests.MIMEDefaultText)
	if m != expect {
		t.Errorf("MIME type of %q should be %q but got %q", in, expect, m)
	}

	in = "foo.html"
	expect = "text/html"
	m = requests.GuessMIMEWithDefault(in, requests.MIMEDefaultText)
	if m != expect {
		t.Errorf("MIME type of %q should be %q but got %q", in, expect, m)
	}
}
