package drs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type (
	Reader struct {
		Header

		Tables []*Table
		Files  []*File

		rd io.ReaderAt // the underlying reader
	}
)

func Open(filename string) (*Reader, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := fh.Stat()
	if err != nil {
		fh.Close()
		return nil, err
	}

	var r Reader
	if err := r.init(fh, fi.Size()); err != nil {
		fh.Close()
		return nil, err
	}
	r.rd = fh
	return &r, nil
}

func NewReader(rd io.ReaderAt, size int64) (*Reader, error) {
	var reader Reader

	if err := reader.init(rd, size); err != nil {
		return nil, err
	}

	return &reader, nil
}

func (reader *Reader) init(f io.ReaderAt, size int64) error {
	rd := io.NewSectionReader(f, 0, size)
	reader.rd = f

	if err := binary.Read(rd, binary.LittleEndian, &reader.Header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if !bytes.HasPrefix(reader.Header.Copyright[:], copyrightHeader) {
		return ErrInvalidCopyright
	}

	tables := make([]Table, reader.TableCount) // improve memory fragmentation
	reader.Tables = make([]*Table, reader.TableCount)

	// read table
	for i := int32(0); i < reader.TableCount; i++ {
		table := &tables[i]
		if err := binary.Read(rd, binary.LittleEndian, &table.TableInfo); err != nil {
			return fmt.Errorf("failed to read header: %w", err)
		}
		reader.Tables[i] = table
	}

	// get number of files
	nFiles := 0
	for i := range reader.Tables {
		nFiles += int(reader.Tables[i].NumFiles)
	}

	reader.Files = make([]*File, 0, nFiles)

	// parse file info
	for i, table := range reader.Tables {

		if _, err := rd.Seek(int64(table.Offset), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to table %d: %w", i, err)
		}

		files := make([]File, table.NumFiles)
		table.Files = make([]*File, table.NumFiles)

		for j := range files {
			file := &files[j]
			file.rd = rd
			file.drs = reader
			file.table = table
			if err := binary.Read(rd, binary.LittleEndian, &file.FileInfo); err != nil {
				return fmt.Errorf("%w: table %d: %w", ErrInvalidTableEntry, i, err)
			}
			table.Files[j] = file
			reader.Files = append(reader.Files, file)
		}
	}

	return nil
}

// Close closes the underlying reader if it implements io.Closer, otherwise it's a noop.
func (rd *Reader) Close() error {
	if r, ok := rd.rd.(io.Closer); ok {
		return r.Close()
	}
	return nil
}
