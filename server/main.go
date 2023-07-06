package main

import (
	"context"
	"io"
	"log"
	"net"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type ChatRoomServer struct {
	pb.UnimplementedChatRoomServiceServer

	rwLock    sync.RWMutex
	msgMemMap map[string]*MsgMemory
	sortedKey []string
}

func (s *ChatRoomServer) DeleteMemory(key string) {
	s.rwLock.Lock()
	delete(s.msgMemMap, key)
	len := len(s.sortedKey)
	for i := 0; i < len; i++ {
		if strings.Compare(s.sortedKey[i], key) == 0 {
			s.sortedKey = append(s.sortedKey[:i], s.sortedKey[i+1:]...)
			break
		}
	}
	s.rwLock.Unlock()
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
	var playerName string
	var e error

	for {
		msgs := s.MakeSortedChatMsg()
		now := time.Now().UnixMilli()
		serverMsg := &pb.ChatServerMsg{
			UnixMil:  now,
			ChatMsgs: msgs,
			Ok:       true,
			ErrMsg:   nil,
		}

		err := stream.Send(serverMsg)
		if err != nil {
			e = err
			break
		}

		clientMsg, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			e = err
			break
		}

		chatMsg := clientMsg.GetChatMsg()
		playerName = chatMsg.GetName()
		if playerName == "" {
			break
		} else {
			mem, exists := s.msgMemMap[playerName]
			if !exists {
				break
			} else {
				msg := chatMsg.GetMsg()
				if len(msg) != common.InputBufferSize {
					break
				} else {
					mem.Set(clientMsg.GetUnixMil(), msg)
				}
			}
		}
	}

	if playerName != "" {
		s.DeleteMemory(playerName)
	}
	return e
}

func main() {
	grpcServer := grpc.NewServer()
	crServer := ChatRoomServer{
		msgMemMap: make(map[string]*MsgMemory),
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
