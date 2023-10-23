package slp

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

type DrawFlags int

const (
	CleanOutline DrawFlags = 1 << iota
)

type bufrd []byte

func (rd *bufrd) getc() (c byte) {
	c, *rd = (*rd)[0], (*rd)[1:]
	return
}

func (rd *bufrd) rshiftOrNext(cmd_byte byte, shift int) int {
	count := int(cmd_byte) >> shift
	if count == 0 {
		count = int(rd.getc())
	}
	return count
}

func (rd *bufrd) lshiftAndNext(cmd_byte byte) int {
	return (int(cmd_byte&0xf0) << 4) + int(rd.getc())
}

func (rd *bufrd) skip(n int) {
	n = min(len(*rd), n)
	*rd = (*rd)[n:]
}

func (r *bufrd) readBytes(dst []byte, n int) {
	copy(dst, (*r)[:n])
	*r = (*r)[n:]
}

func drawLine(r *bufrd, img *image.RGBA, pal color.Palette, x, y int, playerID int) error {
	xOffset, yOffset := img.Rect.Min.X, img.Rect.Min.Y

	setPix := func(x, y int, colorIndex byte) {
		img.Set(x+xOffset, y+yOffset, pal[colorIndex])
	}

mainloop:
	for {
		cmd_byte := r.getc()
		nib := cmd_byte & 0x0f

		if nib&0b11 == 0 {
			n := int(cmd_byte >> 2)
			for i := 0; i < n; i++ {
				setPix(x, y, r.getc())
				x++
			}
		} else if nib&0b11 == 1 { // lesser skip
			n := r.rshiftOrNext(cmd_byte, 2)
			x += n
		} else {
			switch Cmd(nib) {
			case CMD_GREATER_DRAW: // greater draw
				n := r.lshiftAndNext(cmd_byte)
				for i := 0; i < n; i++ {
					setPix(x, y, r.getc())
					x++
				}
			case CMD_GREATER_SKIP: // greater skip
				n := r.lshiftAndNext(cmd_byte)
				x += n
			case CMD_PLAYER_COLOR_DRAW:
				count := r.rshiftOrNext(cmd_byte, 4)
				for i := 0; i < count; i++ {
					setPix(x, y, r.getc()+byte(playerID*16))
					x++
				}
			case CMD_FILL:
				n := r.rshiftOrNext(cmd_byte, 4)
				col := r.getc()
				for i := 0; i < n; i++ {
					setPix(x, y, col)
					x++
				}
			case CMD_FILL_PLAYER_COLOR:
				n := r.rshiftOrNext(cmd_byte, 4)
				col := r.getc() + byte(playerID*16)
				for i := 0; i < n; i++ {
					setPix(x, y, col)
					x++
				}
			case CMD_SHADOW_DRAW:
				n := r.rshiftOrNext(cmd_byte, 4)

				if false {
					for i := 0; i < n; i++ {
						// TODO: pix[x] = 0 // shadow color
						x++
					}
				}

			case CMD_EXTENDED:
				switch Cmd(cmd_byte) {
				case CMD_EXT_FORWARD_DRAW:
					return fmt.Errorf("%w: CMD_EXT_FORWARD_DRAW", ErrNotImplemented)
				case CMD_EXT_REVERSE_DRAW:
					return fmt.Errorf("%w: CMD_EXT_REVERSE_DRAW", ErrNotImplemented)

				case CMD_EXT_NORMAL_TRANSFORM:
				case CMD_EXT_ALTERNATE_TRANSFORM:
					return fmt.Errorf("%w: extended command %+x", ErrNotImplemented, cmd_byte)

				case CMD_EXT_OUTLINE1:
					x++

				case CMD_EXT_OUTLINE1_FILL:
					n := int(r.getc())
					for i := 0; i < n; i++ {
						x++
					}
				case CMD_EXT_OUTLINE2:
					return fmt.Errorf("%w: CMD_EXT_OUTLINE2", ErrNotImplemented)
				case CMD_EXT_OUTLINE2_FILL:
					return fmt.Errorf("%w: CMD_EXT_OUTLINE2_FILL", ErrNotImplemented)
				case Cmd(0x8E):
					return fmt.Errorf("%w: dither", ErrNotImplemented)
				case CMD_EXT_PREMULTIPLIED_ALPHA:
					n := r.getc()
					_ = n
					return fmt.Errorf("%w: CMD_EXT_PREMULTIPLIED_ALPHA", ErrNotImplemented)
				default:
					return fmt.Errorf("%w: Decode(): extended command %#x", ErrNotImplemented, cmd_byte)
				}
			case CMD_END_OF_ROW:
				break mainloop
			default:
				panic(fmt.Errorf("invalid command %#b", nib))
			}
		}
	}
	return nil
}

func drawTo(img *image.RGBA, pal color.Palette, data []byte, f *Frame, playerID int, flags DrawFlags) error {

	maxY := int(f.Height)
	r := bufrd(data)

	for y := 0; y < maxY; y++ {
		x := int(f.Outline[y].LeftSpace)
		drawLine(&r, img, pal, x, y, playerID)
	}
	/*
		if flags&CleanOutline != 0 {
			// xOffset, yOffset := img.Rect.Min.X, img.Rect.Min.Y
			// erase the outline
			//
			// 	for y := 0; y < maxY; y++ {
			// 		for x := 0; x < int(f.Outline[y].LeftSpace); x++ {
			// 			img.SetColorIndex(x+xOffset, y+yOffset, color.RGBA{})
			// 		}
			// 		for x := f.Width - int32(f.Outline[y].RightSpace); x < f.Width; x++ {
			// 			img.SetRGBA(int(x)+xOffset, y+yOffset, color.RGBA{})
			// 		}
			// 	}
			//
		}
	*/
	return nil
}

func DrawTo(img *image.RGBA, pal color.Palette, f *Frame, playerID int, flags DrawFlags) error {
	data, err := io.ReadAll(f.Open())
	if err != nil {
		return err
	}
	return drawTo(img, pal, data, f, playerID, flags)
}
