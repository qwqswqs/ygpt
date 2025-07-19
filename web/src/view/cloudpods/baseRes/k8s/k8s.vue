<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
          <el-form-item label="名称">
            <el-input
                v-model="searchInfo.name"
                class="filter-item"
                placeholder="请输入"
                style="width: 240px;"
            />
          </el-form-item>
          <el-form-item>
            <el-button @click="GetBaseK8SList" type="primary">
              查询
            </el-button>
            <el-button class="button-gap" type="info" @click="searchInfo.name='';GetBaseK8SList()">重置</el-button>
          </el-form-item>

        </el-form>
      </div>
    </div>
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button @click="OpenDrawer" type="primary">
          新增
        </el-button>
        <el-button :disabled="!ids.length" type="grey" @click="HandleDeleteDialog('','deleteIds')">
          删除
        </el-button>
      </div>
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="clusterList"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55"/>
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="OpenNodeDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="发行版">
          <template #default="scope">
            <span>{{ scope.row.version }}</span>
          </template>
        </el-table-column>
        <el-table-column label="节点数量">
          <template #default="scope">
            <span>{{ scope.row.machines }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <i v-if="scope.row.status=='running'" class="status-dot"/>
            <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
            <span>{{ statusType[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button @click="OpenUpdateDrawer(scope.row)" type="text">
              编辑
            </el-button>
            <el-button type="text" link :disabled="scope.row.status=='deleting'" @click="HandleDeleteDialog(scope.row,'delete')">删除
            </el-button>
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

      <el-drawer
          v-model="createVisible"
          :close-on-click-modal="true"
          :close-on-press-escape="true"
          :show-close="false"
          size="40%"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">创建K8S集群</span>
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
          <el-form-item label="名称" prop="name">
            <el-input placeholder="请输入集群名称" v-model="dataInfo.name"/>
          </el-form-item>
          <el-form-item label="KubeConfig 配置" prop="kubeConfig" style="margin-top: 30px">
            <el-input type="textarea" placeholder="请输入KubeConfig 配置" v-model="dataInfo.kubeConfig"/>
          </el-form-item>
        </el-form>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="CloseDrawer">取 消</el-button>
          <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
        </div>
      </el-drawer>
    </div>
    <!--    节点列表对话框-->
    <el-drawer v-model="nodeVisible"
               :close-on-click-modal="true"
               :close-on-press-escape="true" size="70%":show-close="false">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">详情</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="nodeVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <div style="margin: 0 10px">
        <el-tabs v-model="activeName">
          <el-tab-pane name="node" label="节点列表">
            <el-table :data="nodeList"
                      :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
                      style="width: 100%;font-size: 15px;">
              <el-table-column label="名称">
                <template #default="scope">
                  <span>{{ scope.row.name }}</span>
                </template>
              </el-table-column>
              <el-table-column label="IP">
                <template #default="scope">
                  <span>{{ scope.row.address[0].address }}</span>
                </template>
              </el-table-column>
              <el-table-column label="状态">
                <template #default="scope">
                  <span>{{ scope.row.status }}</span>
                </template>
              </el-table-column>
              <el-table-column label="CPU(核)">
                <template #default="scope">
                  <span>{{
                      scope.row.allocatedResources != undefined ? (scope.row.allocatedResources.cpuRequests / 1000) + '/' + (scope.row.allocatedResources.cpuCapacity / 1000) : ''
                    }}</span>
                </template>
              </el-table-column>
              <el-table-column label="内存">
                <template #default="scope">
                  <span>{{
                      scope.row.allocatedResources != undefined ? Number(scope.row.allocatedResources.memoryRequests / 1000 / 1024 / 1024).toFixed(0) + 'G/' + Number(scope.row.allocatedResources.memoryCapacity / 1000 / 1024 / 1024).toFixed(0) + 'G' : ''
                    }}</span>
                </template>
              </el-table-column>
              <el-table-column label="可调度">
                <template #default="scope">
                  <span>{{ scope.row.unschedulable ? '不可调度' : '可调度' }}</span>
                </template>
              </el-table-column>
              <el-table-column label="操作">
                <template #default="scope">
                  <el-button type="text" v-if="scope.row.unschedulable" @click="ModNodeStatus(scope.row,true)">
                    设置为可调度
                  </el-button>
                  <el-button type="text" v-if="!scope.row.unschedulable" @click="ModNodeStatus(scope.row,false)">
                    设置为不可调度
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="gva-pagination">
              <el-pagination
                  :current-page="nodeSearchInfo.pageIndex"
                  :page-size="nodeSearchInfo.pageSize"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="nodeSearchInfo.total"
                  layout="total,prev, pager, next,  sizes, jumper"
                  @current-change="handleCurrentChange"
                  @size-change="handleSizeChange"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane name="log" label="操作日志">
            <Logs :log-data="logData"/>
          </el-tab-pane>
        </el-tabs>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="nodeVisible=false">关闭</el-button>
        </div>
      </div>
    </el-drawer>
    <!--    节点列表对话框-->
    <el-drawer v-model="updateVisible"
               :close-on-click-modal="true"
               :close-on-press-escape="true" size="520" :show-close="false">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">编辑集群</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="updateVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form :model="updateInfo" label-width="60">
        <el-form-item label="名称">
          <el-input v-model="updateInfo.name" disabled/>
        </el-form-item>
        <el-form-item label="服务ip">
          <el-input v-model="updateInfo.ip" disabled/>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model="updateInfo.endpoint"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="updateVisible=false">取 消</el-button>
        <el-button type="primary" @click="handleUpdateConfirm">确 定</el-button>
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
      <span>此操作将永久删除该集群, 是否继续?</span>
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
  deleteBaseCluster,
  deleteBaseClusterByIds,
  getBaseK8SList,
  importBaseCluster,
  updateBaseCluster
} from '@/api/cloudpods/baseRes/k8s'
import {getK8sNodeList, modK8sNodeStatus} from "@/api/cloudpods/k8s/deployment";
import Logs from "@/view/cloudpods/component/log.vue";


export default {
  components: {Logs},
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      updateInfo:{
        id:"",
        ip:"",
        name:"",
        endpoint:"",
      },
      updateVisible:false,
      ids: [],
      activeName: 'node',
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      nodeSearchInfo: {
        pageIndex: 1,
        pageSize: 10,
        clusterID: '',
        total: 0,
      },
      nodeVisible: false,
      nodeList: [],
      statusType: {
        ['running']: '运行中',
        ['lost']: 'lost',
        ['deleting']: '正在删除中',
      },
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
      },
      total: 0,
      clusterList: [],
      dataInfo: {
        name: '',
        kubeConfig: '',
      },
      rules: {
        name: [
          {required: true, message: '请输入仓库名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        kubeConfig: [
          {required: true, message: '请输入配置', trigger: 'blur'}
        ],
      },
      createVisible: false,
      modNodeStatus: {
        cluster: '',
        id: '',
        status: false,
      }
    }
  },
  created() {
    this.GetBaseK8SList()
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
      this.deleteVisible=false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val;
      this.deleteType = type;
      this.deleteVisible = true;
    },
    //编辑确定
    handleUpdateConfirm(){
      updateBaseCluster(this.updateInfo).then(res=>{
        if (res.code==0){
          this.updateVisible=false;
          this.$message({
            type:"success",
            message:"编辑成功",
          })
          this.GetBaseK8SList();
        }
      })
    },
    //打开编辑对话框
    OpenUpdateDrawer(val){
      this.updateVisible=true;
      this.updateInfo.id=val.id;
      this.updateInfo.name=val.name;
      this.updateInfo.ip=val.serverIp;
      this.updateInfo.endpoint=val.port;
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
      }
    },
    //修改调度状态
    ModNodeStatus(val, status) {
      this.modNodeStatus.cluster = val.cluster;
      this.modNodeStatus.id = val.id;
      this.modNodeStatus.status = status;
      modK8sNodeStatus(this.modNodeStatus).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '设置成功',
          })
          this.GetK8sNodeList()
        }
      })
    },
    //获取节点列表
    GetK8sNodeList() {
      getK8sNodeList(this.nodeSearchInfo).then(res => {
        if (res.code == 0) {
          this.nodeList = res.data.data;
          this.nodeSearchInfo.total = res.data.total;
        }
      })
    },
    //打开节点对话框
    OpenNodeDrawer(val) {
      this.nodeSearchInfo.clusterID = val.id;
      this.logData.id = val.id;
      this.GetK8sNodeList();
      this.nodeVisible = true;
    },
    //删除集群
    HandleDelete(val) {
        deleteBaseCluster({id: val.id}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetBaseK8SList(this.searchInfo)
          }
        })
    },
    //批量删除集群
    HandleDeleteByIds() {
        deleteBaseClusterByIds({ids: this.ids}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetBaseK8SList(this.searchInfo)
          }
        })
    },
    //新增对话框
    OpenDrawer() {
      this.dataInfo = {};
      this.createVisible = true;
    },
    //取消
    CloseDrawer() {
      this.createVisible = false;
    },
    //确定
    SubmitDrawer() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          importBaseCluster(this.dataInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createVisible = false;
              this.GetBaseK8SList();
            }
          })
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
    GetBaseK8SList() {
      getBaseK8SList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.clusterList = response.data.data
          this.total = response.data.total
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetBaseK8SList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetBaseK8SList()
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
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}
</style>