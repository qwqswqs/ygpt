<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
            label-position="left"
            style="display: flex;justify-content: space-between"
        >
          <el-form-item>
            <el-form-item label="资源类别">
              <el-select
                  v-model="searchInfo.source"
                  placeholder="请选择"
              >
                <el-option v-for="(item,index) in specType"
                           :label="item"
                           :key="index"
                           :value="index+1"/>
              </el-select>
            </el-form-item>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="GetRenterList"
            >
              查询
            </el-button>
            <el-button
                type="info"
                @click="searchInfo.source=null"
            >重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

    </div>
    <div >
      <el-form
          label-position="left" ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item>
          <el-button  type="primary" @click="HandleExportData">导出数据</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
          :data="tableData"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
<!--        <el-table-column-->
<!--            type="selection"-->
<!--            width="55">-->
<!--        </el-table-column>-->
<!--        <el-table-column-->
<!--           align="left"-->
<!--            type="index"-->
<!--            label="序号"-->
<!--            min-width="60"-->
<!--        />-->
<!--        <el-table-column-->
<!--           align="left"-->
<!--            label="租户编号"-->
<!--            min-width="150"-->
<!--            prop="tenantId"-->
<!--        />-->
        <el-table-column
            align="left"
            label="租户名"
            min-width="150"
        >
          <template #default="scope">
            {{ scope.row.username || '-' }}
          </template>
        </el-table-column>>
        <el-table-column
           align="left"
            label="资源类别"
            min-width="180"
            prop="resourceType"
        >
          <template #default="scope">
            <span>{{resourceType[scope.row.resourceType-1]}}</span>
          </template>
        </el-table-column>
        <el-table-column
           align="left"
            label="关联订单"
            min-width="180"
            prop="orderId"
        />
<!--        <el-table-column-->
<!--           align="left"-->
<!--            label="状态"-->
<!--            min-width="150">-->
<!--          <template #default="scope">-->
<!--            <span>{{statusType[scope.row.status-1]}}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column
           align="left"
            label="生成日期"
            min-width="180">
          <template #default="scope">
            {{FormatDateTime(scope.row.created_at )}}
          </template>
        </el-table-column>
<!--        <el-table-column-->
<!--            fixed="right"-->
<!--            label="操作"-->
<!--            min-width="250"-->
<!--        >-->
<!--          <template #default="scope">-->
<!--            <el-button type="text" @click="OpenHandleDialog(scope.row)" v-if="scope.row.status==1">处理</el-button>-->
<!--            <el-button type="text" @click="OpenEditDialog(scope.row)" v-if="scope.row.status==2">修改</el-button>-->
<!--            <el-button type="text" v-if="scope.row.status==3">回收</el-button>-->
<!--          </template>-->
<!--        </el-table-column>-->

      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="searchInfo.total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-dialog v-model="handleDialog" :before-close="CloseDialog" :title="dialogTitle" width="400px">

      <el-form
          :model="handleDataInfo"
          :rules="rules"
          label-width="80px"
          label-position="left"
      >
        <el-form-item
            label="资源编号"
            prop="resourceID"
        >
          <el-input v-model="handleDataInfo.resourceID"/>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" style="margin: 0 auto;text-align:right;">
        <el-button @click="CloseDialog">取 消</el-button>
        <el-button type="primary" @click="HandleConfirmDialog">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog
        :title="'导出数据'"
        v-model="dialogVisible"
        width="30%"
        @close="CloseDialog"
    >
      <p>请选择导出数据范围</p>
      <div class="dialog-footer"  style="margin: 4px auto;text-align:right;">
      <el-button @click="HandleExport(true)">当前页</el-button>
      <el-button @click="HandleExport(false)">所有页</el-button>
    </div>
    </el-dialog>
  </div>
</template>
<script>
import { getResRenterList,bindResInfo} from "@/api/yunguan/run/renter";
import * as XLSX from "xlsx";
import {deleteImage} from "@/api/cloudpods/image/image";
import {Stamp} from "@element-plus/icons-vue";

export default {
  components: {Stamp},
  data() {
    return {
      exportData:[],
      searchInfo: {
        page: 1,
        pageSize: 10,
        total:0,
      },
      dialogTitle:'',
      tableData:[],
      resourceType:["裸金属","云主机","容器","存储"],
      statusType:["待处理","执行中","执行完","已回收"],//1.待处理 2.执行中 3.执行完 4.已回收
      handleDialog:false,
      dialogVisible:false,
      handleDataInfo:{
        resourceID:0,
        ID:0,
      },
      specType: ["裸金属", "云主机", "容器"],
      rules:{
        resourceID: [
          {required: true, message: '请输入资源编号', trigger: 'blur'}
        ],
      }
    }
  },
  created() {
    this.GetRenterList()
  },
  methods: {

    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        // 如果日期是不合理的，返回空字符串
        return '-';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    //导出
    HandleExport(type){
      if(type) {
        {
          // 创建一个工作簿
          const workbook = XLSX.utils.book_new();
          this.exportData = [],
              //生成数据 todo:待写入表头
              this.exportData.push(["租户账号", "资源类型", "关联订单", "生成日期"])
          for (var i of this.tableData) {
            this.exportData.push([
              //todo:待写入数据
              i.username,
              this.resourceType[i.resourceType - 1],
              i.orderId,
              this.FormatDateTime(i.created_at)
            ])
          }
          // 创建一个工作表
          const worksheet = XLSX.utils.aoa_to_sheet(this.exportData);

          // 将工作表添加到工作簿
          XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1');

          // 生成Excel文件
          XLSX.writeFile(workbook, 'data.xlsx');
        }
      }
      else {
          let obj={};
          obj.page=1;
          obj.pageSize=10000;
          let table;
          getResRenterList(obj).then(response => {
            if (response.code == 0) {
              table = response.data.list
              {
                // 创建一个工作簿
                const workbook = XLSX.utils.book_new();
                this.exportData = [],
                    //生成数据 todo:待写入表头
                    this.exportData.push(["租户账号", "资源类型", "关联订单", "生成日期"])
                for (var i of table) {
                  this.exportData.push([
                    //todo:待写入数据
                    i.username,
                    this.resourceType[i.resourceType - 1],
                    i.orderId,
                    this.FormatDateTime(i.created_at)
                  ])
                }
                // 创建一个工作表
                const worksheet = XLSX.utils.aoa_to_sheet(this.exportData);

                // 将工作表添加到工作簿
                XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1');

                // 生成Excel文件
                XLSX.writeFile(workbook, 'data.xlsx');
              }
            }
          })
      }
    },
    //导出数据
    HandleExportData(){
      this.dialogVisible=true;
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetRenterList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetRenterList()
    },
    //获取列表
    GetRenterList(){
      getResRenterList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.tableData=response.data.list
          this.searchInfo.total=response.data.total
        }
      })
    },
    //处理对话框
    OpenHandleDialog(val){
      this.handleDataInfo.ID=val.id
      this.handleDialog=true
      this.dialogTitle='租户分配资源'
    },
    OpenEditDialog(val){
      this.dialogTitle='租户修改资源'
      this.handleDataInfo.ID=val.id
      this.handleDataInfo.resourceID=val.resourceID
      this.handleDialog=true
    },
    //关闭对话框
    CloseDialog(){
      this.handleDialog=false
    },
    //处理对话框确定
    HandleConfirmDialog(){
      this.handleDataInfo.resourceID=parseInt(this.handleDataInfo.resourceID)
      bindResInfo(this.handleDataInfo).then(response=>{
        if (response.code == 0){
          this.$message({
            message: '分配成功',
            type: 'success'
          })
          this.handleDialog=false
          this.GetRenterList()
        }
      })
    },
  }
}
</script>

<style>

</style>