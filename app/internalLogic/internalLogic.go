package internallogic

import (
	"fmt"
	dbhandler "main/app/dbLogic"
	"main/shared/entry"
	"main/shared/message"
)

var globalItem entry.EntryItem

type Core struct {
	Db dbhandler.DbHandler
}

// Инициализация
func (core *Core) Init() {
	core.Db.Init()
}

func (core *Core) Deinit() {
	core.Db.Deinit()
}

// получает на вход ID юзера и иногда сообщение
// в dbHandler передаёт объект entry (EntryItem или EntryUser)
// из dbHandler получает объекты entry и error
// на выход передаёт объекты message и состояние (srtring)

// Получить текстовое представление предмета
func itemToString(entry.EntryItem) string

// Получить EntryUser
func getUserInfo(ID int64) entry.EntryUser

// Получить из базы список всех предметов
// Вернуть сообщение с инфой о всех предметах [id] name - name @contact
func (core *Core) GetCatalogue(ID int64) (message.Message, string) {
	text := "Каталог\n"
	state := "start"
	items, _ := dbhandler.GetAll()
	for _, item := range items {
		text += fmt.Sprintf("\n[%d] %s - %s @%s", item.ID, item.Name, item.UserInfo.Name, item.UserInfo.Contact)
	}

	var info message.Message
	info.Text = text
	info.Buttons = []string{"Каталог", "Добавить", "Удалить"}

	return info, state
}

// Создаёт пустую структурку EntryItem, пишет её в кэш
// Вызывает функцию AskItemName
func (core *Core) AddItemInit(ID int64) (message.Message, string) {
	var (
		info message.Message
		// new_item entry.EntryItem // кеша нет будет глабольная переменная
		state string
	)
	info, state = core.AskItemName(ID)
	return info, state
}

// Запрашивает у юзера название предмета
// Возвращает состояние add_item_name
func (core *Core) AskItemName(ID int64) (message.Message, string) {
	state := "add_item_name"
	var info message.Message
	info.Text = "Введите название товара: "
	info.Buttons = []string{"Отмена"}

	return info, state
}

// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) AddItemName(ID int64, input string) (message.Message, string) {
	var (
		info message.Message
		state string
	)

	if len(input) <= 30 {
		globalItem.Name = input
		globalItem.UserInfo = getUserInfo(ID)
		info.Text = "Имя успешно добавлено"
		info.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "add_item_wait"
	} else {
		info.Text = "Сорян, длина названия не больше 30 символов"
		info.Buttons = []string{"Отмена"} // может сюда еще отмену добавить?
		state = "add_item_name"
	}
	return info, state
}

// Запрашивает у юзера описание
// Возвращает состояние add_item_desc
func (core *Core) AskItemDescription(ID int64) (message.Message, string) {
	state := "add_item_desk"
	var info message.Message
	info.Text = "Введите описание товара: "
	info.Buttons = []string{"Отмена"}

	return info, state
}

// Пишет описание в структуру в кэше
// Пока ограничиваю описание в 256 символов
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) AddItemDescription(ID int64, input string) (message.Message, string) {
	var (
		info message.Message
		state string
	)

	if len(input) <= 256 {
		globalItem.Desc = input
		info.Text = "Описание успешно добавлено"
		info.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "add_item_wait"
	} else {
		info.Text = "Сорян, длина описания не больше 256 символов"
		info.Buttons = []string{"Отмена"} // может сюда еще отмену добавить?
		state = "add_item_desc"
	}
	return info, state
}

// Удаляет структуру из кеша
// Возвращает состояние start
func (core *Core) AddItemCancel(ID int64) (message.Message, string)

// Вызывает dbcontext.AddItem, передаёт готовую структуру из кеша
// Возвращает состояние start
func (core *Core) AddItemPost(ID int64) (message.Message, string) {
	state := "start"
	dbhandler.AddItem(globalItem)

	var info message.Message
	info.Text = "Товар успешно добавлен"
	info.Buttons = []string{"Каталог", "Добавить", "Удалить"}

	return info, state
}

func (core *Core) RemoveItemInit(ID int64) (message.Message, string)
