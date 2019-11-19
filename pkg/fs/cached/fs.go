package cached

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.beryju.org/BeryJu.org/pixie/pkg/config"
	"git.beryju.org/BeryJu.org/pixie/pkg/fs/base"
	"git.beryju.org/BeryJu.org/pixie/pkg/fs/standard"
	"git.beryju.org/BeryJu.org/pixie/pkg/utils"
	"github.com/allegro/bigcache"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// FileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type FileSystem struct {
	base.FileSystem
	Dir    string
	Cache  *bigcache.BigCache
	Logger *log.Entry
}

// NewFileSystem Initialise new Cached Filesystem
func NewFileSystem() FileSystem {
	cfs := FileSystem{
		Dir:    standard.NewFileSystem().Dir,
		Logger: log.WithField("component", "cached-fs"),
	}
	cacheConfig := bigcache.DefaultConfig(time.Duration(config.Current.CacheEviction) * time.Minute)
	cacheConfig.MaxEntrySize = config.Current.CacheMaxItemSize
	cacheConfig.HardMaxCacheSize = config.Current.CacheMaxSize
	cacheConfig.Logger = cfs.Logger
	cache, err := bigcache.NewBigCache(cacheConfig)
	if err != nil {
		cfs.Logger.Warning(err)
	}
	cfs.Cache = cache
	return cfs
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (cfs FileSystem) Open(name string) (base.ServingFile, error) {
	if utils.ContainsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(cfs.Dir)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	return File{
		Key: fullName,
		FS:  cfs,
	}, nil
}

// GetCacheFallback Wrapper around bigcache.Get and Set
func (cfs FileSystem) GetCacheFallback(key string, populate func() ([]byte, error)) ([]byte, error) {
	ret, err := cfs.Cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		cfs.Logger.WithField("key", key).Debug("Entry not found, calling populate()")
		realValue, err := populate()
		if err != nil {
			return nil, errors.Wrap(err, "Error executing CacheGet populate")
		}
		err = cfs.Cache.Set(key, realValue)
		if err != nil {
			return nil, errors.Wrap(err, "Error during CacheSet")
		}
		return realValue, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Error during CacheGet")
	}
	return ret, nil
}
