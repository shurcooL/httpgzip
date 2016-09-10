package httpgzip_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shurcooL/httpgzip"
)

func ExampleFileServer() {
	http.Handle("/assets/", http.StripPrefix("/assets", httpgzip.FileServer(
		http.Dir("assets"),
		httpgzip.FileServerOptions{
			IndexHTML: true,
		},
	)))
}

// ServeContent should correctly determine the content type as "text/plain",
// not as "application/x-gzip".
func TestServeContent_detectContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		content := "This is some plain text that compresses easily. " +
			strings.Repeat("NaN", 16) + " Batman!"

		httpgzip.ServeContent(w, req, "", time.Time{}, strings.NewReader(content))
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	got := resp.Header.Get("Content-Type")
	want := "text/plain; charset=utf-8"
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}
