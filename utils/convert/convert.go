package convert

import "strconv"

func Atoi(s string) int64 {
	val, _ := strconv.Atoi(s)
	return int64(val)
}
