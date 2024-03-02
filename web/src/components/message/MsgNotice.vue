<template>
  <a-list
      class="list-demo-action-layout"
      :bordered="false"
      :data="dataSource"
      :pagination-props="paginationProps"
      v-if="dataSource.length"
  >
    <template #item="{ item }">
      <a-list-item class="list-demo-item" :style="{padding: '5px'}" action-layout="vertical"
                   @click="router.push(`/article/view/${item.articleId}`)" >
        <template #actions>
          <span class="arco-typography time-text">{{item.createdAt}}</span>
        </template>
        <a-list-item-meta
            :description="item.content"
        >
        </a-list-item-meta>
      </a-list-item>
    </template>
  </a-list>
  <slot name="empty" v-else>
  </slot>
</template>

<script setup>
import {reactive, ref, watch} from 'vue'
import { apiListMsg } from '@/apis/message.js'
import router from "@/router/index.js";
const props = defineProps({
  msgType: {
    type: Number,
    default: 1
  },
  msgState: {
    type: Number,
    default: 1
  }
})
const dataSource = ref([])
const  count = ref()

const paginationProps = reactive({
  defaultPageSize: 15,
  total: count
})
const getMsg = ()=>{
  apiListMsg(props.msgType,props.msgState).then(({data})=>{
    dataSource.value = data.data
    count.value = data.count
  })
}
watch(()=>props, ()=>{
  getMsg()
}, {deep: true, immediate: true})
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
