package main

import (
	"github.com/CezikLikeWhat/PRIR_Sorting/cmd"
	"log"
	"os"
)

func main() {
	logfile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Cannot open file %s (Error: %s)", "logs/app.log", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Cannot close input file (Error: %s)", err)
		}
	}(logfile)

	log.SetOutput(logfile)
	cmd.Execute()
}
