package store

import (
	"context"
	"sync"
)

type Memory struct {
	mu     sync.Mutex
	nextID int64
	byID   map[int64]Link
	byName map[string]int64
}

func NewMemory() *Memory {
	return &Memory{
		byID:   make(map[int64]Link),
		byName: make(map[string]int64),
	}
}

func (m *Memory) List(_ context.Context) ([]Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]Link, 0, len(m.byID))
	for _, link := range m.byID {
		out = append(out, link)
	}
	return out, nil
}

func (m *Memory) Get(_ context.Context, id int64) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	link, ok := m.byID[id]
	if !ok {
		return Link{}, ErrNotFound
	}
	return link, nil
}

func (m *Memory) GetByShortName(_ context.Context, shortName string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, ok := m.byName[shortName]
	if !ok {
		return Link{}, ErrNotFound
	}
	return m.byID[id], nil
}

func (m *Memory) Create(_ context.Context, originalURL, shortName string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.byName[shortName]; ok {
		return Link{}, ErrConflict
	}

	m.nextID++
	link := Link{
		ID:          m.nextID,
		OriginalURL: originalURL,
		ShortName:   shortName,
	}
	m.byID[link.ID] = link
	m.byName[shortName] = link.ID
	return link, nil
}

func (m *Memory) Update(_ context.Context, id int64, originalURL, shortName string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	link, ok := m.byID[id]
	if !ok {
		return Link{}, ErrNotFound
	}
	if otherID, taken := m.byName[shortName]; taken && otherID != id {
		return Link{}, ErrConflict
	}
	if link.ShortName != shortName {
		delete(m.byName, link.ShortName)
		m.byName[shortName] = id
	}
	link.OriginalURL = originalURL
	link.ShortName = shortName
	m.byID[id] = link
	return link, nil
}

func (m *Memory) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	link, ok := m.byID[id]
	if !ok {
		return ErrNotFound
	}
	delete(m.byName, link.ShortName)
	delete(m.byID, id)
	return nil
}
