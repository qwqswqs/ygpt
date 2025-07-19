<template>
  <div class="">
    <div class="flex ">
      <div v-if="type == 6">
        <img style="width: 60px;height: 60px;margin-right: 12px" :src="icons"></img>
      </div>
      <div class=" mb-2" style="width: 100px;">
        <div v-if="type == 6" class="font-bold">
          {{ title }}
        </div>
        <div v-if="type == 6">
          <div class="mt-4 text-gray-600 text-3xl font-mono">
            <span class="statistic" style="font-size: 32px;font-weight: 900;color: rgba(29, 33, 41, 1);margin-right: 4px">{{statistic}}</span>
            <span class="unit" style="font-size: 14px;font-weight: 400;color: rgba(78, 89, 105, 1);">{{ unit }}</span>
          </div>
        </div>
        <slot v-else name="title" />
      </div>
      <div v-if="type == 6 &&border!=='is'" style="height: 60px;border: 1px solid rgba(242, 243, 245, 1);margin-left: 60px"></div>
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
      <charts-pie v-if="type === 5" :data="pieData" :title="title" :unit="unit" height="18rem" />

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
  statistic: {
    type : String,
    default : "0"
  },
  icons: {
    type : String,
    default : ""
  },
  title : {
    type : String,
    default : ""
  },
  unit :{
    type : String,
    default : ""
  },
  border :{
    type : String,
    default : ""
  },
  pieData: {
    type: Array,
    default: () => [
      { value: 80, name: '未使用' },
      { value: 20, name: '已使用' }
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
.statistic{
  @media (max-width: 1540px) {
    font-size: 22px !important;
  }
}
</style>
