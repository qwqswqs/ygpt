<template>
  <div class="app-container">
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="镜像" name="1"><Image v-if="activeName === '1'"/></el-tab-pane>
      <el-tab-pane label="快照" name="2"><Snapshot v-if="activeName === '2'"/></el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>

import Image from "./component/image.vue"
import Snapshot from "./component/snapshot.vue"
export default {
  name: "model",
  components:{
    Snapshot,
    Image
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
