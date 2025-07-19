<template>
  <div>
    <!-- 筛选表单 -->
    <div class="gva-filter-box">
      <div class="gva-btn-list">
        <el-form
            ref="filterForm"
            :inline="true"
            :model="filterInfo"
        >
          <el-form-item label="开始时间">
            <el-date-picker
                v-model="filterInfo.startTime"
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
            <el-input v-model="filterInfo.tenantName" placeholder="输入租户名" />
          </el-form-item>
          <el-form-item label="资源类别">
            <el-select v-model="filterInfo.resourceType">
              <el-option label="云主机" value="云主机" />
              <el-option label="容器" value="容器" />
              <el-option label="裸金属" value="裸金属" />
              <el-option label="存储" value="存储" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="handleFilter"
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
    <!-- 状态栏 -->
    <div class="status-bar">
      <span>已选择记录数: {{ selectedRows.length }}</span>
      <span>费用总数: {{ formatPrice(selectedTotalCost) }}</span>
    </div>
    <!-- 数据表格 -->
    <el-table
        :data="paginatedData"
        row-key="id"
        @selection-change="handleSelectionChange"
        ref="tableRef"
        v-model:selection="selectedRows"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="资源类别" prop="resourceType" />
      <el-table-column label="规格描述" prop="specification" />
      <el-table-column label="开始时间" prop="startTime" />
      <el-table-column label="结束时间" prop="endTime" />
      <el-table-column label="计费方式" prop="billingMethod" />
      <el-table-column label="费用" prop="cost" />
    </el-table>
    <!-- 分页组件 -->
    <div class="gva-pagination">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[10, 20, 30]"
          :page-size="pageSize"
          layout="total,prev, pager, next,  sizes, jumper"
          :total="filteredData.length"
      >
      </el-pagination>
    </div>

  </div>
</template>

<script>


export default {
  data() {
    return {
      // 模拟数据
      dataList: [
        {
          id: 1,
          tenantName: '租户 A',
          tenantCategory: '本地',
          resourceType: '云主机',
          specification: '规格 1',
          startTime: '2024-01-01',
          endTime: '2024-01-02',
          billingMethod: '按小时计费',
          cost: 100
        },
        {
          id: 2,
          tenantName: '租户 B',
          tenantCategory: '外部',
          resourceType: '容器',
          specification: '规格 2',
          startTime: '2024-01-03',
          endTime: '2024-01-04',
          billingMethod: '按天计费',
          cost: 200
        },
        {
          id: 3,
          tenantName: '租户 A',
          tenantCategory: '本地',
          resourceType: '裸金属',
          specification: '规格 3',
          startTime: '2024-01-05',
          endTime: '2024-01-06',
          billingMethod: '按小时计费',
          cost: 300
        },
        {
          id: 4,
          tenantName: '租户 C',
          tenantCategory: '外部',
          resourceType: '存储',
          specification: '规格 4',
          startTime: '2024-01-07',
          endTime: '2024-01-08',
          billingMethod: '按容量计费',
          cost: 400
        }
      ],
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
      pageSize: 10
    };
  },
  computed: {
    // 计算筛选后的数据
    filteredData() {
      return this.dataList.filter(item => {
        const startMatch = !this.filterInfo.startTime || new Date(item.startTime) >= new Date(this.filterInfo.startTime);
        const endMatch = !this.filterInfo.endTime || new Date(item.endTime) <= new Date(this.filterInfo.endTime);
        const tenantMatch = !this.filterInfo.tenantName || item.tenantName.includes(this.filterInfo.tenantName);
        const resourceMatch = !this.filterInfo.resourceType || item.resourceType === this.filterInfo.resourceType;
        return startMatch && endMatch && tenantMatch && resourceMatch;
      });
    },
    // 计算选择记录的费用总数
    selectedTotalCost() {
      return this.selectedRows.reduce((total, row) => {
        return total + Number(row.cost);
      }, 0);
    },
    // 计算当前页要显示的数据
    paginatedData() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.filteredData.slice(start, end);
    }
  },
  methods: {
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
      const exportData = this.selectedRows.length > 0 ? this.selectedRows : this.filteredData;
      const ws = XLSX.utils.json_to_sheet(exportData);
      const wb = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
      XLSX.writeFile(wb, 'table_data.xlsx');
    },
    // 处理每页显示条数变化
    handleSizeChange(newSize) {
      this.pageSize = newSize;
    },
    // 处理页码变化
    handleCurrentChange(newPage) {
      this.currentPage = newPage;
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