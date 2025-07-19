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
    default: '400px'
  },
  data: {
    type: Array,
    default: () => []
  },
  title: {
    type: String,
    default: ''
  },
  unit: {
    type: String,
    default: ''
  }

});

// 自定义颜色渐变，使用115.41deg的线性渐变
const primaryColor = {
  type: 'linear',
  x: 0, y: 1, x2: 1, y2: 0, // 对应115.41deg的方向（左下到右上）
  colorStops: [
    { offset: 0, color: 'rgba(55, 226, 226, 1)' }, // 渐变起始颜色
    { offset: 1, color: 'rgba(255, 255, 255, 1)' } // 渐变结束颜色
  ],
  global: false
};

const secondaryColor = {
  type: 'linear',
  x: 0, y: 1, x2: 1, y2: 0, // 对应115.41deg的方向
  colorStops: [
    {offset: 0, color: 'rgba(0, 115, 255, 1)'}, // 另一种起始颜色
    {offset: 1, color: 'rgba(255, 255, 255, 1)'} // 渐变结束颜色
  ],
  global: false
};

const getColor = (value,index) => {
  const maxValue = Math.max(...props.data.map(item => item.value));
  const minValue = Math.min(...props.data.map(item => item.value));
  const threshold = (maxValue + minValue) / 2;

  return index ? secondaryColor : primaryColor;
};

const {chartOption} = useChartOption(() => {
  return {
    title: {
      text: autoWrapText(props.title),
      left: 'center',
      top: '40%',
      textStyle: {
        fontSize: 16,
        fontWeight: 500,
        color: 'rgba(29, 33, 41, 1)',
        fontFamily: 'Noto Sans SC'
      }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}' + props.unit+ '({d}%)',
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#eee',
      borderWidth: 1,
      textStyle: {
        color: '#333',
        fontFamily: 'Noto Sans SC'
      }
    },
    legend: {
      orient: 'horizontal',
      left: 'center',
      bottom: '0',
      itemGap: 20,
      textStyle: {
        fontSize: 14,
        color: '#333',
        fontFamily: 'Noto Sans SC'
      },
      icon: 'circle',
      itemWidth: 10,
      itemHeight: 10
    },
    series: [
      {
        name: '数据',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: true,
        label: {
          show: true,
          position: 'outside',
          formatter: '{per|{d}%}',
          rich: {
            b: {
              fontSize: 12,
              lineHeight: 20,
              color: '#666'
            },
            per: {
              fontSize: 14,
              lineHeight: 24,
              color: '#333'
            }
          },
          distanceToLabelLine: 10
        },
        data: props.data.map((item,index) => ({
          ...item,
          index,
          itemStyle: {color: getColor(item.value,index)}
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        labelLine: {
          normal: {
            length: 15,
            length2: 20,
            smooth: true,
            lineStyle: {
              width: 1
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

const autoWrapText=(text)=> {
    if (!text || text.length <= 9) return text;
    
    let result = '';
    for (let i = 0; i < text.length; i++) {
        result += text[i];
        if ((i + 1) % 5 === 0 && i !== text.length - 1) {
            result += '\n';
        }
    }
    return result;
}
</script>

<style scoped lang="scss">
/* 无需额外样式，通过 ECharts 配置实现 */
</style>