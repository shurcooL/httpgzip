package httpgzip

import (
	"io"
	"net/http"
)

func (fs *fileServer) findBrotliFile(w http.ResponseWriter, r *http.Request, fpath string) io.ReadSeeker {
	brotliFilePath := fpath + ".br"

	if file, err := fs.root.Open(brotliFilePath); err == nil {
		defer file.Close()

		wHeader := w.Header()
		w.Header().Set("Content-Encoding", "br")
		wHeader.Add("Vary", r.Header.Get("Accept-Encoding"))

		return file
	}

	return nil
}
