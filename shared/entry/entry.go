package entry

type Entry struct {
	ID int
}

type EntryItem struct {
	Entry
	Name     string
	UserInfo EntryUser
}

type EntryUser struct {
	Entry
	Name string
}
