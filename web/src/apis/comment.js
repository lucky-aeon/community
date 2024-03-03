import axios from 'axios';

export function listAllCommentsByArticleId(articleId,page,limit) {
  return axios.get(`/community/comments/allCommentsByArticleId/${articleId}?page=${page}&limit=${limit}`);
}

/**
 * 删除评论
 * @param {number} id 评论id
 * @returns 
 */
export function deleteComment(id) {
  return axios.delete(`/community/comments/${id}`);
}
export function apiReply(data){
  return axios.post("/community/comments/comment",data)
}
/**
 * 获取指定父对象的子评论
 * @param {number} parentId 父评论id
 * @param {number} page 当前页
 * @param {number} limit 条数
 * @param {boolean} root 文章根评论
 * @returns 评论列表
 */
export function apiGetArticleComment(parentId, page = 1, limit = 10, root = true) {
  let url = `/community/comments/byArticleId/${parentId}`
  if (!root) {
    url = `/community/comments/byRootId/${parentId}`
  }
  return axios.get(url, {
      params: {
        page,
        limit,
      }
    })
}

export function apiPublishArticleComment(articleId, parentId, content, root=0) {
  let data = {
    content,
    articleId
  }
  if(root) {
    data.parentId = parentId
    data.rootId = root
  }
  axios.post(`/community/comments/comment`, data)
}
