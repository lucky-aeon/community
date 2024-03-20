const EventIds = [
// id: 1, name: 文章评论
    "/article/view/",
// id: 2, name: 用户关注
    "/user/profile/"
]

export function GetEventUrl(eventId, toId) {
    if(!(eventId >0 && eventId <= EventIds.length)) {
        return ""
    }
    return EventIds[eventId-1]+toId
}

