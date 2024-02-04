<template>
    <div style="border-bottom: 1px solid rgba(215, 215, 215, 0.784); padding-bottom: 20px;">
        <ARow justify="center" :gutter="[20,20]">
            <ACol class="search">
                <ARow :gutter="[0,10]">
                    <ACol>
                        <h1>搜索全站文章</h1>
                        <a-auto-complete :data="data" @search="handleSearch" :style="StyleSearch"
                            placeholder="please enter something">
                            <template #prefix>
                                <icon-search />
                            </template>
                            <template #footer>
                                <div style="padding: 6px 0; text-align: center;">
                                    <a-button>Click Me</a-button>
                                </div>
                            </template>
                        </a-auto-complete>
                    </ACol>
                    <ACol>
                        <a-space>
                            <a-tag color="gray">
                                <template #icon>
                                    <icon-github />
                                </template>
                                Github
                            </a-tag>
                            <a-tag color="orangered">
                                <template #icon>
                                    <icon-gitlab />
                                </template>
                                Gitlab
                            </a-tag>
                            <a-tag color="blue">
                                <template #icon>
                                    <icon-twitter />
                                </template>
                                Twitter
                            </a-tag>
                            <a-tag color="arcoblue">
                                <template #icon>
                                    <icon-facebook />
                                </template>
                                Facebook
                            </a-tag>
                        </a-space>
                    </ACol>
                </ARow>
            </ACol>
        </ARow>
    </div>
    <router-view></router-view>
</template>
<script setup>
import { reactive } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
// search input style
const StyleSearch = reactive({
    width: '360px',
    height: '40px',
    innerHeight: '50px'
})


</script>

<style scoped>
.search {
    width: 100%;
    text-align: center; 
}
</style>