package cachehandler

import (
	"log"
	cacheentry "main/shared/cacheEntry"
	"main/shared/entry"
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
	data := cacheentry.CacheEntry{State: "start", CurrentItem: entry.EntryItem{}, Catalogue: []entry.EntryItem{{}}}
	cache.Set(ID, data)
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
	data, ok := cache.Get(ID)
	return data.Input, ok
}

func (cache *Cache) SetInput(ID int64, input string) {
	// Получаем данные из кеша
	data, _ := cache.Get(ID)
	// Меняем состояние
	data.Input = input
	// Записываем обратно
	cache.Set(ID, data)
}

func (cache *Cache) GetUserState(ID int64) (string, bool) {
	data, ok := cache.Get(ID)
	return data.State, ok
}

func (cache *Cache) SetUserState(ID int64, state string) {
	// Получаем данные из кеша
	data, _ := cache.Get(ID)
	// Меняем состояние
	data.State = state
	// Записываем обратно
	cache.Set(ID, data)
}

func (cache *Cache) GetCurrentItem(ID int64) (entry.EntryItem, bool) {
	data, ok := cache.Get(ID)
	return data.CurrentItem, ok
}

func (cache *Cache) SetCurrentItem(ID int64, item entry.EntryItem) {
	// Получаем данные из кеша
	data, _ := cache.Get(ID)
	// Меняем состояние
	data.CurrentItem = item
	// Записываем обратно
	cache.Set(ID, data)
}

func (cache *Cache) GetCatalogue(ID int64) ([]entry.EntryItem, bool) {
	data, ok := cache.Get(ID)
	return data.Catalogue, ok
}

func (cache *Cache) SetCatalogue(ID int64, catalogue []entry.EntryItem) {
	// Получаем данные из кеша
	data, _ := cache.Get(ID)
	// Меняем состояние
	data.Catalogue = catalogue
	log.Printf("Set catalogue: len %d", len(catalogue))
	// Записываем обратно
	cache.Set(ID, data)
}

func (cache *Cache) Clear(ID int64) {
	delete(cache.storage, ID)
}

func (cache *Cache) ClearAll() {
	for k := range cache.storage {
		delete(cache.storage, k)
	}
}
