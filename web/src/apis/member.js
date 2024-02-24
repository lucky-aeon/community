import axios from 'axios';

export function listAllMember() {
    return axios.get(`/community/admin/member`);
}

export function saveMember(member) {
    console.log(member)
    return axios.post(`/community/admin/member`,member);
}

export function deleteMember(id) {
    return axios.delete(`/community/admin/member/${id}`);
}

