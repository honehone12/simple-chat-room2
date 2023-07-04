package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type KeyInput struct {
	buffer *Buffer

	errCh chan error
}

func NewKeyInput() KeyInput {
	return KeyInput{
		buffer: NewBuffer(),
		errCh:  make(chan error),
	}
}

func (i KeyInput) ErrChan() <-chan error {
	return i.errCh
}

func (i KeyInput) Close() error {
	return keyboard.Close()
}

func (i KeyInput) GetKeys() {
	keyEvents, err := keyboard.GetKeys(10)
	if err != nil {
		i.errCh <- err
		return
	}

	for e := range keyEvents {
		if e.Err != nil {
			i.errCh <- e.Err
			break
		}

		if e.Key == keyboard.KeyEsc {
			i.errCh <- nil
			break
		} else if e.Key == keyboard.KeyBackspace || e.Key == keyboard.KeyBackspace2 {
			i.buffer.Back()
		} else if e.Key == keyboard.KeySpace {
			i.buffer.Add(Space)
		} else if e.Rune != NonAlphaNum {
			// still not sure what will be lost with cast
			i.buffer.Add(byte(e.Rune))
		}

		// then send buffer to server
		fmt.Printf("%v\n", i.buffer)
	}
}
