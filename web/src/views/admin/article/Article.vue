<template>

  <a-space direction="vertical" size="large" fill>
    <a-table row-key="id" :columns="columns" :data="articleData"
             :pagination="pagination" >
      <template #optional="{ record, rowIndex }">
        <a-space>
          <a-button type="primary" @click="delArticle(rowIndex)">删除</a-button>
        </a-space>

      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { } from '@/apis/code.js'
import { reactive, ref } from 'vue';
import {apiAdminDeleteArticles, apiAdminListArticles} from "@/apis/article.js";

const columns = [
  {
    title:"id",
    dataIndex: "id"
  },
  {
    title: '标题',
    dataIndex: 'title',
  },
  {
    title: '分类',
    dataIndex: 'type.title',
  },
  {
    title: '状态',
    dataIndex: 'stateName',
  },
  {
    title: '发布者',
    dataIndex: 'user.name',
  },
  {
    title: '发布时间',
    dataIndex: 'createdAt',
  },
  {
    title: '操作',
    slotName: 'optional'
  }
]
const pagination = ref({
  total: 0,
  current: 1,
  defaultPageSize: 10
})

const articleData = ref([])
const getArticles = (current)=>{
  pagination.value.current = current
  apiAdminListArticles(current,pagination.value.defaultPageSize).then(({data})=>{
    articleData.value = data.list
    pagination.value.total = data.total
  })
}

getArticles()
function delArticle(id) {
  apiAdminDeleteArticles(articleData.value[id].id).then(({ok})=>{
    if (ok) {
      articleData.value.splice(id,1)
    }
  })
}

</script>
