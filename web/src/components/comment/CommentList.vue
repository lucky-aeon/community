<template>
    <a-card :bordered="false">
        <comment-item v-for="comment in articleCommentList" :key="comment.id"/>
    </a-card>
</template>
<script setup>
import { apiGetArticleComment } from '@/apis/comment';
import { onMounted, ref } from 'vue';
import CommentItem from './CommentItem.vue';
const props = defineProps({
    articleId: {
        type: Number,
        default: 0,
    }
})
const articleCommentList = ref([])
function getRootComment() {
    apiGetArticleComment(props.articleId).then(({data, ok})=>{
    if(!ok) {
        return
    }
    articleCommentList.value.push(...data)
})
}
onMounted(()=>{
    getRootComment()
})



</script>