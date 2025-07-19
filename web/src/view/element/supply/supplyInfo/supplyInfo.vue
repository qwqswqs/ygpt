<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box" style="margin-top: 10px">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-input
                v-model="searchInfo.name"
                placeholder="请输入名称"
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
                @click="addSupplyOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="supplyList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="标题" width="150px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ supplyType[scope.row.supplyType-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述" width="200px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="规格描述" width="200px">
          <template #default="scope">
            <div   v-for="(item, index) in JSON.parse(scope.row.specDesc)" :key="index">
              <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ item.key }}: {{ item.value }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="同步状态">
          <template #default="scope">
            <span>{{ scope.row.isSynced==1?'未同步':'已同步' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status==1?'上架':'下架' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="上架时间" width="180px">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.startTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="下架时间" width="180px">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.endTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150px" :fixed="'right'">
          <template #default="scope">
            <el-button type="text" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" link :disabled="scope.row.isSynced==2" @click="syncSupplyInfo(scope.row)">同步</el-button>
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
            <span class="text-lg">{{ IsCreated ? '新增' : '编辑' }}供需</span>
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
          <el-form-item label="标题" prop="title">
            <el-input v-model="dataInfo.title"/>
          </el-form-item>
          <el-form-item label="类别" prop="supplyType">
            <el-select v-model="dataInfo.supplyType" placeholder="请选择供需类别">
              <el-option
                v-for="(item,index) in supplyType"
                :label="item"
                :value="index+1"
                :key="index"/>
            </el-select>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="dataInfo.description"/>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-switch
                v-model="dataInfo.status"
                :active-value="1"
                :inactive-value="2"
                active-text="上架"
                class="ml-2"
                inactive-text="下架"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="上架时间" prop="startTime">
            <el-date-picker style="width: 100%" v-model="dataInfo.startTime" placeholder="请选择上架时间"/>
          </el-form-item>
          <el-form-item label="下架时间" prop="endTime">
            <el-date-picker style="width: 100%" v-model="dataInfo.endTime" placeholder="请选择下架时间"/>
          </el-form-item>
          <el-form-item label="规格描述">
            <div>
              <el-select v-model="currentKey" placeholder="请选择键" style="width: 200px; margin-right: 10px;">
                <el-option v-for="(item,index) in keyOptionsList" :key="index" :label="item" :value="item" />
              </el-select>
              <el-input v-model="currentValue" placeholder="请输入值" style="width: 200px; margin-right: 10px;" />
              <el-button type="primary" @click="addKeyValue">添加</el-button>
            </div>
          </el-form-item>
          <el-form-item>
            <el-table :data="keyValuePairs" style="width: 90%">
              <el-table-column label="键" prop="key"/>
              <el-table-column label="值" prop="value"/>
              <el-table-column label="操作">
                <template #default="{ row, $index }">
                  <el-button type="primary" @click="removeKeyValue($index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import {addSupplyInfo, SyncSupplyDemandInfo,getSupplyInfoList, deleteSupplyInfo,updateSupplyInfo} from "@/api/yunguan/product/productSupply";


export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        type:1,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      dataInfo: {
        id: 0,
        title:'',
        type:1,
        supplyType:'',
        description: '',
        specDesc: '',
        status:1,
        isSynced:1,
        startTime:'',
        endTime:''
      },
      supplyType:["算力","机柜","带宽","硬件"],
      rules: {
        name: [
          {required: true, message: '请输入分区名称', trigger: 'blur'}
        ],
        submitPerson: [
          {required: true, message: '请选择硬件类型', trigger: 'blur'}
        ],
        handlePerson: [
          {required: true, message: '请选择硬件类型', trigger: 'blur'}
        ],
      },
      supplyList: [],
      currentKey:'',
      currentValue:'',
      keyOptionsList:["高度","宽度","颜色","大小","重量"],
      keyValuePairs:[],
    }
  },
  created() {
    this.GetSupplyList()
  },
  methods: {
    addKeyValue(){
      if(this.currentKey&&this.currentValue){
        this.keyValuePairs.push({key: this.currentKey, value: this.currentValue});
        this.currentKey = '';
        this.currentValue = '';
      }
    },
    updateTextArea(){
    },
    removeKeyValue(index){
      this.keyValuePairs.splice(index,1)
    },
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        // 如果日期是不合理的，返回空字符串
        return '';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    },
    handleDelete(val) {
      this.$confirm('此操作将永久删除该分区, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSupplyInfo(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetSupplyList(this.searchInfo)
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
      if (!val || typeof val !== 'object') {
        console.error('Invalid data provided for editing:', val);
        return;
      }

      // 安全地设置 this.dataInfo，并确保 specDesc 是一个有效的对象或数组
      this.dataInfo =val; // 使用解构赋值创建副本，防止直接修改原始数据
      if (typeof this.dataInfo.specDesc === 'string') {
        try {
          // 如果 specDesc 是 JSON 字符串，尝试解析它
          this.dataInfo.specDesc = JSON.parse(this.dataInfo.specDesc);
        } catch (error) {
          console.error('Failed to parse specDesc JSON string:', error);
          this.dataInfo.specDesc = {}; // 或者设置为默认值
        }
      } else if (!this.dataInfo.specDesc || !Array.isArray(this.dataInfo.specDesc)) {
        this.dataInfo.specDesc = {}; // 或者设置为默认值
      }
      this.dataInfo = val
      this.IsCreated = false
      this.openDialog = true
      console.log(this.dataInfo.specDesc)
      this.keyValuePairs=this.dataInfo.specDesc
    },
    syncSupplyInfo(val) {
      this.dataInfo=val
      this.dataInfo.isSynced=2
      SyncSupplyDemandInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '发送同步请求成功，等待响应',
            type: 'success'
          });
          this.GetSupplyList()
        }
      })
    },
    handleSearch() {
      this.GetSupplyList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    addSupplyOpenDialog() {
      this.openDialog = true
      this.IsCreated = true
      this.dataInfo={}
      this.currentValue=''
      this.currentKey=''
      this.keyValuePairs=[]
      this.dataInfo.type=1
      this.dataInfo.isSynced=1
    },
    closeDialog() {
      this.openDialog = false
    },
    UpdateSupply() {
      this.dataInfo.isSynced=1
      updateSupplyInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetSupplyList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddSupply() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addSupplyInfo(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetSupplyList()
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
      this.dataInfo.specDesc=JSON.stringify(this.keyValuePairs)
      if (this.IsCreated) this.AddSupply()
      else this.UpdateSupply()
    },
    // 获取工单列表
    GetSupplyList() {
      getSupplyInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.supplyList = response.data.list
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
