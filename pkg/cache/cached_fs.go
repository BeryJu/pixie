package cache

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.beryju.org/BeryJu.org/pixie/pkg/base"
	"git.beryju.org/BeryJu.org/pixie/pkg/fs"
	"git.beryju.org/BeryJu.org/pixie/pkg/utils"
	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
)

// CachedFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type CachedFileSystem struct {
	base.FileSystem
	Dir    string
	Cache  *bigcache.BigCache
	Logger *log.Entry
}

// NewCachedFileSystem Initialise new Cached Filesystem
func NewCachedFileSystem() CachedFileSystem {
	cfs := CachedFileSystem{
		Dir:    fs.NewFileSystem().Dir,
		Logger: log.WithField("component", "cached-fs"),
	}
	cacheConfig := bigcache.DefaultConfig(10 * time.Minute)
	// cacheConfig.MaxEntrySize = config.CfgCacheMaxItemSize
	// cacheConfig.HardMaxCacheSize = config.CfgCacheMaxItemSize / 1024
	// cfs.Logger.Debug(cacheConfig.MaxEntrySize)
	// cfs.Logger.Debug(cacheConfig.HardMaxCacheSize)
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
func (cfs CachedFileSystem) Open(name string) (base.ServingFile, error) {
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
	f, err := os.Open(fullName)
	if err != nil {
		return nil, utils.MapDirOpenError(err, fullName)
	}
	return CachedFile{
		File: fs.File{
			File: f,
		},
		Key: fullName,
		FS:  cfs,
	}, nil
}
