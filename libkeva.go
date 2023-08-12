package libkeva

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type KeyValueStore struct {
	data         map[string]interface{}
	mutex        sync.RWMutex
	savePath     string
	lastSaved    time.Time
	saveInterval time.Duration
}

func NewKeyValueStore(savePath string, saveInterval time.Duration) *KeyValueStore {
	kvStore := &KeyValueStore{
		data:         make(map[string]interface{}),
		savePath:     savePath,
		saveInterval: saveInterval,
	}
	go kvStore.periodicPersist()
	return kvStore
}

func (store *KeyValueStore) Get(key string) (interface{}, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	value, ok := store.data[key]
	return value, ok
}

func (store *KeyValueStore) Set(key string, value interface{}) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.data[key] = value
}

func (store *KeyValueStore) Delete(key string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.data, key)
}

func (store *KeyValueStore) periodicPersist() {
	for {
		time.Sleep(store.saveInterval)
		if time.Since(store.lastSaved) > store.saveInterval {
			store.persist()
		}
	}
}

func (store *KeyValueStore) persist() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	err := store.SaveToFile(store.savePath)
	if err != nil {
		log.Printf("Error during persistence: %v", err)
	}
	store.lastSaved = time.Now()
}

func (store *KeyValueStore) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(store.data)
	if err != nil {
		return err
	}
	return nil
}

func (store *KeyValueStore) LoadFromFile(filename string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}
		file.WriteString("{}")
		file.Close()
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &store.data)
	if err != nil {
		return err
	}
	return nil
}
