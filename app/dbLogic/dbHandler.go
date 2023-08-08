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
	db.Conninfo = os.Getenv("DB_CONN_STRING")
	con, err := pgx.Connect(context.Background(), db.Conninfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to db failed: %v\n", err)
	}
	db.Conn = *con
}

func (db *DbHandler) Deinit() {
	db.Conn.Close(context.Background())
}

func (db *DbHandler) GetUserState(ID int64) (string, error) {
	var state string
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT state FROM users WHERE id = %d", ID)).Scan(&state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUserState failed: %v\n", err)
	}
	return state, err
}

func (db *DbHandler) GetUserInfo(ID int64) entry.EntryUser {
	var user entry.EntryUser
	err := db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT * FROM users WHERE id = %d", ID)).Scan(&user.ID, &user.State, &user.Name, &user.Username, &user.Contacts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUserInfo failed: %v\n", err)
	}
	return user
}

func (db *DbHandler) UpdateUserState(ID int64, new_state string) error {
	_, err := db.Conn.Exec(context.Background(), fmt.Sprintf("UPDATE users SET state='%s' WHERE id=%d", new_state, ID))
	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateUserState failed: %v\n", err)
	}
	return err
}

func (db *DbHandler) AddItem(item entry.EntryItem) (entry.EntryItem, error) {
	log.Println(fmt.Sprintf("%d : Adding item", item.UserInfo.ID))
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO items (user_id, name, description, type) VALUES ($1, $2, $3, $4) RETURNING ID", item.UserInfo.ID, item.Name, item.Desc, item.Type)
	var id int64
	err := row.Scan(&id)
	item.ID = id
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddItem failed: %v\n", err)
	}
	return item, err
}

func (db *DbHandler) AddUser(user entry.EntryUser) error {
	row := db.Conn.QueryRow(context.Background(), "INSERT INTO users (id, state, name, username, contacts) VALUES ($1, $2, $3, $4, $5)", user.ID, user.State, user.Name, user.Username, user.Contacts)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser failed: %v\n", err)
	}
	return err
}

func (db *DbHandler) EditItem(item entry.EntryItem) error {
	_, err := db.Conn.Exec(context.Background(), "Update items SET user_id=$1, name=$2, description=$3, type=$4 WHERE id=$5", item.UserInfo.ID, item.Name, item.Desc, item.Type, item.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EditItem failed: %v\n", err)
	}
	return err
}

func (db *DbHandler) EditUser(item entry.EntryUser) error {
	_, err := db.Conn.Exec(context.Background(), "Update users SET name=$1, username=$2, contacts=$3 WHERE id=$4", item.Name, item.Username, item.Contacts, item.ID)
	return err
}

func (db *DbHandler) DeleteItem(item entry.EntryItem) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM items WHERE id=$1", item.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteItem failed: %v\n", err)
	}
	return err
}

func (db *DbHandler) DeleteUser(item entry.EntryUser) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", item.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteUser failed: %v\n", err)
	}
	return err
}

func (db *DbHandler) GetAll() ([]entry.EntryItem, error) {

	request := "SELECT * FROM items"
	return db.getItems(request)

}

func (db *DbHandler) Search(substring string) ([]entry.EntryItem, error) {

	request := fmt.Sprintf("SELECT * FROM items WHERE name ILIKE '%%%s%%' OR description ILIKE '%%%s%%'", substring, substring)
	return db.getItems(request)

}

func (db *DbHandler) SearchByUser(ID int64) ([]entry.EntryItem, error) {

	request := fmt.Sprintf("SELECT * FROM items WHERE user_id = %d", ID)
	return db.getItems(request)

}

func (db *DbHandler) getItems(request string) ([]entry.EntryItem, error) {
	items := make([]entry.EntryItem, 0)
	item := entry.EntryItem{}

	rows, _ := db.Conn.Query(context.Background(), request)
	for rows.Next() {
		rows.Scan(&item.ID, &item.UserInfo.ID, &item.Name, &item.Desc, &item.Type)
		items = append(items, item)
		log.Printf("Added item %s", item.Name)
	}

	for index := range items {
		log.Printf("Asking for UserInfo on ID %d", items[index].UserInfo.ID)
		items[index].UserInfo = db.GetUserInfo(items[index].UserInfo.ID)
	}

	return items, nil
}
