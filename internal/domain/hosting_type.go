package domain

import (
	"errors"
	"strings"
)

type HostingType int

const (
	IPFS   HostingType = iota
	HOSTED HostingType = iota
	LOCAL  HostingType = iota
	NONE   HostingType = iota
)

func ParseHostingType(value string) (HostingType, error) {
	switch strings.ToLower(value) {
	case "ipfs":
		return IPFS, nil
	case "hosted":
		return HOSTED, nil
	case "local":
		return LOCAL, nil
	case "none":
		return NONE, nil
	}
	return 0, errors.New("unsupported hosting")
}
func (h HostingType) String() string {
	switch h {
	case IPFS:
		return "ipfs"
	case HOSTED:
		return "hosted"
	case LOCAL:
		return "local"
	case NONE:
		return "none"
	}
	return ""
}
