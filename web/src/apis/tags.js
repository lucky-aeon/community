import axios from "axios";

/**
 * 获取前n的标签
 * @param {number} [limit=10] 条数
 **/
export function apiTagHots(limit = 10) {
    return axios.get(`/community/tags/hot`, {
        params: {
            limit
        }
    })
}

/**
 * 获取标签列表(公开的，所有)，分页
 * @param {number} page 当前页
 * @param {number} limit 条数
 * @returns 结果集
 */
export function apiTags(page = 1, limit = 10, title) {
    return axios.get(`/community/tags`, {
        params: {
            page,
            limit,
            title
        }
    })
}
/**
 * 获取标签下的文章数量
 * @returns object {"code": 200,"data": [{"ArticleCount": 3,"TagId": 1,"TagName": "java"},{"ArticleCount": 2,"TagId": 2,"TagName": "go"}],"msg": "","ok": true} 
 */
export function apiUserTagArticleCount() {
    return axios.get(`/community/tags/getTagArticleCount`)
}