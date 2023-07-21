package queryhandler

import (
	dbhandler "main/app/dbLogic"
	internallogic "main/app/internalLogic"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// разметка кнопочек в начальном состоянии
var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/cat"),
		tgbotapi.NewKeyboardButton("/add"),
		tgbotapi.NewKeyboardButton("/remove"),
	),
)

// начальное состояние ("главное меню")
func startHandle(update *tgbotapi.Update) tgbotapi.MessageConfig {
	// пустой ответ
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "cat":
			response.Text = internallogic.GetCatalogue()
		case "add":
			response.Text = internallogic.AddItem(update.Message.CommandArguments())
		case "remove":
			response.Text = internallogic.RemoveItem(update.Message.CommandArguments())
		}
	} else {
		response.Text = "HelloWorld!"
		response.ReplyMarkup = startKeyboard
	}
	return response
}

// логика обработки запросов
func Process(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// получаем состояние юзера и вызываем соответствующую ему ручку
	state := dbhandler.GetUserState(update.Message.Chat.ID)
	switch state {
	case "start":
		msg := startHandle(update)
		bot.Send(msg)
	}
}
