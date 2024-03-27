package services

import (
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/server/model"
)

type QAAdoption struct {
}

// 查询采纳评论
func (*QAAdoption) SetAdoptionComment(comments []*model.Comments) {

	var ids = make([]int, 0, len(comments))
	for i := range comments {
		ids = append(ids, comments[i].ID)
	}
	if len(ids) == 0 {
		return
	}
	var qaCommentIds []int
	model.QaAdoption().Where("comment_id in ?", ids).Select("comment_id").Find(&qaCommentIds)

	set := mapset.NewSet[int](qaCommentIds...)

	for i := range comments {
		if set.Contains(comments[i].ID) {
			comments[i].AdoptionState = true
		}
	}

}

// 采纳/取消采纳 评论
func (*QAAdoption) Adopt(articleId, commentId int) bool {

	if err := model.QaAdoption().Create(&model.QaAdoptions{ArticleId: articleId, CommentId: commentId}).Error; err != nil {
		model.QaAdoption().Where("comment_id = ? ", commentId).Delete(&model.QaAdoptions{})
		return false
	}

	return true
}

// QA 采纳状态 -> 已解决 or 未解决 变更,  > 0 解决  = 0 未解决
func (*QAAdoption) QAAdoptState(articleId int) bool {
	var count int64
	model.QaAdoption().Where("article_id = ?", articleId).Count(&count)
	return count > 0
}
