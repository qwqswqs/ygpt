<template>
  <div>
    <el-tabs>
      <el-tab-pane label="我的模型">
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
                    @click="OpenDialog">
                  添加
                </el-button>
              </el-form-item>
            </el-form>
          </div>
          <!-- 表格展示 -->
          <el-table
              :data="myDataList"
              row-key="ID"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>

            <el-table-column label="描述">
              <template #default="scope">
                <span>{{ scope.row.description }}</span>
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
                <span>{{FormatDateTime( scope.row.uploadTime) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="共享状态">
              <template #default="scope">
                <span>{{ scope.row.openType==1?'未共享':'已共享' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button type="text" @click="handleShare(scope.row)">{{scope.row.openType == 1?'共享':'取消共享'}}</el-button>
                <el-button type="text" @click="DeleteModelInfo(scope.row)">删除</el-button>
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
        </div>
      </el-tab-pane>
      <el-tab-pane label="共享模型">
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
              </el-form-item>
            </el-form>
          </div>
          <!-- 表格展示 -->
          <el-table
              :data="shareDataList"
              row-key="ID"
          >
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>

            <el-table-column label="描述">
              <template #default="scope">
                <span>{{ scope.row.description }}</span>
              </template>
            </el-table-column>
            <el-table-column label="来源">
              <template #default="scope">
                <span>{{ scope.row.source==1?'节点':'租户' }}</span>
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
                <span>{{FormatDateTime( scope.row.uploadTime) }}</span>
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
        </div>
      </el-tab-pane>
    </el-tabs>
    <el-drawer
        v-model="dialogVisible"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="false"
        size="50%">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ IsCreate ? '新增' : '编辑' }}模型</span>
          <div>
            <el-button @click="CloseDialog">取 消</el-button>
            <el-button type="primary" @click="handleConfirm">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form
          ref="dataInfo"
          :model="dataInfo"
          :rules="IsCreate?modelRules:''"
          label-width="120px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="dataInfo.name" placeholder="请输入模型名称"/>
        </el-form-item>
        <el-form-item label="标签" prop="usage">
          <el-select v-model="dataInfo.usage" multiple placeholder="请选择模型标签" style="width: 100%">
            <el-option
                v-for="(optionItem, optionIndex) in tagsList"
                :key="optionIndex"
                :label="optionItem"
                :value="optionItem"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="dataInfo.description" placeholder="请描述模型用途"/>
        </el-form-item>
        <el-form-item label="文件上传" prop="file">
          <input ref="file" type="file" @change="getFileData"/>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>
<script>
import {
  addElement,
  deleteElement,
  getShareElementList, getUserElementList, syncElement,
  updateElement,
  updateElementFileInfo
} from "@/api/yunguan/element/element";

export default {
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 1,
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      dialogTitle: '创建模型',
      myDataList: [],
      shareDataList:[],
      dialogVisible: false,
      dataInfo: {
        id: '',
        name: '',
        usage: '',
        description: '',
        status: '',
        openType: 1,
        isSynced: 1,
        startTime:'2024-12-24',
        endTime:'2024-12-24',
        source: 2, //来源 1节点 2租户上传
        file: null,
        type: 1,
        price: '0'
      },
      modelRules: {
        name: [
          {required: true, message: '请输入模型名称', trigger: 'blur'}
        ],
        usage: [
          {required: true, message: '请选择模型标签', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入模型描述', trigger: 'blur'}
        ],
        openType: [
          {required: true, message: '请选择模型开放类别', trigger: 'blur'}
        ],
        file: [
          {required: true, message: '请上传文件', trigger: 'blur'}
        ],

      },
      IsCreate: true,
      ShareStatus: ["", "未同步", "已同步"]
    }
  },
  created() {
    this.GetMyModelList()
    this.GetShareModelList()
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
      this.GetMyModelList()
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.dataInfo.file = inputFile
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
      this.dataInfo.id = val.ID
      this.dataInfo.name = val.name
      this.dataInfo.type = val.type
      this.dataInfo.usage = val.usage
      this.dataInfo.description = val.description
      this.dataInfo.openType = val.openType
      this.dataInfo.status = val.status == 1 ? 1 : 2
      this.dataInfo.isSynced = 1
      this.dataInfo.price = val.price
      this.IsCreate = false
      this.dialogVisible = true
    },
    handleShare(val) {
      let obj={}
      obj.id=val.ID
      obj.openType=val.openType==1?2:1
      updateElement(obj).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '修改成功',
            type: 'success'
          })
          this.GetMyModelList(this.searchInfo)
          this.GetShareModelList(this.searchInfo)
          this.dialogVisible = false
        } else {
          this.$message({
            message: '修改失败',
            type: 'error'
          });
        }
      })
    },
    handleConfirm() {
      if (this.IsCreate) this.AddModelInfo()
      else this.UpdateModelInfo()
    },
    AddModelInfo() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addElement(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '创建成功',
                type: 'success'
              })
              this.GetMyModelList(this.searchInfo)
              this.GetShareModelList(this.searchInfo)
              this.dialogVisible = false
            } else {
              this.$message({
                message: '创建失败',
                type: 'error'
              });
            }
          })
        }else {
          console.log('error submit!!');
          return false;
        }
      });
    },
    DeleteModelInfo(val) {
      this.$confirm('此操作将永久删除该模型, 是否继续?', '提示', {
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
            this.GetMyModelList(this.searchInfo)
            this.GetShareModelList(this.searchInfo)
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
    UpdateModelInfo() {
      if (this.dataInfo.file != null) {
        updateElementFileInfo(this.dataInfo).then(response => {
          if (response.code == 0) {
            this.GetMyModelList(this.searchInfo)
          }
        })
      }
      updateElement(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '修改成功',
            type: 'success'
          })
          this.GetMyModelList(this.searchInfo)
          this.dialogVisible = false
        } else {
          this.$message({
            message: '修改失败',
            type: 'error'
          });
        }
      })
    },
    GetMyModelList() {
      getUserElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.myDataList = response.data.list
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
    GetShareModelList() {
      getShareElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.shareDataList = response.data.list
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
      this.GetMyModelList()
      this.GetShareModelList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetMyModelList()
      this.GetShareModelList()
    }
  }
}
</script>

