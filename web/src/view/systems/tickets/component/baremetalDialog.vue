<script setup>

</script>

<template>
  <!--    新增和编辑-->
  <el-drawer v-model="createVisible" :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false"
             size="520px">
    <div class="flex justify-between items-center" style="margin-bottom: 24px">
      <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ '新增'}}裸金属</span>
      <el-icon
          class="cursor-pointer hover:text-gray-500 transition-colors"
          @click="CloseDrawer"
          style="width: 16px; height: 16px; font-size: 16px"
      >
        <Close />
      </el-icon>
    </div>

    <el-form ref="serverDataInfo" :model="serverDataInfo" :rules="rules" label-width="91px">
      <el-form-item label="名称" prop="generate_name">
        <el-input style="width: 320px" placeholder="这是名称" :disabled="isBindRenter" v-model="serverDataInfo.generate_name" />
      </el-form-item>
      <el-form-item label="镜像" prop="imageId">
        <el-select v-model="serverDataInfo.imageId" placeholder="请选择" style="width: 320px" @change="skuChange">
          <el-option v-for="(optionItem, optionIndex) in imageList" :key="optionIndex" :label="optionItem.name"
                     :value="optionItem.id" />
        </el-select>
      </el-form-item>
<!--      <el-form-item label="绑定租户">-->
<!--        <el-switch-->
<!--            v-model="isBindRenter"-->
<!--            :active-value="true"-->
<!--            :inactive-value="false"-->
<!--            active-text="是"-->
<!--            class="ml-2"-->
<!--            inactive-text="否"-->
<!--            inline-prompt-->
<!--            style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"-->
<!--        />-->
<!--      </el-form-item>-->
      <el-form-item prop="renterID" v-if="isBindRenter" label="租户">
        <el-select disabled @change="generateUniqueName(addRenterResInfo.renterID);" placeholder="请选择" style="width: 320px" v-model="addRenterResInfo.renterID">
          <el-option v-for="item in renterList"
                     :label="item.username"
                     :value="item.renterID"/>
        </el-select>
      </el-form-item>
      <el-form-item label="VPC" style="padding-left: 10px" prop="vpcId">
        <el-select
            v-model="serverDataInfo.vpcId"
            style="width: 320px" placeholder="请选择"
            @change="GetIpList();GetExitIpList();serverDataInfo.network='';serverDataInfo.exitIP='';">
          <el-option
              v-for="(optionItem, optionIndex) in vpcList"
              :key="optionIndex"
              :label="'名称:'+optionItem.name"
              :value="optionItem.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item  style="padding-left: 10px" label="公网IP">
        <el-select v-model="serverDataInfo.exitIP" placeholder="请选择公网IP" style="width: 100%">
          <el-option v-for="(optionItem, optionIndex) in exitIpList" :key="optionIndex"
                     :disabled="(optionItem.ports - optionItem.ports_used) == 0"
                     :label="optionItem.name + '(' + optionItem.guest_ip_start + '-' + optionItem.guest_ip_end + ')    ' + '可用：' + (optionItem.ports - optionItem.ports_used)"
                     :value="optionItem.id"/>
        </el-select>
      </el-form-item>
      <el-form-item label="IP子网" prop="network">
        <el-select style="width: 320px"  v-model="serverDataInfo.network" placeholder="请选择" >
          <el-option v-for="(optionItem, optionIndex) in IpList" :key="optionIndex"
                     :disabled="(optionItem.ports - optionItem.ports_used) == 0"
                     :label="optionItem.name + '(' + optionItem.guest_ip_start + '-' + optionItem.guest_ip_end + ')    ' + '可用：' + (optionItem.ports - optionItem.ports_used)"
                     :value="optionItem.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="物理机器" >
        <el-select style="width: 320px" placeholder="请输入"
                   @change="handleChange" v-model="serverDataInfo.host">
          <el-option v-for="baseHost in baseHostList" :disabled="baseHost.server!=undefined" :label="baseHost.name+(baseHost.server!=undefined?' 已分配':' 未分配')" :value="baseHost">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item style="padding-left: 10px"  label="磁盘配置" v-if="serverDataInfo.network !==null">
        <el-text v-show="baseHostDiskInfo === null" type="warning">
          请选择物理机
        </el-text>
        <!-- 渲染物理机的磁盘配置 -->
        <div class=" left-64" v-show="baseHostDiskInfo !== null" id="diskChart">
          <el-table border style="font-size: 10px" :data="hostDiskList">
            <el-table-column label="驱动" width="100px">
              <template #default="scope">
                {{scope.row.driver}}
              </template>
            </el-table-column>
            <el-table-column width="300px" label="型号" >
              <template #default="scope">
                {{scope.row.model}}
              </template>
            </el-table-column>
            <el-table-column label="容量">
              <template #default="scope">
                {{GetImageSize(scope.row.size)}}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-form-item>
      <el-form-item style="padding-left: 10px"  label="root密码">
        <el-input style="width: 320px" placeholder="请选择" type="password" v-model="serverDataInfo.password"></el-input>
      </el-form-item>
    </el-form>
    <div class="flex" style="justify-content: right;margin-top: 24px">
      <el-button type="info" @click="CloseDrawer">取消</el-button>
      <el-button
          type="primary"
          @click="SubmitDrawer"
      >确定</el-button>
    </div>
  </el-drawer>
</template>

<script>
import {addRenterRes, getRenterList} from "@/api/yunguan/run/renter";
import {
  createBakemeats,
} from "@/api/cloudpods/compute/servers";
import {getBaseBareHostList} from "@/api/cloudpods/baseRes/host";
import {getImageList} from "@/api/cloudpods/image/image";
import {getVpcList} from "@/api/cloudpods/network/vpc";
import {getNetIpList} from "@/api/cloudpods/network/ip";

export default {
  props: {
     renterID: {
      type: Number,
      default: 0 // 默认值是一个空对象
    },
  },
  data() {
    return {
      exitIpList: [],
      baseHostDiskInfo:null,
      serverDataInfo: {
        disks: [{
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: 30,
        }],
        generate_name: '',
        nets: [],
        sku: '',
        vcpu_count: 0,
        vmem_size: 0,
        vpcId:'',
      },
      diskDataPrototype: {
        backend: '',
        disk_type: '',
        index: 0,
        medium: '',
        size: 0,
      },
      imageId: '',
      netDataPrototype: {
        network: '',
      },
      IpList: [],
      rules: {
        generate_name: [
          { required: true, message: '请输入云主机名称', trigger: 'blur' }
        ],
        serverID: [
          { required: true, message: '请选择物理机', trigger: 'blur' }
        ],
        sku: [
          { required: true, message: '请选择套餐', trigger: 'blur' }
        ],
        imageId: [
          { required: true, message: '请选择镜像', trigger: 'blur' }
        ],
        network: [
          { required: true, message: '请选择IP子网', trigger: 'blur' }
        ],
        host: [
          { required: true, message: '请选择物理机', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
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
      },
      renterList: [],
      renterVisible: false,
      isBindRenter:true,
      vpcList:[],
      imageList:[],
      baseHostList:[],
      hostDiskList:[],
      createVisible:false,
      addRenterResInfo:{
        renterID:0,
        type:1,
        sshHost:'',
        sshUser:'',
        sshPwd:'',
        resourceID:'',
        privateIp:'',
        aiList:'[]',
      }
    }
  },
  created() {
    this.GetVpcList()
    this.GetAllRenterList()
    this.GetImageList()
  },
  methods: {
    //获取公网IP列表
    GetExitIpList() {
      let obj = {};
      obj.pageSize = 0;
      obj.vpcId = this.serverDataInfo.vpcId;
      obj.server_type = ["eip"];
      getNetIpList(obj).then(response => {
        if (response.code == 0) {
          this.exitIpList = response.data.data
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
    // 切换物理机选择处理
    handleChange() {
      console.log(this.serverDataInfo.host)
      this.baseHostDiskInfo = this.serverDataInfo.host.spec.disk
      this.hostDiskList = this.serverDataInfo.host.spec.storage_info
      console.log(this.hostDiskList)
      // this.drawDiskChart()
    },
    OpenCreateDrawer() {
      this.serverDataInfo.host=null;
      this.serverDataInfo.vpcId='default';
      this.GetIpList();
      this.GetExitIpList();
      this.GetAllRenterList();
      this.baseHostDiskInfo=null;
      // 请求物理机信息
      getBaseBareHostList({
        pageIndex: 1,
        pageSize: 10000,
      }).then((resp) => {
        this.baseHostList = resp.data.data
      })
      this.addRenterResInfo.renterID=this.renterID;
      this.generateUniqueName(this.renterID)
      this.createVisible = true;
    },
    // 获取镜像列表
    GetImageList() {
      let obj={}
      obj.pageIndex = 1;
      obj.pageSize = 10000;
      obj.description = '裸金属';
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
      this.dialogVisible = false
      this.passVisible = false
      this.createVisible = false
      this.renterVisible = false
      this.snapVisible = false
      this.tabActiveName = "详情";
    },
    ParserBaremetalDisk(bastHost){
      //解析格式BaremetalDiskConfigs
      let range = []
      console.log(bastHost)
      bastHost.spec.storage_info.forEach((item)=>{
        range.push(item.slot)
      })
      let diskConfig = [{
        conf:"raid0",
        driver:bastHost.spec.driver,
        count:bastHost.spec.storage_info.length,
        type:bastHost.storage_type,
        adapter:bastHost.spec.storage_info[0].adapter,
        range:range
      }]
      console.log(diskConfig)
      return diskConfig
    },
    //提交对话框
    SubmitDrawer() {
      this.$refs.serverDataInfo.validate(valid=>{
        if (valid){
          this.serverDataInfo.vmem_size = this.serverDataInfo.host.mem_size
          this.serverDataInfo.vcpu_count = this.serverDataInfo.host.spec.cpu
          this.serverDataInfo.disks = [
            {
              size:  this.serverDataInfo.host.storage_size - 30 * 1024,
              image_id: this.serverDataInfo.imageId,
            }, {
              size: -1,
              mountpoint: '/(系统)'
            }
          ]
          this.serverDataInfo.nets = [{
            network: this.serverDataInfo.network,
            exit: false,
          }]
          if (this.serverDataInfo.exitIP != "") {
            this.serverDataInfo.nets.push({
              network: this.serverDataInfo.exitIP,
              exit: true,
            })
          }
          this.serverDataInfo.prefer_host = this.serverDataInfo.host.id
          this.serverDataInfo.baremetal_disk_configs = this.ParserBaremetalDisk(this.serverDataInfo.host)
          createBakemeats(this.serverDataInfo).then((resp)=>{
            if(resp.code===0){
              this.$message({
                type:'success',
                message:'创建成功',
              })
              if (this.isBindRenter){
                this.$emit('dialogSubmit', JSON.parse(resp.data).id,1);
              }
              this.createVisible = false
              this.GetBaremetalList()
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
      this.serverDataInfo.generate_name=uniqueName;
    },
    //创建租户资源
    CreateRenterRes(val){
      this.addRenterResInfo.resourceID=val.id;
      addRenterRes(this.addRenterResInfo).then(res=>{
        if (res.code == 0){
          this.$message({
            type:'success',
            message:'绑定租户成功',
          })
        }
      })
    },
    GetImageSize(val) {
      return (val/1024).toFixed(0)+"GB";
    },
    //选择套餐后
    skuChange(value) {
      const match = value.match(/c(\d+)m(\d+)/);
      if (match) {
        const cpuCount = parseInt(match[1], 10); // 将匹配到的第一个括号内的数字转换为整数
        const memSize = parseInt(match[2], 10); // 将匹配到的第二个括号内的数字转换为整数
        this.serverDataInfo.vcpu_count = cpuCount
        this.serverDataInfo.vmem_size = memSize * 1024
      } else {
        console.log("No match found.");
      }
    },
    //获取VPC列表
    GetVpcList(){
      let obj = {};
      obj.pageSize = 0;
      obj.id='default';
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
      obj.vpcId=this.serverDataInfo.vpcId;
      obj.server_type = ["guest","eip"];
      getNetIpList(obj).then(response => {
        if (response.code == 0) {
          this.IpList = response.data.data
        }
      })
    },
  },
}
</script>
<style scoped lang="scss">

</style>