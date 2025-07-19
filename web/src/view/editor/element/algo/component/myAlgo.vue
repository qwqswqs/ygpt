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
                placeholder="请输入算法名称"
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
          :data="modelList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column label="描述">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.description }}</span>
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
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ GetFileName(scope.row.fileName) }}</span>
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
        <el-table-column label="操作">
          <template #default="scope">
<!--            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>-->
            <el-button type="text" @click="DeleteAlgoInfo(scope.row)">删除</el-button>
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
      <el-dialog
          v-model="dialogVisible"
          :title="IsCreate? '创建算法':'编辑算法'"
          width="40%"
      >
        <el-form
            ref="algoData"
            :model="algoData"
            :rules="IsCreate?modelRules:''"
            label-width="120px"
        >
          <el-form-item label="名称" prop="name">
            <el-input v-model="algoData.name" placeholder="请输入算法名称"/>
          </el-form-item>
          <el-form-item label="标签" prop="usage">
            <el-select v-model="algoData.usage" multiple placeholder="请选择算法标签" style="width: 100%">
              <el-option
                  v-for="(optionItem, optionIndex) in tagsList"
                  :key="optionIndex"
                  :label="optionItem"
                  :value="optionItem"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="algoData.description" placeholder="请描述算法用途"/>
          </el-form-item>
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer" style="margin: 0 auto">
          <el-button type="primary" @click="handleConfirm">确 定</el-button>
          <el-button @click="CloseDialog">取 消</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>
<script>
import {
  addElement,
  deleteElement,
  getUserElementList,
  updateElement,
  updateElementFileInfo
} from "@/api/yunguan/element/element";
import _ from "lodash";

export default {
  name: "MyAlgo",
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 3,
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      dialogTitle: '创建算法',
      modelList: [],
      dialogVisible: false,
      algoData: {
        id: '',
        name: '',
        usage: '',
        description: '',
        status: '',
        source: 2,
        openType: 1,
        isSynced: 1,
        file: null,
        type: 3,
        price: '0'
      },
      modelRules: {
        name: [
          {required: true, message: '请输入算法名称', trigger: 'blur'}
        ],
        usage: [
          {required: true, message: '请选择算法标签', trigger: 'blur'}
        ],
        description: [
          {required: true, message: '请输入算法描述', trigger: 'blur'}
        ],
        openType: [
          {required: true, message: '请选择算法开放类别', trigger: 'blur'}
        ],
        file: [
          {required: true, message: '请上传文件', trigger: 'blur'}
        ],

      },
      IsCreate: true,
    }
  },
  created() {
    this.GetAlgoList()
  },
  methods: {
    FormatDateTime(dateString){
      const date = new Date(dateString);
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
      this.algoData = _.cloneDeep(val)
      this.algoData.id = val.ID
      this.algoData.isSynced = 1
      this.IsCreate = false
      this.dialogVisible = true
    },
    handleConfirm() {
      this.algoData.price=Number(this.algoData.price)
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
      this.$confirm('此操作将永久删除该算法, 是否继续?', '提示', {
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
      getUserElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.modelList = response.data.list
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
