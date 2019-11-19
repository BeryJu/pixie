package api

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FileMeta struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

func GetFileContentType(f http.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func (fm *FileMeta) ForFile(f http.File) error {
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	fm.ContentType = ""
	if stat.Size() >= 512 {
		fm.ContentType, err = GetFileContentType(f)
		f.Seek(0, 0)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Skipping content-type check for %s because file is more than 512 bytes", stat.Name())
	}
	fm.Name = stat.Name()
	fm.Size = stat.Size()
	return nil
}

func HandleFileMeta(w http.ResponseWriter, r *http.Request, f http.File) {
	meta := FileMeta{}
	err := meta.ForFile(f)
	if err != nil {
		log.Error(err)
	}
	if strings.HasPrefix(meta.ContentType, "image/") {
		imageMeta := ImageFileMeta{}
		err = imageMeta.ForFile(f)
		json.NewEncoder(w).Encode(imageMeta)
	} else {
		json.NewEncoder(w).Encode(meta)
	}
}
