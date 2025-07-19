<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
<!--          <el-form-item>-->
<!--            <el-input-->
<!--                v-model="searchInfo.name"-->
<!--                placeholder="请输入区域名称"-->
<!--            />-->
<!--          </el-form-item>-->
<!--          <el-form-item>-->
<!--            <el-button-->
<!--                -->
<!--                type="primary"-->
<!--                @click="handleSearch"-->
<!--            >-->
<!--              查询-->
<!--            </el-button>-->
<!--          </el-form-item>-->
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="logList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="级别">
          <template #default="scope">
            <span>{{ scope.row.level }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作人">
          <template #default="scope">
            <span>{{ scope.row.user }}</span>
          </template>
        </el-table-column>
        <el-table-column label="日志信息">
          <template #default="scope">
            <span>{{ scope.row.birth }}</span>
          </template>
        </el-table-column>
        <el-table-column label="日志源">
          <template #default="scope">
            <span>{{ scope.row.info }}</span>
          </template>
        </el-table-column>
        <el-table-column label="日志产生时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
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
import SelectImage from "@/components/selectImage/selectImage.vue";
import {deleteSysLog, getSysLogList} from "@/api/yunguan/system/log";


export default {
  name: "zone",
  components: {SelectImage},
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        name:'',
      },
      total: 0,
      logList:[],
    }
  },
  created() {
    this.GetLogList()
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
    handleDelete(val){
      this.$confirm('此操作将永久删除该日志, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSysLog(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetLogList(this.searchInfo)
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
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    // 获取分区列表
    GetLogList() {
      getSysLogList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.logList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    }
  }
}
</script>
