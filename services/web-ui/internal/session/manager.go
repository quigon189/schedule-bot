package session

import (
	"encoding/gob"
	"net/http"
	"web-ui/internal/models"

	"github.com/gorilla/sessions"
)

type SessionManager struct {
	store sessions.Store
	name string
}

func NewSessionManager(store sessions.Store, name string) *SessionManager {
	gob.Register(models.JWT{})
	return &SessionManager{
		store: store,
		name: name,
	}
}

func (m *SessionManager) Set(w http.ResponseWriter, r *http.Request, key string, value any) error {
	session, err := m.store.Get(r, m.name)
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r,w)
}

func (m *SessionManager) Get(r *http.Request, key string) any {
	session, _ := m.store.Get(r, m.name)
	return session.Values[key]
}

func (m *SessionManager) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := m.store.Get(r, m.name)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	return sessions.Save(r, w)
}

func (m *SessionManager) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := m.store.Get(r, m.name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return sessions.Save(r, w)
}
