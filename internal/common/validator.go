package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

func ValidateUUIDFormat(str string) error {
	if _, err := uuid.Parse(str); err != nil {
		return errors.New("invalid uuid format")
	}
	return nil
}
func ValidateEquals(str string, allowedWords []string) error {
	isTrue := false
	for _, word := range allowedWords {
		if strings.Contains(str, word) {
			isTrue = true
			break
		}
	}
	if !isTrue {
		return errors.New("string must contain only allowed words")
	}
	return nil
}

func ValidateOnlyAllowedUppercaseLetter(str, fieldName string) error {
	for _, char := range str {
		if char < 'A' || char > 'Z' {
			return errors.New(fieldName + " must contain only allowed uppercase letter (A-Z)")
		}
	}
	return nil
}
