package dao

import (
	"time"
)

var InviteCode inviteCode

type inviteCode struct {
}

func (*inviteCode) SaveCode(code int) error {
	stem, err := DB.Prepare("insert into invite_code (code,state,created_at,updated_at) values(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stem.Close()

	_, err = stem.Exec(code, false, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

// query single code
func (*inviteCode) QueryCode(code int) int {

	row := DB.QueryRow("SELECT code FROM user WHERE code = ?", code)

	var resCode int
	row.Scan(&resCode)

	return resCode
}

func (i *inviteCode) Exist(code int) bool {
	return i.QueryCode(code) == 0
}

func (*inviteCode) Del(code int) error {
	stem, err := DB.Prepare("delete from invite_code where code = ?")
	if err != nil {
		return err
	}
	defer stem.Close()

	if _, err = stem.Exec(code); err != nil {
		return err
	}
	return nil
}
