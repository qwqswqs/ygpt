<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tabChange="HandleTabChange">
      <el-tab-pane name="host" label="云主机">
        <div class="gva-search-box" style="margin-top:16px">
          <div class="gva-btn-list">
            <el-form :inline="true" :model="vmEvsInfo" style="display: flex;justify-content: space-between">
              <div>
                <el-form-item label="名称">
                  <el-input
                      v-model="vmEvsInfo.name"
                      class="filter-item"
                      placeholder="请输入"
                      style="width: 240px;"
                  />
                </el-form-item>
                <el-form-item label="存储类型">
                  <el-select
                      v-model="vmEvsInfo.diskType"
                      class="filter-item"
                      clearable
                      placeholder="请选择"
                      style="width: 240px;"
                  >
                    <el-option label="系统盘" value="sys"/>
                    <el-option label="数据盘" value="data"/>
                  </el-select>
                </el-form-item>
                <el-form-item label="状态">
                  <el-select
                      v-model="vmEvsInfo.unused"
                      class="filter-item"
                      clearable
                      placeholder="请选择挂载"
                      style="width: 240px;"
                  >
                    <el-option label="已挂载" :value="false"/>
                    <el-option label="未挂载" :value="true"/>
                  </el-select>
                </el-form-item>
              </div>
              <el-form-item>
                <el-button @click="GetVmEvsList" type="primary">
                  查询
                </el-button>
                <el-button @click="vmEvsInfo.name='';vmEvsInfo.unused='';vmEvsInfo.diskType='';GetVmEvsList()" type="primary">
                  重置
                </el-button>
              </el-form-item>

            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <div style="margin-bottom: 16px;">
            <el-button @click="OpenVmEvsDrawer" type="primary">
              新 建
            </el-button>
            <el-tooltip :disabled="!deleteList.includes(false)" content="硬盘正被使用中"
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
              :data="vmEvsList"
              @selection-change="handleSelectionChangeData"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column :show-overflow-tooltip="true" label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘容量">
              <template #default="scope">
                <span>{{ scope.row.storage }}</span>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <i v-if="scope.row.status == 'ready'" class="status-dot"/>
                <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
                <i v-else-if="scope.row.status == 'unknown'" style="background-color: rgb(174,185,192)"
                   class="status-dot"/>
                <el-icon v-else class="status-dot-other">
                  <Loading/>
                </el-icon>
                <span>{{ diskStatusType[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘类型">
              <template #default="scope">
                <span>{{ scope.row.type == 'sys' ? '系统盘' : '数据盘' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="云主机" width="160">
              <template #default="scope">
                <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.Mounted">
                  {{ item.mounted_by }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="是否挂载">
              <template #default="scope">
                <span style="color: #13ce66" v-if="scope.row.Mounted!=null">{{ '已挂载' }}</span>
                <span style="color: orange" v-else>{{ '未挂载' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="本地存储" :show-overflow-tooltip="true">
              <template #default="scope">
                <span>{{ scope.row.from_vm_nas }}</span>
              </template>
            </el-table-column>
            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ FormatDateTime(scope.row.age) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="scope">
                <!-- 操作按钮组 -->
                <el-tooltip :disabled="scope.row.type!='sys'" content="系统盘不支持扩容" placement="top">
                  <el-button :disabled="scope.row.type=='sys'||scope.row.status!='ready'" type="text"
                             @click="OpenResizeDrawer(scope.row)">扩容
                  </el-button>
                </el-tooltip>
                <el-button v-if="scope.row.Mounted==null" type="text" @click="OpenAttachDrawer(scope.row)">挂载
                </el-button>
                <el-tooltip :disabled="scope.row.type!='sys'" content="系统盘不支持卸载" placement="top">
                  <el-button v-if="scope.row.Mounted!=null" :disabled="scope.row.type=='sys'" type="text"
                             @click="HandleDeleteDialog(scope.row,'detach')">卸载
                  </el-button>
                </el-tooltip>
                <el-tooltip :disabled="scope.row.Mounted==null" content="正被使用中" placement="top">
                  <el-button :disabled="scope.row.Mounted!=null" type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除
                  </el-button>
                </el-tooltip>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="vmEvsInfo.pageIndex"
                :page-size="vmEvsInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="vmEvsInfo.total"
                layout="total,  prev, pager, next, sizes,jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane name="container" label="裸金属">
        <div class="gva-search-box" style="margin-top:16px">
          <div class="gva-btn-list">
            <el-form :inline="true" :model="vmEvsInfo" style="display: flex;justify-content: space-between">
              <div>
                <el-form-item label="名称">
                  <el-input
                      v-model="vmEvsInfo.name"
                      class="filter-item"
                      placeholder="请输入名称"
                      style="width: 200px;"
                  />
                </el-form-item>
                <el-form-item label="存储类型">
                  <el-select
                      v-model="vmEvsInfo.diskType"
                      class="filter-item"
                      clearable
                      placeholder="请选择"
                      style="width: 200px;"
                  >
                    <el-option label="系统盘" value="sys"/>
                    <el-option label="数据盘" value="data"/>
                  </el-select>
                </el-form-item>
                <el-form-item label="状态">
                  <el-select
                      v-model="vmEvsInfo.unused"
                      class="filter-item"
                      clearable
                      placeholder="请选择"
                      style="width: 200px;"
                  >
                    <el-option label="已挂载" :value="false"/>
                    <el-option label="未挂载" :value="true"/>
                  </el-select>
                </el-form-item>
              </div>
              <el-form-item>
                <el-button @click="GetVmEvsList" type="primary">
                  搜 索
                </el-button>
              </el-form-item>

            </el-form>
          </div>
        </div>
        <div class="gva-table-box">

          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="vmEvsList"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘容量">
              <template #default="scope">
                <span>{{ scope.row.storage }}</span>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <i v-if="scope.row.status == 'ready'" class="status-dot"/>
                <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
                <i v-else-if="scope.row.status == 'unknown'" style="background-color: rgb(174,185,192)"
                   class="status-dot"/>
                <el-icon v-else class="status-dot-other">
                  <Loading/>
                </el-icon>
                <span>{{ diskStatusType[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘类型">
              <template #default="scope">
                <span>{{ scope.row.type == 'sys' ? '系统盘' : '数据盘' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="裸金属">
              <template #default="scope">
                <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.Mounted">
                  {{ item.mounted_by }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="是否挂载">
              <template #default="scope">
                <span style="color: #13ce66" v-if="scope.row.Mounted!=null">{{ '已挂载' }}</span>
                <span style="color: orange" v-else>{{ '未挂载' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="本地存储">
              <template #default="scope">
                <span>{{ scope.row.from_vm_nas }}</span>
              </template>
            </el-table-column>
            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ FormatDateTime(scope.row.age) }}</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="vmEvsInfo.pageIndex"
                :page-size="vmEvsInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="vmEvsInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <!--    新增云主机硬盘-->
    <el-drawer
        v-model="createVmVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="520"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增硬盘</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVmVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="addVmEvsInfo"
          :model="addVmEvsInfo"
          :rules="rules"
          label-position="left"
          label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input style="width: 320px;" placeholder="请输入名称" v-model="addVmEvsInfo.name"/>
        </el-form-item>
        <el-form-item label="类型" prop="backend" label-width="70px" style="padding-left: 10px;margin-top: 30px">
          <el-select style="width: 320px;" v-model="addVmEvsInfo.backend" @change="addVmEvsInfo.storage_id=''">
            <el-option label="本地硬盘" value="local"/>
          </el-select>
        </el-form-item>
        <el-form-item label="容量" label-width="70px" style="padding-left: 10px">
          <div class="input-with-unit">
            <el-input-number controls-position="right" :precision="0"
                             v-model="addVmEvsInfo.tempSize"
                             type="number" style="width: 320px;"
                             :min="1"
                             :max="2024"
                             :step="10"
            />
            <span class="unit">GB</span>
          </div>
        </el-form-item>
        <el-form-item label="本地存储" prop="storage_id">
          <el-select style="width: 320px;" v-model="addVmEvsInfo.storage_id">
            <el-option v-if="addVmEvsInfo.backend=='local'"
                       v-for="item in localStoList"
                       :label="item.name"
                       :value="item.id"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="createVmVisible=false">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>

    <!--    扩容-->
    <el-dialog
        v-model="resizeVisible"
        :show-close="false"
        width="400"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">硬盘扩容</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="resizeVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          :model="resizeInfo"
          label-width="40px"
      >
        <el-form-item label="容量">
          <div class="input-with-unit">
            <el-input-number controls-position="right" :precision="0"
                             v-model="resizeInfo.tempSize"
                             type="number"
                             style="width: 240px"
                             :min="resizeInfo.tempSize"
                             :max="2024"
                             :step="10"
            />
            <span class="unit">GB</span>
          </div>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="resizeVisible=false;">取 消</el-button>
        <el-button type="primary" @click="SubmitResizeDrawer">确 定</el-button>
      </div>
    </el-dialog>
    <!--    挂载-->
    <el-dialog
        v-model="attachVisible"
        :show-close="false"
        width="40%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">挂载云主机</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="attachVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="attachInfo"
          :rules="rules"
          :model="attachInfo"
          label-width="120px"
      >
        <el-form-item label="可挂载云主机" prop="vmId">
          <el-select v-model="attachInfo.vmId" placeholder="请选择云主机">
            <el-option v-for="item in attachHostList"
                       :label="item.name"
                       :value="item.id"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="attachVisible=false;">取 消</el-button>
        <el-button type="primary" @click="SubmitAttachDrawer">确 定</el-button>
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
  attachVmEvs,
  createContainerEvs,
  createVmEvs,
  deleteContainerEvs,
  deleteVmEvs, deleteVmEvsByIds,
  detachVmEvs,
  getAttachableVms,
  getContainerEvsList,
  getVmEvsList,
  resizeVmEvs
} from "@/api/cloudpods/storage/evs";
import {getContainerNasList, getStorageList} from "@/api/cloudpods/storage/storage";
import statusType from "@/locales/zh-CN.json";


export default {
  data() {
    return {
      dialogTitle: '',
      dialogInfo: '',
      deleteId: '',
      deleteVisible: false,
      ids: [],
      deleteList: [],
      activeName: 'host',
      resizeVisible: false,
      resizeInfo: {
        id: '',
        size: 0,
        tempSize: 0,
      },
      vmEvsInfo: {
        pageIndex: 1,
        pageSize: 10,
        storageType: 'local',
        total: 0,
        diskType: '',
        name: '',
        unused: null,
      },
      conEvsInfo: {
        pageIndex: 1,
        pageSize: 10,
        total: 0,
        name: '',
        unused: null,
        cluster: '',
      },
      vmEvsList: [],
      conEvsList: [],
      clusterList: [],
      createVmVisible: false,
      createConVisible: false,
      addVmEvsInfo: {
        backend: '',
        mediumType: '',
        name: '',
        storage_id: '',
        size: 1,
        tempSize: 1,
      },
      addConEvsInfo: {
        cluster: '',
        name: '',
        size: '',
        containerNas: '',
      },
      localStoList: [],
      nfsStoList: [],
      cephStoList: [],
      conNasList: [],
      isShowBlock: false,
      rules: {
        name: [
          {required: true, trigger: 'blur', message: '请输入名称'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        containerNas: [
          {required: true, trigger: 'blur', message: '请选择容器NAS存储'},
        ],
        vmId: [
          {required: true, trigger: 'blur', message: '请选择可挂载云主机'},
        ],
        storage_id: [
          {required: true, trigger: 'blur', message: '请选择本地存储'},
        ],
        cluster: [
          {required: true, trigger: 'blur', message: '请选择集群'},
        ]
      },
      diskStatusType: statusType.status.disk,
      attachHostList: [],
      attachInfo: {
        id: '',
        vmId: '',
      },
      attachVisible: false,
    }
  },
  created() {
    this.GetVmEvsList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDelete(this.deleteId)
          break;
        case 'detach':
          this.HandleDetach(this.deleteId)
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
      switch (type) {
        case 'delete':
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该硬盘，是否继续?';
          break;
        case 'detach':
          this.dialogTitle = '卸载';
          this.dialogInfo = '此操作将取消挂载云主机, 是否继续?';
          break;
      }
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].Mounted!=null){
          this.deleteList.push(false)
        }
      }
    },
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        this.timer = setInterval(() => {
          this.GetVmEvsList();
        }, 5000);
      } else {
      }
    },
    //关闭定时器
    stopInterval() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
      }
    },
    //打开挂载对话框
    OpenAttachDrawer(val) {
      this.GetAttachHostList(val);
      this.attachInfo.id = val.id;
      this.attachVisible = true;
    },
    //挂载确定
    SubmitAttachDrawer() {
      this.$refs.attachInfo.validate(valid => {
        if (valid) {
          attachVmEvs(this.attachInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '挂载成功',
              })
              this.GetVmEvsList();
              this.attachVisible = false;
            }
          })
        }
      })
    },
    //卸载
    HandleDetach(val) {
        detachVmEvs({id: val.id, vm_id: val.Mounted[0].mounted_by_id}).then(res => {
          if (res.code == 0) {
            this.$message({
              type: 'success',
              message: '卸载成功',
            })
            this.GetVmEvsList();
          }
        })
    },
    //获取可挂载云主机列表
    GetAttachHostList(val) {
      getAttachableVms({id: val.id}).then(res => {
        if (res.code == 0) {
          this.attachHostList = JSON.parse(res.data.list);
        }
      })
    },
    //打开扩容对话框
    OpenResizeDrawer(val) {
      this.resizeInfo.id = val.id;
      const match = (val.storage).match(/\d+/);
      this.resizeInfo.tempSize = parseInt(match[0], 10)
      this.resizeVisible = true;
    },
    //扩容确定
    SubmitResizeDrawer() {
      this.resizeInfo.size = this.resizeInfo.tempSize * 1024
      resizeVmEvs(this.resizeInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '扩容成功',
          })
          this.GetVmEvsList();
          this.resizeVisible = false;
        }
      })
    },
    //删除云主机硬盘
    HandleDelete(val) {
        deleteVmEvs({id: val.id}).then(res => {
          if (res.code == 0) {
            this.$message({
              type: 'success',
              message: '删除成功',
            })
            this.GetVmEvsList();
          }
        })
    },
    //删除云主机硬盘
    HandleDeleteByIds(val) {
        deleteVmEvsByIds({ids:this.ids}).then(res => {
          if (res.code == 0) {
            this.$message({
              type: 'success',
              message: '删除成功',
            })
            this.GetVmEvsList();
          }
        })
    },
    //云主机硬盘创建确定
    SubmitDrawer() {
      this.$refs.addVmEvsInfo.validate(valid => {
        if (valid) {
          this.addVmEvsInfo.size = this.addVmEvsInfo.tempSize * 1024;
          createVmEvs(this.addVmEvsInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.GetVmEvsList();
              this.createVmVisible = false;
            }
          })
        }
      })
    },
    //获取本地存储和NAS存储
    GetStorageList() {
      let obj = {};
      obj.pageSize = 0;
      obj.storageType = 'local';
      obj.enable = 'true';
      getStorageList(obj).then(res => {
        if (res.code == 0) {
          this.localStoList = res.data.list;
        }
      })
      obj.storageType = 'nfs';
      getStorageList(obj).then(res => {
        if (res.code == 0) {
          this.nfsStoList = res.data.list;
        }
      })
      obj.storageType = 'rbd';
      getStorageList(obj).then(res => {
        if (res.code == 0) {
          this.cephStoList = res.data.list;
        }
      })
    },
    //打开云主机硬盘创建对话框
    OpenVmEvsDrawer() {
      this.GetStorageList();
      this.addVmEvsInfo = {};
      this.addVmEvsInfo.backend = 'local';
      this.addVmEvsInfo.size = 10;
      this.addVmEvsInfo.tempSize = 10;
      this.addVmEvsInfo.mediumType = 'rotate';
      this.createVmVisible = true;
    },
    //tab切换
    HandleTabChange() {
      switch (this.activeName) {
        case 'host':
          this.vmEvsInfo.storageType = "local"
          this.GetVmEvsList();
          break;
        case 'container':
          this.vmEvsInfo.storageType = "baremetal"
          this.GetVmEvsList();
          break;
      }
    },
    //获取云主机硬盘列表
    GetVmEvsList() {
      getVmEvsList(this.vmEvsInfo).then(res => {
        if (res.code == 0) {
          const tempData = res.data.list;
          if (tempData !== this.vmEvsList) {
            this.vmEvsList = res.data.list;
            this.vmEvsInfo.total = res.data.total;
            this.startInterval()
          }else{
            this.stopInterval()
          }
        }
      })
    },
    handleCurrentChange(val) {
      this.vmEvsInfo.pageIndex = val
      this.GetVmEvsList()
    },
    handleSizeChange(val) {
      this.vmEvsInfo.pageSize = val
      this.GetVmEvsList()
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

.status-dot-other {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-right: 8px;
  animation: rotate 2s linear infinite;
}

.input-with-unit {
  position: relative;
  display: inline-block;
}

.unit {
  position: absolute;
  right: 40px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none; /* 确保用户不能与单位文本交互 */
  color: #606266; /* 与输入框文本颜色一致 */
}
</style>