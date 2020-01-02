package server

import (
	"fmt"
	"net/http"

	"github.com/BeryJu/pixie/internal"
	"github.com/BeryJu/pixie/pkg/config"
	"github.com/BeryJu/pixie/pkg/fs/base"
	"github.com/BeryJu/pixie/pkg/fs/cached"
	"github.com/BeryJu/pixie/pkg/fs/standard"
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
		fsInstance = cached.NewFileSystem()
	} else {
		logger.Debug("Using normal filesystem.")
		fsInstance = standard.NewFileSystem()
	}
	mux := http.NewServeMux()
	if config.Current.Silent {
		mux.Handle("/", internal.FileServer(fsInstance))
	} else {
		mux.Handle("/", logging(logger)(internal.FileServer(fsInstance)))
	}
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
