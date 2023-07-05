package main

import "sync"

type MsgMemory struct {
	joinedAt  int64
	name      string
	timeStamp int64
	msg       string

	rwLock sync.RWMutex
}

func NewMsgMemory(timeStamp int64) *MsgMemory {
	return &MsgMemory{
		joinedAt:  timeStamp,
		timeStamp: timeStamp,
		msg:       "",
	}
}

func (m *MsgMemory) Get() (int64, string, string) {
	return m.timeStamp, m.name, m.msg
}

func (m *MsgMemory) Set(timeStamp int64, msg string) {
	if timeStamp > m.timeStamp {
		m.rwLock.Lock()
		m.timeStamp = timeStamp
		m.msg = msg
		m.rwLock.Unlock()
	}
}
