package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	logFile := "logs/app.log"
	os.MkdirAll("logs", os.ModePerm)

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogDBOperation(operation, query string, args ...interface{}) {
	logMessage := fmt.Sprintf("DB Operation: %s, Query: %s, Args: %v", operation, query, args)
	InfoLogger.Println(logMessage)
}

func LogError(err error) {
	ErrorLogger.Printf("Error occurred: %v", err)
}

func LogWarning(message string) {
	WarningLogger.Println(message)
}

func LogInfo(message string) {
	InfoLogger.Println(message)
}

func LogAPIRequest(method, path, ip string, statusCode int, duration time.Duration) {
	logMessage := fmt.Sprintf("Method: %s, Path: %s, IP: %s, Status: %d, Duration: %v", 
		method, path, ip, statusCode, duration)
	InfoLogger.Println(logMessage)
}