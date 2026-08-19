package common

import userspb "github.com/artni96/GophKeeper/api/proto/users"

type Entity struct {
	Number      uint64
	Title       string
	Description string
}

type Notification struct {
	ActionType   userspb.ActionType
	EntityType   userspb.EntityType
	EntityNumber uint64
}

type DecryptedField struct {
	EncryptedField string
	KeyID          string
	Nonce          []byte
}
