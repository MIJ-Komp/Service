package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevelNames = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
	currentLogLevel = INFO
	logFile         *os.File
	logger          *log.Logger
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func init() {
	logDir := "logs"
	logFileName := "logs/app.log"
	
	// Create logs directory if it doesn't exist
	os.MkdirAll(logDir, os.ModePerm)

	// Open log file with proper permissions
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	logFile = file

	// Create multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger = log.New(multiWriter, "", 0)

	// Set log level from environment variable
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		setLogLevel(envLevel)
	}
}

func setLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		currentLogLevel = DEBUG
	case "INFO":
		currentLogLevel = INFO
	case "WARN", "WARNING":
		currentLogLevel = WARN
	case "ERROR":
		currentLogLevel = ERROR
	case "FATAL":
		currentLogLevel = FATAL
	default:
		currentLogLevel = INFO
	}
}

// Core logging function
func writeLog(level LogLevel, message string, data map[string]interface{}) {
	if level < currentLogLevel {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	var function string
	if ok {
		// Get just the filename, not the full path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		
		// Get function name
		pc, _, _, ok := runtime.Caller(2)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			parts := strings.Split(funcName, ".")
			if len(parts) > 0 {
				function = parts[len(parts)-1]
			}
		}
	}

	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Level:     logLevelNames[level],
		Message:   message,
		File:      file,
		Line:      line,
		Function:  function,
		Data:      data,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple format if JSON marshaling fails
		logger.Printf("[%s] %s %s", logLevelNames[level], time.Now().Format("2006-01-02 15:04:05"), message)
		return
	}

	logger.Println(string(jsonData))
}

// Public logging functions
func LogDebug(message string) {
	writeLog(DEBUG, message, nil)
}

func LogDebugWithData(message string, data map[string]interface{}) {
	writeLog(DEBUG, message, data)
}

func LogInfo(message string) {
	writeLog(INFO, message, nil)
}

func LogInfoWithData(message string, data map[string]interface{}) {
	writeLog(INFO, message, data)
}

func LogWarning(message string) {
	writeLog(WARN, message, nil)
}

func LogWarningWithData(message string, data map[string]interface{}) {
	writeLog(WARN, message, data)
}

func LogError(message string) {
	writeLog(ERROR, message, nil)
}

func LogErrorWithData(message string, data map[string]interface{}) {
	writeLog(ERROR, message, data)
}

func LogFatal(message string) {
	writeLog(FATAL, message, nil)
	os.Exit(1)
}

func LogFatalWithData(message string, data map[string]interface{}) {
	writeLog(FATAL, message, data)
	os.Exit(1)
}

// Specialized logging functions
func LogDBOperation(operation, query string, args ...interface{}) {
	data := map[string]interface{}{
		"operation": operation,
		"query":     query,
		"args":      args,
	}
	writeLog(DEBUG, "Database operation executed", data)
}

func LogAPIRequest(method, path, ip string, statusCode int, duration time.Duration) {
	data := map[string]interface{}{
		"method":      method,
		"path":        path,
		"ip":          ip,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}
	
	level := INFO
	if statusCode >= 400 && statusCode < 500 {
		level = WARN
	} else if statusCode >= 500 {
		level = ERROR
	}
	
	message := fmt.Sprintf("API Request: %s %s - %d (%dms)", method, path, statusCode, duration.Milliseconds())
	writeLog(level, message, data)
}

// Utility functions
func SetLogLevel(level string) {
	setLogLevel(level)
}

func GetLogLevel() string {
	return logLevelNames[currentLogLevel]
}

func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}