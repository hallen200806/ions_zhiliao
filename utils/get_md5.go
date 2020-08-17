package utils

import (
	"fmt"
	"crypto/md5"
)

func GetMd5Str(str string)  string{
	str_to_byte := []byte(str)
	byte_ret := md5.Sum(str_to_byte)
	ret := fmt.Sprintf("%x",byte_ret)
	return ret
}
