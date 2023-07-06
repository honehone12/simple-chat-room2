package main

import (
	"simple-chat-room2/common"
	"sync"
)

type MsgMemory struct {
	timeStamp int64
	msg       string

	rwLock sync.RWMutex
}

func NewMsgMemory(timeStamp int64) *MsgMemory {
	return &MsgMemory{
		timeStamp: timeStamp,
		msg:       common.Space64,
	}
}

func (m *MsgMemory) Set(timeStamp int64, msg string) {
	if timeStamp > m.timeStamp {
		m.rwLock.Lock()
		m.timeStamp = timeStamp
		m.msg = msg
		m.rwLock.Unlock()
	}
}
