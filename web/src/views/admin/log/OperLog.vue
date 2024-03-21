<template>
  <a-space direction="vertical" size="large" >

    <a-form layout="inline">
      <a-form-item field="requestType" label="请求方法" >
        <a-select allow-clear v-model="searchData.requestMethod">
          <a-option value="GET">GET</a-option>
          <a-option value="POST">POST</a-option>
          <a-option value="PUT">PUT</a-option>
          <a-option value="DELETE">DELETE</a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="requestInfo" label="请求信息">
        <a-input v-model="searchData.requestInfo"/>
      </a-form-item>
      <a-form-item field="userId" label="访问者">
        <a-input v-model="searchData.userName"/>
      </a-form-item>
      <a-form-item field="ip" label="ip">
        <a-input v-model="searchData.ip"/>
      </a-form-item>
      <a-form-item field="startTime" label="开始时间">
        <a-date-picker v-model="searchData.startTime" placeholder="Please select ..."/>
      </a-form-item>
      <a-form-item field="endTime" label="范围时间">
        <a-date-picker v-model="searchData.endTime" placeholder="Please select ..."/>
      </a-form-item>
      <a-form-item>
        <a-button @click="search()">搜索</a-button>
      </a-form-item>
    </a-form>

    <a-table :columns="columns" :data="operData"
           :pagination="pagination" :expandable="expandable" @page-change="getOperLogList">
      <template #optional="{ record, rowIndex }">

      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import {Descriptions} from '@arco-design/web-vue'
import {apiOperLogList} from '@/apis/log.js'
import {reactive, ref, h} from 'vue';

const searchData = ref({
  requestMethod:"",
  requestInfo:"",
  userName:null,
  ip:"",
  startTime:"",
  endTime:""
})


const expandable = reactive({
  title: 'Expand',
  width: 80,
});
const columns = [
  {
    title: '请求方法',
    dataIndex: 'requestMethod',
  },
  {
    title: '请求信息',
    dataIndex: 'requestInfo',
  },
  {
    title: '访问者',
    dataIndex: 'userName',
  },
  {
    title: 'ip',
    dataIndex: 'ip',
  },
  {
    title: '执行时间',
    dataIndex: 'execAt',
  },
  {
    title: '操作时间',
    dataIndex: 'createdAt',
  }
]

const pagination = reactive({
  total: 0,
  current: 1,
  defaultPageSize: 10
})
const operData = ref([])

const getOperLogList = async (current,searchData)=>{
   await apiOperLogList(current,pagination.defaultPageSize,searchData).then(({data})=>{
    let temp = data.list.map(e=>{
      return {
        ...e,
        key: e.id,
        expand: h(
            Descriptions,
            {
              data: createLogDetail(e),
              layout: "inline-vertical"
            })
      }
    })
     pagination.current = current
     pagination.total = data.total
     operData.value = temp

  })
}
getOperLogList()
function  createLogDetail(log) {
  return [{
    label:"请求参数",
    value:log.requestBody
  }
    ,{
      label:"响应参数",
      value:log.responseData
    }
  ]
}

function search(){
  getOperLogList(pagination.current,searchData.value)
}
</script>
