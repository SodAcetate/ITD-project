package dbhandler

import (
	"context"
	"fmt"
	entry "main/shared/entry"
	"os"

	"github.com/jackc/pgx/v5"
)

// получает на вход объект entry (EntryItem или EntryUser)
// делает запрос к БД
// из БД получает какую-то структурку, преобразует её в entry
// на выход передаёт объекты entry и error
// Исключение -- логика работы с состояниями

// сгенерить предмет-заглушку
func GetPlaceholderItem() entry.EntryItem {
	var placeholderItem entry.EntryItem
	placeholderItem.ID = 1
	placeholderItem.Name = "PlaceholderItem"
	placeholderItem.UserInfo = GetPlaceholderUser()
	return placeholderItem
}

// сгенерить юзера-заглушку
func GetPlaceholderUser() entry.EntryUser {
	var placeholderUser entry.EntryUser
	placeholderUser.ID = 1
	placeholderUser.Name = "PlaceholderUsername"
	placeholderUser.Contact = "OneVVTG"
	return placeholderUser
}

func ConnectDB(coninfo string) pgx.Conn {
	conn, err := pgx.Connect(context.Background(), coninfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection to db failed: %v\n", err)
	} else {
		fmt.Printf("succes")
	}
	return *conn
}

func CloseDB(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func GetUserState(ID int64, conn *pgx.Conn) string {
	var state string
	err := conn.QueryRow(context.Background(), fmt.Sprintf("SELECT state FROM users WHERE user_id = %d", ID)).Scan(&state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return state
	//return "start"
}

func UpdateUserState(new_state string) error {
	return nil
}

func AddItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func EditItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func DeleteItem(item entry.EntryItem) (entry.EntryItem, error) {
	return item, nil
}

func GetAll() ([]entry.EntryItem, error) {

	items := make([]entry.EntryItem, 0)
	items = append(items, GetPlaceholderItem())

	return items, nil
}
