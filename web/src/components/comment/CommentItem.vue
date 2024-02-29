<template>
    <a-comment author="Socrates" content="Comment body content." datetime="1 hour">
        <template #actions>
            <span class="action" key="heart" @click="onLikeChange">
                <span v-if="like">
                    <IconHeartFill :style="{ color: '#f53f3f' }" />
                </span>
                <span v-else>
                    <IconHeart />
                </span>
                {{ 83 + (like ? 1 : 0) }}
            </span>
            <span class="action" key="reply">
                <IconMessage /> Reply
            </span>
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
import { IconHeart, IconHeartFill } from '@arco-design/web-vue/es/icon';
</script>