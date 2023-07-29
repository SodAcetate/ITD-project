package dbhandler

import (
	"context"
	"fmt"
	"log"
	entry "main/shared/entry"
	"os"

	"github.com/jackc/pgx/v5"
)

type DbHandler struct {
	Conn     pgx.Conn
	Conninfo string
}

func (db *DbHandler) Init() {
	db.Conninfo = "postgresql://postgres:123@localhost:5432/postgres"
	con, err := pgx.Connect(context.Background(), db.Conninfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to db failed: %v\n", err)
	}
	db.Conn = *con
}

func (db *DbHandler) Deinit() {
	db.Conn.Close(context.Background())
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
	placeholderUser.Name = ""
	placeholderUser.Contact = "dorm_market_bot"
	return placeholderUser
}

func (db *DbHandler) GetUserState(ID int64) (string, error) {
	var state string
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT state FROM users WHERE user_id = %d", ID)).Scan(&state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return state, err
	//return "start"
}

func (db *DbHandler) GetUserInfo(ID int64) entry.EntryUser {
	var user entry.EntryUser
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT * FROM users WHERE user_id = %d", ID)).Scan(&user.ID, &user.State, &user.Name, &user.Contact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return user
	//return "start"
}

func (db *DbHandler) UpdateUserState(ID int64, new_state string) error {
	_, err := db.Conn.Exec(context.Background(), "UPDATE users SET state=$1 WHERE user_id=$2", new_state, ID)
	return err
}

func (db *DbHandler) AddItem(item entry.EntryItem) (entry.EntryItem, error) {
	log.Println(fmt.Sprintf("%d : Adding item", item.UserInfo.ID))
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO items (name, user_id, description) VALUES ($1, $2, $3) RETURNING ID", item.Name, item.UserInfo.ID, item.Desc)
	var id int64
	err := row.Scan(&id)
	item.ID = id
	return item, err
	//здесь должно быть добавление юзера если по user_id нет данных в таблице users
	//err = db.conn.QueryRow(context.Background(), "SELECT state FROM users WHERE user_id = $1", ID).Scan(&state)
}

func (db *DbHandler) AddUser(user entry.EntryUser) (entry.EntryUser, error) {
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO users (user_id, state, username, contacts) VALUES ($1, $2, $3, $4) RETURNING user_id", user.ID, user.State, user.Name, user.Contact)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Добавление юзера: %e", err)
	}
	return user, err
}

func (db *DbHandler) EditItem(item entry.EntryItem) error {
	_, err := db.Conn.Exec(context.Background(), "Update items SET user_id=$1, name=$2, description=$3 WHERE id=$4", item.UserInfo.ID, item.Name, item.Desc, item.ID)
	return err
}

func (db *DbHandler) EditUser(item entry.EntryUser) error {
	_, err := db.Conn.Exec(context.Background(), "Update users SET username=$1, contacts=$2 WHERE user_id=$3", item.Name, item.Contact, item.ID)
	return err
}

func (db *DbHandler) DeleteItem(item entry.EntryItem) error {
	if item.ID == 0 {
		err := fmt.Errorf("Не указан ID")
		return err
	}
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM items WHERE id=$1", item.ID)
	return err
}

func (db *DbHandler) DeleteUser(item entry.EntryUser) error {
	if item.ID == 0 {
		err := fmt.Errorf("Не указан ID")
		return err
	}
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM users WHERE user_id=$1", item.ID)
	return err
}

func (db *DbHandler) GetAll() ([]entry.EntryItem, error) {

	items := make([]entry.EntryItem, 0)
	item := entry.EntryItem{}

	rows, _ := db.Conn.Query(context.Background(), fmt.Sprintf("SELECT * FROM items"))
	for rows.Next() {
		//entry, _ := pgx.RowToStructByPos[entry.EntryItem](rows)
		//items = append(items, entry)
		rows.Scan(&item.ID, &item.UserInfo.ID, &item.Name, &item.Desc)
		items = append(items, item)
		log.Printf("Added item %s", item.Name)
	}

	for index := range items {
		items[index].UserInfo = db.GetUserInfo(items[index].UserInfo.ID)
	}

	return items, nil
}
