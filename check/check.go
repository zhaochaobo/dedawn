package check

import (
	"dedawn/card"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Check(server, no, secret string) (card.Card, error) {

	c := card.Card{}
	if a, ok := card.Admin(no, secret); ok {
		return a, nil
	}
	if len(no) == 0 {
		return c, fmt.Errorf("the card no. is required")
	}
	if len(secret) == 0 {
		return c, fmt.Errorf("the secret is required")
	}

	resp, err := http.Get(fmt.Sprintf("%s/check?no=%s&secret=%s", server, no, secret))
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	if resp.StatusCode != http.StatusOK {
		e := struct {
			Error string `json:"error"`
		}{}
		if err := json.Unmarshal(data, &e); err != nil {
			return c, err
		}
		return c, fmt.Errorf("check card failed, %s", e.Error)
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}

var ErrCardDepleted = errors.New("card depleted")
var ErrCardNotFound = errors.New("card not found")

func Deduct(server, no string, amount time.Duration) (card.Card, error) {
	c := card.Card{}

	if len(no) == 0 {
		return c, fmt.Errorf("the card no. is required")
	}

	resp, err := http.Get(fmt.Sprintf("%s/deduct?no=%s&amount=%s", server, no, amount))
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	if resp.StatusCode != http.StatusOK {
		e := struct {
			Error string `json:"error"`
		}{}
		if err := json.Unmarshal(data, &e); err != nil {
			return c, err
		}
		if resp.StatusCode == http.StatusGone {
			return c, ErrCardDepleted
		}
		if resp.StatusCode == http.StatusNotFound {
			return c, ErrCardNotFound
		}
		return c, fmt.Errorf("check card failed, %s", e.Error)
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}
