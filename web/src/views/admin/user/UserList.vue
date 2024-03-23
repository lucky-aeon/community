<template>
  <a-modal :ok-loading="loading" v-model:visible="visible" title="修改用户" @cancel="handleCancel" @before-ok="handleBeforeOk">
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
      <a-form-item field="tags" label="用户标签">
        <a-select v-model="form.tags" :loading="loading" placeholder="请选择用户标签" multiple @search="handleTagSearch"
          :filter-option="false" :show-extra-options="false" :field-names="tagsFieldNames" :options="systemTags" />
      </a-form-item>
    </a-form>
  </a-modal>
  <a-space direction="vertical" size="large" fill>

    <a-table row-key="id" :columns="columns" :data="userData" :row-selection="rowSelection"
      v-model:selectedKeys="selectedKeys" :pagination="pagination" @page-change="getAdminListUsers">
      <template #avatar="{ record }">
        <div style="width: 200px; height: 100px;">
          <a-image :src="apiGetFile(record.avatar)" width="100%" height="100%">123</a-image>
        </div>

      </template>
      <template #optional="{ rowIndex }">
        <a-space>
          <a-button :loading="loading" type="primary" @click="updateUser(rowIndex)">修改</a-button>
        </a-space>
      </template>
    </a-table>
  </a-space>
</template>

<script setup>
import { apiGetFile } from "@/apis/file.js";
import { apiAdminListUserTags, apiAdminListUsers, apiAdminUpdateUsers, apiGetUserTags } from "@/apis/user.js";
import { reactive, ref } from 'vue';

const loading = ref(false)

const visible = ref(false);
const form = reactive({
  id: null,
  name: '',
  account: '',
  desc: '',
  tags: [],
});


const handleBeforeOk = async (done) => {
  await apiAdminUpdateUsers(form)
  done()
  getAdminListUsers()
};

function clearForm() {
  form.id = null
  form.name = null
  form.desc = null
}
const handleCancel = () => {
  visible.value = false;
  clearForm()
}

async function updateUser(id) {
  loading.value = true
  const user = userData.value[id]
  visible.value = true;
  form.id = user.id
  form.name = user.name
  form.account = user.account
  form.desc = user.desc
  form.tags = []
  let res = await apiGetUserTags(user.id)
  if(res.ok && res.data) {
    form.tags = res.data.map(tagItem=> tagItem.id)
  }
  loading.value = false
}

const selectedKeys = ref(['Jane Doe', 'Alisa Ross']);

const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false,
});


const columns = [
  {
    title: "id",
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
const getAdminListUsers = (current) => {
  pagination.current = current
  apiAdminListUsers(current, pagination.defaultPageSize).then(({ data }) => {
    userData.value = data.list
    pagination.total = data.total
  })
}
getAdminListUsers()
const systemTags = ref([])
const tagsFieldNames = ref({
  label: "name",
  value: "id"
})
function handleTagSearch(value) {
  apiAdminListUserTags(value, 10).then(({ data, ok }) => {
    if (!ok) return
    systemTags.value = data.list
  }).finally((()=> loading.value = false));
}

handleTagSearch("")

</script>
