<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="节点算法" name="1"><NodeAlgo v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="我的算法" name="2"><MyAlgo v-if="activeName === '2'"/></el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>

import NodeAlgo from "./component/nodeAlgo.vue"
import MyAlgo from "./component/myAlgo.vue"
export default {
  name: "model",
  components:{
    MyAlgo,
    NodeAlgo
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
