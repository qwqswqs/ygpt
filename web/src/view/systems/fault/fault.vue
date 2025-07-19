<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        class="filter-item"
        placeholder="请输入数据集名称"
        style="width: 200px;"
      />

      <el-button v-waves class="filter-item"  type="primary">
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
      <el-table-column label="故障类型">
        <template #default="scope">
          <span>{{ scope.row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column label="描述">
        <template #default="scope">
          <span>{{ scope.row.desc }}</span>
        </template>
      </el-table-column>
      <el-table-column label="处理状态">
        <template #default="scope">
          <span>{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column label="处理人">
        <template #default="scope">
          <span>{{ scope.row.person }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间">
        <template #default="scope">
          <span>{{ scope.row.createTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="处理时间">
        <template #default="scope">
          <span>{{ scope.row.handleTime }}</span>
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
        { id: 1, type: '网络故障', desc: '主机01的网络中断', status: '待处理', person: '张三', createTime: '2024-09-23', handleTime: '' },
        { id: 2, type: '存储设备', desc: '主机01存储空间不足', status: '已处理', person: '张三', createTime: '2024-09-23', handleTime: '2024-09-24' },
        { id: 3, type: '网络故障', desc: '主机01的网络中断', status: '未分配', person: '', createTime: '2024-09-23', handleTime: '' },
        { id: 4, type: '网络故障', desc: '主机01的网络中断', status: '待处理', person: '张三', createTime: '2024-09-23', handleTime: '' }
      ]
    }
  }
}
</script>
