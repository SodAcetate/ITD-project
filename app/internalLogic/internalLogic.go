package internallogic

import (
	"fmt"
	"log"
	cachehandler "main/app/cacheHandler"
	dbhandler "main/app/dbLogic"
	"main/shared/entry"
	"main/shared/message"
	"strconv"
)

type Core struct {
	Db    dbhandler.DbHandler
	Cache cachehandler.Cache
}

// Инициализация
func (core *Core) Init() {
	core.Db.Init()
	core.Cache.Init()
}

func (core *Core) Deinit() {
	core.Db.Deinit()
	core.Cache.Deinit()
}

// получает на вход ID юзера и иногда сообщение
// в dbHandler передаёт объект entry (EntryItem или EntryUser)
// из dbHandler получает объекты entry и error
// на выход передаёт объекты message и состояние (srtring)

// Получить текстовое представление предмета
func itemToString(entry.EntryItem) string {
	return ""
}

func (core *Core) AddUser(ID int64, name, contact string) {
	core.Db.AddUser(entry.EntryUser{ID: ID, State: "start", Name: name, Contact: contact})
}

// Удаляет структуру из кеша
// Возвращает состояние start
func (core *Core) Cancel(ID int64) (message.Message, string) {
	msg := message.Message{Text: "Операция отменена", Buttons: []string{"Каталог", "Моё", "Поиск"}}
	state := "start"

	return msg, state
}

func (core *Core) MainMenu(ID int64) (message.Message, string) {
	msg := message.Message{Text: "Привет! Выбирай действие!", Buttons: []string{"Каталог", "Моё", "Поиск"}}
	state := "start"

	return msg, state
}

// Получить из базы список всех предметов
// Вернуть сообщение с инфой о всех предметах [id] name - name @contact
func (core *Core) GetCatalogue(ID int64) (message.Message, string) {
	text := "Каталог"
	state := "cat"
	var msg message.Message
	catalogue, _ := core.Db.GetAll()

	if len(catalogue) == 0 {
		msg.Text = "Товаров нет! Можете добавить первый"
		msg.Buttons = []string{"Назад", "Добавить"}
		return msg, state
	}

	//log.Printf("Test: " + catalogue[0].UserInfo.Name)

	core.Cache.SetCatalogue(ID, catalogue)

	for index, item := range catalogue {
		text += fmt.Sprintf("\n\n[%d] %s \n%s \n%s @%s", index+1, item.Name, item.Desc, item.UserInfo.Name, item.UserInfo.Contact)
	}

	msg.Text = text
	msg.Buttons = []string{"Назад", "Поиск"}

	return msg, state
}

func (core *Core) LookForInit(ID int64) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	// msg.Text = "Режим поиска по названию"
	// msg.Buttons = []string{"Отмена"}
	// state = "ask_item_name"

	msg, state = core.AskItemName(ID, "Поиск по названию.\n", "s")

	return msg, state
}

// Создаёт пустую структурку EntryItem, пишет её в кэш
// Вызывает функцию AskItemName
func (core *Core) AddItemInit(ID int64) (message.Message, string) {
	var (
		msg message.Message
		// new_item entry.EntryItem // кеша нет будет глабольная переменная
		state string
	)

	core.Cache.SetCurrentItem(ID, entry.EntryItem{UserInfo: core.Db.GetUserInfo(ID)})

	msg, state = core.AskItemName(ID, "Добавление названия.\n", "a")
	return msg, state
}

// Запрашивает у юзера название предмета
// Возвращает состояние add_item_name
// 'a' - add, 'e' - edit, 's' - search
func (core *Core) AskItemName(ID int64, beg_text string, mode string) (message.Message, string) {
	var (
		state string
		msg   message.Message
	)
	if mode == "a" {
		state = "add_item_name"
	} else if mode == "s" {
		state = "search_for_name"
	}
	msg.Text = beg_text + "Введите название товара: "
	msg.Buttons = []string{"Отмена"}

	return msg, state
}

func (core *Core) SearchItemName(ID int64, substr string) (message.Message, string) {
	var msg message.Message
	text := "Найденные товары: \n"
	msg.Buttons = []string{"Каталог", "Моё", "Поиск"}
	state := "start"

	items, _ := core.Db.SearchByName(substr)

	if len(items) == 0 {
		text = "Увы, товаров не найдено"
	} else {
		core.Cache.SetCatalogue(ID, items)
		for index, item := range items {
			text += fmt.Sprintf("\n\n[%d] %s \n%s \n%s @%s", index+1, item.Name, item.Desc, item.UserInfo.Name, item.UserInfo.Username)
		}
	}

	msg.Text = text
	return msg, state
}


// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) AddItemName(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, _ := core.Cache.GetCurrentItem(ID)

	if len(input) <= 30 {
		entry.Name = input
		core.Cache.SetCurrentItem(ID, entry)

		msg.Text = "Имя успешно добавлено"
		msg.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "add_item_wait"
	} else {
		msg.Text = "Сорян, длина названия не больше 30 символов"
		msg.Buttons = []string{"Отмена"}
		state = "add_item_name"
	}
	return msg, state
}

// Запрашивает у юзера описание
// Возвращает состояние add_item_desc
func (core *Core) AskItemDescription(ID int64) (message.Message, string) {
	state := "add_item_desc"
	var msg message.Message
	msg.Text = "Введите описание товара: "
	msg.Buttons = []string{"Отмена"}

	return msg, state
}

// Пишет описание в структуру в кэше
// Пока ограничиваю описание в 256 символов
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) AddItemDescription(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, _ := core.Cache.GetCurrentItem(ID)

	if len(input) <= 256 {
		entry.Desc = input
		core.Cache.SetCurrentItem(ID, entry)

		msg.Text = "Описание успешно добавлено"
		msg.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "add_item_wait"
	} else {
		msg.Text = "Сорян, длина описания не больше 256 символов"
		msg.Buttons = []string{"Отмена"}
		state = "add_item_desc"
	}
	return msg, state
}

// Вызывает dbcontext.AddItem, передаёт готовую структуру из кеша
// Возвращает состояние start
func (core *Core) AddItemPost(ID int64) (message.Message, string) {
	state := "start"
	entry, _ := core.Cache.GetCurrentItem(ID)
	core.Db.AddItem(entry)

	var msg message.Message
	msg.Text = "Товар успешно добавлен"
	msg.Buttons = []string{"Каталог"}

	return msg, state
}

func (core *Core) RemoveItemInit(ID int64) (message.Message, string) {
	msg := message.Message{Text: "Удаление пока не работает", Buttons: []string{"Каталог"}}
	state := "start"
	return msg, state
}

// Поиск по чему? По названию? Тогда через AskItemName. Пропишу чутка позже
func (core *Core) LookForItem(ID int64) (message.Message, string) {
}

// сюда при состоянии edit_item_wait
// Возвращает кнопку "Отмена" и кнопки для выбора товара для изменения по его ID(ID товара пишется при выборе каталога)
func (core *Core) EditItemInit(ID int64) (message.Message, string) {
	var msg message.Message
	msg.Text = "Выберите предмет для редактирования"

	catalogue, _ := core.Cache.GetCatalogue(ID)

	buttons := make([]string, len(catalogue)+1)
	buttons = append(buttons, "Отмена")
	for index, item := range catalogue {
		buttons = append(buttons, fmt.Sprintf("%d", index+1))
		log.Print(item.ID)
	}
	msg.Buttons = buttons
	state := "edit_item_select"

	return msg, state
}

func (core *Core) EditItemSelect(ID int64, input string) (message.Message, string) {
	var msg message.Message

	index, _ := strconv.Atoi(input)

	catalogue, _ := core.Cache.GetCatalogue(ID)
	log.Printf("Длина каталога %d", len(catalogue))
	if len(catalogue) == 0 {
		catalogue, _ = core.Db.GetAll()
	}

	entry := catalogue[index-1]
	core.Cache.SetCurrentItem(ID, entry)

	msg.Text = fmt.Sprintf("Выбрано: %s", entry.Name)
	msg.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
	state := "edit_item_wait"

	return msg, state
}

// Запрашивает у юзера название предмета
// Возвращает состояние edit_item_name
func (core *Core) AskItemNameEdit(ID int64) (message.Message, string) {
	state := "edit_item_name"
	var msg message.Message
	msg.Text = "Введите название товара: "
	msg.Buttons = []string{"Отмена"}

	return msg, state
}

// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние edit_item_wait
func (core *Core) EditItemName(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, _ := core.Cache.GetCurrentItem(ID)

	if len(input) <= 30 {
		entry.Name = input
		core.Cache.SetCurrentItem(ID, entry)

		msg.Text = "Имя успешно изменено"
		msg.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "edit_item_wait"
	} else {
		msg.Text = "Сорян, длина названия не больше 30 символов"
		msg.Buttons = []string{"Отмена"}
		state = "edit_item_name"
	}
	return msg, state
}

// Запрашивает у юзера описание
// Возвращает состояние edit_item_desc
func (core *Core) AskItemDescriptionEdit(ID int64) (message.Message, string) {
	state := "edit_item_desc"
	var msg message.Message
	msg.Text = "Введите описание товара: "
	msg.Buttons = []string{"Отмена"}

	return msg, state
}

// Пишет описание в структуру в кэше
// Пока ограничиваю описание в 256 символов
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) EditItemDescription(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, _ := core.Cache.GetCurrentItem(ID)

	if len(input) <= 256 {
		entry.Desc = input
		core.Cache.SetCurrentItem(ID, entry)

		msg.Text = "Описание успешно изменено"
		msg.Buttons = []string{"Изменить имя", "Изменить описание", "Отмена", "Готово"}
		state = "edit_item_wait"
	} else {
		msg.Text = "Сорян, длина описания не больше 256 символов"
		msg.Buttons = []string{"Отмена"}
		state = "edit_item_desc"
	}
	return msg, state
}

func (core *Core) EditItemPost(ID int64) (message.Message, string) {
	state := "start"

	entry, _ := core.Cache.GetCurrentItem(ID)
	core.Db.EditItem(entry)

	var msg message.Message
	msg.Text = "Товар успешно изменен"
	msg.Buttons = []string{"Каталог"}

	return msg, state
}

func (core *Core) DeleteItemInit(ID int64) (message.Message, string) {
	var msg message.Message
	msg.Text = "Выберите предмет для удаления"

	catalogue, _ := core.Cache.GetCatalogue(ID)

	buttons := make([]string, len(catalogue)+1)
	buttons = append(buttons, "Отмена")
	for index, item := range catalogue {
		buttons = append(buttons, fmt.Sprintf("%d", index+1))
		log.Print(item.ID)
	}
	msg.Buttons = buttons
	state := "delete_item_select"

	return msg, state
}

func (core *Core) DeleteItemSelect(ID int64, input string) (message.Message, string) {
	state := "start"

	index, _ := strconv.Atoi(input)

	catalogue, _ := core.Cache.GetCatalogue(ID)
	if len(catalogue) == 0 {
		catalogue, _ = core.Db.GetAll()
	}

	entry := catalogue[index-1]
	core.Db.DeleteItem(entry)

	var msg message.Message
	msg.Text = "Товар успешно удалён"
	msg.Buttons = []string{"Каталог"}

	return msg, state
}

func (core *Core) StartMenu(ID int64, input string) (message.Message, string) {
	state := "start"

	var msg message.Message
	msg.Text = "Выберите действие"
	msg.Buttons = []string{"Каталог", "Моё", "Поиск"}

	return msg, state
}
