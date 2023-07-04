package main

import "fmt"

func DisplayMessage(player string, msg string) {
	fmt.Printf("[%s]", player)
	fmt.Printf("%s\n", msg)
}

func ReverseLines(lns int) {
	for i := 0; i < lns; i++ {
		fmt.Printf("\r\033[1A")
	}
}
