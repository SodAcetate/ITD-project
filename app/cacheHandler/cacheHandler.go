package cachehandler

import (
	"log"
	cacheentry "main/shared/cacheEntry"
	"main/shared/data"
	//"github.com/go-redis/redis"
)

type Cache struct {
	storage map[int64]cacheentry.CacheEntry
}

func (cache *Cache) Init() {
	cache.storage = map[int64]cacheentry.CacheEntry{}
}

func (cache *Cache) Deinit() {

}

func (cache *Cache) AddUser(ID int64) {
	entry := cacheentry.CacheEntry{State: "start", CurrentItem: data.EntryItem{}, Catalogue: nil}
	cache.Set(ID, entry)
}

func (cache *Cache) Get(ID int64) (cacheentry.CacheEntry, bool) {
	var entry cacheentry.CacheEntry

	entry, ok := cache.storage[ID]
	log.Printf("CacheHandler : Get : %v, %v", ok, entry)
	if ok != true {
		cache.AddUser(ID)
		entry, _ = cache.storage[ID]
	}

	return entry, ok
}

func (cache *Cache) Set(ID int64, entry cacheentry.CacheEntry) {
	log.Printf("CacheHandler : Set : %v", entry)

	cache.storage[ID] = entry
}

func (cache *Cache) GetInput(ID int64) (string, bool) {
	entry, ok := cache.Get(ID)
	return entry.Input, ok
}

func (cache *Cache) SetInput(ID int64, input string) {
	// Получаем данные из кеша
	entry, _ := cache.Get(ID)
	// Меняем состояние
	entry.Input = input
	// Записываем обратно
	cache.Set(ID, entry)
}

func (cache *Cache) GetUserState(ID int64) (string, bool) {
	entry, ok := cache.Get(ID)
	return entry.State, ok
}

func (cache *Cache) SetUserState(ID int64, state string) {
	// Получаем данные из кеша
	entry, _ := cache.Get(ID)
	// Меняем состояние
	entry.State = state
	// Записываем обратно
	cache.Set(ID, entry)
}

func (cache *Cache) GetCurrentItem(ID int64) (data.EntryItem, bool) {
	entry, ok := cache.Get(ID)
	return entry.CurrentItem, ok
}

func (cache *Cache) SetCurrentItem(ID int64, item data.EntryItem) {
	// Получаем данные из кеша
	entry, _ := cache.Get(ID)
	// Меняем состояние
	entry.CurrentItem = item
	// Записываем обратно
	cache.Set(ID, entry)
}

func (cache *Cache) GetCatalogue(ID int64) ([]data.EntryItem, bool) {
	entry, ok := cache.Get(ID)
	if entry.Catalogue == nil {
		ok = false
	}
	return entry.Catalogue, ok
}

func (cache *Cache) SetCatalogue(ID int64, catalogue []data.EntryItem) {
	// Получаем данные из кеша
	entry, _ := cache.Get(ID)
	// Меняем состояние
	entry.Catalogue = catalogue
	log.Printf("Set catalogue: len %d", len(catalogue))
	// Записываем обратно
	cache.Set(ID, entry)
}

func (cache *Cache) GetFilter(ID int64) (data.ItemFilter, bool) {
	entry, ok := cache.Get(ID)
	return entry.Filter, ok
}

func (cache *Cache) SetFilter(ID int64, filter data.ItemFilter) {
	// Получаем данные из кеша
	entry, _ := cache.Get(ID)
	// Меняем фильтр
	entry.Filter = filter
	// Записываем обратно
	cache.Set(ID, entry)
}

func (cache *Cache) Clear(ID int64) {
	delete(cache.storage, ID)
}

func (cache *Cache) ClearAll() {
	for k := range cache.storage {
		delete(cache.storage, k)
	}
}
