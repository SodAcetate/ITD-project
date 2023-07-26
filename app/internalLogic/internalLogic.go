package internallogic

import (
	"fmt"
	dbhandler "main/app/dbLogic"
	"main/shared/entry"
	"main/shared/message"
)

// получает на вход ID юзера и иногда сообщение
// в dbHandler передаёт объект entry (EntryItem или EntryUser)
// из dbHandler получает объекты entry и error
// на выход передаёт объекты message и состояние (srtring)

func GetCatalogue(ID int64) (message.Message, string) {
	text := "Каталог"
	state := "start"
	items, _ := dbhandler.GetAll()
	for _, item := range items {
		text = text + fmt.Sprintf("\n[%d] %s\n%s @%s", item.ID, item.Name, item.UserInfo.Name, item.UserInfo.Contact)
	}
	var msg message.Message
	msg.Text = text
	// Хардкод временный. Нужно реализовать markupMap.
	msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
	return msg, state
}

func AddItemInit(ID int64) (message.Message, string) {
	new_state := "add_item_name"
	text := "Введите название предмета"
	var msg message.Message
	msg.Text = text
	// Хардкод временный. Нужно реализовать markupMap.
	msg.Buttons = []string{"Отмена"}
	return msg, new_state
}

func AddItemName(ID int64, input string) (message.Message, string) {
	var msg message.Message
	var new_state string

	if len(input) <= 30 {
		var item entry.EntryItem
		item.Name = input
		item.UserInfo = GetUserInfo(ID)
		dbhandler.AddItem(item)

		msg.Text = fmt.Sprintf("\"%s\" успешно добавлен", input)
		msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
		new_state = "start"
	} else {
		msg.Text = "Сори, не больше 30 символов!"
		msg.Buttons = []string{"Отмена"}
		new_state = "add_item_name"
	}
	return msg, new_state

}

func RemoveItemInit(ID int64) (message.Message, string) {
	state := "start"
	text := "Удаление пока не работает :("
	var msg message.Message
	msg.Text = text
	// Хардкод временный. Нужно реализовать markupMap.
	msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
	return msg, state
}

func GetUserInfo(ID int64) entry.EntryUser {
	return dbhandler.GetPlaceholderUser()
}
