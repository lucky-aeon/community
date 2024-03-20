import axios from "axios";


export function apiCodeList(page,limit,startTime,endTime) {
    return axios.get(`/community/admin/code?page=${page}&limit=${limit}`,)
}

export function apiGenerateCode(code) {
    return axios.post(`/community/admin/code/generate`,code)
}

export function apiDeleteCode(code) {
    return axios.delete(`/community/admin/code/${code}`)
}