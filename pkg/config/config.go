package config

import "os"

// Config Config structure for viper marshalling
type Config struct {
	RootDir          string
	Port             int
	Debug            bool
	EXIFPurgeGPS     bool
	CacheEnabled     bool
	CacheEviction    int
	CacheMaxSize     int
	CacheMaxItemSize int
}

// Current Static variable, filled by viper
var Current Config

// Defaults default values
var Defaults Config = Config{
	RootDir:          getCwd(),
	Port:             8080,
	Debug:            false,
	EXIFPurgeGPS:     true,
	CacheEnabled:     false,
	CacheEviction:    10,
	CacheMaxSize:     0,
	CacheMaxItemSize: 500,
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
