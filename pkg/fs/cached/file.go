package cached

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"git.beryju.org/BeryJu.org/pixie/pkg/fs/standard"
)

// CachedFile Same as fs.File, but cache contents in memory
type CachedFile struct {
	standard.File
	Key  string
	FS   CachedFileSystem
	stat *Stat
}

// Serve Write cached contents to http response
func (cf CachedFile) Serve(w http.ResponseWriter, r *http.Request) {
	p, err := cf.FS.GetCacheFallback("", func() ([]byte, error) {
		cf.FS.Logger.Debug("File not in Cache, reading from disk")
		buffer, err := ioutil.ReadAll(cf.File)
		if err != nil {
			cf.FS.Logger.Warningf("ReadAll Err: %s", err)
			return nil, err
		}
		return buffer, nil
	})
	if err != nil {
		cf.FS.Logger.Warning(err)
	}
	re := bytes.NewReader(p)

	d, err := cf.Stat()
	if err != nil {
		cf.FS.Logger.Warning(err)
	} else {
		http.ServeContent(w, r, d.Name(), d.ModTime(), re)
	}
}

// Readdir Wrapper around fs.File.Readdir
func (cf CachedFile) Readdir(n int) (fis []os.FileInfo, err error) {
	return cf.File.Readdir(n)
}

// Stat Cached wrapper around os.File.Stat()
func (cf CachedFile) Stat() (os.FileInfo, error) {
	if cf.stat != nil {
		return cf.stat, nil
	}
	statByte, err := cf.FS.GetCacheFallback("stat-"+cf.Key, func() ([]byte, error) {
		stats, err := cf.File.Stat()
		if err != nil {
			cf.FS.Logger.Warningf("Stat Err: %s", err)
			return nil, err
		}
		marshalled, err := json.Marshal(Stat{
			name:    stats.Name(),
			modTime: stats.ModTime(),
			size:    stats.Size(),
			mode:    stats.Mode(),
		})
		if err != nil {
			cf.FS.Logger.Warningf("JSON Marshal: %s", err)
			return nil, err
		}
		return marshalled, nil
	})
	if err != nil {
		cf.FS.Logger.Debug(err)
	}
	err = json.Unmarshal(statByte, &cf.stat)
	if err != nil {
		cf.FS.Logger.Warning(err)
	}
	return cf.stat, nil
}
