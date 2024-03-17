<template>
  <a-table :columns="columns" :bordered="false" :data="fileData" @page-change="getFileList" :pagination-props="page">


    <template #fileKey="{ record }">
      <div style="width: 200px; height: 100px;">
      <template v-if="record.mimeType.includes('video')">
        <video :src="apiGetFile(record.fileKey)" style="max-width: 100%; max-height: 100%;"></video>
      </template>
      <template v-else>
        <a-image :src="apiGetFile(record.fileKey)"  width="100%" height="100%">123</a-image>
      </template>
      </div>
    </template>

  </a-table>

</template>

<script setup>
import {ref, render} from "vue";
import {apiAdminFile, apiGetFile} from "@/apis/file.js"

const columns = [{
  title: '资源',
  slotName: 'fileKey',
}, {
  title: '文件大小',
  dataIndex: 'sizeName',
}, {
  title: '文件类型',
  dataIndex: 'mimeType',
}, {
  title: '发布者',
  dataIndex: 'userName',
}, {
  title: '上传时间',
  dataIndex: 'createdAt'
}];


const fileData = ref([])
const page = ref({
  defaultPageSize: 15,
  total: 0,
  current:1
})
function getFileList(current){
  page.value.current = current
  apiAdminFile(page.value.current = current,page.value.defaultPageSize).then(({data})=>{
    fileData.value = data.data
    page.value.total = data.count

  })
}
getFileList()
</script>
