package vo

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

type UserID string

const (
	minUserIDLength = 1
	maxUserIDLength = 10
)

func NewUserID(userID string) (UserID, error) {
	if n := utf8.RuneCountInString(userID); n < minUserIDLength || n > maxUserIDLength {
		return "", xerrors.Errorf("user id must be %d or more and %d or less: %s", minUserIDLength, maxUserIDLength, userID)
	}

	if strings.Contains(userID, " ") || strings.Contains(userID, "　") {
		return "", xerrors.Errorf("user id cannot contain spaces: %s", userID)
	}

	return UserID(userID), nil
}

func (i UserID) Value() string {
	return string(i)
}
