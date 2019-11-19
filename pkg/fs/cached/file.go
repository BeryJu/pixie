package cached

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"git.beryju.org/BeryJu.org/pixie/pkg/fs/standard"
)

// File Same as fs.File, but cache contents in memory
type File struct {
	*standard.File
	Key  string
	FS   FileSystem
	stat *Stat
}

// Serve Write cached contents to http response
func (cf File) Serve(w http.ResponseWriter, r *http.Request) {
	p, err := cf.FS.GetCacheFallback(cf.Key, func() ([]byte, error) {
		cf.FS.Logger.WithField("name", cf.Key).Debug("File not in Cache, reading from disk")
		buffer, err := ioutil.ReadAll(cf.openFileIfNeeded())
		if err != nil {
			return nil, errors.Wrap(err, "ReadAll Error")
		}
		return buffer, nil
	})
	if err != nil {
		cf.FS.Logger.Warning(err)
	}
	re := bytes.NewReader(p)

	d, err := cf.Stat()
	if err != nil {
		// http.Error
		cf.FS.Logger.Warning(err)
	} else {
		http.ServeContent(w, r, d.Name(), d.ModTime(), re)
	}
}

// Readdir Wrapper around fs.File.Readdir
func (cf File) Readdir(n int) (fis []os.FileInfo, err error) {
	return cf.openFileIfNeeded().Readdir(n)
}

// Close Close actual file if we needed it
func (cf File) Close() error {
	if cf.File != nil {
		return cf.File.Close()
	}
	return nil
}

// Stat Cached wrapper around os.File.Stat()
func (cf File) Stat() (os.FileInfo, error) {
	if cf.stat != nil {
		return cf.stat, nil
	}
	statByte, err := cf.FS.GetCacheFallback("stat-"+cf.Key, func() ([]byte, error) {
		stats, err := cf.openFileIfNeeded().Stat()
		if err != nil {
			return nil, errors.Wrap(err, "Stat Error (reading into cache)")
		}
		stat := &Stat{
			NameField:    stats.Name(),
			ModTimeField: stats.ModTime(),
			SizeField:    stats.Size(),
			ModeField:    stats.Mode(),
		}
		marshalled, err := json.Marshal(stat)
		if err != nil {
			return nil, errors.Wrap(err, "JSON marshalling into cache")
		}
		return marshalled, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to populate Cache for Stat")
	}
	err = json.Unmarshal(statByte, &cf.stat)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal JSON")
	}
	return cf.stat, nil
}

func (cf File) openFileIfNeeded() *standard.File {
	if cf.File == nil {
		f, err := os.Open(cf.Key)
		if err != nil {
			cf.FS.Logger.Warning(err)
		}
		cf.File = &standard.File{
			File: f,
		}
	}
	return cf.File
}
