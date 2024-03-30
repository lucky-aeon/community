<template>
    <div v-if="route.name != 'qaView'"
        style="border-bottom: 1px solid rgba(215, 215, 215, 0.784); padding-bottom: 20px;">
        <ARow justify="center" :gutter="[20, 20]">
            <ACol>
                <ContentSearchCom v-model="searchData" :handler-search-article="handlerSearchArticle" />
            </ACol>
            <ACol>
                <div style="float: right;">
                <a-radio-group v-model="searchData.state" type="button" :default-value="3">
                        <a-radio :value="3">待解决</a-radio>
                        <a-radio :value="4">已解决</a-radio>
                        <a-radio :value="5" v-if="userStore.userInfo.role == 'admin'">私密提问</a-radio>
                    </a-radio-group> 
                </div>
            </ACol>
        </ARow>
    </div>
    <router-view v-model="searchData" :articleData="currentArticleData" style="margin: auto;max-width: 1100px !important;min-width: 600px;">
    </router-view>
</template>
<script setup>
import ContentSearchCom from '@/components/ContentSearch.vue'
import { useUserStore } from '@/stores/UserStore';
import { reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const searchData = reactive({
        state: route.query.state || 3,
        tags: [],
        context: "",
        orderBy: "created_at",
        descOrder: true
    })
const userStore = useUserStore()
const currentArticleData = ref({})
const handlerSearchArticle = (value) => {
}

</script>

<style scoped>
.search {
    width: 100%;
    text-align: center;
}
</style>