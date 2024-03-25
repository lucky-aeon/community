<template>
  <a-modal v-model:visible="visible" title="生成邀请码" @cancel="handleCancel" @ok="save">
    <a-form :model="generateCode">
      <a-form-item field="number" label="生成数量">
        <a-input-number v-model="generateCode.number"/>
      </a-form-item>
      <a-form-item  field="member" label="绑定等级">
        <a-select :style="{width:'320px'}" v-model="generateCode.memberId" :options="memberData" :field-names="fieldNames">
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
            生成
          </a-button>
        </a-space>
      </a-col>
      <a-col
          :span="12"
          style="display: flex; align-items: center; justify-content: end"
      >
      </a-col>
    </a-row>
    <a-table row-key="name" :columns="columns" :data="codeData" :row-selection="rowSelection"
            :pagination="pagination" @page-change="getCodeList">
      <template #optional="{ record, rowIndex }">
        <a-space>
          <a-popconfirm popup-hover-stay @ok="delCode(rowIndex)" content="你确定要删除该邀请码?">
            <a-button type="primary">删除</a-button>
          </a-popconfirm>
        </a-space>

      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { apiCodeList,apiGenerateCode,apiDeleteCode} from '@/apis/code.js'
import {listAllMember, saveMember} from '@/apis/member.js'
import { reactive, ref } from 'vue';
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';


const fieldNames = {value: 'id', label: 'name'}

const memberData = ref([])

const visible = ref(false);
const generateCode = reactive({
  number: 1,
  memberId:1,
});

const handleClick = () => {
  visible.value = true;
  clearForm()
};
const save = async (done) => {
  await apiGenerateCode(generateCode)
  getCodeList()
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
    title: '邀请码',
    dataIndex: 'code',
  },
  {
    title: '等级',
    dataIndex: 'memberName',
  },

  {
    title: '状态',
    dataIndex: 'state',
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

const codeData = ref([])
const getCodeList = (current)=>{
  pagination.value.current = current
  apiCodeList(current,pagination.value.defaultPageSize).then(({data})=>{
    codeData.value = data.list
    pagination.value.total = data.total
  })
}



getCodeList()
function delCode(id) {
  apiDeleteCode(codeData.value[id].code).then(({ok})=>{
    if (ok) {
      codeData.value.splice(id,1)
    }
  })
}

function clearForm(){
  generateCode.id = null
  generateCode.name = null
  generateCode.desc = null
}

listMember()
function listMember(){
  listAllMember(1,10).then(({data})=>{
    memberData.value = data.list
    generateCode.memberId = data.list[0].id
  })
}
</script>
