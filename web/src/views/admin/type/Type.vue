<template>
  <a-modal v-model:visible="visible" title="添加等级" @cancel="handleCancel" @before-ok="handleBeforeOk">
    <a-form :model="form">
      <a-form-item field="section" label="一级分类">
        <a-select v-model="form.parentId" allow-clear>
          <a-option v-for="item of parentData" :value="item.id" :label="item.title" />
        </a-select>
      </a-form-item>
      <a-form-item field="name" label="名称">
        <a-input v-model="form.title" />
      </a-form-item>
      <a-form-item field="desc" label="描述">
        <a-input v-model="form.desc" />
      </a-form-item>
      <a-form-item field="state" label="状态">
        <a-switch v-model="form.state" />
      </a-form-item>
      <a-form-item field="sort" label="排序">
        <a-input-number v-model="form.sort" />
      </a-form-item>
      <a-form-item field="flagName" label="flagName">
        <a-input v-model="form.flagName" />
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
            新建
          </a-button>
        </a-space>
      </a-col>
      <a-col
          :span="12"
          style="display: flex; align-items: center; justify-content: end"
      >
      </a-col>
    </a-row>
    <a-table row-key="id" :columns="columns" :data="typeData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" >
      <template #optional="{ record, rowIndex }">
        <a-space>
          <a-button type="primary" @click="updateType(record)">修改</a-button>
          <a-button type="primary" @click="delType(record.id)">删除</a-button>
        </a-space>

      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import {apiDeleteType, apilistAllType, apiSaveType,apiListParentAllType} from '@/apis/type.js'
import { saveMember} from '@/apis/member.js'
import {onMounted, reactive, ref} from 'vue';
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';


const visible = ref(false);
const form = reactive({
  id:null,
  parentId:null,
  title: '',
  desc: '',
  state: true,
  sort: 0,
  flagName: '',
});

const handleClick = () => {
  visible.value = true;
  clearForm()
};
const handleBeforeOk = async (done) => {
  await apiSaveType(form)
  done()
  await getTypeList()
};

function clearForm(){
  form.id = null
  form.parentId = null
  form.title = null
  form.desc = null
  form.state = null
  form.sort = null
  form.flagName = null
}
const handleCancel = () => {

  visible.value = false;
  clearForm()
}

function updateType(record){
  const type = record
  visible.value = true;
  form.id = type.id
  if (type.parentId == 0){
    form.parentId = null
  }else{
    form.parentId = type.parentId
  }
  form.title = type.title
  form.desc = type.desc
  form.state = type.state
  form.sort = type.sort
  form.flagName = type.flagName
}



function delType(id) {
  apiDeleteType(id).then(({ok})=>{
    getTypeList()
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
    title: '名称',
    dataIndex: 'title',
  },
  {
    title: '描述',
    dataIndex: 'desc',
  },
  {
    title: '状态',
    dataIndex: 'state',
  },
  {
    title: '排序',
    dataIndex: 'sort',
  },
  {
    title: 'flag_name',
    dataIndex: 'flagName',
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

const typeData = ref([])
const getTypeList = ()=>{
  apilistAllType().then(({data})=>{
    typeData.value = data.list
  })
}
getTypeList()

const parentData = ref([])
const getParentTypeList = ()=>{
  apiListParentAllType().then(({data})=>{
    parentData.value=data
  })
}
getParentTypeList()
</script>
