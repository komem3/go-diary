package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "diary [command]",
		Short: "A tool for managing your diary",
		Long: `Diary is a CLI libray for managing your diary.
This application can format your diary directory, and make index file.
`,
		Args:    cobra.MinimumNArgs(1),
		Version: "0.8.0",
	}
)

type CommandInitializer interface {
	Init() *cobra.Command
}

func Initialize(cmdi ...CommandInitializer) {
	commands := make([]*cobra.Command, 0, len(cmdi))
	for _, c := range cmdi {
		commands = append(commands, c.Init())
	}
	rootCmd.AddCommand(
		commands...,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
