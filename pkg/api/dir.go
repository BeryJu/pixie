package api

import (
	"encoding/json"
	"net/http"
	"sort"
)

type DirectoryFile struct {
	Name string `json:"name"`
}

type DirectoryListing struct {
	Files []DirectoryFile `json:"files"`
}

func HandleDirList(w http.ResponseWriter, r *http.Request, f http.File) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	listing := DirectoryListing{}
	listing.Files = make([]DirectoryFile, len(dirs))
	for i, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		// url := url.URL{Path: name}
		listing.Files[i] = DirectoryFile{
			Name: name,
		}
	}
	json.NewEncoder(w).Encode(listing)
}
