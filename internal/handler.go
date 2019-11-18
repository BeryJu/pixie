package internal

import (
	"net/http"
	"path"
	"strings"

	"git.beryju.org/BeryJu.org/pixie/pkg/api"
	"git.beryju.org/BeryJu.org/pixie/pkg/constant"
	"git.beryju.org/BeryJu.org/pixie/pkg/fs/base"
	"git.beryju.org/BeryJu.org/pixie/pkg/templates"
	log "github.com/sirupsen/logrus"
)

// FileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at root.
//
// To use the operating system's file system implementation,
// use http.Dir:
//
//     http.Handle("/", http.FileServer(http.Dir("/tmp")))
//
// As a special case, the returned file server redirects any request
// ending in "/index.html" to the same path, without the final
// "index.html".
func FileServer(root base.FileSystem) http.Handler {
	return &fileHandler{root}
}

type fileHandler struct {
	root base.FileSystem
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	serveFile(w, r, f.root, path.Clean(upath), true)
}

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs base.FileSystem, name string, redirect bool) {
	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(r.URL.Path, constant.IndexPageFileName) {
		localRedirect(w, r, "./")
		return
	}

	f, err := fs.Open(name)
	if err != nil {
		log.Debug(err)
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		log.Debug(err)
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	// redirect if the directory name doesn't end in a slash
	if d.IsDir() {
		url := r.URL.Path
		if url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}
	}

	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + constant.IndexPageFileName
		ff, err := fs.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = index
				d = dd
				f = ff
			}
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		if checkIfModifiedSince(r, d.ModTime()) == condFalse {
			writeNotModified(w)
			return
		}
		w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
		if _, ok := r.URL.Query()["json"]; ok {
			api.HandleDirList(w, r, f)
		} else {
			galleryTemplate, err := templates.GetTemplate("gallery.html")
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
			}
			stat, err := f.Stat()
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
			}
			galleryTemplate.ServeTemplate(templates.GalleryTemplateContext{
				RelativePath: stat.Name(),
			}, w, r)
		}
		return
	}

	// We're sure it's not a directory, check for the ?meta tag
	if _, ok := r.URL.Query()["meta"]; ok {
		api.HandleFileMeta(w, r, f)
		return
	}

	f.Serve(w, r)
}
