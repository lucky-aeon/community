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
