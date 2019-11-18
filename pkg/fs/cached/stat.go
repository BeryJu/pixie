package cached

import (
	"os"
	"time"
)

type Stat struct {
	NameField    string
	SizeField    int64
	ModTimeField time.Time
	ModeField    os.FileMode
}

func (s Stat) Name() string {
	return s.NameField
}

func (s Stat) Size() int64 {
	return s.SizeField
}

func (s Stat) Mode() os.FileMode {
	return s.ModeField
}

func (s Stat) ModTime() time.Time {
	return s.ModTimeField
}

func (s Stat) IsDir() bool {
	return false
}

func (s Stat) Sys() interface{} {
	return nil
}
