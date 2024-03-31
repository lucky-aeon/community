<template>
    <ALayout>
        <ALayoutHeader v-if="!route.meta.hideSelf">
            <a-card :bordered="false" direction="vertical" size="large"
                :body-style="{ paddingBottom: '0', paddingTop: '0' }">
                <a-space>
                    <a-upload :custom-request="userAvatar.uploadAvatar"
                        :fileList="userAvatar.file ? [userAvatar.file] : []" :show-file-list="false"
                        @change="userAvatar.onChange" @progress="userAvatar.onProgress">
                        <template #upload-button>
                            <div
                                :class="`arco-upload-list-item${userAvatar.file && userAvatar.file.status === 'error' ? ' arco-upload-list-item-error' : ''}`">
                                <div class="arco-upload-list-picture custom-upload-avatar"
                                    v-if="userAvatar.file && userAvatar.file.url">
                                    <img :src="userAvatar.file.url" />
                                    <div class="arco-upload-list-picture-mask">
                                        <IconEdit />
                                    </div>
                                    <a-progress
                                        v-if="userAvatar.file.status === 'uploading' && userAvatar.file.percent < 100"
                                        :percent="userAvatar.file.percent" type="circle" size="mini"
                                        :style="{ position: 'absolute', left: '50%', top: '50%', transform: 'translateX(-50%) translateY(-50%)' }" />
                                </div>
                                <div class="arco-upload-picture-card" v-else>
                                    <div class="arco-upload-picture-card-text">
                                        <IconPlus />
                                        <div style="margin-top: 10px; font-weight: 600">Upload</div>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </a-upload>
                    <a-descriptions :data="userInfo" align="right" />
                </a-space>
                <a-space style="margin: 5px;">
                    <a-tag :color="tagItem.color || 'blue'" v-for="tagItem in userSelfTags" :key="tagItem.id">{{
            tagItem.name }}</a-tag>
                </a-space>
            </a-card>
            <!-- <a-modal simple v-model:visible="userAvatarModal.show">
                
            </a-modal> -->
        </ALayoutHeader>
        <ALayoutContent>
            <RouterView />
        </ALayoutContent>
    </ALayout>

</template>

<script setup>
import { apiGetFile, apiUploadFile } from '@/apis/file';
import { apiUserEditUserAvatar, apiGetUserTags2 as apiGetUserTags } from '@/apis/user';
import { useUserStore } from '@/stores/UserStore';
import { IconEdit, IconPlus } from '@arco-design/web-vue/es/icon';
import { ref, watch } from 'vue';
import { onMounted } from 'vue';
import { reactive } from 'vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const userStore = useUserStore()
const route = useRoute()
const userSelfTags = ref([])
const userAvatar = reactive({
    show: true,
    file: undefined,
    onChange(_, currentFile) {
        userAvatar.file = {
            ...currentFile,
            // url: URL.createObjectURL(currentFile.file),
        };
    },
    onProgress(currentFile) {
        userAvatar.file = currentFile;
    },
    uploadAvatar(options) {
        apiUploadFile(userStore.userInfo.id, options.fileItem.file, (r, key) => {
            if (!r.ok) {
                options.onError(r)
            }
            apiUserEditUserAvatar(key)
        }, (progressed, event) => {
            options.onProgress(progressed, event.event)
        })
    }
})
const userInfo = ref([])
watch(() => userStore.userInfo, () => {
    if (userSelfTags.value.length === 0) {
        apiGetUserTags(userStore.userInfo.id).then(({ data, ok }) => {
            if (!ok) return
            userSelfTags.value = data
        })
    }
    userInfo.value = [{
        label: '账号',
        value: userStore.userInfo.account,
    },
    {
        label: '注册时间',
        value: userStore.userInfo.createdAt,
    }, {
        label: "描述",
        value: userStore.userInfo.desc
    }]
}, {
    immediate: true,
    deep: true
})
userStore.refreshTags()
onMounted(() => {
    userAvatar.file = {
        uid: userStore.userInfo.id,
        name: userStore.userInfo.name,
        status: true,
        url: apiGetFile(userStore.userInfo.avatar)
    }
})
</script>