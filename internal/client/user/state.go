package user

import (
	"errors"

	"github.com/artni96/GophKeeper/internal/client/model/common"
)

var ErrNoAttemptsLeft = errors.New("no attempts left")

type State struct {
	Token          string
	ActionAttempts uint64
	MaxAttempts    uint64
	DataStorage    DataStorage
	IsOnline       bool
}

func NewUserState() *State {
	extendedDataMap := make(map[uint64]any, 10)
	shortDataList := make([]common.Entity, 0, 10)
	dataStorage := DataStorage{
		ExtendedDataMap: extendedDataMap,
		ShortDataList:   shortDataList,
	}
	userState := &State{
		DataStorage: dataStorage,
		MaxAttempts: 3,
	}
	return userState
}

func (u *State) HasAttempts() bool {
	if u.ActionAttempts == 3 {
		return false
	}
	return true
}

func (u *State) AddAttempt() {
	u.ActionAttempts++
}

func (u *State) ResetAttempts() {
	u.ActionAttempts = 0
}

type DataStorage struct {
	ExtendedDataMap map[uint64]any
	ShortDataList   []common.Entity
}
