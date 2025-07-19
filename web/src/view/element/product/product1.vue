<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo">
          <el-form-item>
            <el-input v-model="searchInfo.name" placeholder="请输入产品名称"/>
          </el-form-item>
          <el-form-item>
            <el-button  type="primary" @click="handleSearch">查询</el-button>
            <el-button  type="primary" @click="OpenDialog">添加</el-button>
            <el-button type="primary" @click="modIsSynced()">立即同步</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table :data="algoList" row-key="ID">
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <el-switch
                v-model="scope.row.status"
                :active-value="1"
                :inactive-value="2"
                active-text="启用"
                class="ml-2"
                inactive-text="禁用"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
                @click="modStatus(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="开放类别">
          <template #default="scope">
            <span>{{ openList[scope.row.openType] }}</span>
          </template>
        </el-table-column>

        <el-table-column label="描述">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用次数">
          <template #default="scope">
            <span>{{ scope.row.useTime }}</span>
          </template>
        </el-table-column>

        <el-table-column label="标签">
          <template #default="scope">
            <el-tag
                v-for="(tag,index) in scope.row.usage"
                :key="index"
                :disable-transitions="false"
            >
              <span> {{ tag }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="文件名称">
          <template #default="scope">
            <span>{{ GetFileName(scope.row.fileName) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="文件大小">
          <template #default="scope">
            <span>{{ GetFileSize(scope.row.fileSize) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="上传时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.uploadTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="是否同步">
          <template #default="scope">
            <span>{{ ShareStatus[scope.row.isSynced] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" @click="DeleteAlgoInfo(scope.row)">删除</el-button>
            <el-button type="text">发布</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
      <el-dialog v-model="dialogVisible" :title="IsCreate? '创建算力产品':'编辑算力产品'" width="800px">
        <el-form ref="algoData" :model="algoData" :rules="IsCreate?algoRules:''" label-width="120px">
          <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="产品名称" prop="name">
                    <el-input v-model="algoData.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="产品类别" prop="usage">
                    <el-select v-model="algoData.usage" multiple placeholder="请选择" style="width: 100%">
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="产品说明" prop="description">
                <el-input v-model="algoData.description" placeholder="请描述产品说明" type="textarea"/>
              </el-form-item>
              <div @click="activeName ='second'" style="padding-left:50px;cursor: pointer;text-decoration: underline;">定义规格 >></div>
            </el-tab-pane>
            <el-tab-pane label="规格类型" name="second">
              <el-form-item label="规格" style="margin-top:20px">
                <el-radio-group v-model="algoData.usage">
                  <el-radio value="1">单一规格</el-radio>
                  <el-radio value="2">多重规格</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-row>
                <el-col :span="12">
                  已定义规格数量：2
                </el-col>
                <el-col :span="12">
                  <div style="text-align: right;">
                    <el-button>重新定义</el-button>
                    <el-button type="primary">选择现有</el-button>
                  </div>
                </el-col>
              </el-row>
              <el-table :data="algoList" row-key="ID">
                <el-table-column label="序号" type="index" width="60"></el-table-column>
                <el-table-column label="规格名称">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="规格类别">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="规格内容">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="价格">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="提示">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作">
                  <template #default="scope">
                    <el-button type="text" >删除</el-button>
                    <el-button type="text" >修改</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="规格名称" prop="name">
                    <el-input v-model="algoData.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="规格类别" prop="openType">
                    <el-select v-model="algoData.openType" multiple placeholder="请选择" style="width: 100%">
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="定义价格" prop="description">
                <el-input v-model="algoData.description" placeholder="请填写价格"/>
              </el-form-item>
              <el-form-item label="提示信息" prop="description">
                <el-input v-model="algoData.description" placeholder="请填写提示信息"/>
              </el-form-item>
              <el-form-item label="规格说明" prop="description">
                <el-input v-model="algoData.description" placeholder="请描述规格说明" type="textarea"/>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
      
          <!-- <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item> -->
        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="CloseDialog">取 消</el-button>
          <el-button type="primary" @click="handleConfirm">保 存</el-button>
          <el-button type="warning">提 交</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>
<script>
import {
  addElement,
  deleteElement,
  getNodeElementList, syncElement,
  updateElement,
  updateElementFileInfo
} from "@/api/yunguan/element/element";

export default {
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 2, //1产品 2产品 3数据集
        source: 1, //来源：1节点 2租户
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      statusList: ["", "启用", "禁用"],
      openList: ['请选择开放类别', "不开放", "节点内可用", "算力可见节点内可用", "算力可见算力外可用"],
      dialogTitle: '创建产品',
      algoList: [],
      dialogVisible: false,
      algoData: {
        id: '',
        name: '',
        usage: '',
        description: '',
        status: '',
        openType: 1,
        isSynced: 1,
        source: 1, //来源 1节点 2租户上传
        file: null,
        type: 2,
        price: 0
      },
      algoRules: {
        name: [
          {required: true, message: '请输入产品名称', trigger: 'blur'}
        ],
        usage: [
          {required: true, message: '请选择产品标签', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入产品描述', trigger: 'blur'}
        ],
        openType: [
          {required: true, message: '请选择产品开放类别', trigger: 'blur'}
        ],
        file: [
          {required: true, message: '请上传文件', trigger: 'blur'}
        ],

      },
      IsCreate: true,
      ShareStatus: ["", "未同步", "已同步"],
      activeName:'first',
    }
  },
  created() {
    this.GetAlgoList()
  },
  methods: {
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
    handleSearch() {
      this.GetAlgoList()
    },
    uploadMaterial() {
      this.$refs.file.dispatchEvent(new MouseEvent("click"));
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.algoData.file = inputFile
    },
    GetFileName(val) {
      // 使用逗号分割文件名
      var parts = val.split(',');
      // 返回分割后的第一个元素，即“,”前的文件名
      return parts[0];
    },
    GetFileSize(val) {
      var B = 1024
      var KB = 1024 * 1024
      var MB = 1024 * 1024 * 1024
      var GB = 1024 * 1024 * 1024 * 1024
      if (val < KB) {
        return (val / B).toFixed(2) + "B";
      } else if (val >= KB && val <= MB) {
        return (val / KB).toFixed(2) + "KB";
      } else if (val >= MB && val <= GB) {
        return (val / MB).toFixed(2) + "MB";
      } else if (val > GB) {
        return (val / GB).toFixed(2) + "GB";
      }
    },
    OpenDialog() {
      this.IsCreate = true
      this.dialogVisible = true
    },
    CloseDialog() {
      this.dialogVisible = false
    },
    handleEdit(val) {
      this.algoData.id = val.ID
      this.algoData.name = val.name
      this.algoData.type = val.type
      this.algoData.usage = val.usage
      this.algoData.description = val.description
      this.algoData.openType = val.openType
      this.algoData.status = val.status == 1 ? 1 : 2
      this.algoData.isSynced = 1
      this.algoData.price = val.price
      this.IsCreate = false
      this.dialogVisible = true
    },
    handleConfirm() {
      if (this.IsCreate) this.AddAlgoInfo()
      else this.UpdateAlgoInfo()
    },
    AddAlgoInfo() {
      this.$refs.algoData.validate(valid => {
        if (valid) {
          addElement(this.algoData).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '创建成功',
                type: 'success'
              })
              this.GetAlgoList(this.searchInfo)
              this.dialogVisible = false
            } else {
              this.$message({
                message: '创建失败',
                type: 'error'
              });
            }
          })
        }
      });
    },
    DeleteAlgoInfo(val) {
      this.$confirm('此操作将永久删除该产品, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteElement(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetAlgoList(this.searchInfo)
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
    modStatus(val) {
      this.algoData.id = val.ID
      this.algoData.name = val.name
      this.algoData.type = val.type
      this.algoData.usage = val.usage
      this.algoData.description = val.description
      this.algoData.openType = val.openType
      this.algoData.price = val.price
      this.algoData.status = val.status == 1 ? 1 : 2
      this.algoData.isSynced = 1
      var txt = this.algoData.status == 1 ? '启用' : '禁用'
      updateElement(this.algoData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: txt + '成功',
            type: 'success'
          })
          this.GetAlgoList(this.searchInfo)
        } else {
          this.$message({
            message: txt + '失败',
            type: 'error'
          });
        }
      })
      this.algoData = []
    },
    modIsSynced() {
      this.algoData.type = 2
      this.algoData.isSynced = 2
      syncElement(this.algoData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '同步成功',
            type: 'success'
          })
          this.GetAlgoList(this.searchInfo)
        } else {
          this.$message({
            message: '同步失败',
            type: 'error'
          });
        }
      })
      this.algoData = []
    },
    UpdateAlgoInfo() {
      if (this.algoData.file != null) {
        updateElementFileInfo(this.algoData).then(response => {
          if (response.code == 0) {
            this.GetAlgoList(this.searchInfo)
          }
        })
      }
      updateElement(this.algoData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '修改成功',
            type: 'success'
          })
          this.GetAlgoList(this.searchInfo)
          this.dialogVisible = false
        } else {
          this.$message({
            message: '修改失败',
            type: 'error'
          });
        }
      })
    },
    GetAlgoList() {
      getNodeElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.algoList = response.data.list
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
      this.GetAlgoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetAlgoList()
    }
  }
}
</script>
