package cmd

import (
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type formatCommand struct {
	cmd *cobra.Command
}

var (
	formatCmd = &formatCommand{
		&cobra.Command{
			Use:   "format",
			Short: "Format directory",
			Long: `Format command analys and format target directory.
After format directory, this rewrite directory structure to target file.
`,
			Run: func(cmd *cobra.Command, args []string) {
				logger := diary.NewLogger(verbose)
				if orgMode {
					if file == "" {
						file = "./README.org"
					}
					if tempFile == "" {
						tempFile = orgTmp
					}
				} else {
					if file == "" {
						file = "./README.md"
					}
					if tempFile == "" {
						tempFile = mdTmp
					}
				}
				dGen, err := diary.NewDiaryGenerator(file, tempFile, from, to, orgMode, logger)
				if err != nil {
					logger.Error(
						"when", "generate generator",
						"err", err.Error(),
					)
					return
				}
				logger.Debug(
					"msg", "analys start",
					"dir", from,
				)
				fMap := diary.ParseFileMap(from)
				elem := diary.Map2Elem(fMap)
				dGen.FormatDir(fMap).WriteDirTree(elem)

				if dGen.Err != nil {
					logger.Error(
						"when", "format and write dir tree",
						"err", dGen.Err.Error(),
					)
					os.Exit(1)
				}
			},
		},
	}
)

func (c *formatCommand) init() *cobra.Command {
	c.cmd.PersistentFlags().StringVarP(&from, "dir", "d", ".",
		"Analysis directory.",
	)
	c.cmd.PersistentFlags().StringVar(&to, "copyDir", "",
		"Format directory. \nWhen this option is difference from 'dir', all file will copy to 'copyDir'.",
	)
	c.cmd.PersistentFlags().StringVarP(&file, "file", "f", "",
		"Rewrite file.\nDefault value is './README.md.\nIn org mode value is ./README.org",
	)
	c.cmd.PersistentFlags().StringVar(&tempFile, "tmpl", "",
		"Parse template file.\nDefault is "+mdTmp+".\nIn org mode value is "+orgTmp+".",
	)
	c.cmd.PersistentFlags().BoolVar(&orgMode, "org", false,
		"Use org template.",
	)
	c.cmd.Flags().BoolVar(&verbose, "v", false, "Output verbose.")

	return c.cmd
}
