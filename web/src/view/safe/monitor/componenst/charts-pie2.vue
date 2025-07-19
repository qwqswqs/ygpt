<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
    @desc: 圆环饼图统计图表
!-->

<template>
  <Chart :height="height" :option="chartOption" />
</template>

<script setup>
import Chart from "@/components/charts/index.vue";
import useChartOption from '@/hooks/charts';
import { ref } from 'vue';
import { storeToRefs } from "pinia";
import { useAppStore } from '@/pinia';

const appStore = useAppStore();
const { config } = storeToRefs(appStore);

const props = defineProps({
  height: {
    type: String,
    default: '128px',
  },
  data: {
    type: Array,
    default: () => []
  }
});

const { chartOption } = useChartOption(() => {
  return {
    series: [
      {
        type: 'pie',
        radius: '50%',
        data: props.data,
        label: {
          show: true,
          formatter: '{b}: {c}%', // 显示名称、具体数值和百分比
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };
});
</script>

<style scoped lang="scss">
</style>
