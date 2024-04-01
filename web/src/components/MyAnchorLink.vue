<template>
   <a-anchor-link style="width: 100%;" :href="`#${data.id}`" :key="data.id">{{ data.text }}
    <template v-if="data.children && data.children.length>0" #sublist>
        <MyAnchorLink :data="item" v-for="item in data.children"/>
    </template></a-anchor-link>
    
</template>


<script setup>
import MyAnchorLink from './MyAnchorLink.vue'
const {data} = defineProps({
    data: {
        type: Object,
        required: true
    }
})
</script>