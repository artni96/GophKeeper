package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"

	"google.golang.org/grpc/metadata"
)

func PrepareMDContext(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("authorization", token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

// DecryptField decrypts field data with its aes key and nonce.
func DecryptField(ciphertext, key, nonce []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return []byte(""), nil
	}
	
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
