<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tabChange="HandleTabChange">
      <el-tab-pane name="host" label="云主机">
        <div class="gva-search-box" style="margin-top:16px">
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
                <el-button @click="GetStorageList" type="primary">
                  查询
                </el-button>
                <el-button @click="searchInfo.name='';GetStorageList()" type="info">
                  重置
                </el-button>
              </el-form-item>

            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <div style="margin-bottom: 16px;">
            <el-button @click="OpenDrawer" type="primary">
              新增
            </el-button>
            <el-tooltip :disabled="!deleteList.includes(false)" content="该存储已被分配使用"
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
              :data="nasList"
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
            <el-table-column label="存储类型">
              <template #default="scope">
                <span>{{ scope.row.storageType == 'nfs' ? "NFS" : "Ceph" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <i v-if="scope.row.enabled" class="status-dot"/>
                <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                <span>{{ scope.row.enabled ? "启用" : "禁用" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="物理容量">
              <template #default="scope">
                <span>{{ '已使用:' + scope.row.usedInfo.actualCapacityUsed }}</span><br>
                <span>{{ '总量:' + scope.row.usedInfo.capacity }}</span>
              </template>
            </el-table-column>
            <el-table-column label="虚拟容量">
              <template #default="scope">
                <span>{{ '已分配:' + scope.row.usedInfo.virtualUsedCapacity }}</span><br>
                <span>{{ '总量:' + scope.row.usedInfo.virtualCapacity }}</span>
              </template>
            </el-table-column>
            <!--        <el-table-column label="操作">-->
            <!--          <template #default="scope">-->
            <!--            <template #default="scope">-->
            <!--              <el-button v-if="scope.row.enabled" type="text" @click="HandleDisable(scope.row)">禁用</el-button>-->
            <!--              <el-button v-if="!scope.row.enabled" type="text" @click="HandleEnable(scope.row)">启用</el-button>-->
            <!--            </template>-->
            <!--            <el-tooltip :disabled="scope.row.used==0" content="该存储已被分配使用"-->
            <!--                        placement="top">-->
            <!--              <el-button type="text" :disabled="scope.row.used!=0" @click="HandleDelete(scope.row)">删除</el-button>-->
            <!--            </el-tooltip>-->
            <!--          </template>-->
            <!--        </el-table-column>-->
            <el-table-column label="操作">
              <template #default="scope">
                <!-- 操作按钮组 -->
                <el-button type="text" @click="OpenBindDialog(scope.row)">关联宿主机</el-button>
                <el-button type="text" @click="OpenBindHostDialog(scope.row)">宿主机列表</el-button>
                <el-button v-if="scope.row.enabled" type="text" @click="HandleDeleteDialog(scope.row,'disable')">禁用</el-button>
                <el-button v-if="!scope.row.enabled" type="text" @click="HandleEnable(scope.row)">启用</el-button>

                <!-- 删除按钮及提示 -->
                <el-tooltip :disabled="scope.row.status=='unmount'" content="该存储已被分配使用" placement="top">
                  <el-button type="text" :disabled="scope.row.status!='unmount'&&scope.row.status!='offline'"
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
        </div>
      </el-tab-pane>
      <el-tab-pane name="container" label="容器">
        <div class="gva-search-box" style="margin-top:16px">
          <div class="gva-btn-list">
            <el-form :inline="true" :model="conNasSearchInfo" style="display: flex;justify-content: space-between">
              <div>
                <el-form-item label="集群选择">
                  <el-select style="width: 240px;" @change="GetContainerNasList" v-model="conNasSearchInfo.cluster">
                    <el-option v-for="item in clusterList"
                               :label="item.name"
                               :value="item.id"/>
                  </el-select>
                </el-form-item>
                <el-form-item label="名称">
                  <el-input
                      v-model="conNasSearchInfo.name"
                      class="filter-item"
                      placeholder="请输入"
                      style="width: 240px;"
                  />
                </el-form-item>
              </div>
              <el-form-item>
                <el-button @click="GetContainerNasList" type="primary">
                  查询
                </el-button>
                <el-button @click="conNasSearchInfo.name='';GetContainerNasList()" type="info">
                  重置
                </el-button>
              </el-form-item>

            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <div style="margin-bottom: 16px;">
            <el-button @click="OpenConDrawer" type="primary">
              新 建
            </el-button>
              <el-button :disabled="!ids.length" type="grey" @click="HandleDeleteDialog('','deleteConIds')">
                删除
              </el-button>
          </div>
          <el-table
              :cell-style="{'text-align':'left'}"
              :data="conNasList"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column type="selection" width="55"/>
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="存储类型">
              <template #default="scope">
                <span>{{ "Ceph" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <i v-if="scope.row.status=='active'" class="status-dot"/>
                <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                <span>{{ scope.row.status == 'active' ? "启用" : "禁用" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="文件系统">
              <template #default="scope">
                <span>{{ scope.row.fs_type }}</span>
              </template>
            </el-table-column>
            <el-table-column label="存储池">
              <template #default="scope">
                <span>{{ scope.row.pool }}</span>
              </template>
            </el-table-column>
            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ scope.row.age + ' ago' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button link type="text" @click="HandleDeleteDialog(scope.row,'deleteCon')">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="conNasSearchInfo.pageIndex"
                :page-size="conNasSearchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="conNasSearchInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleConCurrentChange"
                @size-change="handleConSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <!--    新增NAS存储-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="50%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增NAS存储</span>
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
          label-position="left"
          label-width="130px"
      >
        <el-form-item label="名称" prop="name">
          <el-input placeholder="请输入名称" v-model="dataInfo.name"/>
        </el-form-item>
        <el-form-item label="介质类型" label-width="120px" style="padding-left: 10px">
          <el-tag type="primary">机械盘</el-tag>
        </el-form-item>
        <el-form-item label="存储类型" label-width="120px" style="padding-left: 10px">
          <el-radio-group v-model="dataInfo.storageType">
            <el-radio-button label="Ceph" value="rbd"/>
            <el-radio-button label="NFS" value="nfs"/>
          </el-radio-group>
        </el-form-item>
        <!--        nfs-->
        <el-form-item v-if="dataInfo.storageType=='nfs'" label="NFS服务地址" prop="host">
          <el-input placeholder="例如：192.168.1.1" v-model="dataInfo.nfs.host"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.storageType=='nfs'" label="NFS共享目录" prop="dir">
          <el-input placeholder="例如：/nfs_root" v-model="dataInfo.dir"/>
        </el-form-item>
        <!--        ceph-->
        <el-form-item v-if="dataInfo.storageType=='rbd'" label="Ceph Mon Host" prop="rbd_mon_host">
          <el-input placeholder="例如：192.168.1.1,192.168.1.2" v-model="dataInfo.ceph.rbd_mon_host"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.storageType=='rbd'" label="存储池" prop="rbd_pool">
          <el-input placeholder="请输入存储池名称" v-model="dataInfo.rbd_pool"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.storageType=='rbd'" label="认证密钥" prop="rbd_key">
          <el-input placeholder="请输入client.admin的认证密钥，例如：AQCP4L1aORs8HRAAdcVS2y2R1oERNa+3xyacTA=="
                    v-model="dataInfo.rbd_key"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    新增容器NAS存储-->
    <el-drawer
        v-model="createConVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="50%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg"
              style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增容器NAS存储</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createConVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>

      <el-form
          ref="conNasInfo"
          :model="conNasInfo"
          :rules="rules"
          label-position="left"
          label-width="90px"
      >
        <el-form-item label="名称" prop="name">
          <el-input placeholder="请输入名称" v-model="conNasInfo.name"/>
        </el-form-item>
        <el-form-item label="集群选择" prop="cluster">
          <el-select v-model="conNasInfo.cluster" @change="GetSecretList">
            <el-option v-for="item in clusterList"
                       :label="item.name"
                       :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="认证密钥" prop="secretId">
          <el-select clearable @change="GetCephPoolList" v-model="conNasInfo.secretId" placeholder="请选择认证密钥">
            <template #empty>
              <div class="option-content">
                <div v-if="isAdding" @click="isAdding=false;" type="text"
                     style="width: 100%;background-color:rgb(89,158,94); height:40px;line-height: 40px;text-align: center;color: #fff">
                  新增
                </div>
                <template v-else>
                  <el-form-item label="UserId">
                    <el-input
                        style="width: 80%;"
                        v-model="conSecretInfo.userId"
                        placeholder="请输入UserId"
                    />
                  </el-form-item>
                  <el-form-item label="UserKey">
                    <el-input
                        style="width: 80%;"
                        v-model="conSecretInfo.userKey"
                        placeholder="请输入UserKey密钥"
                    />
                  </el-form-item>
                  <el-button type="primary" @click="ConfirmValue">
                    确定
                  </el-button>
                  <el-button @click="isAdding=true;conSecretInfo={}">取消</el-button>
                </template>
              </div>
            </template>
            <div class="option-content">
              <div v-if="isAdding" @click="isAdding=false" type="text"
                   style="width: 100%;background-color:rgb(89,158,94); height:40px;line-height: 40px;text-align: center;color: #fff">
                新增
              </div>
              <template v-else>
                <el-form-item label="UserId">
                  <el-input
                      style="width: 80%;"
                      v-model="conSecretInfo.userId"
                      placeholder="请输入UserId"
                  />
                </el-form-item>
                <el-form-item label="UserKey">
                  <el-input
                      style="width: 80%;"
                      v-model="conSecretInfo.userKey"
                      placeholder="请输入UserKey密钥"
                  />
                </el-form-item>
                <el-button type="primary" @click="ConfirmValue">
                  确定
                </el-button>
                <el-button @click="isAdding=true;conSecretInfo={}">取消</el-button>
              </template>
            </div>
            <el-option style="margin-top: 5px" v-for="(item,index) in secretList"
                       :key="index"
                       :label="item.name"
                       :value="item.id">
              <template #default>
                <div class="option-content">
                  <span class="span-style">{{ item.name }}</span>
                  <el-button :disabled="conNasInfo.secretId==item.id" @click="HandleDeleteDialog(item,'deleteSecret')" type="text"
                             class="delete-button">删除
                  </el-button>
                </div>
              </template>
            </el-option>

          </el-select>
        </el-form-item>
        <el-form-item label="Ceph集群" prop="cephClusterId">
          <el-input disabled placeholder="请选择集群" v-model="conNasInfo.cephClusterId"/>
        </el-form-item>
        <el-form-item label="Pool" prop="pool">
          <el-select placeholder="请先填写上述表单后，Pool会自动更新" v-model="conNasInfo.pool">
            <el-option v-for="item in cephPoolList"
                       :label="item"
                       :value="item"/>
          </el-select>
        </el-form-item>
        <el-form-item label="文件系统" prop="csiFsType">
          <el-select v-model="conNasInfo.csiFsType">
            <el-option label="ext4" value="ext4"/>
            <el-option label="xfs" value="xfs"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="createConVisible=false;">取 消</el-button>
        <el-button type="primary" @click="SubmitConNasDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    已绑定宿主机列表-->
    <el-drawer
        v-model="bindHostVisible"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="true"
        title="已绑定宿主机"
        size="40%"
    >
      <el-table :data="bindHostList">
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="HandleDeleteDialog(scope.row,'cancelBind')">解除绑定</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
    <!--    绑定宿主机-->
    <el-dialog
        v-model="bindVisible"
        :title="'绑定宿主机'"
        width="40%"
    >
      <el-form
          ref="bindData"
          :model="bindData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="宿主机" prop="hostIds">
          <el-select multiple placeholder="请选择宿主机" clearable v-model="bindData.hostIds">
            <el-option v-for="(item,index) in hostList"
                       :value="item.id"
                       :label="item.name"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="bindData.type=='nfs'" label="挂载点" prop="dir">
          <el-input placeholder="例如：/nfs，请确保在宿主机上存在该目录" v-model="bindData.dir"/>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="SubmitBindDrawer">确 定</el-button>
        <el-button @click="CloseDrawer">取 消</el-button>
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
  addNasStorage,
  createContainerNas,
  deleteContainerNas, deleteContainerNasByIds,
  deleteNasStorage,
  deleteNasStorageByIds,
  disableNasStorage,
  enableNasStorage,
  getCephClusterId,
  getCephPoolsList,
  getContainerNasList,
  getSecretList,
  getStorageList,
  nasStorageConnHost,
  nasStorageConnHostList,
  nasStorageDisConnHost,
  nasStorageNoConnHostList,
  secretCreate,
  secretDelete
} from '@/api/cloudpods/storage/storage'
import {getBaseK8SList} from "@/api/cloudpods/baseRes/k8s";


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
      nasId: '',
      cancelData: {
        id: '',
        hostId: '',
      },
      bindData: {
        id: '',
        dir: '',
        type: '',
        hostIds: [],
      },
      bindVisible: false,
      bindHostVisible: false,
      hostList: [],
      bindHostList: [],
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
      },
      total: 0,
      nasList: [],
      dataInfo: {
        name: '',
        storageType: '',
        dir: '',
        rbd_pool: '',
        rbd_key: '',
        nfs: {
          host: '',
          dir: '',
        },
        ceph: {
          rbd_key: '',
          rbd_mon_host: '',
          rbd_pool: '',
        }
      },
      createVisible: false,
      rules: {
        name: [
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        host: [
          {required: true, validator: this.hostCheckValid, trigger: 'change'}
        ],
        csiFsType: [
          {required: true, message: '请选择文件系统', trigger: 'change'}
        ],
        rbd_mon_host: [
          {required: true, validator: this.monCheckValid, trigger: 'change'}
        ],
        rbd_pool: [
          {required: true, message: '请输入存储池', trigger: 'blur'}
        ],
        rbd_key: [
          {required: true, message: '请输入client.admin的认证密钥', trigger: 'blur'}
        ],
        hostIds: [
          {required: true, message: '请选择宿主机', trigger: 'blur'}
        ],
        dir: [
          {
            required: true,
            pattern: /^\/[^/]+/, // 必须以 / 开头，且后面至少有一个非 / 字符
            message: '共享目录必须以 / 开头，且不能仅为 /',
            trigger: 'blur',
          },
        ],
        cluster: [
          {required: true, message: '请选择集群', trigger: 'blur'}
        ],
        secretId: [
          {required: true, message: '请选择认证密钥', trigger: 'blur'}
        ],
        cephClusterId: [
          {required: true, message: '请选择Ceph集群', trigger: 'blur'}
        ],
        pool: [
          {required: true, message: '请选择pool', trigger: 'blur'}
        ],
      },

      conNasSearchInfo: {
        pageIndex: 1,
        pageSize: 10,
        total: 0,
        name: '',
        cluster: '',
      },
      clusterList: [],
      conNasList: [],
      conNasInfo: {
        name: '',
        pool: '',
        cephClusterId: '',
        csiFsType: '',
        cluster: '',
        secretId: '',
      },//新增容器Nas
      conSecretInfo: {
        cluster: '',
        userId: '',
        userKey: '',
      },
      secretList: [],//保密字典列表
      createConVisible: false,
      cephPoolList: [],//ceph集群池
      isAdding: false,
    }
  },
  created() {
    this.GetStorageList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDelete(this.deleteId)
          break;
        case 'deleteCon':
          this.HandleDeleteConNas(this.deleteId)
          break;
        case 'deleteConIds':
          this.HandleDeleteConNasByIds()
          break;
        case 'deleteIds':
          this.HandleDeleteByIds()
          break;
        case 'deleteSecret':
          this.HandleDeleteSecret()
          break;
        case 'cancelBind':
          this.CancelBindHost()
          break;
        case 'disable':
          this.HandleDisable(this.deleteId)
          break;
      }
      this.deleteVisible = false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      this.deleteVisible = true;
      switch (type) {
        case 'delete':
        case 'deleteIds':
        case 'deleteCon':
        case 'deleteConIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该NAS存储，是否继续?';
          break;
        case 'deleteSecret':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该认证密钥, 是否继续?';
          break;
        case 'cancelBind':
          this.dialogTitle = '解绑';
          this.dialogInfo = '此操作将取消关联该NAS存储与宿主机, 是否继续?';
          break;
        case 'disable':
          this.dialogTitle = '禁用';
          this.dialogInfo = '此操作将禁用该存储,禁用后,无法使用该存储创建云主机磁盘, 是否继续?';
          break;
      }
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].status != 'unmount' && val[i].status != 'offline') {
          this.deleteList.push(false)
        }
      }
    },
    //校验
    hostCheckValid(rule, value, _callback) {
      const pattern = /^(?=(\b|\D))(((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))(?=(\b|\D))$/
      console.log(this.dataInfo.nfs.host)
      const ips = this.dataInfo.nfs.host.split(',')
      if (!this.dataInfo.nfs.host || this.dataInfo.nfs.host === '') {
        return _callback(new Error('必须输入Host'))
      }
      ips.forEach((item) => {
        const [ip, dir] = item.split(':')
        if (!pattern.test(ip)) {
          return _callback(new Error("请输入合法的Host，例如： 192.168.1.1,192.168.1.2"))
        } else {
          if (dir) {
            const reg = /^\/.+/
            if (!reg.test(dir)) {
              return _callback(new Error("请输入合法的Host，例如： 192.168.1.1,192.168.1.2"))
            }
          }
        }
      })
      _callback()
    },
    //校验
    monCheckValid(rule, value, _callback) {
      const pattern = /^(?=(\b|\D))(((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))(?=(\b|\D))$/
      console.log(this.dataInfo.ceph.rbd_mon_host)
      const ips = this.dataInfo.ceph.rbd_mon_host.split(',')
      if (!this.dataInfo.ceph.rbd_mon_host || this.dataInfo.ceph.rbd_mon_host === '') {
        return _callback(new Error('必须输入Host'))
      }
      ips.forEach((item) => {
        const [ip, dir] = item.split(':')
        if (!pattern.test(ip)) {
          return _callback(new Error("请输入合法的Host，例如： 192.168.1.1,192.168.1.2"))
        } else {
          if (dir) {
            const reg = /^\/.+/
            if (!reg.test(dir)) {
              return _callback(new Error("请输入合法的Host，例如： 192.168.1.1,192.168.1.2"))
            }
          }
        }
      })
      _callback()
    },
    //创建容器Nas存储确定
    SubmitConNasDrawer() {
      this.$refs.conNasInfo.validate(valid => {
        if (valid) {
          createContainerNas(this.conNasInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createConVisible = false;
              this.GetContainerNasList();
            }
          })
        }
      })
    },
    //创建保密字典
    ConfirmValue() {
      this.conSecretInfo.cluster = this.conNasInfo.cluster;
      secretCreate(this.conSecretInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '创建成功',
          })
          this.conSecretInfo = {};
          this.GetSecretList();
        }
        this.conNasInfo.pool = ""
        this.conNasInfo.secretId = ""
        this.isAdding = true;
      })
    },
    //打开容器创建对话框
    OpenConDrawer() {
      this.createConVisible = true;
      this.conNasInfo.cluster = this.conNasSearchInfo.cluster;
      this.conNasInfo.csiFsType = "ext4";
      this.isAdding = true;
      this.GetSecretList();
      this.GetCephClusterId();
      this.GetCephPoolList();
    },
    //获取ceph集群Id
    GetCephClusterId() {
      getCephClusterId({cluster: this.conNasInfo.cluster}).then(res => {
        if (res.code == 0) {
          this.conNasInfo.cephClusterId = res.data.cephClusterId;
        }
      })
    },
    //获取ceph集群列表
    GetCephPoolList() {
      if (this.conNasInfo.secretId != "" && this.conNasInfo.cephClusterId != "" && this.conNasInfo.cluster != '') {
        getCephPoolsList(this.conNasInfo).then(res => {
          if (res.code == 0) {
            this.cephPoolList = res.data.pools;
          }
        })
      }
    },
    //获取保密字典列表
    GetSecretList() {
      if (this.conNasInfo.cluster != '') {
        this.conNasInfo.secretId = ''
        getSecretList({cluster: this.conNasInfo.cluster}).then(res => {
          if (res.code == 0) {
            this.secretList = res.data.list;
          }
        })
      }
    },
    //删除保密字典
    HandleDeleteSecret(val) {
      secretDelete({id: val, cluster: this.conNasSearchInfo.cluster}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.conNasInfo.secretId = '';
          this.GetSecretList()
        }
      })
    },
    //删除容器Nas
    HandleDeleteConNas(val) {
      deleteContainerNas({id: val, cluster: this.conNasSearchInfo.cluster}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetContainerNasList()
        }
      })
    },
    //删除容器Nas
    HandleDeleteConNasByIds() {
      deleteContainerNasByIds({ids: this.ids, cluster: this.conNasSearchInfo.cluster}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetContainerNasList()
        }
      })
    },
    //获取容器Nas存储列表
    GetContainerNasList() {
      getContainerNasList(this.conNasSearchInfo).then(res => {
        if (res.code == 0) {
          this.conNasList = res.data.list;
          this.conNasSearchInfo.total = res.data.total;
        }
      })
    },
    //tab切换
    HandleTabChange() {
      switch (this.activeName) {
        case 'host':
          this.GetStorageList();
          break;
        case 'container':
          getBaseK8SList({
            pageIndex: 1,
            pageSize: 100000,
          }).then(res => {
            if (res.code == 0) {
              this.clusterList = res.data.data;
              if (this.clusterList.length > 0) {
                this.conNasSearchInfo.cluster = this.clusterList[0].id;
                this.GetContainerNasList();
              }
            }
          })
          break;
      }
    },
    //打开已绑定宿主机对话框
    OpenBindHostDialog(val) {
      this.nasId = val.id;
      this.GetBindHostList()
      this.bindHostVisible = true;
    },
    //获取已绑定宿主机
    GetBindHostList() {
      nasStorageConnHostList({id: this.nasId}).then(res => {
        if (res.code == 0) {
          this.bindHostList = JSON.parse(res.data.list);
        }
      })
    },
    //解除绑定
    CancelBindHost(val) {
      this.cancelData.id = this.nasId;
      this.cancelData.hostId = val;
        nasStorageDisConnHost(this.cancelData).then(res => {
          if (res.code == 0) {
            this.$message({
              type: 'success',
              message: '取消绑定成功',
            })
            this.GetBindHostList()
          }
        })
    },
    //打开绑定对话框
    OpenBindDialog(val) {
      this.GetUnBindHost(val);
      this.bindData.id = val.id;
      this.bindData.hostIds = [];
      this.bindData.type = val.storageType;
      this.bindData.dir = '';
      this.bindVisible = true;
    },
    //获取未绑定宿主机
    GetUnBindHost(val) {
      nasStorageNoConnHostList({id: val.id}).then(res => {
        if (res.code == 0) {
          this.hostList = JSON.parse(res.data.list);
        }
      })
    },
    //禁用
    HandleDisable(val) {
        disableNasStorage({id: val}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '禁用成功',
              type: 'success'
            })
            this.GetStorageList(this.searchInfo)
          }
        })
    },
    //启用
    HandleEnable(val) {
        enableNasStorage({id: val.id}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '启用成功',
              type: 'success'
            })
            this.GetStorageList(this.searchInfo)
          }
        })
    },
    //删除NAS存储
    HandleDelete(val) {
      deleteNasStorage({id: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetStorageList(this.searchInfo)
        }
      })
    },
    //批量删除NAS存储
    HandleDeleteByIds(val) {
      deleteNasStorageByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetStorageList(this.searchInfo)
        }
      })
    },
    //新增对话框
    OpenDrawer() {
      this.dataInfo.name = '';
      this.dataInfo.dir = '';
      this.dataInfo.rbd_pool = '';
      this.dataInfo.rbd_key = '';
      this.dataInfo.storageType = 'rbd';
      this.dataInfo.nfs = {};
      this.dataInfo.ceph = {};
      this.createVisible = true;
    },
    //取消
    CloseDrawer() {
      this.createVisible = false;
      this.bindVisible = false;
    },
    //新增确定
    SubmitDrawer() {
      this.$refs.dataInfo.validate(valid => {
        this.dataInfo.ceph.rbd_pool = this.dataInfo.rbd_pool;
        this.dataInfo.ceph.rbd_key = this.dataInfo.rbd_key;
        // this.dataInfo.nfs.host=this.dataInfo.host;
        this.dataInfo.nfs.dir = this.dataInfo.dir;
        if (valid) {
          addNasStorage(this.dataInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createVisible = false;
              this.GetStorageList();
            }
          })
        }
      })
    },
    //绑定确定
    SubmitBindDrawer() {
      this.$refs.bindData.validate(valid => {
        if (valid) {
          nasStorageConnHost(this.bindData).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '绑定成功',
              })
              this.bindVisible = false;
              this.GetStorageList();
            }
          })
        }
      })
    },
    //格式化容量
    FormatCapacitySize(val) {
      var G = 1024
      var T = 1024 * 1024
      if (val < T) {
        return (val / G).toFixed(1) + "G";
      } else if (val >= G && val <= T) {
        return (val / G).toFixed(1) + "G";
      } else if (val > T) {
        return (val / T).toFixed(1) + "T";
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
    //获取云主机Nas存储
    GetStorageList() {
      getStorageList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.nasList = response.data.list
          this.total = response.data.total
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetStorageList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetStorageList()
    },
    handleConCurrentChange(val) {
      this.conNasSearchInfo.pageIndex = val
      this.GetContainerNasList()
    },
    handleConSizeChange(val) {
      this.conNasSearchInfo.pageSize = val
      this.GetContainerNasList()
    },
  }
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

.span-style {
  float: left;
  margin-right: 8px; /* 右边距 */
}

/* 设置删除按钮样式 */
.delete-button {
  float: right;
}

</style>