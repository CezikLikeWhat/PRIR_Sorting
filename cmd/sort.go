package cmd

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
)

var (
	numberOfThreads       int32
	nameOfInputFileToSort string
	nameOfOutputFile      string
	threshold             int64
)

var sortCmd = &cobra.Command{
	Use:     "sort [flags]",
	Short:   "Command used to sort the transferred dataset",
	Example: "go-sort sort <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		Utils.Sort(nameOfInputFileToSort, nameOfOutputFile, numberOfThreads, threshold)
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)

	sortCmd.Flags().Int32VarP(&numberOfThreads, "threads", "g", 10, "Number of threads")
	sortCmd.Flags().Int64VarP(&threshold, "threshold", "k", 10000, "Number of threshold")
	sortCmd.Flags().StringVarP(&nameOfInputFileToSort, "input", "f", "input.bin", "Name of input file")
	sortCmd.Flags().StringVarP(&nameOfOutputFile, "output", "o", "output.bin", "Name of output file")
}
