<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 py-2 gap-4 md:gap-2 gva-container2">
    <div class="col-span-1 lg:col-span-6 flex flex-wrap gap-4">
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="1" title="当前租户量"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="2" title="物理机数量"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="3" title="当前实例总数"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="2" title="算力产品总数" />
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="2" title="模型总数"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42">
        <GvaChart :data-num="12" :type="3" title="数据集总数"/>
      </GvaCard>
    </div>
       
    <GvaCard custom-class="col-span-1 md:col-span-2 md:row-start-2 lg:col-span-6 col-start-1 row-span-2"
              title="最近运行实例折线趋势图">
      <GvaChart :type="4"/>
    </GvaCard>
    <div class="col-span-1 lg:col-span-4 flex flex-wrap gap-4">
      <GvaCard custom-class="flex-1 h-42" title="云主机使用情况">
        <GvaChart :pie-data="hostData" :type="5"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42" title="裸金属使用情况">
        <GvaChart :pie-data="bareData" :type="5"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42" title="容器使用情况">
        <GvaChart :pie-data="containData" :type="5"/>
      </GvaCard>
    </div>
    <div class="col-span-1 lg:col-span-4 flex flex-wrap gap-4">
      <GvaCard custom-class="flex-1 h-42" title="当前资源使用率">
        <GvaChart :dash-data="resourceUse" :type="6"/>
      </GvaCard>
      <GvaCard custom-class="flex-1 h-42" title="运行实例状态">
        <GvaChart :pie-data="runStatus" :type="7"/>
      </GvaCard>
    </div>
  </div>
</template>

<script>
import {GvaChart, GvaCard} from "./componenst"
import * as echarts from 'echarts';
export default {
  components:{
    GvaChart,
    GvaCard
  },
  data() {
    return {
      renterNum:0,//租户数量
      serverNum:0,//物理机数量
      instanceNum:0,//实例数量
      productNum:0,//算力产品数量
      modelNum:0,//模型数量
      datasetNum:0,//数据集数量
      //云主机使用情况
      hostData:[
        {value: 12, name: '已使用'},
        {value: 5, name: '未使用'}
      ],
      //裸金属使用情况
      bareData:[
        {value: 3, name: '已使用'},
        {value: 8, name: '未使用'}
      ],
      //容器使用情况
      containData:[
        {value: 2, name: '已使用'},
        {value: 5, name: '未使用'}
      ],
      //资源使用率
      resourceUse:[{value: 12, name: '使用率'},],
      //故障、空闲、运行
      runStatus:[
        {value: 2, name: '故障',itemStyle: { color: '#f44336' }},
        {value: 2, name: '空闲',itemStyle: { color: '#ffa726' }},
        {value: 5, name: '运行',itemStyle: { color: '#4caf50' }}
      ],
    }
  },
}
</script>

<style lang="scss" scoped>
.chart-container {
  height: 300px;
}
</style>
