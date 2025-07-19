<template>
  <div style="">
  <div style="background-color: white;margin-bottom: 16px;padding: 24px;margin-top: 16px;">
    <span>资源概况</span>
    <div class="grid grid-cols-1 md:grid-cols-2 grid5 py-2 md:gap-2 gva-container2 flex flex-wrap" style="margin-top: 16px">
      <template v-for="(item,key) in requestData" :key="key">
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'host_num'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key]" :unit="'个'" :icons=Suzhuji title="宿主机数量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'baremetal_num'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key] " :unit="' 个'" :icons=Wuliji title="物理机数量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'cpu'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="'核'" :icons=CPU title="CPU"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'ip'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="' 个'" :icons=IP title="IP总量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'vm'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="'台 '" :icons=Yunzhuji title="云主机数量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'memory'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="'GB '"  :icons=Neicun title="内存总量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'storage'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="'GB '" :border="'is'" :icons=Cunchu title="存储总量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'gpu'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total" :unit="'卡 '" :icons=GPU title="GPU总量"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1  gap-4" v-if="key == 'vpc'">
          <gva-card custom-class="flex-1 h-42">
            <gva-chart :type="6" :statistic="requestData[key].total " :unit="' 个'" :icons=VPC title="VPC总量"/>
          </gva-card>
        </div>
        <template v-if="key == 'k8s'">
          <template v-for="(val,i) in requestData[key]" :key="i">
            <div class="col-span-1 lg:col-span-1  gap-4" v-if="i == 'cluster_num'">
              <gva-card custom-class="flex-1 h-42">
                <gva-chart :type="6" :statistic="requestData[key][i] " :unit="' 个'" :border="'is'" :icons=K8S title="k8s集群数量"/>
              </gva-card>
            </div>
          </template>
        </template>

      </template>
    </div>
  </div>

    <!-- <gva-card custom-class="col-span-1 md:col-span-2 md:row-start-2 lg:col-span-6 col-start-1 row-span-2"
              title="资源购买趋势图">
      <gva-chart :type="4"/>
    </gva-card> -->
  <div style="background-color: white;padding: 24px;margin-top: 16px;">
    <span>资源使用情况</span>
    <div class="grid grid-cols-1 md:grid-cols-2 grid-4 py-2 gap-4 md:gap-2 gva-container2 flex flex-wrap">
      <!-- <template v-for="(item,key) in requestData" :key="key">
        <div class="col-span-1 lg:col-span-1">
          <gva-card custom-class="flex-1 h-42" :title="key+'使用率'">
            <gva-chart :type="5"/>
          </gva-card>
        </div>
      </template> -->
      <template v-for="(item,key) in requestData" :key="key">
        <div class="col-span-1 lg:col-span-1" v-if="key == 'cpu'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart title="CPU使用率" :type="5" :unit="'核'" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'vm'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart title="云主机关机率" :type="5" :unit="'台 '" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'memory'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart :title="'内存使用率'" :type="5" :unit="'GB '" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'storage'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart :title="'存储使用率'" :type="5" :unit="'GB '" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'gpu'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart :title="'GPU使用率'" :unit="'卡 '" :type="5" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'ip'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart :title="'IP使用率'" :type="5" :unit="'个 '" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <div class="col-span-1 lg:col-span-1" v-if="key == 'vpc'">
          <gva-card custom-class="flex-1 h-42" >
            <gva-chart :title="'VPC使用率'" :type="5" :unit="'个 '" :pieData="requestData[key].dataArr"/>
          </gva-card>
        </div>
        <template v-if="key == 'k8s'">
          <template v-for="(val,i) in requestData[key]" :key="i">
            <div class="col-span-1 lg:col-span-1  gap-4" v-if="i == 'cpu'">
              <gva-card custom-class="flex-1 h-42">
                <gva-chart :type="5" :statistic="requestData[key][i].total +' GB'" title="k8s集群cpu资源" :unit="'GB '" :pieData="requestData[key][i].dataArr"/>
              </gva-card>
            </div>
            <div class="col-span-1 lg:col-span-1  gap-4" v-if="i == 'memory'">
              <gva-card custom-class="flex-1 h-42">
                <gva-chart :type="5" :statistic="requestData[key][i].total +' GB'" title="k8s集群的内存资源" :unit="'GB '" :pieData="requestData[key][i].dataArr"/>
              </gva-card>
            </div>
            <div class="col-span-1 lg:col-span-1  gap-4" v-if="i == 'node'">
              <gva-card custom-class="flex-1 h-42">
                <gva-chart :type="5" :statistic="requestData[key][i].total +' GB'" title="k8s集群节点" :unit="'个 '" :pieData="requestData[key][i].dataArr"/>
              </gva-card>
            </div>
          </template>
        </template>
      </template>

    </div>
  </div>

    
  
</div>
</template>

<script setup>
import {GvaChart, GvaCard} from "./componenst";
import { reactive, ref } from 'vue';
import {getcloudpodsStatistic} from "@/api/cloudpods/baseRes/device";
import Wuliji from '@/assets/wuliji.png';
import Suzhuji from '@/assets/suzhuji.png';
import CPU from '@/assets/cpu.png';
import IP from '@/assets/ip.png';
import Yunzhuji from '@/assets/yunzhuji.png';
import Neicun from '@/assets/neicun.png';
import Cunchu from '@/assets/cunchu.png';
import GPU from '@/assets/gpu.png';
import VPC from '@/assets/vpc.png';
import K8S from '@/assets/k8s.png';
const requestData = ref();

// 获取统计数据
const getcloudpods=()=> {
  getcloudpodsStatistic().then(response => {
    if (response.code == 0) {
      let dataObj = response.data;
      for (let key in dataObj) {
        console.log(key)
        let dataArr = [];
        if(key != 'host_num' && key != 'baremetal_num' && key != 'k8s'){
          let obj1 = {},obj2= {};
          if(dataObj[key].ready){
            obj1 = { value: dataObj[key].ready, name: '正常' };
            obj2 = { value: dataObj[key].no_ready, name: '关闭' };
          }else{
            obj1 = { value: dataObj[key].not_used, name: '未使用' };
            obj2 = { value: dataObj[key].used, name: '已使用' };
          }
          dataArr.push(obj1);
          dataArr.push(obj2);
          dataObj[key].dataArr = dataArr;
        }
        if(key == 'k8s'){
          for (let InnerKey in dataObj[key]) {
            let InnerArr = [];
            if(InnerKey != 'cluster_num'){
              let obj1 = {},obj2= {};
              if(dataObj[key][InnerKey].ready!=undefined){
                obj1 = { value: dataObj[key][InnerKey].ready, name: '正常' };
                obj2 = { value: dataObj[key][InnerKey].no_ready, name: '关闭' };
              }else{
                obj1 = { value: dataObj[key][InnerKey].not_used, name: '未使用' };
                obj2 = { value: dataObj[key][InnerKey].used, name: '已使用' };
              }
              InnerArr.push(obj1);
              InnerArr.push(obj2);
              dataObj[key][InnerKey].dataArr = InnerArr;
            }
          }
        }
      }
      requestData.value = dataObj;
    } else {
      this.$message({
        message: '获取失败',
        type: 'error'
      });
    }
  })
}
getcloudpods();
</script>

<style lang="scss" scoped>
sapn{
  font-size: 16px;
  font-weight: 500;
  color: rgba(29, 33, 41, 1);
}

.grid-4{
  grid-template-columns: repeat(4, minmax(0, 1fr)) !important;
  @media (max-width: 1540px) {
    grid-template-columns: repeat(3, minmax(0, 1fr)) !important;
  }
  @media (max-width: 1520px) {
    grid-template-columns: repeat(2, minmax(0, 1fr)) !important;
  }
}

.grid5{
  grid-template-columns: repeat(5, minmax(0, 1fr)) !important;
  @media (max-width: 1540px) {
    grid-template-columns: repeat(4, minmax(0, 1fr)) !important;
  }
  @media (max-width: 1320px) {
    grid-template-columns: repeat(2, minmax(0, 1fr)) !important;
  }
}
</style>
