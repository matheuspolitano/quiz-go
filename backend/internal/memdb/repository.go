package memdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultDBFilename = "%s.data.json"
	baseDir           = "./data"
)

var (
	ErrNotFound     = errors.New("entity not found")
	ErrInvalidEntry = errors.New("invalid Entry")
)

// Identifiable enforces a GetID() method so the generic repository
// knows which string key to use in the `entries` map.
type Identifiable interface {
	GetID() string
}

// Repository Generic repository
type Repository[T Identifiable] struct {
	entries  map[string]*T
	filePath string
	mu       sync.RWMutex
}

// NewRepository creates and returns a new Repository for type
func NewRepository[T Identifiable](filePath string) (*Repository[T], error) {
	repo := &Repository[T]{
		entries:  make(map[string]*T),
		filePath: filePath,
	}

	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load data from file: %s", err.Error())
	}
	return repo, nil
}

// NewRepositoryDefault creates and returns a new Repository for type
// collectionName will be the json filename
func NewRepositoryDefault[T Identifiable](collectionName string) (*Repository[T], error) {
	repo := &Repository[T]{
		entries:  make(map[string]*T),
		filePath: filepath.Join(baseDir, fmt.Sprintf(defaultDBFilename, collectionName)),
	}

	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load data from file: %s", err.Error())
	}
	return repo, nil
}

// createEmptyFile creates an empty JSON file (which starts as an empty array).
func (r *Repository[T]) createEmptyFile() error {
	if err := os.MkdirAll(filepath.Dir(r.filePath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories: %s", err.Error())
	}

	file, err := os.Create(r.filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err.Error())
	}
	defer file.Close()

	if _, err := file.Write([]byte("[]")); err != nil {
		return fmt.Errorf("error initializing empty JSON file: %s", err.Error())
	}
	return nil
}

// loadFromFile reads the slice of T from file, populates r.entries map.
func (r *Repository[T]) loadFromFile() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("Loading data from '%s'", r.filePath)

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if createErr := r.createEmptyFile(); createErr != nil {
				return fmt.Errorf("unable to create data file: %s", createErr.Error())
			}
			return nil
		}
		return fmt.Errorf("unable to open data file: %s", err.Error())
	}
	defer file.Close()

	var items []T
	decoder := json.NewDecoder(file)
	if decodeErr := decoder.Decode(&items); decodeErr != nil {
		if decodeErr.Error() == "EOF" {
			return nil
		}
		return fmt.Errorf("error decoding JSON from file: %s", decodeErr.Error())
	}

	for _, item := range items {
		iCopy := item
		r.entries[iCopy.GetID()] = &iCopy
	}
	return nil
}

// FindByID retrieves an entity by its ID.
func (r *Repository[T]) FindByID(id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.entries[id]
	if !exists {
		var empty T
		return empty, ErrNotFound
	}

	return *entry, nil
}

// ListAll retrieves all items
func (r *Repository[T]) ListAll() ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var responseList []T
	for _, item := range r.entries {
		responseList = append(responseList, *item)
	}
	return responseList, nil
}

// Save adds or updates the entity in the repository.
// The entity must have a valid ID, otherwise it returns ErrInvalidEntry.
func (r *Repository[T]) Save(entity T) error {
	id := (entity).GetID()
	if id == "" {
		return fmt.Errorf("%w: empty ID", ErrInvalidEntry)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	copied := entity
	r.entries[id] = &copied

	if err := r.saveToFile(); err != nil {
		return fmt.Errorf("failed to save entity to file: %w", err)
	}

	return nil
}

// saveToFile writes the slice of T from r.entries to disk.
func (r *Repository[T]) saveToFile() (err error) {
	// Prepare a slice for JSON encoding
	items := make([]T, 0, len(r.entries))
	for _, entry := range r.entries {
		items = append(items, *entry)
	}

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open data file for writing: %s", err.Error())
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("error encoding JSON to file: %s", err.Error())
	}

	return nil
}
