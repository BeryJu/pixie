package fs

import (
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// File is the http.File use in FileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type File struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f File) Readdir(n int) (fis []os.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// Serve Calls http.ServeContent
func (f File) Serve(w http.ResponseWriter, r *http.Request) {
	d, err := f.Stat()
	if err != nil {
		log.Debug(err)
	} else {
		http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	}
}
