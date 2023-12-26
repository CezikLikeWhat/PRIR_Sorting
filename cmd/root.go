package cmd

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-sort [command] [flags]", // TODO: Ogarnąć jak zrobić, żeby wyświetlał się tylko ten przykład w --help
	Short: "Go-sort is a cli tool that allows you to quickly sort large files.",
	Long: `Go-Sort is a tool produced for PRiR (UMK) classes.
Its functionalities are:
	- generation of a sample file
	- single-threaded sorting
	- multi-threaded sorting
Implementations can be found at https://github.com/CezikLikeWhat/PRIR_Sorting`,
	Example: "go-sort <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func init() {
	inputCmd.Flags().BoolVarP(&Utils.IsDebugMode, "debug", "d", false, "Run command with debug messages")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// TODO: Ogarnąć może w jaki sposób wyświetlić error ale w kolorze czerwonym lub coś
	}
}
