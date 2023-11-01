package domain

import "time"

type File struct {
	Name string
	Path string
	Time time.Time
	Size uint64
}

func NewFile(
	name string,
	path string,
	time time.Time,
	size uint64,
) *File {
	return &File{
		Name: name,
		Path: path,
		Time: time,
		Size: size,
	}
}
