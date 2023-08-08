package entry

type EntryItem struct {
	ID       int64
	Name     string
	Desc     string
	UserInfo EntryUser
	Type     int8
	Updated  int64
}

type EntryUser struct {
	ID       int64
	State    string
	Name     string
	Username string
	Contacts string
}
