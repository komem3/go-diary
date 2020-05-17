package cmd

import (
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type initCommand struct {
	cmd *cobra.Command
}

var (
	initCmd = &initCommand{
		&cobra.Command{
			Use:   "init",
			Short: "Initialize directory",
			Long: `Init make template directory.
You need to run this command before running other command.
`,
			Run: func(cmd *cobra.Command, args []string) {
				logger := diary.NewLogger(verbose)
				logger.Debug(
					"msg", "initialize template",
				)
				if err := diary.Initialize(logger, to, tmplDir, []string{mdTmp, orgTmp}); err != nil {
					logger.Error(
						"when", "initialize",
						"err", err,
					)
					os.Exit(1)
				}
			},
		},
	}
)

func (c *initCommand) init() *cobra.Command {
	c.cmd.Flags().StringVarP(
		&to,
		"dir",
		"d",
		".",
		"created template directory path",
	)
	c.cmd.Flags().BoolVar(&verbose, "v", false, "Output verbose.")
	return c.cmd
}
