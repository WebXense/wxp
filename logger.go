package wxp

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

func LogInfo(msg ...any) {
	stdLogger.Println("[INFO]", msg)
	fileLogger.Println("[INFO]", msg)
}

func LogErr(msg ...any) {
	stdLogger.Println("[ERR]", msg)
	fileLogger.Println("[ERR]", msg)
}

func LogDebug(msg ...any) {
	if env.String("GIN_MODE", true) != ginger.GIN_MODE_RELEASE {
		stdLogger.Println("[DEBUG]", msg)
	}
}
