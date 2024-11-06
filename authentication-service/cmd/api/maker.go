package main

import "time"

type TokenMaker interface {
	CreateToken(userId string, email string, expirationDuration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
