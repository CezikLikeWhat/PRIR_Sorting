package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:     "sort [flags]",
	Short:   "Command used to sort the transferred dataset",
	Example: "go-sort sort <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sort command")
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)
}
