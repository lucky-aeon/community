package services

import (
	"strconv"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/constant"
)

type BlacklistService struct {
}

// token 黑名单，如果在 ttl 期间超过 x 次，则将该用户永久封禁小黑屋
func (*BlacklistService) Add(userId int, token string) {
	cache := cache.GetInstance()
	cache.Set(constant.BLACK_LIST+strconv.Itoa(userId), 1, constant.Token_TTl)
	key := constant.BLACK_LIST_COUNT + strconv.Itoa(userId)
	v, b := cache.Get(key)
	if !b {
		cache.Set(key, 1, constant.Token_TTl)
	} else {
		// 如果缓存中有该键，将其转换为int并加1
		intValue, _ := v.(int)
		intValue = intValue + 1
		cache.Set(key, intValue, constant.Token_TTl)
		// 超过 5 则永久关闭小黑屋 todo 先写死
		if intValue > 5 {
			var userService UserService
			userService.BanByUserId(userId)

		}
	}
}

func (*BlacklistService) AddBlackByToken(token string) {
	cache := cache.GetInstance()
	cache.Set(constant.BLACK_LIST+token, 1, constant.Token_TTl)
}

func (*BlacklistService) ExistToken(token string) bool {
	_, b := cache.GetInstance().Get(constant.BLACK_LIST + token)
	return b
}
