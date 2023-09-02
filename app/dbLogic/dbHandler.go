package dbhandler

import (
	"context"
	"fmt"
	"log"
	"main/shared/data"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type DbHandler struct {
	Conn        pgx.Conn
	Conninfo    string
	PageLength  int
	DebugLogger log.Logger
}

func (db *DbHandler) Init() {
	log_path := fmt.Sprintf("%s_%v-%v", os.Getenv("LOG_PATH"), time.Now().Day(), time.Now().Month())
	file, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.Printf("logging err: %v", err)
	db.DebugLogger = *log.New(file, "DEBUG: ", log.Default().Flags())

	db.Conninfo = os.Getenv("DB_CONN_STRING")
	con, err := pgx.Connect(context.Background(), db.Conninfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to db failed: %v\n", err)
	}
	db.Conn = *con
	db.PageLength, _ = strconv.Atoi(os.Getenv("PAGE_LENGTH"))
}

func (db *DbHandler) Deinit() {
	db.Conn.Close(context.Background())
}

func (db *DbHandler) GetUserState(ID int64) (string, error) {

	db.DebugLogger.Printf("DbHandler: GetUserState <- %v", ID)

	var state string
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT state FROM users WHERE id = %d", ID)).Scan(&state)

	db.DebugLogger.Printf("DbHandler: GetUserState -> %v, %v", state, err)

	return state, err
}

func (db *DbHandler) GetUserInfo(ID int64) data.EntryUser {

	db.DebugLogger.Printf("DbHandler: GetUserInfo <- %v", ID)

	var user data.EntryUser
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT * FROM users WHERE id = %d", ID)).Scan(&user.ID, &user.State, &user.Name, &user.Username, &user.Contacts)

	db.DebugLogger.Printf("DbHandler: GetUserState -> %v, %v", user, err)

	return user
}

func (db *DbHandler) UpdateUserState(ID int64, new_state string) error {

	db.DebugLogger.Printf("DbHandler: UpdateUserState <- %v, %v", ID, new_state)

	_, err := db.Conn.Exec(context.Background(), fmt.Sprintf("UPDATE users SET state='%s' WHERE id=%d", new_state, ID))

	db.DebugLogger.Printf("DbHandler: UpdateUserState -> %v", err)

	return err
}

func (db *DbHandler) AddItem(item data.EntryItem) error {

	db.DebugLogger.Printf("DbHandler: AddItem <- %v", item)

	row := db.Conn.QueryRow(context.Background(), "INSERT INTO items (user_id, name, description, type, updated) VALUES ($1, $2, $3, $4, $5) RETURNING ID", item.UserInfo.ID, item.Name, item.Desc, item.Type, item.Updated)
	var id int64
	err := row.Scan(&id)
	item.ID = id

	db.DebugLogger.Printf("DbHandler: AddItem -> %v", err)

	return err
}

func (db *DbHandler) AddUser(user data.EntryUser) error {

	db.DebugLogger.Printf("DbHandler: AddUser <- %v", user)

	row := db.Conn.QueryRow(context.Background(), "INSERT INTO users (id, state, name, username, contacts) VALUES ($1, $2, $3, $4, $5)", user.ID, user.State, user.Name, user.Username, user.Contacts)
	var id int64
	err := row.Scan(&id)

	db.DebugLogger.Printf("DbHandler: AddUser -> %v", err)

	return err
}

func (db *DbHandler) EditItem(item data.EntryItem) error {

	db.DebugLogger.Printf("DbHandler: EditItem <- %v", item)

	_, err := db.Conn.Exec(context.Background(), "Update items SET user_id=$1, name=$2, description=$3, type=$4, updated=$5 WHERE id=$6", item.UserInfo.ID, item.Name, item.Desc, item.Type, item.Updated, item.ID)

	db.DebugLogger.Printf("DbHandler: EditItem -> %v", err)

	return err
}

func (db *DbHandler) EditUser(user data.EntryUser) error {
	db.DebugLogger.Printf("DbHandler: EditUser <- %v", user)

	_, err := db.Conn.Exec(context.Background(), "Update users SET name=$1, username=$2, contacts=$3 WHERE id=$4", user.Name, user.Username, user.Contacts, user.ID)

	db.DebugLogger.Printf("DbHandler: EditUser -> %v", err)
	return err
}

func (db *DbHandler) DeleteItem(item data.EntryItem) error {
	db.DebugLogger.Printf("DbHandler: DeleteItem <- %v", item)

	_, err := db.Conn.Exec(context.Background(), "DELETE FROM items WHERE id=$1", item.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteItem failed: %v\n", err)
	}

	db.DebugLogger.Printf("DbHandler: DeleteItem -> %v", err)
	return err
}

func (db *DbHandler) DeleteUser(user data.EntryUser) error {
	db.DebugLogger.Printf("DbHandler: DeleteUser <- %v", user)

	_, err := db.Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", user.ID)

	db.DebugLogger.Printf("DbHandler: DeleteUser -> %v", err)
	return err
}

func (db *DbHandler) SearchByUser(ID int64) ([]data.EntryItem, error) {
	db.DebugLogger.Printf("DbHandler: SearchByUser <- %v", ID)

	request := fmt.Sprintf("SELECT * FROM items WHERE user_id = %d ORDER BY (updated, id) DESC", ID)
	items, err := db.getItems(request)

	db.DebugLogger.Printf("DbHandler: SearchByUser -> %v, %v", items, err)
	return items, err
}

func (db *DbHandler) GetPage(key data.Key, fwd bool, filter data.ItemFilter) (items []data.EntryItem, isFinal bool, err error) {
	var filterParams, pageParams string

	// Итерируемся по странице вперёд/назад в зависимости от флага
	if fwd && (key != data.Key{}) {
		pageParams = fmt.Sprintf("(updated, id) < (%d, %d)", key.Updated, key.ID)
	} else {
		pageParams = fmt.Sprintf("(updated, id) > (%d, %d)", key.Updated, key.ID)
	}
	// добавляем условия на вхождения подстрок
	if filter.Substrings != nil {
		for _, substring := range filter.Substrings {
			filterParams += fmt.Sprintf(" AND CONCAT(name, ' ', description) ILIKE '%%%s%%'", substring)
		}
	}
	// условие на user_id
	if filter.UserID != 0 {
		filterParams += fmt.Sprintf(" AND user_id=%d", filter.UserID)
	}
	// условие на type
	if filter.ItemType != 0 {
		filterParams += fmt.Sprintf(" AND type=%d", filter.ItemType)
	}

	// создаём запрос
	var request string
	if fwd {
		request = fmt.Sprintf("SELECT * FROM items WHERE %s %s ORDER BY (updated, id) DESC FETCH FIRST %d ROWS ONLY", pageParams, filterParams, db.PageLength)
	} else {
		request = fmt.Sprintf("SELECT * FROM (SELECT * FROM items WHERE %s %s ORDER BY (updated, id) ASC FETCH NEXT %d ROWS ONLY) AS foo ORDER BY (updated, id) DESC", pageParams, filterParams, db.PageLength)
	}

	// получаем вхождения
	items, err = db.getItems(request)

	// считаем, есть ли ещё
	var count int8
	var finalItem data.EntryItem
	if fwd {
		finalItem = items[len(items)-1]
		err = db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM items WHERE (updated, id) < (%d,%d) %s FETCH FIRST 1 ROWS ONLY",
			finalItem.Updated,
			finalItem.ID,
			filterParams)).Scan(&count)
	} else {
		finalItem = items[0]
		err = db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM items WHERE (updated, id) > (%d,%d) %s FETCH FIRST 1 ROWS ONLY",
			finalItem.Updated,
			finalItem.ID,
			filterParams)).Scan(&count)
	}

	if count == 0 {
		isFinal = true
	} else {
		isFinal = false
	}

	return
}

// Работает
func (db *DbHandler) getItems(request string) ([]data.EntryItem, error) {
	db.DebugLogger.Printf("DbHandler: getItems <- %v", request)

	items := make([]data.EntryItem, 0)
	item := data.EntryItem{}

	rows, _ := db.Conn.Query(context.Background(), request)
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.UserInfo.ID, &item.Name, &item.Desc, &item.Type, &item.Updated)
		if err != nil {
			db.DebugLogger.Printf("DbHandler: getItems error : %v", err)
		}
		items = append(items, item)
		db.DebugLogger.Printf("DbHandler: getItems add : %v", item)
	}

	for index := range items {
		items[index].UserInfo = db.GetUserInfo(items[index].UserInfo.ID)
	}

	db.DebugLogger.Printf("DbHandler: getItems -> %v", items)

	return items, nil
}
