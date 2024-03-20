<template>
    <a-layout v-if="currentUserId">
        <a-card
            style="padding: 15px; text-align: center;background: url('https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/6480dbc69be1b5de95010289787d64f1.png~tplv-uwbnlip3yd-webp.webp') center/cover no-repeat;">
            <a-avatar :size="60" :src="userData.avatar">{{ userData.name }}</a-avatar>
            <h3>{{userData.name}}</h3>
            <a-space>
                <a-tag color="red" >{{ userData.roleUp }}</a-tag>
                <a-tag color="blue">{{ userData.createdAt }}</a-tag>
            </a-space>
        </a-card>
        <div style="height: 20px;"></div>
        <a-card title="Latest Article">
            <ArticleList :query-data="queryData" />
        </a-card>

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
import { getUserInfo } from '@/apis/user';
import ArticleList from '@/components/article/ArticleList.vue';
import router from '@/router';
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
const queryData = ref(null)
const userData = ref({})
if (currentUserId.value) {
    getUserInfo(currentUserId.value).then(({data, ok})=>{
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
        orderBy: "created_at",
        descOrder: true,
        userID: currentUserId.value
    }
}

</script>
