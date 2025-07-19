<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="供给" name="1"><Supply v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="需求" name="2"><Demand v-if="activeName === '2'"/></el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>

import Supply from "./supplyInfo/supplyInfo.vue"
import Demand from "./demandInfo/demandInfo.vue"
export default {
  name: "model",
  components:{
    Supply,
    Demand,
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
