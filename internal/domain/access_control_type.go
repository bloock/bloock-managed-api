package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyAccessCode = errors.New("empty access code provided")
)

type AccessControlType int

const (
	TotpAccessControl   AccessControlType = iota
	SecretAccessControl AccessControlType = iota
)

func (t AccessControlType) String() string {
	switch t {
	case TotpAccessControl:
		return "totp"
	case SecretAccessControl:
		return "secret"
	}
	return ""
}

func ParseAccessControlType(value string) (AccessControlType, error) {
	switch strings.ToLower(value) {
	case "totp":
		return TotpAccessControl, nil
	case "secret":
		return SecretAccessControl, nil
	}
	return 0, errors.New("invalid access control type")
}
