import axios from 'axios';

/**
 * 获取文章下的所有评论
 * @param {number} articleId 文章id
 * @returns 
 */
export function listAllCommentsByArticleId(articleId) {
  return axios.get(`/community/comments/allCommentsByArticleId/${articleId}`);
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
  let url = `/byArticleId/${parentId}`
  if (!root) {
    url = `/byRootId/${parentId}`
    return axios.get(url, {
      params: {
        page,
        limit,
      }
    })
  }
}

export function apiPublishArticleComment(parentId, data, root = true) {
  axios.post(`/community/comments/`)
}
