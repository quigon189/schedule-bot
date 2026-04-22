package storage

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	storage map[string]File
	TTL     time.Duration
}

type File struct {
	Data      []byte
	Mime      string
	Filename  string
	ExpiresAt time.Time
}

func NewMemoryStorage(ttl, cleanupInterval time.Duration) *MemoryStorage {
	storage := &MemoryStorage{
		storage: map[string]File{},
		TTL:     ttl,
	}

	go storage.gc(cleanupInterval)

	return storage
}

func (s *MemoryStorage) PutFile(filename, mime string, data []byte) (string, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New()
	exiresAt := time.Now().Add(s.TTL)
	s.storage[id.String()] = File{
		Filename:  filename,
		Mime:      mime,
		Data:      data,
		ExpiresAt: exiresAt,
	}

	return id.String(), exiresAt
}

func (s *MemoryStorage) GetFile(id string) (File, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, found := s.storage[id]
	return file, found
}

func (s *MemoryStorage) gc(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for k, v := range s.storage {
			if v.ExpiresAt.After(now) {
				delete(s.storage, k)
			}
		}
		s.mu.Unlock()
	}
}
