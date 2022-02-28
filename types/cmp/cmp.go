package cmp

// package cmp is a simple helper contain min/max with generic
type Ordered interface {
	int | int8
}

func Max[T Ordered](a, b T) T {
	if a >= b {
		return a
	}
	return b
}

func Min[T Ordered](a, b T) T {
	if a >= b {
		return b
	}
	return a
}

func MaxInSeq[T Ordered](s []T) (res T) {
	for i := range s {
		if res < s[i] {
			res = s[i]
		}
	}
	return res
}
