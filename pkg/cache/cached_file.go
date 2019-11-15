package cache

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/allegro/bigcache"

	"git.beryju.org/BeryJu.org/pixie/pkg/fs"
)

// CachedFile Same as fs.File, but cache contents in memory
type CachedFile struct {
	fs.File
	Key string
	FS  CachedFileSystem
}

// Serve Write cached contents to http response
func (cf CachedFile) Serve(w http.ResponseWriter, r *http.Request) {
	err := cf.readFileIfNeeded()
	if err != nil {
		cf.FS.Logger.Debug(err)
	}
	p, err := cf.FS.Cache.Get(cf.Key)
	if err != nil {
		cf.FS.Logger.Debug(err)
	}
	w.Write(p)
}

// Readdir Wrapper around fs.File.Readdir
func (cf CachedFile) Readdir(n int) (fis []os.FileInfo, err error) {
	return cf.File.Readdir(n)
}

// Stat Stat calls are not cached
func (cf CachedFile) Stat() (os.FileInfo, error) {
	return cf.File.Stat()
}

func (cf CachedFile) readFileIfNeeded() error {
	_, err := cf.FS.Cache.Get(cf.Key)
	if err == bigcache.ErrEntryNotFound {
		cf.FS.Logger.Debug("File not in Cache, reading from disk")
		buffer, err := ioutil.ReadAll(cf.File)
		if err != nil {
			cf.FS.Logger.Warningf("ReadAll Err: %s", err)
			return err
		}
		err = cf.FS.Cache.Set(cf.Key, buffer)
		if err != nil {
			cf.FS.Logger.Warningf("CacheSet: %s", err)
			return err
		}
	} else if err != nil {
		cf.FS.Logger.Warningf("CacheGet: %s", err)
	}
	return nil
}
