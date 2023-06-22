package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "Taskman is a CLI task manager",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
