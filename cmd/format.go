package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/komem3/go-diary"
	"github.com/komem3/go-diary/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type formatter struct {
	from      string
	to        string
	file      string
	templFile string
	verbose   bool
}

func newFormatter() *formatter {
	return &formatter{
		from:      ".",
		to:        "",
		file:      "./README.md",
		templFile: mdTmp(),
		verbose:   false,
	}
}

func NewFormatCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "format",
		Short: "Format directory",
		Long: `Format command format directory.
After format directory, it write directory structure to target file.
`,
	}

	f := newFormatter()
	command.Flags().StringVarP(&f.from, "dir", "d", f.from,
		"Analysis directory.",
	)
	command.Flags().StringVar(&f.to, "copyDir", f.to,
		"Format directory. \nWhen this option is difference from 'dir', all file will copy to 'copyDir'.",
	)
	command.Flags().BoolVar(&f.verbose, "v", f.verbose, "Output verbose.")
	command.Flags().StringVarP(&f.file, "file", "f", f.file,
		"Write file.\nnThe environment variable DIARY_INDEX_FILE is set.",
	)
	utils.ErrorPanic(viper.BindPFlag("format_file", command.Flags().Lookup("file")))
	utils.ErrorPanic(viper.BindEnv("format_file", "DIARY_INDEX_FILE"))
	f.file = viper.GetString("format_file")
	command.Flags().StringVar(&f.templFile, "tmpl", f.templFile,
		"Parse template file.\nThe environment variable DIARY_INDEX_TEMPLATE is set.",
	)
	utils.ErrorPanic(viper.BindPFlag("format_templFile", command.Flags().Lookup("tmpl")))
	utils.ErrorPanic(viper.BindEnv("format_templFile", "DIARY_INDEX_TEMPLATE"))
	f.templFile = viper.GetString("format_templFile")

	command.Run = func(cmd *cobra.Command, args []string) {
		if err := f.Format(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "write index to %s\n", filepath.Join(f.to, f.file))
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
	formatter.FormatDir(fMap, f.to, f.to == f.from).WriteDirTree(elem, f.file, f.templFile, f.to)

	if formatter.Err != nil {
		return fmt.Errorf("format and write dir tree: %w", formatter.Err)
	}
	return nil
}
