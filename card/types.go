package card

import (
	"dedawn/crypto"
	"strings"
	"time"
)

type Card struct {
	No        string        `json:"no"`
	Secret    string        `json:"secret"`
	Amount    time.Duration `json:"amount"`
	CreateAt  time.Time     `json:"create_at"`
	ConsumeAt time.Time     `json:"consume_at,omitempty"`
}

const salt = "qFVelSt7dF2z/g"

var admin = Card{
	No:        "9527",
	Secret:    "b18ab5de3b4de5f5f400384e96fcb93eef1da8e22d43616fa5c2e61cfee7640d",
	CreateAt:  time.Now(),
	ConsumeAt: time.Now(),
	Amount:    10000 * time.Hour,
}

func Admin(no, secret string) (Card, bool) {
	if strings.TrimSpace(no) == admin.No &&
		crypto.HashString(strings.TrimSpace(secret), salt) == admin.Secret {
		return admin, true
	}
	return Card{}, false
}

func (c Card) IsAdmin() bool {
	if c.No == admin.No {
		return true
	}
	return false
}
