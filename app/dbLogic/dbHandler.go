package dbhandler

import (
	"context"
	"fmt"
	"log"
	entry "main/shared/entry"
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

func (db *DbHandler) GetUserInfo(ID int64) entry.EntryUser {

	db.DebugLogger.Printf("DbHandler: GetUserInfo <- %v", ID)

	var user entry.EntryUser
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

func (db *DbHandler) AddItem(item entry.EntryItem) error {

	db.DebugLogger.Printf("DbHandler: AddItem <- %v", item)

	row := db.Conn.QueryRow(context.Background(), "INSERT INTO items (user_id, name, description, type, updated) VALUES ($1, $2, $3, $4, $5) RETURNING ID", item.UserInfo.ID, item.Name, item.Desc, item.Type, item.Updated)
	var id int64
	err := row.Scan(&id)
	item.ID = id

	db.DebugLogger.Printf("DbHandler: AddItem -> %v", err)

	return err
}

func (db *DbHandler) AddUser(user entry.EntryUser) error {

	db.DebugLogger.Printf("DbHandler: AddUser <- %v", user)

	row := db.Conn.QueryRow(context.Background(), "INSERT INTO users (id, state, name, username, contacts) VALUES ($1, $2, $3, $4, $5)", user.ID, user.State, user.Name, user.Username, user.Contacts)
	var id int64
	err := row.Scan(&id)

	db.DebugLogger.Printf("DbHandler: AddUser -> %v", err)

	return err
}

func (db *DbHandler) EditItem(item entry.EntryItem) error {

	db.DebugLogger.Printf("DbHandler: EditItem <- %v", item)

	_, err := db.Conn.Exec(context.Background(), "Update items SET user_id=$1, name=$2, description=$3, type=$4, updated=$5 WHERE id=$6", item.UserInfo.ID, item.Name, item.Desc, item.Type, item.Updated, item.ID)

	db.DebugLogger.Printf("DbHandler: EditItem -> %v", err)

	return err
}

func (db *DbHandler) EditUser(user entry.EntryUser) error {
	db.DebugLogger.Printf("DbHandler: EditUser <- %v", user)

	_, err := db.Conn.Exec(context.Background(), "Update users SET name=$1, username=$2, contacts=$3 WHERE id=$4", user.Name, user.Username, user.Contacts, user.ID)

	db.DebugLogger.Printf("DbHandler: EditUser -> %v", err)
	return err
}

func (db *DbHandler) DeleteItem(item entry.EntryItem) error {
	db.DebugLogger.Printf("DbHandler: DeleteItem <- %v", item)

	_, err := db.Conn.Exec(context.Background(), "DELETE FROM items WHERE id=$1", item.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteItem failed: %v\n", err)
	}

	db.DebugLogger.Printf("DbHandler: DeleteItem -> %v", err)
	return err
}

func (db *DbHandler) DeleteUser(user entry.EntryUser) error {
	db.DebugLogger.Printf("DbHandler: DeleteUser <- %v", user)

	_, err := db.Conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", user.ID)

	db.DebugLogger.Printf("DbHandler: DeleteUser -> %v", err)
	return err
}

func (db *DbHandler) GetAll() ([]entry.EntryItem, error) {
	db.DebugLogger.Printf("DbHandler: GetAll")

	request := "SELECT * FROM items ORDER BY (updated, id) DESC"
	items, err := db.getItems(request)

	db.DebugLogger.Printf("DbHandler: GetAll -> %v, %v", items, err)
	return items, err
}

func (db *DbHandler) GetCatalogueFirstPage() ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetCatalogueFirstPage")

	params := ""
	items, err, isLastPage := db.firstPage(params)

	db.DebugLogger.Printf("DbHandler: GetCatalogueFirstPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage

}

func (db *DbHandler) GetCatalogueNextPage(key_upd, key_id int64) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetCatalogueNextPage <- %v, %v", key_upd, key_id)

	params := ""
	items, err, isLastPage := db.nextPage(key_upd, key_id, params)

	db.DebugLogger.Printf("DbHandler: GetCatalogueFirstPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage
}

func (db *DbHandler) GetCataloguePrevPage(key_upd, key_id int64) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetCataloguePrevPage <- %v, %v", key_upd, key_id)

	params := ""
	items, err, isFirstPage := db.prevPage(key_upd, key_id, params)

	db.DebugLogger.Printf("DbHandler: GetCataloguePrevPage -> %v, %v, %v", items, err, isFirstPage)
	return items, err, isFirstPage
}

func (db *DbHandler) GetSearchFirstPage(substring string) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetSearchFirstPage <- %v", substring)

	params := fmt.Sprintf("name ILIKE '%%%s%%' OR description ILIKE '%%%s%%'", substring, substring)
	items, err, isLastPage := db.firstPage(params)

	db.DebugLogger.Printf("DbHandler: GetSearchFirstPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage
}

func (db *DbHandler) GetSearchNextPage(key_upd, key_id int64, substring string) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetSearchNextPage <- %v, %v, %v", key_upd, key_id, substring)

	params := fmt.Sprintf("AND (name ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", substring, substring)
	items, err, isLastPage := db.nextPage(key_upd, key_id, params)

	db.DebugLogger.Printf("DbHandler: GetSearchNextPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage
}

func (db *DbHandler) GetSearchPrevPage(key_upd, key_id int64, substring string) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: GetSearchPrevPage <- %v, %v, %v", key_upd, key_id, substring)

	params := fmt.Sprintf("AND (name ILIKE '%%%s%%' OR description ILIKE '%%%s%%')", substring, substring)
	items, err, isFirstPage := db.prevPage(key_upd, key_id, params)

	db.DebugLogger.Printf("DbHandler: GetSearchPrevPage -> %v, %v, %v", items, err, isFirstPage)
	return items, err, isFirstPage
}

func (db *DbHandler) SearchByUser(ID int64) ([]entry.EntryItem, error) {
	db.DebugLogger.Printf("DbHandler: SearchByUser <- %v", ID)

	request := fmt.Sprintf("SELECT * FROM items WHERE user_id = %d ORDER BY (updated, id) DESC", ID)
	items, err := db.getItems(request)

	db.DebugLogger.Printf("DbHandler: SearchByUser -> %v, %v", items, err)
	return items, err
}

func (db *DbHandler) firstPage(params string) ([]entry.EntryItem, error, bool) {
	var request string
	if params != "" {
		request = fmt.Sprintf("SELECT * FROM items WHERE %s ORDER BY (updated, id) DESC FETCH FIRST %d ROWS ONLY", params, db.PageLength)
	} else {
		request = fmt.Sprintf("SELECT * FROM items ORDER BY (updated, id) DESC FETCH FIRST %d ROWS ONLY", db.PageLength)
	}

	db.DebugLogger.Printf("DbHandler: firstPage : %v", request)

	items, err := db.getItems(request)
	if len(items) == 0 {
		db.DebugLogger.Printf("DbHandler: firstPage : len[items] == 0 -> %v, %v, %v", items, err, true)
		return items, err, true
	}

	var isLastPage bool

	var count int8

	if params != "" {
		params = "AND (" + params + ")"
	}

	err = db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM items WHERE (updated, id) < (%d, %d) %s FETCH FIRST 1 ROWS ONLY",
		items[len(items)-1].Updated,
		items[len(items)-1].ID,
		params)).Scan(&count)

	if count == 0 {
		isLastPage = true
	} else {
		isLastPage = false
	}

	db.DebugLogger.Printf("DbHandler: firstPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage
}

func (db *DbHandler) nextPage(key_upd, key_id int64, params string) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: nextPage <- %v, %v", key_upd, key_id)

	var isLastPage bool
	request := fmt.Sprintf("SELECT * FROM items WHERE (updated, id) < (%d, %d) %s ORDER BY (updated, id) DESC FETCH FIRST %d ROWS ONLY", key_upd, key_id, params, db.PageLength)

	db.DebugLogger.Printf("DbHandler: nextPage : %v", request)

	items, err := db.getItems(request)
	if len(items) == 0 {
		db.DebugLogger.Printf("DbHandler: nextPage : len[items] == 0 -> %v, %v, %v", items, err, true)
		return items, err, true
	}

	var count int8
	err = db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM items WHERE (updated, id) < (%d, %d) %s FETCH FIRST 1 ROWS ONLY",
		items[len(items)-1].Updated,
		items[len(items)-1].ID,
		params)).Scan(&count)

	if count == 0 {
		isLastPage = true
	} else {
		isLastPage = false
	}

	db.DebugLogger.Printf("DbHandler: nextPage -> %v, %v, %v", items, err, isLastPage)
	return items, err, isLastPage
}

func (db *DbHandler) prevPage(key_upd, key_id int64, params string) ([]entry.EntryItem, error, bool) {
	db.DebugLogger.Printf("DbHandler: prevPage <- %v, %v", key_upd, key_id)

	var isFirstPage bool
	request := fmt.Sprintf("SELECT * FROM (SELECT * FROM items WHERE (updated, id) > (%d, %d) %s FETCH NEXT %d ROWS ONLY) AS foo ORDER BY (updated, id) DESC", key_upd, key_id, params, db.PageLength)

	db.DebugLogger.Printf("DbHandler: prevPage : %v", request)

	items, err := db.getItems(request)
	if len(items) == 0 {
		db.DebugLogger.Printf("DbHandler: prevPage : len[items] == 0 -> %v, %v, %v", items, err, true)
		return items, err, true
	}

	var count int8
	err = db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM items WHERE (updated, id) > (%d, %d) %s FETCH FIRST 1 ROWS ONLY",
		items[0].Updated,
		items[0].ID,
		params)).Scan(&count)

	if count == 0 {
		isFirstPage = true
	} else {
		isFirstPage = false
	}

	db.DebugLogger.Printf("DbHandler: nextPage -> %v, %v, %v", items, err, isFirstPage)
	return items, err, isFirstPage
}

func (db *DbHandler) getItems(request string) ([]entry.EntryItem, error) {
	db.DebugLogger.Printf("DbHandler: getItems <- %v", request)

	items := make([]entry.EntryItem, 0)
	item := entry.EntryItem{}

	rows, _ := db.Conn.Query(context.Background(), request)
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.UserInfo.ID, &item.Name, &item.Desc, &item.Type, &item.Updated)
		if err != nil {
			db.DebugLogger.Printf("DbHandler: getItems error : %v", err)
		}
		items = append(items, item)
	}

	for index := range items {
		items[index].UserInfo = db.GetUserInfo(items[index].UserInfo.ID)
	}

	db.DebugLogger.Printf("DbHandler: getItems -> %v", items)

	return items, nil
}
