// Package files opens the SVG file to check on the local file system.
package files

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

// ErrFolder means the name given is a folder, not a file.
var ErrFolder = errors.New("it is a folder, not a file: give the SVG file inside it")

// Open opens a file for reading. Errors leave out the file name, which
// the caller already shows, and a folder is refused up front.
func Open(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if pe := (*fs.PathError)(nil); errors.As(err, &pe) {
		err = pe.Err
	}
	if err != nil {
		return nil, err
	}
	// If Stat itself fails, reading will fail too and say why.
	if info, err := f.Stat(); err == nil && info.IsDir() {
		f.Close()
		return nil, ErrFolder
	}
	return f, nil
}
