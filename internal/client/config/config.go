package config

import (
	"github.com/google/uuid"
)

type Config struct {
	Token         string
	ServerAddress string
	ClientID      uuid.UUID
	MaxAttempts   uint64
}

func NewConfig() *Config {
	return &Config{
		ClientID:    uuid.New(),
		MaxAttempts: 3,
	}
}

// State represent the client state metrics.
type State struct {
	ActionAttempts uint64
	IsOnline       bool
}

func NewState() *State {
	return &State{ActionAttempts: 0, IsOnline: false}
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
