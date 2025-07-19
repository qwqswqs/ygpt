<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-input
                v-model="searchInfo.name"
                placeholder="请输入区域名称"
            />
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="handleSearch"
            >
              查询
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="snapshotList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
        <el-table-column label="快照名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属实例">
          <template #default="scope">
            <span>{{ scope.row.instanceID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="快照大小">
          <template #default="scope">
            <span>{{ scope.row.size }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <span>{{ scope.row.renterID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="scope">
            <span>{{ scope.row.specDesc }}</span>
          </template>
        </el-table-column>
        <el-table-column label="快照类型">
          <template #default="scope">
            <span>{{ scope.row.snapshotType }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,  prev, pager, next, sizes,jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>
<script>


import {deleteSnapshot, getSnapshotList} from "@/api/yunguan/backup/snapshot";

export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      snapshotType:["磁盘","内存","实例"],
      snapshotList: [],
    }
  },
  created() {
    this.GetSnapshotList()
  },
  methods: {
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
    handleSearch() {
      this.GetImageList()
    },
    handleDelete(val){
      this.$confirm('此操作将永久删除该快照, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSnapshot(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetSnapshotList(this.searchInfo)
          } else {
            this.$message({
              message: '删除失败',
              type: 'error'
            });
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetSnapshotList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetSnapshotList()
    },
    // 获取快照列表
    GetSnapshotList() {
      getSnapshotList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.snapshotList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
  }
}
</script>
