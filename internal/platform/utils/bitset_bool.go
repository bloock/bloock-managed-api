package utils

type BitsetBool []bool

func (b *BitsetBool) GetBit(i int) bool {
	if i >= len(*b) {
		return false
	}
	return (*b)[i]
}

func (b *BitsetBool) SetBit(i int, value bool) {
	if i >= len(*b) {
		b.grow(1 + i)
	}
	(*b)[i] = value
}

func (b *BitsetBool) grow(size int) {
	b2 := make(BitsetBool, size)
	copy(b2, *b)
	*b = b2
}

func (b *BitsetBool) Len() int {
	return len(*b)
}

func (b *BitsetBool) ToBytes() []byte {
	length := b.Len()
	bytes := length >> 3
	if (length & 0x07) != 0 {
		bytes++
	}

	arr := make([]byte, bytes)

	bitIndex := 0
	byteIndex := 0

	for i := 0; i < b.Len(); i++ {
		if b.GetBit(i) {
			arr[byteIndex] |= byte(1) << byte(7-bitIndex)
		}
		bitIndex++

		if bitIndex == 8 {
			bitIndex = 0
			byteIndex++
		}
	}
	return arr
}

func BitsetBoolFromBytes(b []byte) BitsetBool {
	bitset := make(BitsetBool, len(b)*8)
	for i, b := range b {
		for j := 0; j < 8; j++ {
			bitset.SetBit(i*8+j, b>>uint(7-j)&0x01 == 0x01)
		}
	}
	return bitset
}
