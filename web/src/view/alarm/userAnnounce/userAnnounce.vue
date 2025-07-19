<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-filter-box" style="padding-top: 10px">
      <div class="gva-btn-list">
        <el-form
            ref="filterForm"
            :inline="true"
            :model="searchInfo"
            label-width="110"
        >
          <el-form-item label="资源类别">
            <el-select v-model="searchInfo.resourceType" clearable>
              <el-option label="裸金属" :value="1"></el-option>
              <el-option label="云主机" :value="2"></el-option>
              <el-option label="容器" :value="3"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="告警类别">
            <el-select v-model="searchInfo.alertType" clearable>
              <el-option label="CPU使用率" :value="1"></el-option>
              <el-option label="内存使用率" :value="2"></el-option>
              <el-option label="磁盘使用率" :value="3"></el-option>
              <el-option label="网络入带宽" :value="4"></el-option>
              <el-option label="网络出带宽" :value="5"></el-option>
              <el-option label="GPU使用率" :value="10"></el-option>
              <el-option label="GPU显存使用率" :value="11"></el-option>
              <el-option label="GPU温度" :value="12"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="最早告警时间">
            <el-date-picker
                clearable
                v-model="searchInfo.startTime"
                type="date"
                placeholder="选择日期">
            </el-date-picker>
          </el-form-item>
          <el-form-item label="最晚告警时间">
            <el-date-picker
                clearable
                v-model="searchInfo.endTime"
                type="date"
                placeholder="选择日期">
            </el-date-picker>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="getAlertList">
              搜索
            </el-button>
          </el-form-item>
        </el-form>
      </div>

    </div>
    <el-table
        :data="alarmList"
        row-key="id"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
    >
      <el-table-column label="实例名" prop="instanceName">
        <template #default="scope">
          {{ scope.row.resourceName }}
        </template>
      </el-table-column>
      <el-table-column label="租户名">
        <template #default="scope">
          {{ GetRenterName(scope.row.renterID) || "——" }}
        </template>
      </el-table-column>
      <el-table-column label="资源类别" prop="resourceType">
        <template #default="scope">
          {{ resourceType[scope.row.resourceType - 1] }}
        </template>
      </el-table-column>
      <el-table-column label="告警类型" prop="alarmType">
        <template #default="scope">
          {{ alertType[scope.row.alertType] }}
        </template>
      </el-table-column>
      <el-table-column label="告警值" prop="alertData">
        <template #default="scope">
          {{ getMonitorData(scope.row.alertData,scope.row)}}
        </template>
      </el-table-column>
      <el-table-column label="等级" prop="alertLevel">
        <template #default="scope">
          <el-tag v-if="scope.row.alertLevel===1" type="danger">
            普通
          </el-tag>
          <el-tag v-else type="danger">
            严重
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="告警时间" prop="alertTime">
        <template #default="scope">
          {{ FormatDateTime(scope.row.alertTime) }}
        </template>
      </el-table-column>
      <el-table-column label="是否查看">
        <template #default="scope">
          <span
              :class="{ 'text-red-500':!scope.row.isCheck }"
              @click="handleView(scope.row)">
            {{ scope.row.isCheck? '已查看' : '未查看' }}
          </span>
        </template>
      </el-table-column>
    </el-table>
    <div class="gva-pagination">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="searchInfo.pageIndex"
          :page-sizes="[5, 10, 20]"
          :page-size="searchInfo.pageSize"
          layout="total,  prev, pager, next, sizes,jumper"
          :total="searchInfo.total">
      </el-pagination>
    </div>

  </div>
</template>

<script setup>
import {ref, computed} from 'vue';
import {getSystemAlertList} from "@/api/yunguan/config/alertInfo";
import {getAllRenterList} from "@/api/yunguan/run/renter";

// 模拟告警数据
const alarmList = ref([]);

const handleView = (row) => {
  row.isCheck = true;
  // 这里可以添加将查看状态保存到后端的逻辑
};
const handleSizeChange = (newSize) => {
  searchInfo.value.pageSize = newSize;
  getAlertList()
};

const handleCurrentChange = (newPage) => {
  searchInfo.value.pageIndex = newPage;
  getAlertList()
};
const FormatDateTime = (dateString) => {
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
}
const searchInfo = ref({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
  resourceType: null,
  alertType: null,
  alertTime: null,
  startTime:null,
  endTime:null,
  isRenterData: true,
})
const GetRenterName = (id) => {
  const tempData = renterList.value.find(t => t.ID === id);
  console.log(tempData)
  return tempData ? tempData.userName : null;
};
const resourceType = ref(["裸金属", "云主机", "容器", "存储"]);
const alertType={
  1:'CPU使用率',
  2:'内存使用率',
  3:'磁盘使用率',
  // 4:'磁盘读速率',
  // 5:'磁盘写速率',
  6:'网络入带宽',
  7:'网络出带宽',
  // 8:'网络收包速率',
  // 9:'网络发包速率',
  10:'GPU使用率',
  11:'GPU显存使用率',
  12:'GPU温度'
}
const renterList = ref([])
//获取监控值
const getMonitorData=(data,val)=>{
  if (val.alertType<4||val.alertType===10||val.alertType===11){
    return data+'%';
  } else if (val.alertType===12){
    return data+'°C';
  } else {
    if (data<1024){
      return data+'b/s';
    }
    else if (data<1024*1024){
      return (data/1024)+'Kb/s';
    }
    else if (data<1024*1024*1024){
      return (data/1024/1024)+'Mb/s';
    }
    else{
      return (data/1024/1024/1024)+'Gb/s';
    }
  }
};
//获取租户列表
const getRenterList = () => {
  getAllRenterList().then(response => {
    if (response.code == 0) {
      renterList.value = response.data.list
    }
  })
};
const getAlertList = () => {
  getSystemAlertList(searchInfo.value).then(res => {
    if (res.code == 0) {
      alarmList.value = res.data.list;
      searchInfo.value.total = res.data.total;
    }
  })
}
getAlertList()
getRenterList()
</script>

<style scoped>
.gva-filter-box {
  margin-bottom: 16px;
}

.text-red-500 {
  color: red;
}
</style>