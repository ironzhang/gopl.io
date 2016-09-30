package mem_session_map

import (
	"ac-mqtt/mq/session"
	"sync"
)

var _ session.Map = &msmap{}

type msmap struct {
	mu sync.RWMutex
	st map[string]*session.Session
}

func (m *msmap) Set(s *session.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.st[s.Id()] = s
	return nil
}

func (m *msmap) Get(id string) (*session.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.st[id]
	return s, ok
}

func (m *msmap) Del(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.st, id)
}

func (m *msmap) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.st)
}

func mapFactoryFunc(param string) (session.Map, error) {
	return &msmap{st: make(map[string]*session.Session)}, nil
}

func init() {
	session.Register("mem-session-map", mapFactoryFunc)
}
