package dbhandler

import (
	entry "main/shared/entry"
)

type DbHandler struct {
}

func (db *DbHandler) Init() {

}

func (db *DbHandler) Deinit() {

}

// получает на вход объект entry (EntryItem или EntryUser)
// делает запрос к БД
// из БД получает какую-то структурку, преобразует её в entry
// на выход передаёт объекты entry и error
// Исключение -- логика работы с состояниями

// сгенерить предмет-заглушку
func (db *DbHandler) GetPlaceholderItem() entry.EntryItem {
	var placeholderItem entry.EntryItem
	placeholderItem.ID = 1
	placeholderItem.Name = "PlaceholderItem"
	placeholderItem.UserInfo = db.GetPlaceholderUser()
	return placeholderItem
}

// сгенерить юзера-заглушку
func (db *DbHandler) GetPlaceholderUser() entry.EntryUser {
	var placeholderUser entry.EntryUser
	placeholderUser.ID = 1
	placeholderUser.Name = "PlaceholderUsername"
	placeholderUser.Contact = "OneVVTG"
	return placeholderUser
}

func (db *DbHandler) GetUserState(ID int64) string {
	return "start"
}

func (db *DbHandler) UpdateUserState(ID int64, new_state string) error {
	return nil
}

func (db *DbHandler) AddItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func (db *DbHandler) EditItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func (db *DbHandler) DeleteItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func (db *DbHandler) GetAll() ([]entry.EntryItem, error) {

	items := make([]entry.EntryItem, 0)
	items = append(items, db.GetPlaceholderItem())

	return items, nil
}
