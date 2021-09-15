package client

import (
	"go_im/im"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/message"
)

type Manager struct {
	appId int64
	m     client.IClientManager
}

func NewManager() *Manager {
	ret := &Manager{}
	ret.m = im.NewClientManager()
	return ret
}

func (m *Manager) ClientConnected(conn conn.Connection) int64 {
	connId := m.m.ClientConnected(conn)
	return connId
}

func (m *Manager) ClientSignIn(oldUid int64, uid int64, device int64) {
	m.m.ClientSignIn(oldUid, uid, device)
}

func (m *Manager) UserLogout(uid int64) {
	m.m.UserLogout(uid)
}

func (m *Manager) HandleMessage(from int64, message *message.Message) error {
	return m.m.HandleMessage(from, message)
}

func (m *Manager) EnqueueMessage(uid int64, message *message.Message) {
	m.m.EnqueueMessage(uid, message)
}

func (m *Manager) AllClient() []int64 {
	return m.m.AllClient()
}