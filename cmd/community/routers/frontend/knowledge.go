package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/request"
	service "xhyovo.cn/community/server/service/llm"
)

func InitKnowledgeRouters(ctx *gin.Engine) {

	group := ctx.Group("/community/knowledge")
	group.GET("/like", QueryLikeQuestion)
	group.Use(middleware.OperLogger())
	group.GET("/query", QueryKnowledgies)

}

func QueryLikeQuestion(ctx *gin.Context) {
	question := ctx.Query("question")
	var questionCaches service.QuestionCacheService
	queryQuestion := questionCaches.QueryQuestion(question)
	result.Ok(queryQuestion, "").Json(ctx)
}

func QueryKnowledgies(ctx *gin.Context) {

	question := ctx.Query("question")
	refreshCache := ctx.Query("refreshCache")

	parseBool, err := strconv.ParseBool(refreshCache)
	if err != nil {
		parseBool = false
	}

	var knowS service.KnowledgeBaseService
	userId := middleware.GetUserId(ctx)
	documents, err := knowS.QueryKnowledgies(question, parseBool, userId)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.Ok(documents, "").Json(ctx)

}

func queryKnowledge(ctx *gin.Context) {
	var knowRequest request.KnowledgeRequest
	if err := ctx.ShouldBindJSON(&knowRequest); err != nil {
		msg := utils.GetValidateErr(&knowRequest, err)
		result.Err(msg).Json(ctx)
		return
	}
	var knowS service.KnowledgeBaseService
	knowledge, err := knowS.QueryKnowledge(knowRequest.Question, knowRequest.Document)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.Ok(knowledge, "").Json(ctx)
}
