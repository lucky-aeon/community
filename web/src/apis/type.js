import axios from 'axios';

export function apiListParentAllType() {
    return axios.get(`/community/admin/type/parent`);
}
export function apilistAllType() {
    return axios.get(`/community/admin/type`);
}

export function apiSaveType(member) {
    return axios.post(`/community/admin/type`,member);
}



export function apiDeleteType(id) {
    return axios.delete(`/community/admin/type/${id}`);
}

