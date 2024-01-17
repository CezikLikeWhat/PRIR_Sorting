package cmd

import (
	"errors"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/CezikLikeWhat/PRIR_Sorting/configuration"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	numberOfCPUs          int32
	numberOfBuckets       int32
	nameOfInputFileToSort string
	nameOfOutputFile      string
)

var sortCmd = &cobra.Command{
	Use:          "sort [flags]",
	Short:        "Command used to sort the transferred dataset.\nCommand is executed concurrently using sample sort and min heap.",
	Example:      "go-sort sort -f input.bin -o output.bin -t 8 -b 8",
	Version:      "1.0.0",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateSortArgs(cmd); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := Utils.Sort(nameOfInputFileToSort, nameOfOutputFile, numberOfCPUs, numberOfBuckets)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)

	sortCmd.Flags().Int32VarP(&numberOfCPUs, "cpu", "t", int32(runtime.NumCPU()), "Number of CPU's used in Go environment")
	sortCmd.Flags().Int32VarP(&numberOfBuckets, "buckets", "b", 8, "Number of buckets in sample sort")
	sortCmd.Flags().Int64VarP(&configuration.THRESHOLD, "size", "s", int64(10_000_000), "Size of fragment of data in every goroutine")
	sortCmd.Flags().StringVarP(&nameOfInputFileToSort, "file", "f", "input.bin", "Name of input file")
	sortCmd.Flags().StringVarP(&nameOfOutputFile, "output", "o", "output.bin", "Name of output file")
}

func validateSortArgs(cmd *cobra.Command) error {
	cpus, err := cmd.Flags().GetInt32("cpu")
	if err != nil || cpus <= 0 {
		return errors.New("invalid number of CPU's")
	}

	buckets, err := cmd.Flags().GetInt32("buckets")
	if err != nil || buckets <= 0 {
		return errors.New("invalid number of buckets")
	}

	size, err := cmd.Flags().GetInt64("size")
	if err != nil || size <= 0 {
		return errors.New("invalid size of fragment data")
	}

	inputFile, _ := cmd.Flags().GetString("file")
	if _, err := os.Stat(inputFile); errors.Is(err, os.ErrNotExist) {
		return errors.New("failed to detect input file or file does not exist")
	}

	return nil
}
