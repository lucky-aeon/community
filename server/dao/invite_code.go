package dao

var InviteCode inviteCode

type inviteCode struct {
}

func (d *inviteCode) SaveCode(code int) error {

	return nil
}

func (i *inviteCode) Exist(code uint16) bool {

	return true
}

func (*inviteCode) Del(code int) error {

	return nil
}

func (i *inviteCode) SetState(code uint16) error {

	return nil
}
