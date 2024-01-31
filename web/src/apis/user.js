import axios from 'axios';

export function getUserMenu() {
  return axios.get('/api/community/user/menu');
}
