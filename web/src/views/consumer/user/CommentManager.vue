<template>
  <a-list
      class="list-demo-action-layout"
      :bordered="false"
      :data="commentData"
      :pagination-props="page"
  >
    <template #item="{ item ,index}">
      <a-list-item class="list-demo-item" action-layout="vertical">
        <template #actions>
          <span>{{item.createdAt}}</span>
          <a-button @click="reply(item)" type="text">
            <icon-message/>回复
          </a-button>
          <a-button @click="del(index)" type="text">
            <icon-delete />
          </a-button>
        </template>

        <a-list-item-meta
            :title="`${item.fromUserName}  ${item.toUserName==''? '' : '回复@'+item.toUserName} `"
            :description="item.content"
        >
          <template #avatar>
            <a-avatar shape="square">
              <img alt="avatar" :src="item.fromUserAvatar" />
            </a-avatar>
          </template>
        </a-list-item-meta>

      </a-list-item>
    </template>
  </a-list>
  <a-modal title="评论回复" v-model:visible="replyEdit.show" :modal-style="{ minWidth: '800px', maxHeight: '90%' }" :body-style="{ padding: '0 0 0', height: '100%'}"
           ok-text="回复" mask-closable="false" @ok="replyRq">
    <div style="height: 350px">
      <markdown-edit v-model="replyEdit.data" />
    </div>
  </a-modal>
</template>

<script setup>
import { IconMessage,IconDelete } from '@arco-design/web-vue/es/icon';
import { listAllCommentsByArticleId,deleteComment,apiReply } from '@/apis/comment.js';
import { reactive, ref } from 'vue';
import MarkdownEdit from "@/components/MarkdownEdit.vue";


const replyEdit = reactive({
  show: false,
  data: "",
  item:{}
})


function replyRq(){
  const data1 = replyEdit.item
  const replyComment = {
    parentId:data1.parentId,
    rootId:data1.rootId,
    content:replyEdit.data,
    articleId:data1.articleId,
    toUserId:data1.FromUserId,
  }
  apiReply(replyComment)
}

function reply(item){
  replyEdit.item = item
  replyEdit.show = true
}

function del(index){
  deleteComment(commentData.value[index].id).then(({ok})=>{
    if (ok) {
      commentData.value.splice(index,1)
    }
  })


}

const commentData = ref([])
const count = ref([])
listAllCommentsByArticleId(0).then(({data})=>{
  commentData.value = data.data
  count.value = data.count
})
const page = reactive({
  defaultPageSize: 10,
  total: count.value
})
</script>
