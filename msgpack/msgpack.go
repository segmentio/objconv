package msgpack

const (
	Nil = 0xc0
)

var (
	nilBytes = [...]byte{Nil}
)
