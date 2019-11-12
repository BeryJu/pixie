package api

import (
	"net/http"
	"strings"

	"git.beryju.org/BeryJu.org/pixie/pkg/config"

	log "github.com/sirupsen/logrus"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type ImageFileMeta struct {
	FileMeta
	Exif map[string]string `json:"exif"`
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
	if config.CfgPurgeExifGPS && !strings.HasPrefix(stringName, "GPS") {
		w.meta.Exif[stringName] = value
	}
	return nil
}

func (fm *ImageFileMeta) ForFile(f http.File) error {
	err := fm.FileMeta.ForFile(f)
	if err != nil {
		return err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return err
	}
	fm.Exif = make(map[string]string)
	walker := walker{
		file: f,
		meta: fm,
	}
	err = x.Walk(walker)
	return err
}
