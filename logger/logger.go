package logger

import (
	"log"
	"os"

	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
)

var stdLogger = func() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}()

var fileLogger = func() *log.Logger {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return log.New(file, "", log.LstdFlags)
}()

func Info(msg ...any) {
	stdLogger.Println("[INFO]", msg)
	fileLogger.Println("[INFO]", msg)
}

func Warn(msg ...any) {
	stdLogger.Println("[WARN]", msg)
}

func Err(msg ...any) {
	stdLogger.Println("[ERR]", msg)
	fileLogger.Println("[ERR]", msg)
}

func Debug(msg ...any) {
	if env.String("GIN_MODE", true) != ginger.GIN_MODE_RELEASE {
		stdLogger.Println("[DEBUG]", msg)
	}
}
