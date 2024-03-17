import axios from 'axios';




export function apiAdminFile(page, limit) {
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
export async function apiGetUploadPolicy() {
  return axios.get(`/community/file/policy`)
}
export async function apiUploadFile(userId, file, callback) {
  let policy = await apiGetUploadPolicy()
  console.log(policy)
  if (!policy.ok) return Promise.reject("无法获取授权")
  let key = `${userId}/${new Date().getTime()}`
  let result = await axios.postForm(`https://luckly-community.${policy.data.host}`, {
    OSSAccessKeyId: policy.data.accessid,
    policy: policy.data.policy,
    Signature: policy.data.signature,
    key: key,
    callback: policy.data.callback,
    file: file
  }, {
    timeout: 999999
  })
  console.log(result )
  callback(result, key)
}