package main

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/cmd"
	"log"
	"os"
)

func main() {
	logfile, err := os.Create("logs/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panicf("Cannot close input file | Error: %s", err)
		}
	}(logfile)

	log.SetOutput(logfile)
	cmd.Execute()
}
