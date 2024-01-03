package cmd

import (
	"errors"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	numberOfThreads       int32
	nameOfInputFileToSort string
	nameOfOutputFile      string
)

var sortCmd = &cobra.Command{
	Use:          "sort [flags]",
	Short:        "Command used to sort the transferred dataset.\nCommand is executed concurrently using sample sort and min heap.",
	Example:      "go-sort sort -f input.bin -o output.bin -t 8",
	Version:      "1.0.0",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateSortArgs(cmd); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		Utils.Sort(nameOfInputFileToSort, nameOfOutputFile, numberOfThreads)
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)

	sortCmd.Flags().Int32VarP(&numberOfThreads, "threads", "t", int32(runtime.NumCPU()), "Number of threads")
	sortCmd.Flags().StringVarP(&nameOfInputFileToSort, "file", "f", "input.bin", "Name of input file")
	sortCmd.Flags().StringVarP(&nameOfOutputFile, "output", "o", "output.bin", "Name of output file")
}

func validateSortArgs(cmd *cobra.Command) error {
	threads, err := cmd.Flags().GetInt32("threads")
	if err != nil || threads <= 0 {
		return errors.New("invalid number of threads")
	}

	inputFile, _ := cmd.Flags().GetString("file")
	if _, err := os.Stat(inputFile); errors.Is(err, os.ErrNotExist) {
		return errors.New("failed to detect input file or file does not exist")
	}

	return nil
}
