<template>
  <div class="">
    <div class="flex items-center justify-between mb-2">
      <div  v-if="title" class="font-bold">
        {{ title }}
      </div>
      <slot v-else name="title" />
    </div>
    <div class=" w-full relative">
      <div v-if="type < 4">
        <div class="mt-4 text-gray-600 text-3xl font-mono">
          <el-statistic :value="dataNum" />
        </div>
      </div>
<!--      <div class=" absolute top-0 right-2 w-[50%] h-20">-->
<!--        <charts-people-number v-if="type === 1" :data="data[0]" height="100%" />-->
<!--        <charts-people-number v-if="type === 2" :data="data[1]" height="100%" />-->
<!--        <charts-people-number v-if="type === 3" :data="data[2]" height="100%" />-->
<!--      </div>-->
      <charts-content-number v-if="type === 4" height="14rem" />
      <charts-pie v-if="type === 5" :data="pieData" height="14rem" />
      <chartDash v-if="type === 6" :data="dashData" height="22rem" />
      <chartsPie2 v-if="type === 7" :data="pieData" height="28rem" />
      <chartLine v-if="type === 8" :title="lineTitle" :x-data="xLineData" :y-data="yLineData" :data-value="chartValue" height="22rem" />
    </div>
  </div>
</template>

<script setup>


import chartsPeopleNumber from "./charts-people-numbers.vue"
import chartsContentNumber from "./charts-content-numbers.vue"
import chartsPie from "./charts-pie.vue"
import chartsPie2 from "./charts-pie2.vue"
import chartDash from "./charts-dash.vue"
import chartLine from "./charts-line.vue"
defineProps({
  type :{
    type : Number,
    default : 1
  },
  dataNum:{
    type:Number,
    default:0
  },
  title : {
    type : String,
    default : ""
  },
  pieData: {
    type: Array,
    default: () => [
      { value: 80, name: '未使用' },
      { value: 20, name: '已使用' }
    ]
  },
  dashData:{
    type: Array,
    default: () => [
      { value: 80, name: '未使用' },
      { value: 20, name: '已使用' }
    ]
  },
  lineTitle: {
    type : String,
    default : ""
  },
  xLineData:{
    type: Array,
    default: () => [
      2,5
    ]
  },
  yLineData:{
    type: Array,
    default: () => [
      1,6
    ]
  },
  chartValue:{
    type : String,
    default : '%'
  }
})
</script>

<style scoped lang="scss">

</style>
