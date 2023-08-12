package cacheentry

import "main/shared/entry"

type CacheEntry struct {
	Input       string
	State       string
	CurrentItem entry.EntryItem
	Catalogue   []entry.EntryItem
}
