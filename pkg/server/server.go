package server

import (
	"fmt"
	"net/http"

	"git.beryju.org/BeryJu.org/pixie/internal"
	"git.beryju.org/BeryJu.org/pixie/pkg/base"
	"git.beryju.org/BeryJu.org/pixie/pkg/cache"
	"git.beryju.org/BeryJu.org/pixie/pkg/config"
	"git.beryju.org/BeryJu.org/pixie/pkg/fs"
	log "github.com/sirupsen/logrus"
)

// Server Holds all components needed for the main HTTP Server
type Server struct {
	Logger *log.Entry
	Mux    *http.ServeMux
}

// NewServer Initialise new HTTP-Server
func NewServer() *Server {
	logger := log.WithField("component", "http-server")
	var fsInstance base.FileSystem
	if config.Current.CacheEnabled {
		logger.Debug("Using cached filesystem.")
		fsInstance = cache.NewCachedFileSystem()
	} else {
		logger.Debug("Using normal filesystem.")
		fsInstance = fs.NewFileSystem()
	}
	mux := http.NewServeMux()
	mux.Handle("/", logging(logger)(internal.FileServer(fsInstance)))
	mux.HandleFunc("/-/ping", Ping)
	return &Server{
		Logger: logger,
		Mux:    mux,
	}
}

// Run Start HTTP Server
func Run() {
	server := NewServer()
	server.Logger.Infof("Serving '%s'", config.Current.RootDir)
	listen := fmt.Sprintf(":%d", config.Current.Port)
	if config.Current.Debug {
		listen = fmt.Sprintf("localhost:%d", config.Current.Port)
	}
	server.Logger.Fatal(http.ListenAndServe(listen, server.Mux))
}
