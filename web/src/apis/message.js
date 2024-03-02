import axios from 'axios';

export function apiListMsg(type, state) {
    return axios.get(`/community/message?type=${type}&state=${state}`);
}


