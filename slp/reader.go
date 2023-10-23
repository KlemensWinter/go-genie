package slp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type (
	Reader struct {
		Header Header
		Frames []*Frame

		rd   io.ReaderAt
		size int64
	}

	ReadCloser struct {
		Reader

		fh *os.File
	}
)

func New(rd io.ReaderAt, size int64) (*Reader, error) {
	r := &Reader{}
	if err := r.init(rd, size); err != nil {
		return nil, err
	}
	return r, nil
}

func Open(filename string) (*ReadCloser, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, _ := fh.Stat()

	reader := &ReadCloser{
		fh: fh,
	}
	if err := reader.init(fh, fi.Size()); err != nil {
		fh.Close()
		return nil, err
	}
	return reader, nil
}

// GetMaxSize returns the max size a frame will have.
func (reader *Reader) GetMaxSize() (w, h int) {
	for _, f := range reader.Frames {
		w = max(w, int(f.Width))
		h = max(h, int(f.Height))
	}
	return
}

func (rd *Reader) NumFrames() int {
	return len(rd.Frames)
}

func (reader *Reader) init(rd io.ReaderAt, size int64) error {
	r := io.NewSectionReader(rd, 0, size)
	if err := binary.Read(r, binary.LittleEndian, &reader.Header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if !bytes.Equal(reader.Header.Version[:], []byte{'2', '.', '0', 'N'}) {
		return fmt.Errorf("invalid version: %q", reader.Header.Version)
	}

	reader.size = size

	for i := 0; i < int(reader.Header.NumFrames); i++ {
		frame := &Frame{
			slpr: reader,
		}
		if err := binary.Read(r, binary.LittleEndian, &frame.FrameInfo); err != nil {
			return fmt.Errorf("failed to read frame info: %w", err)
		}
		if frame.OutlineTableOffset >= frame.CmdTableOffset {
			panic("invalid offset!")
		}
		reader.Frames = append(reader.Frames, frame)
	}
	for _, frame := range reader.Frames {
		if _, err := r.Seek(int64(frame.OutlineTableOffset), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to outline table: %w", err)
		}
		frame.Outline = make([]Outline, frame.Height)
		if err := binary.Read(r, binary.LittleEndian, &frame.Outline); err != nil {
			return fmt.Errorf("failed to read outline table: %w", err)
		}
	}

	// cmd offsets
	for i, frame := range reader.Frames {
		if _, err := r.Seek(int64(frame.CmdTableOffset), io.SeekStart); err != nil {
			return fmt.Errorf("frame %d: failed to seek to cmd table: %w", i, err)
		}
		frame.CmdOffsets = make([]uint32, frame.Height)
		if err := binary.Read(r, binary.LittleEndian, &frame.CmdOffsets); err != nil {
			return fmt.Errorf("frame %d: failed to read cmd table: %w", i, err)
		}
	}

	// calculate size
	for i, frame := range reader.Frames {
		endPos := size
		if i+1 < len(reader.Frames) {
			endPos = int64(reader.Frames[i+1].OutlineTableOffset)
		}
		// log.Printf("Next frame: %s", reader.Frames[i+1])
		if endPos > size {
			panic("invalid end pos")
		}
		frame.dataSize = int64(endPos) - int64(frame.OutlineTableOffset)
		if frame.dataSize <= 0 {
			panic("invalid data size")
		}
	}

	reader.rd = r

	return nil
}
