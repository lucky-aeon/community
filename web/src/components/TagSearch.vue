<template>
    <a-auto-complete :data="data" @search="handleSearch" :style="{width:'360px'}" placeholder="please enter something"/>
</template>

<script setup>
import { ref } from 'vue';
defineProps({
    show:{
        typeof: 'boolean',
        default: false
    }
})
const data = ref([])
const handleSearch = (e)=>{
    console.log("awdad", e)
}
</script>