package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "Taskman is a CLI task manager",
}
