<template>
  <a-row gutter="20" v-if="articleData">
    <a-col :span="17">
      <a-spin :loading="loadingState.article" tip="This may take a while..." style="width: 100%;">
        <a-page-header :show-back="false" :style="{ background: 'var(--color-bg-2)', paddingBottom: '0px' }"
          :title="articleData.title">
          <template #breadcrumb>
            <a-breadcrumb>
              <router-link to="/"><a-breadcrumb-item>首页</a-breadcrumb-item></router-link>
              <router-link
                :to="`/article/${articleData.type.flag}`"><a-breadcrumb-item>文章列表</a-breadcrumb-item></router-link>
              <a-breadcrumb-item>文章详情</a-breadcrumb-item>
            </a-breadcrumb>
          </template>
          <template #subtitle>
            <a-space>

              <span><a-avatar :image-url="getFileUrl(articleData.user.avatar)"
                  :style="{ marginRight: '8px', backgroundColor: '#165DFF' }" :size="28">
                  load
                </a-avatar>
                <RouterLink :to="`/user/${articleData.user.id}`">{{ articleData.user.name }}</RouterLink>
              </span>
              <span>时间: {{ articleData.updatedAt }}</span>
            </a-space>
          </template>
          <template #extra>
            <a-space>
              <AButtonGroup>
                <AButton @click="like(articleData.id)">
                  <icon-thumb-up-fill size="large" v-if="likeState" />
                  <icon-thumb-up size="large" v-else />{{ articleData.like }}
                </AButton>

              </AButtonGroup>
              <a-dropdown-button @click="subscribe(1, articleData.id)">
                {{ articleSubscribe }}
                <template #icon>
                  <icon-down />
                </template>
                <template #content>
                  <a-doption @click="subscribe(2, articleData.user.id)">{{ userSubscribe }}</a-doption>
                </template>
              </a-dropdown-button>

            </a-space>
          </template>
          <div>
            <a-tag color=blue size="large">
              {{ articleData.type.title || "未知" }}
            </a-tag>
            <a-divider direction="vertical" />
            <a-space>
              <a-tooltip v-for="tagItem in articleData.tags" :key="tagItem.name" :content="tagItem.description"><a-tag
                  color="red">{{ tagItem.name }}</a-tag></a-tooltip>
            </a-space>
          </div>
        </a-page-header>
        <MarkdownEdit :render-nav="setArticleNavs" preview-only :show-nav="false" v-model="articleData.content" />
      </a-spin>
      <a-card :loading="loadingState.comment" :bordered="false" style="margin-top: 10px;">
        <comment-edit :callback="getRootComment" :article-id="articleData.id" style="margin-bottom: 15px;" />
        <template v-if="articleCommentList.length">
          <comment-item :article="articleData" :callback="getRootComment" v-for="comment in articleCommentList"
            :comment="comment" :key="comment.id" />
          <a-divider />
          <a-pagination @change="getRootComment" :total="paginationData.total"
            v-model:page-size="paginationData.pageSize" v-model:current="paginationData.current" show-page-size />
        </template>
        <a-result v-else title="没有评论" subtitle="快来参与互动呀">
        </a-result>
      </a-card>
    </a-col>
    <a-col span="7">
      <a-grid row-gap="5" cols="1">
        <a-col>
          <UserInfoCard :user-id="articleData.user.id" />
        </a-col>
        <a-col>
          <a-affix v-if="articleNavs.length > 0">
            <a-card style="margin-top: 10px;max-height: 400px;" :bordered="false" title="目录">
              <a-anchor smooth style="width: 100%;max-height: 300px;">
                <MyAnchorLink v-for="item in articleNavs" :data="item" style="width: 100%;" />
              </a-anchor>
            </a-card>
          </a-affix>
        </a-col>
        <a-col>
          <a-card :bordered="false" body-style="padding: 0px;" title="用户文章">
            <a-list :loading="loadingState.userlink" :split="false" :size="'small'" :bordered="false"
              :data="userLinkArticles">
              <template #item="{ item }">
                <a-list-item class="list-demo-item" style="padding: 0 9px;" action-layout="vertical"
                  @click="router.push(`/article/view/${item.id}`)">
                  <a-list-item-meta :title="item.title">
                    <template #description>
                      <small>{{ item.like }}点赞</small>
                    </template>
                  </a-list-item-meta>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-grid>
    </a-col>
  </a-row>
</template>

<script setup>
import MyAnchorLink from '@/components/MyAnchorLink.vue'
import UserInfoCard from '@/components/user/UserInfoCard.vue'
import { apiSubscribe, apiSubscribeState, } from "@/apis/apiSubscribe.js";
import { apiArticleLike, apiArticleLikeState, apiArticleList, apiArticleView } from '@/apis/article';
import { apiGetArticleComment } from '@/apis/comment';
import { apiGetFile } from "@/apis/file";
import CommentEdit from '@/components/comment/CommentEdit.vue';
import CommentItem from '@/components/comment/CommentItem.vue';
import router from "@/router";
import { IconDown, IconThumbUp, IconThumbUpFill, IconHeart } from '@arco-design/web-vue/es/icon';
import 'cherry-markdown/dist/cherry-markdown.css';
import { watch } from "vue";
import { useRoute } from 'vue-router';
import MarkdownEdit from '@/components/MarkdownEdit.vue';
import { ref } from 'vue';
import { reactive } from 'vue';
const getFileUrl = (fileKey) => apiGetFile(fileKey)

const articleData = ref(null)
const articleNavs = ref([])
const route = useRoute()
const likeState = ref(true)
const userLinkArticles = ref([])
const loadingState = reactive({
  comment: true,
  userlink: true,
  article: true
})

const articleSubscribe = ref("订阅文章")

const userSubscribe = ref("订阅用户")
const articleCommentList = ref([])
const paginationData = ref({
  current: 1,
  total: 0,
  pageSize: 5

})
function getRootComment() {
  apiGetArticleComment(articleData.value.id, paginationData.value.current, paginationData.value.pageSize).then(({ data, ok }) => {
    if (!ok || !data.data) {
      return
    }
    paginationData.value.total = data.count
    articleCommentList.value = data.data
  }).finally(() => loadingState.comment = false)
}
watch(() => route.fullPath, () => {
  loadingState.userlink = true
  loadingState.comment = true
  loadingState.article = true
  apiArticleView(route.params.id).then(({ data }) => {
    articleData.value = data
    const id = articleData.value.id
    apiArticleLikeState(id).then((data) => {
      likeState.value = data.data
    })
    apiSubscribeState(1, id).then((data) => {
      articleSubscribe.value = data.data ? "文章已订阅" : "订阅文章"
    })
    apiSubscribeState(2, articleData.value.user.id).then((data) => {
      userSubscribe.value = data.data ? "用户已订阅" : "订阅用户"
    })
    apiArticleList({ userId: data.user.id, state: 2 }, 1, 5).then(({ data, ok }) => {
      if (!ok) return
      userLinkArticles.value = data.list
    }).finally(() => loadingState.userlink = false)
    getRootComment()
  }).catch(() => {
    router.back()
  }).finally(() => loadingState.article = false)
}, { immediate: true })

function like(id) {
  apiArticleLike(id).then((data) => {
    articleData.value.like += (data.data ? 1 : -1)
    likeState.value = data.data
  })

}

function subscribe(eventId, businessId) {
  apiSubscribe(eventId, businessId).then((data) => {
    if (eventId == 1) {
      articleSubscribe.value = data.data ? "文章已订阅" : "订阅文章"
    } else {
      userSubscribe.value = data.data ? "用户已订阅" : "订阅用户"
    }
  })
}
function arrayToTree(arr, level = 1) {
  const tree = [];

  for (let i = 0; i < arr.length; i++) {
    const item = arr[i];
    if (level != 1 && item.level == 1) break
    if (item.level === level) {
      const node = {
        id: item.id,
        text: item.text,
        children: arrayToTree(arr.slice(i + 1), level + 1)
      };
      tree.push(node);
    }
  }

  return tree;
}

function setArticleNavs(arr) {
  console.log(arr)
  articleNavs.value = arrayToTree(arr)
}
</script>

<style scoped>
.arco-dropdown-open .arco-icon-down {
  transform: rotate(180deg);
}

.cherry {
  box-shadow: none;
}

.cherry-flex-toc {
  position: fixed
}

.list-demo-item:hover {
  transition: background-color 0.3s ease;
  background-color: #f4f4f491;
  cursor: pointer;
}
</style>