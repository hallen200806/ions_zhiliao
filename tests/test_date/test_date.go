package main

import (
	"time"
	"fmt"
)

func main() {
	a:= time.Now().Format("2006-01-02")

	b,_ := time.Parse("2006-01-02",a)

	fmt.Println(b)

	c := "2020-08-13"
	d,_ := time.Parse("2006-01-02",c)
	fmt.Println(d)
}
