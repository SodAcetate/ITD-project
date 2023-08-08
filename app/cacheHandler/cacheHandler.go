package cachehandler

import (
	"encoding/json"
	"fmt"
	"log"
	cacheentry "main/shared/cacheEntry"
	"main/shared/entry"

	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func (cache *Cache) Init() {
	cache.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (cache *Cache) Deinit() {
	cache.client.Close()
}

func (cache *Cache) AddUser(ID int64) {
	data := cacheentry.CacheEntry{State: "start", CurrentItem: entry.EntryItem{}, Catalogue: nil}
	cache.Set(ID, data)
}

func (cache *Cache) Get(ID int64) (cacheentry.CacheEntry, error) {
	js, err := cache.client.Get(fmt.Sprint(ID)).Result()
	if err != nil {
		cache.AddUser(ID)
	}
	var data cacheentry.CacheEntry
	err = json.Unmarshal([]byte(js), &data)
	return data, err
}

func (cache *Cache) Set(ID int64, data cacheentry.CacheEntry) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cache.client.Set(fmt.Sprint(ID), string(js), 0)
	return nil
}

func (cache *Cache) GetUserState(ID int64) (string, error) {
	data, err := cache.Get(ID)
	return data.State, err
}

func (cache *Cache) SetUserInfo(ID int64, state string) error {
	// Получаем данные из кеша
	data, err := cache.Get(ID)
	if err != nil {
		return err
	}
	// Меняем состояние
	data.State = state
	// Записываем обратно
	err = cache.Set(ID, data)
	return err
}

func (cache *Cache) GetCurrentItem(ID int64) (entry.EntryItem, error) {
	data, err := cache.Get(ID)
	return data.CurrentItem, err
}

func (cache *Cache) SetCurrentItem(ID int64, item entry.EntryItem) error {
	// Получаем данные из кеша
	data, err := cache.Get(ID)
	if err != nil {
		return err
	}
	// Меняем состояние
	data.CurrentItem = item
	// Записываем обратно
	err = cache.Set(ID, data)
	return err
}

func (cache *Cache) GetCatalogue(ID int64) ([]entry.EntryItem, error) {
	data, err := cache.Get(ID)
	return data.Catalogue, err
}

func (cache *Cache) SetCatalogue(ID int64, catalogue []entry.EntryItem) error {
	// Получаем данные из кеша
	data, err := cache.Get(ID)
	if err != nil {
		log.Println("SetCatalogue error: " + err.Error())
		return err
	}
	// Меняем состояние
	data.Catalogue = catalogue
	log.Printf("Set catalogue: len %d", len(catalogue))
	// Записываем обратно
	err = cache.Set(ID, data)
	return err
}

func (cache *Cache) Clear(ID int64) {
	cache.client.Del(fmt.Sprint(ID))
}
