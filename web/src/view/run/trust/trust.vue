<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入租户名称"
            style="width: 200px;"
        />

        <el-button v-waves class="filter-item" type="primary">
          搜索
        </el-button>
      </div>

      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="list"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column label="序号">
          <template #default="scope">
            <span>{{ scope.row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="租户名称">
          <template #default="scope">
            <span>{{ scope.row.userName }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="到期时间">
          <template #default="scope">
            <span>{{ scope.row.dueTime }}</span>
          </template>
        </el-table-column>
        <el-table-column label="欠费金额">
          <template #default="scope">
            <span>{{ scope.row.oweMoney }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <el-button type="text" link>提醒</el-button>
          <el-button type="text" link>停止服务</el-button>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'ComplexTable',
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
      },
      total: 0,
      list: [
        {id: 1, userName: '张三', description: '水印服务', status: '已到期', dueTime: '2024-09-08', oweMoney: 0.00},
        {id: 2, userName: '李四', description: '智能服务', status: '未到期', dueTime: '2024-10-15', oweMoney: 0.00},
        {id: 3, userName: '李四', description: '欠费', status: '欠费', dueTime: '', oweMoney: 1.08}
      ]
    }
  }
}
</script>
