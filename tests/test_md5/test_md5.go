package main

import (
	"crypto/md5"
	"fmt"
)

func GetMd5File(str string)  string{
	str_to_byte := []byte(str)
	byte_ret := md5.Sum(str_to_byte)
	ret := fmt.Sprintf("%x",byte_ret)
	return ret
}



func main() {
	md5_zhiliao := GetMd5File("zhiliao")
	fmt.Println(md5_zhiliao)
}
