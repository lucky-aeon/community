<template>
  <a-row gutter="20" v-if="articleData">
    <a-col :span="17">

      <a-page-header :show-back="false" :style="{ background: 'var(--color-bg-2)', paddingBottom: '0px' }"
        :title="articleData.title">
        <template #breadcrumb>
          <a-breadcrumb>
            <router-link to="/"><a-breadcrumb-item>首页</a-breadcrumb-item></router-link>
            <router-link :to="`/qa/${articleData.type.flag}`"><a-breadcrumb-item>问题列表</a-breadcrumb-item></router-link>
            <a-breadcrumb-item>问题详情</a-breadcrumb-item>
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
          <a-space>
            <a-tag color="green" size="large" v-if="articleData.state == 4">已解决</a-tag>
            <a-tag color="red" size="large" v-else>未解决</a-tag>
            <a-tag color=blue size="large">
              {{ articleData.type.title || "未知" }}
            </a-tag>
          </a-space>
          <a-divider direction="vertical" />
          <a-space>
            <a-tooltip v-for="tagItem in articleData.tags" :key="tagItem.name" :content="tagItem.description"><a-tag
                color="red">{{ tagItem.name }}</a-tag></a-tooltip>
          </a-space>
        </div>
      </a-page-header>
      <div id="markdown-container"></div>
      <a-card :bordered="false" style="margin-top: 10px;">
        <comment-edit :callback="getRootComment" :article-id="articleData.id"  style="margin-bottom: 15px;"/>
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
      <UserInfoCard :user-id="articleData.user.id" />
    </a-col>
  </a-row>
</template>

<script setup>
import { apiSubscribe, apiSubscribeState, } from "@/apis/apiSubscribe.js";
import { apiArticleLike, apiArticleLikeState, apiArticleView } from '@/apis/article';
import { apiGetArticleComment } from '@/apis/comment';
import { apiGetFile } from "@/apis/file";
import CommentEdit from '@/components/comment/CommentEdit.vue';
import CommentItem from '@/components/comment/CommentItem.vue';
import UserInfoCard from "@/components/user/UserInfoCard.vue";
import router from "@/router";
import { IconDown, IconThumbUp, IconThumbUpFill } from '@arco-design/web-vue/es/icon';
import Cherry from 'cherry-markdown';
import 'cherry-markdown/dist/cherry-markdown.css';
import { nextTick, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
var CustomHookA = Cherry.createSyntaxHook('codeBlock', Cherry.constants.HOOKS_TYPE_LIST.PAR, {
  makeHtml(str) {
    console.warn('custom hook', 'hello');
    return str;
  },
  rule() {
    const regex = {
      begin: '',
      content: '',
      end: '',
    };
    regex.reg = new RegExp(regex.begin + regex.content + regex.end, 'g');
    return regex;
  },
});
const getFileUrl = (fileKey) => apiGetFile(fileKey)

var cherryConfig = {
  id: 'markdown-container',
  value: '# welcome to cherry editor! \n awdwadad',
  externals: {
    echarts: window.echarts,
    katex: window.katex,
    MathJax: window.MathJax,
  },
  engine: {
    global: {
      urlProcessor(url, srcType) {
        console.log(`url-processor`, url, srcType);
        return url;
      },
    },
    syntax: {
      fontEmphasis: {
        allowWhitespace: true, // 是否允许首尾空格
      },
      mathBlock: {
        engine: 'MathJax', // katex或MathJax
        src: 'https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js', // 如果使用MathJax将js在此处引入，katex则需要将js提前引入
      },
      inlineMath: {
        engine: 'MathJax', // katex或MathJax
      },
      emoji: {
        useUnicode: false,
        customResourceURL: 'https://github.githubassets.com/images/icons/emoji/unicode/${code}.png?v8',
        upperCase: true,
      },
      // toc: {
      //     tocStyle: 'nested'
      // }
      // 'header': {
      //   strict: false
      // }
    },
    customSyntax: {
      // SyntaxHookClass
      CustomHook: {
        syntaxClass: CustomHookA,
        force: false,
        after: 'br',
      },
    },
  },
  toolbars: {
    toolbar: false,
    toc: {
      updateLocationHash: true, // 要不要更新URL的hash
      defaultModel: 'full', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
    },
  },
  editor: {
    defaultModel: 'previewOnly',
  },
  callback: {
    onClickPreview: function (e) {
      const { target } = e;
      if (target.tagName === 'IMG') {
        console.log('click img', target);
        // eslint-disable-next-line no-undef
        const tmp = new Viewer(target, {
          button: false,
          navbar: false,
          title: [1, (image, imageData) => `${image.alt.replace(/#.+$/, '')} (${imageData.naturalWidth} × ${imageData.naturalHeight})`],
          hidden() {
            tmp.destroy()
          },
        });
        tmp.show();
      }
    }
  },
  previewer: {
    // 自定义markdown预览区域class
    // className: 'markdown'
  },
  keydown: [],
  //extensions: [],
};

const articleData = ref(null)
const route = useRoute()
const likeState = ref(true)

const articleSubscribe = ref("订阅问题")

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
  })
}

watch(() => route.fullPath, () => {
  apiArticleView(route.params.id).then(({ data }) => {
    articleData.value = data
    nextTick(() => {
      cherryConfig.value = articleData.value.content
      new Cherry(cherryConfig);
      setTimeout(() => window.scrollTo(0, 0), 1)
    })
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
    getRootComment()
  }).catch(() => {
    router.back()
  })
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


</script>

<style>
.arco-dropdown-open .arco-icon-down {
  transform: rotate(180deg);
}

.cherry {
  box-shadow: none;
}

.cherry-flex-toc {
  position: fixed
}
</style>