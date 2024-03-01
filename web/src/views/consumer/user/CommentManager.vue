<template>
  <a-space direction="vertical" size="large" fill>
    <div>
      <span>OnlyCurrent: </span>
      <a-switch v-model="rowSelection.onlyCurrent" />
    </div>
    <a-table row-key="name" :columns="columns" :data="commentData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" />
  </a-space>
</template>

<script setup>
import { listAllCommentsByArticleId } from '@/apis/comment.js';
import { reactive, ref } from 'vue';

const selectedKeys = ref(['Jane Doe', 'Alisa Ross']);

const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false,
});
const pagination = {pageSize: 5}

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
]
// 你data是什么数据？对的上的数据 哈哈哈哈哈哈哈？咋这么多data，哪个data ？？？？ 行不行，
const commentData = ref([])
listAllCommentsByArticleId(0).then(({data})=>{
  // 你是会写变量名 ？
  // data.value = data.data 都是dat你要给谁？ 我一开始写的是body我不管 ？？？？？我tm
  commentData.value = data.data
})
</script>
