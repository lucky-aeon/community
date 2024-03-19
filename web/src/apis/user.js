import axios from 'axios';

export function getUserMenu() {
  return axios.get('/community/user/menu');
}

export function getUserInfo() {
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

export function apiAdminListUsers() {
  return axios.get(`/community/admin/user`)
}

export function apiAdminUpdateUsers(user) {
  return axios.post(`/community/admin/user`,user)
}
