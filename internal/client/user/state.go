package user

import (
	"errors"
)

var ErrNoAttemptsLeft = errors.New("no attempts left")

type UserState struct {
	CurrentMenu    string
	InTyping       bool
	Token          string
	ActionAttempts uint64
	MaxAttempts    uint64
}

func (u *UserState) HasAttempts() bool {
	if u.ActionAttempts == 3 {
		return false
	}
	return true
}

func (u *UserState) AddAttempt() {
	u.ActionAttempts++
}

func (u *UserState) ResetAttempts() {
	u.ActionAttempts = 0
}

func NewUserState() *UserState {
	return &UserState{
		MaxAttempts: 3,
	}
}

type DataStorage struct {
	Logins[]
	Texts
}
