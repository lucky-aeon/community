<template>
    <a-card :bordered="false" style="min-width: 300px">
        <a-card-meta v-if="userInfo.id && userInfo.name">
            <template #title>
                <a-comment :datetime="userInfo.createdAt">
                    <template #author>
                        <a-typography-text style="cursor: pointer;"
                            @click="router.push({ path: `/user/${userInfo.id}` })">
                            {{ userInfo.name }}
                        </a-typography-text>
                    </template>
                    <template #avatar>
                        <a-avatar :image-url="apiGetFile(userInfo.avatar)">
                        </a-avatar>
                    </template>
                    <template #content>
                        <a-typography-text type="secondary">
                            {{ userInfo.desc || "👋 TA 很懒，没有留下任何东西" }}
                        </a-typography-text>
                        <br />
                        <a-space v-if="userTags.length">
                            <!-- <a-tag color="red">{{ userData.roleUp }}</a-tag> -->
                            <a-tag size="small" :color="tagItem.color || 'blue'" v-for="tagItem in userTags"
                                :key="tagItem.id">{{ tagItem.name
                                }}</a-tag>
                        </a-space>
                        <a-tag>没标签的哈~</a-tag>
                    </template>
                </a-comment>
            </template>
        </a-card-meta>
        <a-skeleton-line :rows="3" v-else/>
    </a-card>
</template>

<script setup>
import { apiGetFile } from '@/apis/file';
import { apiGetUserTags2, getUserInfo } from '@/apis/user';
import router from '@/router';
import { ref } from 'vue';
import { watch } from 'vue';
const props = defineProps({
    userId: {
        type: Number,
        default: 0,
        required: true,
    }
})
const userInfo = ref({
    id: 0,
    avatar: '',
    name: undefined,
    email: '',
    desc: ''
})

const userTags = ref([])

watch(() => props.userId, (newV) => {
    if (newV) {
        getUserInfo(props.userId).then(({ data, ok }) => {
            if (!ok) {
                currentUserId.value = null
                return
            }
            userInfo.value = data
            userInfo.value.roleUp = data.role.toUpperCase()
        })
        apiGetUserTags2(newV).then(({ data, ok }) => {
            if (!ok) return
            userTags.value = data
        })
    }
}, { immediate: true })
</script>