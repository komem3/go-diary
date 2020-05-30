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
	orgMode  bool
	verbose  bool
}

func newFormatter() *formatter {
	return &formatter{
		from:     ".",
		to:       "",
		file:     "",
		tempFile: "",
		orgMode:  false,
		verbose:  false,
	}
}

func NewFormatCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "format",
		Short: "Format directory",
		Long: `Format command analys and format target directory.
After format directory, this rewrite directory structure to target file.
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
		"Rewrite file.\nDefault value is './README.md.\nIn org mode value is ./README.org",
	)
	command.PersistentFlags().StringVar(&f.tempFile, "tmpl", f.tempFile,
		"Parse template file.\nDefault is "+mdTmp()+".\nIn org mode value is "+orgTmp()+".",
	)
	command.PersistentFlags().BoolVar(&f.orgMode, "org", f.orgMode,
		"Use org template.",
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
	if f.orgMode {
		if f.file == "" {
			f.file = "./README.org"
		}
		if f.tempFile == "" {
			f.tempFile = orgTmp()
		}
	} else {
		if f.file == "" {
			f.file = "./README.md"
		}
		if f.tempFile == "" {
			f.tempFile = mdTmp()
		}
	}
	dGen, err := diary.NewDiaryGenerator(f.file, f.tempFile, f.from, f.to, f.orgMode, logger)
	if err != nil {
		return fmt.Errorf("generate generator: %w", err)
	}
	logger.Debug(
		"msg", "analys start",
		"dir", f.from,
	)
	fMap := diary.ParseFileMap(f.from)
	elem := diary.Map2Elem(fMap)
	dGen.FormatDir(fMap).WriteDirTree(elem)

	if dGen.Err != nil {
		return fmt.Errorf("format and write dir tree: %w", dGen.Err)
	}
	return nil
}
