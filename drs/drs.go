package drs

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"strings"

	"slices"
)

const (
	InvalidFileID = FileID(math.MaxUint32)
)

var (
	ErrInvalidCopyright  = errors.New("drs: invalid copyright string")
	ErrInvalidTableEntry = errors.New("drs: invalid table entry")
)

var (
	extBin          = []byte("anib")
	copyrightHeader = []byte("Copyright (c) 1997 Ensemble Studios.")
)

type (
	FileID uint32

	Header struct {
		Copyright [40]byte
		Version   [4]byte
		Ftype     [12]byte

		TableCount int32
		FileOffset int32
	}

	TableInfo struct {
		FileExtension [4]byte
		Offset        int32
		NumFiles      int32
	}

	Table struct {
		TableInfo

		Files []*File
	}
)

func (id FileID) String() string {
	return strconv.Itoa(int(id))
}

func (id FileID) IsValid() bool {
	return id != InvalidFileID && id != 0
}

func FormatExtension(ext []byte) string {
	if len(ext)%2 != 0 {
		panic("unreachable")
	}
	if bytes.Equal(ext, extBin) {
		return "bin"
	}
	tmp := slices.Clone(ext)
	n := len(ext)
	for i := 0; i < n/2; i++ {
		tmp[i], tmp[n-i-1] = tmp[n-i-1], tmp[i]
	}
	return strings.TrimSpace(string(tmp))
}

func (ti *TableInfo) Extension() string {
	return FormatExtension(ti.FileExtension[:])
}
