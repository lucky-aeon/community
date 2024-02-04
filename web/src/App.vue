<script setup>
import { RouterView } from 'vue-router';
import { useUserStore } from './stores/UserStore';

const  userStore = useUserStore()
userStore.getMenu()
</script>

<template>
  <RouterView/>
</template>

<style>
html, body ,#app{
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>
