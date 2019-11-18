package cached

import (
	"os"
	"time"
)

type Stat struct {
	name    string
	size    int64
	modTime time.Time
	mode    os.FileMode
}

func (s Stat) Name() string {
	return s.name
}

func (s Stat) Size() int64 {
	return s.size
}

func (s Stat) Mode() os.FileMode {
	return s.mode
}

func (s Stat) ModTime() time.Time {
	return s.modTime
}

func (s Stat) IsDir() bool {
	return false
}

func (s Stat) Sys() interface{} {
	return nil
}
