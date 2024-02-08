<template>
  <a-card direction="vertical" size="large">
    <a-space :size="54">
      <a-avatar
          :trigger-icon-style="{ color: '#3491FA' }"
          :auto-fix-font-size="false"
          @click="toast"
          :style="{ backgroundColor: '#168CFF' }"
      >
        A
        <template #trigger-icon>
          <IconCamera />
        </template>
      </a-avatar>
      <a-descriptions :data="userInfo" align="right" />
    </a-space>
  </a-card>

  <a-tabs default-active-key="1">
    <a-tab-pane key="1" title="Tab 1">
      <a-space direction="horizontal" size="large" :style="{width: '600px'}" >
        <a-form :model="form" :style="{ width: '600px' }" @submit="editUserInfo">
          <a-form-item field="name" label="昵称" :rules="[{required:true,message:'name is required'}]">
            <a-input
                v-model="form.name"
            />
          </a-form-item>
          <a-form-item field="desc" label="描述" >
            <a-textarea placeholder="Please enter something" :max-length="{length:200,errorOnly:true}" allow-clear
                        v-model="form.desc"
                        show-word-limit />
          </a-form-item>
          <a-form-item>
            <a-space>
              <a-button html-type="submit">Submit</a-button>
              <a-button @click="$refs.formRef.resetFields()">Reset</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-space>
    </a-tab-pane>
    <a-tab-pane key="2" title="Tab 2">
      <a-space size="large" :style="{width: '600px'}" >
        <a-form :model="form1"  :style="{ width: '600px' }" @submit="editPswd">
          <a-form-item field="oldPassword" label="旧密码" :rules="[{required:true,message:'请输入旧密码'}]">
            <a-input
                v-model="form1.oldPassword"
            />
          </a-form-item>
          <a-form-item field="newPassword" label="新密码" :rules="[{required:true,message:'请输入新密码'}]">
            <a-input
                v-model="form1.newPassword"
            />
          </a-form-item>
          <a-form-item field="confirmPassword" label="确认密码" :rules="[{required:true,message:'请输入二次密码'}]">
            <a-input
                v-model="form1.confirmPassword"
            />
          </a-form-item>
          <a-form-item>
            <a-space>
              <a-button html-type="submit">Submit</a-button>
              <a-button @click="$refs.formRef.resetFields()">Reset</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-space>
    </a-tab-pane>

  </a-tabs>

</template>

<script setup>
import {getUserInfo, saveUserInfo} from '@/apis/user'
import { reactive, ref } from 'vue';
import { IconCamera, IconEdit, IconUser } from '@arco-design/web-vue/es/icon';


const toast = function () {
  this.$message.info('Uploading...');
}
const userInfo = ref([])
const form = ref({
  name: ""
})
const form1 = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: ""
})
getUserInfo().then(({data})=>{
  userInfo.value = [{
    label: '账号',
    value: data.account,
  }, {
    label: '注册时间',
    value: data.createdAt,
  }];
  form.value = data
})
const editUserInfo = ({values, errors}) => {
  saveUserInfo("info",values).then(({data})=>{

  })
}
const editPswd = ({values, errors}) => {
  console.log(values)

  saveUserInfo("pass",values).then(({data})=>{
  })
}
</script>

