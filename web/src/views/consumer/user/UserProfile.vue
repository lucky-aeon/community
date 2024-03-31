<template>
  <a-card style="margin-top: 5px;" :bordered="false">
  <a-tabs default-active-key="1">
    <a-tab-pane key="1" title="个人信息">
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
          <a-form-item field="desc" label="订阅站点消息" >
              <a-switch v-model="form.subscribe" :checked-value="2" :unchecked-value="1" type="round">
                <template #checked>
                  订阅
                </template>
                <template #unchecked>
                  未订阅
                </template>
              </a-switch>
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
    <a-tab-pane key="2" title="账号管理">
      <a-space size="large" :style="{width: '600px'}" >
        <a-form :model="form1"  :style="{ width: '600px' }" @submit="editPswd">
          <a-form-item field="oldPassword" label="旧密码" :rules="[{required:true,message:'请输入旧密码'}]">
            <a-input
                v-model="form1.oldPassword" type="password"
            />
          </a-form-item>
          <a-form-item field="newPassword" label="新密码" :rules="[{required:true,message:'请输入新密码'}]">
            <a-input
                v-model="form1.newPassword" type="password"
            />
          </a-form-item>
          <a-form-item field="confirmPassword" label="确认密码" :rules="[{required:true,message:'请输入二次密码'}]">
            <a-input
                v-model="form1.confirmPassword" type="password"
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
  </a-card>


</template>

<script setup>
import { getUserInfo, saveUserInfo } from '@/apis/user';
import { reactive, ref } from 'vue';


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
const editUserInfo = ({values}) => {
  saveUserInfo("info",values).then(()=>{

  })
}
const editPswd = ({values}) => {
  console.log(values)

  saveUserInfo("pass",values).then(()=>{
  })
}
</script>

