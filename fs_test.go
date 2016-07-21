package httpgzip_test

import (
	"net/http"

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
