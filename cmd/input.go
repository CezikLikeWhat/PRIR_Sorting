package cmd

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
)

var (
	numberOfElements int64
	nameOfInputFile  string
)

var inputCmd = &cobra.Command{
	Use:     "input [flags]",
	Short:   "Command used to generate an input file for sorting based on the arguments",
	Example: "go-sort input <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		Utils.WriteToFile(nameOfInputFile, numberOfElements)
	},
}

func init() {
	RootCmd.AddCommand(inputCmd)

	inputCmd.Flags().Int64VarP(&numberOfElements, "number", "n", 100, "Number of elements to generate in input file")
	inputCmd.Flags().StringVarP(&nameOfInputFile, "file", "f", "input.bin", "Name of binary file to generate")
}
