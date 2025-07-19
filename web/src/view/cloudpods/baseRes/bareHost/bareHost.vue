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
            <el-button  @click="GetBaseBareHostList" type="primary">
              查询
            </el-button>
            <el-button class="button-gap" type="info" @click="searchInfo.name='';GetBaseBareHostList()">重置</el-button>
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

        <el-tooltip :disabled="!deleteList.includes(false)" content="物理机必须在禁用状态下才可进行此操作"
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
          :data="hostList"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="30" />
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="OpenMoreDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <i v-if="scope.row.status=='running'" class="status-dot"/>
            <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
            <i v-else-if="scope.row.status == 'ready'||scope.row.status == 'unknown'" style="background-color: rgb(174,185,192)" class="status-dot"/>
            <el-icon v-else class="status-dot-other">
              <Loading />
            </el-icon>
            <span>{{ hostStatus[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="启用状态"width="80">
          <template #default="scope">
            <i v-if="scope.row.enabled" class="status-dot"/>
            <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
            <span>{{ scope.row.enabled?'启用':'禁用' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP"align="left" width="180px">
          <template #default="scope">
            <span>{{ scope.row.ipmi_ip }}(带外)</span><br/>
            <span v-if="scope.row.status!='prepare'">{{ scope.row.access_ip }}(管理)</span>
          </template>
        </el-table-column>
        <el-table-column label="IPMI">
          <template #default="scope">
            <el-popover placement="top" :width="popoverWidth" trigger="click">
              <template #reference>
                <el-button @click="HandleGetImpiInfo(scope.row)" type="text" style="margin-right: 16px">查看详情</el-button>
              </template>
              <span>IP:{{impiData.ip}}</span><br/>
              <span>用户名:{{impiData.username}}</span><br/>
              <span>密码:{{impiData.password }}</span>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="初始账号" width="120">
          <template #default="scope">
            <el-popover placement="top" :width="popoverWidth" trigger="click">
              <template #reference>
                <el-button @click="HandleGetLoginInfo(scope.row)" type="text" style="margin-right: 16px">查看详情</el-button>
              </template>
              <span>用户名:{{impiData.username}}</span><br/>
              <span>密码:{{impiData.password }}</span>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="规格" width="180px">
          <template #default="scope">
            <span>CPU:{{ scope.row.cpu_count }}</span><br/>
            <span>内存:{{ FormatCapacitySize(scope.row.mem_size) }}</span><br/>
            <span>磁盘:{{ FormatCapacitySize(scope.row.storage_size)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="品牌">
          <template #default="scope">
            <span>{{ scope.row.sys_info.manufacture }}</span>
          </template>
        </el-table-column>
        <el-table-column label="型号">
          <template #default="scope">
            <span>{{ scope.row.sys_info.model }}</span>
          </template>
        </el-table-column>
        <el-table-column label="分配">
          <template #default="scope">
            <span>{{ scope.row.server }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="scope">
            <el-button link v-if="scope.row.enabled" @click="HandleModStatus(scope.row,'disable','禁用')" type="text">禁用</el-button>
            <el-button link v-if="!scope.row.enabled" @click="HandleModStatus(scope.row,'enable','启用')" type="text">启用</el-button>
            <el-button link v-if="scope.row.status=='running'" @click="HandleModStatus(scope.row,'stop','关机')" type="text">关机</el-button>
            <el-button link v-if="scope.row.status!='running'" :disabled="scope.row.status!='ready'" @click="HandleModStatus(scope.row,'start','开机')" type="text">开机</el-button>
            <el-tooltip :disabled="!scope.row.enabled" content="物理机必须在禁用状态下才可进行此操作"
                        placement="top">
              <el-button link :disabled="scope.row.enabled" @click="HandleDeleteDialog(scope.row,'delete')" type="text">删除</el-button>
            </el-tooltip>
            <el-tooltip :disabled="scope.row.status=='running'" content="物理机必须在开机状态下才可进行此操作"
                        placement="top">
              <el-button link :disabled="scope.row.status!='running'" @click="HandleRemoteConn(scope.row)" type="text">远程连接</el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="searchInfo.total"
            layout="total,  prev, pager, next, sizes,jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!--新增物理机-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="520"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增物理机</span>
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
          :rules="rules"
          label-width="130px"
          label-position="left"
      >
        <el-form-item label="添加方式" label-width="120px" style="padding-left: 10px">
          <el-tag type="primary">PXE引导注册</el-tag>
        </el-form-item>
        <el-form-item label="物理机名称" prop="name">
          <el-input placeholder="请输入名称" v-model="dataInfo.name"/>
        </el-form-item>
        <el-form-item label="IPMI地址" prop="ip" style="margin-top: 30px">
          <el-input placeholder="请输入已配置好的BMC的信息" v-model="dataInfo.ip"/>
        </el-form-item>
        <el-form-item label="IPMI用户名" prop="username">
          <el-input placeholder="请输入已配置好的BMC的信息" v-model="dataInfo.username"/>
        </el-form-item>
        <el-form-item label="IPMI密码" prop="password">
          <el-input placeholder="请输入已配置好的BMC的信息" v-model="dataInfo.password"/>
        </el-form-item>
        <el-form-item label="管理口MAC地址" prop="mac">
          <el-input placeholder="请输入进行PXE启动网卡的MAC地址" v-model="dataInfo.mac"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-tooltip :disabled="licenseConfig.deviceNum>searchInfo.total" content="当前系统设备数量已达上限"
                    placement="top">
          <el-button type="primary" :disabled="licenseConfig.deviceNum<=searchInfo.total" @click="SubmitDrawer">确 定</el-button>
        </el-tooltip>
      </div>
    </el-drawer>
    <!--      详情-->
    <el-drawer v-model="visible" :close-on-press-escape="true" :show-close="false" size="920px">

      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">详情</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="visible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <div style="margin: 0 10px">
        <el-tabs v-model="activeName">
          <el-tab-pane name="detail" label="详情">
            <el-collapse v-model="activeNames" expand-icon-position="left" accordion>
              <!-- 基本信息折叠面板 -->
              <el-collapse-item :name="1">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">基本信息</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="名称">{{detailData.name}}</el-descriptions-item>
                    <el-descriptions-item label="IP">{{detailData.access_ip}}</el-descriptions-item>
                    <el-descriptions-item label="MAC地址">{{detailData.access_mac}}</el-descriptions-item>
                    <el-descriptions-item label="创建时间">{{FormatDateTime(detailData.created_at)}}</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>

              <!-- 品牌信息折叠面板 -->
              <el-collapse-item :name="2">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">品牌信息</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="名称">{{detailData.sys_info.manufacture}}</el-descriptions-item>
                    <el-descriptions-item label="型号">{{detailData.sys_info.model}}</el-descriptions-item>
                    <el-descriptions-item label="序列号">{{detailData.sys_info.sn}}</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>

              <!-- CPU信息折叠面板 -->
              <el-collapse-item :name="3">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">CPU</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="物理CPU">{{detailData.cpu_count}}核</el-descriptions-item>
                    <el-descriptions-item label="虚拟CPU">{{detailData.cpu_count*detailData.cpu_commit_bound}}核</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>

              <!-- 内存信息折叠面板 -->
              <el-collapse-item :name="4">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">内存</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="容量">{{FormatCapacitySize(detailData.mem_size)}}</el-descriptions-item>
                    <el-descriptions-item label="系统预留">{{FormatCapacitySize(detailData.mem_reserved)}}</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>

              <!-- 存储信息折叠面板 -->
              <el-collapse-item :name="5">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">存储</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="容量">{{FormatCapacitySize(detailData.storage_size)+'(物理机总量)'}}</el-descriptions-item>
                    <el-descriptions-item label="类型">{{detailData.storage_type=='rotate'?'机械盘':'数据盘'}}</el-descriptions-item>
                  </el-descriptions>
                  <div style="margin-top: 10px">
                    <el-table border style="font-size: 10px" :data="detailData.spec.storage_info">
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
                          {{FormatCapacitySize(scope.row.size)}}
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </el-tab-pane>
          <el-tab-pane name="log" label="操作日志">
            <Logs :log-data="logData"/>
          </el-tab-pane>
        </el-tabs>
      </div>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="visible=false">取 消</el-button>
      </div>
    </el-drawer>

  </div>
</template>

<script>
import {
  addBareHost,
  deleteBareHost, deleteBareHostByIds,
  getBaseBareHostList, getImpiInfo, getLoginInfo, getWebConsoleSSH, updateBareHostStatus
} from "@/api/cloudpods/baseRes/host"
import Logs from "@/view/cloudpods/component/log.vue";
import statusType from "@/locales/zh-CN.json";
import {queryLicense} from "@/api/yunguan/config/licenseConfig";


export default {
  components: {Logs},
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
      deleteType: '',
      ids: [],
      deleteList: [],
      licenseConfig:{},
      detailData:{},
      activeName:'detail',
      activeNames:[1,2,3,4,5],
      logData:{
        id:'',
        isContainer:false,
        ownerKind:'',
        ownerName:'',
        namespace:'',
        clusterId:'',
      },
      visible:false,
      timer:null,
      popoverWidth:280,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name:'',
        total: 0,
      },
      hostList:[],
      hostStatus: statusType.status.host,
      impiData:{},
      updateData:{
        id:'',
        status:'',
      },
      dataInfo:{
        mac:'',
        ip:'',
        password:'',
        username:'',
        name:'',
      },
      createVisible:false,
      rules:{
        name:[
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        mac:[
          {required: true, message: '请输入进行PXE启动网卡的MAC地址', trigger: 'blur'},
          {
            pattern: /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/, // 正则表达式，匹配标准的MAC地址格式
            message: '无效的MAC地址格式,请输入正确的MAC地址',
            trigger: 'blur'
          }
        ],
        ip:[
          {required: true, message: '请输入已配置好的BMC的信息', trigger: 'blur'},
          {
            pattern: /^((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))\.((\d{1,2})|(1\d{2})|(2[0-4]\d)|(25[0-5]))$/, // 正则表达式，匹配标准的IPv4地址格式
            message: '无效的IP地址格式，请输入正确的IP地址',
            trigger: 'blur'
          }
        ],
        password:[
          {required: true, message: '请输入已配置好的BMC的信息', trigger: 'blur'}
        ],
        username:[
          {required: true, message: '请输入已配置好的BMC的信息', trigger: 'blur'}
        ],
      },
    }
  },
  created(){
    this.GetBaseBareHostList()
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
      this.deleteId = val;
      this.deleteType = type;
      this.deleteVisible = true;
      switch (type) {
        case 'delete':
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该物理机，是否继续?';
          break;
      }
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].scope.row.enabled){
          this.deleteList.push(false)
        }
      }
    },
    //获取License
    GetLicenseConfig() {
      queryLicense().then(res => {
        if (res.code == 0) {
          this.licenseConfig = res.data;
        }
      })
    },
    //打开详情对话框
    OpenMoreDrawer(val){
      this.detailData=val;
      this.logData.id=val.id;
      this.visible=true;
    },
    //开启定时器
    startInterval() {
      if (this.timer==null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetBaseBareHostList();
        }, 5000);
      } else {
        console.log('定时器已经在运行');
      }
    },
    //关闭定时器
    stopInterval() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
        console.log('定时器已清除');
      }
    },
    //新增对话框
    OpenDrawer(){
      this.dataInfo={};
      this.GetLicenseConfig();
      this.createVisible=true;
    },
    //取消
    CloseDrawer(){
      this.createVisible=false;
      this.visible=false;
    },
    //确定
    SubmitDrawer(){
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addBareHost(this.dataInfo).then(res=>{
            if (res.code == 0){
              this.$message({
                type: 'success',
                message:'创建成功',
              })
              this.createVisible=false;
              this.GetBaseBareHostList();
            }
          })
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
    //格式化容量
    FormatCapacitySize(val) {
      var G = 1024
      var T = 1024 * 1024
      if (val < G) {
        return (val / G).toFixed(1) + "G";
      } else if (val >= G && val <= T) {
        return (val / G).toFixed(0) + "G";
      } else if (val > T) {
        return (val / T).toFixed(1) + "T";
      }
    },
    GetBaseBareHostList() {
      getBaseBareHostList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData=response.data.data;
          if(JSON.stringify(tempData) !== JSON.stringify(this.hostList)){
            this.hostList = response.data.data
            this.searchInfo.total = response.data.total
            this.startInterval()
          }else{
            this.stopInterval()
          }
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetBaseBareHostList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetBaseBareHostList()
    },
    //获取IPMI信息
    HandleGetImpiInfo(val){
      getImpiInfo({id:val.id}).then(res=>{
        if (res.code == 0){
          this.impiData = JSON.parse(res.data)
          const usernameLength = this.impiData.username.length;
          const passwordLength = this.impiData.password.length;

          // 根据用户名和密码的长度计算宽度
          const baseWidth = 10; // 基础宽度
          const additionalWidthPerChar = 9; // 每个字符增加的宽度

          const totalLength = usernameLength + passwordLength;
          const dynamicWidth = baseWidth + (totalLength * additionalWidthPerChar);
          this.popoverWidth=dynamicWidth;
        }
      })
    },
    //获取登录信息
    HandleGetLoginInfo(val){
      let obj={}
      obj.id="";
      if (val.server!=undefined){
        obj.id=val.server_id;
        obj.status="server"
      }else{
        obj.id=val.id;
        obj.status="host";
      }
      getLoginInfo(obj).then(res=>{
        if (res.code == 0){
          this.impiData = JSON.parse(res.data)
          const usernameLength = this.impiData.username.length;
          const passwordLength = this.impiData.password.length;

          // 根据用户名和密码的长度计算宽度
          const baseWidth = 10; // 基础宽度
          const additionalWidthPerChar = 10; // 每个字符增加的宽度

          const totalLength = usernameLength + passwordLength;
          const dynamicWidth = baseWidth + (totalLength * additionalWidthPerChar);
          this.popoverWidth=dynamicWidth;
        }
      })
    },
    //修改物理机状态
    HandleModStatus(val,status,text){
      this.updateData.id=val.id;
      this.updateData.status=status;
      updateBareHostStatus(this.updateData).then(res=>{
        if (res.code == 0){
          this.GetBaseBareHostList()
          this.$message({
            type: 'success',
            message:text+'成功',
          })
        }
      })
    },
    //删除物理机
    HandleDelete(val){
        this.updateData.id=val.id;
        deleteBareHost(this.updateData).then(res=>{
          if (res.code == 0){
            this.GetBaseBareHostList()
            this.$message({
              type: 'success',
              message:'删除成功',
            })
          }
        })
    },
    //删除物理机
    HandleDeleteByIds(val){
      deleteBareHostByIds({ids:this.ids}).then(res=>{
          if (res.code == 0){
            this.GetBaseBareHostList()
            this.$message({
              type: 'success',
              message:'删除成功',
            })
          }
        })
    },
    //远程连接
    HandleRemoteConn(val){
      getWebConsoleSSH({id:val.id}).then(res=>{
        if (res.code == 0){
          const data = JSON.parse(res.data)
          // const params = new URLSearchParams(data).toString();
          window.open(`${data.access_url}&session_id=${data.session}`, '_blank');
        }
      })
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
.brand-info {
  margin: 10px;
  margin-top: 5px;
}

.boxer {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}

.el-collapse-item__header {
  display: flex;
  align-items: center;
}

.el-collapse-item__icon {
  margin-right: 8px; /* 图标与标题文字之间的间距 */
}
.form-section {
  border: 1px solid #dcdfe6;
  /* 边框颜色 */
  border-radius: 4px;
  /* 圆角 */
  padding: 16px;
  /* 内边距 */
  margin-bottom: 20px;
  /* 下外边距 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  /* 阴影 */
  background-color: #fff;
  /* 背景色 */
}
</style>