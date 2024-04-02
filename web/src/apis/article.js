import axios from "axios";

/**
 * 获取文章列表(多条件查询)
 * @param {{tags?:number[],context?:string,orderBy?:string,descOrder?:boolean,userId?:number}} data - 条件
 * @param {number} page - 当前页
 * @param {number} limit - 条数
 */
export function apiArticleList(data, page, limit) {
    return axios({ method: "POST", url: "/community/articles", params: { page, limit }, data: data })
}

/**
 * 获取文章
 * @param {number} id 文章id
 * @returns article object
 */
export function apiArticleView(id) {
    return axios.get("/community/articles/" + id)
}

/**
 * 文章添加或更新接口
 * @param {object} data 文章数据
 * @param {boolean} add 添加文章？
 */
export function apiArticleUpdate(data, add = false) {
    if (add) {
        data.id = 0
    }
    if(data.tags && data.tags.length>0) {
        data.tags = data.tags.map(item=> item.value.TagId)
    }
    return axios.post(`/community/articles/update`, data)
}

export function apiArticleDelete(id) {
    return axios.delete(`/community/articles/${id}`)
}

export function apiArticleLike(id){
    return axios.post(`/community/articles/like?articleId=${id}`)
}
export function apiArticleLikeState(id){
    return axios.get(`/community/articles/like/state/${id}`)
}

export function apiAdminListArticles(page,limit){
    return axios.get(`/community/admin/article/page?page=${page}&limit=${limit}`)
}

export function apiAdminDeleteArticles(id){
    return axios.delete(`/community/admin/article/${id}`)
}

export function apiAdminUpdateArticleState(article){
    return axios.post(`/community/admin/article/state`,article)
}
export function apiAdminListArticleStates(){
    return axios.get(`/community/admin/article/states`)
}

export function apiGetTopArticle(type, page=1, limit=10) {
    return axios.get(`/community/articles/top`, {
        params:{
            type,
            page,
            limit
        }
    })
}

export function apiAutoSaveArticle(
    /** 文章id */
    articleId = 0, 
    /** 正文 */
    content="", 
    /** 标签 */
    labels=[], 
    /** 文章类型 */
    type=0) {
    return axios.post(`/community/draft`, {
        articleId,
        content,
        labels,
        type
    })
}

export function apiGetAutoSaveArticle(articleId=0) {
    let req = {
        articleId
    }
    return axios.get(`/community/draft`, {
        params: req
    })
}