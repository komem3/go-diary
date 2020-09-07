package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/komem3/go-diary"
	"github.com/komem3/go-diary/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidSortType = errors.New("sort type is invalid")
)

type sortOption struct {
	year  string
	month string
	day   string
}

type formatter struct {
	from      string
	to        string
	file      string
	templFile string
	sort      sortOption
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

	command.Flags().StringVar(&f.sort.year, "yearSort", f.sort.year,
		"Optional year order. Can specify asc or desc.",
	)
	command.Flags().StringVar(&f.sort.month, "monthSort", f.sort.month,
		"Optional month order. Can specify asc or desc.",
	)
	command.Flags().StringVar(&f.sort.day, "daySort", f.sort.day,
		"Optional day order. Can specify asc or desc.",
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

	options, err := f.createMap2ElemOption()
	if err != nil {
		return fmt.Errorf("create order options: %w", err)
	}

	formatter := diary.NewFormatter(logger)
	logger.Debug(
		"msg", "analys start",
		"dir", f.from,
	)
	fMap := formatter.ParseFileMap(f.from)
	elem := formatter.Map2Elem(fMap, options...)
	formatter.FormatDir(fMap, f.to, f.to == f.from).WriteDirTree(elem, f.file, f.templFile, f.to)

	if formatter.Err != nil {
		return fmt.Errorf("format and write dir tree: %w", formatter.Err)
	}
	return nil
}

func (f formatter) createMap2ElemOption() (options []diary.Map2ElemOptionFunc, err error) {
	if f.sort.year != "" {
		sortType, err := parseSortOption(f.sort.year)
		if err != nil {
			return nil, fmt.Errorf("parse year: %w", err)
		}
		options = append(options, diary.YearSort(sortType))
	}
	if f.sort.month != "" {
		sortType, err := parseSortOption(f.sort.month)
		if err != nil {
			return nil, fmt.Errorf("parse month: %w", err)
		}
		options = append(options, diary.MonthSort(sortType))
	}
	if f.sort.day != "" {
		sortType, err := parseSortOption(f.sort.day)
		if err != nil {
			return nil, fmt.Errorf("parse day: %w", err)
		}
		options = append(options, diary.DaySort(sortType))
	}
	return options, nil
}

func parseSortOption(o string) (diary.SortType, error) {
	if strings.EqualFold(o, "asc") {
		return diary.ASCSort, nil
	}
	if strings.EqualFold(o, "desc") {
		return diary.DESCSort, nil
	}
	return diary.InValid, fmt.Errorf("parse %s: %w", o, ErrInvalidSortType)
}
