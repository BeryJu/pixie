package cache

import (
	"bytes"
	"encoding/json"
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
	p, err := cf.readFileIfNeeded()
	if err != nil {
		cf.FS.Logger.Debug(err)
	}
	re := bytes.NewReader(p)
	statByte, err := cf.FS.Cache.Get("stat-" + cf.Key)
	if err == bigcache.ErrEntryNotFound {
		statByte, err = cf.readStatsIntoCache()
		if err != nil {
			cf.FS.Logger.Warning(err)
		}
	} else if err != nil {
		cf.FS.Logger.Warning(err)
	}
	var stats Stat
	err = json.Unmarshal(statByte, &stats)
	if err != nil {
		cf.FS.Logger.Warning(err)
	}
	http.ServeContent(w, r, stats.Name, stats.ModTime, re)
}

// Readdir Wrapper around fs.File.Readdir
func (cf CachedFile) Readdir(n int) (fis []os.FileInfo, err error) {
	return cf.File.Readdir(n)
}

// Stat Stat calls are not cached
func (cf CachedFile) Stat() (os.FileInfo, error) {
	return cf.File.Stat()
}

func (cf CachedFile) readFileIfNeeded() ([]byte, error) {
	data, err := cf.FS.Cache.Get(cf.Key)
	if err == bigcache.ErrEntryNotFound {
		cf.FS.Logger.Debug("File not in Cache, reading from disk")
		data, err := cf.readIntoCache()
		if err != nil {
			cf.FS.Logger.Warning(err)
		}
		// We also cache the stats of the file
		_, err = cf.readStatsIntoCache()
		if err != nil {
			cf.FS.Logger.Warning(err)
		}
		return data, nil
	} else if err != nil {
		cf.FS.Logger.Warningf("CacheGet: %s", err)
	}
	return data, nil
}

func (cf CachedFile) readIntoCache() ([]byte, error) {
	buffer, err := ioutil.ReadAll(cf.File)
	if err != nil {
		cf.FS.Logger.Warningf("ReadAll Err: %s", err)
		return nil, err
	}
	err = cf.FS.Cache.Set(cf.Key, buffer)
	if err != nil {
		cf.FS.Logger.Warningf("CacheSet: %s", err)
		return nil, err
	}
	return buffer, nil
}

func (cf CachedFile) readStatsIntoCache() ([]byte, error) {
	stats, err := cf.File.Stat()
	if err != nil {
		cf.FS.Logger.Warningf("Stat Err: %s", err)
		return nil, err
	}
	marshalled, err := json.Marshal(Stat{
		Name:    stats.Name(),
		ModTime: stats.ModTime(),
	})
	if err != nil {
		cf.FS.Logger.Warningf("JSON Marshal: %s", err)
		return marshalled, err
	}
	err = cf.FS.Cache.Set("stat-"+cf.Key, marshalled)
	if err != nil {
		cf.FS.Logger.Warningf("CacheSet: %s", err)
		return marshalled, err
	}
	return marshalled, nil
}
