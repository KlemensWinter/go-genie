package cmd

import (
	"fmt"
	"image/color"

	"github.com/KlemensWinter/go-genie/palette"
	"github.com/spf13/cobra"
)

var paletteCmd = &cobra.Command{
	Use:  "palette <FILE>",
	RunE: runPaletteShow,
	Args: cobra.ExactArgs(1),
}

func init() {
	fl := paletteCmd.Flags()

	fl.Bool("hex", false, "print as hex values")
	fl.Bool("go", false, "print as go struct")
}

func runPaletteShow(cmd *cobra.Command, args []string) error {
	filename := args[0]
	pal, err := palette.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open palette: %w", err)
	}

	out := cmd.OutOrStdout()

	asHex, _ := cmd.Flags().GetBool("hex")
	asGo, _ := cmd.Flags().GetBool("go")

	for _, col := range pal {
		c := col.(color.RGBA)
		switch {
		case asHex:
			fmt.Fprintf(out, "#%02x%02x%02x\n", c.R, c.G, c.B)
		case asGo:
			fmt.Fprintf(out, "%#v\n", col)
		default:
			fmt.Fprintf(out, "%d %d %d\n", c.R, c.G, c.B)
		}
	}

	return nil
}
