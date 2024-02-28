import axios from 'axios';




export function apiAdminFile() {
  return axios.get(`/community/admin/file`)
}