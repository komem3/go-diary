package cmd

import (
	"fmt"
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type initCommand struct {
	to      string
	verbose bool
	cmd     *cobra.Command
}

func NewInitCommand() CommandInitializer {
	command := &initCommand{
		cmd: &cobra.Command{
			Use:   "init",
			Short: "Initialize directory",
			Long: `Init make template directory.
You need to run this command before running other command.
`,
		},
	}
	command.cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := command.initialize(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	return command
}

func (c *initCommand) Init() *cobra.Command {
	c.cmd.Flags().StringVarP(
		&c.to,
		"dir",
		"d",
		".",
		"created template directory path",
	)
	c.cmd.Flags().BoolVar(&c.verbose, "v", false, "Output verbose.")
	return c.cmd
}

func (c initCommand) initialize() error {
	logger := diary.NewLogger(c.verbose)
	logger.Debug(
		"msg", "initialize template",
	)
	if err := diary.Initialize(logger, c.to, tmplDir, []string{mdTmp, orgTmp}); err != nil {
		return fmt.Errorf("initialize template: %w", err)
	}
	return nil
}
