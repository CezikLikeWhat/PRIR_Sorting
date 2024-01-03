package cmd

import (
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/configuration"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var RootCmd = &cobra.Command{
	Use:   "go-sort [command] [flags]",
	Short: "Go-sort is a cli tool that allows you to quickly sort large files.",
	Long: `Go-Sort is a tool produced for PRiR (UMK) classes.
Its functionalities are:
	- generation of a sample file
	- single-threaded sorting
	- multi-threaded sorting
	- sorting verification
Implementations can be found at: https://github.com/CezikLikeWhat/PRIR_Sorting`,
	Example:       "go-sort [input | sort | verify] [flags]",
	Version:       "1.0.0",
	SilenceErrors: true,
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
			fmt.Printf("Action took: %fs (%s)\n", configuration.AppTimer.Elapsed.Seconds(), configuration.AppTimer.Elapsed)
		}
	},
}

func init() {

	RootCmd.PersistentFlags().BoolVarP(&configuration.IsDebugMode, "debug", "d", false, "Run command with debug messages")
	RootCmd.PersistentFlags().BoolVarP(&configuration.TimeMeasuring, "clock", "c", false, "Run command with time measuring")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Printf("ERROR")
		color.Unset()
		fmt.Printf(": %s\n", err)
	}
}
