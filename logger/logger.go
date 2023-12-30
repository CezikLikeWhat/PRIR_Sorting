package logger

import (
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/configuration"
	"log"
)

func Debug(message string, variables ...any) {
	if configuration.IsDebugMode {
		log.Printf("[DEBUG] %s\n", fmt.Sprintf(message, variables...))
	}
}

func Info(message string, variables ...any) {
	if configuration.IsDebugMode {
		log.Printf("[INFO] %s\n", fmt.Sprintf(message, variables...))
	}
}

func Fatal(message string, variables ...any) {
	log.Fatalf("[FATAL] %s\n", fmt.Sprintf(message, variables...))
}
