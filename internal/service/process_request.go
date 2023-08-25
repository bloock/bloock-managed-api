package service

import (
	"bloock-managed-api/internal/domain"
	"errors"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"strconv"
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

func parseKeyType(value string) (KeyType, error) {
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

type HostingType int

const (
	IPFS   HostingType = iota
	HOSTED HostingType = iota
	NONE   HostingType = iota
)

func parseHostingType(value string) (HostingType, error) {
	switch strings.ToLower(value) {
	case "ipfs":
		return IPFS, nil
	case "hosted":
		return HOSTED, nil
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
	case NONE:
		return "none"
	}
	return ""
}

type ProcessRequest struct {
	file                  []byte
	isIntegrityEnabled    bool
	isAuthenticityEnabled bool
	keyType               KeyType
	kty                   key.KeyType
	keyID                 uuid.UUID
	hostingType           HostingType
	useEnsResolution      bool
}

func NewProcessRequest(data []byte, integrityEnabled string, authenticityEnabled string, keyType string, kty string, kid string, availabilityType string, ensRes string) (*ProcessRequest, error) {
	isIntegrityEnabled, err := strconv.ParseBool(integrityEnabled)
	if err != nil {
		return nil, err
	}

	isAuthenticityEnabled, err := strconv.ParseBool(authenticityEnabled)
	if err != nil {
		return nil, err
	}

	useEnsResolution, err := strconv.ParseBool(ensRes)
	if err != nil {
		return nil, err
	}

	hostingType, err := parseHostingType(availabilityType)
	if err != nil {
		return nil, err
	}

	keyID, err := uuid.Parse(kid)
	if err != nil {
		return nil, err
	}

	ktyp, err := ValidateKeyType(kty)
	if err != nil {
		return nil, err
	}

	authenticityKeyType, err := parseKeyType(keyType)
	if err != nil {
		return nil, err
	}
	return &ProcessRequest{
		file:                  data,
		isIntegrityEnabled:    isIntegrityEnabled,
		isAuthenticityEnabled: isAuthenticityEnabled,
		keyType:               authenticityKeyType,
		kty:                   ktyp,
		keyID:                 keyID,
		hostingType:           hostingType,
		useEnsResolution:      useEnsResolution,
	}, nil
}

func (s ProcessRequest) KeyID() uuid.UUID {
	return s.keyID
}

func (s ProcessRequest) UseEnsResolution() bool {
	return s.useEnsResolution
}

func (s ProcessRequest) Data() []byte {
	return s.file
}
func (s ProcessRequest) ReplaceDataWith(newData []byte) {
	s.file = newData
}
func (s ProcessRequest) IsAuthenticityEnabled() bool {
	return s.isAuthenticityEnabled
}

func (s ProcessRequest) IsIntegrityEnabled() bool {
	return s.isIntegrityEnabled
}

func (s ProcessRequest) HostingType() HostingType {
	return s.hostingType
}

func (s ProcessRequest) KeyType() KeyType {
	return s.keyType
}

func (s ProcessRequest) Kty() key.KeyType {
	return s.kty
}

func ValidateKeyType(kty string) (key.KeyType, error) {
	switch kty {
	case "EcP256k":
		return key.EcP256k, nil
	case "Rsa2048":
		return key.Rsa2048, nil
	case "Rsa3072":
		return key.Rsa3072, nil
	case "Rsa4096":
		return key.Rsa4096, nil
	default:
		return -1, domain.ErrInvalidKeyType(kty)
	}
}
