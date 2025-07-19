<template>

</template>
<script setup>
import { getCloudPodsConsoleUrl } from '@/api/system/cloudPods';
import router from '@/router';
import { ElMessage } from 'element-plus';
import { onMounted } from 'vue';

onMounted(()=>{
    getCloudPodsConsoleUrl().then(res=>{
        window.open(res.data)
    }).catch((err)=>{
        ElMessage.error(err)
    }).finally(()=>{
        router.push('/')
    })
})
</script>

<style scoped>

</style>