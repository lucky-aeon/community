<template>
    <div v-if="route.name != 'articleView'"
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
import { apiGetTopArticle } from '@/apis/article';
import ContentSearchCom from '@/components/ContentSearch.vue'
import { ref } from 'vue';
import { reactive } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const currentArticleData = ref({})
const searchData = reactive({
    tags: [],
    context: "",
    state: 2,
    orderBy: "created_at",
    descOrder: true
})
apiGetTopArticle(route.params.classfily, 1,)
const handlerSearchArticle = (value) => {
}


</script>
