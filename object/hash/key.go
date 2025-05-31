package hash

type Key uint64

func (h Key) Bytes() []byte {
	bytes := [4]byte{}
	bytes[0] = byte(h >> 24)
	bytes[1] = byte(h >> 16)
	bytes[2] = byte(h >> 8)
	bytes[3] = byte(h)
	return bytes[:]
}
