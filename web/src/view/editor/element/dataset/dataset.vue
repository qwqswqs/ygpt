<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="节点数据集" name="1"><NodeDataset v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="我的数据集" name="2"><MyDataset v-if="activeName === '2'"/></el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>

import NodeDataset from "./component/nodeDataset.vue"
import MyDataset from "./component/myDataset.vue"
export default {
  name: "model",
  components:{
    MyDataset,
    NodeDataset
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
