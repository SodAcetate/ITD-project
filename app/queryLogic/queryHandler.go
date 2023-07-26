package queryhandler

import (
	"log"
	internallogic "main/app/internalLogic"
	"main/shared/message"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type QueryHandler struct {
	Core internallogic.Core
}

func (qHandler *QueryHandler) Init() {
	qHandler.Core.Init()
}

func (qHandler *QueryHandler) Deinit() {
	qHandler.Core.Deinit()
}

// получает на вход объект tgbotapi.Update
// в internalLogic передаёт ID юзера и сообщение
// из internalLogic получает объект message и новое состояние (string)
// в бота передаёт сообщение tgbotapi.MessageConfig

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

var stateMap = map[string]func(*tgbotapi.Update) (message.Message, string){
	"start":         startHandle,
	"add_item_name": addItemNameHandle,
}

// начальное состояние ("главное меню")
func startHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Каталог":
		msg, new_state = internallogic.GetCatalogue(update.Message.Chat.ID)
	case "Добавить":
		msg, new_state = internallogic.AddItemInit(update.Message.Chat.ID)
	case "Удалить":
		msg, new_state = internallogic.RemoveItemInit(update.Message.Chat.ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
		new_state = "start"
	}

	return msg, new_state
}

// Добавление предмета [1] -- имя
func addItemNameHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		return startHandle(update)
	}

	msg, new_state = internallogic.AddItemName(update.Message.Chat.ID, update.Message.Text)

	return msg, new_state
}

// логика обработки запросов
func (qHandler *QueryHandler) Process(update *tgbotapi.Update) tgbotapi.MessageConfig {
	// получаем айди юзера и состояние
	ID := update.Message.Chat.ID
	state := qHandler.Core.Db.GetUserState(ID)
	log.Printf("ID %d: state %s", ID, state)
	// создаём пустой респонс
	response := tgbotapi.NewMessage(ID, "")
	// вызываем соответствующую ручку
	msg, new_state := stateMap[state](update)
	// конвертим Message -> MessageConfig
	response.Text = msg.Text
	response.ReplyMarkup = buildMarkup(msg.Buttons)

	qHandler.Core.Db.UpdateUserState(update.Message.Chat.ID, new_state)
	return response
}
