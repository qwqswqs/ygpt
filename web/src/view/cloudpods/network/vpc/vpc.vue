<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form label-width="100" ref="searchForm" :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="名称">
              <el-input style="width: 240px" v-model="searchInfo.name" placeholder="请输入"/>
            </el-form-item>
            <el-form-item label="状态">
              <el-select style="width: 240px"  v-model="searchInfo.status" placeholder="请选择" clearable>
                <el-option label="可用" value="available"/>
                <el-option label="禁用" value="unavailable"/>
              </el-select>
            </el-form-item>
            <el-form-item label="允许外网访问">
              <el-select style="width: 240px"  v-model="searchInfo.isExternal" placeholder="请选择" clearable>
                <el-option label="启用" value="eip-distgw"/>
                <el-option label="禁用" value="none"/>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item>
            <el-button  type="primary" @click="GetVpcList">查询</el-button>
            <el-button  class="button-gap" type="info" @click="searchInfo.name='';searchInfo.isExternal='';GetVpcList()">重置</el-button>
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
        <el-tooltip :disabled="!deleteList.includes(false)" content="VPC不为空,请先删除其中的网络"
                    placement="top">
          <el-button :disabled="!ids.length||deleteList.includes(false)"
                     type="grey" @click="HandleDeleteDialog('','deleteIds')">
            删除
          </el-button>
        </el-tooltip>
      </div>
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="vpcList"
          @sort-change="handleSortChange"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <div style="display: flex;align-items: center">
              <div :class="handleTagType(scope.row.status)"> </div>
              <span>{{ scope.row.status=="available"?"可用":"禁用" }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="IPv4目标网段">
          <template #default="scope">
            <span>{{ scope.row.cidr_block!=undefined?scope.row.cidr_block:'-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="允许外网访问">
          <template #default="scope">
            <span>{{ scope.row.external_access_mode!="none"?"启用":"禁用" }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP子网数量" sortable="custom">
          <template #default="scope">
            <el-button link type="text" @click="OpenIpDrawer(scope.row)">{{ scope.row.network_count }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-tooltip :disabled="scope.row.network_count==0" content="VPC不为空,请先删除其中的网络" placement="top">

            <el-button :disabled="scope.row.network_count!=0" type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除</el-button>
            </el-tooltip>
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
    <!--      新增对话框-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="520"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增VPC</span>
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
          label-width="110px"
          label-position="left"
      >
        <el-form-item label="名称" prop="name">
          <el-input style="width: 320px" placeholder="请输入名称" v-model="dataInfo.name"/>
        </el-form-item>
        <el-form-item label="IPv4目标网段" prop="url" label-width="100px" style="padding-left: 10px;margin-top: 30px">
          <el-select style="width: 320px" placeholder="请选择IPv4目标网段" v-model="dataInfo.cidr">
            <el-option v-for="(item,index) in cidrList"
                       :value="item"
                       :label="item"/>
          </el-select>
        </el-form-item>
<!--        <el-form-item label="允许外网访问" label-width="100px" style="padding-left: 10px">-->
<!--          <el-switch-->
<!--              v-model="dataInfo.isExternal"-->
<!--              :active-value="true"-->
<!--              :inactive-value="false"-->
<!--              active-text="启用"-->
<!--              class="ml-2"-->
<!--              inactive-text="禁用"-->
<!--              inline-prompt-->
<!--              style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"-->
<!--          />-->
<!--        </el-form-item>-->
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--      子网列表-->
    <el-drawer
        v-model="ipVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="false"
        :show-close="false"
        size="65%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">子网数量</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="ipVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <div v-if="ipVisible">
        <IP style="padding: 0" :vpcID="vpcID"></IP>
      </div>
    </el-drawer>
    <!--    删除弹框-->
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
          </el-icon>删除</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="deleteVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <span>此操作将永久删除该VPC, 是否继续?</span>
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
  addVpc, deleteVpc, deleteVpcByIds,
  getVpcList
} from '@/api/cloudpods/network/vpc'

import IP from '@/view/cloudpods/network/ip/ip.vue'


export default {
  components: {
    IP
  },
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      ids: [],
      deleteList: [],
      ipVisible: false,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name:'',
      },
      total: 0,
      vpcList:[],
      createVisible:false,
      dataInfo:{
        name: '',
        cidr:'',
        isExternal:true,
      },
      rules:{
        name:[
          {required: true, message: '请输入VPC名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        cidr:[
          {required: true, message: '请输入IPv4目标网段', trigger: 'blur'}
        ],
      },
      cidrList:["192.168.0.0/16","172.16.0.0/12","10.0.0.0/8"],
      vpcID: '',
    }
  },
  created(){
    this.GetVpcList()
  },
  methods:{
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
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      this.deleteVisible = true;
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].network_count>0){
          this.deleteList.push(false)
        }
      }
    },
    handleSortChange({ column, prop, order }) {
      switch (order) {
        case "ascending":
          this.searchInfo.ipNumDesc="asc";
          break;
        case "descending":
          this.searchInfo.ipNumDesc="desc";
          break;
        default:
          this.searchInfo.ipNumDesc="";
          break;
      }
      this.GetVpcList();
    },
    handleTagType(status) {
      switch (status) {
        case 1:
          return 'danger'
        case 'available':
          return 'success'
        default:
          return 'danger'
      }
    },
    //打开子网列表
    OpenIpDrawer(val){
      this.ipVisible = true;
      this.vpcID=val.id;
    },
    //删除
    HandleDelete(val){
        deleteVpc({id:val}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetVpcList(this.searchInfo)
          }
        })
    },
    //批量删除
    HandleDeleteByIds(val){
        deleteVpcByIds({ids:this.ids}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetVpcList(this.searchInfo)
          }
        })
    },
    //新增对话框
    OpenDrawer(){
      this.dataInfo={};
      this.dataInfo.isExternal=true;
      this.createVisible=true;
    },
    //取消
    CloseDrawer(){
      this.createVisible=false;
    },
    //确定
    SubmitDrawer(){
      addVpc(this.dataInfo).then(res=>{
        if (res.code == 0){
          this.$message({
            type: 'success',
            message:'创建成功',
          })
          this.createVisible=false;
          this.GetVpcList();
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
    GetVpcList() {
      getVpcList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.vpcList = response.data.data
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
      this.GetVpcList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetVpcList()
    }
  }

}
</script>

<style>
.danger {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(245, 63, 63, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}

.success {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(0, 180, 42, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}
</style>