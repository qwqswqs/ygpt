<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
         <div>
           <el-form-item label="名称">
             <el-input style="width: 240px" v-model="searchInfo.name" placeholder="请输入"/>
           </el-form-item>
           <el-form-item label="类型">
             <el-select style="width: 240px" v-model="searchInfo.serverType" placeholder="请选择" clearable>
               <el-option label="云主机" value="guest"/>
               <el-option label="PXE" value="pxe"/>
               <el-option label="IPMI" value="ipmi"/>
               <el-option label="弹性公网" value="eip"/>
               <el-option label="物理机" value="baremetal"/>
             </el-select>
           </el-form-item>
         </div>
          <el-form-item>
            <el-button  type="primary" @click="GetIpList">查询</el-button>
            <el-button  class="button-gap" type="info" @click="searchInfo.name='';searchInfo.serverType='';GetIpList()">重置</el-button>
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
        <el-tooltip :disabled="!deleteList.includes(false)" content="存在子网已被分配使用"
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
          :data="netIpList"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="30" />

        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column label="状态">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.status=="available"?"正常":"禁用" }}</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column label="类型">
          <template #default="scope">
            <span>{{ serverType[scope.row.server_type] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP地址" width="200px">
          <template #default="scope">
            <span> 起:{{ scope.row.guest_ip_start }}</span><br>
            <span> 止:{{ scope.row.guest_ip_end }}</span>
          </template>
        </el-table-column>
        <el-table-column label="所属VPC">
          <template #default="scope">
            <span>{{ scope.row.vpc }}</span>
          </template>
        </el-table-column>
        <el-table-column label="掩码位数" width="60">
          <template #default="scope">
            <span>{{ scope.row.guest_ip_mask }}</span>
          </template>
        </el-table-column>
        <el-table-column label="网关地址">
          <template #default="scope">
            <span>{{ scope.row.guest_gateway }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用情况">
          <template #default="scope">
            <span> 总计:{{ scope.row.ports }}</span><br>
            <span> 使用:{{ scope.row.ports_used }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button v-if="scope.row.is_default_vpc" link type="text" @click="HandleEdit(scope.row)">编辑</el-button>
            <el-tooltip :disabled="scope.row.ports_used==0" content="该子网已被分配使用"
                        placement="top">
              <el-button link :disabled="scope.row.ports_used!=0" type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除
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
    <!--      新增对话框-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="false"
        size="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{IsCreate?'新增IP子网':'编辑IP子网'}}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form
          ref="dataInfo"
          :model="dataInfo"
          :rules="IsCreate?rules:''"
          label-width="100px"
          label-position="left"
      >
        <el-form-item label="名称" prop="name">
          <el-input style="width: 320px;" placeholder="请输入名称" v-model="dataInfo.name"/>
        </el-form-item>
        <el-form-item v-if="IsCreate" label="VPC" prop="vpcId" style="margin-top: 30px">
          <el-select style="width: 320px;" placeholder="请选择VPC" @change="HandleSelectChange" v-model="dataInfo.vpcId">
            <el-option v-for="(item,index) in vpcList"
                       :value="item.id"
                       :label="item.name+'('+(item.id!='default'?item.cidr_block:'经典网络')+')'"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataInfo.type=='default'" label="服务器类型" prop="serverType"  >

          <div style="display: flex;align-items: center;background:  #F2F3F5;">
            <el-radio-group style="width: 320px;" :disabled="!IsCreate" v-model="dataInfo.serverType">
              <el-radio-button value="guest">云主机</el-radio-button>
              <el-radio-button value="baremetal">物理机</el-radio-button>
              <el-radio-button value="pxe">PXE</el-radio-button>
              <el-radio-button value="ipmi">IPMI</el-radio-button>
              <el-radio-button value="eip">弹性公网IP</el-radio-button>
            </el-radio-group>
          </div>

        </el-form-item>
        <el-form-item v-if="dataInfo.type=='default'" label="起始地址" prop="ipStart">
          <el-input style="width: 320px;" :disabled="dataInfo.vpcId!='default'" placeholder="起始IP地址" v-model="dataInfo.ipStart"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.type=='default'" label="结束地址" prop="ipEnd">
          <el-input style="width: 320px;" :disabled="dataInfo.vpcId!='default'" placeholder="结束IP地址" v-model="dataInfo.ipEnd"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.type=='default'" label="掩码位数" prop="ipMask">
          <el-select style="width: 320px;" :disabled="dataInfo.vpcId!='default'" placeholder="请选择掩码" v-model="dataInfo.ipMask">
            <el-option v-for="index in range(16,30)"
                       :value="index"
                       :label="index"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataInfo.type=='default'" label="默认网关" prop="ipGateway">
          <el-input style="width: 320px;" :disabled="dataInfo.vpcId!='default'" placeholder="请输入默认网关" v-model="dataInfo.ipGateway"/>
        </el-form-item>
        <el-form-item v-if="dataInfo.type=='other'" label="子网网段" prop="ipPrefix">
          <el-input style="width: 320px;" placeholder="请输入子网网段" v-model="dataInfo.ipPrefix"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
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
import {addNetIP, deleteNetIP, deleteNetIpByIds, getNetIpList, updateNetIP,} from "@/api/cloudpods/network/ip"
import {getVpcList} from "@/api/cloudpods/network/vpc";


export default {
  props: {
    vpcID: {
      type: String,
      default: ''
    },
  },
  computed: {
    rules() {
      const self = this;
      return {
        name: [
          {required: true, message: '请输入', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        ipPrefix: [
          {required: true, message: '请输入IP前缀', trigger: 'blur'},
          {
            validator(rule, value, callback) {
              console.log(self.vpcData)
              let result = self.validateIpAgainstSubnet(value, self.vpcData);
              if (result !== '有效的IP地址') {
                callback(new Error(result));
              } else {
                callback();
              }
            },
            trigger: 'blur'
          }
        ],
        ipMask: [
          {required: true, message: '请输入', trigger: 'blur'}
        ],
        serverType: [
          {required: true, message: '请选择', trigger: 'blur'}
        ],
        ipGateway: [
          {required: true, message: '请输入', trigger: 'blur'}
        ],
        ipStart: [
          {required: true, message: '请输入', trigger: 'blur'},
          {
            pattern: /^(((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))$/, // 正则表达式，匹配正确的IP地址
            message: '请输入正确的IP地址',
            trigger: 'blur'
          }
        ],
        ipEnd: [
          {required: true, message: '请输入', trigger: 'blur'},
          {
            pattern: /^(((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))$/, // 正则表达式，匹配正确的IP地址
            message: '请输入正确的IP地址',
            trigger: 'blur'
          }
        ],
        vpcId: [
          {required: true, message: '请选择', trigger: 'blur'}
        ],
      }
    }
  },
  data() {
    return {
      ids:[],
      deleteList:[],
      timer:null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name:'',
        vpcId: this.vpcID,
        server_type:[],
        serverType:"",
      },
      total: 0,
      netIpList: [],
      dialogVisible: false,
      IsCreate: true,
      serverType: {
        ['ipmi']: 'IPMI',
        ['guest']: '云主机',
        ['baremetal']: '物理机',
        ['eip']: '弹性公网IP',
        ['pxe']: 'PXE',
      },
      createVisible: false,
      dataInfo: {
        id: '',
        name: '',
        serverType: '',
        ipPrefix: '',
        ipMask: '',
        ipGateway: '',
        ipStart: '',
        ipEnd: '',
        type: '',
        vpcId: '',
        wireId: '',
      },
      vpcList: [],
      vpcData: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
    }
  },
  created() {
    this.GetIpList()
  },
  methods: {
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].ports_used>0){
          this.deleteList.push(false)
        }
      }
    },
    //范围
    range(start, end) {
      return Array(end - start + 1).fill().map((_, idx) => start + idx);
    },
    //校验
    validateIpAgainstSubnet(ipWithMask, targetVpc) {
      const ipv4Pattern = /^(((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))$/;
      const maskPattern = /^(\d{1,2})$/;

      const [targetNetwork, targetMask] = targetVpc.split('/');
      const targetMaskLength = parseInt(targetMask);

      const [inputNetwork, inputMask] = ipWithMask.split('/');

      if (!ipv4Pattern.test(inputNetwork)) {
        return '无效的IP地址格式';
      }

      if (!maskPattern.test(inputMask) || parseInt(inputMask) > 32 || parseInt(inputMask) < 0) {
        return '无效的子网掩码';
      }

      const inputMaskLength = parseInt(inputMask);

      if (inputMaskLength < targetMaskLength) {
        return `子网掩码必须大于或等于${targetMask}`;
      }

      const ipToInt = (ip) => ip.split('.').reduce((int, octet) => (int << 8) + parseInt(octet, 10), 0);
      const intToIp = (int) => [(int >>> 24) & 0xFF, (int >>> 16) & 0xFF, (int >>> 8) & 0xFF, int & 0xFF].join('.');

      const inputIpInt = ipToInt(inputNetwork);
      const targetNetworkInt = ipToInt(targetNetwork);
      const targetMaskInt = ((1 << targetMaskLength) - 1) << (32 - targetMaskLength);

      if ((inputIpInt & targetMaskInt) === (targetNetworkInt & targetMaskInt)) {
        return '有效的IP地址';
      } else {
        return `IP地址不属于${targetVpc}网段`;
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
    GetIpList() {
      if (this.searchInfo.serverType!=null&&this.searchInfo.serverType!=""){
        this.searchInfo.server_type[0]=this.searchInfo.serverType;
      }else{
        this.searchInfo.server_type=[];
      }
      getNetIpList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.netIpList =response.data.data
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
      this.GetIpList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetIpList()
    },
    //获取VPC列表
    GetVpcList() {
      let obj = {}
      obj.pageIndex = 1
      obj.pageSize = 10000
      getVpcList(obj).then(response => {
        if (response.code == 0) {
          this.vpcList = response.data.data
        } else {
          this.$message({
            message: '创建失败',
            type: 'warning'
          });
        }
      })
    },
    //VPC类型
    HandleSelectChange() {
      if (this.dataInfo.vpcId == 'default') {
        this.dataInfo.type = 'default';
        this.dataInfo.serverType = 'guest';
      } else {
        this.dataInfo.type = 'other';
      }
      let temp = this.vpcList.filter(item => item.id == this.dataInfo.vpcId)
      this.vpcData = temp[0].cidr_block;
      console.log(this.vpcData)
    },
    //新增
    OpenDrawer() {
      this.IsCreate = true;
      this.dataInfo = {}
      this.GetVpcList()
      this.createVisible = true
    },
    //取消
    CloseDrawer() {
      this.createVisible = false;
    },
    //确定
    SubmitDrawer() {
      //新增
      if (this.IsCreate) {
        this.$refs.dataInfo.validate(valid => {
          if (valid) {
            addNetIP(this.dataInfo).then(res => {
              if (res.code == 0) {
                this.$message({
                  message: '创建成功',
                  type: 'success'
                })
                this.createVisible = false;
                this.GetIpList()
              }
            })
          }
        })
      }
      //编辑
      else {
        this.dataInfo.ipMask=Number(this.dataInfo.ipMask)
        updateNetIP(this.dataInfo).then(res => {
          if (res.code == 0) {
            this.$message({
              message: '修改成功',
              type: 'success'
            })
            this.createVisible = false;
            this.GetIpList()
          }
        })
      }
    },
    //删除
    HandleDelete(val) {
      deleteNetIP({id: val.id}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetIpList(this.searchInfo)
        }
      })
    },
    //批量删除
    HandleDeleteByIds(val) {
      deleteNetIpByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetIpList(this.searchInfo)
        }
      })
    },

    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      switch (type) {
        case 'delete':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该云主机，是否继续?';
          break;
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该云主机，是否继续?';
          break;
      }
      this.deleteVisible = true;
    },
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
    //编辑
    HandleEdit(val) {
      this.GetVpcList()
      this.IsCreate = false;
      this.dataInfo.name = val.name;
      this.dataInfo.ipStart = val.guest_ip_start;
      this.dataInfo.ipEnd = val.guest_ip_end;
      this.dataInfo.ipMask = String(val.guest_ip_mask);
      this.dataInfo.ipGateway = val.guest_gateway;
      this.dataInfo.vpcId = val.vpc_id;
      this.dataInfo.serverType = val.server_type;
      this.dataInfo.type = 'default';
      this.dataInfo.id = val.id;
      this.dataInfo.wireId = val.wire_id;
      this.createVisible = true;
    },
  },

}
</script>
