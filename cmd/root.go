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

func Initialize() {
	rootCmd.AddCommand(
		initCmd.init(),
		formatCmd.init(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
