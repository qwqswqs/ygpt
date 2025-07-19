<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
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
        <el-table-column label="租户名称">
          <template #default="scope">
            <span>{{ scope.row.username }}</span>
          </template>
        </el-table-column>
        <el-table-column label="缓存区名称">
          <template #default="scope">
            <span>{{ scope.row.cacheName }}</span>
          </template>
        </el-table-column>
        <el-table-column label="缓存类型">
          <template #default="scope">
            <span>{{ scope.row.cacheType }}</span>
          </template>
        </el-table-column>
        <el-table-column label="缓存大小(GB)">
          <template #default="scope">
            <span>{{ scope.row.cacheTotal }}</span>
          </template>
        </el-table-column>
        <el-table-column label="已使用(GB)">
          <template #default="scope">
            <span>{{ scope.row.cacheUsed }}</span>
          </template>
        </el-table-column>
        <el-table-column label="文件数量">
          <template #default="scope">
            <span>{{ scope.row.fileNum }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <el-button type="text" link>扩容</el-button>
          <el-button type="text" link>缩容</el-button>
          <el-button type="text" link>启用</el-button>
          <el-button type="text" link>禁用</el-button>
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
        {
          username: '张三',
          cacheName: 'cache01',
          status: '启用',
          cacheType: '磁盘',
          cacheTotal: '100',
          cacheUsed: '60',
          fileNum: '600'
        },
        {
          username: '李四',
          cacheName: 'cache02',
          status: '启用',
          cacheType: '内存',
          cacheTotal: '100',
          cacheUsed: '60',
          fileNum: '600'
        },
        {
          username: '王五',
          cacheName: 'cache03',
          status: '启用',
          cacheType: '磁盘',
          cacheTotal: '100',
          cacheUsed: '60',
          fileNum: '600'
        },
        {
          username: '赵六',
          cacheName: 'cache04',
          status: '启用',
          cacheType: '磁盘',
          cacheTotal: '100',
          cacheUsed: '60',
          fileNum: '600'
        }
      ]
    }
  }
}
</script>
