package common

type Float interface {
	~float32 | ~float64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Integer interface {
	Signed | Unsigned
}
type Ordered interface {
	Integer | Float | ~string
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
