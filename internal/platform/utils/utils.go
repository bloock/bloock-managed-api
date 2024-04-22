package utils

func Uint16ToBytes(x uint16) []byte {
	bytes := make([]byte, 2)
	bytes[0] = byte(x >> 8)
	bytes[1] = byte(x)
	return bytes
}
