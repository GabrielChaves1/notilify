package types

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Email struct {
	Address string
}

func NewEmail(address string) (*Email, error) {
	address = strings.TrimSpace(address)

	if !emailRegex.MatchString(address) {
		return nil, NewInvalidEmailError(address)
	}

	return &Email{Address: address}, nil
}

func (e Email) Value() string {
	return e.Address
}

func NewInvalidEmailError(email string) error {
	return fmt.Errorf("invalid email address: %s", email)
}
