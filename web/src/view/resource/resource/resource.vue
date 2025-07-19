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
                @click="addInfoOpenDialog"
            >
              新增
            </el-button>
            <el-upload
                action
                style="margin-left: 5px"
                accept=".xlsx,.xls"
                :show-file-list="false"
                :auto-upload="false"
                :on-change="HandleImportFile"
            >
              <template #trigger>
                <el-button type="primary">导入文件</el-button>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="dataList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="资源编号">
          <template #default="scope">
            <span>{{ scope.row.ID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="资源类别">
          <template #default="scope">
            <span>{{ resourceType[scope.row.resourceType - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="启用时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.startTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="结束时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.endTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="访问地址">
          <template #default="scope">
            <span>{{ scope.row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SSH地址">
          <template #default="scope">
            <span>{{ scope.row.sshHost }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SSH端口">
          <template #default="scope">
            <span>{{ scope.row.sshPort }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SSH用户名">
          <template #default="scope">
            <span>{{ scope.row.sshUser }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属服务器ID">
          <template #default="scope">
            <span>{{ GetServerName(this.serverList, scope.row.serverID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="规格描述">
          <template #default="scope">
            <el-button @click="OpenSpecDialog(scope.row.specDesc)" type="text">详情</el-button>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <span>{{ GetRenterName(scope.row.renterID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="当前状态">
          <template #default="scope">
            <span>{{ statusType[scope.row.status - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属子网">
          <template #default="scope">
            <span>{{ GetDataName(this.subList, scope.row.ipID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所用镜像">
          <template #default="scope">
            <span>{{ GetDataName(this.imageList, scope.row.imageID) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作"align="left" :min-width="150" :fixed="'right'">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
<!--            <el-button type="text" @click="handleMonitor(scope.row)">监控</el-button>-->
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
            <span class="text-lg">{{ IsCreated ? '新增' : '编辑' }}资源</span>
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
          <el-form-item label="资源类别" prop="resourceType">
            <el-select v-model="dataInfo.resourceType">
              <el-option
                  v-for="(item,index) in resourceType"
                  :key="index"
                  :label="item"
                  :value="index+1"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="启用时间" prop="startTime">
            <el-date-picker type="datetime" v-model="dataInfo.startTime" style="width:100%"/>
          </el-form-item>
          <el-form-item label="结束时间" prop="endTime">
            <el-date-picker type="datetime" v-model="dataInfo.endTime" style="width:100%"/>
          </el-form-item>
          <el-form-item label="访问地址" prop="url">
            <el-input v-model="dataInfo.url"/>
          </el-form-item>
          <el-form-item label="SSH地址" prop="sshHost">
            <el-input v-model="dataInfo.sshHost"/>
          </el-form-item>
          <el-form-item label="SSH端口" prop="sshPort">
            <el-input v-model="dataInfo.sshPort"/>
          </el-form-item>
          <el-form-item label="SSH用户名" prop="sshUser">
            <el-input v-model="dataInfo.sshUser"/>
          </el-form-item>
          <el-form-item label="SSH密码" prop="sshPwd">
            <el-input v-model="dataInfo.sshPwd"/>
          </el-form-item>
          <el-form-item label="所属服务器" prop="serverID">
            <el-select v-model="dataInfo.serverID">
              <el-option :value="0" :label="'无'"/>
              <el-option
                  v-for="(item,index) in serverList"
                  :key="index"
                  :label="'cpu类型: '+item.cpuType"
                  :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所属租户" prop="renterID">
            <el-select v-model="dataInfo.renterID">
              <el-option :value="-1" :label="'无'"/>
              <el-option
                  v-for="(item,index) in renterList"
                  :key="index"
                  :label="item.userName"
                  :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="dataInfo.status">
              <el-option
                  v-for="(item,index) in statusType"
                  :key="index"
                  :label="item"
                  :value="index+1"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所属子网" prop="ipID">
            <el-select v-model="dataInfo.ipID">
              <el-option :value="0" :label="'无'"/>
              <el-option
                  v-for="(item,index) in subList"
                  :key="index"
                  :label="item.name"
                  :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所用镜像" prop="imageID">
            <el-select v-model="dataInfo.imageID">
              <el-option :label="'无'" :value="0"/>
              <el-option
                  v-for="(item,index) in imageList"
                  :key="index"
                  :label="item.name"
                  :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="规格描述" prop="specDesc">
            <el-button type="text" @click="OpenSpecConfigDialog">定义</el-button>
          </el-form-item>
          <el-form-item label="规格描述">
            <el-table :data="dataInfo.specDesc">
              <el-table-column label="规格名称" prop="name"/>
              <el-table-column label="规格内容" prop="value"/>
            </el-table>
          </el-form-item>
        </el-form>
      </el-drawer>
      <el-dialog v-model="specDialog" :before-close="closeDialog" :title="'规格详情'" width="500px">
        <el-table :data="specData">
          <el-table-column label="规格名称" prop="name"/>
          <el-table-column label="规格内容" prop="value"/>
        </el-table>
      </el-dialog>
      <el-dialog v-model="specConfigDialog" :before-close="CloseSpecConfigDialog" :title="'规格定义'" width="500px">
        <el-table
            :data="configList"
        >
          <el-table-column label="名称">
            <template #default="scope">
              <span>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column label="规格内容">
            <template #default="scope">
              <el-select v-model="scope.row.value[0]">
                <el-option
                    v-for="(item, index) in configAllList[scope.$index].value"
                    :key="index"
                    :label="item"
                    :value="item"/>
              </el-select>
            </template>
          </el-table-column>
        </el-table>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="CloseSpecConfigDialog">取 消</el-button>
          <el-button type="warning" @click="ConfirmConfig">确 定</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>
<script>
import {addInfo, deleteInfo, getInfoList, updateInfo} from "@/api/yunguan/resource/resInfo";
import {getAllServerList} from "@/api/yunguan/resource/resDevice";
import {getAllRenterList} from "@/api/yunguan/run/renter";
import {getBackImageAllList} from "@/api/yunguan/backup/image";
import {getAllSubList} from "@/api/yunguan/network/netSub";
import {getAllConfigInfoList} from "@/api/yunguan/product/productConfig";
import _ from "lodash";
import * as XLSX from "xlsx";
import {monitorInstance} from "@/api/yunguan/cloudpods/instance/instance";

export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      dataInfo: {
        id: 0,
        resourceType: '',
        serverID: '',
        startTime: '',
        endTime:'',
        specDesc: [],
        renterID: '',
        status: '',
        ipID: '',
        imageID: '',
        url: '',
        sshHost: '',
        sshPort: '',
        sshUser: '',
        sshPwd: '',
      },
      resourceType:  ["裸金属", "云主机", "容器", "存储"],
      statusType: ["空闲", "使用", "损坏", "注销"],
      rules: {
        resourceType: [
          {required: true, message: '请选择资源类型', trigger: 'blur'}
        ],
        startTime: [
          {required: true, message: '请选择开始时间', trigger: 'blur'}
        ],
        endTime: [
          {required: true, message: '请选择结束时间', trigger: 'blur'}
        ],
        serverID: [
          {required: true, message: '请选择所属服务器', trigger: 'blur'}
        ],
        renterID: [
          {required: true, message: '请选择所属租户', trigger: 'blur'}
        ],
        ipID: [
          {required: true, message: '请选择子网', trigger: 'blur'}
        ],
        status: [
          {required: true, message: '请选择状态', trigger: 'blur'}
        ],
        imageID: [
          {required: true, message: '请选择镜像', trigger: 'blur'}
        ],
      },
      dataList: [],
      serverList: [],
      renterList: [],
      subList: [],
      imageList: [],
      specDialog: false,
      specConfigDialog: false,
      specData: [],
      configList:[],
      configAllList:[],

      //导入表格相关数据
      titleName:["名称","类别","时间"],
      tableTitle:[],
      tableData:[],
    }
  },
  created() {
    this.GetInfoList()
    this.GetImageList()
    this.GetRenterList()
    this.GetSubList()
    this.GetServerList()

  },
  methods: {
    // 把文件按照二进制进行读取
    ReadFileBinary(file){
      return new Promise((resolve) => {
        let reader = new FileReader();
        reader.readAsBinaryString(file);
        reader.onload = (ev) => {
          resolve(ev.target.result);
        };
      });
    },
    //导入文件
    async HandleImportFile(file) {
      this.tableTitle = [];
      this.tableData = [];
      let fileData = file.raw;
      let data = await this.ReadFileBinary(fileData)
      let workbook = XLSX.read(data, {type: "binary"});
      console.log("workbook", workbook); //这里就是可读取的文件了
      // 最后把数据转成json格式的
      let worksheet = workbook.Sheets[workbook.SheetNames[0]]; //这里是表格的名字,这里取第一个表格,1就是第二个表格数据
      console.log("worksheet", worksheet);
      //将得到的worksheet转化为json格式
      data = XLSX.utils.sheet_to_json(worksheet);

      console.log("data", data);

      // 这个是表格所有的标题
      for (const key in data[0]) {
        this.tableTitle.push(key); //获取的是标题,即每个对象的键名
      }
      console.log("title", this.tableTitle);
      data.forEach((item, index) => {
        let obj = {};
        this.tableTitle.forEach((item2, index2) => {
          console.log("item", item[item2]);
          obj[this.titleName[index2]] = item[item2];
        });
        this.tableData.push(obj);
      });
      console.log("tableData.value", this.tableData);
    },
    //获取规格配置
    GetAllConfigInfoList() {
      getAllConfigInfoList().then((response) => {
        if (response.code == 0) {
          this.configList = response.data.list
        } else {
          this.$message({
            message: "获取失败",
            type: "error",
          });
        }
      });
    },
    OpenSpecDialog(val) {
      this.specDialog = true
      this.specData = JSON.parse(val)
    },
    OpenSpecConfigDialog(){
      this.specConfigDialog=true
      this.configList=this.configList.filter(item => item.type.includes(this.dataInfo.resourceType))
      this.configAllList=_.cloneDeep(this.configList)
      if (this.configList.length>0){
        for (let i = 0; i < this.configList.length; i++) {
          this.configList[i].value=[this.configList[i].value[0]]
        }
      }
      console.log(this.configAllList)
      console.log(this.configList)
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
    //获取租户名称
    handleDelete(val) {
      this.$confirm('此操作将永久删除该分区, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteInfo(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetInfoList(this.searchInfo)
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
    //监控
    handleMonitor(val){
      let obj={}
      obj.sshHost=val.sshHost
      monitorInstance(obj).then(res=>{
        if (res.code == 0){

        }
      })
    },
    handleEdit(val) {
      this.dataInfo = val
      this.dataInfo.specDesc=JSON.parse(val.specDesc)
      this.IsCreated = false
      this.openDialog = true
      this.GetAllConfigInfoList()
    },
    handleSearch() {
      this.GetInfoList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetInfoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetInfoList()
    },
    addInfoOpenDialog() {
      this.dataInfo = {}
      this.specData=[],
      this.openDialog = true
      this.IsCreated = true
      this.GetAllConfigInfoList()
    },
    closeDialog() {
      this.openDialog = false
      this.specDialog = false
    },
    CloseSpecConfigDialog(){
      this.specConfigDialog=false
    },
    ConfirmConfig(){
      this.dataInfo.specDesc=this.configList
      this.specConfigDialog=false
    },
    UpdateInfo() {
      this.dataInfo.sshPort = Number(this.dataInfo.sshPort)
      this.dataInfo.specDesc = JSON.stringify(this.dataInfo.specDesc)
      updateInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetInfoList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    // 新增分区
    AddInfo() {
      this.dataInfo.sshPort = Number(this.dataInfo.sshPort)
      this.dataInfo.specDesc = JSON.stringify(this.dataInfo.specDesc)
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addInfo(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false
              this.GetInfoList()
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
      if (this.IsCreated) this.AddInfo()
      else this.UpdateInfo()
    },
    // 获取分区列表
    GetInfoList() {
      getInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取镜像列表
    GetImageList() {
      getBackImageAllList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.imageList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取服务器列表
    GetServerList() {
      getAllServerList().then(response => {
        if (response.code == 0) {
          this.serverList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取租户列表
    GetRenterList() {
      getAllRenterList().then(response => {
        if (response.code == 0) {
          this.renterList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取子网列表
    GetSubList() {
      getAllSubList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.subList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    GetRenterName(id) {
      const image = this.renterList.find(t => t.ID === id);
      return image ? image.userName : null;
    },
    GetDataName(list, id) {
      const image = list.find(t => t.ID === id);
      return image ? image.name : null;
    },
    GetServerName(list, id) {
      const image = list.find(t => t.ID === id);
      return image ? image.snNum : null;
    }
  }
}
</script>
