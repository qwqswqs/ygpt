<template>
  <div>
<!--    表格-->
<!--    <el-tabs v-model="activeName" @click="handleTabsClick">-->
<!--      <el-tab-pane :label="'全部('+searchInfo.totalNum+')'" name="all"></el-tab-pane>-->
<!--      <el-tab-pane :label="'待分配('+searchInfo.unAssignNum+')'" name="unAssign"></el-tab-pane>-->
<!--      <el-tab-pane :label="'待处理('+searchInfo.unHandleNum+')'" name="unHandle"></el-tab-pane>-->
<!--      <el-tab-pane :label="'已处理('+searchInfo.handleNum+')'" name="handled"></el-tab-pane>-->
<!--    </el-tabs>-->
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
                placeholder="请输入名称"
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
          :data="ticketList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="工单类别">
          <template #default="scope">
            <span>{{ ticketType[scope.row.ticketType-1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <span>{{ GetRenterName(scope.row.renterID)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="scope">
            <el-button type="text" @click="GetDescDetail(scope.row.description)">查看详情</el-button>
          </template>
        </el-table-column>
<!--        <el-table-column label="分配人">-->
<!--          <template #default="scope">-->
<!--            <span>{{ GetUserName(scope.row.assignPerson) }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--        <el-table-column label="处理人">-->
<!--          <template #default="scope">-->
<!--            <span>{{ GetUserName(scope.row.handlePerson) }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--        <el-table-column label="状态">-->
<!--          <template #default="scope">-->
<!--            <el-tag :type="tagStatus[scope.row.status]">{{ handleStatus[scope.row.status] }}</el-tag>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column label="创建时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
<!--        <el-table-column label="分配时间">-->
<!--          <template #default="scope">-->
<!--            <span>{{ FormatDateTime(scope.row.assignTime) }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column label="操作"align="left">
          <template #default="scope">
<!--            <el-button type="text" v-if="scope.row.status!=3"  @click="handleAssign(scope.row)">分配</el-button>-->
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
    </div>
<!--    描述详情对话框-->
    <el-dialog v-model="descDialog" :before-close="closeDialog" :title="'分配详情'" width="500px">
      <el-descriptions
          :column="2"
          style="margin-top: 10px"
          title="基本信息"
          border
      >
        <el-descriptions-item label="租户">
          <span>{{detailData.phone}}</span>
        </el-descriptions-item>
        <el-descriptions-item label="资源类型">
          <span>{{resourceType[detailData.type-1]}}</span>
        </el-descriptions-item>
        <el-descriptions-item label="计费方式">
          <span>{{GetBillingWay(detailData.jfWay) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="使用时长">
          <span>{{detailData.useTime}}</span>
        </el-descriptions-item>
      </el-descriptions>
      <el-descriptions
          :column="1"
          style="margin-top: 10px"
          title="规格配置"
      >
        <el-descriptions-item>
          <el-table
              :data="specDetail"
          >
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="规格内容">
              <template #default="scope">
                <span>{{ scope.row.value[scope.row.valueIndex] }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
<!--    处理分配对话框-->
    <el-dialog v-model="handleDialog" :before-close="closeDialog" :title="'工单分配'" width="500px">
      <el-form>
      <el-form-item>
        <el-select v-model="dataInfo.handlePerson" placeholder="请选择处理人">
          <el-option
            v-for="(item,index) in userList"
            :label="item.userName"
            :value="item.ID"/>
        </el-select>
      </el-form-item></el-form>
      <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="warning" @click="ConfirmAssign">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import {
  getSysTicketList,
  deleteSysTicket,
  assignSysTicket,
} from "@/api/yunguan/system/ticket";
import {getAllRenterList,getAllUserList} from "@/api/yunguan/run/renter";
import * as XLSX from 'xlsx';

export default {
  data() {
    return {
      activeName:'all',
      searchInfo: {
        page: 1,
        pageSize: 10,
        type:2,
        status:0,
        totalNum:0,
        unAssignNum:0,
        unHandleNum:0,
        handleNum:0,
      },
      total: 0,
      tagStatus:['','warning','error','success'],
      handleStatus:["","待分配","未处理","已处理"],
      ticketList: [],
      userList:[],
      renterList:[],
      dataInfo:{
        id:'',
        handlePerson:'',
      },
      ticketType:["故障","资源分配","资源扩大"],
      descDialog:false,
      handleDialog:false,
      specDetail:[],
      detailData:'',
      resourceType: ["裸金属", "云主机", "容器", "存储"],
      billingWay:[
        {value:1,label:'计时'},
        {value:2,label:'按日'},
        {value:4,label:'按月'},
        {value:8,label:'按年'},
      ]
    }
  },
  created() {
    this.GetTicketList()
    this.GetUserList()
    this.GetRenterList()
  },
  methods: {
    //获取计费方式值
    GetBillingWay(val){
      const item = this.billingWay.find(item => item.value === val);
      return item ? item.label : null;
    },
    //tab切换
    handleTabsClick(){
      switch (this.activeName) {
        case "all":
          this.searchInfo.status = 0
          break;
        case "unAssign":
          this.searchInfo.status = 1
          break;
        case "unHandle":
          this.searchInfo.status = 2
          break;
        case "handled":
          this.searchInfo.status = 3
          break;
      }
      this.GetTicketList()
    },
    //分配对话框
    handleAssign(val){
      this.handleDialog=true;
      this.dataInfo={}
      this.dataInfo.id=val.ID;
    },
    ConfirmAssign(){
      assignSysTicket(this.dataInfo).then(res=>{
        if (res.code == 0){
          this.$message({
            message:'分配成功',
            type:'success',
          })
          this.handleDialog=false;
          this.GetTicketList()
        }
      })
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
        deleteSysTicket(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetTicketList(this.searchInfo)
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
    handleSearch() {
      this.GetTicketList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetTicketList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetTicketList()
    },
    //关闭描述详情对话框
    closeDialog(){
      this.descDialog=false
      this.handleDialog=false
    },
    //获取描述详情
    GetDescDetail(val){
      this.descDialog=true
      this.detailData=JSON.parse(val)
      this.specDetail= JSON.parse(JSON.parse(val).specJson)
    },
    // 获取工单列表
    GetTicketList() {
      getSysTicketList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.ticketList = response.data.list
          this.total = response.data.total
          this.searchInfo.totalNum = response.data.totalNum
          this.searchInfo.unAssignNum = response.data.unAssignNum
          this.searchInfo.unHandleNum = response.data.unHandleNum
          this.searchInfo.handleNum = response.data.handleNum
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //获取分区名称
    GetUserName(id) {
      const tenant = this.userList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
    //获取分区名称
    GetRenterName(id) {
      const tenant = this.renterList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
    // 获取用户列表
    GetUserList() {
      getAllUserList().then(response => {
        if (response.code == 0) {
          this.userList = response.data.list
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    // 获取用户列表
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
    }
  }
}
</script>
