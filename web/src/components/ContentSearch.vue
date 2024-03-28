<template>
    <div class="search">
        <ARow :gutter="[0, 10]">
            <ACol>
                <h1>搜索内容</h1>
                <a-input-group style="height: 36px;max-height: 36px;width: 400px;">
                    <a-badge :count="model.tags.length" :offset="[-50, 15]">
                        <AButton size="large" type="dashed" @click="searchData.tagModal.show = true">
                            <IconTags />
                        </AButton>
                    </a-badge>
                    <a-auto-complete :model-value="model.context" allowClear
                        @select="(e) => router.push(`/article/view/${e}`)" :data="searchData.result"
                        :style="StyleSearch" @search="searchArticle" placeholder="please enter something">

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
                </a-space>
            </ACol>
        </ARow>
        <a-modal hideCancel v-model:visible="searchData.tagModal.show" @ok="() => { searchData.tagModal.show = false }">
            <template #title>
                搜索标签
            </template>
            <div>
                <a-auto-complete v-model:model-value="searchData.tagModal.select" allowClear
                    @select="searchData.tagModal.handlerSelect" placeholder="搜索标签"
                    :data="noselectTags" @search="handleSearchTags">
                    <template #data="{ data }">
                        {{ data.title }}</template>
                </a-auto-complete>
                <a-space v-if="model.tags.length" style="margin: 10px;">
                    <a-tag v-on:close="searchData.tagModal.rmTag(index)" closable
                        v-for="(tagitem, index) in [...new Set(model.tags)]" :key="tagitem">{{ tagitem
                        }}</a-tag>
                </a-space>
                <AResult v-else title="no tag condition">
                    <template #icon>
                        <IconEmpty />
                    </template>
                </AResult>
            </div>
        </a-modal>
    </div>

</template>

<script setup>
import { apiArticleList } from '@/apis/article';
import { apiTagHots, apiTags } from '@/apis/tags';
import router from '@/router';
import { IconEmpty, IconSearch, IconTags } from '@arco-design/web-vue/es/icon';
import { watch, onMounted,computed,reactive,nextTick } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const StyleSearch = reactive({
    width: '100%',
    height: '100%',
})
const { handlerSearchArticle } = defineProps({
    handlerSearchArticle: {
        type: Function,
        default(content) { }
    }
})
const model = defineModel({
    required: true,
    type: {
        tags: [],
        context: "",
        state: 2,
        orderBy: "created_at",
        descOrder: true
    },
    default: {
        tags: [],
        context: "",
        state: 2,
        orderBy: "created_at",
        descOrder: true
    }
})
const searchData = reactive({
    hotTags: [],
    tagModal: {
        show: false,
        searchResult: [],
        select: "",
        handlerSelect(e) {
            model.value.tags.push(e)
            model.value.tags = [...new Set(model.value.tags)]
            nextTick(() => searchData.tagModal.select = "")
        },
        rmTag(i) {
            model.value.tags.splice(i, 1)
        }
    },
    result: []
})
const noselectTags = computed(()=> searchData.tagModal.searchResult.filter((tag) => {
            return model.value.tags.find(item=> item == tag) == undefined
        }))
const handleSearchTags = (value) => {
    apiTags(1, 15, value).then(({ data, ok }) => {
        if (!ok) {
            return
        }
        searchData.tagModal.searchResult = data.list
    })
}

function ToSearchRoute() {
    let target = {
        path: route.path,
        query: {
            tags: model.value.tags,
            context: model.value.context
        }
    }
    if(model.value.state != 2) target.query.state = model.value.state
    router.push(target)
}


const searchArticle = (value) => {
    model.value.context = value

    apiArticleList(model.value, 1, 5).then(({ data, ok }) => {
        handlerSearchArticle(data, ok)
        if (!ok || !data || !data.list) {
            return
        }
        searchData.result = data.list.map(data => ({
            label: data.title,
            value: data.id
        }))
    })
}

onMounted(() => {
    apiTagHots().then(({ data }) => {
        searchData.hotTags = data
        searchData.tagModal.searchResult = data.map(({ tag }) => {
            return tag
        })
    })
    ToSearchRoute()
    searchArticle(model.value.context)
})

watch(()=> model.value.state, ()=>{
    nextTick(()=> ToSearchRoute())
})
watch(() => route.query, (newV) => {
    if (route.query.tags){
        let tags = route.query.tags
        if(typeof tags == 'string') {
            tags = [tags]
        }
        model.value.tags = [...new Set(tags)]
    }
    if(route.query.state) {
        model.value.state = parseInt(route.query.state)
    }
    if (route.query.context)
        model.value.context = route.query.context
}, { immediate: true })
</script>

<style scoped>
.search {
    width: 100%;
    text-align: center;
}
</style>