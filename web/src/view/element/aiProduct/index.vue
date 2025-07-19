<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="模型" name="1"><Model v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="数据集" name="2"><Dataset v-if="activeName === '2'"/></el-tab-pane>
<!--      <el-tab-pane label="算法" name="3"><Algo v-if="activeName === '3'"/></el-tab-pane>-->
    </el-tabs>
  </div>
</template>
<script>

import Model from "./model/model.vue"
import Dataset from "./dataset/dataset.vue"
import Algo from "./algo/algo.vue"
export default {
  name: "model",
  components:{
    Model,
    Algo,
    Dataset
  },
  data(){
    return{
      activeName: '1',
    }
  },
  methods:{
    handleClick(tab,event){
      console.log(tab,event);
    }
  }
}
</script>
<style scoped>
.app-container {
  padding: 20px;
}

</style>
