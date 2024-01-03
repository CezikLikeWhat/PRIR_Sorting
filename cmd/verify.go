package cmd

import (
	"errors"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var nameOfFileToVerify string

var verifyCmd = &cobra.Command{
	Use:          "verify [flags]",
	Short:        "Command used to validate sorting of numbers in a binary file.\nCommand works sequentially.",
	Example:      "go-sort verify -f output.bin",
	Version:      "1.0.0",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateVerifyArgs(cmd); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		isSortedCorrectly, err := Utils.VerifySortedFile(nameOfFileToVerify)
		if err != nil {
			return err
		}

		if isSortedCorrectly {
			color.Green("The specified file (%s) is sorted correctly\n", nameOfFileToVerify)
			return nil
		}

		color.Red("The specified file (%s) is NOT sorted correctly\n", nameOfFileToVerify)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&nameOfFileToVerify, "file", "f", "output.bin", "Name of the binary file to check for sorting correctness")
}

func validateVerifyArgs(cmd *cobra.Command) error {
	outputFile, _ := cmd.Flags().GetString("file")
	if _, err := os.Stat(outputFile); errors.Is(err, os.ErrNotExist) {
		return errors.New("failed to detect file or file does not exist")
	}

	return nil
}
