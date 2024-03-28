import axios from 'axios';

export function getUserMenu() {
  return axios.get('/community/user/menu');
}

export function getUserInfo(userId=undefined) {
  if (userId) {
    return axios.get('/community/user/info', {
      params: {userId}
    });
  }
  return axios.get('/community/user/info');
}

export function saveUserInfo(tab,user) {
  return axios.post(`/community/user/edit/${tab}`,user);
}

/**
 * 登录或注册，code.length==8是注册
 * @param {{account: string,password: string, code?: string, name?: string}} authForm 
 * @returns 
 */
export function apiAuthAccount(authForm) {
  if(authForm.code) {
    return axios.post(`/community/register`, authForm)
  }
  return axios.post(`/community/login`, authForm)
}

/**
 * 用户数据统计
 * @returns {"code":200,"data":{"articleCount":2,"likeCount":1},"msg":"","ok":true}
 */
export function apiGetUserStatistics() {
  return axios.get(`/community/user/statistics`)
}

export function apiAdminListUsers(page,limit) {
  return axios.get(`/community/admin/user?page=${page}$limit=${limit}`)
}

export function apiAdminUpdateUsers(user) {
  return axios.post(`/community/admin/user`,user)
}

export function apiAdminListUserTags(page,limit) {
  return axios.get(`/community/admin/user/tag?page=${page}&limit=${limit}`)
}

export function apiAdminSaveUserTags(userTags) {
  return axios.post(`/community/admin/user/tag?`,userTags)
}
export function apiAdminDeleteUserTags(id) {
  return axios.delete(`/community/admin/user/tag/${id}`)
}

/**
 * 获取用户标签
 * @param {number} id 用户id
 * @returns 标签集
 */
export function apiGetUserTags(id) {
  return axios.get(`/community/admin/user/tag/${id}`)
}

export function apiGetUserTags2(id) {
  return axios.get(`/community/user/tags/${id}`)
}

export function apiSearchUserByName(name="") {
  return axios.get(`/community/user`, {
    params: {
      name
    }
  })
}

export function apiUserEditUserAvatar(avatar) {
  return axios.post(`/community/user/edit/avatar`, {
    avatar
  })
}