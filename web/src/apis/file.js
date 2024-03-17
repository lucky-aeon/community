import axios from 'axios';




export function apiAdminFile(page,limit) {
  return axios.get(`/community/admin/file?page=${page}&limit=${limit}`)
}

export function apiGetFile(fileKey) {
  return axios.getUri({
    method: "GET",
    url: `/community/file/singUrl`,
    params: {
      fileKey
    }
  })
}