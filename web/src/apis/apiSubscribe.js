import axios from "axios";

export function getList(page) {
    return axios.get(`/community/subscription`, {
        params: {
            page: page.current,
            limit: page.defaultPageSize
        }
    })
}

export function apiSubscribe(eventId, businessId) {
    return axios.post(`/community/subscribe`, {eventId, businessId})
}

export function apiSubscribeState(eventId,businessId) {
    const subscriptions = {
        EventId : eventId,
        BusinessId :businessId
    }
    return axios.post(`/community/subscription/state`, subscriptions)
}