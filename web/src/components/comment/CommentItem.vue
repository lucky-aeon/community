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
            <template v-if="article.user.id == userStore.userInfo.id">
                <a-button v-if="comment.adoptionState" type="secondary" @click="adoption(comment.id)">
                取消采纳
            </a-button>
            <a-button v-else type="primary" @click="adoption(comment.id)">
                采纳
            </a-button>
            </template>
            
        </template>

        <template #avatar>
            <a-avatar :image-url="comment.fromUserAvatar" @click="router.push({path: `/user/${comment.FromUserId}`})"
            </a-avatar>
        </template>
        <CommentEdit :callback="getSubCommentPage" :parent-id="comment.id" :article-id="comment.articleId" :root-comment="comment.rootId" v-if="replyEdit.show"/>
        <CommentItem :article="article" v-for="item in subCommentData" :key="item.id" :comment="item" :callback="getSubCommentPage"/>
        <template v-if="comment.rootId==comment.id && comment.childCommentNumber>5">
            <a-button type="text" long v-if="!showPage" @click="getSubCommentPage()">展开评论</a-button>
            <APagination v-else-if="comment.childCommentNumber>pageData.pageSize" :page-size="pageData.pageSize" @change="getSubCommentPage" size="small" v-model:current="pageData.current" :total="pageData.total" />
        </template>
        
    </a-comment>
</template>

<script setup>
import { deleteComment as apiCommentDelete, apiGetArticleComment, apiAdoptionComment } from '@/apis/comment';
import { useUserStore } from '@/stores/UserStore';
import { IconMessage } from '@arco-design/web-vue/es/icon';
import 'cherry-markdown/dist/cherry-markdown.css';
import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core';
import { computed, reactive, ref } from 'vue';
import CommentEdit from './CommentEdit.vue';
import router from '@/router';
const props = defineProps({
    comment: {
        type: Object,
        default: () => { }
    },
    callback: {
        type: Function,
        default(){}
    },
    article: {
        type: Object,
        default: ()=>{}
    }
})
const replyEdit = ref({
    show: false,
    comment: {}
})
const showPage = ref(false)
const userStore = useUserStore()
const subComment = ref([])
const subCommentData = computed(()=> {
    if(subComment.value.length>0) {
        return subComment.value
    }
    return props.comment.childComments
})
const pageData = reactive({
    pageSize: 5,
    current: 1,
    total: props.comment.childCommentNumber
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
    apiGetArticleComment(props.comment.id, pageData.current, pageData.pageSize, false).then(({data, ok})=>{
        if(!ok || !data) return
        showPage.value = true
        subComment.value = data.data
        pageData.total = data.count
    })
}
function deleteComment() {
    apiCommentDelete(props.comment.id).then(({ok})=>{
        if(!ok) return
        props.callback()
    })
}
function adoption(id) {
    apiAdoptionComment(id).then(({ok})=>{
        if(!ok) return
        props.callback()
    })
}
</script>