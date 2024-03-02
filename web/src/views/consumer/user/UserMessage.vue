

<template>
  <a-tabs v-model:active-key="msgType">
    <a-tab-pane :key="1" title="通知">
    </a-tab-pane>
    <a-tab-pane :key="2" title="@我">
    </a-tab-pane>
  </a-tabs>
  <msg-notice :msg-type="msgType"/>

</template>
<script setup>
import MsgNotice from "@/components/message/MsgNotice.vue";
import {ref} from "vue";

const msgType = ref(1)
</script>
<style scoped>

</style>