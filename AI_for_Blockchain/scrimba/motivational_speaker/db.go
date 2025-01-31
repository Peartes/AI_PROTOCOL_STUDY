package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// MemStore represents an in-memory key-value store with file persistence
type MemStore struct {
	mu       sync.RWMutex
	data     map[string]string
	filePath string
}

// NewMemStore initializes a new MemStore, loading data from the local file if available
func NewMemStore(filePath string) *MemStore {
	store := &MemStore{
		data:     make(map[string]string),
		filePath: filePath,
	}

	// Load existing data from file if it exists
	store.loadFromFile()
	return store
}

// Set adds or updates a key-value pair
func (m *MemStore) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	m.saveToFile()
}

// Get retrieves a value by key
func (m *MemStore) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.data[key]
	return val, exists
}

// Delete removes a key from the store
func (m *MemStore) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	m.saveToFile()
}

// List returns all key-value pairs
func (m *MemStore) List() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copyData := make(map[string]string)
	for k, v := range m.data {
		copyData[k] = v
	}
	return copyData
}

// saveToFile persists data to the local file system
func (m *MemStore) saveToFile() {
	dataJSON, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		fmt.Println("Error saving data:", err)
		return
	}
	err = os.WriteFile(m.filePath, dataJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

// loadFromFile loads data from the local file system
func (m *MemStore) loadFromFile() {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return // File doesn't exist, start with empty data
	}

	dataJSON, err := os.ReadFile(m.filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(dataJSON, &m.data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
}
