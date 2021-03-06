package cache

type ByteView struct {
	data []byte
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.data)
}

func (b ByteView) String() string {
	return string(b.data)
}

func cloneBytes(b []byte) []byte {
	r := make([]byte, len(b))
	copy(r, b)
	return r
}
