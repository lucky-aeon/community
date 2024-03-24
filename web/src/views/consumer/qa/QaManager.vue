<template>
    <ARow style="margin-top: 10px;" :gutter="10">
        <ACol :span="19">
            <a-row>
                <a-col :span="24">
                    <a-space size="small">
                        <ACard hoverable><a-statistic title="讨论人数" :value="userStatistics.likeCount"
                                show-group-separator />
                        </ACard>
                        <ACard hoverable><a-statistic title="问答数量" :value="userStatistics.articleCount"
                                show-group-separator />
                        </ACard>
                    </a-space>
                </a-col>
                <a-col :span="24">
                    <div style="float: right;padding: 5px;">
                    <a-radio-group v-model="queryData.state" type="button">
                        <a-radio :value="3">待解决</a-radio>
                        <a-radio :value="5">私密提问</a-radio>
                        <a-radio :value="4">已解决</a-radio>
                        <a-radio :value="1">草稿</a-radio>
                    </a-radio-group>
                </div>
                </a-col>
                <a-col :span="24">
                    <QaListCom v-if="showList" style="margin-top: 10px;" :queryData="queryData" />
                </a-col>
            </a-row>
        </ACol>
        <ACol :span="5">
            <ASpace direction="vertical" style="width: 100%;">
                <ACard size="small" hoverable>
                    <AButton long type="primary" @click="editArticle.show = true">发布问答</AButton>
                </ACard>

                <ACard title="标签" :extra="userStore.userTags.length">
                    <ASpace wrap>
                        <ATag v-for="tagItem in userStore.userTags" :key="tagItem.TagId">{{ tagItem.TagName }}({{
        tagItem.ArticleCount }})</ATag>
                    </ASpace>
                </ACard>
            </ASpace>
        </ACol>
    </ARow>
    <QaEditCom v-model="editArticle.show" :article-id="editArticle.currentId" :call-response="refreshPage" />
</template>
<script setup>
import { apiGetUserStatistics } from "@/apis/user";
import QaEditCom from "@/components/qa/QaEdit.vue";
import QaListCom from "@/components/qa/QaList.vue";
import { useUserStore } from "@/stores/UserStore";
import { nextTick, reactive, ref } from "vue";
// const IconFont = Icon.addFromIconFontCn({ src: 'https://at.alicdn.com/t/c/font_4443211_9ji0rjy60kw.js' });
const userStore = useUserStore()
const queryData = ref({
    tags: [],
    context: "",
    orderBy: "created_at",
    descOrder: true,
    state: 3,
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
    apiGetUserStatistics().then(({ data }) => {
        userStatistics.value = data
    })
    showList.value = false
    nextTick(() => showList.value = true)
}
refreshPage()
</script>