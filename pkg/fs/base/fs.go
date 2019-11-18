package base

// FileSystem Base interface for FileSystem
type FileSystem interface {
	Open(name string) (ServingFile, error)
}
