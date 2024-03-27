<script setup>
import { RouterView } from 'vue-router';
import { useUserStore } from './stores/UserStore';
const  userStore = useUserStore()
userStore.refreshInfo()
</script>

<template>
  <RouterView/>
  <a-card hoverable >
    <div
      :style="{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
      }"
    >
      <span
        :style="{ display: 'flex', alignItems: 'center', color: '#1D2129' }"
      >
        <a-avatar
          :style="{ marginRight: '8px', backgroundColor: '#165DFF' }"
          :size="28"
          image-url="/icp.png"
        >
          A
        </a-avatar>
        <a-typography-text>    <a-link href="https://beian.mps.gov.cn/#/query/webSearch?code=42088102000211" rel="noreferrer" target="_blank">鄂公网安备42088102000211</a-link>
</a-typography-text>
      </span>
      <a-link>Power by Xhy, 2024-至今</a-link>
    </div>
  </a-card>
</template>

<style>
html, body ,#app{
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>
