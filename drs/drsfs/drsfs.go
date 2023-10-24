package drsfs

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"

	"gopkg.in/KlemensWinter/go-genie.v1/drs"
)

type FS struct {
	rd *drs.Reader

	tables map[string]*tableFile
}

var (
	_ fs.FS = (*FS)(nil)
	// _ fs.ReadDirFS  = (*FS)(nil)
	//_ fs.ReadFileFS = (*FS)(nil)

	errUnreadable = errors.New("can not read a directory")
)

func OpenFS(filename string) (*FS, error) {
	rd, err := drs.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewFS(rd)
}

func NewFS(f *drs.Reader) (*FS, error) {
	fsys := &FS{
		rd:     f,
		tables: make(map[string]*tableFile),
	}

	for _, table := range f.Tables {
		tab := &tableFile{tab: table}
		fsys.tables[table.Extension()] = tab
	}

	return fsys, nil
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fmt.Errorf("invalid path %q", name),
		}
	}
	if name == "." {
		// open root file
		return &rootFile{fsys}, nil
	}
	dir, fname := path.Split(name)

	if dir == "" {
		// we have a table!
		tab, found := fsys.tables[fname]
		if !found {
			return nil, &fs.PathError{
				Op:   "Open",
				Path: name,
				Err:  fs.ErrNotExist,
			}
		}
		return tab, nil
	}

	log.Printf("OPEN: %q - dir=%q fname=%q", name, dir, fname)
	panic("Open: implement me!")
}

func (f *FS) OpenID(id drs.FileID) (fs.File, error) {
	for _, fh := range f.rd.Files {
		if fh.ID == id {
			return NewFile(fh)
		}
	}
	return nil, fs.ErrNotExist
}
