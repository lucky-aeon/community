<template>
  <a-row gutter="10">
    <a-col :span="16">
      <a-card title="问答列表" :bordered="false" :style="{ width: '100%' }">
        <template #extra>
          <a-radio-group v-model="model.state" type="button" :default-value="3">
            <a-radio :value="3">待解决</a-radio>
            <a-radio :value="4">已解决</a-radio>
            <a-radio :value="5" v-if="userStore.userInfo.role == 'admin'">私密提问</a-radio>
          </a-radio-group>
        </template>
        <QaListCom :queryData="model" />
      </a-card>
    </a-col>
    <a-col :span="8">
      <TopList />
    </a-col>
  </a-row>
</template>
<script setup>
import TopList from '@/components/TopList.vue'
import QaListCom from "@/components/qa/QaList.vue";
import { useUserStore } from "@/stores/UserStore";
const userStore = useUserStore()
const model = defineModel({
  required: true, default: {
    tags: [],
    context: "",
    state: 3,
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