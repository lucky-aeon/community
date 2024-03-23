package services

import (
	"xhyovo.cn/community/server/model"
)

type UserTag struct {
}

func (*UserTag) Page(page, limit int) (userTags []model.UserTags, count int64) {
	db := model.UserTag()
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Find(&userTags)
	return
}

func (*UserTag) Save(userTag model.UserTags) {
	if userTag.ID != 0 {
		model.UserTag().Updates(&userTag)
	} else {
		model.UserTag().Create(&userTag)
	}
}

func (*UserTag) DeleteById(id int) {
	model.UserTag().Where("id = ?", id).Delete(&model.UserTags{})
}

func (*UserTag) AssignUserLabel(userId int, tags []int) (okIds []int) {

	// 删除用户的标签
	model.UserTagRelation().Where("user_id = ?", userId).Delete(&model.UserTagRelations{})
	if len(tags) != 0 {
		var userTags []model.UserTagRelations
		for i := range tags {
			userTags = append(userTags, model.UserTagRelations{UserId: userId, UserTagId: tags[i]})
		}
		// 分配标签
		model.UserTagRelation().Create(userTags)
		model.UserTag().Where("id in ?", tags).Select("id").Find(&okIds)
	}
	return
}

func (*UserTag) GetTagsByUserId(userId int) (tagNames []model.UserTags) {
	var tagsIds []int
	model.UserTagRelation().Where("user_id = ?", userId).Select("user_tag_id").Find(&tagsIds)
	if len(tagsIds) == 0 {
		return
	}
	model.UserTag().Where("id in ?", tagsIds).Find(&tagNames)
	return tagNames
}
