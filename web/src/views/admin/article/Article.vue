<template>
  <a-modal v-model:visible="visible" title="修改文章" @cancel="handleCancel" @before-ok="handleBeforeOk">
    <a-form :model="form">
      <a-form-item field="name" label="标题">
        <a-input v-model="form.title" />
      </a-form-item>
      <a-form-item  field="name" label="状态">
        <a-select :style="{width:'320px'}" v-model="form.state" :options="articleStates" :field-names="fieldNames">
        </a-select>
      </a-form-item>
      <a-form-item field="desc" label="置顶排序">
        <a-input-number v-model="form.topNumber" />
      </a-form-item>
    </a-form>
  </a-modal>
  <a-space direction="vertical" size="large" fill>
    <a-table row-key="id" :columns="columns" :data="articleData"
             :pagination="pagination" >
      <template #optional="{ record, rowIndex }">
        <a-space>
        <a-button type="primary" @click="updateState(rowIndex)">修改</a-button>
        <a-popconfirm popup-hover-stay @ok="delArticle(rowIndex)" content="你确定要删除该文章?">
          <a-button type="primary">删除</a-button>
        </a-popconfirm>
        </a-space>
      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { } from '@/apis/code.js'
import { reactive, ref } from 'vue';
import {
  apiAdminDeleteArticles,
  apiAdminListArticles,
  apiAdminListArticleStates,
  apiAdminUpdateArticleState
} from "@/apis/article.js";

const articleStates = ref([])

const fieldNames = {value: 'id', label: 'name'}

const visible = ref(false);
const form = reactive({
  id:null,
  name: '',
  state: 0,
  topNumber: 0
});


const handleBeforeOk = (done) => {
  apiAdminUpdateArticleState(form)
  done()
  getArticles()
};

const handleCancel = () => {
  visible.value = false;
}

function updateState(id){
  const article = articleData.value[id]
  visible.value = true;
  form.id = article.id
  form.title = article.title
  form.state = article.state
  form.topNumber = article.topNumber
}

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
    title: '置顶排序',
    dataIndex: 'topNumber',
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

listStates()
function listStates(){
  apiAdminListArticleStates().then(({data})=>{
    articleStates.value = data
  })
}

</script>
