<template>
    <a-card :bordered="false">
        <comment-item :callback="getRootComment" v-for="comment in articleCommentList" :comment="comment" :key="comment.id"/>
        <a-divider/>
        <a-pagination @change="getRootComment" :total="paginationData.total" v-model:page-size="paginationData.pageSize" v-model:current="paginationData.current" show-page-size/>
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
const paginationData = ref({
    current: 1,
    total: 0,
    pageSize: 5
    
})
function getRootComment() {
    apiGetArticleComment(props.articleId, paginationData.value.current, paginationData.value.pageSize).then(({data, ok})=>{
    if(!ok || !data.data) {
        return
    }
    paginationData.value.total = data.count
    articleCommentList.value = data.data
})
}
onMounted(()=>{
    getRootComment()
})



</script>