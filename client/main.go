package main

import "log"

func main() {
	i := NewKeyInput()
	e := i.ErrChan()
	defer i.Close()
	go i.GetKeys(true)

	err := <-e
	if err != nil {
		log.Fatal(err)
	}
}
