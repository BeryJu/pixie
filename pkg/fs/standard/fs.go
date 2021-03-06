package standard

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BeryJu/pixie/pkg/config"
	"github.com/BeryJu/pixie/pkg/fs/base"
	"github.com/BeryJu/pixie/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// FileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type FileSystem struct {
	base.FileSystem
	Dir    string
	Logger *log.Entry
}

// NewFileSystem Initialise new FileSystem
func NewFileSystem() FileSystem {
	return FileSystem{
		Dir:    config.Current.RootDir,
		Logger: log.WithField("component", "fs"),
	}
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fs FileSystem) Open(name string) (base.ServingFile, error) {
	if utils.ContainsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(fs.Dir)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, utils.MapDirOpenError(err, fullName)
	}
	return File{f}, nil
}
