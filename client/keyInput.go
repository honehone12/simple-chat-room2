package main

import (
	pb "simple-chat-room2/pb"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	keyBufferSize = 10
	nonAlphaNum   = 0x0
	space         = 0x20
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

func (i KeyInput) Input(name string, stream pb.ChatRoomService_ChatClient) {
	keyEvents, err := keyboard.GetKeys(keyBufferSize)
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
			i.buffer.Add(space)
		} else if e.Rune != nonAlphaNum {
			// still not sure what will be lost with cast
			i.buffer.Add(byte(e.Rune))
		}

		s, err := i.buffer.String()
		if err != nil {
			i.errCh <- err
			break
		}

		err = stream.Send(&pb.ChatClientMsg{
			UnixMil: time.Now().UnixMilli(),
			ChatMsg: &pb.ChatMsg{
				Name: name,
				Msg:  s,
			},
		})
		if err != nil {
			i.errCh <- err
			break
		}
	}
}
