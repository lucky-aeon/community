<template>
    <a-spin :loading="loading" tip="This may take a while..." style="width: 100%;height: 100%;">
        <div :id="`${markId}`"></div>
        <a-modal v-model:visible="searchAt.show" :ok-loading="loading"  v-on:ok="handlerAtOk">
    <template #title>
      请搜索要提起的人
    </template>
    <div>
        <a-auto-complete v-model:model-value="searchAt.select" :data="searchAt.data" :loading="loading" v-on:search="searchAt.search"/>
    </div>
  </a-modal>
    </a-spin>
</template>
<script setup>
import { apiGetFile, apiUploadFile } from '@/apis/file';
import { useUserStore } from '@/stores/UserStore';
import {apiSearchUserByName} from '@/apis/user'
import { Message } from '@arco-design/web-vue';
import Cherry from 'cherry-markdown';
import 'cherry-markdown/dist/cherry-markdown.css';
import { nextTick } from 'vue';
import { reactive } from 'vue';
import { onMounted, ref, watch } from 'vue';
const markId = ref("markdownedit-container" + (Math.random().toString()))
const userStore = useUserStore()
const loading = ref(false)
const searchAt = reactive({
    show: false,
    data: [],
    select: "",
    search(v){
        apiSearchUserByName(v).then(({data, ok})=> {
            if(!ok) return
            searchAt.data = data.map(item=> ({
                value: `@(${item.name})[${item.id}]`,
                label: item.name
            }))
        })
    }
})
const handlerAtOk = ()=> {
    searchAt.data = []
    model.value = model.value.substring(0, model.value.length-1).concat(searchAt.select)
    searchAt.select = ""
}
const props = defineProps({
    previewOnly: {
        type: Boolean,
        default: false
    },
    showNav: {
        type: Boolean,
        default: true
    }
})
/**
 * 自定义一个语法，识别形如 ***ABC*** 的内容，并将其替换成 <span style="color: red"><strong>ABC</strong></span>
 */
var CustomHookA = Cherry.createSyntaxHook('important', Cherry.constants.HOOKS_TYPE_LIST.SEN, {
    makeHtml(str) {
        return str.replace(this.RULE.reg, function (whole, m1, m2) {

            return `<a class="chip" href="/user/${m2}" target="_blank">@${m1}</a>`
            // h(
            //     "a-tag",
            //     {
            //         color: "blue",
            //         icon: "icon-user"
            //     },
            //     {
            //         m1
            //     }
            // );
        });
    },
    rule(str) {
        return { reg: /@\((.*?)\)\[(.*?)\]/g };
    },
});
var cherryInstance = undefined
var cherryConfig = {
    id: markId.value,
    value: '# welcome to cherry editor! \n awdwadad',
    externals: {
        echarts: window.echarts,
        katex: window.katex,
        MathJax: window.MathJax,
    },
    engine: {
        global: {
            urlProcessor(url, srcType) {
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
                force: true,
            },
        },
    },
    toolbars: {
        toc: (props.showNav ? {
            updateLocationHash: false, // 要不要更新URL的hash
            defaultModel: 'full', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
        } : undefined),
        toolbarRight: ['fullScreen', '|', 'theme'],

    },
    editor: {
        defaultModel: 'editOnly',
    },
    callback: {
        onClickPreview: function (e) {
            const { target } = e;
            if (target.tagName === 'IMG') {
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
        },
        afterChange(e, b, c) {
            model.value = e
            if (e[e.length - 1] == '@') {
                nextTick(() => { 
                    searchAt.show = true
                })
            } else {
                searchAt.show = false
            }
        }
    },
    previewer: {
        // 自定义markdown预览区域class
        // className: 'markdown'
    },
    keydown: [],
    //extensions: [],
    async fileUpload(file, callback) {
        loading.value = true
        await apiUploadFile(userStore.userInfo.id, file, (data, key) => {
            if (!data.ok) {
                Message.error(data.msg)
                return
            }
            loading.value = false
            callback(apiGetFile(key))
        })
    }
};
const [model] = defineModel({
    default: "",
    required: true,
})
function refreshMarkdown(preview) {
    if (!cherryInstance) return
    if (preview) {
        props.previewOnly ? cherryInstance.switchModel("previewOnly") : cherryInstance.switchModel("editOnly")
        return
    }
    if (cherryInstance.getValue() == model.value) return
    cherryInstance.setValue(model.value, true)
}
onMounted(() => {

    cherryInstance = new Cherry(cherryConfig);
    refreshMarkdown()
})

watch(() => props.previewOnly, () => {
    refreshMarkdown(true)
}, { immediate: true, deep: true })
watch(() => props.showNav, () => {
    refreshMarkdown()
}, { immediate: true, deep: true })
watch(model, () => {
    refreshMarkdown()
})


</script>

<style>
.cherry {
    box-shadow: none;
    width: 100%;
    height: 100%;
}

.cherry-previewer {
    background-color: white;
    padding-top: 5px;
}

.cherry-flex-toc {
    position: fixed
}

.cherry video {
    width: 50%;
}

.chip {
    display: inline-block;
    background-color: #a4c3f9bd;
    color: #4f1afc;
    border-radius: 5px;
    margin-right: 4px;
}

.chip:hover {
    background-color: #86aff5bd;
    cursor: pointer;
}
</style>