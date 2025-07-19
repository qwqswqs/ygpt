<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo"
                 style="display: flex;justify-content: space-between">
          <el-form-item label="名称">
            <el-input style="width: 240px;" v-model="searchInfo.name" placeholder="请输入"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="GetEipList">查询</el-button>
            <el-button class="button-gap" type="info" @click="searchInfo.name='';GetEipList()">重置</el-button>
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
        <el-tooltip :disabled="!deleteList.includes(false)" content="必须先解除绑定才可进行此操作"
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
          :data="netEipList"
          @sort-change="handleSortChange"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55"/>

        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP" sortable="custom">
          <template #default="scope">
            <span>{{ scope.row.ip_addr }}</span>
          </template>
        </el-table-column>
        <el-table-column label="带宽">
          <template #default="scope">
            <span>{{ scope.row.bandwidth }}</span>
          </template>
        </el-table-column>
        <el-table-column label="绑定的资源">
          <template #default="scope">
            <span>{{
                scope.row.can_delete ? '-' : scope.row.associate_name + '(' + resType[scope.row.associate_type] + ')'
              }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link :disabled="scope.row.can_delete" @click="OpenDialog(scope.row,false)">修改带宽
            </el-button>
            <el-button type="text" link v-if="scope.row.can_delete" @click="OpenDialog(scope.row,true)">绑定</el-button>
            <el-button type="text" link v-if="!scope.row.can_delete" @click="HandleUnbind(scope.row)">解绑</el-button>
            <el-tooltip :disabled="scope.row.can_delete" content="必须先解除绑定才可进行此操作"
                        placement="top">
              <el-button link :disabled="!scope.row.can_delete" type="text"
                         @click="HandleDeleteDialog(scope.row,'delete')">删除
              </el-button>
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
      <!--      新增对话框-->
      <el-drawer
          v-model="createVisible"
          :close-on-click-modal="true"
          :close-on-press-escape="true"
          :show-close="false"
          size="520"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg"
                style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增弹性公网IP</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="createVisible=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <el-form
            ref="dataInfo"
            :model="dataInfo"
            :rules="rules"
            label-width="130px"
            label-position="left"
        >
          <el-form-item label="公网IP名称" prop="name">
            <el-input style="width: 320px" placeholder="请输入名称" v-model="dataInfo.name"/>
          </el-form-item>
          <el-form-item label="指定IP子网" prop="network" style="margin-top: 30px">
            <el-select style="width: 320px" placeholder="请选择IP子网" v-model="dataInfo.network">
              <el-option v-for="(item,index) in ipList"
                         :disabled="item.ports-item.ports_used==0"
                         :value="item.id"
                         :label="item.name+'('+item.guest_ip_start+' - '+item.guest_ip_end+' 可用:'+(item.ports-item.ports_used)+')'"/>
            </el-select>
          </el-form-item>
          <el-form-item label="带宽峰值(Mbps)" prop="bandwidth">
            <el-input style="width: 320px" type="number" placeholder="请输入带宽峰值"
                      v-model.number="dataInfo.bandwidth"/>
          </el-form-item>
        </el-form>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="CloseDrawer">取 消</el-button>
          <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
        </div>
      </el-drawer>
      <!--    绑定对话框-->
      <el-dialog
          v-model="bindVisible"
          :show-close="false"
          width="40%"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg"
                style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{isBindType?'绑定':'修改带宽'}}</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="bindVisible=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <el-form
            ref="eipData"
            :model="eipData"
            :rules="rules"
            label-width="120px"
        >
          <el-form-item v-if="isBindType" label="云主机" prop="serverId">
            <el-select placeholder="请选择云主机" v-model="eipData.serverId">
              <el-option v-for="(item,index) in bindServerList"
                         :value="item.id"
                         :label="item.name+'('+item.ips+')'"/>
            </el-select>
          </el-form-item>
          <el-form-item v-if="!isBindType" label="带宽峰值" prop="bandwidth">
            <el-input style="width: 150px" type="number" placeholder="请输入带宽峰值" v-model="eipData.bandwidth"/>
            Mbps
          </el-form-item>
        </el-form>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="CloseDrawer">取 消</el-button>
          <el-button type="primary" @click="HandleConfirm">确 定</el-button>
        </div>
      </el-dialog>
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
        <span>此操作将永久删除该弹性公网IP, 是否继续?</span>
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
  </div>
</template>

<script>
import {getNetIpList,} from "@/api/cloudpods/network/ip"
import {
  addNetEip,
  bindEipServer,
  deleteNetEip,
  deleteNetEipByIds,
  getBindServerList,
  getNetEipList,
  unbindEipServer
} from "@/api/cloudpods/network/eip";


export default {
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      ids: [],
      deleteList: [],
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
      },
      total: 0,
      netEipList: [],
      dialogVisible: false,
      IsCreate: true,
      createVisible: false,
      dataInfo: {
        id: '',
        name: '',
        network: '',
        bandwidth: 30,
      },
      rules: {
        name: [
          {required: true, message: '请输入', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        network: [
          {required: true, message: '请输入', trigger: 'blur'}
        ],
        bandwidth: [
          {required: true, message: '请输入', trigger: 'blur'}
        ],
        serverId: [
          {required: true, message: '请输入', trigger: 'blur'}
        ],
      },
      ipList: [],
      resType: {
        ['server']: '云主机',
      },
      isBindType: '',
      bindVisible: false,
      eipData: {
        id: '',
        serverId: '',
        bandwidth: '',
      },
      bindServerList: [],
    }
  },
  created() {
    this.GetEipList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDelete(this.deleteId)
          break;
        case 'deleteIds':
          this.HandleDeleteByIds()
          break;
      }
      this.deleteVisible = false;
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
        if (!val[i].can_delete) {
          this.deleteList.push(false)
        }
      }
    },
    handleSortChange({column, prop, order}) {
      switch (order) {
        case "ascending":
          this.searchInfo.ipDesc = "asc";
          break;
        case "descending":
          this.searchInfo.ipDesc = "desc";
          break;
        default:
          this.searchInfo.ipDesc = "";
          break;
      }
      this.GetEipList();
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
    GetEipList() {
      getNetEipList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.netEipList = response.data.data
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'warning'
          });
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetEipList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetEipList()
    },
    //获取IP列表
    GetIpList() {
      let obj = {}
      obj.pageIndex = 1
      obj.pageSize = 0
      obj.server_type = ["eip"]
      getNetIpList(obj).then(response => {
        if (response.code == 0) {
          this.ipList = response.data.data
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //新增
    OpenDrawer() {
      this.dataInfo = {};
      this.dataInfo.bandwidth = 30;
      this.GetIpList()
      this.createVisible = true
    },
    //取消
    CloseDrawer() {
      this.createVisible = false;
      this.bindVisible = false;
    },
    //确定
    SubmitDrawer() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addNetEip(this.dataInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                message: '创建成功',
                type: 'success'
              })
              this.createVisible = false;
              this.GetEipList()
            }
          })
        }
      })
    },
    //删除
    HandleDelete(val) {
      deleteNetEip({id: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetEipList(this.searchInfo)
        }
      })
    },
    //批量删除
    HandleDeleteByIds(val) {
      deleteNetEipByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetEipList(this.searchInfo)
        }
      })
    },
    //打开对话框
    OpenDialog(val, type) {
      this.eipData = {}
      this.eipData.id = val.id;
      this.isBindType = type;
      if (type) {
        this.eipData.bandwidth = 0;
        this.GetBindServerList(val)
      } else {
        this.eipData.bandwidth = val.bandwidth;
      }
      this.bindVisible = true;
    },
    //获取可绑定云主机
    GetBindServerList(val) {
      getBindServerList({id: val.id}).then(res => {
        if (res.code == 0) {
          this.bindServerList = JSON.parse(res.data).data
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    //绑定确定/修改带宽
    HandleConfirm() {
      this.$refs.eipData.validate(valid => {
        if (valid) {
          bindEipServer(this.eipData).then(res => {
            if (res.code == 0) {
              this.$message({
                message: '绑定成功',
                type: 'success'
              })
              this.bindVisible = false;
              this.GetEipList()
            }
          })
        }
      })
    },
    //解绑
    HandleUnbind(val) {
      this.$confirm('你所选的1个弹性公网IP将执行解绑操作，你是否确认操作？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        unbindEipServer({id: val.id}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '解绑成功',
              type: 'success'
            })
            this.GetEipList(this.searchInfo)
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消解绑'
        });
      });
    },
  }

}
</script>
