package main

import (
	"io"
	"log"
	"net"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"
	"time"

	"google.golang.org/grpc"
)

type ChatRoomServer struct {
	pb.UnimplementedChatRoomServiceServer
	msgMemMap MsgMemMap
}

func (s *ChatRoomServer) Chat(stream pb.ChatRoomService_ChatServer) error {
	for {
		clientMsg, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		var serverMsg *pb.ChatServerMsg
		chatMsg := clientMsg.GetChatMsg()
		playerName := chatMsg.GetName()
		now := time.Now().UnixMilli()
		if playerName == "" {
			serverMsg = &pb.ChatServerMsg{
				UnixMil:  now,
				ChatMsgs: nil,
				Ok:       false,
				ErrMsg:   &pb.ErrorMsg{Msg: "player name is empty"},
			}
		} else {
			mem, exists := s.msgMemMap[playerName]
			timeStamp := clientMsg.GetUnixMil()
			msg := chatMsg.GetMsg()
			if !exists {
				s.msgMemMap[playerName] = NewMsgMemory(timeStamp, msg)
			} else {
				mem.Set(timeStamp, msg)
			}

			msgs := s.msgMemMap.ToChatMsgs()

			serverMsg = &pb.ChatServerMsg{
				UnixMil:  now,
				ChatMsgs: msgs,
				Ok:       true,
				ErrMsg:   nil,
			}
		}

		err = stream.Send(serverMsg)
		if err != nil {
			return err
		}
	}
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
