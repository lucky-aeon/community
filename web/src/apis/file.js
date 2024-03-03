import axios from 'axios';




export function apiAdminFile(page,limit) {
  return axios.get(`/community/admin/file?page=${page}&limit=${limit}`)
}