<template>
  <a-list :bordered="false" :data="commentData" :pagination-props="page" @page-change="getListData">
    <template #item="{ item, index }">
      <a-list-item>
        <a-list-item-meta :title="`${item.fromUserName}  ${item.toUserName == '' ? '' : '回复【@' + item.toUserName+'】'} 时间:${ item.createdAt }`"
          :description="item.content">
          <template #avatar>
            <a-avatar shape="square" :image-url="getFileUrl(item.fromUserAvatar)" @click="router.push({path: `/user/${item.FromUserId}`})">
            </a-avatar>
          </template>
        </a-list-item-meta>
        <template #actions>
          <a-dropdown-button @click="reply(item)">
            回复评论
            <template #icon>
              <icon-down />
            </template>
            <template #content>
              <a-doption @click="router.push(`/article/view/${item.articleId}`)"><icon-eye />查看原文</a-doption>
              <a-doption @click="del(index)"><icon-delete />删除评论</a-doption>
            </template>
          </a-dropdown-button>
        </template>
      </a-list-item>
    </template>
  </a-list>
  <a-modal title="评论回复" v-model:visible="replyEdit.show" :modal-style="{ minWidth: '800px', maxHeight: '90%' }"
    :body-style="{ padding: '0 0 0', height: '100%' }" ok-text="回复" :mask-closable="false" @ok="replyRq">
    <div style="height: 350px">
      <markdown-edit v-model="replyEdit.data" />
    </div>
  </a-modal>
</template>

<script setup>
import { apiReply, deleteComment, listAllCommentsByArticleId } from '@/apis/comment.js';
import { apiGetFile } from '@/apis/file';
import MarkdownEdit from "@/components/MarkdownEdit.vue";
import router from '@/router';
import { IconDelete, IconDown } from '@arco-design/web-vue/es/icon';
import { reactive, ref } from 'vue';


const replyEdit = reactive({
  show: false,
  data: "",
  item: {}
})

const getFileUrl = (fileKey) => {
    return apiGetFile(fileKey)

}
function replyRq() {
  const data1 = replyEdit.item
  const replyComment = {
    parentId: data1.parentId,
    rootId: data1.rootId,
    content: replyEdit.data,
    articleId: data1.articleId,
    toUserId: data1.FromUserId,
  }
  apiReply(replyComment)
}

function reply(item) {
  replyEdit.item = item
  replyEdit.show = true
}

function del(index) {
  deleteComment(commentData.value[index].id).then(({ ok }) => {
    if (ok) {
      commentData.value.splice(index, 1)
    }
  })


}
const page = ref({
  defaultPageSize: 10,
  total: 0,
  current: 1
})
const commentData = ref([])
function getListData(current) {
  page.value.current = current
  listAllCommentsByArticleId(0, current, page.value.defaultPageSize).then(({ data }) => {
    commentData.value = data.data
    page.value.total = data.count
  })
}
getListData()
</script>
