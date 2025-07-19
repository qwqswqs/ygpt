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
                @click="addNatOpenDialog"
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
        <el-table-column label="NAT名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="用途类别">
          <template #default="scope">
            <span>{{useType [scope.row.useType-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属VPC或者子网名称">
          <template #default="scope">
            <span>{{ scope.row.belongName }}</span>
          </template>
        </el-table-column>
        <el-table-column label="NAT规则">
          <template #default="scope">
            <span>{{ scope.row.natRule }}</span>
          </template>
        </el-table-column>
        <el-table-column label="内网IP地址">
          <template #default="scope">
            <span>{{ GetIpName(scope.row.innerIpID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="外网IP地址">
          <template #default="scope">
            <span>{{ scope.row.outerIpID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所在设备">
          <template #default="scope">
            <span>{{ scope.row.deviceID }}</span>
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
          <el-form-item label="NAT名称" prop="name">
            <el-input v-model="dataInfo.name"/>
          </el-form-item>
          <el-form-item label="说明" prop="description">
            <el-input v-model="dataInfo.description"/>
          </el-form-item>
          <el-form-item label="用途类别" prop="useType">
            <el-select v-model="dataInfo.useType" placeholder="请选择归属类别">
              <el-option
                  v-for="(item,index) in useType"
                  :label="item"
                  :value="index"
                  :key="index"/>
            </el-select>
          </el-form-item>
          <el-form-item label="所属VPC或者子网名称" prop="belongName">
            <el-input v-model="dataInfo.belongName"/>
          </el-form-item>
          <el-form-item label="NAT规则" prop="natRule">
            <el-input v-model="dataInfo.natRule"/>
          </el-form-item>
          <el-form-item label="内网IP地址" prop="innerIpID">
            <el-select v-model="dataInfo.innerIpID" placeholder="请选择归属类别">
              <el-option
                  v-for="(item,index) in dataIpList"
                  :label="'IP名称:'+item.name+' IP类别:'+belongType[item.belongType-1]"
                  :value="item.ID"
                  :key="index"/>
            </el-select>
          </el-form-item>
          <el-form-item label="外网IP地址" prop="outerIpID">
            <el-select v-model="dataInfo.outerIpID" placeholder="请选择归属类别">
              <el-option
                  v-for="(item,index) in dataIpList"
                  :label="'IP名称:'+item.name+' IP类别:'+belongType[item.belongType-1]"
                  :value="item.startIp"
                  :key="index"/>
            </el-select>
          </el-form-item>
<!--          <el-form-item label="所在设备" prop="deviceID">-->
<!--            <el-input v-model="dataInfo.deviceID"/>-->
<!--          </el-form-item>-->
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import {addNat, deleteNat, getNatList, updateNat} from "@/api/yunguan/network/netNat";
import {getAllIpList} from "@/api/yunguan/network/netIp";


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
        description:'',
        useType:'',
        belongName:'',
        natRule:'',
        innerIpID:'',
        outerIpID:'',
        // deviceID:'',
      },
      rules: {
        name: [
          {required: true, message: '请输入分区名称', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入分区说明', trigger: 'blur'}
        ],
      },
      dataList:[],
      dataIpList:[],
      useType:["vpc","子网","空闲"],
      belongType:["VPC内","子网内","VPC外部","子网外"],
    }
  },
  created() {
    this.GetNatList()
    this.GetAllIpList()
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
        deleteNat(val).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetNatList(this.searchInfo)
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
      this.GetNatList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetNatList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetNatList()
    },
    addNatOpenDialog() {
      this.openDialog = true
      this.IsCreated = true
      this.dataInfo={}
    },
    closeDialog() {
      this.openDialog = false
    },
    UpdateNat() {
      updateNat(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetNatList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddNat() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addNat(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetNatList()
            }
          })
        }
      });
    },
    // 提交
    submitDialog() {
      if (this.IsCreated)this.AddNat()
      else this.UpdateNat()
    },
    // 获取分区列表
    GetNatList() {
      getNatList(this.searchInfo).then(response => {
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
    },
    // 获取分区列表
    GetAllIpList() {
      getAllIpList().then(response => {
        if (response.code == 0) {
          this.dataIpList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取分区列表
    GetIpName(val) {
      const ip = this.dataIpList.find(t => t.ID === val)
      return ip ? '起:'+ip.startIp+'\n止:'+ip.endIp : null;
    }
  }
}
</script>
