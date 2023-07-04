package main

import (
	"log"
	"time"
)

func mainLoop() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for range ticker.C {

	}
}

func main() {
	i := NewKeyInput()
	e := i.ErrChan()
	defer i.Close()
	go i.GetKeys()

	go mainLoop()

	err := <-e
	if err != nil {
		log.Fatal(err)
	}
}
