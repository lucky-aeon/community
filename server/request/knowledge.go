package request

type KnowledgeRequest struct {
	Question string `binding:"required" msg:"问题不能为空"`
	Document string `binding:"required" msg:"内容不能为空"`
}
