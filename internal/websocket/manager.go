package websocket

import (
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	rooms map[uuid.UUID]*Room
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[uuid.UUID]*Room),
	}
}

func (m *Manager) GetOrCreateRoom(roomID uuid.UUID) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, ok := m.rooms[roomID]; ok {
		return room
	}

	room := NewRoom(roomID)
	m.rooms[roomID] = room
	go room.Run()
	return room
}
