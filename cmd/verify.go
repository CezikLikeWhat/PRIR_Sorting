package cmd

import (
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
)

var nameOfFileToVerify string

var verifyCmd = &cobra.Command{
	Use:     "verify [flags]",
	Short:   "Command used to validate sorting of numbers in a binary file",
	Example: "go-sort verify <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		isSortedCorrectly := Utils.VerifySortedFile(nameOfFileToVerify)

		if isSortedCorrectly {
			fmt.Printf("The specified file (%s) is sorted correctly\n", nameOfFileToVerify)
			return
		}

		fmt.Printf("The specified file (%s) is NOT sorted correctly\n", nameOfFileToVerify)
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&nameOfFileToVerify, "file", "f", "output.bin", "Name of the binary file to check for sorting correctness")
}
