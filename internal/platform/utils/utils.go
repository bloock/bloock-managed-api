package utils

type Bitset []bool

func (b *Bitset) GetBit(i int) bool {
	if i >= len(*b) {
		return false
	}
	return (*b)[i]
}

func (b *Bitset) SetBit(i int, value bool) {
	if i >= len(*b) {
		b.grow(1 + i)
	}
	(*b)[i] = value
}

func (b *Bitset) grow(size int) {
	b2 := make(Bitset, size)
	copy(b2, *b)
	*b = b2
}

func (b *Bitset) Len() int {
	return len(*b)
}

func (b *Bitset) ToBytes() []byte {

	len := b.Len()
	bytes := len >> 3
	if (len & 0x07) != 0 {
		bytes++
	}

	arr := make([]byte, bytes)

	bitIndex := 0
	byteIndex := 0

	for i := 0; i < b.Len(); i++ {

		if b.GetBit(i) {
			arr[byteIndex] |= (byte(1) << byte(7-bitIndex))
		}

		bitIndex++

		if bitIndex == 8 {
			bitIndex = 0
			byteIndex++
		}
	}

	return arr
}

func BitsetFromBytes(b []byte) Bitset {
	bitset := make(Bitset, len(b)*8)
	for i, b := range b {
		for j := 0; j < 8; j++ {
			bitset.SetBit(i*8+j, b>>uint(7-j)&0x01 == 0x01)
		}
	}
	return bitset
}

func Uint16ToBytes(x uint16) []byte {
	bytes := make([]byte, 2)
	bytes[0] = byte(x >> 8)
	bytes[1] = byte(x)
	return bytes
}
