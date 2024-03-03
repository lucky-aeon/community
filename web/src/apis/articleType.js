import axios from "axios";

/**
 * 获取文章分类列表
 */
export function apiGetArticleTypes() {
    return axios.get('/community/classfiys?parentId=1')
}