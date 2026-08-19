package common

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateEntityI interface {
	GetUserID() uuid.UUID
	SetNumber(number uint64)
}

type PBEncryptedField struct {
	Value *wrapperspb.BytesValue
	Nonce *wrapperspb.BytesValue
	KeyID *wrapperspb.UInt64Value
}
type EncryptedField struct {
	Value []byte
	Nonce []byte
	KeyID uint64
}
