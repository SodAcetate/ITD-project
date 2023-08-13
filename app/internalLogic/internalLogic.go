package internallogic

import (
	"fmt"
	"log"
	cachehandler "main/app/cacheHandler"
	dbhandler "main/app/dbLogic"
	"main/shared/entry"
	"main/shared/message"
	"strconv"
	"time"
)

type Core struct {
	Db          dbhandler.DbHandler
	Cache       cachehandler.Cache
	MarkupMap   map[string][]string
	finalStates map[string]bool
}

// Инициализация
func (core *Core) Init() {
	core.Db.Init()
	core.Cache.Init()
	core.MarkupMap = map[string][]string{
		"start":              {"Каталог", "Моё", "Поиск"},
		"cat":                {"Выйти", "Поиск"},
		"cat_my":             {"Выйти", "Добавить", "Изменить", "Удалить", "Указать контакты"},
		"edit_item":          {"Изменить имя", "Изменить описание", "Отмена", "Готово"},
		"search":             {"Выйти"},
		"ask_item_name":      {"Отмена"},
		"ask_item_desc":      {"Отмена"},
		"ask_search":         {"Отмена"},
		"ask_contact":        {"Отмена"},
		"delete_item_select": {"Отмена"},
		"edit_item_select":   {"Отмена"},
	}
	core.finalStates = map[string]bool{
		"new":                true,
		"start":              true,
		"ask_item_name":      false,
		"ask_item_desc":      false,
		"ask_contact":        false,
		"edit_item_select":   false,
		"edit_item":          false,
		"delete_item_select": false,
		"cat":                true,
		"ask_search":         false,
		"cat_my":             true,
		"search":             true,
	}
}

func (core *Core) Deinit() {
	core.Db.Deinit()
	core.Cache.Deinit()
}

func (core *Core) GetUserState(ID int64) (string, error) {
	state, ok := core.Cache.GetUserState(ID)
	var err error
	if !ok {
		state, err = core.Db.GetUserState(ID)
	}
	return state, err
}

func (core *Core) UpdateUserState(ID int64, state string) {
	core.Cache.SetUserState(ID, state)
	if core.finalStates[state] {
		core.Db.UpdateUserState(ID, state)
	}
}

// получает на вход ID юзера и иногда сообщение
// в dbHandler передаёт объект entry (EntryItem или EntryUser)
// из dbHandler получает объекты entry и error
// на выход передаёт объекты message и состояние (srtring)

func userToString(user entry.EntryUser) string {
	text := fmt.Sprintf("<i>%s @%s</i>", user.Name, user.Username)
	if user.Contacts != "" {
		text += fmt.Sprintf("\n<i>%s</i>", user.Contacts)
	}
	return text
}

// Получить текстовое представление предмета
func itemToString(item entry.EntryItem, userInfoNeeded bool) string {
	text := fmt.Sprintf("<b>%s</b>\n%s\n", item.Name, item.Desc)
	if userInfoNeeded {
		text += fmt.Sprintf("%s\n", userToString(item.UserInfo))
	}
	return text
}

func catalogueToString(catalogue []entry.EntryItem, header string, userInfoNeeded bool) string {
	text := ""
	if header != "" {
		text = header + "\n"
	}
	for index, item := range catalogue {
		text += fmt.Sprintf("\n<b>%d.</b> %s", index+1, itemToString(item, userInfoNeeded))
	}
	return text
}

func (core *Core) AddUser(ID int64, name, username string) {
	core.Db.AddUser(entry.EntryUser{ID: ID, State: "new", Name: name, Username: username})
}

func (core *Core) EditUser(ID int64, name, username string) {
	core.Db.EditUser(entry.EntryUser{ID: ID, Name: name, Username: username})
}

// Удаляет структуру из кеша
// Возвращает состояние start
func (core *Core) Cancel(ID int64) (message.Message, string) {
	state := "start"
	core.Cache.SetCurrentItem(ID, entry.EntryItem{})
	msg := message.Message{Text: "Операция отменена", Buttons: core.MarkupMap[state]}
	return msg, state
}

func (core *Core) Start(ID int64) (message.Message, string) {
	state := "start"
	msg := message.Message{Text: "Привет! Выбирай действие!\n\n" +
		"<b>Каталог</b> — список всех предметов\n" +
		"<b>Поиск</b> — поиск по названию и описанию\n" +
		"<b>Моё</b> — список твоих предметов. Здесь можно добавлять, удалять и изменять предметы.", Buttons: core.MarkupMap[state]}
	return msg, state
}

func (core *Core) Echo(ID int64, state string, reply string) (message.Message, string) {
	text := "Ой, чёт не то: " + reply
	msg := message.Message{Text: text, Buttons: core.MarkupMap[state]}
	return msg, state
}

// Получить из базы список всех предметов
// Вернуть сообщение с инфой о всех предметах [id] name - name @contact
func (core *Core) GetCatalogue(ID int64) (message.Message, string) {
	state := "cat"
	var msg message.Message
	catalogue, _, isLastPage := core.Db.GetCatalogueFirstPage()

	if len(catalogue) == 0 {
		msg.Text = "Товаров нет! Можете добавить первый"
		msg.Buttons = core.MarkupMap[state]
		return msg, state
	}

	core.Cache.SetCatalogue(ID, catalogue)

	msg.Text = catalogueToString(catalogue, "Каталог", true)
	msg.Buttons = core.MarkupMap[state]

	if isLastPage == false {
		msg.Buttons = append(msg.Buttons, "Вперёд")
	}

	return msg, state
}

func (core *Core) CatNextPage(ID int64) (message.Message, string) {
	state := "cat"
	var msg message.Message

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}
	key := []int64{catalogue[len(catalogue)-1].Updated, catalogue[len(catalogue)-1].ID}
	log.Println(key)

	catalogue, _, isLastPage := core.Db.GetCatalogueNextPage(key[0], key[1])

	if len(catalogue) == 0 {
		msg, state = core.Echo(ID, "start", "Это последняя страница!")
		return msg, state
	} else {
		core.Cache.SetCatalogue(ID, catalogue)
		msg.Text = catalogueToString(catalogue, "", true)
		msg.Buttons = core.MarkupMap[state]
		msg.Buttons = append(msg.Buttons, "Назад")
		if isLastPage == false {
			msg.Buttons = append(msg.Buttons, "Вперёд")
		}

		return msg, state
	}
}

func (core *Core) CatPrevPage(ID int64) (message.Message, string) {
	state := "cat"
	var msg message.Message

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	key := []int64{catalogue[0].Updated, catalogue[0].ID}
	log.Println(key)

	catalogue, _, isFirstPage := core.Db.GetCataloguePrevPage(key[0], key[1])

	if len(catalogue) == 0 {
		msg, state = core.Echo(ID, "start", "Это первая страница!")
		return msg, state
	} else {
		core.Cache.SetCatalogue(ID, catalogue)

		msg.Text = catalogueToString(catalogue, "Каталог", true)
		msg.Buttons = core.MarkupMap[state]
		if isFirstPage == false {
			msg.Buttons = append(msg.Buttons, "Назад")
		}
		msg.Buttons = append(msg.Buttons, "Вперёд")

		return msg, state
	}
}

// поиск по вхождению в название
// запрашивает у юзера подстроку
func (core *Core) SearchInit(ID int64) (message.Message, string) {
	var (
		state string
		msg   message.Message
	)

	state = "ask_search"
	msg.Text = "Введите запрос:"
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// получает повары с подстрокой, пишет их в кеш и на экран
// возвращает состояние start
func (core *Core) Search(ID int64, input string) (message.Message, string) {
	var msg message.Message
	var text string
	state := "search"
	core.Cache.SetInput(ID, input)

	items, err, isLastPage := core.Db.GetSearchFirstPage(input)

	if err != nil {
		return core.Echo(ID, "start", "Непредвиденная ошибка")
	} else if len(items) == 0 {
		text = "Увы, товаров не найдено"
		state = "start"
	} else {
		core.Cache.SetCatalogue(ID, items)
		text = catalogueToString(items, "Результаты поиска:", true)
	}

	msg.Text = text
	msg.Buttons = core.MarkupMap[state]
	if isLastPage == false {
		msg.Buttons = append(msg.Buttons, "Вперёд")
	}

	return msg, state
}

func (core *Core) SearchNextPage(ID int64) (message.Message, string) {
	var msg message.Message
	var text string
	state := "search"
	input, ok := core.Cache.GetInput(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	key := []int64{catalogue[len(catalogue)-1].Updated, catalogue[len(catalogue)-1].ID}
	log.Println(key)

	items, _, isLastPage := core.Db.GetSearchNextPage(key[0], key[1], input)

	if len(items) == 0 {
		msg, state = core.Echo(ID, "start", "Это последняя страница!")
		return msg, state
	} else {
		core.Cache.SetCatalogue(ID, items)
		text = catalogueToString(items, "", true)
	}

	msg.Text = text
	msg.Buttons = core.MarkupMap[state]
	msg.Buttons = append(msg.Buttons, "Назад")
	if isLastPage == false {
		msg.Buttons = append(msg.Buttons, "Вперёд")
	}

	return msg, state
}

func (core *Core) SearchPrevPage(ID int64) (message.Message, string) {
	var msg message.Message
	var text string
	state := "search"
	input, ok := core.Cache.GetInput(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	key := []int64{catalogue[0].Updated, catalogue[0].ID}
	log.Println(key)

	items, _, isLastPage := core.Db.GetSearchPrevPage(key[0], key[1], input)

	if len(items) == 0 {
		msg, state = core.Echo(ID, "start", "Это первая страница!")
		return msg, state
	} else {
		core.Cache.SetCatalogue(ID, items)
		text = catalogueToString(items, "", true)
	}

	msg.Text = text
	msg.Buttons = core.MarkupMap[state]
	if isLastPage == false {
		msg.Buttons = append(msg.Buttons, "Назад")
	}
	msg.Buttons = append(msg.Buttons, "Вперёд")

	return msg, state
}

func (core *Core) GetUsersItems(ID int64) (message.Message, string) {
	var text string
	state := "cat_my"
	var msg message.Message
	catalogue, _ := core.Db.SearchByUser(ID)

	text = userToString(core.Db.GetUserInfo(ID)) + "\n\n"

	if len(catalogue) == 0 {
		text = text + "Товаров нет! Можете добавить первый"
		msg.Buttons = []string{"Выйти", "Добавить", "Указать контакты"}
	} else {
		log.Printf("Test: " + catalogue[0].UserInfo.Name)
		core.Cache.SetCatalogue(ID, catalogue)
		text = text + catalogueToString(catalogue, "<b>Ваши товары:</b>", false)
		msg.Buttons = core.MarkupMap[state]
	}

	msg.Text = text

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

// сюда при состоянии edit_item_wait
// Возвращает кнопку "Отмена" и кнопки для выбора товара для изменения по его ID(ID товара пишется при выборе каталога)
func (core *Core) EditItemSelect(ID int64) (message.Message, string) {
	var msg message.Message
	msg.Text = "Выберите предмет для редактирования"
	state := "edit_item_select"

	items, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	buttons := core.MarkupMap[state]

	for index, item := range items {
		buttons = append(buttons, fmt.Sprintf("%d", index+1))
		log.Printf("EditItemInit: %d", item.ID)
	}
	msg.Buttons = buttons

	return msg, state
}

// вызывает AskItemName для получение нового имени
func (core *Core) EditItemInit(ID int64, input string) (message.Message, string) {
	var msg message.Message

	index, _ := strconv.Atoi(input)

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	log.Printf("Длина каталога %d", len(catalogue))

	if len(catalogue) == 0 {
		catalogue, _ = core.Db.GetAll()
	}

	entry := catalogue[index-1]
	core.Cache.SetCurrentItem(ID, entry)

	msg.Text = fmt.Sprintf("Выбрано: %s", entry.Name)
	state := "edit_item"
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Запрашивает у юзера название предмета
func (core *Core) AskItemName(ID int64, mode string) (message.Message, string) {
	var (
		state string
		msg   message.Message
	)

	state = "ask_item_name"
	msg.Text = "Введите название товара: "
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Запрашивает у юзера описание
// Возвращает состояние add_item_desc
func (core *Core) AskItemDescription(ID int64, mode string) (message.Message, string) {
	var state string
	var msg message.Message

	state = "ask_item_desc"
	msg.Text = "Введите описание товара: "
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) SetItemName(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, ok := core.Cache.GetCurrentItem(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	if len(input) > 60 {
		return core.Echo(ID, "ask_item_name", fmt.Sprintf("длина названия не больше 60 символов, введено: %d", len(input)))
	}

	entry.Name = input
	core.Cache.SetCurrentItem(ID, entry)

	state = "edit_item"

	msg.Text = "Имя успешно добавлено"
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Пишет описание в структуру в кэше
// Пока ограничиваю описание в 256 символов
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) SetItemDescription(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	entry, ok := core.Cache.GetCurrentItem(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	if len(input) > 512 {
		return core.Echo(ID, "add_item_desc", fmt.Sprintf("длина описания не больше 512 символов, введено: %d", len(input)))
	}

	entry.Desc = input
	core.Cache.SetCurrentItem(ID, entry)

	state = "edit_item"

	msg.Text = "Описание успешно добавлено"
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Вызывает dbcontext.AddItem, передаёт готовую структуру из кеша
// Возвращает состояние start
func (core *Core) ItemPost(ID int64) (message.Message, string) {
	var msg message.Message
	state := "start"

	entry, ok := core.Cache.GetCurrentItem(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}
	entry.Updated = time.Now().Unix()

	if entry.ID == 0 {
		core.Db.AddItem(entry)
		msg.Text = "Товар успешно добавлен:\n" + itemToString(entry, true)
	} else {
		core.Db.EditItem(entry)
		msg.Text = "Товар успешно изменён:\n" + itemToString(entry, true)
	}

	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

func (core *Core) DeleteItemSelect(ID int64) (message.Message, string) {
	var msg message.Message
	msg.Text = "Выберите предмет для удаления"
	state := "delete_item_select"

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	buttons := core.MarkupMap[state]

	for index, item := range catalogue {
		buttons = append(buttons, fmt.Sprintf("%d", index+1))
		log.Printf("DeleteItemInit: %d", item.ID)
	}
	msg.Buttons = buttons

	return msg, state
}

func (core *Core) DeleteItem(ID int64, input string) (message.Message, string) {
	state := "start"

	index, _ := strconv.Atoi(input)

	catalogue, ok := core.Cache.GetCatalogue(ID)
	if !ok {
		return core.Echo(ID, "start", "")
	}

	if len(catalogue) == 0 {
		catalogue, _ = core.Db.GetAll()
	}

	entry := catalogue[index-1]
	core.Db.DeleteItem(entry)

	var msg message.Message
	msg.Text = "Товар успешно удалён"
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

func (core *Core) AskContact(ID int64) (message.Message, string) {
	var state string
	var msg message.Message

	state = "ask_contact"
	msg.Text = "Опишите, как с вами можно связаться: "
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}

// Пишет имя в структуру в кэше
// Даёт пользователю кнопки: Изменить имя, Изменить описание, Отмена, Готово
// Возвращает состояние add_item_wait
func (core *Core) SetContact(ID int64, input string) (message.Message, string) {
	var (
		msg   message.Message
		state string
	)

	user := core.Db.GetUserInfo(ID)

	if len(input) > 512 {
		return core.Echo(ID, "ask_contact", fmt.Sprintf("не больше 512 символов, введено: %d", len(input)))
	}

	user.Contacts = input
	core.Db.EditUser(user)

	state = "start"

	msg.Text = "Контактные данные успешно обновлены: \n" + userToString(user)
	msg.Buttons = core.MarkupMap[state]

	return msg, state
}
