package main

import (
	"log"

	"github.com/kspiros/xlib"

	"github.com/joho/godotenv"
)

const (
	connPort = "3000"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	//TODO: remove later
	err := godotenv.Load("../../.env")
	//TODO: need for email notif when something goes wrong
	//emailServer := infra.NewMailSender("", "", "")

	//Create logger
	logger, tidylogger := xlib.NewLogger()
	defer tidylogger()

	//Create redis
	_, err = xlib.NewMemCash()
	if err != nil {
		logger.Fatal(err)
		return err
	}

	//Setup db connection
	_, dbtidy, err := xlib.InitDatabase()
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer dbtidy()

	return nil
}
