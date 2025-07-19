<template>
  <!--    新增和编辑-->
  <el-drawer
      v-model="createVisible"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      size="65%"
  >
    <div class="flex justify-between items-center" style="margin-bottom: 24px">
      <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{
          IsCreated ? '新增' : '编辑'
        }}云主机</span>
      <el-icon
          class="cursor-pointer hover:text-gray-500 transition-colors"
          @click="CloseDrawer"
          style="width: 16px; height: 16px; font-size: 16px"
      >
        <Close/>
      </el-icon>
    </div>
    <el-form
        ref="serverDataInfo"
        :model="serverDataInfo"
        :rules="IsCreated?rules:''"
        label-width="120px"
        label-position="left"
    >
      <el-form-item label="名称" prop="generate_name">
        <el-input class="formStyle" :disabled="true" v-model="serverDataInfo.generate_name"/>
      </el-form-item>
      <el-form-item label="CPU(核)" prop="generate_name">
        <el-input-number style="width: 200px;" controls-position="right" :precision="0" :min="1"
                         v-model="skuData.cpu"/>
      </el-form-item>
      <el-form-item label="内存(GB)" prop="generate_name">
        <el-input-number style="width: 200px;" controls-position="right" :precision="0" :min="1"
                         v-model="skuData.mem"/>
      </el-form-item>

      <!--      <el-form-item label="到期时间">-->
      <!--        <div style="display: flex">-->
      <!--          <el-input-number disabled style="flex: 2;width: 100%" controls-position="right" :precision="0" :min="1" v-model="serverDataInfo.timeNum"/>-->
      <!--          <el-select disabled style="flex: 1;width: 100%" v-model="serverDataInfo.timeType">-->
      <!--            <el-option label="小时" :value="'h'"/>-->
      <!--            <el-option label="天" :value="'d'"/>-->
      <!--            <el-option label="月" :value="'m'"/>-->
      <!--            <el-option label="年" :value="'y'"/>-->
      <!--          </el-select>-->
      <!--        </div>-->
      <!--      </el-form-item>-->
      <el-form-item label="GPU" style="padding-left: 10px" label-width="100">
        <el-tooltip
            :disabled="!deviceInfoArray.length==0 || openDevice"
            content="当前节点GPU资源已用完，目前无法配置GPU实例"
            placement="top"
        >
          <el-switch
              :disabled="deviceInfoArray.length==0"
              v-model="openDevice"
              :active-value="true"
              :inactive-value="false"
              active-text="启用"
              class="ml-2"
              inactive-text="禁用"
              inline-prompt
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-tooltip>
      </el-form-item>

      <div v-for="(item,index) in selectDevice" style="display: flex;" :key="index">
        <el-form-item v-if="openDevice" label="GPU类型" style="flex: 1">
          <el-select placeholder="请选择GPU类型" style="width: 200px;" v-model="item.index">
            <el-option v-for="(item1,index1) in deviceInfoArray"
                       :label="'类型: '+item1.model"
                       :value="index1"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="openDevice" label="数量" style="flex: 1;">
          <el-radio-group v-model="item.num">
            <div v-for="n in deviceInfoArray[item.index].count" :key="n">
              <el-radio-button :label="n" :value="n"/>
            </div>
          </el-radio-group>
        </el-form-item>
      </div>
      <el-form-item prop="renterID" v-if="isBindRenter" label="租户">
        <el-select disabled placeholder="请选择租户" @change="generateUniqueName(serverDataInfo.renterID);"
                   v-model="serverDataInfo.renterID">
          <el-option v-for="item in renterList"
                     :label="item.username"
                     :value="item.renterID"/>
        </el-select>
      </el-form-item>
      <el-form-item label="镜像" prop="imageId">
        <el-select
            v-model="serverDataInfo.imageId"
            placeholder="请选择镜像"
            style="width: 100%"
            @change="skuChange">
          <el-option
              v-for="(optionItem, optionIndex) in imageList"
              :key="optionIndex"
              :label="'系统:'+(optionItem.properties!=undefined?optionItem.properties.distro:'未知')+' 名称:'+optionItem.name"
              :value="optionItem.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="VPC" prop="vpcId">
        <el-select
            v-model="serverDataInfo.vpcId"
            placeholder="请选择vpc"
            style="width: 100%"
            @change="GetIpList();serverDataInfo.network='';">
          <el-option
              v-for="(optionItem, optionIndex) in vpcList"
              :key="optionIndex"
              :label="'名称:'+optionItem.name"
              :value="optionItem.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="IP子网" prop="network">
        <el-select
            v-model="serverDataInfo.network"
            placeholder="请选择IP子网"
            style="width: 100%">
          <el-option
              v-for="(optionItem, optionIndex) in IpList"
              :key="optionIndex"
              :disabled="(optionItem.ports-optionItem.ports_used)==0"
              :label="optionItem.name+'('+optionItem.guest_ip_start+'-'+optionItem.guest_ip_end+')    '+'可用：'+(optionItem.ports-optionItem.ports_used)"
              :value="optionItem.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="绑定弹性公网IP" v-if="serverDataInfo.vpcId!=''&&serverDataInfo.vpcId!='default'">
        <el-switch
            v-model="isEip"
            :active-value="true"
            :inactive-value="false"
            active-text="是"
            class="ml-2"
            inactive-text="否"
            inline-prompt
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            @change="handleBindRenterChange"
        />
      </el-form-item>
      <div v-if="isEip">
        <el-form-item label="网络计费方式">
          <el-tag>按带宽计费</el-tag>
        </el-form-item>
        <el-form-item label="带宽(Mbps)" prop="bandwidth">
          <el-input class="formStyle"  type="number" placeholder="请输入带宽峰值" v-model.number="serverDataInfo.eip_bw"/>
        </el-form-item>
      </div>
      <el-form-item label="系统盘(GB)" prop="disks" style="display: flex; align-items: center; gap: 0px; padding-left: 10px" label-width="110">
        <el-table :data="serverDataInfo.disks"
                  :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
          <el-table-column label="硬盘类别">
            <template #default="scope">
              <el-select v-model="scope.row.backend"
                         @change="HandleDiskChange('local')">
                <el-option :value="'local'" label="本地硬盘"/>
                <el-option :value="'nfs'" v-if="nfsStoList.length>0" label="云硬盘(NFS)"/>
                <el-option :value="'rbd'" v-if="cephStoList.length>0" label="云硬盘(Ceph)"/>
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="容量(GB)">
            <template #default="scope">
              <el-tooltip
                  v-if="scope.row.size === minDiskSize"
                  :content="'该镜像系统盘最小容量为'+ minDiskSize +'GB'"
                  placement="top">
                <el-input-number

                    controls-position="right"
                    :precision="0"
                    v-model="scope.row.size"
                    type="number"
                    style="width: 150px;"
                    :min="minDiskSize"
                    :step="10"
                />
              </el-tooltip>
              <el-input-number
                  v-else
                  controls-position="right"
                  :precision="0"
                  v-model="scope.row.size"
                  type="number"
                  style="width: 150px;"
                  :min="minDiskSize"
                  :step="10"
              />
            </template>
          </el-table-column>
          <el-table-column label="指定块存储">
            <template #default="scope">
              <div v-if="scope.$index==0" style="display: flex">
                <el-switch v-model="isUseLocalBlock"
                           :active-value="true"
                           :inactive-value="false"
                           class="ml-2"
                           inline-prompt
                           @click="modStatus('local')"/>
                <el-select style="width: 160px;margin-left: 10px" v-if="'storage_id' in scope.row"
                           v-model="scope.row.storage_id">
                  <el-option
                      v-for="(item,index) in (scope.row.backend=='local'?localStoList:scope.row.backend=='nfs'?nfsStoList:cephStoList)"
                      :label="item.name+' 总量:'+item.usedInfo.capacity+'/已分配:'+item.usedInfo.virtualUsedCapacity+'/已使用:'+item.usedInfo.actualCapacityUsed"
                      :value="item.id"/>
                </el-select>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item label="数据盘(GB)" style="padding-left: 10px" label-width="110" prop="disks">
        <div style="display: flex;flex-direction: column;width: 100%">
          <div style="display: block;width: 100%">
            <el-table :data="dataDiskList"
                      :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
              <el-table-column label="硬盘类别">
                <template #default="scope">
                  <el-select v-if="scope.$index==0" style="width: 200px;" v-model="scope.row.backend"
                             @change="HandleDiskChange('data')">
                    <el-option :value="'local'" label="本地硬盘"/>
                    <el-option :value="'nfs'" v-if="nfsStoList.length>0" label="云硬盘(NFS)"/>
                    <el-option :value="'rbd'" v-if="cephStoList.length>0" label="云硬盘(Ceph)"/>
                  </el-select>
                  <el-tag v-if="scope.$index>0" style="width: 200px;" type="primary">
                    {{
                      scope.row.backend == 'local' ? '本地硬盘' : scope.row.backend == 'nfs' ? '云硬盘(NFS)' : '云硬盘(Ceph)'
                    }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="容量(GB)">
                <template #default="scope">
                  <el-input-number controls-position="right" :precision="0" v-model="scope.row.size" :min="10"
                                   :step="10"></el-input-number>
                </template>
              </el-table-column>
              <el-table-column label="指定块存储">
                <template #default="scope">
                  <div v-if="scope.$index==0" style="display: flex">
                    <el-switch v-model="isUseBlock"
                               :active-value="true"
                               :inactive-value="false"
                               class="ml-2"
                               inline-prompt
                               @click="modStatus('data')"/>
                    <el-select @change="SetDataStorage" style="width: 160px;margin-left: 10px"
                               v-if="'storage_id' in scope.row"
                               v-model="scope.row.storage_id">
                      <el-option
                          v-for="(item1,index1) in (scope.row.backend=='local'?localStoList:scope.row.backend=='nfs'?nfsStoList:cephStoList)"
                          :label="item1.name+' 总量:'+item1.usedInfo.capacity+'/已分配:'+item1.usedInfo.virtualUsedCapacity+'/已使用:'+item1.usedInfo.actualCapacityUsed"
                          :value="item1.id"/>
                    </el-select>
                  </div>
                  <el-tag v-if="scope.$index>0&&'storage_id' in scope.row"
                          style="width: 200px;height: 30px;margin-left: 10px">
                    {{ dataStorageBlock.name }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
          <div style="display: flex;width: 100%;margin-top: 5px;justify-content: space-between">
            <el-button style="width: 100%" @click="AddNewDisk">添加新磁盘</el-button>
            <el-button style="width: 100%" v-if="dataDiskList.length != 0" @click="RemoveNewDisk">移除磁盘</el-button>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <div class="flex" style="justify-content: right;margin-top: 24px">
      <el-button type="info" @click="CloseDrawer">取消</el-button>
      <el-button
          type="primary"
          @click="SubmitDrawer"
      >确定
      </el-button>
    </div>
  </el-drawer>
</template>

<script>
import {addHost, getGpuList,} from "@/api/cloudpods/compute/servers";
import {getNetIpList} from "@/api/cloudpods/network/ip";
import {getImageList} from "@/api/cloudpods/image/image";
import {getStorageList} from "@/api/cloudpods/storage/storage";
import {getRenterList} from "@/api/yunguan/run/renter";
import {getVpcList} from "@/api/cloudpods/network/vpc";
import _ from "lodash";

export default {
  props: {
    hostData: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      isUseBlock: false,
      isUseLocalBlock: false,
      isUpdateRenterRes: false,
      modRenterRes: {
        id: '',
        resourceID: '',
        startTime: '',
        endTime: '',
      },
      imageList: [],
      renterList: [],
      minDiskSize: 10,
      skuData: {
        cpu: 2,
        mem: 2,
      },
      isEip:false,
      createVisible: false,
      IsCreated: true,
      serverDataInfo: {
        disks: [{
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: 10,
        }],
        duration: '',
        timeType: 'h',
        timeNum: 1,
        network: '',
        imageId: '',
        generate_name: '',
        nets: [],
        sku: '',
        eip_bw: 30,
        renterID: '',
        vcpu_count: 0,
        vmem_size: 0,
        isolated_devices: [],
        vpcId: 'default',
      },
      addServerDataInfo: {
        disks: [{
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: 10,
        }],
        generate_name: '',
        nets: [],
        sku: '',
        eip_bw: 30,
        vcpu_count: 0,
        vmem_size: 0,
        isolated_devices: [],
        vpcId: 'default',
      },
      diskDataPrototype: {
        backend: '',
        disk_type: '',
        index: 0,
        medium: '',
        size: 0,
      },
      netDataPrototype: {
        network: '',
      },
      IpList: [],
      rules: {
        generate_name: [
          {required: true, message: '请输入云主机名称', trigger: 'blur'},
        ],
        serverID: [
          {required: true, message: '请选择物理机', trigger: 'blur'}
        ],
        renterID: [
          {required: true, message: '请选择租户', trigger: 'blur'}
        ],
        sku: [
          {required: true, message: '请选择套餐', trigger: 'blur'}
        ],
        vpcId: [
          {required: true, message: '请选择VPC', trigger: 'blur'}
        ],
        imageId: [
          {required: true, message: '请选择镜像', trigger: 'blur'}
        ],
        network: [
          {required: true, message: '请选择IP子网', trigger: 'blur'}
        ],
        username: [
          {required: true, message: '请输入用户名', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'},
          {
            validator: (rule, value, callback) => {
              if (value.length < 12 || value.length > 30) {
                callback(new Error('密码长度必须在12到30个字符之间'));
              } else if (value.startsWith('/')) {
                callback(new Error('密码不能以“/”开头'));
              } else {
                const hasUpperCase = /[A-Z]/.test(value);
                const hasLowerCase = /[a-z]/.test(value);
                const hasNumbers = /\d/.test(value);
                const hasSpecialChars = /[~`!@#$%^&*()\-_=+$${}|:'",.<>?]/.test(value);

                if (!hasUpperCase || !hasLowerCase || !hasNumbers || !hasSpecialChars) {
                  callback(new Error('密码必须同时包含大写字母、小写字母、数字和特殊字符'));
                } else {
                  callback();
                }
              }
            },
            trigger: 'blur'
          }
        ],
        name: [
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],

      },
      serverList: [],
      dataDiskList: [],
      dataStorageBlock: null,
      resetData: {
        id: '',
        username: '',
        password: '',
        type: false,
      },
      localStoList: [],
      nfsStoList: [],
      cephStoList: [],
      openDevice: false,
      selectDevice: [
        {
          index: 0,
          num: 1,
        }
      ],
      deviceInfoArray: [],
      specDetail: [],
      renterVisible: false,
      isBindRenter: false,
      vpcList: [],
    }
  },
  created() {
    this.GetAllRenterList()
    this.GetVpcList()
    this.GetImageList()
  },
  methods: {
    modStatus(type) {
      if (type == 'local') {
        if (this.isUseLocalBlock) {
          this.SetSystemStorage('local');
        } else if (!this.isUseLocalBlock) {
          this.RemoveSystemStorage('local');
        }
      } else {
        if (this.isUseBlock) {
          this.SetSystemStorage('data');
        } else if (!this.isUseBlock) {
          this.RemoveSystemStorage('data');
        }
      }
    },
    //获取可使用的GPU信息
    GetGpuList() {
      getGpuList().then(res => {
        if (res.code == 0) {
          let dataTemp = res.data;
          for (const category in dataTemp) {
            if (dataTemp.hasOwnProperty(category)) {
              const items = dataTemp[category];
              let maxKey = null;
              let maxValue = -Infinity;

              // 找到最大值及其对应的键
              for (const key in items) {
                if (items.hasOwnProperty(key)) {
                  if (items[key] > maxValue) {
                    maxValue = items[key];
                    maxKey = key;
                  }
                }
              }

              // 将最大值的结果存储到 dataList 中
              if (maxKey !== null) {
                this.deviceInfoArray.push({
                  model: category, // 分类标识
                  key: maxKey,        // 最大值对应的键名
                  count: maxValue,    // 最大值
                  vendor: 'NVIDIA'
                });
              }
            }
          }
        }
      })
    },
    //获取所有租户列表
    GetAllRenterList() {
      let obj = {};
      obj.page = 1;
      obj.pageSize = 1000;
      getRenterList(obj).then(res => {
        if (res.code == 0) {
          this.renterList = res.data.list;
        }
      })
    },
    //磁盘类型切换
    HandleDiskChange(type, index) {
      if (type == 'local') {
        if ('storage_id' in this.serverDataInfo.disks[0]) {
          delete this.serverDataInfo.disks[0].storage_id;
        }
      } else if (type == 'data') {
        for (let i = 0; i < this.dataDiskList.length; i++) {
          this.dataDiskList[i].backend = this.dataDiskList[0].backend;
          if ('storage_id' in this.dataDiskList[i]) {
            delete this.dataDiskList[i].storage_id;
          }
        }
      }
    },
    //指定本地存储
    SetSystemStorage(type, index) {
      if (type == 'local') {
        this.serverDataInfo.disks[0].storage_id = ''
      } else if (type == 'data') {
        this.dataDiskList[0].storage_id = '';
      }
    },
    //设置数据盘本地存储
    SetDataStorage() {
      for (let i = 0; i < this.dataDiskList.length; i++) {
        this.dataDiskList[i].storage_id = this.dataDiskList[0].storage_id;
      }
      const tempList = this.dataDiskList[0].backend == 'local' ? this.localStoList : this.nfsStoList;
      this.dataStorageBlock = tempList.find(item => item.id == this.dataDiskList[0].storage_id);
    },
    //移除指定本地存储
    RemoveSystemStorage(type, index) {
      if (type == 'local') {
        delete this.serverDataInfo.disks[0].storage_id;
      } else if (type == 'data') {
        for (let i = 0; i < this.dataDiskList.length; i++) {
          delete this.dataDiskList[i].storage_id;
        }
      }
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
    OpenCreateDrawer() {
      this.GetStorageList();
      this.GetGpuList();
      this.GetImageList();
      this.GetIpList();
      this.serverDataInfo = {};
      this.serverDataInfo.isolated_devices = [];
      this.openDevice = false;
      this.serverDataInfo.generate_name = "";
      this.serverDataInfo.imageId = "";
      this.serverDataInfo.vpcId = "default";
      this.serverDataInfo.network = "";
      this.skuData.cpu = this.hostData.cpu;
      this.skuData.mem = this.hostData.mem;
      this.dataDiskList = [];
      this.serverDataInfo.disks = [];
      this.serverDataInfo.nets = [];
      this.serverDataInfo.eip_bw = 30;
      this.serverDataInfo.disks = [
        {
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: this.hostData.sysDisk,
        }
      ];
      this.serverDataInfo.renterID = this.hostData.renterID;
      this.generateUniqueName(this.hostData.renterID)
      this.isBindRenter = true;
      if (this.hostData.dataDisk > 0) {
        this.dataDiskList.push({
          backend: 'local',
          disk_type: 'data',
          index: this.dataDiskList.length + 1,
          medium: 'rotate',
          size: this.hostData.dataDisk,
        })
      }
      this.serverDataInfo.timeNum = this.hostData.timeNum;
      switch (this.hostData.timeType) {
        case 1:
          this.serverDataInfo.timeType = 'h';
          break;
        case 2:
          this.serverDataInfo.timeType = 'd';
          break;
        case 8:
          this.serverDataInfo.timeType = 'm';
          break;
        case 32:
          this.serverDataInfo.timeType = 'y';
          break;
        default:
          this.serverDataInfo.timeType = 'h';
          break;
      }
      this.minDiskSize = 10;
      this.createVisible = true;
    },
    // 获取镜像列表
    GetImageList() {
      let obj = {}
      obj.pageIndex = 1;
      obj.pageSize = 10000;
      obj.status = 'active';
      getImageList(obj).then(response => {
        if (response.code == 0) {
          this.imageList = response.data.data
        }
      })
    },
    //关闭对话框
    CloseDrawer() {
      this.visible = false
      this.passVisible = false
      this.dialogVisible = false
      this.createVisible = false
      this.snapVisible = false
      this.renterVisible = false
      this.secGroupVisible = false
    },
    //提交对话框
    SubmitDrawer() {
      //格式化套餐名称
      {
        let skuText = "ecs.g1.";
        skuText = skuText + "c" + this.skuData.cpu + "m" + this.skuData.mem;
        this.serverDataInfo.sku = skuText;
        this.serverDataInfo.vcpu_count = this.skuData.cpu * 1
        this.serverDataInfo.vmem_size = this.skuData.mem * 1024
      }
      this.serverDataInfo.isolated_devices = [];
      if (this.openDevice) {
        for (let item = 0; item < this.selectDevice.length; item++) {
          for (let i = 0; i < this.selectDevice[item].num; i++) {
            this.serverDataInfo.isolated_devices.push({
              model: this.deviceInfoArray[this.selectDevice[item].index].model,
              vendor: this.deviceInfoArray[this.selectDevice[item].index].vendor,
            })
          }
        }
      }
      if (!this.isEip){
        this.serverDataInfo.eip_bw=0;
      }
      this.netDataPrototype.network = this.serverDataInfo.network;
      this.serverDataInfo.nets = [];
      this.serverDataInfo.nets.push(this.netDataPrototype)
      this.serverDataInfo.duration = this.serverDataInfo.timeNum + this.serverDataInfo.timeType;
      this.serverDataInfo.disks[0].image_id = this.serverDataInfo.imageId;
      this.addServerDataInfo = _.cloneDeep(this.serverDataInfo)
      for (let i = 0; i < this.addServerDataInfo.disks.length; i++) {
        this.addServerDataInfo.disks[i].size = this.addServerDataInfo.disks[i].size * 1024
      }
      this.AddHost()
    },
    //创建云主机
    AddHost() {
      this.$refs.serverDataInfo.validate(valid => {
        if (valid) {
          addHost(this.addServerDataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createVisible = false
              if (this.isBindRenter) {
                this.$emit('dialogSubmit', JSON.parse(response.data).id, 2);
              }
            }
          })
        }
      })
    },
    generateRandomLetter() {
      // 随机生成一个字母（小写字母）
      const letters = "abcdefghijklmnopqrstuvwxyz";
      const randomIndex = Math.floor(Math.random() * letters.length);
      return letters[randomIndex];
    },

    timestampToLetters(timestamp) {
      // 将时间戳转换为字符串，然后通过 ASCII 码转换为字母
      let timestampStr = timestamp.toString();
      // 确保至少有8位数字，不足的部分在前面补零
      if (timestampStr.length < 8) {
        timestampStr = timestampStr.padStart(8, '0');
      }
      // 截取后8位数字
      if (timestampStr.length > 8) {
        timestampStr = timestampStr.slice(-8);
      }
      let result = "";
      for (let char of timestampStr) {
        // 将字符转换为整数
        const digit = parseInt(char, 10);
        // 加上97得到新的ASCII值
        const asciiVal = digit + 97;
        // 转换为对应的ASCII字符
        result += String.fromCharCode(asciiVal);
      }
      return result;
    },

    generateUniqueName(renterId) {
      // 生成纳秒级时间戳并转换为字符串，取前 8 位以控制长度
      const timeStr = this.timestampToLetters(Date.now());

      // 将租户 ID 转换为字符串
      const renterIdStr = renterId.toString();

      // 组合时间戳、随机字符串和租户 ID
      let uniqueName = `${timeStr}-${renterIdStr}-`;

      // 确保名称不超过 20 个字符
      if (uniqueName.length > 20) {
        uniqueName = uniqueName.slice(0, 20);
      }

      // 随机生成一个字母作为结尾
      if (!/[a-zA-Z]$/.test(uniqueName)) {
        uniqueName = uniqueName.slice(0, -1) + this.generateRandomLetter();
      }
      this.serverDataInfo.generate_name = uniqueName;
    },
    //选择套餐后
    skuChange(value) {
      const minDiskSizeTemp = Number(this.imageList.find(item => item.id == this.serverDataInfo.imageId).min_disk / 1024).toFixed(0)
      this.minDiskSize = minDiskSizeTemp < 10 ? 10 : minDiskSizeTemp;
      if (this.hostData.sysDisk < this.minDiskSize) {
        this.serverDataInfo.disks[0].size = this.minDiskSize;
      } else {
        this.serverDataInfo.disks[0].size = this.hostData.sysDisk;
      }
    },
    //获取VPC列表
    GetVpcList() {
      let obj = {};
      obj.pageSize = 0;
      getVpcList(obj).then(response => {
        if (response.code == 0) {
          this.vpcList = response.data.data
        }
      })
    },
    //获取IP列表
    GetIpList() {
      let obj = {};
      obj.pageSize = 0;
      obj.vpcId = this.serverDataInfo.vpcId;
      obj.server_type = ["guest"];
      getNetIpList(obj).then(response => {
        if (response.code == 0) {
          this.IpList = response.data.data
        }
      })
    },
  },
}
</script>
<style scoped>

.input-with-unit {
  position: relative;
  display: inline-block;
  margin-left: 5px;
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