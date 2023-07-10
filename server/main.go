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

	rwLock     sync.RWMutex
	msgMemMap  map[string]*MsgMemory
	sortedKeys []string
}

func (s *ChatRoomServer) DeleteMemory(key string) {
	s.rwLock.Lock()
	delete(s.msgMemMap, key)
	len := len(s.sortedKeys)
	for i := 0; i < len; i++ {
		if strings.Compare(s.sortedKeys[i], key) == 0 {
			s.sortedKeys = append(s.sortedKeys[:i], s.sortedKeys[i+1:]...)
			break
		}
	}
	s.rwLock.Unlock()
}

func (s *ChatRoomServer) MakeSortedChatMsg() []*pb.ChatMsg {
	len := len(s.sortedKeys)
	msgs := make([]*pb.ChatMsg, 0, len)
	for i := 0; i < len; i++ {
		key := s.sortedKeys[i]
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
			s.rwLock.Lock()
			s.msgMemMap[playerName] = NewMsgMemory(time.Now().UnixMilli())
			s.sortedKeys = append(s.sortedKeys, playerName)
			res = &pb.JoinResponse{
				Ok:     true,
				ErrMsg: nil,
			}
			s.rwLock.Unlock()
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
			log.Println("received empty player name")
			break
		} else {
			mem, exists := s.msgMemMap[playerName]
			if !exists {
				log.Printf("received unknown player name: %s", playerName)
				playerName = ""
				break
			} else {
				msg := chatMsg.GetMsg()
				if len(msg) > common.InputBufferSize {
					log.Println("received oversize message")
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
		msgMemMap:  make(map[string]*MsgMemory),
		sortedKeys: make([]string, 0),
	}
	pb.RegisterChatRoomServiceServer(grpcServer, &crServer)

	l, err := net.Listen(common.Transport, common.Localhost)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server at %s\n", common.Localhost)
	log.Fatal(grpcServer.Serve(l))
}
