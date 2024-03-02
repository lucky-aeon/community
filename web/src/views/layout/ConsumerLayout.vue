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
          <AButtonGroup style="float: right;" >

            <a-tooltip content="消息通知">
            <AButton type="text"><icon-notification size="large"
                                                    class="nav-btn"
                                                    type="outline"
                                                    :shape="'circle'"
                                                    @click="setPopoverVisible"/></AButton>
            </a-tooltip>
            <a-popover
                trigger="click"
                :arrow-style="{ display: 'none' }"
                :content-style="{ padding: 0, minWidth: '400px' }"
                content-class="message-popover"
            >
              <div ref="refBtn" class="ref-btn"></div>
              <template #content>
                <a-spin style="display: block">
                  <a-tabs default-active-key="notice" type="rounded" destroy-on-hide>
                    <a-tab-pane v-for="item in tabList" :key="item.key">
                      <template #title>
                        <span> {{ item.title }} </span>
                      </template>
                      <a-card :body-style="{padding: '10px'}">
                        <MsgNotice>
                          <template #empty>
                            <a-empty/>
                          </template>
                        </MsgNotice>
                        <template #actions>
                          <a-space
                              fill
                              :size="0"
                          >
                            <div class="footer-wrap">
                              <a-link >已读</a-link>
                            </div>
                            <div class="footer-wrap">
                              <a-link>查看更多</a-link>
                            </div>
                          </a-space>
                        </template>
                      </a-card>
                    </a-tab-pane>
                    <template #extra>
                      <a-button type="text">
                        清空
                      </a-button>
                    </template>
                  </a-tabs>

                </a-spin>
              </template>
            </a-popover>
            <AButton type="text">退出</AButton>
          </AButtonGroup>

        </a-menu>
      </a-layout-header>
      <a-layout>
        <a-layout-sider style="height: 100%;" :width="220" collapsible>
          <a-menu :style="{ height: '100%' }" :default-open-keys="['0']" :default-selected-keys="['0_2']">
            <a-sub-menu v-for="item in userStore.menu" :key="item.name">
              <template #icon><icon-apps></icon-apps></template>
              <template #title>{{ item.meta.locale }}</template>
              <router-link v-for="child in item.children" :key="child.name" :to="child.path">
                <a-menu-item>{{ child.meta.locale }}</a-menu-item></router-link>
            </a-sub-menu>
            <a-sub-menu key="0">
              <template #icon><icon-apps></icon-apps></template>
              <template #title>个人空间</template>
              <router-link to="/user">
                <a-menu-item>工作台</a-menu-item>
              </router-link>
              <router-link to="/user/profile">
                <a-menu-item>用户信息</a-menu-item>
              </router-link>
              <router-link to="/user/article">
                <a-menu-item>文章管理</a-menu-item>
              </router-link>
              <router-link to="/user/comment">
                <a-menu-item>评论</a-menu-item>
              </router-link>
              <router-link to="/user/message">
                <a-menu-item>消息</a-menu-item>
              </router-link>
            </a-sub-menu>
            <a-sub-menu key="3">
              <template #icon><icon-apps></icon-apps></template>
              <template #title>等级</template>
              <router-link to="/member">
                <a-menu-item>等级</a-menu-item>
              </router-link>

            </a-sub-menu>
            <a-sub-menu key="4">
              <template #icon><icon-apps></icon-apps></template>
              <template #title>文件</template>
              <router-link to="/file">
                <a-menu-item>文件</a-menu-item>
              </router-link>

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
<script lang="ts"setup>
import {reactive, ref} from 'vue'
import { Tooltip, Typography } from '@arco-design/web-vue';
import { useUserStore } from "@/stores/UserStore";
import { IconApps, IconNotification } from "@arco-design/web-vue/es/icon";
import { RouterView } from "vue-router";
import MsgNotice from "@/components/message/MsgNotice.vue";

interface TabItem {
  key: string;
  title: string;
  avatar?: string;
}
const tabList: TabItem[] = [
  {
    key: 'notice',
    title: '通知',
  },
  {
    key: 'todo',
    title: '@',
  },
];

const refBtn = ref();
const setPopoverVisible = () => {
  const event = new MouseEvent('click', {
    view: window,
    bubbles: true,
    cancelable: true,
  });
  refBtn.value.dispatchEvent(event);
};

const userStore = useUserStore()
userStore.getMenu()
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