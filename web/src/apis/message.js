import axios from 'axios';

export function apiListMsg(page,limit,type, state) {
    return axios.get(`/community/message?type=${type}&state=${state}&page=${page}&limit=${limit}`);
}

export function apiClearUnReadMsg(type){
    return axios.delete(`/community/message/UnReadMsg/${type}`);
}

export function apiGetUnReadCount() {
    return axios.get(`/community/message/unReader/count`)
}
export function apiPostRead(ids) {
    return axios.post(`/community/message/read`, ids)
}