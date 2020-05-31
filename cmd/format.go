package cmd

import (
	"fmt"
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type formatter struct {
	from     string
	to       string
	file     string
	tempFile string
	verbose  bool
}

func newFormatter() *formatter {
	return &formatter{
		from:     ".",
		to:       "",
		file:     "./README.md",
		tempFile: mdTmp(),
		verbose:  false,
	}
}

func NewFormatCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "format",
		Short: "Format directory",
		Long: `Format command analys and format directory.
After format directory, it write directory structure to target file.
`,
	}

	f := newFormatter()
	command.PersistentFlags().StringVarP(&f.from, "dir", "d", f.from,
		"Analysis directory.",
	)
	command.PersistentFlags().StringVar(&f.to, "copyDir", f.to,
		"Format directory. \nWhen this option is difference from 'dir', all file will copy to 'copyDir'.",
	)
	command.PersistentFlags().StringVarP(&f.file, "file", "f", f.file,
		"Write file.",
	)
	command.PersistentFlags().StringVar(&f.tempFile, "tmpl", f.tempFile,
		"Parse template file.",
	)
	command.Flags().BoolVar(&f.verbose, "v", f.verbose, "Output verbose.")

	command.Run = func(cmd *cobra.Command, args []string) {
		if err := f.Format(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	return command
}

func (f formatter) Format() error {
	logger := diary.NewLogger(f.verbose)
	if f.to == "" {
		f.to = f.from
	}

	formatter := diary.NewFormatter(logger)
	logger.Debug(
		"msg", "analys start",
		"dir", f.from,
	)
	fMap := formatter.ParseFileMap(f.from)
	elem := formatter.Map2Elem(fMap)
	formatter.FormatDir(fMap, f.to, f.to == f.from).WriteDirTree(elem, f.file, f.tempFile, f.to)

	if formatter.Err != nil {
		return fmt.Errorf("format and write dir tree: %w", formatter.Err)
	}
	return nil
}
