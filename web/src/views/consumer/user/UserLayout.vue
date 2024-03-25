<template>
    <ALayout>
        <ALayoutHeader v-if="!route.meta.hideSelf">
            <a-card direction="vertical" size="large" :body-style="{paddingBottom: '2px'}">
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
                <a-space style="margin: 5px;">
                <a-tag :color="tagItem.color || 'blue'" v-for="tagItem in userSelfTags" :key="tagItem.id">{{
            tagItem.name }}</a-tag>
            </a-space>
            </a-card>
            
        </ALayoutHeader>
        <ALayoutContent>
            <RouterView />
        </ALayoutContent>
    </ALayout>

</template>

<script setup>
import { apiGetUserTags } from '@/apis/user';
import { useUserStore } from '@/stores/UserStore';
import { isLogin } from '@/utils/auth';
import { IconCamera } from '@arco-design/web-vue/es/icon';
import { ref } from 'vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const userStore = useUserStore()
const route = useRoute()
const userSelfTags = ref([])
const userInfo = computed(() => {
    if (!isLogin()) {
        return {}
    }
    if (userSelfTags.value.length == 0) {
        apiGetUserTags(userStore.userInfo.id).then(({ data, ok }) => {
            if (!ok) return
            userSelfTags.value = data
        })
    }
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