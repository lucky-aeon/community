import axios from 'axios';

export function listMsg(type,state) {
    return axios.get(`/community/message?type=${type}&state=${state}`);
}


