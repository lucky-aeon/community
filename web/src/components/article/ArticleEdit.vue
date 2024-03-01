<template>
    <a-modal :mask-closable="false" :fullscreen="fullScreen" :body-style="{ height: '100%' }" v-model:visible="visible"
        @ok="handleOk" @cancel="handleCancel" draggable
        :modal-style="{ minWidth: '800px', maxHeight: fullScreen ? '' : '90%' }">
        <template #title>
            {{ modalTitile }}
        </template>
        <template #footer>
            <a-row>
                <a-col flex="0">
                    <a-button-group>
                        <a-button @click="fullScreen = !fullScreen">{{fullScreen?"窗口":"全屏"}}</a-button>
                        <a-button @click="() => previewOnly = !previewOnly">{{ previewOnly ? "纯编辑" : "纯预览" }}</a-button>
                    </a-button-group>
                </a-col>
                <a-col flex="auto">
                </a-col>
                <a-col flex="100px">
                    <a-button-group type="primary">
                        <a-button>发布</a-button>
                        <a-dropdown @select="handleSelect" :popup-max-height="false">
                            <a-button>
                                <template #icon>
                                    <icon-down />
                                </template>
                            </a-button>
                            <template #content>
                                <a-doption>保存草稿</a-doption>
                            </template>
                        </a-dropdown>
                    </a-button-group>

                </a-col>
            </a-row>

        </template>
        <div>
            <a-input-group style="width: 100%;">
                <a-auto-complete placeholder="分类" style="width: 150px;" />
                <a-input placeholder="标题" style="width: 100%;"></a-input>
            </a-input-group>
            <div style="margin-top: 5px;"></div>
            <tag-search/>
            <div style="margin-top: 5px;"></div>
            <a-scrollbar :style="`height: ${fullScreen ? '' : '350px'};overflow: auto;`">
                <markdown-edit v-model="data.desc" :preview-only="previewOnly" />
            </a-scrollbar>
        </div>
    </a-modal>
</template>
<script setup>
import { IconDown } from '@arco-design/web-vue/es/icon';
import { computed, ref } from 'vue';
import MarkdownEdit from '../MarkdownEdit.vue';
import TagSearch from '../TagSearch.vue';
const { articleId } = defineProps({
    articleId: {
        type: Number,
        default: 0
    }
})
const data = ref({ desc: "awdwad" })
const visible = ref(true);
const previewOnly = ref(false);
const fullScreen = ref(false);
const modalTitile = computed(() => articleId > 0 ? "编辑文章" : "添加文章")
const handleOk = () => {
    visible.value = false;
};
const handleCancel = () => {
    visible.value = false;
}
</script>