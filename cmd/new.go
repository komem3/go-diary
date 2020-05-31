package cmd

import (
	"fmt"
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type newer struct {
	tmplFile   string
	dir        string
	nameFormat string
	verbose    bool
}

func newNewer() *newer {
	return &newer{
		verbose:    false,
		tmplFile:   "template/diary.template.md",
		dir:        ".",
		nameFormat: "20060102.md",
	}
}

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "new",
		Short: "Generate new diary",
		Long:  "Generate new today diary from template file.",
	}

	n := newNewer()
	command.Flags().StringVar(&n.tmplFile, "tmpl", n.tmplFile,
		"Parse template file.",
	)
	command.Flags().StringVarP(&n.dir, "dir", "d", n.dir,
		"Destination directory.",
	)
	command.Flags().StringVarP(&n.nameFormat, "format", "f", n.nameFormat,
		"File name format.\nRefer to https://golang.org/src/time/format.go",
	)

	command.Run = func(cmd *cobra.Command, args []string) {
		if err := n.New(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	return command
}

func (n newer) New() error {
	logger := diary.NewLogger(n.verbose)

	generator := diary.NewGenerator(logger)
	generator.NewDiary(n.tmplFile, n.dir, n.nameFormat)
	return generator.Err
}
