package internallogic

import (
	dbhandler "main/app/dbLogic"
	"main/shared/entry"
	"main/shared/message"
)

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
func ItemToString(entry.EntryItem) string

// Получить из базы список всех предметов
// Вернуть сообщение с инфой о всех предметах
func GetCatalogue(ID int64) (message.Message, string)

// Создаёт пустую структурку EntryItem, пишет её в кэш
// Вызывает функцию AskItemName
func AddItemInit(ID int64) (message.Message, string)

// Запрашивает у юзера название предмета
// Возвращает состояние add_item_name
func AskItemName(ID int64) (message.Message, string)

// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func AddItemName(ID int64, input string) (message.Message, string)

// Запрашивает у юзера описание
// Возвращает состояние add_item_desc
func AskItemDescription(ID int64) (message.Message, string)

// Пишет описание в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func AddItemDescription(ID int64, input string) (message.Message, string)

// Удаляет структуру из кеша
// Возвращает состояние start
func AddItemCancel(ID int64) (message.Message, string)

// Вызывает dbcontext.AddItem, передаёт готовую структуру из кеша
// Возвращает состояние start
func AddItemPost(ID int64) (message.Message, string)

func RemoveItemInit(ID int64) (message.Message, string)

func GetUserInfo(ID int64) entry.EntryUser
