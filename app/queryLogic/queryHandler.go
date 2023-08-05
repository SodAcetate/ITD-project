package queryhandler

import (
	"fmt"
	"log"
	internallogic "main/app/internalLogic"
	"main/shared/message"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type QueryHandler struct {
	Core     internallogic.Core
	stateMap map[string]func(*tgbotapi.Update) (message.Message, string)
}

func (qHandler *QueryHandler) Init() {
	qHandler.Core.Init()
	qHandler.stateMap = map[string]func(*tgbotapi.Update) (message.Message, string){
		"start":              qHandler.startHandle,
		"add_item_wait":      qHandler.addItemWaitHandle,
		"add_item_name":      qHandler.addItemNameHandle,
		"add_item_desc":      qHandler.addItemDescHandle,
		"edit_item_select":   qHandler.editItemSelectHandle,
		"edit_item_wait":     qHandler.editItemWaitHandle,
		"edit_item_name":     qHandler.editItemNameHandle,
		"edit_item_desc":     qHandler.editItemDescHandle,
		"delete_item_select": qHandler.deleteItemSelectHandle,
		"cat":                qHandler.catHandle,
	}
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

// начальное состояние ("главное меню")
func (qHandler *QueryHandler) startHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Каталог":
		msg, new_state = qHandler.Core.GetCatalogue(update.Message.Chat.ID)
	case "Моё":
		msg, new_state = qHandler.Core.GetUsersItems(update.Message.Chat.ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог", "Моё"}
		new_state = "start"
	}

	return msg, new_state
}

// Каталог
func (qHandler *QueryHandler) catHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Добавить":
		msg, new_state = qHandler.Core.AddItemInit(update.Message.Chat.ID)
	case "Изменить":
		msg, new_state = qHandler.Core.EditItemInit(update.Message.Chat.ID)
	case "Удалить":
		msg, new_state = qHandler.Core.DeleteItemInit(update.Message.Chat.ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог"}
		new_state = "start"
	}

	return msg, new_state
}

// Добавление предмета [0] -- ожидание команды
func (qHandler *QueryHandler) addItemWaitHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Изменить имя":
		msg, new_state = qHandler.Core.AskItemName(update.Message.Chat.ID)
	case "Изменить описание":
		msg, new_state = qHandler.Core.AskItemDescription(update.Message.Chat.ID)
	case "Отмена":
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	case "Готово":
		msg, new_state = qHandler.Core.AddItemPost(update.Message.Chat.ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог"}
		new_state = "start"
	}

	return msg, new_state
}

// Добавление предмета [1] -- имя
func (qHandler *QueryHandler) addItemNameHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.AddItemName(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Добавление предмета [2] -- описание
func (qHandler *QueryHandler) addItemDescHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.AddItemDescription(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Изменение предмета [0] -- выбор
func (qHandler *QueryHandler) editItemSelectHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.EditItemSelect(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Изменение предмета [1] -- ожидание команды
func (qHandler *QueryHandler) editItemWaitHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Изменить имя":
		msg, new_state = qHandler.Core.AskItemNameEdit(update.Message.Chat.ID)
	case "Изменить описание":
		msg, new_state = qHandler.Core.AskItemDescriptionEdit(update.Message.Chat.ID)
	case "Отмена":
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	case "Готово":
		msg, new_state = qHandler.Core.EditItemPost(update.Message.Chat.ID)
	default:
		msg.Text = "HelloWorld!"
		msg.Buttons = []string{"Каталог"}
		new_state = "start"
	}

	return msg, new_state
}

// Изменение предмета [2] -- имя
func (qHandler *QueryHandler) editItemNameHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.EditItemName(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Изменение предмета [3] -- описание
func (qHandler *QueryHandler) editItemDescHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.EditItemDescription(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Удаление предмета
func (qHandler *QueryHandler) deleteItemSelectHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.DeleteItemSelect(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// логика обработки запросов
func (qHandler *QueryHandler) Process(update *tgbotapi.Update) tgbotapi.MessageConfig {
	// получаем айди юзера и состояние
	ID := update.Message.Chat.ID
	state, err := qHandler.Core.Db.GetUserState(ID)
	if err != nil {
		log.Printf("Adding user %s", update.Message.Chat.FirstName)
		qHandler.Core.AddUser(ID, fmt.Sprintf("%s %s", update.Message.Chat.FirstName, update.Message.Chat.LastName), update.Message.Chat.UserName)
		state, err = qHandler.Core.Db.GetUserState(ID)
	}
	log.Printf("ID %d: state %s", ID, state)
	// создаём пустой респонс
	response := tgbotapi.NewMessage(ID, "")
	// вызываем соответствующую ручку
	msg, new_state := qHandler.stateMap[state](update)
	// конвертим Message -> MessageConfig
	response.Text = msg.Text
	response.ReplyMarkup = buildMarkup(msg.Buttons)

	qHandler.Core.Db.UpdateUserState(update.Message.Chat.ID, new_state)
	return response
}
