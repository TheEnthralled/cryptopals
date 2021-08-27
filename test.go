package main

import (
	"fmt"
	"cryptopals/morestrings"
	"cryptopals/utils"
)

func main(){
	fmt.Println(morestrings.ReverseRunes("Hello World"))
	i , _ := utils.RandInt(-50, 10)
	fmt.Println(i)
}
