package dbhandler

import (
	"context"
	"fmt"
	entry "main/shared/entry"
	"os"

	"github.com/jackc/pgx/v5"
)

type DbHandler struct {
	Conn     pgx.Conn
	Conninfo string
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

func (db *DbHandler) Init() {
	con, err := pgx.Connect(context.Background(), db.Conninfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to db failed: %v\n", err)
	}
	db.Conn = *con
}

func (db *DbHandler) Deinit() {
	db.Conn.Close(context.Background())
}

func (db *DbHandler) GetUserState(ID int64) string {
	var state string
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT state FROM users WHERE user_id = %d", ID)).Scan(&state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return state
	//return "start"
}

func (db *DbHandler) UpdateUserState(ID int64, new_state string) error {
	_, err := db.Conn.Exec(context.Background(), "UPDATE users SET state=$1 WHERE user_id=$2", new_state, ID)
	return err
}

func (db *DbHandler) AddItem(item entry.EntryItem) error {
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO items (user_id, description) VALUES ($1, $2) RETURNING ID", item.UserInfo.ID, item.Name)
	var id int64
	err := row.Scan(&id)
	return err
	//здесь должно быть добавление юзера если по user_id нет данных в таблице users
	//err = db.conn.QueryRow(context.Background(), "SELECT state FROM users WHERE user_id = $1", ID).Scan(&state)
}

func (db *DbHandler) AddUser(item entry.EntryUser) error {
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO users (user_id, user-name, contacts) VALUES ($1, $2, $3) RETURNING user_id", item.ID, item.Name, item.Contact)
	var id int64
	err := row.Scan(&id)
	return err
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
