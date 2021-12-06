package apidep

import (
	"go_im/im/client"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

// ClientManager 客户端连接相关接口
var ClientManager ClientManagerInterface = client.Manager

// GroupManager 群相关接口
var GroupManager GroupManagerInterface = &groupInterface{}

func SendMessage(uid int64, device int64, m *message.Message) {
	ClientManager.EnqueueMessage(uid, device, m)
	logger.D("uid=%d, device=%d, msg=%v", uid, device, m)
}

func SendMessageIfOnline(uid int64, device int64, m *message.Message) {
	if ClientManager.IsDeviceOnline(uid, device) {
		SendMessage(uid, device, m)
	}
}

type ClientManagerInterface interface {
	ClientSignIn(oldUid int64, uid int64, device int64)
	ClientLogout(uid int64, device int64)
	EnqueueMessage(uid int64, device int64, message *message.Message)
	IsDeviceOnline(uid, device int64) bool
	IsOnline(uid int64) bool
	AllClient() []int64
}

type GroupManagerInterface interface {
	MemberOnline(gid int64, uid int64) error
	MemberOffline(gid int64, uid int64) error
	PutMember(gid int64, mb map[int64]int32) error
	RemoveMember(gid int64, uid ...int64) error
	CreateGroup(gid int64) error
	DissolveGroup(gid int64) error
	MuteGroup(gid int64, mute bool) error
	DispatchNotifyMessage(gid int64, message *message.Message) error
}

type groupInterface struct{}

func (g *groupInterface) MemberOnline(gid int64, uid int64) error {
	u := []group.MemberUpdate{
		{
			Uid:   uid,
			Flag:  group.FlagMemberOnline,
			Extra: nil,
		},
	}
	return group.Manager.UpdateMember(gid, u)
}

func (g *groupInterface) MemberOffline(gid int64, uid int64) error {
	u := []group.MemberUpdate{
		{
			Uid:   uid,
			Flag:  group.FlagMemberOffline,
			Extra: nil,
		},
	}
	return group.Manager.UpdateMember(gid, u)
}

func (g *groupInterface) PutMember(gid int64, mb map[int64]int32) error {

	var u []group.MemberUpdate
	for uid := range mb {
		u = append(u, group.MemberUpdate{
			Uid:  uid,
			Flag: group.FlagMemberAdd,
		})
	}
	return group.Manager.UpdateMember(gid, u)
}

func (g *groupInterface) RemoveMember(gid int64, uid ...int64) error {
	var u []group.MemberUpdate
	for _, id := range uid {
		u = append(u, group.MemberUpdate{
			Uid:  id,
			Flag: group.FlagMemberDel,
		})
	}
	return group.Manager.UpdateMember(gid, u)
}

func (g *groupInterface) CreateGroup(gid int64) error {
	return group.Manager.UpdateGroup(gid, group.Update{Flag: group.FlagGroupCreate})
}

func (g *groupInterface) DissolveGroup(gid int64) error {
	return group.Manager.UpdateGroup(gid, group.Update{Flag: group.FlagGroupDissolve})
}

func (g *groupInterface) MuteGroup(gid int64, mute bool) error {
	var f = group.FlagGroupMute
	if !mute {
		f = group.FlagGroupCancelMute
	}
	return group.Manager.UpdateGroup(gid, group.Update{Flag: int64(f)})
}

func (g *groupInterface) DispatchNotifyMessage(gid int64, message *message.Message) error {
	return group.Manager.DispatchNotifyMessage(gid, message)
}
