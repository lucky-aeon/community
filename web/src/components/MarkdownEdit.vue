<template>
    <div id="markdownedit-container"></div>
</template>
<script setup>
import Cherry from 'cherry-markdown';
import 'cherry-markdown/dist/cherry-markdown.css';
import { onMounted, watch } from 'vue';
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
var cherryInstance = undefined
var cherryConfig = {
    id: 'markdownedit-container',
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
        toc: (props.showNav ? {
            updateLocationHash: true, // 要不要更新URL的hash
            defaultModel: 'full', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
        } : undefined),
    },
    editor: {
        defaultModel: 'editOnly',
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
        },
        afterChange(e) {
            model.value = e
        }
    },
    previewer: {
        // 自定义markdown预览区域class
        // className: 'markdown'
    },
    keydown: [],
    //extensions: [],
};
const [model] = defineModel({
    default: "",
    required: true,
})
function refreshMarkdown() {
    if (!cherryInstance) return
    cherryInstance.setValue(model.value)
    props.previewOnly ? cherryInstance.switchModel("previewOnly") : cherryInstance.switchModel("editOnly")
}
onMounted(() => {

    cherryInstance = new Cherry(cherryConfig);
    refreshMarkdown()
})

watch(() => props.previewOnly, (n) => {
    console.log(n)
    refreshMarkdown()
}, { immediate: true, deep: true })
watch(() => props.showNav, (n) => {
    console.log(n)
    refreshMarkdown()
}, { immediate: true, deep: true })
</script>

<style>
.cherry {
    box-shadow: none;
    width: 100%;
    height: 100%;
}

.cherry-flex-toc {
    position: fixed
}
</style>