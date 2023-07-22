package internallogic

import (
	"fmt"
	dbhandler "main/app/dbLogic"
)

func GetCatalogue() string {
	msg := "Каталог"
	items, _ := dbhandler.GetAll()
	for _, item := range items {
		msg = msg + fmt.Sprintf("\n[%d] %s %s", item.ID, item.Name, item.UserInfo.Name)
	}
	return msg
}

func AddItem() string {
	msg := "Добавление пока не работает :("
	return msg
}

func RemoveItem() string {
	msg := "Удаление пока не работает :("
	return msg
}
