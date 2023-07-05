package main

import (
	"log"
	"net"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"

	"google.golang.org/grpc"
)

type ChatRoomServer struct {
	pb.UnimplementedChatRoomServiceServer
	msgMemMap MsgMemMap
}

func (s *ChatRoomServer) Chat(stream pb.ChatRoomService_ChatServer) error {
	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	crServer := ChatRoomServer{msgMemMap: make(MsgMemMap)}
	pb.RegisterChatRoomServiceServer(grpcServer, &crServer)

	l, err := net.Listen(common.Transport, common.Localhost)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server at %s\n", common.Localhost)
	log.Fatal(grpcServer.Serve(l))
}
