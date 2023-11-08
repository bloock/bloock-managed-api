package domain

type LocalStorageStrategy int64

const (
	LocalStorageStrategyUnknown LocalStorageStrategy = iota
	LocalStorageStrategyHash
	LocalStorageStrategyFilename
)

func LocalStorageStrategyFromString(strategy string) LocalStorageStrategy {
	switch strategy {
	case "HASH":
		return LocalStorageStrategyHash
	case "FILENAME":
		return LocalStorageStrategyFilename
	default:
		return LocalStorageStrategyHash
	}
}

func (s LocalStorageStrategy) String() string {
	switch s {
	case LocalStorageStrategyHash:
		return "HASH"
	case LocalStorageStrategyFilename:
		return "FILENAME"
	}

	return "UNKNOWN"
}
