package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Protocol represents the supported remote access mechanism for a session.
type Protocol string

const (
	ProtocolSSH  Protocol = "ssh"
	ProtocolSFTP Protocol = "sftp"
	ProtocolFTP  Protocol = "ftp"
)

// AuthMethod indicates how a user authenticates against a remote system.
type AuthMethod string

const (
	AuthPassword   AuthMethod = "password"
	AuthPrivateKey AuthMethod = "private_key"
)

// Session contains the metadata needed to establish a remote connection. The
// struct deliberately omits secret material such as passwords. These must be
// provided at runtime to avoid storing sensitive data on disk.
type Session struct {
	Alias        string     `json:"alias"`
	Protocol     Protocol   `json:"protocol"`
	Host         string     `json:"host"`
	Port         int        `json:"port"`
	Username     string     `json:"username"`
	AuthMethod   AuthMethod `json:"authMethod"`
	KeyPath      string     `json:"keyPath,omitempty"`
	UseTLS       bool       `json:"useTls,omitempty"`
	Description  string     `json:"description,omitempty"`
	RequiresPass bool       `json:"requiresPass"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// SessionStore provides CRUD operations for session definitions.
type SessionStore struct {
	Sessions map[string]Session `json:"sessions"`
}

// LoadSessions reads the session registry from disk or creates an empty store
// when the file does not yet exist.
func LoadSessions() (*SessionStore, error) {
	path, err := sessionsFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &SessionStore{Sessions: map[string]Session{}}, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to stat sessions file: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read sessions file: %w", err)
	}

	store := &SessionStore{}
	if err := json.Unmarshal(data, store); err != nil {
		return nil, fmt.Errorf("invalid sessions file: %w", err)
	}

	if store.Sessions == nil {
		store.Sessions = map[string]Session{}
	}

	return store, nil
}

// Save writes the session store to disk using pretty printed JSON to ease
// manual inspection.
func (s *SessionStore) Save() error {
	path, err := sessionsFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialise sessions: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write sessions: %w", err)
	}

	return nil
}

// Upsert adds a new session or updates an existing entry. The Alias is used as
// unique identifier and is normalised to lower case.
func (s *SessionStore) Upsert(session Session) Session {
	if s.Sessions == nil {
		s.Sessions = map[string]Session{}
	}

	key := strings.ToLower(session.Alias)
	now := time.Now().UTC()
	session.Alias = key

	if existing, ok := s.Sessions[key]; ok {
		session.CreatedAt = existing.CreatedAt
	} else {
		session.CreatedAt = now
	}

	session.UpdatedAt = now
	s.Sessions[key] = session

	return session
}

// Remove deletes the session with the provided alias. It returns an error when
// the alias does not exist to inform the caller that no change happened.
func (s *SessionStore) Remove(alias string) error {
	key := strings.ToLower(alias)
	if _, ok := s.Sessions[key]; !ok {
		return fmt.Errorf("session '%s' not found", alias)
	}
	delete(s.Sessions, key)
	return nil
}

// Get fetches a session by alias.
func (s *SessionStore) Get(alias string) (Session, bool) {
	session, ok := s.Sessions[strings.ToLower(alias)]
	return session, ok
}

// List returns a deterministic list of sessions sorted by alias. This is used
// for display purposes.
func (s *SessionStore) List() []Session {
	sessions := make([]Session, 0, len(s.Sessions))
	for _, session := range s.Sessions {
		sessions = append(sessions, session)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Alias < sessions[j].Alias
	})

	return sessions
}
