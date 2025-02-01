package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JsonTypes interface {
	string | int | float64 | bool | interface{}
}

// MemStore represents an in-memory key-value store with file persistence
type MemStore[V JsonTypes] struct {
	mu       sync.RWMutex
	data     map[string]map[string]V
	filePath string
}

// NewMemStore initializes a new MemStore, loading data from the local file if available
func NewMemStore[V JsonTypes](filePath string) *MemStore[V] {
	store := &MemStore[V]{
		data:     make(map[string]map[string]V),
		filePath: filePath,
	}

	// Load existing data from file if it exists
	store.loadFromFile()
	return store
}

// Set adds or updates a key-value pair to a table
func (m *MemStore[V]) Set(key string, value V, table string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, tableExists := m.data[table]; !tableExists {
		m.data[table] = make(map[string]V)
	}
	m.data[table][key] = value
	m.saveToFile()
}

// Get retrieves a value by key
func (m *MemStore[V]) Get(key, table string) (*V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, tableExists := m.data[table]; !tableExists {
		return nil, false
	}
	val, exists := m.data[table][key]
	return &val, exists
}

// Delete removes a key from the store
func (m *MemStore[V]) Delete(key, table string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, tableExists := m.data[table]; !tableExists {
		return
	}
	delete(m.data[table], key)
	m.saveToFile()
}

// List returns all key-value pairs
func (m *MemStore[V]) List(table string) map[string]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, tableExists := m.data[table]; !tableExists {
		return map[string]V{}
	}
	copyData := make(map[string]V)
	for k, v := range m.data[table] {
		copyData[k] = v
	}
	return copyData
}

// saveToFile persists data to the local file system
func (m *MemStore[V]) saveToFile() {
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
func (m *MemStore[V]) loadFromFile() {
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
