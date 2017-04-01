package httpgzip_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shurcooL/httpgzip"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
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

// Test that there are no infinite redirects at root path even if the entire
// req.URL.Path is stripped, e.g., via an overly aggressive http.StripPrefix.
// See https://github.com/shurcooL/httpgzip/pull/3
// and https://github.com/golang/go/commit/3745716bc3940f471137bf06fbe8c042257a43d3.
func TestFileServer_implicitLeadingSlash(t *testing.T) {
	fs := httpfs.New(mapfs.New(map[string]string{
		"foo.txt": "Hello world",
	}))
	ts := httptest.NewServer(http.StripPrefix("/bar/", httpgzip.FileServer(fs, httpgzip.FileServerOptions{})))
	defer ts.Close()
	get := func(suffix string) string {
		res, err := http.Get(ts.URL + suffix)
		if err != nil {
			t.Fatalf("Get %s: %v", suffix, err)
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("ReadAll %s: %v", suffix, err)
		}
		return string(b)
	}
	if got := get("/bar/"); !strings.Contains(got, ">foo.txt<") {
		t.Errorf("got:\n%v\nwant a directory listing with foo.txt\n", got)
	}
	if got, want := get("/bar/foo.txt"), "Hello world"; got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}
