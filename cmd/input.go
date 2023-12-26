package cmd

import (
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/Utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	numberOfElements     int64
	nameOfInputFile      string
	isMultithreadingMode bool
)

var inputCmd = &cobra.Command{
	Use:     "input [flags]",
	Short:   "Command used to generate an input file for sorting based on the arguments",
	Example: "go-sort input <TODO>",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now() // TODO: Dodać opcje na wyświetlanie tego
		file := Utils.CreateFile(nameOfInputFile)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Panicf("Cannot close input file | Error: %s", err)
			}
		}(file)

		if isMultithreadingMode {
			// MultiThread
			Utils.WriteToFileAsync(file, numberOfElements) // TODO: Sprawdzić, bo coś tam jest odjebane srogo i program się zacina/trwa za długo przez to
			//Utils.ReadFromFileAsync(file, numberOfElements)
			elapsed := time.Since(start)
			fmt.Printf("Took: %s\n", elapsed)
			return
		}

		// Sequential
		Utils.WriteToFile(file, numberOfElements)
		//Utils.ReadFromFile(file, numberOfElements)
		elapsed := time.Since(start)
		fmt.Printf("Took: %s\n", elapsed)
	},
}

func init() {
	RootCmd.AddCommand(inputCmd)

	inputCmd.Flags().Int64VarP(&numberOfElements, "number", "n", 100, "Number of elements to generate in input file")
	inputCmd.Flags().StringVarP(&nameOfInputFile, "file", "f", "input.bin", "Name of binary file to generate")
	inputCmd.Flags().BoolVarP(&isMultithreadingMode, "multithread", "m", false, "Run command on multiple threads")
}
