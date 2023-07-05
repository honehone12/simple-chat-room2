package main

import (
	"context"
	"flag"
	"log"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	playerName := flag.String("name", "", "player's name")
	flag.Parse()

	if *playerName == "" {
		log.Fatalln("name is needed")
	}

	input := NewKeyInput()
	display := NewDisplay()

	conn, err := grpc.Dial(common.Localhost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	crClient := pb.NewChatRoomServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := crClient.Chat(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go input.Input(*playerName, stream)
	go display.Display(stream)

	log.Fatal(catch(
		input.ErrChan(),
		display.ErrChan(),
	))
}

func catch(inputErr <-chan error, displayErr <-chan error) error {
	var err error
	select {
	case err = <-inputErr:
	case err = <-displayErr:
	}
	return err
}
