package card

import (
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

var admin = Card{
	No:        "9527",
	Secret:    "dedawn@9527",
	CreateAt:  time.Now(),
	ConsumeAt: time.Now(),
	Amount:    10000 * time.Hour,
}

func Admin(no, secret string) (Card, bool) {
	if strings.TrimSpace(no) == admin.No && strings.TrimSpace(secret) == admin.Secret {
		return admin, true
	}
	return Card{}, false
}

func IsAdmin() {

}
