package config

// CfgRootDir Root directory to serve
var CfgRootDir string

// CfgPurgeExifGPS Purge GPS-related EXIF tags from images
var CfgPurgeExifGPS bool

// CfgDebug Enable debug mode (verbose logging, etc)
var CfgDebug bool

// CfgCacheEnabled Enable in-memory cache
var CfgCacheEnabled bool

// CfgCacheMaxItems Maximum Items to cache
var CfgCacheMaxItems int

// CfgCacheMaxItemSize Maximum Item size to be cached (bytes)
var CfgCacheMaxItemSize int
