<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="节点模型" name="1"><NodeModel v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="我的模型" name="2"><MyModel v-if="activeName === '2'"/></el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>

import NodeModel from "./component/nodeModel.vue"
import MyModel from "./component/myModel.vue"
export default {
  name: "model",
  components:{
    MyModel,
    NodeModel
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
