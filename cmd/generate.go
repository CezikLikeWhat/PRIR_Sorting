package cmd

import (
	"errors"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
)

var (
	numberOfElements int64
	nameOfInputFile  string
)

var generateCmd = &cobra.Command{
	Use:          "generate [flags]",
	Short:        "Command used to generate an input file for sorting based on the arguments.\nCommand works sequentially.",
	Example:      "go-sort generate -f input.bin -n 1000000",
	Version:      "1.0.0",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateGenerateArgs(cmd); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := Utils.GenerateInputBinaryFile(nameOfInputFile, numberOfElements)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().Int64VarP(&numberOfElements, "number", "n", 100, "Number of elements to generate in input file")
	generateCmd.Flags().StringVarP(&nameOfInputFile, "file", "f", "input.bin", "Name of binary file to generate")
}

func validateGenerateArgs(cmd *cobra.Command) error {
	number, err := cmd.Flags().GetInt64("number")
	if err != nil || number <= 0 {
		return errors.New("invalid number of elements to generate")
	}

	return nil
}
