package entry

type EntryItem struct {
	ID       int64
	Name     string
	Desc     string
	UserInfo EntryUser
}

type EntryUser struct {
	ID      int64
	Name    string
	Contact string
}
