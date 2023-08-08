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
	Db           dbhandler.DbHandler
	Cache        cachehandler.Cache
	StateHandler map[string][]string
}

// Инициализация
func (core *Core) Init() {
	core.Db.Init()
	core.Cache.Init()
	core.StateHandler = map[string][]string{
		"start":     {"Каталог", "Моё", "Поиск"},
		"cat":       {"Назад", "Поиск"},
		"cat_my":    {"Назад", "Добавить", "Изменить", "Удалить"},
		"edit_item": {"Изменить имя", "Изменить описание", "Отмена", "Готово"},
	}
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
func itemToString(item entry.EntryItem) string {
	return fmt.Sprintf("<b>%s</b>\n%s\n<i>%s @%s</i>\n", item.Name, item.Desc, item.UserInfo.Name, item.UserInfo.Username)
}

func catalogueToString(catalogue []entry.EntryItem, header string) string {
	text := header + "\n"
	for index, item := range catalogue {
		text += fmt.Sprintf("\n<b>%d.</b> %s", index+1, itemToString(item))
	}
	return text
}

func (core *Core) AddUser(ID int64, name, username string) {
	core.Db.AddUser(entry.EntryUser{ID: ID, State: "start", Name: name, Username: username})
}

// Удаляет структуру из кеша
// Возвращает состояние start
func (core *Core) Cancel(ID int64) (message.Message, string) {
	state := "start"
	core.Cache.SetCurrentItem(ID, entry.EntryItem{})
	msg := message.Message{Text: "Операция отменена", Buttons: core.StateHandler[state]}
	return msg, state
}

func (core *Core) Start(ID int64) (message.Message, string) {
	state := "start"
	msg := message.Message{Text: "Привет! Выбирай действие!", Buttons: core.StateHandler[state]}
	return msg, state
}

func (core *Core) Echo(ID int64, state string) (message.Message, string) {
	state = "start"
	msg := message.Message{Text: "Сори чё-то пошло не так", Buttons: core.StateHandler[state]}
	return msg, state
}

// Получить из базы список всех предметов
// Вернуть сообщение с инфой о всех предметах [id] name - name @contact
func (core *Core) GetCatalogue(ID int64) (message.Message, string) {
	state := "cat"
	var msg message.Message
	catalogue, _ := core.Db.GetAll()
	msg.Buttons = core.StateHandler[state]

	if len(catalogue) == 0 {
		msg.Text = "Товаров нет! Можете добавить первый"
		return msg, state
	}

	core.Cache.SetCatalogue(ID, catalogue)

	msg.Text = catalogueToString(catalogue, "Каталог")

	return msg, state
}

// поиск по вхождению в название
// запрашивает у юзера подстроку
func (core *Core) SearchInit(ID int64) (message.Message, string) {
	return core.AskItemName(ID, "search")
}

// получает повары с подстрокой, пишет их в кеш и на экран
// возвращает состояние start
func (core *Core) SearchName(ID int64, input string) (message.Message, string) {
	var msg message.Message
	text := "Найденные товары: \n"

	state := "start"

	items, _ := core.Db.SearchByName(input)

	if len(items) == 0 {
		text = "Увы, товаров не найдено"
	} else {
		core.Cache.SetCatalogue(ID, items)
		for index, item := range items {
			text += fmt.Sprintf("\n\n[%d] %s \n%s \n%s @%s", index+1, item.Name, item.Desc, item.UserInfo.Name, item.UserInfo.Username)
		}
	}

	msg.Text = text
	msg.Buttons = core.StateHandler[state]
	return msg, state
}

func (core *Core) GetUsersItems(ID int64) (message.Message, string) {
	text := "Ваши товары:"
	state := "cat_my"
	var msg message.Message
	catalogue, _ := core.Db.SearchByUser(ID)

	if len(catalogue) == 0 {
		msg.Text = "Товаров нет! Можете добавить первый"
		msg.Buttons = []string{"Назад", "Добавить"}
		return msg, state
	}

	log.Printf("Test: " + catalogue[0].UserInfo.Name)

	core.Cache.SetCatalogue(ID, catalogue)

	for index, item := range catalogue {
		text += fmt.Sprintf("\n\n[%d] %s \n%s \n%s @%s", index+1, item.Name, item.Desc, item.UserInfo.Name, item.UserInfo.Username)
	}

	msg.Text = text
	msg.Buttons = core.StateHandler[state]

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

	msg, state = core.AskItemName(ID, "add")
	return msg, state
}

// Запрашивает у юзера название предмета
func (core *Core) AskItemName(ID int64, mode string) (message.Message, string) {
	var (
		state string
		msg   message.Message
	)

	switch mode {
	case "add":
		state = "add_item_name"
	case "edit":
		state = "edit_item_name"
	case "search":
		state = "search_name"
	}

	msg.Text = "Введите название товара: "
	msg.Buttons = []string{"Отмена"}

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
		msg.Buttons = core.StateHandler["edit_item"]
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
func (core *Core) AskItemDescription(ID int64, mode string) (message.Message, string) {
	var state string
	var msg message.Message

	switch mode {
	case "add":
		state = "add_item_desc"
	case "edit":
		state = "edit_item_desc"
	}

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
		msg.Buttons = core.StateHandler["edit_item"]
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
	msg.Buttons = core.StateHandler[state]

	return msg, state
}

// сюда при состоянии edit_item_wait
// Возвращает кнопку "Отмена" и кнопки для выбора товара для изменения по его ID(ID товара пишется при выборе каталога)
func (core *Core) EditItemInit(ID int64) (message.Message, string) {
	var msg message.Message
	msg.Text = "Выберите предмет для редактирования"

	items, _ := core.Cache.GetCatalogue(ID)

	buttons := []string{}
	buttons = append(buttons, "Отмена")
	for index, item := range items {
		buttons = append(buttons, fmt.Sprintf("%d", index+1))
		log.Printf("EditItemInit: %d", item.ID)
	}
	msg.Buttons = buttons
	state := "edit_item_select"

	return msg, state
}

// вызывает AskItemName для получение нового имени
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
	msg.Buttons = core.StateHandler["edit_item"]
	state := "edit_item_wait"

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
		msg.Buttons = core.StateHandler["edit_item"]
		state = "edit_item_wait"
	} else {
		msg.Text = "Сорян, длина названия не больше 30 символов"
		msg.Buttons = []string{"Отмена"}
		state = "edit_item_name"
	}
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
		msg.Buttons = core.StateHandler["edit_item"]
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
	msg.Buttons = core.StateHandler[state]

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
		log.Printf("DeleteItemInit: %d", item.ID)
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
	msg.Buttons = core.StateHandler[state]

	return msg, state
}
