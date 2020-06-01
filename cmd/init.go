package cmd

import (
	"fmt"
	"os"

	"github.com/komem3/go-diary"
	"github.com/spf13/cobra"
)

type initializer struct {
	to      string
	verbose bool
}

func newInitializer() *initializer {
	return &initializer{
		to:      ".",
		verbose: false,
	}
}

func NewInitCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize directory",
		Long: `Init command make template directory.
You need to run this command before running other command.
`,
	}

	i := newInitializer()
	command.Flags().StringVarP(
		&i.to,
		"dir",
		"d",
		i.to,
		"Created template directory path",
	)
	command.Flags().BoolVar(&i.verbose, "v", i.verbose, "Output verbose.")

	command.Run = func(cmd *cobra.Command, args []string) {
		if err := i.init(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	return command
}

func (i initializer) init() error {
	logger := diary.NewLogger(i.verbose)
	logger.Debug(
		"msg", "initialize template",
	)
	if err := diary.Initialize(logger, i.to, tmplDir, []string{mdTmp(), diaryTmp()}); err != nil {
		return fmt.Errorf("initialize template: %w", err)
	}
	return nil
}
