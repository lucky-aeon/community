import ConsumerLayoutVue from "@/views/layout/ConsumerLayout.vue";
import { createRouter, createWebHistory } from "vue-router";
import ConsumerRouters from "./consumer";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: ConsumerLayoutVue,
            children: ConsumerRouters
        }
    ]
})

export default router