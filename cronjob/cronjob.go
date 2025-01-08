package cronjob

import (
	"DegenAIBot/config"
	"DegenAIBot/helper"
	"github.com/robfig/cron"
	"log"
	random "math/rand"
	"time"
)

type TaskService interface {
	SetUpTask()
	AutoSendTweetPNLMessage()
}

type Task struct {
	cfg    config.Config
	helper helper.HelperInterface
}

func NewTask(cfg config.Config, helper helper.HelperInterface) TaskService {
	return &Task{
		cfg:    cfg,
		helper: helper,
	}
}

func (t *Task) SetUpTask() {

	var err error

	c := cron.New()

	err = c.AddFunc("@every 3h0m00s", t.AutoSendTweetPNLMessage)
	if err != nil {
		log.Fatal("Error adding cron job:", err)
	}

	c.Start()
}

func (t *Task) AutoSendTweetPNLMessage() {

	randSource := random.NewSource(time.Now().UnixNano()) // Correct usage of NewSource
	randGenerator := random.New(randSource)               // Create a new random generator

	randomNum := randGenerator.Intn(len(t.cfg.SolanaAddresses))

	address := t.cfg.SolanaAddresses[randomNum]

	_, err := t.helper.AddAddressTransactions(address)
	if err != nil {
		log.Println("Error getting address status: ", err.Error())
		return
	}

	pnl, err := t.helper.GetAddressPNL(address)
	if err != nil {
		log.Println("Error getting address PNL: ", err.Error())
		return
	}

	message, err := t.helper.CalculatePortfolio(address)
	if err != nil {
		log.Println("Error getting user portfolio: ", err.Error())
		return
	}

	err = t.helper.SendTweet(message)
	if err != nil {
		log.Println("Error sending tweet: ", err.Error())
		return
	}

	time.Sleep(time.Duration(20) * time.Second)

	err = t.helper.SendTweet(pnl)
	if err != nil {
		log.Println("Error sending tweet: ", err.Error())
		return
	}

}
