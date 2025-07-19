<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
          <el-form-item label="名称">
            <el-input style="width: 240px;" v-model="searchInfo.name"
                placeholder="请输入"
            />
          </el-form-item>
          <el-form-item>
            <el-button  @click="GetImageReposList" type="primary">
              查询
            </el-button>
            <el-button type="info" @click="searchInfo.name='';;GetImageReposList()">
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
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="imageReposList"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="OpenImageDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="仓库地址">
          <template #default="scope">
            <span>{{ scope.row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="HandleDeleteDialog(scope.row,'delete')">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,  prev, pager, next, sizes,jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!--新增镜像仓库-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增镜像仓库</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>

      <el-form
          ref="serverDataInfo"
          :model="dataInfo"
          :rules="rules"
          label-width="80px"
          label-position="left"
      >
        <el-form-item label="名称" prop="name">
          <el-input  style="width: 320px;" placeholder="请输入仓库名称" v-model="dataInfo.name"/>
        </el-form-item>
        <el-form-item label="仓库地址" prop="url" style="margin-top: 30px">
          <el-input  style="width: 320px;" placeholder="请输入仓库地址" v-model="dataInfo.url"/>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input  style="width: 320px;" placeholder="请输入用户名" v-model="dataInfo.username"/>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input  style="width: 320px;" placeholder="请输入密码" v-model="dataInfo.password"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--镜像列表展示-->
    <el-drawer
        v-model="showImage"
        :close-on-click-modal="true"
        :close-on-press-escape="false"
        :show-close="false"
        size="520px"
        title=""
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">容器镜像列表</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="showImage=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>


      <el-table
          v-loading="false"
          :data="imageList"
          style="width: 100%;font-size: 15px;"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="镜像版本">
          <template #default="scope">
            <el-tag
                v-for="(tag,index) in scope.row.tags"
                :key="index"
                :disable-transitions="false"
                style="margin-left: 2px"
            >
              <span> {{ tag }}</span>
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
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
  getImageReposList,
  addImageRepos, deleteImageRepos, getImageReposImageList, deleteImageReposByIds
} from '@/api/cloudpods/image/imageRepos'


export default {
  data() {
    return {
      ids:[],
      timer:null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
      },
      total: 0,
      imageReposList:[],
      imageList:[],
      dataInfo:{
        name: '',
        url:'',
        username:'',
        password:'',
      },
      rules:{
        name:[
          {required: true, message: '请输入仓库名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        url:[
          {required: true, message: '请输入仓库地址', trigger: 'blur'}
        ],
        username:[
          {required: true, message: '请输入用户名', trigger: 'blur'}
        ],
        password:[
          {required: true, message: '请输入密码', trigger: 'blur'}
        ],
      },
      createVisible:false,
      showImage:false,
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
    }
  },
  created(){
    this.GetImageReposList()
  },
  methods:{
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
      }
    },
    //开启定时器
    startInterval() {
      if (this.timer==null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetImageReposList();
        }, 5000);
      } else {
        console.log('定时器已经在运行');
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
    //删除镜像仓库
    HandleDelete(val){
      deleteImageRepos({id:val.id}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetImageReposList(this.searchInfo)
        }
      })
    },
    //删除镜像仓库
    HandleDeleteByIds(val){
      deleteImageReposByIds({ids:this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetImageReposList(this.searchInfo)
        }
      })
    },
    //新增对话框
    OpenDrawer(){
      this.dataInfo={};
      this.createVisible=true;
    },
    //镜像对话框
    OpenImageDrawer(val){
      this.showImage=true;
      getImageReposImageList({id:val.id}).then(res => {
        if (res.code == 0) {
          this.imageList = res.data.list
        }
      })
    },
    //取消
    CloseDrawer(){
      this.createVisible=false;
      this.showImage=false;
    },
    //确定
    SubmitDrawer(){
      addImageRepos(this.dataInfo).then(res=>{
        if (res.code == 0){
          this.$message({
            type: 'success',
            message:'创建成功',
          })
          this.createVisible=false;
          this.GetImageReposList();
        }
      })
    },
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
    GetImageReposList() {
      getImageReposList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData=JSON.parse(response.data).data;
          if(JSON.stringify(tempData) !== JSON.stringify(this.imageReposList)){
            this.imageReposList = JSON.parse(response.data).data
            this.total = JSON.parse(response.data).total
          }
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
      this.GetImageReposList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetImageReposList()
    }
  },
  mounted() {
    this.startInterval()
  },
  unmounted() {
    this.stopInterval()
  },
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