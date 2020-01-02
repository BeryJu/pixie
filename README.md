# pixie

```
Usage:
  pixie [directory to serve] [flags]

Flags:
      --cache-enabled             Enable in-memory cache
      --cache-eviction int        Time after which entry can be evicted (in minutes) (default 10)
      --cache-max-item-size int   Maximum Item size to cache (in bytes) (default 500)
      --cache-max-size int        Maximum Cache size in MB (0 disables the limit)
      --debug                     Enable debug-mode
      --exif-purge-gps            Purge GPS-Related EXIF metadata (default true)
  -h, --help                      help for pixie
      --silent                    Enable silent mode (no access logs)
```

Demo: https://i.beryju.org/pixie-demo/

## Running

### Docker

Run the container like this:

```
docker run -v "whatever directory you want to share":/data -w /data beryju/pixie:latest-amd64
```

Now you can access pixie on http://localhost:8080

### Binary

Download a binary from [GitHub](https://github.com/BeryJu/pixie/releases) and run it:

```
./pixie /data
```

Now you can access pixie on http://localhost:8080

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.
