
<template>
    <a-comment author="Socrates" datetime="发表评论">
        <template #content>
            <MarkdownEdit :show-nav="false" style="margin-top: 5px;height: 100%;" v-model="commentData"/>
        </template>
        <template #avatar>
            <a-avatar>
                <img alt="avatar"
                    src="https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp" />
            </a-avatar>
        </template>
    </a-comment>
</template>
<script setup>
import { ref } from 'vue';
import MarkdownEdit from '../MarkdownEdit.vue';
const props = defineProps({
    articleId: {
        type: Number,
        default: 0,
    }
})
const commentData = ref("")

</script>