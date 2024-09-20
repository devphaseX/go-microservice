package main

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	symmetricKey []byte
	paster       *paseto.V2
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d size", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		symmetricKey: []byte(symmetricKey),
		paster:       paseto.NewV2(),
	}, nil
}

func (p *PasetoMaker) CreateToken(userId string, email string, expirationDuration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, email, expirationDuration)

	if err != nil {
		return "", nil, err
	}

	var tokenStr string
	tokenStr, err = p.paster.Encrypt(p.symmetricKey, payload, nil)
	return tokenStr, payload, err
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var (
		payload Payload
		err     error
	)

	if err = p.paster.Decrypt(token, p.symmetricKey, &payload, nil); err != nil {
		return nil, ErrInvalidToken
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return &payload, nil
}
