package main

import (
	"github.com/komem3/diary/cmd"
	_ "github.com/komem3/diary/statik"
)

func main() {
	cmd.Initialize(cmd.NewInitCommand(), cmd.NewFormatCommand())
	cmd.Execute()
}
