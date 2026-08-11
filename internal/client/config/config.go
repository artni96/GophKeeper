package config

import (
	"github.com/google/uuid"
)

type Config struct {
	ServerAddress   string `json:"server_addr"`
	ClientID        uuid.UUID
	MaxAttempts     uint64
	DefaultServAddr string
}

func NewConfig() *Config {
	newCfg := &Config{
		ClientID:        uuid.New(),
		MaxAttempts:     3,
		DefaultServAddr: ":3200",
	}
	return newCfg
}

// State represent the client state metrics.
type State struct {
	Token          string
	ActionAttempts uint64
	IsOnline       bool
	Password       string
	Keys           map[uint64][]byte
	ActiveKey      []byte
	ActiveKeyID    uint64
}

func NewState() *State {
	return &State{ActionAttempts: 0, IsOnline: false, Keys: make(map[uint64][]byte)}
}

func (s *State) HasAttempts() bool {
	if s.ActionAttempts == 3 {
		return false
	}
	return true
}

func (s *State) AddAttempt() {
	s.ActionAttempts++
}

func (s *State) ResetAttempts() {
	s.ActionAttempts = 0
}
