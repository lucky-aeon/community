import axios from 'axios';

export function getUserMenu() {
  return axios.get('/community/user/menu');
}
