<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        class="filter-item"
        placeholder="请输入数据集名称"
        style="width: 200px;"
      />

      <el-button v-waves class="filter-item" type="primary">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary">
        添加
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
      <el-table-column label="请求时间">
        <template #default="scope">
          <span>{{ scope.row.time }}</span>
        </template>
      </el-table-column>
      <el-table-column label="请求用户名称">
        <template #default="scope">
          <span>{{ scope.row.applyPerson }}</span>
        </template>
      </el-table-column>
      <el-table-column label="请求状态">
        <template #default="scope">
          <span>{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column label="资源类型">
        <template #default="scope">
          <span>{{ scope.row.resourceType }}</span>
        </template>
      </el-table-column>
      <el-table-column label="调度策略">
        <template #default="scope">
          <span>{{ scope.row.way }}</span>
        </template>
      </el-table-column>
      <el-table-column label="批准人">
        <template #default="scope">
          <span>{{ scope.row.permitPerson }}</span>
        </template>
      </el-table-column>
      <el-table-column label="批准时间">
        <template #default="scope">
          <span>{{ scope.row.permitTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <el-button type="text">编辑</el-button>
        <el-button type="text">删除</el-button>
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
</template>

<script>

export default {
  name: 'ComplexTable',
  data() {
    return {
      searchInfo: {
        pageIndex:1,
        pageSize:10,
      },
      total: 0,
      list: [
        { id: 1, time: '2024-09-24', applyPerson: '张三', applyStatus: '已批准', resourceType: '云主机', way: 'FIFO', permitPerson: '李四', permitTime: '2024-09-24' }
      ]
    }
  }
}
</script>
