package Utils

import (
	"golang.org/x/exp/constraints"
	"strconv"
)

func DecToBin[T constraints.Integer](number T) string {
	return strconv.FormatInt(int64(number), 2)
}
