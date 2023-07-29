package cacheentry

import "main/shared/entry"

type CacheEntry struct {
	State       string
	CurrentItem entry.EntryItem
	Catalogue   []entry.EntryItem
}
