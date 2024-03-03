<template>
    <ARow style="margin-top: 10px;" :gutter="10">
        <ACol :span="20">
            <a-space size="small">
                <ACard hoverable><a-statistic title="喜欢人数" :value="userStatistics.likeCount" show-group-separator /></ACard>
                <ACard hoverable><a-statistic title="文章数量" :value="userStatistics.articleCount" show-group-separator />
                </ACard>
            </a-space>
            <ArticleListCom style="margin-top: 10px;" :queryData="queryData" />
        </ACol>
        <ACol :span="4">
            <ASpace direction="vertical" style="width: 100%;">
                <ACard size="small" hoverable>
                    <AButton long type="primary" @click="editArticle.show=true">发布文章</AButton>
                </ACard>

                <ACard title="标签" extra="2 个">
                    <ASpace>
                        <ATag>标签</ATag>
                    </ASpace>
                </ACard>
            </ASpace>
        </ACol>
    </ARow>
    <ArticleEditCom v-model="editArticle.show" :article-id="editArticle.currentId"/>
</template>
<script setup>
import { apiGetUserStatistics } from "@/apis/user";
import ArticleEditCom from "@/components/article/ArticleEdit.vue";
import ArticleListCom from "@/components/article/ArticleList.vue";
import { useUserStore } from "@/stores/UserStore";
import { reactive, ref } from "vue";
// const IconFont = Icon.addFromIconFontCn({ src: 'https://at.alicdn.com/t/c/font_4443211_9ji0rjy60kw.js' });
const userStore = useUserStore()
const queryData = ref({
    tags: [],
    context: "",
    orderBy: "created_at",
    descOrder: true,
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
apiGetUserStatistics().then(({ data }) => {
    userStatistics.value = data
})
</script>