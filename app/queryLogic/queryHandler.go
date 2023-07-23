package queryhandler

import (
	"log"
	dbhandler "main/app/dbLogic"
	internallogic "main/app/internalLogic"
	"main/shared/message"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// получает на вход объект tgbotapi.Update
// в internalLogic передаёт ID юзера и сообщение
// из internalLogic получает объект message и новое состояние (string)
// в бота передаёт сообщение tgbotapi.MessageConfig

// разметка кнопочек в начальном состоянии
var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Каталог"),
		tgbotapi.NewKeyboardButton("Добавить"),
		tgbotapi.NewKeyboardButton("Удалить"),
	),
)

func buildMarkup(buttons []string) tgbotapi.ReplyKeyboardMarkup {
	kb_buttons := []tgbotapi.KeyboardButton{}
	for _, button := range buttons {
		kb_buttons = append(kb_buttons, tgbotapi.NewKeyboardButton(button))
	}
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(kb_buttons...),
	)
	return kb
}

var stateMap = map[string]func(*tgbotapi.Update) (tgbotapi.MessageConfig, string){
	"start": startHandle,
}

// начальное состояние ("главное меню")
func startHandle(update *tgbotapi.Update) (tgbotapi.MessageConfig, string) {
	// пустой ответ
	ID := update.Message.Chat.ID
	response := tgbotapi.NewMessage(ID, "")

	var new_state string
	var msg message.Message

	switch update.Message.Text {
	case "Каталог":
		msg, new_state = internallogic.GetCatalogue(ID)
	case "Добавить":
		msg, new_state = internallogic.AddItemInit(ID)
	case "Удалить":
		msg, new_state = internallogic.RemoveItemInit(ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
		new_state = "start"
	}
	response.Text = msg.Text
	response.ReplyMarkup = buildMarkup(msg.Buttons)
	return response, new_state
}

// логика обработки запросов
func Process(update *tgbotapi.Update) tgbotapi.MessageConfig {
	// получаем состояние юзера и вызываем соответствующую ему ручку
	state := dbhandler.GetUserState(update.Message.Chat.ID)
	log.Printf("ID %d: state %s", update.Message.Chat.ID, state)
	msg, new_state := stateMap[state](update)
	dbhandler.UpdateUserState(new_state)
	return msg
}
