<template>
    <div v-if="route.name != 'articleView'"
        style="border-bottom: 1px solid rgba(215, 215, 215, 0.784); padding-bottom: 20px;">
        <ARow justify="center" :gutter="[20, 20]">
            <ACol>
                <ContentSearchCom v-model="searchData" :handler-search-article="handlerSearchArticle"/>
            </ACol>
        </ARow>
    </div>

    <router-view :articleData="currentArticleData" style="margin: auto;max-width: 1100px !important;min-width: 600px;">
    </router-view>

</template>
<script setup>
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

const handlerSearchArticle = (value) => {
}


</script>

