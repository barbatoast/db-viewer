package dbviewer

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// NeuteredFileSystem is used to prevent directory listing of static assets
type NeuteredFileSystem struct {
	Fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	var err error

	// check if path exists
	f, err := nfs.Fs.Open(path)
	if err != nil {
		logrus.Error("Fs.Open ", err)
		return nil, err
	}

	// If path exists, check if it is a file or a directory.
	// If it is a directory, stop here with an error saying that file
	// does not exist. User will get a 404 error code for a file/directory
	// that does not exist, and for directories that exist.
	s, err := f.Stat()
	if err != nil {
		logrus.Error("f.Stat() ", err)
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	// if file exists and the path is not a directory, return the file
	return f, nil
}
