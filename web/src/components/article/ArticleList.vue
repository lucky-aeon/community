<template>
    <a-list v-if="dataSource.length" class="list-demo-action-layout" :bordered="false">
        <a-list-item  v-for="(item, index) in dataSource" :key="item.id" @click="router.push(`/article/view/${item.id}`)"
            class="list-demo-item" action-layout="horizontal">
            <template #actions>
                <a-button-group @click.stop="() => { }" v-if="userStore.userInfo.id == item.user.id">
                    <a-button @click="editArticle = { id: item.id, show: true, index: index }">编辑</a-button>
                    <a-dropdown :hide-on-select="false" :popup-max-height="false">
                        <a-button :id="`editmoreButton${index}`">
                            <template #icon>
                                <icon-down />
                            </template>
                        </a-button>

                        <template #content>
                            <a-popconfirm popup-hover-stay @ok="deleteArticle(index)" content="你确定要删除该文章?">
                                <a-doption>删除文章</a-doption>
                            </a-popconfirm>
                        </template>
                    </a-dropdown>
                </a-button-group>
            </template>
            <a-list-item-meta :title="item.title">

                <template #description>

                    <a-tag size="small">
                        分类: {{ item.type.title || "未知" }}
                    </a-tag>
                    <a-divider direction="vertical" />
                    <a-tag v-if="!item.tags" color="blue" size="small">
                        无标签
                    </a-tag>
                    <template v-else>
                        <a-space>
                            <a-tag v-for="tagItem in item.tags.split(',')" :key="tagItem" color="blue" size="small">{{
        tagItem }}</a-tag>
                        </a-space>
                    </template>
                    <br />
                    <span><icon-user />{{ item.user.name || "未知" }}</span>
                    <a-divider direction="vertical" />
                    <span><icon-heart />{{item.like}}</span>
                    <a-divider direction="vertical" />
                    <span><icon-calendar />{{ item.updatedAt }}</span>
                </template>

                <template #avatar>
                    <a-avatar shape="square" :image-url="item.user.avatar">
                    </a-avatar>
                </template>
            </a-list-item-meta>
        </a-list-item>
    </a-list>
    <AResult v-else title="no articles">

        <template #icon>
            <IconEmpty />
        </template>
    </AResult>
    <ArtilceEdit v-model="editArticle.show" :article-id="editArticle.id" :call-response="refreshList" />

</template>

<script setup>
import { apiArticleDelete, apiArticleList } from '@/apis/article';
import ArtilceEdit from '@/components/article/ArticleEdit.vue';
import router from '@/router';
import { useUserStore } from '@/stores/UserStore';
import { IconCalendar, IconDown, IconEmpty, IconHeart, IconUser } from '@arco-design/web-vue/es/icon';
import { onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps({
    queryData: {
        type: Object,
        default: () => ({
            tags: [],
            context: "",
            orderBy: "created_at",
            descOrder: true
        })
    },

})
const editArticle = ref({
    show: false,
    id: 0,
    index: 0
})
const userStore = useUserStore()
const dataSource = ref([])
dataSource.value = []
const paginationProps = reactive({
    page: 1,
    noReq: false,
    defaultPageSize: 10,
    total: 9999
})
const route = useRoute()
function refreshList(data, ok) {
    if (!ok) {
        return
    }
    let currentItem = dataSource.value[editArticle.value.index]
    if (currentItem.id != data.id) {
        return
    }
    currentItem.title = data.title
    currentItem.type = data.type
    currentItem.tags = data.tags.map(item => item.name).join(',')
}
function getArticleList(clean=false) {

    console.trace("1")
    apiArticleList(Object.assign(props.queryData, { context: route.query.context, tags: (typeof route.query.tags == "string" ? [route.query.tags] : route.query.tags) || [] }), paginationProps.page, paginationProps.defaultPageSize).then(({ data }) => {
        if(clean) dataSource.value = []
        paginationProps.total = data.total
        if(data.list == null) {
            return
        }
        dataSource.value.push(...data.list)
    })
}
function deleteArticle(index) {
    apiArticleDelete(dataSource.value[index].id).then(({ok})=>{
        if(!ok) return
        dataSource.value.splice(index, 1)
    })
}
const scroll = () => {
    const scrollHeight = document.documentElement.scrollHeight // 可滚动区域的高
    const scrollTop = document.documentElement.scrollTop // 已经滚动区域的高
    const clientHeight = document.documentElement.clientHeight // 可视区高度
    // 以滚动高度 + 当前视口高度  >= 可滚动高度 = 触底
    if (clientHeight + scrollTop >= scrollHeight - 0.5 && !paginationProps.noReq) {
        // 此处可书写触底刷新代码

        paginationProps.noReq = true
        if (paginationProps.page >= Math.ceil(paginationProps.total / paginationProps.defaultPageSize)) {
            return
        }
        setTimeout(() => paginationProps.noReq = false, 1000)
        paginationProps.page++
        getArticleList()
    }
}
onMounted(() => {
    window.addEventListener('scroll', scroll)
})
// 页面销毁移除scroll事件
onUnmounted(() => window.removeEventListener('scroll', scroll))
watch(() => route.fullPath, () => {
    dataSource.value = []
    paginationProps.page = 1
    getArticleList(true)
}, {immediate: true})
watch(()=> props.queryData.state, ()=>{
    dataSource.value = []
    getArticleList()
})
</script>

<style scoped>
.list-demo-action-layout .image-area {
    width: 183px;
    height: 119px;
    border-radius: 2px;
    overflow: hidden;
}

.list-demo-action-layout .list-demo-item {
    padding: 20px 0;
    border-bottom: 1px solid var(--color-fill-3);
}

.list-demo-action-layout .image-area img {
    width: 100%;
}

.list-demo-action-layout .arco-list-item-action .arco-icon {
    margin: 0 4px;
}


.list-demo-item:hover {
    transition: background-color 0.3s ease;
    background-color: #f4f4f491;
    cursor: pointer;
}
</style>