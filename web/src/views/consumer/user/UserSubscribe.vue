<template>
    <a-table :columns="columns" :data="dataSources" :pagination="pagination">
        <template #optional="{record }">
            <a-popconfirm content="确定要取消该订阅?" @ok="subscribe(record.eventId, record.businessId)">
                <a-button>取消订阅</a-button>
            </a-popconfirm>
        </template>
    </a-table>
</template>

<script setup>
import * as apiSubscribe from '@/apis/subscribe';
import { ref } from 'vue';
const columns = [
    {
        title: '订阅事件',
        dataIndex: 'eventName',
    },
    {
        title: '订阅对象',
        dataIndex: 'businessName',
    },
    {
        title: '订阅时间',
        dataIndex: 'createdAt',
    },
    {
        title: '操作',
        slotName: 'optional'
    }
];
const dataSources = ref([]);
const pagination = ref({
    total: 0,
    current: 1,
    defaultPageSize: 10
})
function refreshList() {
    apiSubscribe.getList(pagination.value).then(({ data, ok }) => {
        if (!ok) return
        dataSources.value = data
        pagination.value.total = data.count
    })
}
refreshList()

function subscribe(eventId, businessId) {
    apiSubscribe.subscribe(eventId, businessId).then(({ok})=>{
        if(!ok) return
        refreshList()
    })
}
</script>
