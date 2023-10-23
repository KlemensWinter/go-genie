package icm

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"os"
)

const (
	numRows = 10 // number of brightness levels
)

type (
	Map [numRows][32][32][32]byte
)

func Open(filename string) (f Map, err error) {
	err = f.Open(filename)
	return
}

func New(rd io.Reader) (f Map, err error) {
	err = f.Load(rd)
	return
}

func (f *Map) Open(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return f.Load(fh)
}

func (f *Map) Load(rd io.Reader) error {
	if err := binary.Read(rd, binary.LittleEndian, f); err != nil {
		return fmt.Errorf("icm: %w", err)
	}
	return nil
}

func (f Map) Index(brightness uint8, c color.Color) int {
	r, g, b, _ := c.RGBA()
	r = (r >> 10) & 0xff
	g = (g >> 10) & 0xff
	b = (b >> 10) & 0xff
	return int(f[brightness][r][g][b])
}
