<template>
  <a-list class="list-demo-action-layout" :bordered="false" :data="dataSource" :pagination-props="paginationProps"
    v-if="dataSource.length" @page-change="getMsg">
    <template #item="{ item, index }">
      <a-list-item class="list-demo-item" :style="{ padding: '5px' }" action-layout="vertical" @click="handlerMsg(item, index)">

        <template #actions>
          <span class="arco-typography time-text">{{ item.createdAt }}</span>
        </template>
        <a-badge :count="item.state" dot>
          <a-list-item-meta :description="item.content">
          </a-list-item-meta></a-badge>
      </a-list-item>
    </template>
  </a-list>
  <a-empty v-else />
</template>

<script setup>
import { apiListMsg, apiPostRead } from '@/apis/message.js';
import router from "@/router/index.js";
import { useMessage } from '@/stores/MessageStore';
import { reactive, ref, watch } from 'vue';
const props = defineProps({
  msgType: {
    type: Number,
    default: 1
  },
  msgState: {
    type: Number,
    default: 0
  },
  reload: {
    type: Number,
    default: 0
  }
})
const dataSource = ref([])
const count = ref()

const paginationProps = reactive({
  defaultPageSize: 15,
  total: count,
  current:1
})
const messageStore = useMessage()
const getMsg = (current) => {
  messageStore.getUnReadCount()
  paginationProps.current = current
  apiListMsg(current,10,props.msgType, props.msgState).then(({ data }) => {
    dataSource.value = data.data
    count.value = data.count
  })
}
const handlerMsg = (item) => {
  if(item.state) {
    apiPostRead([item.id]).then(({ok})=>{
      if(!ok) return
      getMsg()
    })
  }
  router.push(`/article/view/${item.articleId}`)
}
watch(() => props, () => {
  getMsg()
}, { deep: true, immediate: true })
</script>

<style scoped>
.list-demo-action-layout .image-area {
  width: 183px;
  height: 119px;
  border-radius: 2px;
  overflow: hidden;
}

.list-demo-action-layout .list-demo-item {
  padding: 20px 0;
  border-bottom: 1px solid var(--color-fill-3);
}

.list-demo-action-layout .image-area img {
  width: 100%;
}

.list-demo-action-layout .arco-list-item-action .arco-icon {
  margin: 0 4px;
}

.list-demo-item:hover {
  transition: background-color 0.3s ease;
  background-color: #f4f4f491;
  cursor: pointer;
}
</style>
