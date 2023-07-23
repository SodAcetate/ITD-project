package dbhandler

import (
	entry "main/shared/entry"
)

// получает на вход объект entry (EntryItem или EntryUser)
// делает запрос к БД
// из БД получает какую-то структурку, преобразует её в entry
// на выход передаёт объекты entry и error
// Исключение -- логика работы с состояниями

// сгенерить предмет-заглушку
func GetPlaceholderItem() entry.EntryItem {
	var placeholderItem entry.EntryItem
	placeholderItem.ID = 1
	placeholderItem.Name = "PlaceholderItem"
	placeholderItem.UserInfo = GetPlaceholderUser()
	return placeholderItem
}

// сгенерить юзера-заглушку
func GetPlaceholderUser() entry.EntryUser {
	var placeholderUser entry.EntryUser
	placeholderUser.ID = 1
	placeholderUser.Name = "PlaceholderUsername"
	placeholderUser.Contact = "OneVVTG"
	return placeholderUser
}

func GetUserState(ID int64) string {
	return "start"
}

func UpdateUserState(new_state string) error {
	return nil
}

func AddItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func EditItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func DeleteItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func GetAll() ([]entry.EntryItem, error) {

	items := make([]entry.EntryItem, 0)
	items = append(items, GetPlaceholderItem())

	return items, nil
}
