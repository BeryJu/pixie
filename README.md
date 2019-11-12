# pixie

```
Usage:
  pixie [flags]

Flags:
      --debug             Enable debug-mode.
  -h, --help              help for pixie
      --purge-exif-gps    Purge GPS-Relateed EXIF metadata. (default true)
  -r, --root-dir string   Root directory to serve. (default ".")
```

Demo: https://i.beryju.org/pixie-demo/

## Running

### Docker

Run the container like this:

```
docker run -v "whatever directory you want to share":/data docker.beryju.org/pixie/server:latest -r=/data
```

Now you can access pixie on http://localhost:8080

### Binary

Download a binary from [GitLab](https://git.beryju.org/BeryJu.org/pixie/pipelines) and run it:

```
./pixie -r=/data
```

Now you can access pixie on http://localhost:8080

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.
