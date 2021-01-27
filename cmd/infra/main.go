package main

import (
	"infra"
	"log"

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
	logger, tidylogger := infra.NewLogger()
	defer tidylogger()

	//Create redis
	//memcash, err := infra.NewMemCash()
	//if err != nil {
	//	logger.Fatal(err)
	//	return err
	//	}

	//Setup db connection
	dbconn, dbtidy, err := infra.InitDatabase()
	if err != nil {
		logger.Fatal(err)
		return err
	}
	defer dbtidy()

	//SetupDB
	err = infra.SetupDB(dbconn)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	return nil
}
