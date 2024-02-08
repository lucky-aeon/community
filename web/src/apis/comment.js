import axios from 'axios';

export function listAllCommentsByArticleId(articleId) {
  return axios.get(`/community/comments/allCommentsByArticleId/${articleId}`);
}


