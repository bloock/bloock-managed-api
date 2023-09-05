package domain

import (
	"errors"
	"strings"
)

type KeyType int

const (
	LOCAL_KEY           KeyType = iota
	MANAGED_KEY         KeyType = iota
	LOCAL_CERTIFICATE   KeyType = iota
	MANAGED_CERTIFICATE KeyType = iota
)

func (t KeyType) String() string {
	switch t {
	case LOCAL_KEY:
		return "local_key"
	case MANAGED_KEY:
		return "managed_key"
	case LOCAL_CERTIFICATE:
		return "local_certificate"
	case MANAGED_CERTIFICATE:
		return "managed_certificate"
	}
	return ""
}

func ParseKeySource(value string) (KeyType, error) {
	switch strings.ToLower(value) {
	case "local_key":
		return LOCAL_KEY, nil
	case "managed_key":
		return MANAGED_KEY, nil
	case "local_certificate":
		return LOCAL_CERTIFICATE, nil
	case "managed_certificate":
		return MANAGED_CERTIFICATE, nil
	}
	return 0, errors.New("invalid key type")
}
