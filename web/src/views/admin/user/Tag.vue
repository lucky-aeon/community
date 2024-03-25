<template>
  <a-modal v-model:visible="visible" title="添加用户标签" @cancel="handleCancel" @before-ok="handleBeforeOk">
    <a-form :model="form">
      <a-form-item field="name" label="标签名">
        <a-input v-model="form.name" />
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
    <a-table row-key="id" :columns="columns" :data="tagData"
            :pagination="pagination" @page-change="getAdminListUserTags">
    <template  #optional="{ record, rowIndex }">
      <a-space>
        <a-button type="primary" @click="updateTag(rowIndex)">修改</a-button>
        <a-popconfirm popup-hover-stay @ok="delTag(rowIndex)" content="你确定要删除该标签?">
          <a-button type="primary">删除</a-button>
        </a-popconfirm>
      </a-space>
    </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { reactive, ref } from 'vue';
import {apiAdminDeleteUserTags, apiAdminListUserTags, apiAdminSaveUserTags} from "@/apis/user.js";
import {IconPlus} from "@arco-design/web-vue/es/icon/index.js";
import {apiDeleteCode} from "@/apis/code.js";

const visible = ref(false);
const form = reactive({
  id:null,
  name: '',
});


const handleBeforeOk = async (done) => {
  await apiAdminSaveUserTags(form)
  done()
  await getAdminListUserTags()
};
const handleClick = () => {
  visible.value = true;
  clearForm()
};
function clearForm(){
  form.id = null
  form.name = null
}
const handleCancel = () => {
  visible.value = false;
  clearForm()
}

function updateTag(id){
  const tag = tagData.value[id]
  visible.value = true;
  form.id = tag.id
  form.name = tag.name
}


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
    title: '创建时间',
    dataIndex: 'createdAt',
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
const tagData = ref([])
const getAdminListUserTags = (current)=>{
  pagination.current = current
  apiAdminListUserTags(current, pagination.defaultPageSize).then(({data})=>{
    tagData.value = data.list
    pagination.total = data.total
  })
}
getAdminListUserTags()

function delTag(id) {
  apiAdminDeleteUserTags(tagData.value[id].id).then(({ok})=>{
    if (ok) {
      tagData.value.splice(id,1)
    }
  })
}
</script>
