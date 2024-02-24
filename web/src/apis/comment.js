import axios from 'axios';

export function listAllCommentsByArticleId(articleId) {
  return axios.get(`/community/comments/allCommentsByArticleId/${articleId}`);
}

export function deleteComment(id) {
  return axios.delete(`/community/comments/${id}`);
}
