<template>
    <a-card title="置顶内容" :bordered="false" :style="{ width: '100%' }" body-style="padding-top: 0px">
        <template #extra>
            <a-button @click="page++" :disabled="page==0">换一换</a-button>
        </template>
        <a-list :bordered="false" size="small">
            <a-list-item style="cursor: pointer;" @click="router.push(`/${item.state==2?'article':'qa'}/view/${item.id}`)" v-for="item in listData" :key="item.id">{{ item.title }}</a-list-item>
        </a-list>
    </a-card>

</template>

<script setup>
import { apiGetTopArticle } from '@/apis/article';
import router from '@/router';
import { ref } from 'vue';
import { computed } from 'vue';
import { watch } from 'vue';
import { useRoute } from 'vue-router';
const props = defineProps({
    classfily: {
        type: String,
        default: undefined
    }
})
const page = ref(1)
const route = useRoute()
const listData = ref([])
const currentClassfily = computed(() => props.classfily || (route.params.classfily || ""))
function getTopList() {
    if(page.value == 0) return
    apiGetTopArticle(currentClassfily.value, page.value).then(({ data, ok }) => {
        if(data.total/10.0 <= page.value) {
            page.value = 0
        }
        if (!ok || !data.list) return
        listData.value = data.list
    })
}
watch(() => page.value, () => {
    getTopList()
}, { immediate: true })
watch(() => route.fullPath, () => {
    getTopList()
})
</script>