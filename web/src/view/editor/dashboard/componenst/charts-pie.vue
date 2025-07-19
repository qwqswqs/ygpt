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
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '数据',
        type: 'pie',
        radius: ['40%', '70%'],
        data: props.data.map(item => ({
          ...item,
          itemStyle: { color: getColor(item.value) }
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          normal: {
            formatter: '{b}: {c}',
            position: 'outside'
          }
        },
        labelLine: {
          normal: {
            length: 10,
            length2: 10,
            smooth: true,
            lineStyle: {
              width: 0
            }
          }
        },
        itemStyle: {
          borderWidth: 0
        }
      }
    ]
  };
});
</script>

<style scoped lang="scss">
</style>
