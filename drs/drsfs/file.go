package drsfs

import (
	"io"
	"io/fs"
	"time"

	"gopkg.in/KlemensWinter/go-genie.v1/drs"
)

type (
	file struct {
		fh *drs.File

		rd *io.SectionReader
	}
)

var (
	_ fs.ReadDirFile = (*rootFile)(nil)
)

func NewFile(fh *drs.File) (f *file, err error) {
	f = &file{
		fh: fh,
	}
	f.rd, err = fh.Open()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *file) Size() int64 {
	return f.rd.Size()
}
func (f *file) Sys() any {
	return f.fh
}

func (f *file) Stat() (fs.FileInfo, error) {
	ret := &fileInfo{
		name: f.fh.ID.String(),
		size: int64(f.fh.Size),
		// mode: fs.Mode
	}
	return ret, nil
}

func (f *file) Read(p []byte) (int, error) {
	return f.rd.Read(p)
}

func (f *file) ReadAt(p []byte, off int64) (int, error) {
	return f.rd.ReadAt(p, off)
}

func (f *file) Close() error { return nil }

// this struct represents the root entry.
type rootFile struct {
	fsys *FS
}

func (rf *rootFile) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name: ".",
		mode: fs.ModeDir,
	}, nil
}
func (rf *rootFile) Read([]byte) (int, error) { return 0, errUnreadable }
func (rf *rootFile) Close() error             { return nil }

// ReadDir reads the contents of the directory and returns
// a slice of up to n DirEntry values in directory order.
// Subsequent calls on the same file will yield further DirEntry values.
//
// If n > 0, ReadDir returns at most n DirEntry structures.
// In this case, if ReadDir returns an empty slice, it will return
// a non-nil error explaining why.
// At the end of a directory, the error is io.EOF.
// (ReadDir must return io.EOF itself, not an error wrapping io.EOF.)
//
// If n <= 0, ReadDir returns all the DirEntry values from the directory
// in a single slice. In this case, if ReadDir succeeds (reads all the way
// to the end of the directory), it returns the slice and a nil error.
// If it encounters an error before the end of the directory,
// ReadDir returns the DirEntry list read until that point and a non-nil error.
func (rf *rootFile) ReadDir(n int) ([]fs.DirEntry, error) {
	// list tables
	var entries []fs.DirEntry
	for _, tableInfo := range rf.fsys.tables {
		// entries = append(entries, fs.FileInfoToDirEntry(tableInfo))
		info, _ := tableInfo.Stat()
		entries = append(entries, fs.FileInfoToDirEntry(info))
	}
	return entries, nil
}

type fileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi fileInfo) ModTime() time.Time { return time.Time{} }
func (fi fileInfo) IsDir() bool        { return fi.mode.IsDir() }
func (fi fileInfo) Sys() any           { return nil }

// a table in the DRS archive
type tableFile struct {
	tab *drs.Table
}

var (
	_ fs.ReadDirFile = (*tableFile)(nil)
)

func (tf *tableFile) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name: tf.tab.Extension(),
		mode: fs.ModeDir,
	}, nil
}

func (tf *tableFile) Read([]byte) (int, error) { return 0, io.EOF }
func (tf *tableFile) Close() error             { return nil }

func (tf *tableFile) ReadDir(n int) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	for _, file := range tf.tab.Files {
		info := fileInfo{
			name: file.ID.String(),
			size: int64(file.Size),
		}
		entries = append(entries, fs.FileInfoToDirEntry(info))
	}
	return entries, nil
}
