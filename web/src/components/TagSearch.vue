<template>
    <div style="width: 100%;">
        <a-dropdown :style="{ width: '100%' }" popup-container="#tagSearchInputPop" @select="selectResult"
            :popup-max-height="true">
            <a-input-tag @press-enter="pressEnter" ref="inputRef" allow-clear unique-value v-model="tags"
                v-model:input-value="inputData" placeholder="文章标签" />
            <div style="position: relative;" id="tagSearchInputPop"></div>
            <template #content>
                <template v-if="userTags.length">
                    <a-doption v-for="item in userTags" :key="item.label" :value="item">{{ item.label }}</a-doption>
                </template>
                <a-empty v-else />
            </template>
        </a-dropdown>
    </div>
</template>

<script setup>
import { useUserStore } from '@/stores/UserStore';
import { computed, nextTick, onMounted, ref } from 'vue';
defineProps({
    show: {
        typeof: 'boolean',
        default: false
    }
})
const inputRef = ref("inputRef")
const inputData = ref("")
const tags = ref([])
const userStore = useUserStore()

const userTags = computed(() => userStore.userTags.filter(item => {
    // 匹配用户输入的字符串
    let leftStr = item.TagName.toLocaleLowerCase()
    let rightStr = inputData.value
    if (leftStr.indexOf(rightStr) !== -1) {
        return !tags.value.find((elementB) => item.TagId === elementB.value.TagId)
    }
    return false
    // 去除已被选择的元素
}).map(item => ({
    label: `${item.TagName} (${item.ArticleCount})`,
    value: item,
    closable: true
})))

function selectResult(value) {
    nextTick(() => inputData.value = "")
    inputRef.value.focus()
    tags.value.push(value)
}
function pressEnter(v) {
    if (userTags.value.length > 0) {
        return
    }
    console.log(v)
    return
}
onMounted(() => {
    console.log()
})
</script>