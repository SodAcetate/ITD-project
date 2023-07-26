package main

import (
	"log"
	queryhandler "main/app/queryLogic"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	token := os.Args[1]
	bot, err := tgbotapi.NewBotAPI(token)
	log.Println("Connected via token " + token)
	if err != nil {
		log.Fatal(err)
	}

	var qHandler queryhandler.QueryHandler
	qHandler.Init()
	defer qHandler.Deinit()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		// Process -- логика обработки запросов
		log.Printf("Update from %d [%s]", update.Message.Chat.ID, update.Message.Chat.UserName)
		msg := qHandler.Process(&update)
		bot.Send(msg)
	}

}
