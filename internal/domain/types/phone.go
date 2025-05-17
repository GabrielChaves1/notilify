package types

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	phoneNumberRegex = regexp.MustCompile(`^(\+?55)?\d{10,11}$`)
)

type PhoneNumber struct {
	Number string
}

func NewPhoneNumber(number string) (*PhoneNumber, error) {
	number = strings.TrimSpace(number)

	if !phoneNumberRegex.MatchString(number) {
		return nil, NewInvalidPhoneNumberError(number)
	}

	return &PhoneNumber{Number: number}, nil
}

func (p PhoneNumber) Value() string {
	return p.Number
}

func NewInvalidPhoneNumberError(phoneNumber string) error {
	return fmt.Errorf("invalid phone number: %s", phoneNumber)
}
