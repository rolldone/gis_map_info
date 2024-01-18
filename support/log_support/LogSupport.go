package log_support

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Init() *os.File {
	currentTime := time.Now()
	path := "./storage/log"
	os.MkdirAll(path, 0755)
	fileName := fmt.Sprint(path, "/", currentTime.Format("2006_01_02"), "_app.log")
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("This is a test log entry")
	log.Println()
	return f
}
