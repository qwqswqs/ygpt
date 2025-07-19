<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
    @desc: 曲线统计图表
!-->

<template>
  <Chart :height="height" :option="chartOption" />
</template>

<script setup>
import Chart from "@/components/charts/index.vue";
import useChartOption from '@/hooks/charts';

const props = defineProps({
  height: {
    type: String,
    default: '128px',
  },
  title: {
    type: String,
    default: '标题',
  },
  xData: {
    type: Array,
    default: () => []
  },
  yData: {
    type: Array,
    default: () => []
  },
  dataValue:{
    type:String,
    default:'%',
  }
});

const { chartOption } = useChartOption(() => {
  return {
    grid: {
      left: 10,
      top: 50,
      right: 10,
      bottom: 0,
      containLabel: true,
    },
    legend: {
      right: 5,
      type: 'scroll',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        lineStyle: {
          color: '#ddd',
        },
      },
      backgroundColor: 'rgba(255,255,255,1)',
      padding: [5, 10],
      textStyle: {
        color: '#7588E4',
      },
      extraCssText: 'box-shadow: 0 0 5px rgba(0,0,0,0.3)',
      formatter: function(params) { // 自定义tooltip内容
        let tooltipText = '';
        params.forEach(item => {
          tooltipText += `<div style="margin-bottom: 5px;">
                            时间: ${item.name}<br/>
                            ${item.seriesName}: ${item.value}${props.dataValue}
                          </div>`;
        });
        return tooltipText;
      }
    },
    xAxis: {
      type: 'category',
      data: props.xData,
      splitLine: {
        show: true,
        interval: 'auto',
        lineStyle: {
          color: ['#D4DFF5'],
        },
      },
      axisTick: {
        show: false,
      },
      axisLine: {
        lineStyle: {
          color: '#999',
        },
      },
      axisLabel: {
        showMaxLabel: false,
        margin: 10,
        align: 'left',
        textStyle: {
          fontSize: 12,
        },
      },
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          color: ['#D4DFF5'],
        },
      },
      axisTick: {
        show: false,
      },
      axisLine: {
        lineStyle: {
          color: '#999',
        },
      },
      axisLabel: {
        margin: 10,
        textStyle: {
          fontSize: 12,
        },
        formatter: '{value}'+props.dataValue,
      },
    },
    series: props.yData,
  };
});
</script>

<style scoped lang="scss">
</style>
