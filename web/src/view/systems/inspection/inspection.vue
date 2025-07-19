<template>
  <div>
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
          :data="operateList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="操作类别">
          <template #default="scope">
            <span>{{ operateType[scope.row.type-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="开始时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.startTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="结束时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.endTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="发现问题数量">
          <template #default="scope">
            <span>{{ scope.row.problemNum }}</span>
          </template>
        </el-table-column>
        <el-table-column label="问题记录文本">
          <template #default="scope">
            <span>{{ scope.row.problemText }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作人">
          <template #default="scope">
            <span>{{ GetTenantName(scope.row.operator) }}</span>
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
import {deleteSysOperate, getSysOperateList} from "@/api/yunguan/system/operate";
import {getAllUserList} from "@/api/yunguan/run/renter"

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
      operateList:[],
      userList:[],
      operateType:["一键巡查","病毒扫描","漏洞扫描"]
    }
  },
  created() {
    this.GetOperateList()
    this.GetUserList()
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
        deleteSysOperate(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetOperateList(this.searchInfo)
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
      this.GetOperateList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetOperateList()
    },
    // 获取分区列表
    GetOperateList() {
      getSysOperateList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.operateList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取分区列表
    GetUserList() {
      getAllUserList().then(response => {
        if (response.code == 0) {
          this.userList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //获取租户名称
    GetTenantName(id) {
      const tenant = this.userList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
  }
}
</script>
