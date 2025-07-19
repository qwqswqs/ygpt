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
            <el-button type="primary" @click="modAllIsSynced">立即同步
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
        <el-table-column label="名称" width="150px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span v-if="scope.row.auditStatus==0">{{ '待审核' }}</span>
            <el-button type="text" @click="auditVisible=true;auditOpinion=scope.row.auditOpinion" v-else>{{ scope.row.auditStatus==-1?'驳回':'通过' }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="是否同步">
          <template #default="scope">
            <span>{{ ShareStatus[scope.row.isSynced] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="开放类别" min-width="180">
          <template #default="scope">
            <span>{{ openList[scope.row.openType] }}</span>
          </template>
        </el-table-column>

        <el-table-column label="描述" width="200px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="价格">
          <template #default="scope">
            <span>{{ scope.row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="标签" width="180px">
          <template #default="scope">
            <el-tag
                style="margin-left: 2px"
                v-for="(tag,index) in scope.row.usage"
                :key="index"
                :disable-transitions="false"
            >
              <span> {{ tag }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :fixed="'right'" label="操作" min-width="150"align="left">
          <template #default="scope">
            <el-button link type="text" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button link v-if="scope.row.auditStatus==1" type="text" :disabled="scope.row.isSynced==2" @click="modIsSynced(scope.row)">同步</el-button>
            <el-button link type="text" @click="DeleteDatasetInfo(scope.row)">删除</el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="150" :fixed="'right'"align="left">
          <template #default="scope">
            <el-button link type="text" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button link type="text" :disabled="scope.row.isSynced==2" @click="modIsSynced(scope.row)">同步</el-button>
            <el-button link type="text" @click="DeleteModelInfo(scope.row)">删除</el-button>
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
      <el-drawer
          v-model="dialogVisible"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
          size="50%">
        <template #header>
          <div class="flex justify-between items-center">
            <span class="text-lg">{{ IsCreate ? '新增' : '编辑' }}算法</span>
            <div>
              <el-button @click="CloseDialog">取 消</el-button>
              <el-button type="primary" @click="handleConfirm">确 定</el-button>
            </div>
          </div>
        </template>
        <el-form
            ref="modelData"
            :model="modelData"
            :rules="IsCreate?modelRules:''"
            label-width="120px"
        >
          <el-form-item label="名称" prop="name">
            <el-input v-model="modelData.name" placeholder="请输入算法名称"/>
          </el-form-item>
          <el-form-item label="标签" prop="usage">
            <el-select v-model="modelData.usage" multiple placeholder="请选择算法标签" style="width: 100%">
              <el-option
                  v-for="(optionItem, optionIndex) in tagsList"
                  :key="optionIndex"
                  :label="optionItem"
                  :value="optionItem"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="modelData.description" placeholder="请描述算法用途"/>
          </el-form-item>
          <el-form-item label="价格" prop="price">
            <el-input type="number" v-model="modelData.price" placeholder="请输入价格"/>
          </el-form-item>
          <el-form-item label="开放类别" prop="openType">
            <el-select v-model="modelData.openType" placeholder="请选择开放类别" style="width: 100%">
              <el-option
                  v-for="(optionItem, optionIndex) in openList"
                  v-if="optionIndex !== 0"
                  :disabled="optionIndex==0"
                  :key="optionIndex"
                  :label="optionItem"
                  :value="optionIndex"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="上架时间" prop="startTime">
            <el-date-picker style="width: 100%" type="date" value-format="YYYY-MM-DD" v-model="modelData.startTime" placeholder="请选择上架时间"/>
          </el-form-item>
          <el-form-item label="下架时间" prop="endTime">
            <el-date-picker style="width: 100%" type="date" value-format="YYYY-MM-DD"  v-model="modelData.endTime" placeholder="请选择下架时间"/>
          </el-form-item>
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
        </el-form>
      </el-drawer>
      <!--      审核意见-->
      <el-dialog v-model="auditVisible" :show-close="true" style="min-width: 500px" title="审核意见" width="40%" center>
        <span>{{auditOpinion!=''?auditOpinion:'无'}}</span>
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
import _ from "lodash";

export default {
  data() {
    return {
      auditVisible:false,
      auditOpinion:'',
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 3,
        source: 1, //来源：1节点 2租户
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      statusList: ["", "启用", "禁用"],
      openList: ['请选择开放类别', "不开放", "节点内可用", "算力可见节点内可用", "算力可见算力外可用"],
      dialogTitle: '创建算法',
      modelList: [],
      dialogVisible: false,
      modelData: {
        id: '',
        name: '',
        usage: '',
        description: '',
        status: '',
        openType: 1,
        isSynced: 1,
        startTime:'',
        endTime:'',
        source: 1, //来源 1节点 2租户上传
        file: null,
        type: 3,
        price: 0
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
      ShareStatus: ["", "未同步", "已同步"]
    }
  },
  created() {
    this.GetModelList()
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
      this.GetModelList()
    },
    uploadMaterial() {
      this.$refs.file.dispatchEvent(new MouseEvent("click"));
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.modelData.file = inputFile
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
      this.modelData = _.cloneDeep(val)
      this.modelData.id = val.ID
      this.modelData.isSynced = 1
      this.IsCreate = false
      this.dialogVisible = true
    },
    handleConfirm() {
      this.modelData.price=Number(this.modelData.price)
      if (this.IsCreate) this.AddModelInfo()
      else this.UpdateModelInfo()
    },
    AddModelInfo() {
      this.$refs.modelData.validate(valid => {
        if (valid) {
          addElement(this.modelData).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '创建成功',
                type: 'success'
              })
              this.GetModelList(this.searchInfo)
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
            this.GetModelList(this.searchInfo)
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
      this.modelData.id = val.ID
      this.modelData.name = val.name
      this.modelData.type = val.type
      this.modelData.usage = val.usage
      this.modelData.description = val.description
      this.modelData.openType = val.openType
      this.modelData.price = val.price
      this.modelData.status = val.status == 1 ? 1 : 2
      this.modelData.isSynced = 1
      var txt = this.modelData.status == 1 ? '启用' : '禁用'
      updateElement(this.modelData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: txt + '成功',
            type: 'success'
          })
          this.GetModelList(this.searchInfo)
        } else {
          this.$message({
            message: txt + '失败',
            type: 'error'
          });
        }
      })
      this.modelData = []
    },
    //同步单个数据
    modIsSynced(val) {
      this.modelData=val
      this.modelData.id=val.ID
      syncElement(this.modelData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '同步成功',
            type: 'success'
          })
        }
      })
      this.GetModelList(this.searchInfo)
    },
    //同步所有
    modAllIsSynced() {
      this.modelData.type = 3
      this.modelData.isSynced = 2
      syncElement(this.modelData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '同步成功',
            type: 'success'
          })
          this.GetModelList(this.searchInfo)
        } else {
          this.$message({
            message: '同步失败，请检查网络连接',
            type: 'error'
          });
        }
      })
    },
    UpdateModelInfo() {
      if (this.modelData.file != null) {
        updateElementFileInfo(this.modelData).then(response => {
          if (response.code == 0) {
            this.GetModelList(this.searchInfo)
          }
        })
      }
      updateElement(this.modelData).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '修改成功',
            type: 'success'
          })
          this.GetModelList(this.searchInfo)
          this.dialogVisible = false
        } else {
          this.$message({
            message: '修改失败',
            type: 'error'
          });
        }
      })
    },
    GetModelList() {
      getNodeElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.modelList = response.data.list
          this.searchInfo.pageIndex = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.total = response.data.total
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetModelList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetModelList()
    }
  }
}
</script>

