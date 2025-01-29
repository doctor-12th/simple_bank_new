package token

import (
	"errors"
	// "fmt"
	"time"

	"github.com/google/uuid"
)

type Payload struct{
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`

	// jwt.StandardClaims
}
var(
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")	
)



func NewPayload(username string,duration time.Duration) (*Payload,error){
	payload := &Payload{
		ID: uuid.New(),
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt){
		// fmt.Println("paylaod is expired")
		return ErrExpiredToken
	}
	return nil
}

