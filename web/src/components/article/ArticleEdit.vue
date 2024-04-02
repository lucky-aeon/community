<template>
    <a-modal :disabled="loading" :mask-closable="false" :fullscreen="fullScreen" :body-style="{ height: '100%' }"
        v-model:visible="visible" @cancel="handleCancel" draggable
        :modal-style="{ minWidth: '800px', maxHeight: fullScreen ? '' : '90%' }">
        <template #title>
            {{ modalTitile }}
        </template>

        <template #footer>
            <a-row>
                <a-col flex="0">
                    <a-button-group>
                        <a-button @click="fullScreen = !fullScreen">{{ fullScreen ? "窗口" : "全屏" }}</a-button>
                    </a-button-group>
                </a-col>
                <a-col flex="auto">
                </a-col>
                <a-col flex="100px">
                    <a-button-group type="primary">
                        <a-button @click="updateArticle(2)" :disabled="loading">发布</a-button>
                        <a-dropdown @select="handleSelect" :popup-max-height="false">
                            <a-button :disabled="loading">
                                <template #icon>
                                    <icon-down />
                                </template>
                            </a-button>

                            <template #content>
                                <a-doption @click="updateArticle(1)">保存草稿</a-doption>
                            </template>
                        </a-dropdown>
                    </a-button-group>

                </a-col>
            </a-row>

        </template>
        <div>
            <a-input-group style="width: 100%;">
                <a-select :default-value="data.defaultType" v-model="data.type" placeholder="分类" style="width: 150px;">
                    <a-option v-for="item in artilceTypes" :key="item.value" :value="item">{{ item.label
                        }}</a-option>
                </a-select>
                <a-input placeholder="标题" style="width: 100%;" v-model:model-value="data.title"></a-input>
            </a-input-group>
            <div style="margin-top: 5px;"></div>
            <tag-search v-model="data.tags" :default-data="data.tagIds" />
            <div style="margin-top: 5px;"></div>
            <a-scrollbar :style="`height: ${fullScreen ? '' : '350px'};overflow: auto;`">
                <markdown-edit v-model="data.content" />
            </a-scrollbar>
        </div>
    </a-modal>
</template>

<script setup>
import { apiArticleUpdate, apiArticleView, apiAutoSaveArticle, apiGetAutoSaveArticle } from '@/apis/article';
import { apiGetArticleTypes } from '@/apis/articleType';
import { IconDown } from '@arco-design/web-vue/es/icon';
import { computed, ref, watch } from 'vue';
import MarkdownEdit from '../MarkdownEdit.vue';
import TagSearch from '../TagSearch.vue';
import { Message } from '@arco-design/web-vue';
import { onMounted } from 'vue';
import { onUnmounted } from 'vue';
const props = defineProps({
    articleId: {
        type: Number,
        default: 0
    },
    callResponse: {
        type: Function,
        default() { }
    }
})
const visible = defineModel({ require: true })
const defaultData = { id: 0, content: "", type: null, tags: [], tagIds: [] }
const data = ref(Object.assign({}, defaultData))
const fullScreen = ref(false);
const artilceTypes = ref([])
const modalTitile = computed(() => props.articleId > 0 ? "编辑文章" : "添加文章")
const loading = ref(false)
const updateArticle = (state = 1) => {
    let postData = Object.assign({}, data.value)
    if (!postData.type) {
        Message.error("请选择分类")
        return
    }
    postData.type = postData.type.value
    delete postData.createdAt
    delete postData.updatedAt
    loading.value = true
    apiArticleUpdate({ ...postData, state }, props.articleId == 0).then(({ data, ok }) => {
        loading.value = false
        props.callResponse(data, ok)
        if (!ok) return
        visible.value = false;
        data.id
    })
};

const handleCancel = () => {
    visible.value = false;
    data.value = Object.assign({}, defaultData)
}

apiGetArticleTypes().then(({ data, ok }) => {
    if (!ok) return
    artilceTypes.value = data.map(item => ({ label: item.title, value: item.id }))
})
watch(() => visible.value, (newV) => {
    if (newV && props.articleId > 0) {
        apiArticleView(props.articleId).then((res) => {
            if (!res.ok) return
            let result = Object.assign({}, res.data)
            result.tagIds = result.tags.map(item => item.name)
            result.tags = []
            result.defaultType = result.type.id
            result.type = {
                label: result.type.title,
                value: result.type.id
            }
            data.value = result
        }).finally(() => {
            apiGetAutoSaveArticle(props.articleId).then(({ data, ok }) => {
                if (!ok) return
                data.value.content = data.content
                data.value.tagIds = data.labelIds
                data.value.defaultType = data.type
            })
        })
    } else {
        data.value = Object.assign({}, defaultData)
        apiGetAutoSaveArticle().then(({ data, ok }) => {
            if (!ok) return
            data.value.content = data.content
            data.value.tagIds = data.labelIds
            data.value.defaultType = data.type
        })
    }
})
// 自动保存定时任务
let autoSaveInterval = undefined;
onMounted(() => {
    autoSaveInterval = setInterval(() => {
        console.log("auto save");
        autoSave()
    }, 1000 * 10)
})

onUnmounted(() => {
    if (autoSaveInterval)
        clearInterval(autoSaveInterval)
})

function autoSave() {
    if (!visible.value) return
    if (!data.value.content && !data.value.tagIds.length == 0 && !data.value.type) return
    apiAutoSaveArticle(props.articleId, data.value.content, data.value.tags.map(tagItem => tagItem.value.TagId), data.value.type ? data.value.type.value : 0).then(({ ok }) => {
        if (!ok) return
        Message.success("自动保存成功")
    })
}
</script>