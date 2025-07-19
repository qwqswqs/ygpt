<template>
  <div>
    <!-- 筛选表单 -->
    <div class="gva-filter-box">
      <div class="gva-btn-list">
        <el-form
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item label="开始时间">
            <el-date-picker
                v-model="searchInfo.startTime"
                type="date"
                placeholder="选择开始日期"
            />
          </el-form-item>
          <el-form-item label="结束时间">
            <el-date-picker
                v-model="filterInfo.endTime"
                type="date"
                placeholder="选择结束日期"
            />
          </el-form-item>
          <el-form-item label="租户名">
            <el-input v-model="filterInfo.tenantName" placeholder="输入租户名"/>
          </el-form-item>
          <el-form-item label="资源类别">
            <el-select clearable v-model="searchInfo.type">
              <el-option label="裸金属" :value="1"/>
              <el-option label="云主机" :value="2"/>
              <el-option label="容器" :value="3"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="getDataList"
            >
              筛选
            </el-button>
            <el-button
                type="primary"
                @click="exportTable"
            >
              导出
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <!--    &lt;!&ndash; 状态栏 &ndash;&gt;-->
    <!--    <div class="status-bar">-->
    <!--      <span>已选择记录数: {{ selectedRows.length }}</span>-->
    <!--      <span>费用总数: {{ formatPrice(selectedTotalCost) }}</span>-->
    <!--    </div>-->
    <!-- 数据表格 -->
    <el-table
        :data="dataList"
        row-key="id"
        @selection-change="handleSelectionChange"
        ref="tableRef"
        v-model:selection="selectedRows"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
    >
      <el-table-column type="selection" width="55"/>
      <el-table-column label="租户名">
        <template #default="scope">
          {{ GetRenterName(scope.row.renterID) }}
        </template>
      </el-table-column>
      <el-table-column label="数据来源">
        <template #default="scope">
          {{ scope.row.dataSource == 1 ? '本地' : '算力圈' }}
        </template>
      </el-table-column>
      <el-table-column label="资源类别" prop="resourceType">
        <template #default="scope">
          {{ resourceType[scope.row.type - 1] }}
        </template>
      </el-table-column>
      <el-table-column label="规格描述" prop="resourceType">
        <template #default="scope">
          <el-button v-if="scope.row.dataSource != 1" type="text" @click="OpenDetailDialog(scope.row)">详情</el-button>
          <span v-if="scope.row.dataSource == 1" >无</span>
        </template>
      </el-table-column>
      <el-table-column label="开始时间">
        <template #default="scope">
          {{ FormatDateTime(scope.row.startTime) }}
        </template>
      </el-table-column>
      <el-table-column label="释放时间">
        <template #default="scope">
          {{ FormatDateTime(scope.row.endTime) }}
        </template>
      </el-table-column>
      <el-table-column label="计费方式">
        <template #default="scope">
          {{GetPayTypeName(scope.row)}}
        </template>
      </el-table-column>
      <el-table-column label="购买时长">
        <template #default="scope">
          {{scope.row.description!=""?JSON.parse(scope.row.description).duration:0}}
        </template>
      </el-table-column>
      <el-table-column label="费用">
        <template #default="scope">
          {{scope.row.description!=""?JSON.parse(scope.row.description).paidAmount!=undefined?JSON.parse(scope.row.description).paidAmount:0+"元":0+"元"}}
        </template>
      </el-table-column>
    </el-table>
    <!-- 分页组件 -->
    <div class="gva-pagination">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="searchInfo.pageIndex"
          :page-sizes="[10, 20, 30]"
          :page-size="searchInfo.pageSize"
          layout="total,prev, pager, next,  sizes, jumper"
          :total="searchInfo.total"
      >
      </el-pagination>
    </div>

    <el-drawer v-model="detailVisible" :show-close="true" :title="'详情'" size="50%">
      <div v-if="specJson.length==0">
        <h3>该条记录不存在规格描述</h3>
      </div>
      <el-descriptions v-if="specJson.length>0" size="small" border>
        <el-descriptions-item
            v-for="(item, index) in specJson"
            :key="index"
            :label="item.name"
            :span="specJson.length"
            label-align="right"
            align="left"
        >
          {{ item.value[item.valueIndex != undefined ? item.valueIndex : 0] }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script>


import {getAllRenterList, queryRenterResList} from "@/api/yunguan/run/renter";
import * as XLSX from "xlsx";

export default {
  data() {
    return {
      // 模拟数据
      dataList: [],
      detailVisible: false,
      //资源类型
      resourceType: ["裸金属", "云主机", "容器", "存储"],
      // 筛选条件
      filterInfo: {
        startTime: '',
        endTime: '',
        tenantName: '',
        resourceType: ''
      },
      // 选择的行
      selectedRows: [],
      // 表格引用
      tableRef: null,
      // 当前页码
      currentPage: 1,
      // 每页显示条数
      pageSize: 10,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        total: 0,
        type: null,
      },
      renterList: [],
      specJson:{},
    };
  },
  created() {
    this.getDataList();
    this.GetRenterList();
  },
  methods: {
    //打开详情对话框
    OpenDetailDialog(val){
      if (val.description==""){
        this.specJson=[];
        this.detailVisible=true;
        return;
      }
      const tempData=JSON.parse(val.description)
      this.specJson = JSON.parse(tempData.specJson)
      this.detailVisible=true;
    },
    //获取收费类型
    GetPayTypeName(val){
      if (val.description=="")return "无";
      const tempData=JSON.parse(val.description)
      if (tempData.billingType==undefined)return "";
      switch (tempData.billingType){
        case 1:
          return "按时计费";
        case 2:
          return "按日计费";
        case 8:
          return "按月计费";
        case 32:
          return "按年计费";
        default:
          return "不限制"
      }
    },
    //格式化日期
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
    //获取租户列表
    GetRenterList() {
      getAllRenterList().then(response => {
        if (response.code == 0) {
          this.renterList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //获取租户名称
    GetRenterName(id) {
      const tenant = this.renterList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
    //获取对账列表
    getDataList() {
      queryRenterResList(this.searchInfo).then(res => {
        if (res.code == 0) {
          this.dataList = res.data.list;
          this.searchInfo.total = res.data.total;
        }
      })
    },
    // 处理筛选逻辑
    handleFilter() {
      this.currentPage = 1; // 筛选后回到第一页
    },
    // 处理选择变化
    handleSelectionChange(val) {
      this.selectedRows = val;
    },
    // 格式化价格
    formatPrice(price) {
      return Number(price).toFixed(2);
    },
    // 导出表格
    exportTable() {
      const exportData = this.dataList;
      const ws = XLSX.utils.json_to_sheet(exportData);
      const wb = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
      XLSX.writeFile(wb, 'table_data.xlsx');
    },
    // 处理每页显示条数变化
    handleSizeChange(newSize) {
      this.searchInfo.pageSize = newSize;
      this.getDataList()
    },
    // 处理页码变化
    handleCurrentChange(newPage) {
      this.searchInfo.pageIndex = newPage;
      this.getDataList()
    }
  }
};
</script>

<style scoped>
.gva-filter-box {
  margin-bottom: 16px;
}

.status-bar {
  margin-bottom: 10px;
  display: flex;
  gap: 20px;
}
</style>