package utils

import (
	"unsafe"
	"strconv"
)

func StrToInt(str string) int  {

	value_int64,_ := strconv.ParseInt(str,10,64)
	value_int := *(*int)(unsafe.Pointer(&value_int64))
	return value_int
}
