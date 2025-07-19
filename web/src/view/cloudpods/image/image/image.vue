<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo"  style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="名称">
              <el-input style="width: 240px;" v-model="searchInfo.name" placeholder="请输入"></el-input>
            </el-form-item>
            <el-form-item label="系统">
              <el-select style="width: 240px;" v-model="searchInfo.distributions" placeholder="请选择" clearable>
                <el-option label="CentOS" value="CentOS"/>
                <el-option label="Ubuntu" value="Ubuntu"/>
                <el-option label="Debian" value="Debian"/>
              </el-select>
            </el-form-item>
            <el-form-item label="CPU架构">
              <el-select style="width: 240px;" v-model="searchInfo.osArch" placeholder="请选择" clearable>
                <el-option label="x86_64" value="x86"/>
                <el-option label="aarch64" value="arm"/>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item>
            <el-button type="primary" @click="GetImageList">
              查询
            </el-button>
            <el-button type="info" @click="searchInfo.name='';searchInfo.osArch='';searchInfo.distributions='';GetImageList()">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button @click="OpenDrawer"
                   type="primary">
          新增
        </el-button>
        <el-button :disabled="!ids.length" type="grey" @click="HandleDeleteDialog('','deleteIds')">删除</el-button>
      </div>
      <!-- 表格展示 -->
      <el-table
          @selection-change="handleSelectionChangeData" :data="imageList" @sort-change="handleSortChange"  :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
        <el-table-column type="selection" width="30" />
        <el-table-column :show-overflow-tooltip='true' label="镜像名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="标签">
          <template #default="scope">
            <el-tag
                :type="getTagType(scope.row.description)"
                :class="getTagClass(scope.row.description)"
                v-if="scope.row.description"
            >
              {{ scope.row.description }}
            </el-tag>
            <el-tag type="info" v-else>无</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="系统">
          <template #default="scope">
            <span>{{ scope.row.properties!=undefined?scope.row.properties.distro!=undefined?scope.row.properties.distro:scope.row.properties.os_type:'未知' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU架构">
          <template #default="scope">
            <span>{{ scope.row.os_arch }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <div v-if="scope.row.status == 'active'" >
              <i class="status-dot" />
              <span>{{ imageStatus[scope.row.status] }}</span>
            </div>
            <div v-else-if="scope.row.status == 'killed'||scope.row.status == 'save_fail'" >
              <i class="status-dot" style="background-color: red"/>
              <span>{{ imageStatus[scope.row.status] }}</span>
            </div>
            <div v-else>
              <el-icon class="status-dot-other">
                <Loading />
              </el-icon>
              <span>{{ imageStatus[scope.row.status] +(scope.row.progress!=0?(' ('+Number(scope.row.progress).toFixed(2)+'%)'):'')}}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="镜像大小" sortable="custom" >
          <template #default="scope">
            <span>{{ GetImageSize(scope.row.size) }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column label="系统">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{-->
        <!--                scope.row.properties.os_distribution != null ? (scope.row.properties.os_distribution + ' ' + scope.row.properties.os_version) : scope.row.properties.os_type-->
        <!--              }}</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column label="格式">
          <template #default="scope">
            <span>{{ scope.row.disk_format }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="HandleEdit(scope.row)">编辑</el-button>
            <el-button :disabled="scope.row.status=='saving'" type="text" link @click="HandleDeleteDialog(scope.row,'delete')">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="searchInfo.pageIndex" :page-size="searchInfo.pageSize"
                       :page-sizes="[10, 30, 50, 100]" :total="total" layout="total,  prev, pager, next, sizes,jumper"
                       @current-change="handleCurrentChange" @size-change="handleSizeChange" />
      </div>
    </div>

    <!--新增镜像-->
    <el-drawer v-model="createVisible" :close-on-click-modal="true" :close-on-press-escape="true" :show-close="false"
               size="520px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">上传镜像</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>

      <el-form ref="formRef" :model="dataInfo" :rules="rules" label-position="left" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input style="width: 320px;" placeholder="请输入镜像名称" @input="handleInput" v-model="dataInfo.name" />
        </el-form-item>
        <el-form-item label="上传方式" prop="uploadType" style="margin-top: 30px">
          <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
            <el-radio-group v-model="dataInfo.uploadType">
              <el-radio-button :value="'url'">镜像Url</el-radio-button>
              <el-radio-button :value="'file'">镜像文件</el-radio-button>
            </el-radio-group>
          </div>
        </el-form-item>
        <el-form-item label="CPU架构" prop="osArch">
          <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
            <el-radio-group v-model="dataInfo.osArch">
              <el-radio-button :value="'x86'">x86_64</el-radio-button>
              <el-radio-button :value="'arm'">aarch64</el-radio-button>
            </el-radio-group>
          </div>

        </el-form-item>
        <el-form-item label="操作系统" prop="osType">
          <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
            <el-radio-group v-model="dataInfo.osType">
              <el-radio-button :value="'Linux'">Linux</el-radio-button>
              <el-radio-button :value="'Windows'">Windows</el-radio-button>
            </el-radio-group>
          </div>

        </el-form-item>
        <el-form-item v-if="dataInfo.uploadType=='file'" label="镜像文件" prop="imageFile">
          <input style="width: 320px;" ref="file" type="file" @change="getFileData" />
        </el-form-item>
        <el-form-item v-if="dataInfo.uploadType=='url'" label="镜像Url" prop="url">
          <el-input style="width: 320px;" placeholder="请输入镜像Url"v-model="dataInfo.url" />
        </el-form-item>
      </el-form>

      <!-- 新增的上传进度条 -->
      <div v-if="uploading" class="mt-4 px-4">
        <el-progress :percentage="uploadProgress" status="active"></el-progress>
        <div class="text-center mt-2 text-sm">上传进度: {{ uploadProgress }}%</div>
      </div>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer" :disabled="uploading">确 定</el-button>
      </div>
    </el-drawer>
    <!--    修改镜像-->
    <el-dialog v-model="modVisible" :show-close="false" width="520px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">修改镜像</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="modVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form ref="resetData" :model="resetData" :rules="editRule" label-width="120px">
        <el-form-item label="镜像名称" prop="name">
          <el-input style="width: 320px;" v-model="resetData.name" placeholder="请输入镜像名称" />
        </el-form-item>
        <el-form-item label="镜像标签" label-width="110" style="padding-left: 10px">

          <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
            <el-radio-group v-model="resetData.description">
              <el-radio-button :value="'裸金属'">裸金属</el-radio-button>
              <el-radio-button :value="'其他'">其他</el-radio-button>
            </el-radio-group>
          </div>

        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="primary" @click="HandleConfirm">确 定</el-button>
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog
        v-model="deleteVisible"
        :show-close="false"
        width="520px"
        @close="deleteVisible=false"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg"
              style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);display: flex;align-items: center;">
          <el-icon :style="{ color: 'rgba(255, 125, 0, 1)', fontSize: '1.25em',marginRight: '8px' }">
            <WarningFilled/>
          </el-icon>{{ dialogTitle }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="deleteVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <span>{{ dialogInfo }}</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="info" @click="deleteVisible = false">取消</el-button>
          <el-button
              type="primary"
              @click="confirmDelete"
          >
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script>
import {
  deleteImage,
  deleteImageByIds,
  getImageByName,
  getImageList,
  updateImageTag,
  uploadImage
} from "@/api/cloudpods/image/image"
import { onUnmounted } from "vue";
import statusType from "@/locales/zh-CN.json";
import axios from 'axios';

export default {
  data() {
    return {
      ids:[],
      timer:null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name:'',
      },
      imageStatus: statusType.status.image,
      total: 0,
      openDialog: false,
      dataInfo: {
        name: '',
        osArch: '',
        osType: '',
        uploadType: '',
        url:'',
        image: null,
      },
      editRule:{
        name: [
          { required: true, message: '请输入镜像名称', trigger: 'blur' },
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9\.\-]{1,19}$/,
            message: '名称必须以字母开头，只能包含英文字符、数字、点(.)和连字符(-)，且长度应在2到20个字符之间',
            trigger: 'blur'
          }
        ],
      },
      rules: {
        name: [
          { required: true, message: '请输入镜像名称', trigger: 'blur' },
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9\.\-]{1,19}$/,
            message: '名称必须以字母开头，只能包含英文字符、数字、点(.)和连字符(-)，且长度应在2到20个字符之间',
            trigger: 'blur'
          },
          { validator: this.validateNameUnique, trigger: 'input' }
        ],
        osArch: [
          { required: true, message: '请选择镜像状态', trigger: 'blur' }
        ],
        uploadType: [
          { required: true, message: '请选择镜像上传方式', trigger: 'blur' }
        ],
        osType: [
          { required: true, message: '请选择镜像状态', trigger: 'blur' }
        ],
        image: [
          { required: true, message: '请选择镜像状态', trigger: 'blur' }
        ],
        imageFile: [
          { required: true, message: '请上传镜像文件', trigger: 'blur' }
        ],
        url: [
          { required: true, message: '请输入镜像URL', trigger: 'blur' },
          {
            pattern: /^https?:\/\/.+/i, // 正则表达式，匹配以http://或https://开头的字符串
            message: 'URL必须以http://或https://开头',
            trigger: 'blur'
          },
        ],
      },
      imageList: [],
      createVisible: false,
      modVisible: false,
      imageUploadProgress: 0,
      progressShow: false,
      intervalId: null,
      formRef:null,
      resetData:{
        id:'',
        name:'',
        description:'',
      },
      // 新增的上传相关状态
      uploading: false,
      uploadProgress: 0,
      uploadCancelToken: null,
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
    }
  },
  created() {
    this.GetImageList()
  },
  methods: {
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
      }
    },
    handleSortChange({ column, prop, order }) {
      switch (order) {
        case "ascending":
          this.searchInfo.sizeDesc="asc";
          break;
        case "descending":
          this.searchInfo.sizeDesc="desc";
          break;
        default:
          this.searchInfo.sizeDesc="";
          break;
      }
      this.GetImageList();
    },
    getTagType(description) {
      if (!description) return 'info'; // 无描述使用默认样式
      if (description.includes('裸金属')) return 'success'; // 裸金属使用绿色
      return 'primary'; // 其他情况使用橙色
    },
    getTagClass(description) {
      return description && !description.includes('裸金属') ? 'custom-orange-tag' : '';
    },
    //开启定时器
    startInterval() {
      if (this.timer==null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetImageList();
        }, 5000);
      }
    },
    //关闭定时器
    stopInterval() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
        console.log('定时器已清除');
      }
    },
    // 处理输入事件
    handleInput(){
      this.$refs.formRef.validateField('name');
    },
    async checkNameExists (name){
      try {
        const response = await getImageByName({ name });
        if (response.code === 0 && response.data) {
          let dataTemp = JSON.parse(response.data);
          return dataTemp !== null;
        }
        return false;
      } catch (error) {
        console.error(error);
        return false;
      }
    },
    async validateNameUnique(rule, value, callback) {
      if (!value) {
        callback(new Error('请输入名称'));
        return;
      }
      try {
        const exists = await this.checkNameExists(value);
        if (exists) {
          callback(new Error('该名称已存在'));
        } else {
          callback();
        }
      } catch (error) {
        callback(new Error('检查名称失败，请重试'));
      }
    },
    //编辑
    HandleEdit(val) {
      this.resetData.id = val.id;
      this.resetData.name = val.name;
      this.resetData.description = val.description;
      this.modVisible = true;
    },
    //编辑确定
    HandleConfirm(){
      this.$refs.resetData.validate(valid => {
        if (valid){
          updateImageTag(this.resetData).then(response => {
            if (response.code === 0) {
              this.$message({
                type: 'success',
                message:'修改成功',
              })
            }
            this.GetImageList();
            this.modVisible=false;
          })
        }
      })
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      switch (type) {
        case 'delete':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该云主机，是否继续?';
          break;
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该云主机，是否继续?';
          break;
      }
      this.deleteVisible = true;
    },
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDelete(this.deleteId)
          break;
        case 'deleteIds':
          this.HandleDeleteByIds()
          break;
      }
      this.deleteVisible=false;
    },
    //删除
    HandleDelete(val) {
      deleteImage({ id: val.id }).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetImageList(this.searchInfo)
        }
      })
    },
    //批量删除
    HandleDeleteByIds() {
      deleteImageByIds({ ids: this.ids }).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetImageList(this.searchInfo)
        }
      })
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.dataInfo.image = inputFile
    },
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
        deleteImage(val.ID).then(response => {
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
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetImageList()

    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    // 获取分区列表
    GetImageList() {
      getImageList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData=response.data.data;
          if(tempData!== this.imageList){
            this.imageList = response.data.data
            this.total = response.data.total
            this.startInterval()
          }else{
            this.stopInterval()
          }
        }
      })
    },
    // 获取分区列表
    GetImageByName() {
      getImageByName({name:this.dataInfo.name}).then(response => {
        if (response.code == 0) {
          let dataTemp = JSON.parse(response.data)
          if (dataTemp!=null) {
            return true;
          }
          return false;
        }
      })
    },
    //新增对话框
    OpenDrawer() {
      this.dataInfo = {
        name: '',
        osArch: '',
        osType: '',
        uploadType: '',
        url: '',
        image: null
      };
      this.createVisible = true;
      this.uploading = false;
      this.uploadProgress = 0;
    },
    //取消
    CloseDrawer() {
      this.createVisible = false;
      this.modVisible = false;
      // 如果上传中，取消上传
      if (this.uploading && this.uploadCancelToken) {
        this.uploadCancelToken.cancel('用户取消上传');
      }
      this.uploading = false;
    },
    //确定
    SubmitDrawer() {
      this.$refs.formRef.validate((valid) => {
        if (valid) {
          const formData = new FormData();
          formData.append('name', this.dataInfo.name);
          formData.append('image', this.dataInfo.image);
          formData.append('osArch', this.dataInfo.osArch);
          formData.append('osType', this.dataInfo.osType);
          formData.append('uploadType', this.dataInfo.uploadType);
          formData.append('url', this.dataInfo.url);

          // 开始上传
          this.uploading = true;
          this.uploadProgress = 0;

          // 使用axios的cancel token支持取消上传
          const CancelToken = axios.CancelToken;
          this.uploadCancelToken = CancelToken.source();

          // 调用上传API并监听进度
          uploadImage(formData, {
            onUploadProgress: progressEvent => {
              const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
              this.uploadProgress = percentCompleted;
            },
            cancelToken: this.uploadCancelToken.token
          }).then(res => {
            this.uploading = false;
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '上传成功',
              });
              this.createVisible = false;
              this.GetImageList(); // 刷新列表
            } else {
              this.$message({
                type: 'error',
                message: '上传失败: ' + res.message,
              });
            }
          }).catch(error => {
            this.uploading = false;
            if (this.isCancel(error)) {
              this.$message({
                type: 'info',
                message: '已取消上传',
              });
            } else {
              this.$message({
                type: 'error',
                message: '上传发生错误: ' + error.message,
              });
            }
          });
        }
      });
    },
  },
  mounted() {
    this.startInterval()
  },
  unmounted() {
    this.stopInterval()
    // 确保组件销毁时取消正在进行的上传
    if (this.uploading && this.uploadCancelToken) {
      this.uploadCancelToken.cancel('组件销毁，取消上传');
    }
  },
}
</script>
<style scoped>
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50;
  /* Green color */
  margin-right: 5px;
}
.status-dot-other {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-right: 8px;
  animation: rotate 2s linear infinite;
}

.custom-orange-tag {
  --el-tag-text-color: #ff7d00 !important;
  --el-tag-bg-color: #fff7e8 !important;
  --el-tag-border-color: #fff7e8 !important;
}

</style>