package palette

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
)

const (
	Header     = "JASC-PAL"
	palVersion = "0100"
)

var (
	ErrInvalidVersion = errors.New("invalid version")
	ErrInvalidHeader  = errors.New("invalid header")
	ErrInvalidLine    = errors.New("invalid line")
)

// Parse parses a JASC palette.
func Parse(rd io.Reader) (color.Palette, error) {
	brd, ok := rd.(*bufio.Reader)
	if !ok {
		brd = bufio.NewReader(rd)
	}

	// JASC-PAL
	ln, err := brd.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	if strings.TrimSpace(ln) != Header {
		return nil, fmt.Errorf("%w: %q", ErrInvalidHeader, strings.TrimSpace(ln))
	}
	// 0100
	ln, err = brd.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}
	ver := strings.TrimSpace(ln)
	if ver != palVersion {
		return nil, fmt.Errorf("%w: %q", ErrInvalidVersion, ver)
	}

	var numColors int
	if _, err := fmt.Fscanf(brd, "%d\n", &numColors); err != nil {
		return nil, fmt.Errorf("failed to parse color count: %w", err)
	}

	lineNr := 3
	pal := make(color.Palette, numColors)
	for i := 0; i < int(numColors); i++ {
		var col color.RGBA
		col.A = 255

		_, err := fmt.Fscanf(brd, "%d %d %d\n", &col.R, &col.G, &col.B)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse line %d: %w", ErrInvalidLine, i, err)
		}

		lineNr++
		pal[i] = col
	}
	return pal, nil
}

// Open a JASC palette from the local filesystem.
func Open(filename string) (color.Palette, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return Parse(fh)
}

// Marshal returns the given palette JASC file.
func Marshal(pal color.Palette) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(Header)
	buf.WriteString("\r\n0100\r\n")
	fmt.Fprintf(&buf, "%d\r\n", len(pal))
	for i := range pal {
		r, g, b, _ := pal[i].RGBA()
		fmt.Fprintf(&buf, "%d %d %d\r\n", r>>8, g>>8, b>>8)
	}
	return buf.Bytes(), nil
}
