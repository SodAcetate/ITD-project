package cacheentry

import "main/shared/data"

type CacheEntry struct {
	Input       string
	State       string
	CurrentItem data.EntryItem
	Catalogue   []data.EntryItem
	Filter      data.ItemFilter
}
