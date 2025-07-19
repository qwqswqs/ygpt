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
                placeholder="请输入配置名称"
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
                @click="addIpOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="dataList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="地址段名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="开始地址">
          <template #default="scope">
            <span>{{ scope.row.startIp }}</span>
          </template>
        </el-table-column>
        <el-table-column label="结束IP">
          <template #default="scope">
            <span>{{ scope.row.endIp }}</span>
          </template>
        </el-table-column>
        <el-table-column label="网关地址">
          <template #default="scope">
            <span>{{ scope.row.gatewayIp }}</span>
          </template>
        </el-table-column>
        <el-table-column label="掩码位数">
          <template #default="scope">
            <span>{{ scope.row.maskBits }}</span>
          </template>
        </el-table-column>
        <el-table-column label="地址个数">
          <template #default="scope">
            <span>{{ scope.row.ipNum }}</span>
          </template>
        </el-table-column>
        <el-table-column label="归属类别">
          <template #default="scope">
            <span>{{ belongType[scope.row.belongType-1]  }}</span>
          </template>
        </el-table-column>
        <el-table-column label="归属">
          <template #default="scope">
            <span>{{ scope.row.belongID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="handleEdit(scope.row)">编辑</el-button>
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
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
      <el-drawer
          v-model="openDialog"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
          size="60%"
      >
        <template #header>
          <div class="flex justify-between items-center">
            <span class="text-lg">{{ IsCreated ? '新增' : '编辑' }}任务</span>
            <div>
              <el-button @click="closeDialog">取 消</el-button>
              <el-button type="primary" @click="submitDialog">确 定</el-button>
            </div>
          </div>
        </template>

        <el-form
            ref="dataInfo"
            :model="dataInfo"
            :rules="IsCreated?rules:''"
            label-width="100px"
        >
          <el-form-item label="地址段名称" prop="name">
            <el-input v-model="dataInfo.name"/>
          </el-form-item>
          <el-form-item label="开始地址" prop="startIp">
            <el-input v-model="dataInfo.startIp"/>
          </el-form-item>
          <el-form-item label="结束地址" prop="endIp">
            <el-input v-model="dataInfo.endIp"/>
          </el-form-item>
          <el-form-item label="网关地址" prop="gatewayIp">
            <el-input v-model="dataInfo.gatewayIp"/>
          </el-form-item>
          <el-form-item label="掩码位数" prop="maskBits">
            <el-input-number style="width: 100%" :min="8" :step="8" :max="24" v-model="dataInfo.maskBits"/>
          </el-form-item>
          <el-form-item label="地址个数" prop="ipNum">
            <el-input-number style="width: 100%" :min="1" v-model="dataInfo.ipNum"/>
          </el-form-item>
          <el-form-item label="归属类别" prop="belongType">
            <el-select v-model="dataInfo.belongType" placeholder="请选择归属类别">
              <el-option
                  v-for="(item,index) in belongType"
                  :label="item"
                  :value="index+1"
                  :key="index"/>
            </el-select>
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import {addIp, deleteIp, getIpList, updateIp} from "@/api/yunguan/network/netIp";


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
      dataInfo: {
        id:0,
        name: '',
        startIp:'',
        endIp:'',
        gatewayIp:'',
        maskBits:'',
        ipNum:'',
        belongType:'',
      },
      belongType:["VPC内","子网内","VPC外部","子网外"],
      rules: {
        name: [
          {required: true, message: '请输入分区名称', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入分区说明', trigger: 'blur'}
        ],
      },
      dataList:[],
    }
  },
  created() {
    this.GetIpList()
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
      this.$confirm('此操作将永久删除该分区, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteIp(val).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetIpList(this.searchInfo)
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
    handleEdit(val) {
      this.dataInfo=val
      this.IsCreated = false
      this.openDialog = true
    },
    handleSearch() {
      this.GetIpList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    addIpOpenDialog() {
      this.openDialog = true
      this.IsCreated = true
      this.dataInfo={}
    },
    closeDialog() {
      this.openDialog = false
    },
    UpdateIp() {
      updateIp(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetIpList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddIp() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addIp(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetIpList()
            }
          })
        }
      });
    },
    // 提交
    submitDialog() {
      if (this.IsCreated)this.AddIp()
      else this.UpdateIp()
    },
    // 获取分区列表
    GetIpList() {
      getIpList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
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
