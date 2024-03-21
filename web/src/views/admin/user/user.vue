<template>
  <a-modal v-model:visible="visible" title="添加等级" @cancel="handleCancel" @before-ok="handleBeforeOk">
    <a-form :model="form">
      <a-form-item field="name" label="名称">
        <a-input v-model="form.name" />
      </a-form-item>
      <a-form-item field="account" label="账户">
        <a-input v-model="form.account" />
      </a-form-item>
      <a-form-item field="desc" label="描述">
        <a-input v-model="form.desc" />
      </a-form-item>
    </a-form>
  </a-modal>
  <a-space direction="vertical" size="large" fill>

    <a-table row-key="id" :columns="columns" :data="userData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" @page-change="getAdminListUsers">
      <template #avatar="{ record, rowIndex }">
        <div style="width: 200px; height: 100px;">
        <a-image :src="apiGetFile(record.avatar)"  width="100%" height="100%">123</a-image>
        </div>

      </template>
    <template  #optional="{ record, rowIndex }">
      <a-space>
        <a-button type="primary" @click="updateUser(rowIndex)">修改</a-button>
      </a-space>
    </template>
    </a-table>
  </a-space>
</template>

<script setup>
import {apiAdminListUsers, apiAdminUpdateUsers} from '@/apis/user.js'
import { saveMember} from '@/apis/member.js'
import { reactive, ref } from 'vue';
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';
import {apiGetFile} from "@/apis/file.js";

const visible = ref(false);
const form = reactive({
  id:null,
  name: '',
  account: '',
  desc: ''
});


const handleBeforeOk = async (done) => {
  await apiAdminUpdateUsers(form)
  done()
  await getAdminListUsers()
};

function clearForm(){
  form.id = null
  form.name = null
  form.desc = null
}
const handleCancel = () => {
  visible.value = false;
  clearForm()
}

function updateUser(id){
  const user = userData.value[id]
  visible.value = true;
  form.id = user.id
  form.name = user.name
  form.account = user.account
  form.desc = user.desc
  console.log(userData.value[id])
}

const selectedKeys = ref(['Jane Doe', 'Alisa Ross']);

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
    title: '昵称',
    dataIndex: 'name',
  },
  {
    title: '账户',
    dataIndex: 'account',
  },
  {
    title: '邀请码',
    dataIndex: 'inviteCode',
  },
  {
    title: '描述',
    dataIndex: 'desc',
  },
  {
    title: '头像',
    slotName: 'avatar',
  },
  {
    title: '操作',
    slotName: 'optional'
  }
]
const pagination = reactive({
  total: 0,
  current: 1,
  defaultPageSize: 10
})
const userData = ref([])
const getAdminListUsers = (current)=>{
  pagination.current = current
  apiAdminListUsers(current, pagination.defaultPageSize).then(({data})=>{
    userData.value = data.list
    pagination.total = data.total
  })
}
getAdminListUsers()
</script>
