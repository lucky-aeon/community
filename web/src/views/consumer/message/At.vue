<template>
  <a-list
      class="list-demo-action-layout"
      :bordered="false"
      :data="dataSource"
      :pagination-props="paginationProps"
  >
    <template #item="{ item }">
      <a-list-item class="list-demo-item" action-layout="vertical">
        <template #actions>
          <span><icon-pen-fill />2024-01-01</span>
        </template>
        <a-list-item-meta
            :title="item.content"
        >
        </a-list-item-meta>
      </a-list-item>
    </template>
  </a-list>
</template>

<script setup>
import {reactive, ref} from 'vue'
// import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';
import { apiListMsg } from '@/apis/message.js'

const dataSource = ref([])
const  count = ref()
apiListMsg(2,0).then(({data})=>{
  dataSource.value = data.data
  count.value = data.count
})
const paginationProps = reactive({
  defaultPageSize: 15,
  total: count
})
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
</style>
