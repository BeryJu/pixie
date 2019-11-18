# pixie

```
Usage:
  pixie [flags]

Flags:
      --cache-enabled             Enable in-memory cache
      --cache-eviction int        Time after which entry can be evicted (in minutes) (default 10)
      --cache-max-item-size int   Maximum Item size to cache (in bytes) (default 500)
      --cache-max-size int        Maximum Cache size in MB (0 disables the limit)
      --debug                     Enable debug-mode
      --exif-purge-gps            Purge GPS-Related EXIF metadata (default true)
  -h, --help                      help for pixie
  -r, --root-dir string           Root directory to serve (default "`cwd`")
      --silent                    Enable silent mode (no access logs)
```

Demo: https://i.beryju.org/pixie-demo/

## Running

### Docker

Run the container like this:

```
docker run -v "whatever directory you want to share":/data docker.beryju.org/pixie/server:latest -r /data
```

Now you can access pixie on http://localhost:8080

### Binary

Download a binary from [GitLab](https://git.beryju.org/BeryJu.org/pixie/pipelines) and run it:

```
./pixie -r /data
```

Now you can access pixie on http://localhost:8080

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.
