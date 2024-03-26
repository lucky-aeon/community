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

export function apiMessageTemplateList(page,limit) {
    return axios.get(`/community/admin/message/template?page=${page}&limit=${limit}`)
}

export function apiMessageTemplateEventList() {
    return axios.get(`/community/admin/message/template/event`)
}

export function apiMessageTemplateVarList(page,limit) {
    return axios.get(`/community/admin/message/template/var`)
}

export function apiMessageTemplateSave(template) {
    return axios.post(`/community/admin/message/template`,template)
}

export function apiMessageTemplateListDelete(id) {
    return axios.delete(`/community/admin/message/template/${id}`)
}
