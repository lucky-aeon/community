import axios from "axios";

/**
 * 获取文章列表(多条件查询)
 * @param {{tags?:number[],context?:string,orderBy?:string,descOrder?:boolean}} data - 条件
 * @param {number} page - 当前页
 * @param {number} limit - 条数
 */
export function apiArticleList(data, page, limit) {
    return axios({method: "POST", url: "/community/articles", params: {page, limit}, data: data})
}

/**
 * 获取文章
 * @param {number} id 文章id
 * @returns article object
 */
export function apiArticleView(id) {
    return axios.get("/community/articles/"+id)
}