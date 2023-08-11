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
		"ask_item_name":      qHandler.askItemNameHandle,
		"ask_item_desc":      qHandler.askItemDescHandle,
		"ask_contact":        qHandler.AskContactHandle,
		"edit_item_select":   qHandler.editItemSelectHandle,
		"edit_item":          qHandler.editItemHandle,
		"delete_item_select": qHandler.deleteItemSelectHandle,
		"cat":                qHandler.catHandle,
		"search":             qHandler.searchHandler,
		"cat_my":             qHandler.CatMyHandle,
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
	kb_rows := []([]tgbotapi.KeyboardButton){}
	kb_row := []tgbotapi.KeyboardButton{}
	for index, button := range buttons {
		kb_row = append(kb_row, tgbotapi.NewKeyboardButton(button))
		log.Printf("buildMarkup: %s", button)
		if index%3 == 2 || index == len(buttons)-1 {
			kb_rows = append(kb_rows, kb_row)
			kb_row = []tgbotapi.KeyboardButton{}
		}
	}

	kb := tgbotapi.NewReplyKeyboard(kb_rows...)

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
	case "Поиск":
		msg, new_state = qHandler.Core.SearchInit(update.Message.Chat.ID)
	default:
		msg, new_state = qHandler.Core.Echo(update.Message.Chat.ID, "start")
	}

	return msg, new_state
}

// Каталог
func (qHandler *QueryHandler) catHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Поиск":
		msg, new_state = qHandler.Core.SearchInit(update.Message.Chat.ID)
	case "Назад":
		msg, new_state = qHandler.Core.Start(update.Message.Chat.ID)
	default:
		msg, new_state = qHandler.Core.Echo(update.Message.Chat.ID, "cat")
	}

	return msg, new_state
}

func (qHandler *QueryHandler) CatMyHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Добавить":
		msg, new_state = qHandler.Core.AddItemInit(update.Message.Chat.ID)
	case "Изменить":
		msg, new_state = qHandler.Core.EditItemSelect(update.Message.Chat.ID)
	case "Удалить":
		msg, new_state = qHandler.Core.DeleteItemSelect(update.Message.Chat.ID)
	case "Назад":
		msg, new_state = qHandler.Core.Start(update.Message.Chat.ID)
	case "Указать контакты":
		msg, new_state = qHandler.Core.AskContact(update.Message.Chat.ID)
	default:
		msg, new_state = qHandler.Core.Echo(update.Message.Chat.ID, "cat")
	}

	return msg, new_state
}

// Каталог
func (qHandler *QueryHandler) searchHandler(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Отмена":
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	default:
		msg, new_state = qHandler.Core.Search(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Добавление предмета [0] -- ожидание команды
func (qHandler *QueryHandler) editItemHandle(update *tgbotapi.Update) (message.Message, string) {
	var msg message.Message
	var new_state string

	switch update.Message.Text {
	case "Изменить имя":
		msg, new_state = qHandler.Core.AskItemName(update.Message.Chat.ID, "add")
	case "Изменить описание":
		msg, new_state = qHandler.Core.AskItemDescription(update.Message.Chat.ID, "add")
	case "Отмена":
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	case "Готово":
		msg, new_state = qHandler.Core.ItemPost(update.Message.Chat.ID)
	default:
		msg, new_state = qHandler.Core.Echo(update.Message.Chat.ID, "add_item_wait")
	}

	return msg, new_state
}

// Добавление предмета [1] -- имя
func (qHandler *QueryHandler) askItemNameHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.SetItemName(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

// Добавление предмета [2] -- описание
func (qHandler *QueryHandler) askItemDescHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.SetItemDescription(update.Message.Chat.ID, update.Message.Text)
	}

	return msg, new_state
}

func (qHandler *QueryHandler) AskContactHandle(update *tgbotapi.Update) (message.Message, string) {
	var new_state string
	var msg message.Message

	qHandler.Core.EditUser(update.Message.Chat.ID, fmt.Sprintf("%s %s", update.Message.Chat.FirstName, update.Message.Chat.LastName), update.Message.Chat.UserName)

	if update.Message.Text == "Отмена" {
		msg, new_state = qHandler.Core.Cancel(update.Message.Chat.ID)
	} else {
		msg, new_state = qHandler.Core.SetContact(update.Message.Chat.ID, update.Message.Text)
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
		msg, new_state = qHandler.Core.EditItemInit(update.Message.Chat.ID, update.Message.Text)
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
		msg, new_state = qHandler.Core.DeleteItem(update.Message.Chat.ID, update.Message.Text)
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
	response.ParseMode = tgbotapi.ModeHTML
	response.ReplyMarkup = buildMarkup(msg.Buttons)

	qHandler.Core.Db.UpdateUserState(update.Message.Chat.ID, new_state)
	return response
}
