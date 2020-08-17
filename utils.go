package main

import (
	"fmt"
)

func panicCheck(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}
