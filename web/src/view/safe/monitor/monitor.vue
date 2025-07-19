<template>
  <div>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <GvaCard custom-class="h-42" title="服务器CPU利用率">
        <GvaChart :dash-data="cpuUse" :type="6"/>
      </GvaCard>
      <GvaCard custom-class="h-42" title="服务器磁盘使用率">
        <GvaChart :dash-data="diskUse" :type="6"/>
      </GvaCard>
      <GvaCard custom-class="h-42" title="服务器内存占用率">
        <GvaChart :dash-data="memoryUse" :type="6"/>
      </GvaCard>
    </div>
    <el-card style="margin-top: 10px">
      <el-descriptions border :column="2" title="连接平台">
        <el-descriptions-item label="连接状态">
          <el-tag :type="isConnect?'success':'danger'">{{isConnect?'已连接':'未连接'}}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作">
          <el-button type="text">{{isConnect?'断开':'连接'}}</el-button>
          <el-button type="text" :disabled="isConnect">设置</el-button>
        </el-descriptions-item>
        <div v-if="isConnect">
          <el-descriptions-item label="平台IP">10.10.100.12</el-descriptions-item>
          <el-descriptions-item label="端口">8889</el-descriptions-item>
          <el-descriptions-item label="关联码">10086</el-descriptions-item>
          <el-descriptions-item label="是否自动连接">是</el-descriptions-item>
          <el-descriptions-item label="自动连接时间间隔">60s</el-descriptions-item>
        </div>
      </el-descriptions>
    </el-card>
    <el-card style="margin-top: 10px">
      <el-descriptions :column="1" title="运行中的实例">
        <el-descriptions-item>
          <!-- 表格展示 -->
          <el-table
              :data="dataList"
              row-key="ID"
          >
            <el-table-column label="访问地址">
              <template #default="scope">
                <span>{{ scope.row.url }}</span>
              </template>
            </el-table-column>
            <el-table-column label="SSH地址">
              <template #default="scope">
                <span>{{ scope.row.sshHost }}</span>
              </template>
            </el-table-column>
            <el-table-column label="SSH端口">
              <template #default="scope">
                <span>{{ scope.row.sshPort }}</span>
              </template>
            </el-table-column>
            <el-table-column label="SSH用户名">
              <template #default="scope">
                <span>{{ scope.row.sshUser }}</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchInfo.page"
                :page-size="searchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="searchInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script>
import {GvaChart, GvaCard} from "./componenst"
import * as echarts from 'echarts';
import {getInfoList} from "@/api/yunguan/resource/resInfo";
import {getSystemState} from "@/api/system/system";
export default {
  components:{
    GvaChart,
    GvaCard
  },
  data() {
    return {
      //CPU利用率
      cpuUse:[{value: 13, name: '利用率'},],
      //磁盘使用率
      diskUse:[{value: 15, name: '使用率'},],
      //内存占用率
      memoryUse:[{value: 56, name: '占用率'},],
      isConnect:true,
      searchInfo: {
        page: 1,
        pageSize: 10,
        total:0,
      },
      dataList:[],
      timer:'',
    }
  },
  mounted() {
    this.timer = setInterval(() => {
      setTimeout(() => {
        this.GetServerStatus() //调用接口的方法
      }, 0)
    }, 5000)
  },
  created() {
    this.GetInfoList()
    this.GetServerStatus()
  },
  methods:{
    //定时器
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetInfoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetInfoList()
    },
    // 获取实例列表
    GetInfoList() {
      getInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.searchInfo.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //获取服务器状态
    GetServerStatus(){
      getSystemState().then(res=>{
        if (res.code==0){
          let tempCpu = res.data.server.cpu.cpus
          this.memoryUse[0].value = res.data.server.ram.usedPercent
          this.diskUse[0].value = res.data.server.disk[0].usedPercent
          let totalCpu = tempCpu.reduce((sum, cpu) => sum + cpu, 0)
          this.cpuUse[0].value = (totalCpu / tempCpu.length).toFixed(1);

        }
      })
    }
  },
  beforeDestroy(){
    clearInterval(this.timer);
    this.timer = null;
  }
}
</script>

<style lang="scss" scoped>
.chart-container {
  height: 300px;
}
</style>
