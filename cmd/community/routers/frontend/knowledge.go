package frontend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/request"
	service "xhyovo.cn/community/server/service/llm"
)

func InitKnowledgeRouters(ctx *gin.Engine) {

	group := ctx.Group("/community/knowledge")
	group.Use(middleware.OperLogger())
	group.GET("/query", QueryKnowledgies)
	group.GET("/	", queryKnowledge)

}

func QueryKnowledgies(ctx *gin.Context) {

	question := ctx.Query("question")
	var knowS service.KnowledgeBaseService
	documents, err := knowS.QueryKnowledgies(question)
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
