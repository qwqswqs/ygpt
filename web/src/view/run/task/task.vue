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
                placeholder="请输入模型名称"
            />
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="handleSearch"
            >
              查询
            </el-button>
            <el-button
                type="primary"
                @click="addRenterTaskOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="taskList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="任务名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="任务描述">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="模型名称">
          <template #default="scope">
            <span>{{ GetElementName(scope.row.modelID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="数据集名称" width="100px">
          <template #default="scope">
            <span>{{ GetElementName(scope.row.datasetID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="算法名称">
          <template #default="scope">
            <span>{{ GetElementName(scope.row.algoID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="租户名称">
          <template #default="scope">
            <span>{{ GetTenantName(scope.row.renterID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="任务类型">
          <template #default="scope">
            <span>{{ taskType[scope.row.type - 1] }}</span>
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
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link
          @click="handleDelete(scope.row)">删除</el-button>
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
import {getAllShareElementList,} from "@/api/yunguan/element/element";
import SelectImage from "@/components/selectImage/selectImage.vue";
import {deleteRenterTask, getRenterTaskList} from "@/api/yunguan/run/task";
import {getAllRenterList} from "@/api/yunguan/run/renter";

export default {
  name: "NodeModel",
  components: {SelectImage},
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      taskType: ["推理", "训练", "其他"],
      dataInfo: {
        id: 0,
        name: '',
        description: '',
        type: '', //1推理 2训练
        modelID: '',
        datasetID: '',
        algoID: '',
        instanceID: '',
      },
      rules: {
        taskName: [
          {required: true, message: '请输入任务名称', trigger: 'blur'}
        ],
        taskType: [
          {required: true, message: '请输入任务类型', trigger: 'blur'}
        ],
        elementIds: [
          {required: true, message: '请选择元素', trigger: 'blur'}
        ],
      },
      modelType: 1, //1节点模型 2我的模型
      datasetType: 1, //1节点数据集 2我的数据集
      algoType: 1, //1节点算法 2我的算法
      taskList: [],
      elementList: [], //节点模型列表
    }
  },
  created() {
    this.GetRenterTaskList()
    this.GetTenantList()
    this.GetShareElementList()
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
    handleDelete(val) {
      this.$confirm('此操作将永久删除该任务, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteRenterTask(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetRenterTaskList(this.searchInfo)
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
    handleSearch() {
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    // 获取节点元素列表
    GetShareElementList() {
      getAllShareElementList().then(response => {
        if (response.code == 0) {
          //1模型 2算法 3数据集
          this.elementList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取任务列表
    GetRenterTaskList() {
      getRenterTaskList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.taskList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //获取租户列表
    GetTenantList() {
      getAllRenterList().then(response => {
        if (response.code == 0) {
          this.tenantList = response.data.list
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
      if (Array.isArray(this.tenantList) && this.tenantList.length > 0) {
        const tenant = this.tenantList.find(t => t.ID === id);
        return tenant? tenant.userName : null;
      } else {
        console.log('tenantList is not available or empty.');
        return null;
      }
    },
    //获取元素名称
    GetElementName(id) {
      const tenant = this.elementList.find(t => t.ID === id);
      return tenant ? tenant.name : null;
    }
  }
}
</script>

