import axios from "axios";


export function apiCodeList() {
    return axios.get(`/community/admin/code`)
}

export function apiGenerateCode(code) {
    return axios.post(`/community/admin/code/generate`,code)
}

export function apiDeleteCode(code) {
    return axios.delete(`/community/admin/code/${code}`)
}