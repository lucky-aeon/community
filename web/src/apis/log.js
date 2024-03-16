import axios from "axios";



export function apiOperLogList(page,limit,searchData) {
    return axios({ method: "GET", url: "/community/admin/oper/log", params: { page, limit,...searchData}})
}



