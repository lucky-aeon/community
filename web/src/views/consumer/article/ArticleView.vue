<template>
  <div v-if="articleData">
    <a-page-header :show-back="false" :style="{ background: 'var(--color-bg-2)', paddingBottom: '0px' }"
      :title="articleData.title">
      <template #breadcrumb>
        <a-breadcrumb>
          <router-link to="/"><a-breadcrumb-item>首页</a-breadcrumb-item></router-link>
          <router-link :to="`/article/${articleData.type.flag}`"><a-breadcrumb-item>文章列表</a-breadcrumb-item></router-link>
          <a-breadcrumb-item>文章详情</a-breadcrumb-item>
        </a-breadcrumb>
      </template>
      <template #subtitle>
        <a-space>
          <span>作者: {{ articleData.user.name }}</span>
          <span>时间: {{ articleData.updatedAt }}</span>
        </a-space>
      </template>
      <template #extra>
        <a-space>
          <AButtonGroup>
            <AButton @click="like(articleData.id)">
              <icon-thumb-up-fill  size="large"  v-if="likeState" />
              <icon-thumb-up  size="large" v-else />{{articleData.like}}
            </AButton>

          </AButtonGroup>
          <a-dropdown-button>
            订阅
            <template #icon>
              <icon-down />
            </template>
            <template #content>
              <a-doption @click="subscribe(1,articleData.id)">{{articleSubscribe}}</a-doption>
              <a-doption @click="subscribe(2,articleData.user.id)">{{userSubscribe}}</a-doption>
            </template>
          </a-dropdown-button>

        </a-space>
      </template>
      <div>
         <a-tag color=blue size="large">
                {{ articleData.type.title || "未知" }}
              </a-tag>
              <a-divider direction="vertical"/>
        <a-space>
          <a-tooltip v-for="tagItem in articleData.tags" :key="tagItem.name" :content="tagItem.description"><a-tag
              color="red">{{ tagItem.name }}</a-tag></a-tooltip>
        </a-space>
      </div>
    </a-page-header>
    <div id="markdown-container"></div>
    <comment-edit :article-id="articleData.id"/>
    <comment-list style="margin-top: 10px;" :article-id="articleData.id"/>
  </div>
</template>

<script setup>
import { IconThumbUp,IconThumbUpFill, IconDown } from '@arco-design/web-vue/es/icon';
import {apiArticleLike, apiArticleLikeState, apiArticleView} from '@/apis/article';
import CommentEdit from '@/components/comment/CommentEdit.vue';
import CommentList from '@/components/comment/CommentList.vue';
import Cherry from 'cherry-markdown';
import 'cherry-markdown/dist/cherry-markdown.css';
import { nextTick, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import {IconNotification} from "@arco-design/web-vue/es/icon/index.js";
import {apiSubscribe, apiSubscribeState,} from "@/apis/apiSubscribe.js";
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

const articleSubscribe = ref("订阅文章")

const userSubscribe = ref("订阅用户")

onMounted(() => {
  apiArticleView(route.params.id).then(({ data }) => {
    articleData.value = data
    nextTick(() => {
      cherryConfig.value = articleData.value.content
      new Cherry(cherryConfig);
    })
    const id = articleData.value.id
    apiArticleLikeState(id).then((data)=>{
      likeState.value = data.data
    })
    apiSubscribeState(1,id).then((data)=>{
      articleSubscribe.value = data.data ? "文章已订阅" : "订阅文章"
    })
    apiSubscribeState(2,id).then((data)=>{
      userSubscribe.value = data.data ? "用户已订阅" : "订阅用户"
    })
  })

})

function like(id){
  apiArticleLike(id).then((data)=>{
    articleData.value.like += (data.data ? 1 :-1)
    likeState.value = data.data
    console.log(articleData.value.like)
  })

}

function subscribe(eventId,businessId){
  apiSubscribe(eventId,businessId).then((data)=>{
    if (eventId == 1){
      articleSubscribe.value = data.data ? "文章已订阅" : "订阅文章"
    }else {
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
}</style>