<template>
    <a-comment style="margin-top: 0px;" :author="comment.fromUserName" :datetime="comment.updatedAt">
        <template #content >
            <div v-html="htmlContent">
            </div>
        </template>

        <template #actions>
            <IconMessage />
            <span class="action" @click="onLikeChange">
                {{ comment.childCommentNumber }}

            </span>
            <a-popconfirm v-if="comment.FromUserId == userStore.userInfo.id" content="你确定要删除该评论?" @ok="deleteComment()">
            <a-button size="small">
                <span class="action" key="deleteComment">
                    删除评论
                </span>
            </a-button>
            </a-popconfirm>
            <a-button v-else @click="replyEdit = {show: !replyEdit.show, comment: comment}" size="small">
                <span class="action" key="reply">
                    {{ replyEdit.show?"取消回复":"回复" }}
                </span>
            </a-button>
            
        </template>

        <template #avatar>
            <a-avatar :image-url="comment.fromUserAvatar">
            </a-avatar>
        </template>
        <CommentEdit :parent-id="comment.id" :article-id="comment.articleId" :root-comment="comment.rootId" v-if="replyEdit.show"/>
        <CommentItem :callback="getSubCommentPage" v-for="item in subCommentData" :key="item.id" :comment="item"/>
        <template v-if="comment.rootId==comment.id && comment.childCommentNumber>5">
            <a-button type="text" long v-if="!showPage" @click="getSubCommentPage()">展开评论</a-button>
            <APagination v-else-if="comment.childCommentNumber>5" @change="getSubCommentPage" size="small" :current="currentPage" :total="comment.childCommentNumber" />
        </template>
        
    </a-comment>
</template>

<script setup>
import { apiGetArticleComment, deleteComment as apiCommentDelete } from '@/apis/comment';
import { useUserStore } from '@/stores/UserStore';
import { IconMessage } from '@arco-design/web-vue/es/icon';
import 'cherry-markdown/dist/cherry-markdown.css';
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core';
import { computed, ref } from 'vue';
import CommentEdit from './CommentEdit.vue';
const props = defineProps({
    comment: {
        type: Object,
        default: () => { }
    },
    callback: {
        type: Object,
        default: ()=> {}
    }
})
const replyEdit = ref({
    show: false,
    comment: {}
})
const currentPage = ref(1)
const showPage = ref(false)
const userStore = useUserStore()
const subComment = ref([])
const subCommentData = computed(()=> {
    if(subComment.value.length>0) {
        return subComment.value
    }
    return props.comment.childComments
})
const commentMark = new CherryEngine();
const htmlContent = computed(() => {
    let content = props.comment.content
    if(props.comment.parentId) {
    content = `@${props.comment.toUserName||"未知用户"} 说: \n${content}`
    }
    return commentMark.makeHtml(content)
})
function getSubCommentPage() {
    apiGetArticleComment(props.comment.id, currentPage.value, 10, false).then(({data, ok})=>{
        if(!ok) return
        showPage.value = true
        subComment.value = data.data
    })
}
function deleteComment() {
    apiCommentDelete(props.comment.id).then(({ok})=>{
        if(!ok) return
        props.callback()
    })
}
</script>