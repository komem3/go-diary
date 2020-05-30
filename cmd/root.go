package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(commands ...*cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "diary [command]",
		Short: "A tool for managing your diary",
		Long: `Diary is a CLI libray for managing your diary.
This application can format your diary directory, and make index file.
`,
		Args:    cobra.MinimumNArgs(1),
		Version: "0.8.0",
	}
	rootCmd.AddCommand(
		commands...,
	)
	return rootCmd
}
