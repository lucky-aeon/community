<template>
    <a-layout v-if="currentUserId">
        <a-card
            style="padding: 15px; text-align: center;background: url('https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/6480dbc69be1b5de95010289787d64f1.png~tplv-uwbnlip3yd-webp.webp') center/cover no-repeat;">
            <a-avatar v-if="userData.name" :size="60" :image-url="getFileUrl(userData.avatar)">{{ userData.name
                }}</a-avatar>
            <h3>{{ userData.name }}</h3>
            <span style="display: block;">{{ userData.desc }}</span>
            <a-space>
                <!-- <a-tag color="red">{{ userData.roleUp }}</a-tag> -->
                <a-tag :color="tagItem.color || 'blue'" v-for="tagItem in userTags" :key="tagItem.id">{{ tagItem.name
                    }}</a-tag>
                <a-tag>{{ userData.createdAt }}</a-tag>
            </a-space>
        </a-card>
        <div style="height: 20px;"></div>
        <a-row v-if="false" gutter="20" :style="{ marginBottom: '20px' }">
            <a-col :span="8">
                <a-card title="Arco Card" :bordered="false" :style="{ width: '100%' }">
                    <template #extra>
                        <a-link>More</a-link>
                    </template>
                    Card content
                </a-card>
            </a-col>
            <a-col :span="8">
                <a-card title="Arco Card" :bordered="false" :style="{ width: '100%' }">
                    <template #extra>
                        <a-link>More</a-link>
                    </template>
                    Card content
                </a-card>
            </a-col>
            <a-col :span="8">
                <a-card title="Arco Card" :bordered="false" :style="{ width: '100%' }">
                    <template #extra>
                        <a-link>More</a-link>
                    </template>
                    Card content
                </a-card>
            </a-col>
        </a-row>
        <a-row :gutter="20">
            <a-col :span="16">
                <a-card title="Latest Article" :bordered="false" :style="{ width: '100%' }">
                    <template #extra>
                        <a-radio-group v-model="queryData.state" type="button">
                            <a-radio :value="2">文章</a-radio>
                            <a-radio :value="3">待解决</a-radio>
                            <a-radio :value="4">已解决</a-radio>
                            <a-radio :value="5" v-if="userStore.userInfo.role == 'admin'">私密提问</a-radio>
                        </a-radio-group>
                    </template>

                    <ArticleList :query-data="queryData" />
                </a-card>
            </a-col>
            <a-col :span="8">
                <a-card title="个人描述" :bordered="false" :style="{ width: '100%' }">
                    <template #extra>
                        <!-- <a-link>订阅</a-link> -->
                    </template>
                    {{ userData.desc }}
                </a-card>
            </a-col>
        </a-row>
    </a-layout>
    <a-result status="error" title="未找到该用户信息 " v-else>
        <template #icon>
            <IconFaceFrownFill />
        </template>
        <template #subtitle> NOT FIND THE USER </template>

        <template #extra>
            <a-button type="primary" @click="router.back()">返回上一页</a-button>
        </template>
        <a-typography style="background: var(--color-fill-2); padding: 24px;">
            <a-typography-paragraph>可能的原因:</a-typography-paragraph>
            <ul>
                <li> 该用户将自己信息设为私有 </li>
                <li> 并不存在该用户 </li>
                <li> 服务发生错误(但不可能) </li>
            </ul>
        </a-typography>
    </a-result>

</template>
<script setup>
import { apiGetFile } from '@/apis/file';
import { apiGetUserTags2, getUserInfo } from '@/apis/user';
import ArticleList from '@/components/article/ArticleList.vue';
import router from '@/router';
import { useUserStore } from '@/stores/UserStore';
import { IconFaceFrownFill } from '@arco-design/web-vue/es/icon';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
const route = useRoute()
const currentUserId = computed(() => {
    let id = route.params.userId
    try {
        return parseInt(id)
    } catch (_) {
        return null
    }
})
const userStore = useUserStore()
const queryData = ref(null)
const userData = ref({})
const userTags = ref([])
const getFileUrl = (fileKey) => {
    return apiGetFile(fileKey)

}
if (currentUserId.value) {
    apiGetUserTags2(currentUserId.value).then(({ data, ok }) => {
        if (!ok) return
        userTags.value = data
    })
    getUserInfo(currentUserId.value).then(({ data, ok }) => {
        if (!ok) {
            currentUserId.value = null
            return
        }
        userData.value = data
        userData.value.roleUp = data.role.toUpperCase()
    })
    queryData.value = {
        tags: [],
        context: "",
        state: 2,
        orderBy: "created_at",
        descOrder: true,
        userID: currentUserId.value
    }
}

</script>
