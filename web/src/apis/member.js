import axios from 'axios';

export function listAllMember(page,limit) {
    return axios.get(`/community/admin/member?page=${page}&limit=${limit}`);
}

export function saveMember(member) {
    return axios.post(`/community/admin/member`,member);
}

export function deleteMember(id) {
    return axios.delete(`/community/admin/member/${id}`);
}

