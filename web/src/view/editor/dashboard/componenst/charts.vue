<template>
  <div class="">
    <div class="flex items-center justify-between mb-2">
      <div v-if="title" class="font-bold">
        {{ title }}
      </div>
      <slot v-else name="title" />
    </div>
    <div class=" w-full relative">
      <div v-if="type < 4">
        <div class="mt-4 text-gray-600 text-3xl font-mono">
          <el-statistic :value="268500" />
        </div>
        <div class="mt-2 text-green-600 text-sm font-bold font-mono">
          +80% <el-icon><TopRight /></el-icon>
        </div>
      </div>
      <div class=" absolute top-0 right-2 w-[50%] h-20">
        <charts-people-number v-if="type === 1" :data="data[0]" height="100%" />
        <charts-people-number v-if="type === 2" :data="data[1]" height="100%" />
        <charts-people-number v-if="type === 3" :data="data[2]" height="100%" />
      </div>
      <charts-content-number v-if="type === 4" height="14rem" />
      <charts-pie v-if="type === 5" :data="pieData" height="14rem" />
      <charts-pie v-if="type === 6" :data="pieData2" height="14rem" />

    </div>
  </div>
</template>

<script setup>


import chartsPeopleNumber from "./charts-people-numbers.vue"
import chartsContentNumber from "./charts-content-numbers.vue"
import chartsPie from "./charts-pie.vue"
defineProps({
  type :{
    type : Number,
    default : 1
  },
  title : {
    type : String,
    default : ""
  },
  pieData: {
    type: Array,
    default: () => [
      { value: 10, name: '运行中' },
      { value: 20, name: '未运行' }
    ]
  },
  pieData2: {
    type: Array,
    default: () => [
      { value: 20, name: '已使用' },
      { value: 80, name: '未使用' }
    ]
  }
})

const data = [
    [12, 22, 32, 45, 32, 78, 89, 92],
    [1, 2, 43, 5, 67, 78, 89, 12],
    [12, 22, 32, 45, 32, 78, 89, 92],
]
</script>

<style scoped lang="scss">

</style>
