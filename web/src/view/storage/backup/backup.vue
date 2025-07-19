<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入名称"
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
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="存储类型">
          <template #default="scope">
            <span>{{ scope.row.type }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属域">
          <template #default="scope">
            <span>{{ scope.row.zone }}</span>
          </template>
        </el-table-column>
        <el-table-column label="共享范围">
          <template #default="scope">
            <span>{{ scope.row.share }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <el-button type="text" link>同步状态</el-button>
          <el-dropdown split-button style="margin-left: 10px" type="text">
            更多
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>设置共享</el-dropdown-item>
              <el-dropdown-item>删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
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
        {name: 'minio', status: '在线', type: '对象存储', zone: 'Default', share: '全局共享'}
      ]
    }
  }
}
</script>
