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
            <el-button
                type="primary"
                @click="addDeviceOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="deviceList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ deviceType[scope.row.type-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="品牌">
          <template #default="scope">
            <span>{{ scope.row.brand }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所在位置">
          <template #default="scope">
            <span>{{ scope.row.location }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SN号">
          <template #default="scope">
            <span>{{ scope.row.snNum }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status==1?'正常':'异常' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU类型">
          <template #default="scope">
            <span>{{ scope.row.cpuType }}</span>
          </template>
        </el-table-column>
        <el-table-column label="GPU类型">
          <template #default="scope">
            <span>{{ scope.row.gpuType }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="规格描述">
          <template #default="scope">
            <span>{{ scope.row.specDesc }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
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
      <el-drawer
          v-model="openDialog"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
          size="50%"
      >
        <template #header>
          <div class="flex justify-between items-center">
            <span class="text-lg">{{ IsCreated ? '新增' : '编辑' }}设备</span>
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
          <el-form-item label="类别" prop="type">
            <el-select v-model="dataInfo.type">
              <el-option
                v-for="(item,index) in deviceType"
                :value="index+1"
                :label="item"
                :key="index"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="品牌" prop="brand">
            <el-input v-model="dataInfo.brand"/>
          </el-form-item>
          <el-form-item label="所在位置" prop="location">
            <el-input v-model="dataInfo.location"/>
          </el-form-item>
          <el-form-item label="SN号" prop="snNum">
            <el-input v-model="dataInfo.snNum"/>
          </el-form-item>
          <el-form-item label="CPU类型" prop="cpuType">
            <el-input v-model="dataInfo.cpuType"/>
          </el-form-item>
          <el-form-item label="GPU类型" prop="gpuType">
            <el-input v-model="dataInfo.gpuType"/>
          </el-form-item>
          <el-form-item label="说明" prop="description">
            <el-input v-model="dataInfo.description"/>
          </el-form-item>
          <el-form-item label="规格描述" prop="specDesc">
            <el-input v-model="dataInfo.specDesc"/>
          </el-form-item>
          <el-form-item label="安装日期" prop="installDate">
            <el-date-picker style="width:100%" v-model="dataInfo.installDate"/>
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import SelectImage from "@/components/selectImage/selectImage.vue";
import {addDevice, deleteDevice, getDeviceList, updateDevice} from "@/api/yunguan/resource/resDevice";


export default {
  name: "zone",
  components: {SelectImage},
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        status:0,
        type:0,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      dataInfo: {
        id:0,
        type: '', // 1.物理机 2.网络设备 3.基础设备
        brand:'', //品牌
        location:'',// 所在机房、区域、机柜、机位
        snNum: '', // SN号
        installDate:'',//安装日期
        status:1, //状态 1.正常 2.异常
        gpuType:'',//GPU类型
        cpuType:'',//cpu类型
        description: '',//说明
        specDesc:'',//规格描述
      },
      deviceType:["物理机","网络设备","基础设备"],
      rules: {
        name: [
          {required: true, message: '请输入分区名称', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入分区说明', trigger: 'blur'}
        ],
      },
      deviceList:[],
    }
  },
  created() {
    this.GetDeviceList()
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
        deleteDevice(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetDeviceList(this.searchInfo)
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
      this.GetDeviceList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    addDeviceOpenDialog() {
      this.dataInfo={}
      this.openDialog = true
      this.IsCreated = true
    },
    closeDialog() {
      this.openDialog = false
      this.dataInfo={}
      this.dataInfo.status=1
    },
    UpdateDevice() {
      updateDevice(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetDeviceList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddDevice() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addDevice(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetDeviceList()
            } else {
              this.$message({
                message: '添加失败',
                type: 'error'
              });
            }
          })
        }
      });
    },
    // 提交
    submitDialog() {
      if (this.IsCreated)this.AddDevice()
      else this.UpdateDevice()
    },
    // 获取分区列表
    GetDeviceList() {
      getDeviceList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.deviceList = response.data.list
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
