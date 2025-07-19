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
                @click="addSubOpenDialog"
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
        <el-table-column label="内网IP地址">
          <template #default="scope">
            <span>{{ GetIpName(scope.row.innerIpID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="网关地址">
          <template #default="scope">
            <span>{{scope.row.gatewayIp }}</span>
          </template>
        </el-table-column>
        <el-table-column label="DNS">
          <template #default="scope">
            <span>{{ scope.row.dns }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="负责人">
          <template #default="scope">
            <span>{{ scope.row.header }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status==1?'空闲':'运行' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="租户">
          <template #default="scope">
            <span>{{ GetTenantName(scope.row.renterID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="租用时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.rentTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属VPC">
          <template #default="scope">
            <span>{{ GetName(scope.row.belongVpc) }}</span>
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
          <el-form-item label="内网IP地址" prop="innerIpID">
            <el-select v-model="dataInfo.innerIpID" placeholder="请选择内网IP">
              <el-option
                  v-for="(item,index) in dataIpList"
                  :label="'IP名称:'+item.name+' IP类别:'+belongType[item.belongType-1]"
                  :value="item.ID"
                  :key="index"/>
            </el-select>
          </el-form-item>
          <el-form-item label="网关地址" prop="gatewayIp">
            <el-input v-model="dataInfo.gatewayIp"/>
          </el-form-item>
          <el-form-item label="DNS" prop="dns">
            <el-input v-model="dataInfo.dns"/>
          </el-form-item>
          <el-form-item label="说明" prop="description">
            <el-input v-model="dataInfo.description"/>
          </el-form-item>
          <el-form-item label="负责人" prop="header">
            <el-input v-model="dataInfo.header"/>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="dataInfo.status" placeholder="请选择内网">
              <el-option :label="'空闲'" :value="1"/>
              <el-option :label="'运行'" :value="2"/>
            </el-select>
          </el-form-item>
          <el-form-item label="租户" prop="renterID">
            <el-select v-model="dataInfo.renterID" placeholder="请选择内网IP">
              <el-option
                  v-for="(item,index) in tenantList"
                  :label="item.userName"
                  :value="item.ID"
                  :key="index"/>
            </el-select>
          </el-form-item>
          <el-form-item label="租户时间" prop="rentTime">
            <el-date-picker style="width:100%" v-model="dataInfo.rentTime"/>
          </el-form-item>
          <el-form-item label="所属VPC" prop="belongVpc">
            <el-select v-model="dataInfo.belongVpc" placeholder="请选择归属VPC">
              <el-option :label="'独立'" :value="0"/>
              <el-option
                  v-for="(item,index) in dataVpcList"
                  :label="'VPC名称:'+item.name"
                  :value="item.ID"
                  :key="index"/>
            </el-select>
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import {addSub, deleteSub, getSubList, updateSub} from "@/api/yunguan/network/netSub";
import {getAllIpList} from "@/api/yunguan/network/netIp";
import {getAllVpcList} from "@/api/yunguan/network/netVpc";
import {getAllRenterList} from "@/api/yunguan/run/renter";


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
        innerIpID:'',
        gatewayIp:'',
        dns:'',
        description:'',
        header:'',
        status:'',
        renterID:'',
        rentTime:'',
        belongVpc:'',
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
      dataVpcList:[],
      tenantList:[],
      useType:["vpc","子网","空闲"],
      belongType:["VPC内","子网内","VPC外部","子网外"],
    }
  },
  created() {
    this.GetSubList()
    this.GetAllIpList()
    this.GetAllVpcList()
    this.GetTenantList()
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
        deleteSub(val).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetSubList(this.searchInfo)
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
      this.GetSubList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetSubList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetSubList()
    },
    addSubOpenDialog() {
      this.openDialog = true
      this.IsCreated = true
      this.dataInfo={}
    },
    closeDialog() {
      this.openDialog = false
    },
    UpdateSub() {
      updateSub(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetSubList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddSub() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addSub(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetSubList()
            }
          })
        }
      });
    },
    // 提交
    submitDialog() {
      if (this.IsCreated)this.AddSub()
      else this.UpdateSub()
    },
    // 获取分区列表
    GetSubList() {
      getSubList(this.searchInfo).then(response => {
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
    GetAllVpcList() {
      getAllVpcList().then(response => {
        if (response.code == 0) {
          this.dataVpcList = response.data.list
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
      const tenant = this.tenantList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
    //获取租户名称
    GetName(id) {
      const tenant = this.dataVpcList.find(t => t.ID === id);
      return tenant ? tenant.name : null;
    },
    // 获取分区列表
    GetIpName(val) {
      const ip = this.dataIpList.find(t => t.ID === val)
      return ip ? '起:'+ip.startIp+'\n止:'+ip.endIp : null;
    }
  }
}
</script>
