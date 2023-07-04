package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

type KeyInput struct {
	buffer []byte

	errCh chan error
}

func NewKeyInput() KeyInput {
	return KeyInput{
		buffer: make([]byte, 0, 1024),
		errCh:  make(chan error),
	}
}

func (i *KeyInput) ErrChan() <-chan error {
	return i.errCh
}

func (i *KeyInput) Buffer() []byte {
	return i.buffer
}

func (i *KeyInput) Close() error {
	return keyboard.Close()
}

func (i *KeyInput) GetKeys(print bool) {
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
			if print {
				fmt.Print("\n")
			}
			i.errCh <- nil
			break
		}

		var b byte
		if e.Rune != 0x00 {
			// still not sure what will be lost with cast
			b = byte(e.Rune)
		} else if e.Key == keyboard.KeySpace {
			b = 0x20
		}
		i.buffer = append(i.buffer, b)

		if print {
			fmt.Printf("%c", b)
		}
	}
}
