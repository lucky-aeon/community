<template>
    <div v-if="route.name != 'qaView'"
        style="padding-bottom: 20px;">
        <ARow justify="center" :gutter="[20, 20]">
            <ACol>
                <ContentSearchCom v-model="searchData" :handler-search-article="handlerSearchArticle" />
            </ACol>
        </ARow>
    </div>
    <router-view v-model="searchData" :articleData="currentArticleData"
                    style="margin: auto;">
                </router-view>
   
</template>
<script setup>
import ContentSearchCom from '@/components/ContentSearch.vue'
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