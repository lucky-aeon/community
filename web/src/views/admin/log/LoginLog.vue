<template>
  <a-space direction="vertical" size="large" fill>
    <a-form layout="inline">
      <a-form-item field="requestInfo" label="账户">
        <a-input v-model="searchData.account"/>
      </a-form-item>
      <a-form-item field="ip" label="ip">
        <a-input v-model="searchData.ip"/>
      </a-form-item>
      <a-form-item field="startTime" label="开始时间">
        <a-date-picker v-model="searchData.startTime" placeholder="Please select ..."/>
      </a-form-item>
      <a-form-item field="endTime" label="结束时间">
        <a-date-picker v-model="searchData.endTime" placeholder="Please select ..."/>
      </a-form-item>
      <a-form-item>
        <a-button @click="search()">搜索</a-button>
      </a-form-item>
    </a-form>


    <a-table :columns="columns" :data="loginData"
           :pagination="pagination" @page-change="getLoginLogList">
      <template #optional="{ record, rowIndex }">


      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import {apiLoginLogList} from '@/apis/log.js'
import {reactive, ref, h} from 'vue';

const searchData = ref({
  account:null,
  ip:"",
  startTime:"",
  endTime:""
})

const pagination = {pageSize: 15}

const columns = [
  {
    title: '账户',
    dataIndex: 'account',
  },
  {
    title: '状态',
    dataIndex: 'state',
  },
  {
    title: '浏览器',
    dataIndex: 'browser',
  },
  {
    title: '设备',
    dataIndex: 'equipment',
  },
  {
    title: 'ip',
    dataIndex: 'ip',
  },
  {
    title: '登录时间',
    dataIndex: 'createdAt',
  }
]

const loginData = ref([])
const getLoginLogList = (searchData)=>{
  apiLoginLogList(1,15,searchData).then(({data})=>{
    loginData.value = data.list
  })
}
getLoginLogList()

function search(){
  getLoginLogList(searchData.value)
}
</script>
