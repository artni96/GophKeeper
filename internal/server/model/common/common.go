package common

import "google.golang.org/protobuf/types/known/wrapperspb"

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
