package cmd

import (
	"fmt"
	"os"

	"github.com/komem3/diary"
	"github.com/spf13/cobra"
)

type formatCommand struct {
	cmd      *cobra.Command
	from     string
	to       string
	file     string
	tempFile string
	orgMode  bool
	verbose  bool
}

func NewFormatCommand() CommandInitializer {
	command := &formatCommand{
		cmd: &cobra.Command{
			Use:   "format",
			Short: "Format directory",
			Long: `Format command analys and format target directory.
After format directory, this rewrite directory structure to target file.
`,
		},
	}
	command.cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := command.Format(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	return command
}

func (c *formatCommand) Init() *cobra.Command {
	c.cmd.PersistentFlags().StringVarP(&c.from, "dir", "d", ".",
		"Analysis directory.",
	)
	c.cmd.PersistentFlags().StringVar(&c.to, "copyDir", "",
		"Format directory. \nWhen this option is difference from 'dir', all file will copy to 'copyDir'.",
	)
	c.cmd.PersistentFlags().StringVarP(&c.file, "file", "f", "",
		"Rewrite file.\nDefault value is './README.md.\nIn org mode value is ./README.org",
	)
	c.cmd.PersistentFlags().StringVar(&c.tempFile, "tmpl", "",
		"Parse template file.\nDefault is "+mdTmp+".\nIn org mode value is "+orgTmp+".",
	)
	c.cmd.PersistentFlags().BoolVar(&c.orgMode, "org", false,
		"Use org template.",
	)
	c.cmd.Flags().BoolVar(&c.verbose, "v", false, "Output verbose.")

	return c.cmd
}

func (c formatCommand) Format() error {
	logger := diary.NewLogger(c.verbose)
	if c.orgMode {
		if c.file == "" {
			c.file = "./README.org"
		}
		if c.tempFile == "" {
			c.tempFile = orgTmp
		}
	} else {
		if c.file == "" {
			c.file = "./README.md"
		}
		if c.tempFile == "" {
			c.tempFile = mdTmp
		}
	}
	dGen, err := diary.NewDiaryGenerator(c.file, c.tempFile, c.from, c.to, c.orgMode, logger)
	if err != nil {
		return fmt.Errorf("generate generator: %w", err)
	}
	logger.Debug(
		"msg", "analys start",
		"dir", c.from,
	)
	fMap := diary.ParseFileMap(c.from)
	elem := diary.Map2Elem(fMap)
	dGen.FormatDir(fMap).WriteDirTree(elem)

	if dGen.Err != nil {
		return fmt.Errorf("format and write dir tree: %w", dGen.Err)
	}
	return nil
}
