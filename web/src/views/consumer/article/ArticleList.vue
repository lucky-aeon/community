<template>
  <a-row gutter="10">
    <a-col :span="17">
      <a-card title="文章列表" :bordered="false" :style="{ width: '100%' }">
        <template #extra>
          <a-link>排序</a-link>
        </template>
        <ArticleListCom :queryData="queryData" />
      </a-card>
    </a-col>
    <a-col :span="7">
      <TopList />
    </a-col>
  </a-row>
</template>
<script setup>
import TopList from '@/components/TopList.vue';
import ArticleListCom from "@/components/article/ArticleList.vue";
const queryData = defineModel({
  required: true,
  default: {
    tags: [],
    context: "",
    state: 2,
    orderBy: "created_at",
    descOrder: true
  }
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


.list-demo-item:hover {
  transition: background-color 0.3s ease;
  background-color: #f4f4f491;
  cursor: pointer;
}
</style>