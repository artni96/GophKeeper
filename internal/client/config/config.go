package config

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	Token         string
	ServerAddress string `json:"server_addr"`
	ClientID      uuid.UUID
	MaxAttempts   uint64
}

func NewConfig() *Config {
	newCfg := &Config{
		ClientID:    uuid.New(),
		MaxAttempts: 3,
	}
	err := newCfg.SetServerAddress()
	if err != nil {
		panic(err)
	}
	return newCfg
}

func (c *Config) SetServerAddress() error {
	file, err := os.ReadFile("./internal/server/config/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	return nil
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
