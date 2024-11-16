package dao

import (
	"fmt"
	"github.com/pgvector/pgvector-go"
	"log"
	model "xhyovo.cn/community/server/model/knowledge"
)

type VectorDao struct{}

type SimilarityQA struct {
	Id      int
	QaTitle string
}

// 获取向量
func (v *VectorDao) Get(id int) (*model.Vectors, error) {
	var vector model.Vectors
	if err := model.Vector().First(&vector, id).Error; err != nil {
		return nil, err
	}
	return &vector, nil
}

func (v *VectorDao) List(embeddings []float32) ([]model.Vectors, error) {
	var vectors []model.Vectors
	// 这里假设你有一个适当的向量搜索实现
	if err := model.Vector().Where("embedding @> ?", embeddings).Find(&vectors).Error; err != nil {
		return nil, err
	}
	return vectors, nil
}

func (v *VectorDao) Insert(vectorModel model.Vectors) error {
	return model.Vector().Create(&vectorModel).Error
}

func (v *VectorDao) Inserts(vectorModel []model.Vectors) error {
	return model.Vector().Create(&vectorModel).Error
}

func (v *VectorDao) Delete(id int) error {
	return model.Vector().Delete(&model.Vectors{}, id).Error
}

func (v *VectorDao) DeleteByDocumentId(documentId int) error {
	return model.Vector().Where("document_id", documentId).Delete(model.Documents{}).Error
}

func (v *VectorDao) Query(embedding []float32, limit int, threshold float64, typee int) ([]model.Vectors, error) {
	var vectors []model.Vectors
	vector := pgvector.NewVector(embedding)

	// 构建查询，使用子查询计算余弦相似度，并在外层设置阈值筛选
	query := model.Vector().Raw(`
		SELECT *
		FROM (
			SELECT content, document_id,type, embedding <=> ? AS similarity
			FROM vectors
			ORDER BY similarity
			LIMIT ?
		) AS subquery
		where type = ?
		ORDER BY similarity
		LIMIT ?`, vector, limit, typee, limit)

	// 执行查询
	err := query.Scan(&vectors).Error
	if err != nil {
		log.Printf("Error querying vectors with threshold: %v", err)
		return vectors, fmt.Errorf("failed to query vectors: %w", err) // 更清晰的错误处理
	}

	return vectors, nil
}

func (v *VectorDao) Delete2(content string, documentId int) error {
	return model.Vector().Where("document_id", documentId).Where("content", content).Where("type", 2).Delete(model.Documents{}).Error

}
