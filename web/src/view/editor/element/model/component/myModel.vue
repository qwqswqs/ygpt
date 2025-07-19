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
            <el-button type="text" @click="OpenDownDialog(scope.row)">下载</el-button>
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
      <el-dialog
          v-model="dialogVisible"
          :title="IsCreate? '创建模型':'编辑模型'"
          width="40%"
      >
        <el-form
            ref="modelData"
            :model="modelData"
            :rules="IsCreate?modelRules:''"
            label-width="120px"
        >
          <el-form-item label="名称" prop="name">
            <el-input v-model="modelData.name" placeholder="请输入模型名称"/>
          </el-form-item>
          <el-form-item label="标签" prop="usage">
            <el-select v-model="modelData.usage" multiple placeholder="请选择模型标签" style="width: 100%">
              <el-option
                  v-for="(optionItem, optionIndex) in tagsList"
                  :key="optionIndex"
                  :label="optionItem"
                  :value="optionItem"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="modelData.description" placeholder="请描述模型用途"/>
          </el-form-item>
          <el-form-item label="上传类型" prop="typed">
            <el-select v-model="modelData.uploadType" placeholder="请选择上传类型" style="margin-top: -12px !important;">
              <el-option label="本地上传" :value="1" />
              <el-option label="文件链接" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="modelData.uploadType===1" label="附件" prop="attachment">
            <el-upload
                :action="uploadUrl"
                :file-list="modelData.files"
                multiple
                @success="handleUploadSuccess"
            >
              <el-button type="primary">
                上传文件
              </el-button>
            </el-upload>
          </el-form-item>
          <el-form-item v-if="modelData.uploadType===2" label="文件链接" prop="url">
            <el-input v-model="modelData.filePath" placeholder="请输入文件链接" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer" style="margin: 0 auto">
          <el-button type="primary" @click="handleConfirm">确 定</el-button>
          <el-button @click="CloseDialog">取 消</el-button>
        </div>
      </el-dialog>
    </div>
<!--    下载对话框-->
    <el-dialog
        v-model="downDialog"
        title="下载"
        width="40%"
    >
      <el-tabs v-model="activeName" @tab-change="handleTabsClick">
        <el-tab-pane :label="'裸金属'" name="baremetal">
        </el-tab-pane>
        <el-tab-pane :label="'云主机'" name="virtual">
        </el-tab-pane>
        <el-tab-pane :label="'容器'" name="container">
        </el-tab-pane>
      </el-tabs>
      <el-form
          style="margin-top: 10px"
          ref="downloadReq"
          :model="downloadReq"
          :rules="modelRules"
          label-width="120px"
      >
        <el-form-item label="实例" prop="privateIp">
          <el-select v-model="downloadReq.privateIp" placeholder="请选择实例">
            <el-option
                v-for="(item,index) in instanceList"
                :value="item.privateIp"
                :label="'名称:'+item.name+' IP地址:'+item.privateIp"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="HandleDownloadElem">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
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
import {getRenterResList} from "@/api/yunguan/resource/resInfo";
import {DownloadElement} from "@/api/yunguan/cloudpods/instance/instance";
import _ from "lodash";
import {getBaseUrl} from "@/utils/format";

export default {
  name: "MyModel",
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
      modelList: [],
      dialogVisible: false,
      modelData: {
        id: '',
        name: '',
        usage: '',
        description: '',
        status: '',
        source: 2,
        openType: 1,
        isSynced: 1,
        file: null,
        type: 1,
        price: '0',
        fileName: '',
        fileSize: '',
        files:[],
        url:'',
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

        privateIp: [
          {required: true, message: '请选择实例', trigger: 'blur'}
        ],

      },
      uploadUrl: '', // 初始化为空字符串
      IsCreate: true,
      instanceList:[],
      downDialog:false,
      downloadReq:{
        privateIp:'',
        DownloadDetails:[]
      },
      instanceInfo:{
        page:1,
        pageSize:10,
        type:1,
      },
      activeName:'baremetal',
      downDetail:{
        DownloadSource:2,
        ProductID:'',
        ProductType:'',
        ProductName:'',
        FileName:'',
      },
      baremetalList:[],
      hostList:[],
      containerList:[],
    }
  },
  created() {
    this.GetModelList()
  },
  computed: {
    uploadUrl() {
      return `${getBaseUrl()}product/elem/upload`; // 动态生成上传接口地址
    }
  },
  methods: {
    //文件上传成功
    handleUploadSuccess(response, file, fileList) {
      if (response.code === 0) {
        this.$message({
          message: '文件上传成功',
          type: 'success'
        });
        // 处理上传成功后的逻辑，例如更新文件列表或显示上传结果
        this.modelData.fileName = response.data.file.Filename;
        this.modelData.fileSize = response.data.file.Size;
      } else {
        this.$message({
          message: '文件上传失败',
          type: 'error'
        });
      }
    },
    //tab切换
    handleTabsClick(){
      switch (this.activeName) {
        case "baremetal":
          this.instanceInfo.type = 1
          this.instanceList=this.baremetalList;
          break;
        case "virtual":
          this.instanceInfo.type = 2
          this.instanceList=this.hostList;
          break;
        case "container":
          this.instanceInfo.type = 3
          this.instanceList=this.containerList;
          break;
      }
      this.downloadReq.privateIp=''
    },
    //获取所有实例列表
    GetRenterInfoList() {
      this.instanceInfo.type = 1;
      getRenterResList(this.instanceInfo).then(response => {
        if (response.code == 0) {
          this.baremetalList = response.data.list.filter(item=>item.status=="running")
        }
      })
      this.instanceInfo.type = 2;
      getRenterResList(this.instanceInfo).then(response => {
        if (response.code == 0) {
          this.hostList = response.data.list.filter(item=>item.status=="running")
        }
      })
      this.instanceInfo.type = 3;
      getRenterResList(this.instanceInfo).then(response => {
        if (response.code == 0) {
          this.containerList = response.data.list.filter(item=>item.status=="Running")
        }
      })
    },
    //打开下载对话框
    OpenDownDialog(val){
      this.GetRenterInfoList()
      this.downloadReq.privateIp=''
      this.downDetail.FileName=val.fileName;
      this.downDetail.ProductID=val.ID;
      this.downDetail.ProductType=val.type;
      this.downDetail.ProductName=val.name;
      this.downloadReq.DownloadDetails.push(this.downDetail)
      this.downDialog=true
    },
    //下载
    HandleDownloadElem(){
      this.$refs.downloadReq.validate(valid => {
        if (valid) {
          DownloadElement(this.downloadReq).then(res=>{
            if (res.code == 0){
              this.$message({
                type: "success",
                message:'下载成功',
              })
              this.downDialog = false
            }
          })
        }
      })
    },
    CloseDialog() {
      this.dialogVisible = false
      this.downDialog = false
    },
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
      this.GetModelList()
    },
    uploadMaterial() {
      this.$refs.file.dispatchEvent(new MouseEvent("click"));
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
      this.GetModelList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetModelList()
    }
  }
}
</script>
