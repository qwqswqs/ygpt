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
                @click="addDiskOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="diskList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ statusType[scope.row.status-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="登记时间">
          <template #default="scope">
            <span>{{FormatDateTime(scope.row.registerTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="容量大小">
          <template #default="scope">
            <span>{{ scope.row.capacity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="存储类别">
          <template #default="scope">
            <span>{{ storageType[scope.row.storageType-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="品牌">
          <template #default="scope">
            <span>{{ scope.row.brand }}</span>
          </template>
        </el-table-column>
        <el-table-column label="规格描述">
          <template #default="scope">
            <span>{{ scope.row.specification }}</span>
          </template>
        </el-table-column>
        <el-table-column label="磁盘格式">
          <template #default="scope">
            <span>{{ scope.row.diskFormat }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属服务器">
          <template #default="scope">
            <span>{{ scope.row.serverID }}</span>
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
          <el-form-item label="状态" prop="status">
            <el-select v-model="dataInfo.status">
              <el-option
                  v-for="(item,index) in statusType"
                  :value="index+1"
                  :label="item"
                  :key="index"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="容量" prop="capacity">
            <el-input-number style="width:100%" :min="0" v-model="dataInfo.capacity"/>
          </el-form-item>
          <el-form-item label="存储类别" prop="storageType">
            <el-select v-model="dataInfo.storageType">
              <el-option
                  v-for="(item,index) in storageType"
                  :value="index+1"
                  :label="item"
                  :key="index"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="品牌" prop="brand">
            <el-input v-model="dataInfo.brand"/>
          </el-form-item>
          <el-form-item label="编号" prop="serialNumber">
            <el-input v-model="dataInfo.serialNumber"/>
          </el-form-item>
          <el-form-item label="磁盘格式" prop="diskFormat">
            <el-input v-model="dataInfo.diskFormat"/>
          </el-form-item>
          <el-form-item label="规格描述" prop="specification">
            <el-input v-model="dataInfo.specification"/>
          </el-form-item>
<!--          <el-form-item label="所属服务器" prop="serverID">-->
<!--            <el-input v-model="dataInfo.serverID"/>-->
<!--          </el-form-item>-->
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import SelectImage from "@/components/selectImage/selectImage.vue";
import {addDisk, deleteDisk, getDiskList, updateDisk} from "@/api/yunguan/resource/resDisk";


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
        status: '', // 1.使用 2.空闲 3.损坏
        capacity:'', //品牌
        storageType:'',// 存储类别
        brand: '', // 品牌
        serialNumber: '', // 编号
        specification:'',//规格描述
        diskFormat:'',//磁盘格式
        serverID:0,//服务器ID
      },
      statusType:["使用","空闲","损坏"],
      storageType:["HDD","SSD"],
      rules: {
        name: [
          {required: true, message: '请输入分区名称', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入分区说明', trigger: 'blur'}
        ],
      },
      iskList:[],
    }
  },
  created() {
    this.GetDiskList()
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
        deleteDisk(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetDiskList(this.searchInfo)
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
      this.GetDiskList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    addDiskOpenDialog() {
      this.dataInfo={}
      this.openDialog = true
      this.IsCreated = true
    },
    closeDialog() {
      this.openDialog = false
      this.dataInfo={}
      this.dataInfo.status=1
    },
    UpdateDisk() {
      updateDisk(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetDiskList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddDisk() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addDisk(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetDiskList()
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
      if (this.IsCreated)this.AddDisk()
      else this.UpdateDisk()
    },
    // 获取分区列表
    GetDiskList() {
      getDiskList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.diskList = response.data.list
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
