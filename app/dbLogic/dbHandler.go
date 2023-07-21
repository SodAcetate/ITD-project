package dbhandler

import (
	entry "main/shared/entry"
)

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
	return placeholderUser
}

func GetUserState(ID int64) string {
	return "start"
}

func AddItem() (entry.EntryItem, error) {
	var item entry.EntryItem
	return item, nil
}

func EditItem() (entry.EntryItem, error) {
	var item entry.EntryItem
	return item, nil
}

func GetAll() ([]entry.EntryItem, error) {

	items := make([]entry.EntryItem, 0)
	items = append(items, GetPlaceholderItem())

	return items, nil
}
