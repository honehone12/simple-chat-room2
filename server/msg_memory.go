package main

import "sync"

type MsgMemory struct {
	timeStamp int64
	msg       string

	rwLock sync.RWMutex
}

func NewMsgMemory(timeStamp int64, msg string) *MsgMemory {
	return &MsgMemory{
		timeStamp: timeStamp,
		msg:       msg,
	}
}

func (m *MsgMemory) Get() (int64, string) {
	return m.timeStamp, m.msg
}

func (m *MsgMemory) Set(timeStamp int64, msg string) {
	if timeStamp > m.timeStamp {
		m.rwLock.Lock()
		m.timeStamp = timeStamp
		m.msg = msg
		m.rwLock.Unlock()
	}
}
