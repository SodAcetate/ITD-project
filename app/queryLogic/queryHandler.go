package queryhandler

import (
	"log"
	dbhandler "main/app/dbLogic"
	internallogic "main/app/internalLogic"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// разметка кнопочек в начальном состоянии
var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Каталог"),
		tgbotapi.NewKeyboardButton("Добавить"),
		tgbotapi.NewKeyboardButton("Удалить"),
	),
)

var stateMap = map[string]func(*tgbotapi.Update) (tgbotapi.MessageConfig, string){
	"start": startHandle,
}

var stateMarkups = map[string]tgbotapi.ReplyKeyboardMarkup{
	"start": startKeyboard,
}

// начальное состояние ("главное меню")
func startHandle(update *tgbotapi.Update) (tgbotapi.MessageConfig, string) {
	// пустой ответ
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	new_state := "start"
	switch update.Message.Text {
	case "Каталог":
		response.Text = internallogic.GetCatalogue()
	case "Добавить":
		response.Text = internallogic.AddItem()
	case "Удалить":
		response.Text = internallogic.RemoveItem()
	default:
		response.Text = "HelloWorld!"
		response.ReplyMarkup = startKeyboard
	}
	return response, new_state
}

// логика обработки запросов
func Process(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// получаем состояние юзера и вызываем соответствующую ему ручку
	state := dbhandler.GetUserState(update.Message.Chat.ID)
	log.Printf("ID %d: state %s", update.Message.Chat.ID, state)
	msg, new_state := stateMap[state](update)
	msg.ReplyMarkup = stateMarkups[new_state]
	bot.Send(msg)
}
