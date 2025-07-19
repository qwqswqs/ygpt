<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入数据集名称"
            style="width: 200px;"
        />

        <el-button v-waves class="filter-item"  @click="handleSearch" type="primary">
          搜索
        </el-button>
        <el-button class="filter-item" style="margin-left: 10px;" @click="OpenDialog" type="primary">
          添加
        </el-button>
      </div>

      <!-- 表格展示 -->
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="modelList"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left',}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小">
          <template #default="scope">
            <span>{{ scope.row.size }}</span>
          </template>
        </el-table-column>
        <el-table-column label="关联租户">
          <template #default="scope">
            <span>{{ scope.row.renterID}}</span>
          </template>
        </el-table-column>
        <el-table-column label="上传时间">
          <template #default="scope">
            <span>{{ scope.row.uploadTime }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="DeleteBackInfo(scope.row)">删除</el-button>
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
import {deleteBack,getBackList} from "@/api/yunguan/backup/backup";



export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      total: 0,
      modelList: [],
    }
  },
  created(){
    this.GetBackList()
  },
  methods:{
    FormatDateTime(dateString){
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
    handleSearch(){
      this.GetBackList()
    },
    DeleteBackInfo(val) {
      this.$confirm('此操作将永久删除该数据, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteBack(val).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetBackList(this.searchInfo)
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
    GetBackList() {
      getBackList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.modelList = response.data.list
          this.searchInfo.pageIndex = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.total = response.data.total
        } else {
          this.$message({
            message: '创建失败',
            type: 'warning'
          });
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetBackList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetBackList()
    }

  }
}


</script>
