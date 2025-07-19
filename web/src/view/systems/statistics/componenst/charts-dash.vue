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


const primaryColor = '#0065ff';
const secondaryColor = '#b0d0ff';


const getColor = (value) => {

  const maxValue = Math.max(...props.data.map(item => item.value));
  const minValue = Math.min(...props.data.map(item => item.value));
  const threshold = (maxValue + minValue) / 2;
  return value > threshold ? primaryColor : secondaryColor;
};

const { chartOption } = useChartOption(() => {
  return {
    series: [
      {
        type: 'gauge',
        detail: {formatter: '{value}%'},
        data: props.data,
        axisLine: {
          lineStyle: {
            color: [[0.3, '#f44336'], [0.7, '#ffa726'], [1, '#4caf50']],
            width: 30,
          },
        },
      },
    ],
  };
});
</script>

<style scoped lang="scss">
</style>
