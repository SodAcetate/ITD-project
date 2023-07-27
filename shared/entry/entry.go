package entry


type Entry struct {
	ID int64
}

type EntryItem struct {
	Entry
	Name     string
	Desc     string
	UserInfo EntryUser
}

type EntryUser struct {
	Entry
	Name    string
	Contact string
}
