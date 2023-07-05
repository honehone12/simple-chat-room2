package main

import (
	"context"
	"io"
	"log"
	"net"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"
	"time"

	"google.golang.org/grpc"
)

type MsgMemMap map[string]*MsgMemory

type ChatRoomServer struct {
	pb.UnimplementedChatRoomServiceServer
	msgMemMap MsgMemMap
	sortedKey []string
}

func (s *ChatRoomServer) MakeSortedChatMsg() []*pb.ChatMsg {
	len := len(s.sortedKey)
	msgs := make([]*pb.ChatMsg, 0, len)
	for i := 0; i < len; i++ {
		key := s.sortedKey[i]
		msgs = append(msgs, &pb.ChatMsg{
			Name: key,
			Msg:  s.msgMemMap[key].msg,
		})
	}
	return msgs
}

func (s *ChatRoomServer) Join(
	ctx context.Context, req *pb.JoinRequest,
) (*pb.JoinResponse, error) {
	var res *pb.JoinResponse
	playerName := req.GetName()
	if playerName == "" {
		res = &pb.JoinResponse{
			Ok:     false,
			ErrMsg: &pb.ErrorMsg{Msg: "player name is empty"},
		}
	} else {
		_, exists := s.msgMemMap[playerName]
		if exists {
			res = &pb.JoinResponse{
				Ok:     false,
				ErrMsg: &pb.ErrorMsg{Msg: "player name is already used"},
			}
		} else {
			s.msgMemMap[playerName] = NewMsgMemory(time.Now().UnixMilli())
			s.sortedKey = append(s.sortedKey, playerName)
			res = &pb.JoinResponse{
				Ok:     true,
				ErrMsg: nil,
			}
		}
	}
	return res, nil
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
		now := time.Now().UnixMilli()
		chatMsg := clientMsg.GetChatMsg()
		playerName := chatMsg.GetName()
		if playerName == "" {
			serverMsg = &pb.ChatServerMsg{
				UnixMil:  now,
				ChatMsgs: nil,
				Ok:       false,
				ErrMsg:   &pb.ErrorMsg{Msg: "player name is empty"},
			}
		} else {
			mem, exists := s.msgMemMap[playerName]
			if !exists {
				serverMsg = &pb.ChatServerMsg{
					UnixMil:  now,
					ChatMsgs: nil,
					Ok:       false,
					ErrMsg:   &pb.ErrorMsg{Msg: "no such player"},
				}
			} else {
				mem.Set(clientMsg.GetUnixMil(), chatMsg.GetMsg())
				msgs := s.MakeSortedChatMsg()

				serverMsg = &pb.ChatServerMsg{
					UnixMil:  now,
					ChatMsgs: msgs,
					Ok:       true,
					ErrMsg:   nil,
				}
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
	crServer := ChatRoomServer{
		msgMemMap: make(MsgMemMap),
		sortedKey: make([]string, 0),
	}
	pb.RegisterChatRoomServiceServer(grpcServer, &crServer)

	l, err := net.Listen(common.Transport, common.Localhost)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server at %s\n", common.Localhost)
	log.Fatal(grpcServer.Serve(l))
}
