<template>
  <a-space direction="vertical" size="large" fill>
    <div>
      <span>OnlyCurrent: </span>
      <a-switch v-model="rowSelection.onlyCurrent" />
    </div>
    <a-table row-key="name" :columns="columns" :data="commentData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" >
      <template #optional="{ record }">
        <a-button @click="ttt(record.id)">删除</a-button>
      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { listAllCommentsByArticleId} from '@/apis/comment.js'
import { deleteComment} from '@/apis/comment.js'
import { reactive, ref } from 'vue';

function ttt(id) {
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
    title: 'fromUserName',
    dataIndex: 'fromUserName',
  },
  {
    title: 'toUserName',
    dataIndex: 'toUserName',
  },
  {
    title: 'content',
    dataIndex: 'content',
  },
  {
    title: 'articleTitle',
    dataIndex: 'articleTitle',
  },
  {
    title: 'Optional',
    slotName: 'optional'
  }
]

const commentData = ref([])
listAllCommentsByArticleId(0).then(({data})=>{
  commentData.value = data.data
})
</script>
