package xlib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//ILogger logging interface
type ILogger interface {
	Fatal(msg interface{})
	Close()
}

type logger struct {
	log *log.Logger
}

func (l *logger) Fatal(msg interface{}) {
	fmt.Println(msg)
	l.log.Println(msg)
}

func (l *logger) Close() {
	l.Close()
}

//NewLogger Create Logger to report errors
func NewLogger() (ILogger, func(), error) {
	lp := os.Getenv("LOGGER_FILENAME")
	if len(lp) == 0 {
		lp = "./log/log.log"
	}
	d := filepath.Dir(lp)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		os.Mkdir(d, os.ModePerm)
	}
	openLogfile, err := os.OpenFile(lp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil, err
	}
	l := log.New(openLogfile, "Error Auth:\t", log.Ldate|log.Ltime|log.Lshortfile)
	return &logger{log: l}, func() { openLogfile.Close() }, nil
}
