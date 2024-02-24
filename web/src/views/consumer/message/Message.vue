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
            :title="item.title"
            :description="item.description"
        >
        </a-list-item-meta>
      </a-list-item>
    </template>
  </a-list>
</template>

<script>
import { reactive } from 'vue'
import { IconPlus, IconCheckCircle } from '@arco-design/web-vue/es/icon';
import { listMsg} from '@/apis/message.js'


const names = ['Socrates', 'Balzac', 'Plato'];
const dataSource = new Array(5).fill(null).map((_, index) => {
  return {
    index: index,
    title: names[index % names.length],
    description:
        '123',
  };
});

export default {
  setup() {
    return {
      IconPlus,
      IconCheckCircle,
      dataSource,
      listMsg,
      paginationProps: reactive({
        defaultPageSize: 15,
        total: dataSource.length
      })
    }
  },
}
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
