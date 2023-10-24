package blendomatic

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// see https://github.com/aap/geniedoc/blob/master/blendomatic.txt
type (
	TileBitmask uint32

	BlendingMode struct {
		TileSize     uint32    // number of pixels, always 2353 since we only have flat tiles
		TileHasAlpha [31]uint8 // 1 if the tile has any alpha pixels

		TileBits  []TileBitmask // bitmask if the pixel has a alpha value; the bit for pixel n of tile m is `tileBits[n] & 1<<m`
		TileAlpha [][]uint8     // the pixels alpha values per tile
	}

	Blendomatic struct {
		Header struct {
			NrBlendingModes uint32
			NrTiles         uint32
		}

		Modes []BlendingMode // nr_blending_modes
	}
)

func (b *BlendingMode) IsAlphaPixel(tile, pixel int) bool {
	// the bit for pixel n of tile m is `tileBits[n] & 1<<m`
	return b.TileBits[pixel]&(1<<tile) == 0
}

func (b *BlendingMode) HasAlpha(tileNr int) bool {
	return b.TileHasAlpha[tileNr] != 0
}

func (b *BlendingMode) GetAlphaValues(tileNr int) []uint8 {
	return b.TileAlpha[tileNr]
}

func Open(filename string) (*Blendomatic, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	rd := bufio.NewReader(fh)
	return New(rd)
}

func (mode *BlendingMode) parse(rd io.Reader, nr_tiles int) error {
	if err := binary.Read(rd, binary.LittleEndian, &mode.TileSize); err != nil {
		return fmt.Errorf("failed to decode blendingmode (TileSize): %w", err)
	}
	if _, err := io.ReadAtLeast(rd, mode.TileHasAlpha[:], len(mode.TileHasAlpha)); err != nil {
		return fmt.Errorf("failed to decode blendingmode (TileFlags): %w", err)
	}

	mode.TileBits = make([]TileBitmask, mode.TileSize)
	if err := binary.Read(rd, binary.LittleEndian, &mode.TileBits); err != nil {
		return fmt.Errorf("failed to decode blendingmode (TileBits): %w", err)
	}

	mode.TileAlpha = make([][]byte, len(mode.TileHasAlpha))
	for i := 0; i < len(mode.TileHasAlpha); i++ {
		if mode.TileHasAlpha[i] == 0 {
			panic("implement me!")
		}
		alpha := make([]byte, mode.TileSize)

		if _, err := io.ReadAtLeast(rd, alpha, len(alpha)); err != nil {
			return fmt.Errorf("failed to decode blendingmode (TileFlags): %w", err)
		}
		mode.TileAlpha[i] = alpha
	}

	return nil
}

func New(rd io.Reader) (*Blendomatic, error) {
	var bm Blendomatic

	if err := binary.Read(rd, binary.LittleEndian, &bm.Header); err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}
	for i := 0; i < int(bm.Header.NrBlendingModes); i++ {
		var mode BlendingMode
		if err := mode.parse(rd, int(bm.Header.NrTiles)); err != nil {
			return nil, fmt.Errorf("error decoding blendingmode %d: %w", i, err)
		}
		bm.Modes = append(bm.Modes, mode)
	}

	return &bm, nil
}
