package config

// Config Config structure for viper marshalling
type Config struct {
	RootDir          string
	Port             int
	Debug            bool
	EXIFPurgeGPS     bool
	CacheEnabled     bool
	CacheMaxItems    int
	CacheMaxItemSize int
}

// Current Static variable, filled by viper
var Current Config

// Defaults default values
var Defaults Config = Config{
	RootDir:          ".",
	Port:             8080,
	Debug:            false,
	EXIFPurgeGPS:     true,
	CacheEnabled:     false,
	CacheMaxItems:    2000,
	CacheMaxItemSize: 1024,
}
