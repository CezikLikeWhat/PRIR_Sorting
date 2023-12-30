package cmd

import (
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/configuration"
	"github.com/spf13/cobra"
	"log"
	"time"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if configuration.TimeMeasuring {
			configuration.AppTimer.Start = time.Now()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalf("Cannot display help message | Error: %s", err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if configuration.TimeMeasuring {
			configuration.AppTimer.Elapsed = time.Since(configuration.AppTimer.Start)
			fmt.Printf("Action took: %s\n", configuration.AppTimer.Elapsed)
		}
	},
}

func init() {
	inputCmd.Flags().BoolVarP(&configuration.IsDebugMode, "debug", "d", false, "Run command with debug messages")
	inputCmd.Flags().BoolVarP(&configuration.TimeMeasuring, "time", "t", false, "Run command with time measuring")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// TODO: Ogarnąć może w jaki sposób wyświetlić error ale w kolorze czerwonym lub coś
	}
}
