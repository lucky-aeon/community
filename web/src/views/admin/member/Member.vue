<template>
  <a-modal v-model:visible="visible" title="添加等级" @cancel="handleCancel" @before-ok="handleBeforeOk">
    <a-form :model="form">
      <a-form-item field="name" label="名称">
        <a-input v-model="form.name" />
      </a-form-item>
      <a-form-item field="desc" label="描述">
        <a-input v-model="form.desc" />
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
    <a-table row-key="name" :columns="columns" :data="memberData" :row-selection="rowSelection"
             v-model:selectedKeys="selectedKeys" :pagination="pagination" @page-change="getMemberList">
      <template #optional="{ record, rowIndex }">
        <a-space>
          <a-button type="primary" @click="updateComment(rowIndex)">修改</a-button>
          <a-popconfirm popup-hover-stay @ok="delComment(rowIndex)" content="你确定要删除该等级?">
            <a-button type="primary">删除</a-button>
          </a-popconfirm>
        </a-space>

      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { listAllMember,deleteMember} from '@/apis/member.js'
import { saveMember} from '@/apis/member.js'
import { reactive, ref } from 'vue';
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';

const visible = ref(false);
const form = reactive({
  id:null,
  name: '',
  desc: ''
});

const handleClick = () => {
  visible.value = true;
  clearForm()
};
const handleBeforeOk = (done) => {
  saveMember(form)
  done()
  getMemberList()
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

function updateComment(id){
  const comment = memberData.value[id]
  visible.value = true;
  form.id = comment.id
  form.name = comment.name
  form.desc = comment.desc
}

function delComment(id) {
  deleteMember(memberData.value[id].id).then(({ok})=>{
    if (ok) {
      memberData.value.splice(id,1)
    }
  })
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
    title: '名称',
    dataIndex: 'name',
  },
  {
    title: '描述',
    dataIndex: 'desc',
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

const memberData = ref([])
const getMemberList = (current)=>{
  pagination.current = current
  console.log(current)
  listAllMember(current,pagination.defaultPageSize).then(({data})=>{
    memberData.value = data.list
    pagination.total = data.total
  })
}
getMemberList()
</script>
