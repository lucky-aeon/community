<template>
    <div v-if="route.name !== 'qaView'"
        style="border-bottom: 1px solid rgba(215, 215, 215, 0.784); padding-bottom: 20px;">
        <ARow justify="center" :gutter="[20, 20]">
            <ACol class="search">
                <ARow :gutter="[0, 10]">
                    <ACol>
                        <h1>搜索全站问题</h1>
                        <a-input-group style="height: 36px;max-height: 36px;width: 400px;">
                            <a-badge :count="searchData.search.tags.length" :offset="[-50, 15]">
                                <AButton size="large" type="dashed" @click="searchData.tagModal.show = true">
                                    <IconTags />
                                </AButton>
                            </a-badge>
                            <a-auto-complete allowClear @select="(e) => router.push(`/qa/view/${e}`)"
                                :data="searchData.result" :style="StyleSearch" @search="handlerSearchArticle"
                                placeholder="please enter something">

                            </a-auto-complete>
                            <AButton size="large" type="primary" @click="ToSearchRoute()">
                                <IconSearch />
                            </AButton>
                        </a-input-group>

                    </ACol>
                    <ACol>
                        <a-space>
                            <a-tooltip mini :content="tagItem.description" position="bottom"
                                v-for="tagItem in searchData.hotTags" :key="tagItem.tag">
                                <a-tag color="red">
                                    {{ tagItem.tag }}
                                </a-tag>
                            </a-tooltip>
                            <!-- <a-tag color="blue">
                                <template #icon>
                                    <icon-tags />
                                </template>
                                更多
                            </a-tag> -->
                        </a-space>
                    </ACol>
                </ARow>
            </ACol>
        </ARow>
    </div>

        <router-view :articleData="currentArticleData" style="margin: auto;max-width: 1100px !important;min-width: 600px;">
        </router-view>
    <a-modal hideCancel v-model:visible="searchData.tagModal.show" @ok="() => { searchData.tagModal.show = false }">
        <template #title>
            搜索标签
        </template>
        <div>
            <a-auto-complete allowClear @select="(e) => { searchData.search.tags.push(e) }" placeholder="搜索标签"
                :data="searchData.tagModal.searchResult" @search="handleSearchTags">
                <template #data="{ data }">
                    {{ data.title }}</template>
            </a-auto-complete>
            <a-space v-if="searchData.search.tags.length" style="margin: 10px;">
                <a-tag v-for="tagitem in [...new Set(searchData.search.tags)]" :key="tagitem">{{ tagitem }}</a-tag>
            </a-space>
            <AResult v-else title="no tag condition">
                <template #icon>
                    <IconEmpty />
                </template>
            </AResult>
        </div>
    </a-modal>
</template>
<script setup>
import { apiArticleList } from '@/apis/article';
import { apiTagHots, apiTags } from '@/apis/tags';
import router from '@/router';
import { IconEmpty, IconSearch, IconTags } from '@arco-design/web-vue/es/icon';
import { onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
// const route = useRoute()
// search input style
const StyleSearch = reactive({
    width: '100%',
    height: '100%',
})
const searchData = reactive({
    search: {
        state: 0,
        tags: [],
        context: "",
        orderBy: "created_at",
        descOrder: true
    },
    hotTags: [],
    tagModal: {
        show: false,
        searchResult: []
    },
    result: []
})


const route = useRoute()

const currentArticleData = ref({})
const handleSearchTags = (value) => {
    apiTags(1, 15, value).then(({ data, ok }) => {
        if (!ok) {
            return
        }
        searchData.tagModal.searchResult = data.list.map(({ tag }) => {
            return tag
        })
    })
}
const handlerSearchArticle = (value) => {
    searchData.search.context = value
    apiArticleList(searchData.search, 1, 5).then(({ data, ok }) => {
        if (!ok) {
            return
        }
        searchData.result = data.list.map(data => ({
            label: data.title,
            value: data.id
        }))
    })
}

function ToSearchRoute() {
    router.push({
        path: route.path,
        query: {
            tags: searchData.search.tags,
            context: searchData.search.context
        }
    })
}

onMounted(() => {
    apiTagHots().then(({ data }) => {
        searchData.hotTags = data
        searchData.tagModal.searchResult = data.map(({ tag }) => {
            return tag
        })
    })
    handlerSearchArticle("")
})

</script>

<style scoped>
.search {
    width: 100%;
    text-align: center;
}
</style>