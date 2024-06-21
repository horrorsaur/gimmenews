package utils

import (
	"log"
	"os"
	"path/filepath"
)

// Creates (truncates) fileName at logPath
func NewLogFile(fileName string, logPath string) *os.File {
	logFilePath := filepath.Join(logPath, fileName)
	file, err := os.Create(logFilePath)
	if err != nil {
		log.Printf("couldnt create log file: %s", err)
	}
	log.Printf("created log file '%s' at '%s'", fileName, logPath)
	return file
}
