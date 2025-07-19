<template>
  <Chart :height="height" :option="chartOption"/>
</template>

<script setup>
import Chart from "@/components/charts/index.vue";
import useChartOption from '@/hooks/charts';
import {ref} from 'vue';
import {storeToRefs} from "pinia";
import {useAppStore} from '@/pinia';

const appStore = useAppStore();
const {config} = storeToRefs(appStore);

const props = defineProps({
  height: {
    type: String,
    default: '128px',
  },
  diskData: {
  }
});
const userTotalInfo = ref({})

const primaryColor = '#0065ff';
const secondaryColor = '#b0d0ff';

const {chartOption} = useChartOption(() => {
  return {
    title: {
      text: props.diskData.title, // 主标题
      x: 'left' // x轴方向对齐方式
    },
    grid: {containLabel: true},
    tooltip: {
      trigger: 'item'
    },
    // backgroundColor: '#ffffff',
    legend: {
      orient: 'vertical',
      icon: 'circle',
      align: 'left',
      x: 'right',
      y: 'bottom',
      data: props.diskData.name
    },
    series: [
      {
        name: props.diskData.title,
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        center: ['50%', '50%'],
        itemStyle: {
          emphasis: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          },
          color: function (params) {
            // 自定义颜色
            var colorList = props.diskData.color
            return colorList[params.dataIndex]
          }
        },
        label: {
          normal: {
            show: true,
            position: 'center',
            color: '#4c4a4a',
            formatter: '{total|' + props.diskData.total + '}' + '\n\r' + '{active|' + props.diskData.totalText + '}',
            rich: {
              total: {
                fontSize: 35,
                fontFamily: '微软雅黑',
                color: '#454c5c'
              },
              active: {
                fontFamily: '微软雅黑',
                fontSize: 16,
                color: '#6c7a89',
                lineHeight: 30
              }
            }
          },
          emphasis: { // 中间文字显示
            show: true
          }
        },
        data: props.diskData.value
      }
    ]
  }
});

</script>
