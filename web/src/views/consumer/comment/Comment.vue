<template>
  <a-space direction="vertical" size="large" fill>

    <a-table row-key="name" :columns="columns" :data="commentData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" >
      <template #optional="{ record }">
        <a-button @click="delComment(record.id)">删除</a-button>
      </template>
    </a-table>
  </a-space>
</template>

import { listAllCommentsByArticleId } from '@/apis/comment.js';
import { reactive, ref } from 'vue';
import { deleteComment} from '@/apis/comment.js'
import { reactive, ref } from 'vue';

function delComment(id) {
  deleteComment(id).then(({msg})=>{
    console.log(msg)
  })
}

const selectedKeys = ref(['Jane Doe', 'Alisa Ross']);

const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false,
});
const pagination = {pageSize: 15}

const columns = [
  {
    title:"id",
    dataIndex: "id"
  },
  {
    title: '评论者',
    dataIndex: 'fromUserName',
  },
  {
    title: '回复人',
    dataIndex: 'toUserName',
  },
  {
    title: '评论内容',
    dataIndex: 'content',
  },
  {
    title: '文章标题',
    dataIndex: 'articleTitle',
  },
  {
    title: '操作',
    slotName: 'optional'
  }
]

const commentData = ref([])
listAllCommentsByArticleId(0).then(({data})=>{
  commentData.value = data.data
})
</script>
