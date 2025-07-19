<template>
  <div>
    <div class="gva-table-box">
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="logList"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column label="操作时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.ops_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作类型">
          <template #default="scope">
            <span>{{ actionType[scope.row.action] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="风险类型">
          <template #default="scope">
            <span>{{ scope.row.severity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="事件类型">
          <template #default="scope">
            <span>{{ scope.row.kind }}</span>
          </template>
        </el-table-column>
        <el-table-column label="结果">
          <template #default="scope">
            <span v-if="scope.row.success" style="color: green">{{ '成功' }}</span>
            <span v-if="!scope.row.success" style="color: red">{{ '失败' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="openDialog(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px">
      <el-button type="text" v-if="nextMarker>0" @click="GetLogsList">{{ '加载更多' }}</el-button>
      <span v-if="nextMarker==0">{{'没有更多了'}}</span>
      </div>
    </div>
    <el-dialog v-model="dialogVisible"  :show-close="false" :before-close="closeDialog" style="width: 520px;">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">
            日志查看
          </span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="dialogVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-descriptions style="margin: 10px 5px"  :column="1" label-width="100px">
        <el-descriptions-item label="ID">{{detailInfo.id}}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{FormatDateTime(detailInfo.start_time)}}</el-descriptions-item>
        <el-descriptions-item label="资源类型">{{detailInfo.obj_type}}</el-descriptions-item>
        <el-descriptions-item label="资源名称">{{detailInfo.obj_name}}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{actionType[detailInfo.action]}}</el-descriptions-item>
        <el-descriptions-item label="风险级别">{{detailInfo.severity}}</el-descriptions-item>
        <el-descriptions-item label="事件类型">{{detailInfo.kind}}</el-descriptions-item>
        <el-descriptions-item label="结果">{{detailInfo.success?'成功':'失败'}}</el-descriptions-item>
        <el-descriptions-item label="详情">
<!--          <div class="scrollable-content">{{ formattedJson(detailInfo.notes) }}</div>-->
          {{ formattedJson(detailInfo.notes) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>

import { getLogsList} from '@/api/cloudpods/common/log'


export default {
  props: {
    logData: {
      type: Object,
      default: () => ({}) // 默认值是一个空对象
    },
  },
  data() {
    return {
      dialogVisible:false,
      searchInfo: {
        page: 1,
        pageSize: 10,
        id: this.logData.id,
        objType:this.logData.objType,
        ownerKind:this.logData.ownerKind,
        ownerName:this.logData.ownerName,
        namespace:this.logData.namespace,
        clusterId:this.logData.clusterId,
        nextMarker:0,
        isContainer:this.logData.isContainer,
      },
      nextMarker:0,
      logList: [],
      resType:{
        ["server"]:'云主机',
      },
      actionType:{
        ["update_status"]:'更新状态',
        ["vm_change_bandwidth"]:'调整带宽',
        ["allocate"]:'分配',
        ["create"]:'创建',
        ["Vm Start"]:'开机',
        ["vm_start"]:'开机',
        ["vm_stop"]:'关机',
        ["vm_deploy"]:'部署',
        ["online"]:'上线',
        ["offline"]:'下线',
        ["enable"]:'启用',
        ["disable"]:'禁用',
        ["sync_conf"]:'同步配置',
        ["update"]:'更新',
        ["delete"]:'删除',
        ["同步"]:'同步',
        ["webssh"]:'webssh',
        ["创建集群"]:'创建集群',
        ["创建资源"]:'创建资源',
        ["同步资源"]:'同步资源',
        ["console"]:'console',
        ["public"]:'设为共享',
        ["create_scaling_policy"]:'创建伸缩策略',
        ["delete_scaling_policy"]:'删除伸缩策略',
        ["disk_create_snapshot"]:'磁盘创建快照',
        ["remove_guest"]:'移除实例',
      },
      detailInfo:{}
    }
  },
  created() {
    this.GetLogsList()
  },
  methods: {
    formattedJson(val) {
      try {
        const parsedJson = JSON.parse(val); // 解析 JSON 字符串
        return JSON.stringify(parsedJson, null, 2); // 格式化 JSON，缩进为 2 个空格
      } catch (error) {
        console.error("Invalid JSON string:", val);
        return val; // 如果不是有效的 JSON，直接返回原始字符串
      }
    },
    //打开对话框
    openDialog(val){
      this.detailInfo=val;
      this.dialogVisible=true;
    },
    //关闭对话框
    closeDialog(){
      this.dialogVisible=false;
    },
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        // 如果日期是不合理的，返回空字符串
        return '';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    GetLogsList() {
      getLogsList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.logList=this.logList.concat(JSON.parse(response.data).data);
          if(JSON.parse(response.data).next_marker!=undefined){
            this.nextMarker=Number(JSON.parse(response.data).next_marker);
            this.searchInfo.nextMarker=this.nextMarker
          }else{
            this.nextMarker=0;
            this.searchInfo.nextMarker=0
          }
        }
      })
    },
  },
}
</script>
<style scoped>
.scrollable-content {
  height: 300px;
  overflow-y: auto; /* 垂直滚动   */
  word-break: break-word; /* 自动换行 */
  white-space: pre-wrap; /* 保留换行符 */
}
</style>