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

func AddItem(args string) string {
	msg := "Добавить предмет" + args
	return msg
}

func RemoveItem(args string) string {
	msg := "Удалить предмет" + args
	return msg
}
