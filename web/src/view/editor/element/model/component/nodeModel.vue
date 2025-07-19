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
        <el-table-column label="描述" width="200px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.description }}</span>
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
            <span>{{FormatDateTime( scope.row.uploadTime)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="OpenDownDialog(scope.row)">下载</el-button>
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
  getNodeElementList,
} from "@/api/yunguan/element/element";
import {getRenterResList} from "@/api/yunguan/resource/resInfo";
import {DownloadElement} from "@/api/yunguan/cloudpods/instance/instance";


export default {
  name: "NodeModel",
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 1, //1模型 2算法 3数据集
        openType: 2,//不开放1，节点内可见可用2，开放到算力节点节点内用3，开放到算力节点且可外部用4
        status: 1,//1正常、2禁用
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      modelList: [],
      instanceList:[],
      downDialog:false,
      downloadReq:{
        privateIp:'',
        DownloadDetails:[]
      },
      modelRules:{
        privateIp: [
          {required: true, message: '请选择实例', trigger: 'blur'}
        ],
      },
      instanceInfo:{
        page:1,
        pageSize:10,
        type:1,
      },
      activeName:'baremetal',
      baremetalList:[],
      hostList:[],
      containerList:[],
      downDetail:{
        DownloadSource:2,
        ProductID:'',
        ProductType:'',
        ProductName:'',
        FileName:'',
      }
    }
  },
  created() {
    this.GetModelList()
  },
  methods: {
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
    GetModelList() {
      getNodeElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.modelList = response.data.list
          this.searchInfo.pageIndex = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
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
    },
    //下载

  }
}
</script>
