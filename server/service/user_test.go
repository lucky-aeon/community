package services

import (
	"testing"
)

func TestEncryptPwd(t *testing.T) {
	password := "123123"
	ret, err := GetPwd(password)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("EncryptPwd: %s -> %s", password, ret)
}
