package cmd

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	numberOfThreads       int32
	nameOfInputFileToSort string
	nameOfOutputFile      string
)

var sortCmd = &cobra.Command{
	Use:     "sort [flags]",
	Short:   "Command used to sort the transferred dataset",
	Example: "go-sort sort <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		Utils.Sort(nameOfInputFileToSort, nameOfOutputFile, numberOfThreads)
	},
}

func init() {
	RootCmd.AddCommand(sortCmd)

	sortCmd.Flags().Int32VarP(&numberOfThreads, "threads", "g", int32(runtime.NumCPU()), "Number of threads")
	sortCmd.Flags().StringVarP(&nameOfInputFileToSort, "input", "f", "input.bin", "Name of input file")
	sortCmd.Flags().StringVarP(&nameOfOutputFile, "output", "o", "output.bin", "Name of output file")
}
