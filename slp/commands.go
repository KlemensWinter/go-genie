//go:generate stringer -type=Cmd
package slp

type Cmd byte

const (
	CMD_INVALID           = Cmd(^byte(0))
	CMD_LESSER_DRAW       = Cmd(0b00)
	CMD_LESSER_SKIP       = Cmd(0b01)
	CMD_GREATER_DRAW      = Cmd(0b0010)
	CMD_GREATER_SKIP      = Cmd(0b0011)
	CMD_PLAYER_COLOR_DRAW = Cmd(0b0110)
	CMD_FILL              = Cmd(0b0111)
	CMD_FILL_PLAYER_COLOR = Cmd(0b1010)
	CMD_SHADOW_DRAW       = Cmd(0b1011)
	CMD_EXTENDED          = Cmd(0x0e)
	CMD_END_OF_ROW        = Cmd(0x0f)

	CMD_EXT_FORWARD_DRAW        = Cmd(0x0e)
	CMD_EXT_REVERSE_DRAW        = Cmd(0x1e)
	CMD_EXT_NORMAL_TRANSFORM    = Cmd(0x2e)
	CMD_EXT_ALTERNATE_TRANSFORM = Cmd(0x3e)
	CMD_EXT_OUTLINE1            = Cmd(0x4e)
	CMD_EXT_OUTLINE1_FILL       = Cmd(0x5e)
	CMD_EXT_OUTLINE2            = Cmd(0x6e)
	CMD_EXT_OUTLINE2_FILL       = Cmd(0x7e)

	CMD_EXT_PREMULTIPLIED_ALPHA = Cmd(0x9e)
)
