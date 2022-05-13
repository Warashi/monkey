package converter

import "golang.org/x/exp/constraints"

func Ints[T, U constraints.Integer](from []U) []T {
	t := make([]T, len(from))
	for i := range from {
		t[i] = T(from[i])
	}
	return t
}
