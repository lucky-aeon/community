<template>
    <QaListCom :queryData="queryData"/>
  </template>
  <script setup>
  import QaListCom from "@/components/qa/QaList.vue";
  import { ref } from "vue";
  const queryData = ref({
    tags: [],
    context: "",
    state: 3,
    orderBy: "created_at",
    descOrder: true
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