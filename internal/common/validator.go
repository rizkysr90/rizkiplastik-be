package common

import (
	"errors"
	"strconv"
)

func ValidateMaxLengthStr(str string, maxLength int) error {
	if len(str) > maxLength {
		return errors.New("string length must be less than " + strconv.Itoa(maxLength))
	}
	return nil
}

func ValidateMinLengthStr(str string, minLength int) error {
	if len(str) < minLength {
		return errors.New("string length must be greater than " + strconv.Itoa(minLength))
	}
	return nil
}

func ValidateLenghtEqual(str string, length int) error {
	if len(str) != length {
		return errors.New("string length must be equal to " + strconv.Itoa(length))
	}
	return nil
}
