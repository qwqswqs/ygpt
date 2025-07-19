<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="instanceName" label="实例名字" show-overflow-tooltip />
      <el-table-column prop="category" label="类别" show-overflow-tooltip />
      <el-table-column prop="activationDate" label="启用日期" show-overflow-tooltip />
      <el-table-column prop="remainingTime" label="剩余时间" show-overflow-tooltip />
      <el-table-column prop="cpuUsage" label="CPU" show-overflow-tooltip />
      <el-table-column prop="gpuUsage" label="GPU" show-overflow-tooltip />
      <el-table-column prop="memoryUsage" label="内存使用率" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" show-overflow-tooltip />
    </el-table>
  </div>
</template>

<script setup>
import { ref } from 'vue';

// 示例数据
const tableData = ref([
  {
    instanceName: '实例1',
    category: '云主机',
    activationDate: '2024-05-08',
    remainingTime: '30天',
    cpuUsage: '50%',
    gpuUsage: '30%',
    memoryUsage: '60%',
    status: '运行'
  },
  {
    instanceName: '实例2',
    category: '容器',
    activationDate: '2024-05-10',
    remainingTime: '20天',
    cpuUsage: '70%',
    gpuUsage: '40%',
    memoryUsage: '80%',
    status: '运行'
  },
  // 更多数据...
]);
</script>

<style scoped lang="scss">
/* 根据需要添加样式 */
</style>
