package services

import model "xhyovo.cn/community/server/model/knowledge"

type QuestionCacheService struct{}

func (q *QuestionCacheService) GetAnswer(question string) []model.AnswerCaches {

	var ids []int
	model.QuestionCache().Where("question = ?", question).Select("id").Find(&ids)

	var answers []model.AnswerCaches

	if len(ids) == 0 {
		return answers
	}
	model.AnswerCache().Where("question_cache_id in ?", ids).Find(&answers)

	return answers

}

func (q *QuestionCacheService) QueryQuestion(question string) []string {
	var questions []string
	model.QuestionCache().Where("question like ?", "%"+question+"%").Select("question").Find(&questions)
	return questions
}

func (q *QuestionCacheService) Put(question string, caches []model.AnswerCaches) {

	// 先删除缓存
	q.DeleteCache(question)

	questionModel := &model.QuestionCaches{
		Question: question,
	}

	model.QuestionCache().Create(&questionModel)

	for i := range caches {
		caches[i].QuestionCacheId = questionModel.ID
	}

	model.AnswerCache().Create(&caches)
}

func (q *QuestionCacheService) DeleteCache(question string) {

	var id int
	model.QuestionCache().Where("question = ?", question).Select("id").Find(&id)

	model.QuestionCache().Delete("id = ?", id).Delete(&model.QuestionCaches{})
	model.AnswerCache().Where("question_cache_id = ?", id).Delete(&model.AnswerCaches{})
}
