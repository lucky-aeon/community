<template>
    <a-menu style="font-size: 16px;" class="line" mode="horizontal" :default-selected-keys="['1']">
        <a-menu-item key="0" :style="{ padding: 0, marginRight: '38px' }" disabled>
            <div :style="{
                width: '80px',
                height: '30px',
                borderRadius: '2px',
                background: 'var(--color-fill-3)',
                cursor: 'text',
            }" />
        </a-menu-item>
        <a-menu-item key="1">Home</a-menu-item>
        <a-menu-item key="2">services</a-menu-item>
        <a-menu-item key="3">about</a-menu-item>
        <a-menu-item key="4">other</a-menu-item>
    </a-menu>
    <a-divider />
    <div style="height: 100%; width: 100%;margin-top: 100px;">
        <a-card hoverable style="margin: auto; width: 550px;" title="身份验证">
            <template #extra>
                <a-link>帮助</a-link>
            </template>
            <a-form :model="authForm" :style="{ width: '100%' }" @submit="handleSubmit">
                <a-form-item feedback field="account" tooltip="联系站长获取邀请码" label="Email"
                    :rules="[{ required: true, message: 'account is required' }, { minLength: 6, message: 'must be greater than 6 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.account" placeholder="please enter your emial..." />
                </a-form-item>
                <a-form-item feedback field="password" label="Password"
                    :rules="[{ required: true, message: 'password is required' }, { minLength: 6, message: 'must be greater than 6 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.password" type="password" placeholder="please enter your password..." />
                </a-form-item>
                <a-form-item field="code" tooltip="第一次登录需填写" label="InviteCode">
                    <a-verification-code :formatter="(inputValue) => /^\d*$/.test(inputValue) ? inputValue : false"
                        :length="8" v-model="authForm.code" style="width: 200px" />
                </a-form-item>
                <a-form-item feedback v-if="authForm.code.length == 8" field="name" label="Nickname"
                    :rules="[{ required: true, message: 'name is required' }, { minLength: 2, message: 'must be greater than 2 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.name" placeholder="please enter your name..." />
                </a-form-item>
                <a-form-item feedback field="isRead" :rules="[{ type:'boolean', true:true, message: '请阅读用户协议并同意' }]"
                    :validate-trigger="['change', 'input']">
                    <a-checkbox v-model="authForm.isRead"> I have read the manual </a-checkbox>
                </a-form-item>
                <a-form-item>
                    <a-button type="primary" html-type="submit">login / register</a-button>
                </a-form-item>
            </a-form>
        </a-card>
    </div>
</template>
<script setup>
import { useUserStore } from '@/stores/UserStore';
import { reactive } from 'vue';

const authForm = reactive({
    account: "liuscraft@qq.com",
    password: "123123",
    code: "",
    name: "",
    isRead: false
})
const userStore = useUserStore()
const handleSubmit = () => {
    userStore.login(authForm)
}
</script>