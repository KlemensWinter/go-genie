package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/KlemensWinter/go-genie/drs"
	"github.com/spf13/cobra"
)

var drsCmd = &cobra.Command{
	Use: "drs <command>",
}

func listDRS(cmd *cobra.Command, args []string) error {
	filename := args[0]
	rd, err := drs.Open(filename)
	if err != nil {
		return err
	}
	defer rd.Close()

	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "copyright: %q\n", rd.Header.Copyright)
	fmt.Fprintf(out, "version: %q\n", rd.Header.Version)
	fmt.Fprintf(out, "ftype: %q\n", rd.Header.Ftype)
	fmt.Fprintf(out, "tableCount: %d\n", rd.Header.TableCount)
	fmt.Fprintf(out, "fileOffset: %#x\n", rd.Header.FileOffset)
	fmt.Fprintln(out)

	for _, table := range rd.Tables {
		fmt.Fprintf(out, "table: %q:\n", table.FileExtension)
		fmt.Fprintf(out, "numFiles: %d\n", table.TableInfo.NumFiles)
		fmt.Fprintf(out, "offset: %#x\n", table.TableInfo.Offset)
		fmt.Fprintln(out)
		w := tabwriter.NewWriter(out, 10, 4, 1, ' ', 0)
		fmt.Fprintf(w, "ID\tOffset\tSize\n")
		fmt.Fprintf(w, "--\t------\t----\n")
		for _, file := range table.Files {
			if file == nil {
				fmt.Fprintf(w, "-\n")
			} else {
				fmt.Fprintf(w, "%d\t%#x\t%d\n", file.ID, file.Offset, file.Size)
			}
		}
		if err := w.Flush(); err != nil {
			return err
		}
		fmt.Fprintln(out)
	}

	return nil
}

func init() {
	drsCmd.AddCommand(
		&cobra.Command{
			Use:   "list <FILE>",
			Short: "list the content of a DRS file",
			Args:  cobra.ExactArgs(1),
			RunE:  listDRS,
		},
	)
}
