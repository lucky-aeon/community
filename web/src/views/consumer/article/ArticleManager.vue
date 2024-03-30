<template>
    <ARow style="margin-top: 10px;" :gutter="10">
        <ACol :span="19">
            <a-row>
                <a-col :span="24">
                    <a-space size="small">
                        <ACard :bordered="false" hoverable><a-statistic title="喜欢人数" :value="userStatistics.likeCount"
                                show-group-separator />
                        </ACard>
                        <ACard :bordered="false" hoverable><a-statistic title="文章数量" :value="userStatistics.articleCount"
                                show-group-separator />
                        </ACard>
                    </a-space>
                </a-col>
                <a-col :span="24">
                    <div style="float: right;padding: 5px;">
                    <a-radio-group v-model="queryData.state" type="button">
                        <a-radio :value="2">已发布</a-radio>
                        <a-radio :value="1">草稿</a-radio>
                    </a-radio-group>
                </div>
                </a-col>
                <a-col :span="24">
                    <a-card :bordered="false">
                    <ArticleListCom v-if="showList" style="margin-top: 10px;" :queryData="queryData" />
                </a-card>
                </a-col>
            </a-row>

        </ACol>
        <ACol :span="5">
            <ASpace direction="vertical" style="width: 100%;">
                <ACard size="small" hoverable :bordered="false">
                    <AButton long type="primary" @click="editArticle.show = true">发布文章</AButton>
                </ACard>

                <ACard :bordered="false" title="标签" :extra="userStore.userTags.length">
                    <ASpace wrap>
                        <ATag v-for="tagItem in userStore.userTags" :key="tagItem.TagId">{{ tagItem.TagName }}({{
        tagItem.ArticleCount }})</ATag>
                    </ASpace>
                </ACard>
            </ASpace>
        </ACol>
    </ARow>
    <ArticleEditCom v-model="editArticle.show" :article-id="editArticle.currentId" :call-response="refreshPage" />
</template>
<script setup>
import { apiGetUserStatistics } from "@/apis/user";
import ArticleEditCom from "@/components/article/ArticleEdit.vue";
import ArticleListCom from "@/components/article/ArticleList.vue";
import { useUserStore } from "@/stores/UserStore";
import { nextTick, reactive, ref } from "vue";
// const IconFont = Icon.addFromIconFontCn({ src: 'https://at.alicdn.com/t/c/font_4443211_9ji0rjy60kw.js' });
const userStore = useUserStore()
const queryData = ref({
    tags: [],
    context: "",
    orderBy: "created_at",
    descOrder: true,
    state: 2,
    userID: userStore.userInfo.id
})
const userStatistics = ref({
    likeCount: 0,
    articleCount: 0
});
const editArticle = reactive({
    show: false,
    currentId: 0
})
const showList = ref(false)
function refreshPage() {
    apiGetUserStatistics(true).then(({ data }) => {
        userStatistics.value = data
    })
    showList.value = false
    nextTick(()=> showList.value = true)
}
refreshPage()
</script>