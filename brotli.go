package httpgzip

import (
	"io"
)

func (fs *fileServer) findBrotliFile(fpath string) io.ReadSeeker {
	brotliFilePath := fpath + ".br"

	if file, err := fs.root.Open(brotliFilePath); err == nil {
		return file
	}

	return nil
}
