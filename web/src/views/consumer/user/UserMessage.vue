

<template>
   <a-card :bordered="false" style="margin-top: 10px;">
  <a-tabs v-model:active-key="msgType">
    <a-tab-pane :key="1" title="通知">
    </a-tab-pane>
    <a-tab-pane :key="2" title="@我">
    </a-tab-pane>
    <template #extra>
      <a-button type="text" @click="clearMsg">
        清空
      </a-button>
    </template>
  </a-tabs>
  <msg-notice :msg-type="msgType" :reload="msgReload"/>
</a-card>
</template>
<script setup>
import { apiClearUnReadMsg } from '@/apis/message.js';
import MsgNotice from "@/components/message/MsgNotice.vue";

import { ref } from "vue";

const msgType = ref(1)
const msgReload = ref(0)
function clearMsg(){
  apiClearUnReadMsg(msgType.value).then(({ok})=> {if(ok) {msgReload.value++}})
}
</script>
<style scoped>

</style>