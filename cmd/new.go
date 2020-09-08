package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/komem3/go-diary"
	"github.com/komem3/go-diary/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type newer struct {
	tmplFile   string
	dir        string
	date       string
	nameFormat string
	verbose    bool
}

func newNewer() *newer {
	return &newer{
		verbose:    false,
		tmplFile:   "template/diary.template.md",
		dir:        ".",
		date:       "today",
		nameFormat: "20060102.md",
	}
}

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "new",
		Short: "Create new diary",
		Long:  "New command create new today diary from template file.",
	}

	n := newNewer()
	command.Flags().StringVar(&n.dir, "dir", n.dir,
		"Destination directory.",
	)
	command.Flags().StringVarP(&n.date, "date", "d", n.date,
		`Date of making diary.
Format: YYYY/MM/dd(2010/01/31) or today(t) or yesterday(y) or tomorrow(tm).`,
	)
	command.Flags().BoolVar(&n.verbose, "v", n.verbose, "Output verbose.")
	command.Flags().StringVarP(&n.nameFormat, "format", "f", n.nameFormat,
		"File name format.\nRefer to https://golang.org/src/time/format.go",
	)
	command.Flags().String("tmpl", n.tmplFile,
		"Parse template file.\nThe environment variable DIARY_NEW_TEMPLATE is set.",
	)
	utils.ErrorPanic(viper.BindPFlag("tmpl", command.Flags().Lookup("tmpl")))
	utils.ErrorPanic(viper.BindEnv("tmple", "DIARY_NEW_TEMPLATE"))

	command.Run = func(cmd *cobra.Command, args []string) {
		diaryPath, err := n.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "generated %s\n", diaryPath)
	}
	return command
}

func (n newer) New() (string, error) {
	n.tmplFile = viper.GetString("tmpl")

	logger := diary.NewLogger(n.verbose)

	generator := diary.NewCreator(logger)
	if err := n.formatAndSetNow(generator); err != nil {
		return "", err
	}

	diaryPath := generator.NewDiary(n.tmplFile, n.dir, n.nameFormat)
	return diaryPath, generator.Err
}

func (n newer) formatAndSetNow(creator *diary.Creator) error {
	switch n.date {
	case "today", "t":
		creator.SetNowFunc(time.Now)
		return nil
	case "yesterday", "y":
		creator.SetNowFunc(func() time.Time { return time.Now().AddDate(0, 0, -1) })
		return nil
	case "tomorrow", "tm":
		creator.SetNowFunc(func() time.Time { return time.Now().AddDate(0, 0, 1) })
		return nil
	}

	re := regexp.MustCompile(`^((1|2)[0-9]{3})/([0-9]{1,2})/([0-9]{1,2})$`)
	submathes := re.FindStringSubmatch(n.date)
	if len(submathes) == 0 {
		return nil
	}

	year, err := strconv.Atoi(submathes[1])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(submathes[3])
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(submathes[4])
	if err != nil {
		return err
	}
	creator.SetNowFunc(func() time.Time { return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local) })
	return nil
}
