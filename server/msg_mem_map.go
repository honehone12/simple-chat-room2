package main

import (
	pb "simple-chat-room2/pb"
)

type MsgMemMap map[string]*MsgMemory

func (m MsgMemMap) ToChatMsgs() []*pb.ChatMsg {
	msgs := make([]*pb.ChatMsg, 0, len(m))
	for _, v := range m {
		msgs = append(msgs, &pb.ChatMsg{
			Name: v.name,
			Msg:  v.msg,
		})
	}
	return msgs
}
