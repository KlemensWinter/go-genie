package lang

import (
	"encoding/binary"
	"fmt"

	"github.com/saferwall/pe"
	"golang.org/x/text/encoding/unicode"
)

const (
	MaxStringsPerLeaf = 16
)

type (
	Entry struct {
		ID      uint32
		Strings []string
	}
)

// https://github.com/sandsmark/pcrio/blob/master/pcrio.c

func ParseID(id uint32) (dir, offset int) {
	dir = int(id/MaxStringsPerLeaf) + 1
	offset = int(id % MaxStringsPerLeaf)
	return
}

func DecodeUTF16(data []byte) (string, error) {
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	buf, err := dec.Bytes(data)
	return string(buf), err
}

func readStrings(buf []byte) ([]string, error) {
	var res []string
	// read strings
	for len(buf) > 2 {
		// get length
		l := binary.LittleEndian.Uint16(buf[:2]) // 2 byte length
		buf = buf[2:]
		// data
		n := l * 2
		if n == 0 {
			// empty
			res = append(res, "")
		} else {
			data := buf[:n]
			if int(n) >= len(buf) { // TODO: ok?
				buf = nil
			} else {
				buf = buf[n:]
			}
			str, err := DecodeUTF16(data)
			if err != nil {
				return nil, fmt.Errorf("failed to read string: %w", err)
			}
			res = append(res, str)
		}
	}
	return res, nil
}

// isLeaf returns true if the hight bit of address is 0 (a leaf)
func isLeaf(entry pe.ResourceDirectoryEntry) bool {
	return entry.Struct.OffsetToData&(uint32(1)<<31) == 0
}

func readEntry(file *pe.File, entry *pe.ResourceDataEntry) ([]string, error) {
	sec := file.Sections[pe.ImageDirectoryEntryResource]

	pos := entry.Struct.OffsetToData
	size := entry.Struct.Size
	codePage := entry.Struct.CodePage

	if codePage != 0x4e4 {
		panic(fmt.Errorf("invalid codepage %#x", codePage))
	}

	data := sec.Data(pos, size, file)
	return readStrings(data)
}

func Open(filename string) ([]*Entry, error) {
	fh, err := pe.New(filename, &pe.Options{})
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	if err = fh.Parse(); err != nil {
		return nil, err
	}
	var entries []*Entry
	for _, entry := range fh.Resources.Entries {
		for _, child := range entry.Directory.Entries {
			dir := child.Directory
			if len(dir.Entries) != 1 {
				panic("invalid entry")
			}
			st, err := readEntry(fh, &dir.Entries[0].Data)
			if err != nil {
				return nil, err
			}
			e := &Entry{ID: child.ID,
				Strings: st,
			}
			entries = append(entries, e)
		}
	}
	return entries, nil
}
