<template>
  <div class="layout-demo">
    <a-layout style="height: 100%;">
      <a-layout-header>
        <a-menu class="line" mode="horizontal" :default-selected-keys="['1']">
          <a-menu-item key="0" :style="{ padding: 0, marginRight: '38px' }" disabled>
            <div :style="{
              width: '80px',
              height: '30px',
              borderRadius: '2px',
              background: 'var(--color-fill-3)',
              cursor: 'text',
            }" />
          </a-menu-item>
          <a-menu-item key="1">Home</a-menu-item>
          <a-menu-item key="2">Solution</a-menu-item>
          <a-menu-item key="3">Cloud Service</a-menu-item>
          <a-menu-item key="4">Cooperation</a-menu-item>
        </a-menu>
      </a-layout-header>
      <a-layout>
        <a-layout-sider style="height: 100%;" :width="220" collapsible :collapsed="collapsed" @collapse="onCollapse">
          <a-menu :style="{ height: '100%' }" :default-open-keys="['0']" :default-selected-keys="['0_2']">
            <a-sub-menu key="0">
              <template #icon><icon-apps></icon-apps></template>
              <template #title>Navigation 1</template>
              <a-menu-item key="0_0">Menu 1</a-menu-item>
              <a-menu-item key="0_1">Menu 2</a-menu-item>
              <a-menu-item key="0_2">Menu 3</a-menu-item>
              <a-menu-item key="0_3">Menu 4</a-menu-item>
            </a-sub-menu>
            <a-sub-menu key="1">
              <template #icon><icon-bug></icon-bug></template>
              <template #title>Navigation 2</template>
              <a-menu-item key="1_0">Menu 1</a-menu-item>
              <a-menu-item key="1_1">Menu 2</a-menu-item>
              <a-menu-item key="1_2">Menu 3</a-menu-item>
            </a-sub-menu>
            <a-sub-menu key="2">
              <template #icon><icon-bulb></icon-bulb></template>
              <template #title>Navigation 3</template>
              <a-menu-item key="2_0">Menu 1</a-menu-item>
              <a-menu-item key="2_1">Menu 2</a-menu-item>
              <a-sub-menu key="2_2" title="Navigation 4">
                <a-menu-item key="2_2_0">Menu 1</a-menu-item>
                <a-menu-item key="2_2_1">Menu 2</a-menu-item>
              </a-sub-menu>
            </a-sub-menu>
          </a-menu>
        </a-layout-sider>
        <a-layout-content style="padding: 15px;">
          <RouterView />
        </a-layout-content>
      </a-layout>

    </a-layout>
  </div>
</template>
<script setup>
import { RouterView } from "vue-router";
// import { getUserMenu } from "../../apis/user";
// const {data} = await getUserMenu()
// console.log(data)
</script>
<style scoped>
.layout-demo {
  height: 100%;
}
.layout-demo :deep(.arco-layout-header),
.layout-demo :deep(.arco-layout-footer) {
  display: flex;
  flex-direction: column;
  justify-content: center;
  font-size: 16px;
  font-stretch: condensed;
}

.line {
  border: 1px solid rgba(211, 211, 211, 0.841);
}
</style>