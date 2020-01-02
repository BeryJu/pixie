package api

import (
	"net/http"
	"strings"

	"github.com/BeryJu/pixie/pkg/config"

	log "github.com/sirupsen/logrus"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

// ImageFileMeta Extends FileMeta with a [string]string map for EXIF attributes
type ImageFileMeta struct {
	FileMeta
	EXIF map[string]string `json:"exif"`
}

type walker struct {
	file http.File
	meta *ImageFileMeta
}

func (w walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	value, err := tag.StringVal()
	if err != nil {
		log.Debug(err)
	}
	stringName := string(name)
	if config.Current.EXIFPurgeGPS && !strings.HasPrefix(stringName, "GPS") {
		w.meta.EXIF[stringName] = value
	}
	return nil
}

// ForFile get meta data for file
func (fm *ImageFileMeta) ForFile(f http.File) error {
	err := fm.FileMeta.ForFile(f)
	if err != nil {
		return err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return err
	}
	fm.EXIF = make(map[string]string)
	walker := walker{
		file: f,
		meta: fm,
	}
	err = x.Walk(walker)
	return err
}
