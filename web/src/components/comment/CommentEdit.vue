<template>
    <a-comment :author="userStore.userInfo.name" datetime="发表评论">
        <template #content>
            <MarkdownEdit :show-nav="false" style="margin-top: 5px;height: 100%;" v-model="commentData" />
            <a-button long @click="pushComment()">发布评论</a-button>
        </template>

        <template #avatar>
            <a-avatar :image-url="userStore.userInfo.avatar">
            </a-avatar>
        </template>
    </a-comment>
</template>

<script setup>
import { apiPublishArticleComment } from '@/apis/comment';
import { useUserStore } from '@/stores/UserStore';
import { ref } from 'vue';
import MarkdownEdit from '../MarkdownEdit.vue';
const props = defineProps({
    articleId: {
        type: Number,
        default: 0,
    },
    parentId: {
        type: Number,
        default: 0,
    },
    rootComment: {
        type: Number,
        default: () => 0
    }
})
const userStore = useUserStore()
const commentData = ref("")
function pushComment() {
    console.log(props)
    apiPublishArticleComment(props.articleId, props.parentId, commentData.value, props.rootComment)
}
</script>