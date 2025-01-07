package main

import (
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

	cronjobHelper.SetUpTask()

	select {}

}

//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"os"
//)
//
//type Transaction struct {
//	FromAddress string  `json:"from_address"`
//	Value       float64 `json:"value"`
//}
//
//type Data struct {
//	Data []Transaction `json:"data"`
//}
//
//func main() {
//	// Read the JSON file
//	file, err := os.Open("bb.json")
//	if err != nil {
//		log.Fatalf("Failed to open file: %v", err)
//	}
//	defer file.Close()
//
//	byteValue, _ := ioutil.ReadAll(file)
//
//	var data Data
//	err = json.Unmarshal(byteValue, &data)
//	if err != nil {
//		log.Fatalf("Failed to unmarshal JSON: %v", err)
//	}
//
//	// Filter and print addresses with large values
//	for _, transaction := range data.Data {
//		if transaction.Value > 100 { // Adjust the threshold as needed
//			fmt.Println(transaction.FromAddress)
//		}
//	}
//}
