package internallogic

import (
	"fmt"
	dbhandler "main/app/dbLogic"
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
	state := "start"
	text := "Добавление пока не работает :("
	var msg message.Message
	msg.Text = text
	// Хардкод временный. Нужно реализовать markupMap.
	msg.Buttons = []string{"Каталог", "Добавить", "Удалить"}
	return msg, state
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
