import { apiGetUnReadCount } from "@/apis/message";
import { defineStore } from "pinia";

export const useMessage = defineStore('MessageStore', {
    state: ()=>({
        unReadCount: 0,
        msgs: []
    }),
    actions: {
        getUnReadCount() {
            apiGetUnReadCount().then(({data, ok})=>{
                if(!ok) return
                this.unReadCount = data
            })
        }
    }
})