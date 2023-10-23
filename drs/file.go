package drs

import (
	"fmt"
	"io"
)

type FileInfo struct {
	ID     FileID
	Offset int32
	Size   int32
}

type File struct {
	FileInfo

	drs *Reader
	rd  io.ReaderAt

	table *Table
}

func (f *File) Table() *Table {
	return f.table
}

func (file *File) Open() (*io.SectionReader, error) {
	rd := io.NewSectionReader(file.rd, int64(file.Offset), int64(file.Size))
	return rd, nil
}

func (file *File) Data() ([]byte, error) {
	buf := make([]byte, file.Size)
	n, err := file.rd.ReadAt(buf, int64(file.Offset))
	if err != nil {
		return nil, err
	}
	if n != int(file.Size) {
		return nil, fmt.Errorf("not enough data! want=%d, got=%d", file.Size, n)
	}
	return buf, nil
}
