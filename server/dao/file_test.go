package dao_test

import (
	"testing"
	"xhyovo.cn/community/server/config"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

var file dao.File

func TestMain(t *testing.M) {
	config.InitConfig()
	t.Run()
}

func TestSave(t *testing.T) {
	file.Save(&model.File{FileKey: "sadasd", Size: 500})
}
