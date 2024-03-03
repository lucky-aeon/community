import axios from "axios";

export function getList(page) {
    return axios.get(`/community/subscription`, {
        params: {
            page: page.current,
            limit: page.defaultPageSize
        }
    })
}

export function subscribe(event, businessId) {
    return axios.post(`/community/subscribe`, {event, businessId})
}