package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pgvector/pgvector-go"
	"strconv"
	"strings"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	dao "xhyovo.cn/community/server/dao/llm"
	model "xhyovo.cn/community/server/model/knowledge"
)

type KnowledgeBaseService struct {
}

/*
*
添加知识
content:原文
typee:类型
linke:链接
*/
func (k *KnowledgeBaseService) AddKnowledge(content, link, remake string, typee constant.ContentType, businessId int) error {

	// 先删除已存在的知识
	err := k.DeleteKnowledge(businessId, typee)
	if err != nil {
		return err
	}

	// 分片知识
	var llm LLMService
	texts, err := llm.Chat(constant.KnowledgeSplitPrompt, content)
	if err != nil {
		return err
	}
	array, err := textToArray(texts)
	if err != nil {
		return err
	}

	if len(array) == 0 {
		return nil
	}

	document := &model.Documents{
		Content:    content,
		Type:       typee,
		Link:       link,
		Remark:     remake,
		BusinessId: businessId,
	}

	var documentDao dao.DocumentDao
	documentId, err := documentDao.Create(document)
	if err != nil {
		return err
	}

	var embeddingService EmbeddingService
	// 向量化,一行一行向量化，避免因为某行原因导致整个结果返回不了
	for j := range array {
		embeddings, err := embeddingService.GetTextEmbeddings([]string{array[j]})
		if err != nil {
			continue
		}

		var vectors []model.Vectors
		items := embeddings.Output.Embeddings
		for i := range items {
			item := items[i]
			vectors = append(vectors, model.Vectors{
				Content:    array[j],
				DocumentId: documentId,
				Embedding:  pgvector.NewVector(item.Embedding),
				Type:       1,
			})

		}
		var vectorDao dao.VectorDao
		err = vectorDao.Inserts(vectors)
		if err != nil {
			return err
		}
	}

	// 存储
	return nil

}

// 删除知识，link 是唯一性的
func (k *KnowledgeBaseService) DeleteKnowledge(businessId int, typee constant.ContentType) error {
	// 查出文档 id
	var documentdao dao.DocumentDao
	if !documentdao.Exist(businessId, typee) {
		return nil
	}

	document, err := documentdao.GetByLink(businessId, typee)
	if err != nil {
		return nil
	}

	err = documentdao.Delete(businessId, typee)
	if err != nil {
		return err
	}
	var vectorDao dao.VectorDao
	err = vectorDao.DeleteByDocumentId(document.ID)
	if err != nil {
		return err
	}
	return nil
}

func textToArray(texts string) ([]string, error) {

	// 检查是否包含 [ 和 ] 字符
	startIndex := strings.Index(texts, "[")
	endIndex := strings.LastIndex(texts, "]")

	// 安全检查: 如果没有找到 [ 或 ]
	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return nil, fmt.Errorf("输入无效: 未找到有效的 JSON 数组,text：%v", texts)
	}

	// 安全检查: 如果字符串切片超出范围
	if startIndex < 0 || endIndex+1 > len(texts) {
		return nil, fmt.Errorf("输入无效: 索引超出范围,text：%v", texts)
	}

	// 将原始文本字符串转换为 JSON 数组
	var jsonArray []string
	err := json.Unmarshal([]byte(texts[startIndex:endIndex+1]), &jsonArray)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v,text：%v", err, texts)
	}

	return jsonArray, nil
}

func (k *KnowledgeBaseService) QueryKnowledgies(question string, refreshCache bool, userId int) ([]model.Documents, error) {

	// 限制用户问题不可超过50个字
	if len(question) > 50 || question == "" {
		return nil, errors.New("问题不合规")
	}

	var questionCache QuestionCacheService
	answers := questionCache.GetAnswer(question)
	if !refreshCache {
		if len(answers) > 0 {
			var documents []model.Documents
			for i := range answers {
				answer := answers[i]
				documents = append(documents, model.Documents{
					Content: answer.Content,
					Type:    answer.Type,
					Link:    answer.Link,
					ID:      answer.DocumentId,
					Remark:  answer.Remark,
					Answer:  answer.Answer,
				})
			}
			return documents, nil
		}
	}

	// 向量化
	var embeddingService EmbeddingService
	// 向量化
	embeddings, err := embeddingService.GetTextEmbeddings([]string{question})

	if err != nil {
		return nil, nil
	}

	item := embeddings.Output.Embeddings[0]
	var vectorDao dao.VectorDao
	query, err := vectorDao.Query(item.Embedding, 20, 0.2, 1)
	if err != nil {
		return nil, err
	}

	// 组装 A
	output := "Q:" + question + "\n"
	output += "A:\n"

	for i := range query {
		vectors := query[i]

		output += "-" + vectors.Content + "\t" + strconv.Itoa(vectors.DocumentId) + "\n"

	}

	var llm LLMService
	chat, err := llm.Chat(constant.QAPrompt, output)
	if err != nil {
		return nil, err
	}

	var documentDao dao.DocumentDao
	split := strings.Split(chat, ",")
	ids := make([]int, len(split))
	for i := range split {
		s := strings.TrimSpace(split[i])
		atoi, err := strconv.Atoi(s)
		if err != nil {
			log.Warn("chat QA 返回 id 无法解析：%v", err.Error())
			continue
		}
		ids[i] = atoi
	}
	documents, err := documentDao.ListById(ids)
	if err != nil {
		return nil, err
	}

	// 创建 channel 来收集结果
	type result struct {
		index  int
		answer string
		err    error
	}
	resultCh := make(chan result, len(documents))

	// 遍历 documents，启动 Goroutine 处理每个文档
	for i := range documents {
		go func(i int) {
			m := documents[i]
			q := "问题：" + question + "\n"
			content := "内容：" + m.Content
			answer, err := llm.Chat(constant.KnowledgeQAPrompt, q+content)
			documents[i].Content = ""
			resultCh <- result{i, answer, err} // 将结果传回 channel
		}(i)
	}

	// 清洗
	var results []model.Documents

	// 收集每个 Goroutine 的结果
	for i := 0; i < len(documents); i++ {
		res := <-resultCh
		if res.err != nil {
			log.Warnf("Error in LLM chat: %v", res.err.Error())
			continue
		}
		if !strings.Contains(res.answer, "没有检索到相关数据") {
			documents[res.index].Answer = res.answer // 更新文档的答案
			results = append(results, documents[res.index])
		}
	}

	close(resultCh) // 关闭 channel

	// 添加缓存
	if len(results) > 0 {
		go func() {
			var answers []model.AnswerCaches
			for i := range results {
				r := results[i]
				answers = append(answers, model.AnswerCaches{
					Answer:     r.Answer,
					Type:       r.Type,
					Link:       r.Link,
					Content:    r.Content,
					Remark:     r.Remark,
					DocumentId: r.ID,
				})
			}
			questionCache.Put(question, answers)
		}()

	}

	return results, nil
}

func (k *KnowledgeBaseService) QueryDocumentDetail(id int) (model.Documents, error) {
	var documentDao dao.DocumentDao
	document, err := documentDao.Get(id)
	return document, err
}

func (k *KnowledgeBaseService) QueryKnowledge(question, content string) (string, error) {

	question = "问题：" + question
	content = "内容：" + content

	var llm LLMService
	chat, err := llm.Chat(constant.KnowledgeQAPrompt, question+content)
	if err != nil {
		return "", err
	}
	return chat, nil
}

// 加载相似 qa 的函数
func (k *KnowledgeBaseService) AddQA(title string, id int) error {

	var embeddingService EmbeddingService

	embeddings, err := embeddingService.GetTextEmbeddings([]string{title})
	if err != nil {
		return err
	}

	var vectors []model.Vectors
	items := embeddings.Output.Embeddings
	for i := range items {
		item := items[i]
		vectors = append(vectors, model.Vectors{
			Content:    title,
			DocumentId: id,
			Embedding:  pgvector.NewVector(item.Embedding),
			Type:       2,
		})

	}
	var vectorDao dao.VectorDao
	// 先删除
	if err = vectorDao.Delete2(title, id); err != nil {
		return err
	}

	err = vectorDao.Inserts(vectors)
	return err
}
