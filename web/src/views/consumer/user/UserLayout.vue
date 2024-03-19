<template>
    <ALayout>
        <ALayoutHeader>
            <a-card direction="vertical" size="large">
        <a-space :size="54">
            <a-avatar :trigger-icon-style="{ color: '#3491FA' }" :auto-fix-font-size="false" @click="toast"
                :style="{ backgroundColor: '#168CFF' }">
                A
                <template #trigger-icon>
                    <IconCamera />
                </template>
            </a-avatar>
            <a-descriptions :data="userInfo" align="right" />
        </a-space>
    </a-card>
        </ALayoutHeader>
        <ALayoutContent>
            <RouterView />
        </ALayoutContent>
    </ALayout>
    
</template>

<script setup>
import { useUserStore } from '@/stores/UserStore';
import { isLogin } from '@/utils/auth';
import { IconCamera } from '@arco-design/web-vue/es/icon';
import { computed } from 'vue';

const userStore = useUserStore()
const userInfo = computed(() => {
    if(!isLogin()) {
        return {}
    }
    console.log("user info")
  console.log(userStore.userInfo)
    return [{
        label: '账号',
        value: userStore.userInfo.account,
    },
    {
        label: '注册时间',
        value: userStore.userInfo.createdAt,
    }]
})
const toast = function () {
  this.$message.info('Uploading...');
}
userStore.refreshTags()
</script>