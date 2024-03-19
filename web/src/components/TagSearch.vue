<template>
    <div style="width: 100%;">
        <a-dropdown :style="{ width: '100%' }" :popup-container="`#${randomId}`" @select="selectResult"
            :popup-max-height="true">

            <a-input-tag @press-enter="pressEnter" ref="inputRef" allow-clear unique-value v-model="tags"
                v-model:input-value="inputData" placeholder="文章标签" />


            <div style="position: relative;" :id="randomId"></div>

            <template #content>
                <template v-if="userTags.length">
                    <a-doption v-for="item in userTags" :key="randomId+item.label" :value="item">{{ item.label }}</a-doption>
                </template>
                <a-empty v-else />
            </template>
        </a-dropdown>
        <a-modal :ok-loading="askCreate.loading" v-model:visible="askCreate.show" @ok="handlerAskOk"
            @cancel="askCreate.show = false">

            <template #title>
                新增标签
            </template>
            <div>
                <a-form :model="newTagData">
                    <a-form-item field="tag" label="标签名">
                        <a-input v-model:model-value="newTagData.tag" placeholder="标签名" />
                    </a-form-item>
                    <a-form-item field="desc" label="标签描述">
                        <a-textarea v-model:model-value="newTagData.desc" placeholder="标签描述" />
                    </a-form-item>
                </a-form>
            </div>
        </a-modal>
    </div>
</template>

<script setup>
import { apiTagCreate } from '@/apis/tags';
import { useUserStore } from '@/stores/UserStore';
import { computed, nextTick, reactive, ref, watch } from 'vue';
const porps = defineProps({
    defaultData: {
        type: Array,
        default: () => []
    }
})
const randomId = ref(`tagSearchInputPop-${Math.floor((Math.random()*100)+1)}`)
const newTagData = ref({
    tag: "",
    desc: ""
})
const tags = defineModel({
    required: true, default: []
})
const inputRef = ref("inputRef")
const inputData = ref("")
const userStore = useUserStore()
const askCreate = reactive({
    show: false,
    loading: false
})
const userTags = computed(() => userStore.userTags.filter(item => {
    // 匹配用户输入的字符串
    let leftStr = item.TagName.toLocaleLowerCase()
    let rightStr = inputData.value
    if (leftStr.indexOf(rightStr) !== -1) {
        return !tags.value.find((elementB) => {
            if (typeof elementB != 'object') return true
            return item.TagId === elementB.value.TagId
        })
    }
    return false
    // 去除已被选择的元素
}).map(item => createTagItem(item, true)))
function createTagItem(value, closable = true) {
    let lowerName = value.TagName.toLocaleLowerCase()
    return {
        key: lowerName,
        label: `${value.TagName} (${value.ArticleCount})`,
        value,
        lowerName: lowerName,
        closable
    }
}
function selectResult(value) {
    nextTick(() => inputData.value = "")
    inputRef.value.focus()
    tags.value.push(value)
}
function changeUsedTags(value, add = true) {
    if (add) {
        let temp = Object.assign([], tags.value)
        temp.push(value)
        tags.value = [...new Set(temp)]
    }
}
function findUsedTag(value) {
    return userTags.value.find(item => item.lowerName === value.toLocaleLowerCase())
}
function pressEnter(v) {
    nextTick(()=>{
        tags.value.pop()
    })
    if (userTags.value.length > 0) {
        let temp = findUsedTag(v)
        if (temp) {
            changeUsedTags(temp)
            return
        }
    }
    askCreate.loading = false
    askCreate.show = true
    newTagData.value = {
        tag: v,
        desc: ""
    }

    return
}
function handlerAskOk() {
    askCreate.loading = true
    apiTagCreate(newTagData.value.tag, newTagData.value.desc).then(({ data, ok }) => {
        askCreate.loading = false
        if (!ok) return
        let tagItem = {
            TagName: data.tag,
            ArticleCount: 0,
            TagId: data.id
        }
        userStore.userTags.push(tagItem)
        changeUsedTags(createTagItem(tagItem))
        askCreate.show = false
    })
}
watch(() => porps.defaultData, (newV) => {
    tags.value =  userTags.value.filter(item => newV.find(sub => sub.toLocaleLowerCase() == item.lowerName))
})
</script>