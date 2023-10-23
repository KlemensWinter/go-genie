package slp

import (
	"errors"
	"image"
	"io"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

// see also https://github.com/SFTtech/openage/blob/master/doc/media/slp-files.md
type (
	Header struct {
		Version   [4]byte
		NumFrames int32
		Comment   [24]byte
	}

	FrameInfo struct {
		CmdTableOffset     uint32
		OutlineTableOffset uint32
		PaletteOffset      uint32
		Properties         uint32
		Width              int32
		Height             int32
		HotspotX           int32
		HotspotY           int32
	}

	Outline struct {
		LeftSpace  uint16
		RightSpace uint16
	}

	Frame struct {
		FrameInfo

		Outline    []Outline
		CmdOffsets []uint32

		dataSize int64
		slpr     *Reader
	}
)

func (fi FrameInfo) Size() image.Point {
	return image.Point{X: int(fi.Width), Y: int(fi.Height)}
}

func (fi FrameInfo) Hotspot() image.Point {
	return image.Point{X: int(fi.HotspotX), Y: int(fi.HotspotY)}
}

// Width returns the width of this outline (left + right space)
func (ol Outline) Width() int {
	return int(ol.LeftSpace) + int(ol.RightSpace)
}

func (frame *Frame) Open() *io.SectionReader {
	startCommandData := int64(frame.CmdTableOffset) + 4*int64(frame.Height)
	return io.NewSectionReader(frame.slpr, startCommandData, frame.dataSize)
}

/*
func (frame *Frame) DrawTo(img *image.RGBA, pal color.Palette, playerID int) error {
	return DrawTo(img, pal, frame, playerID, 0)
}

func (frame *Frame) GetImage(pal color.Palette, playerID int) (img *image.RGBA, err error) {
	img = image.NewRGBA(image.Rect(0, 0, int(frame.Width), int(frame.Height)))
	err = DrawTo(img, pal, frame, playerID, 0)
	return
}
*/

func (r *Reader) ReadAt(p []byte, off int64) (n int, err error) {
	return r.rd.ReadAt(p, off)
}

func (rc *ReadCloser) Close() error {
	return rc.fh.Close()
}

func (f *Frame) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(f.Width), int(f.Height))
}
