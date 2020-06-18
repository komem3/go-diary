package main

import (
	"fmt"
	"os"

	"github.com/komem3/go-diary/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	command := cmd.NewRootCommand(cmd.NewInitCommand(), cmd.NewFormatCommand(), cmd.NewCommand())
	if err := doc.GenMarkdownTree(command, "."); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
