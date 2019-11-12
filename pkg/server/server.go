package server

import (
	"net/http"

	"git.beryju.org/BeryJu.org/pixie/internal"
	"git.beryju.org/BeryJu.org/pixie/pkg/abstract"
	"git.beryju.org/BeryJu.org/pixie/pkg/config"

	log "github.com/sirupsen/logrus"
)

func Run() {
	fs := abstract.CFileSystem{http.Dir(config.CfgRootDir)}
	log := log.WithField("component", "http-server")
	log.Infof("Serving '%s'", config.CfgRootDir)
	http.Handle("/", logging(log)(internal.FileServer(fs)))
	http.HandleFunc("/-/ping", Ping)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
