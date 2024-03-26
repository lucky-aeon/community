<template>
  <a-modal v-model:visible="visible" title="创建事件消息模板" @cancel="handleCancel" @ok="save">
    <a-form :model="form">
      <a-form-item field="number" label="内容">
        <a-textarea  v-model="form.content"/>
      </a-form-item>
      <a-form-item  field="member" label="事件">
        <a-select :style="{width:'320px'}" v-model="form.eventId" :options="eventData" :field-names="fieldNames">
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
  <a-space direction="vertical" size="large" fill>
    <a-row style="margin-bottom: 16px">
      <a-col :span="12">
        <a-space>
          <a-button type="primary" @click="handleClick">
            <template #icon>
              <icon-plus />
            </template>
            创建
          </a-button>
        </a-space>
      </a-col>
      <a-col
          :span="12"
          style="display: flex; align-items: center; justify-content: end"
      >
      </a-col>
    </a-row>
    <a-table row-key="name" :columns="columns" :data="MsgTempData" :row-selection="rowSelection"
             :pagination="pagination" @page-change="getMsgTempList">
      <template #optional="{ record, rowIndex }">
        <a-space>
          <a-button type="primary" @click="updateMsgTemp(rowIndex)">修改</a-button>
          <a-popconfirm popup-hover-stay @ok="delMsgTemp(rowIndex)" content="你确定要删除该事件吗?">
            <a-button type="primary">删除</a-button>
          </a-popconfirm>
        </a-space>

      </template>
    </a-table>
  </a-space>
</template>

<script setup>

import { reactive, ref } from 'vue';
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';
import {
  apiMessageTemplateEventList,
  apiMessageTemplateList,
  apiMessageTemplateListDelete,
  apiMessageTemplateSave
} from "@/apis/message.js";
import {listAllMember} from "@/apis/member.js";


const fieldNames = {value: 'id', label: 'msg'}

const eventData = ref([])

const visible = ref(false);


const handleClick = () => {
  visible.value = true;
  clearForm()
};
const save = async (done) => {
  await apiMessageTemplateSave(form)
  getMsgTempList()
};

const handleCancel = () => {
  visible.value = false;
  clearForm()
}


const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false,
});


const columns = [
  {
    title:"id",
    dataIndex: "id"
  },
  {
    title: '内容',
    dataIndex: 'content',
  },
  {
    title: '事件',
    dataIndex: 'eventName',
  },
  {
    title: '创建时间',
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

const MsgTempData = ref([])
const getMsgTempList = (current)=>{
  pagination.value.current = current
  apiMessageTemplateList(current,pagination.value.defaultPageSize).then(({data})=>{
    MsgTempData.value = data.list
    pagination.value.total = data.total
  })
}
const form = reactive({
  id:null,
  content: '',
  eventId: null
});
function updateMsgTemp(id){
  const temp = MsgTempData.value[id]
  visible.value = true;
  form.id = temp.id
  form.content = temp.content
  form.eventId = temp.eventId
}


getMsgTempList()
function delMsgTemp(id) {
  apiMessageTemplateListDelete(MsgTempData.value[id].id).then(({ok})=>{
    if (ok) {
      MsgTempData.value.splice(id,1)
    }
  })
}

function clearForm(){
  form.id = null
  form.content = null
  form.eventId = null
}
listEvent()
function listEvent(){
  apiMessageTemplateEventList().then(({data})=>{
    eventData.value = data
    eventData.value.splice(0,1)
  })
}
</script>
