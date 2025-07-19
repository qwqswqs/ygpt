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
                @click="addImageOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="imageList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column :show-overflow-tooltip='true' label="镜像名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="系统">
          <template #default="scope">
            <span>{{ scope.row.system }}</span>
          </template>
        </el-table-column>
        <el-table-column label="格式">
          <template #default="scope">
            <span>{{ scope.row.frame}}</span>
          </template>
        </el-table-column>
        <el-table-column label="语言">
          <template #default="scope">
            <span>{{ scope.row.language }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小">
          <template #default="scope">
            <span>{{  GetImageSize(scope.row.size) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="生成方式">
          <template #default="scope">
            <span>{{ scope.row.generateWay==1?'上传':'实例' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="来源">
          <template #default="scope">
            <span>{{ scope.row.sourceID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="handleDownload(scope.row)">下载</el-button>
            <el-button type="text" @click="handleShowHis(scope.row)">使用记录</el-button>
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
            layout="total,  prev, pager, next, sizes,jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
      <!--      新增镜像-->
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
          <el-form-item label="镜像名称" prop="name">
            <el-input v-model="dataInfo.name"/>
          </el-form-item>
          <el-form-item label="系统" prop="system">
            <el-input v-model="dataInfo.system"/>
          </el-form-item>
          <el-form-item label="框架" prop="frame">
            <el-input v-model="dataInfo.frame"/>
          </el-form-item>
          <el-form-item label="语言" prop="language">
            <el-input v-model="dataInfo.language"/>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="dataInfo.description"/>
          </el-form-item>
          <el-form-item label="是否公开" prop="isPublic">
            <el-switch
                v-model="dataInfo.isPublic"
                :active-value="1"
                :inactive-value="2"
                active-text="是"
                class="ml-2"
                inactive-text="否"
                inline-prompt
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
        </el-form>
      </el-drawer>
      <!--      使用记录-->
      <el-drawer
          v-model="useDialog"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
          size="60%"
      >
        <template #header>
          <div class="flex justify-between items-center">
            <span class="text-lg">使用记录</span>
            <div>
              <el-button @click="closeDialog">取 消</el-button>
            </div>
          </div>
        </template>

        <el-table
            :data="imageUseHisList"
            row-key="ID"
        >

          <el-table-column label="资源类别">
            <template #default="scope">
              <span>{{ resourceType[scope.row.resourceType-1] }}</span>
            </template>
          </el-table-column>
          <el-table-column label="启用时间">
            <template #default="scope">
              <span>{{ FormatDateTime(scope.row.startTime) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="访问地址">
            <template #default="scope">
              <span>{{ scope.row.url }}</span>
            </template>
          </el-table-column>
          <el-table-column label="访问方式">
            <template #default="scope">
              <span>{{ scope.row.visitWay }}</span>
            </template>
          </el-table-column>
          <el-table-column label="访问密码">
            <template #default="scope">
              <span>{{ scope.row.visitPwd }}</span>
            </template>
          </el-table-column>
          <el-table-column label="所属服务器ID">
            <template #default="scope">
              <span>{{ scope.row.serverID }}</span>
            </template>
          </el-table-column>
          <el-table-column label="规格描述">
            <template #default="scope">
              <span>{{ scope.row.specDesc }}</span>
            </template>
          </el-table-column>
          <el-table-column label="所属租户">
            <template #default="scope">
              <span>{{ scope.row.renterID }}</span>
            </template>
          </el-table-column>
          <el-table-column label="当前状态">
            <template #default="scope">
              <span>{{ scope.row.status==1?'空闲':'使用' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="所属子网">
            <template #default="scope">
              <span>{{ scope.row.ipID }}</span>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
              :current-page="searchUseHis.page"
              :page-size="searchUseHis.pageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="total"
              layout="total,prev, pager, next,  sizes, jumper"
              @current-change="handleCurrentUseChange"
              @size-change="handleUseSizeChange"
          />
        </div>
      </el-drawer>
    </div>
  </div>
</template>
<script>
import {
  addBackImage,
  addBackImageFile,
  updateBackImage,
  deleteBackImage,
  getBackImageList
} from "@/api/yunguan/backup/image";
import {getInfoList} from "@/api/yunguan/resource/resInfo";
import {ElMessage} from "element-plus";


export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        sourceType:2
      },
      searchUseHis:{
        page:1,
        pageSize:10,
        imageID:0,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      useDialog: false,
      dataInfo: {
        id: 0,
        system:'',
        frame:'',
        language:'',
        name:'',
        isPublic:1,
        description:'',
        generateWay:1,
        sourceType:2,
      },
      fileInfo:'',
      rules: {
        name: [
          {required: true, message: '请输入镜像名称', trigger: 'blur'}
        ],
        status: [
          {required: true, message: '请选择镜像状态', trigger: 'blur'}
        ],
        imageType: [
          {required: true, message: '请选择镜像类型', trigger: 'blur'}
        ],
      },
      resourceType:["存储","云主机","容器","裸金属"],
      imageList: [],
      imageUseHisList: [],
    }
  },
  created() {
    this.GetImageList()
  },
  methods: {
    GetImageSize(val) {
      var B = 1
      var KB = 1024
      var MB = 1024 * 1024
      var GB = 1024 * 1024 * 1024
      if (val < KB) {
        return (val / B).toFixed(1) + "B";
      } else if (val >= KB && val <= MB) {
        return (val / KB).toFixed(1) + "KB";
      } else if (val >= MB && val <= GB) {
        return (val / MB).toFixed(1) + "MB";
      } else if (val > GB) {
        return (val / GB).toFixed(1) + "GB";
      }
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
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    handleDelete(val) {
      this.$confirm('此操作将永久删除该分区, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteBackImage(val).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetImageList(this.searchInfo)
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
          message: '已取消导出'
        });
      });
    },
    handleEdit(val) {
      this.dataInfo.id = val.ID
      this.dataInfo.name = val.name
      this.dataInfo.status = val.status
      this.dataInfo.size = val.size
      this.dataInfo.imageType = val.imageType
      this.dataInfo.description = val.description
      this.dataInfo.contentDescription = val.contentDescription
      this.dataInfo.sourceType = val.sourceType
      this.IsCreated = false
      this.openDialog = true
    },
    handleDownload(val) {
      let obj={}
      obj.id=val.ID;
      try {

        // 使用 window.location.origin 构造下载 URL，这样会自动包含端口号（如果有的话）
        const downloadUrl = `${window.location.origin}/api/backup/image/downloadImage?id=${obj.id}`;

        // 直接跳转到下载链接
        window.location.href = downloadUrl;

        ElMessage.success('下载开始');
        this
      } catch (error) {
        console.error('下载失败', error);
        ElMessage.error('下载失败');
      }

    },
    handleShowHis(val){
      this.searchUseHis.imageID=val.ID
      this.GetImageUseHistory()
      this.useDialog=true
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.fileInfo = inputFile
    },
    handleSearch() {
      this.GetImageList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetImageList()
    },
    handleCurrentUseChange(val) {
      this.searchUseHis.page = val
      this.GetImageUseHistory()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetImageList()
    },
    handleUseSizeChange(val) {
      this.searchUseHis.pageSize = val
      this.GetImageUseHistory()
    },
    addImageOpenDialog() {
      this.dataInfo={}
      this.dataInfo.sourceType=2
      this.dataInfo.isPublic=1
      this.dataInfo.generateWay=1
      this.openDialog = true
      this.IsCreated = true
    },
    closeDialog() {
      this.openDialog = false
      this.useDialog = false
    },
    UpdateImage() {
      updateBackImage(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetImageList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddImage() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addBackImage(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              let obj={}
              obj.id=response.data.ID
              obj.file=this.fileInfo
              addBackImageFile(obj).then(response => {
                if (response.code == 0){
                  this.openDialog=false
                  this.GetImageList()
                }
              })
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
      if (this.IsCreated) this.AddImage()
      else this.UpdateImage()
    },
    // 获取分区列表
    GetImageList() {
      getBackImageList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.imageList = response.data.list
          this.total = response.data.total
        }
      })
    },
    GetImageUseHistory(){
      getInfoList(this.searchUseHis).then(response => {
        if (response.code == 0) {
          this.imageUseHisList = response.data.list
          this.total = response.data.total
        }
      })
    }
  }
}
</script>
<style scoped>

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}
</style>
