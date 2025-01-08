package main

import (
	//"DegenAIBot/cronjob"
	//"DegenAIBot/helper"
	"DegenAIBot/config"
	"DegenAIBot/cronjob"
	"DegenAIBot/helper"
	"log"
)

func main() {

	cfg, err := config.Load("")
	if err != nil {
		log.Println("Error loading config file: ", err.Error())
	}

	helperAdapter := helper.NewHelper(cfg)

	cronjobHelper := cronjob.NewTask(cfg, helperAdapter)
	cronjobHelper.AutoSendTweetPNLMessage()

	//cronjobHelper.SetUpTask()
	//select {}

}
