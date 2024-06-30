package token

import (
	"time"
)

type Maker interface {
	CreateToken(username string , duration time.Duration) (string , error) // createstoken for specific user
	VerifyToken(token string) (*Payload , error) // verifies token and returns payload
}

