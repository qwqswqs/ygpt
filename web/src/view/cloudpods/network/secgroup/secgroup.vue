<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo"
                 style="display: flex;justify-content: space-between">
          <el-form-item label="名称">
            <el-input style="width: 320px" v-model="searchInfo.name" placeholder="请输入"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="GetSecgroupList">查询</el-button>
            <el-button class="button-gap" type="info" @click="searchInfo.name='';GetSecgroupList()">重置</el-button>
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
        <el-tooltip :disabled="!deleteList.includes(false)" content="安全组关联实例不为空,请先取消关联"
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
          :data="secList"
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
        <el-table-column label="状态">
          <template #default="scope">
            <i v-if="scope.row.status == 'ready'" class="status-dot"/>
            <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
            <el-icon v-else class="status-dot-other">
              <Loading/>
            </el-icon>
            <span>{{ statusType[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="关联实例" sortable="custom">
          <template #default="scope">
            <span>{{ scope.row.total_cnt }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="OpenDetailDrawer(scope.row,'in')">配置规则</el-button>
            <el-button type="text" link @click="OpenDetailDrawer(scope.row,'host')">管理云主机</el-button>
            <el-tooltip v-if="scope.row.id!='default'" :disabled="scope.row.total_cnt==0"
                        content="安全组关联实例不为空,请先取消关联" placement="top">

              <el-button link :disabled="scope.row.total_cnt!=0" type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除
              </el-button>
            </el-tooltip>
            <el-tooltip v-if="scope.row.id=='default'" content="默认安全组不能删除" placement="top">

              <el-button link :disabled="true" type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除</el-button>
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
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="40%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增安全组</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="serverDataInfo"
          :model="dataInfo"
          :rules="rules"
          label-width="60px"
          label-position="left"
      >
        <el-form-item label="名称" prop="name">
          <el-input placeholder="请输入名称" v-model="dataInfo.name"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--      详情列表-->
    <el-drawer
        v-model="ipVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        :before-close="closeDrawer"
        size="80%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">详情</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="ipVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-tabs v-model="activeName" @tabChange="tabChange">
        <el-tab-pane label="关联实例" name="host">
          <div v-if="ipVisible">
            <Server style="padding: 0;margin: 0" :sec-group-i-d="secgroupID"></Server>
          </div>
        </el-tab-pane>
        <el-tab-pane label="入方向" name="in">
          <warning-bar
              title="提示：安全组未设置任何自定义放行规则，会导致无法访问云主机端口，若需访问请添加安全组规则放行对应的端口"/>
          <el-button type="primary" @click="OpenRuleDrawer('in')">新增</el-button>
          <el-button :disabled="!ids.length" type="info" @click="handleDeleteByIds">删除</el-button>
          <el-table
              @selection-change="handleSelectionChangeData" style="margin-top: 10px" :data="secRuleList">
            <el-table-column type="selection" width="30"/>
            <el-table-column label="CIDR">
              <template #default="scope">
                {{ scope.row.cidr != undefined ? scope.row.cidr : '任意IP' }}
              </template>
            </el-table-column>
            <el-table-column label="协议">
              <template #default="scope">
                {{ scope.row.protocol }}
              </template>
            </el-table-column>
            <el-table-column label="端口">
              <template #default="scope">
                {{ scope.row.ports != undefined ? scope.row.ports : '任意端口' }}
              </template>
            </el-table-column>
            <el-table-column label="策略">
              <template #default="scope">
                {{ scope.row.action == 'allow' ? '允许' : '拒绝' }}
              </template>
            </el-table-column>
            <el-table-column label="优先级">
              <template #default="scope">
                {{ scope.row.priority }}
              </template>
            </el-table-column>
            <el-table-column label="备注">
              <template #default="scope">
                {{ scope.row.description }}
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button type="text" @click="handleDelete(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="secRuleInfo.pageIndex"
                :page-size="secRuleInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="secRuleInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </el-tab-pane>
        <el-tab-pane label="出方向" name="out">
          <warning-bar title="提示：安全组出方向默认允许所有访问，即从安全组内云主机访问外部都是放行的"/>
          <el-button type="primary" @click="OpenRuleDrawer('out')">新增</el-button>
          <el-button :disabled="!ids.length" type="info" @click="handleDeleteByIds">删除</el-button>
          <el-table @selection-change="handleSelectionChangeData" style="margin-top: 10px" :data="secRuleList">
            <el-table-column type="selection" width="30"/>
            <el-table-column label="CIDR">
              <template #default="scope">
                {{ scope.row.cidr != undefined ? scope.row.cidr : '任意IP' }}
              </template>
            </el-table-column>
            <el-table-column label="协议">
              <template #default="scope">
                {{ scope.row.protocol }}
              </template>
            </el-table-column>
            <el-table-column label="端口">
              <template #default="scope">
                {{ scope.row.ports != undefined ? scope.row.ports : '任意端口' }}
              </template>
            </el-table-column>
            <el-table-column label="策略">
              <template #default="scope">
                {{ scope.row.action == 'allow' ? '允许' : '拒绝' }}
              </template>
            </el-table-column>
            <el-table-column label="优先级">
              <template #default="scope">
                {{ scope.row.priority }}
              </template>
            </el-table-column>
            <el-table-column label="备注">
              <template #default="scope">
                {{ scope.row.description }}
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button type="text" @click="handleDelete(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="secRuleInfo.pageIndex"
                :page-size="secRuleInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="secRuleInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="log">
          <Log :log-data="logData"></Log>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
    <!--      新增和编辑规则对话框-->
    <el-drawer :show-close="false" :close-on-press-escape="true" v-model="secRuleVisible" size="720">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg"
                style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{
              isCreate ? '新建规则' : '编辑规则'
            }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="secRuleVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form label-position="left" :rules="rules" ref="secRuleData" :model="secRuleData" label-width="90px">
        <el-form-item label="类型" v-if="isCreate" label-width="80px" style="padding-left: 10px">
          <el-select style="width: 320px" @change="secRuleTypeChange" v-model="secRuleData.type"
                     placeholder="请选择类型">
            <el-option
                v-for="(item,index) in secRuleType"
                :label="item[0]"
                :value="item[1]"/>
          </el-select>
        </el-form-item>
        <el-form-item label="来源" label-width="80px" style="padding-left: 10px" prop="cidr">
          <div style="display: flex;flex-direction: column">
            <div class="cell-content">
              <el-select :disabled="secRuleData.isAnyCidr" v-model="secRuleData.cidr" filterable allow-create
                         style="width: 320px" :placeholder="secRuleData.isAnyCidr?'ALL':'请选择来源'">
                <el-option
                    v-for="(item,index) in cidrType"
                    :label="item"
                    :value="item"/>
              </el-select>
              <el-checkbox @click="secRuleData.isAnyCidr?secRuleData.cidr='0.0.0.0/0':secRuleData.cidr='';"
                           v-model="secRuleData.isAnyCidr" style="margin-left: 10px;width: 5%">任意IP
              </el-checkbox>
            </div>
            <div>
              <p class="hint">来源支持以下格式: 单个IP: 192.168.0.1(IPv4)或fd:3200::1(IPv6) CIDR:
                192.168.1.0/24(IPv4)或fd:3200::/64(IPv6)</p>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="协议" label-width="80px" style="padding-left: 10px">
          <el-select style="width: 320px" :disabled="proDisable" @change="secRuleProChange"
                     v-model="secRuleData.protocol"
                     placeholder="请选择协议">
            <el-option
                v-for="(item,index) in protocolType"
                :label="item[0]"
                :value="item[1]"/>
          </el-select>
        </el-form-item>
        <el-form-item label="端口" prop="ports">
          <div style="display: flex;flex-direction: column">
            <div class="cell-content">
              <el-input :disabled="secRuleData.isAnyPorts||portDisable" style="width: 320px" v-model="secRuleData.ports"
                        :placeholder="secRuleData.isAnyPorts?'ALL':'请输入端口号'"></el-input>
              <el-checkbox :disabled="portDisable" @click="secRuleData.ports='';" v-model="secRuleData.isAnyPorts"
                           style="margin-left: 10px;width: 5%">任意端口
              </el-checkbox>
            </div>
            <div>
              <p class="hint">协议端口支持以下格式: 单个端口: 80 多个端口: 80,443 连续端口: 3306-20000</p>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="类型策略" label-width="80px" style="padding-left: 10px">
          <el-radio-group v-model="secRuleData.action">
            <el-radio-button label="允许" value="allow"></el-radio-button>
            <el-radio-button label="拒绝" value="deny"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="优先级" label-width="80px" style="padding-left: 10px">
          <div style="display: flex;flex-direction: column">
            <div>
              <el-input-number style="width: 320px" controls-position="right" :min="1" :max="100"
                               v-model.number="secRuleData.priority"/>
            </div>
            <div>
              <span style="width: 100%" class="hint">1-100, 数字越大优先级越高</span>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="备注" label-width="80px" style="padding-left: 10px">
          <el-input style="width: 320px" placeholder="请输入备注" v-model="secRuleData.description"/>
        </el-form-item>
      </el-form>

      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="secRuleVisible=false">取 消</el-button>
        <el-button type="primary" @click="handleConfirm">提 交</el-button>
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
      <span>此操作将永久删除该项, 是否继续?</span>
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

import Log from '@/view/cloudpods/component/log.vue'
import Server from '@/view/cloudpods/compute/servers/servers.vue'
import {
  addSecgroup,
  addSecgroupRule,
  deleteSecgroup,
  deleteSecgroupByIds,
  deleteSecgroupRule,
  deleteSecgroupRuleByIds,
  getSecgroupList,
  getSecgroupRuleList,
  updateSecgroupRule
} from "@/api/cloudpods/network/secgroup";
import statusType from "@/locales/zh-CN.json";
import Info from "@/view/info/index.vue";
import WarningBar from "@/components/warningBar/warningBar.vue";
import _ from "lodash";


export default {
  components: {
    WarningBar,
    Info,
    Log, Server,
  },
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      ids: [],
      deleteList: [],
      isCreate: true,
      activeName: 'host',
      statusType: statusType.status.secgroup,
      ipVisible: false,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
      },
      logData: {
        id: '',
        isContainer: false,
        objType: 'secgroup',
      },
      total: 0,
      secList: [],
      createVisible: false,
      dataInfo: {
        name: '',
      },
      rules: {
        name: [
          {required: true, message: '请输入安全组名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
        ],
        cidr: [
          {required: false, message: '请输入IP地址', trigger: 'change'},
          {
            pattern: /^((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?:\/(?:\d|[1-2]\d|3[0-2]))?|(?:(?:[a-fA-F0-9]{1,4}:){1,7}[a-fA-F0-9]{1,4})(?:\/(?:\d|[1-9]\d|1[0-1]\d|12[0-8]))?)$/,
            message: 'IP 格式不正确，请输入正确格式',
            trigger: 'change'
          }
        ],
        ports: [
          {
            required: true,
            validator: this.validatePorts,
            trigger: 'change'
          }
        ]
      },
      secgroupID: '',
      secRuleInfo: {
        pageIndex: 1,
        pageSize: 10,
        total: 0,
        id: '',
        direction: '',
      },
      secRuleList: [],
      secRuleData: {
        id: '',
        direction: '',
        action: '',
        priority: 1,
        secgroup: '',
        type: '',
        ports: '',
        protocol: '',
        description: '',
        cidr: '',
        isAnyCidr: true,
        isAnyPorts: false,
      },
      secRuleVisible: false,
      secRuleType: {
        ['custom']: ["自定义", "custom"],
        ["windows"]: ['微软远程桌面(3389)', 'windows'],
        ["ssh"]: ['SSH(22)', 'ssh'],
        ["http"]: ['HTTP(80)', 'http'],
        ["https"]: ['HTTPS(443)', 'https'],
        ["ping"]: ['Ping', 'ping'],
      },
      cidrType: [
        '0.0.0.0/0',
        '10.0.0.0/8',
        '172.16.0.0/12',
        '192.168.0.0/16',
        '::/0',
        'fd::/64',
      ],
      proDisable: false,
      portDisable: false,
      protocolType: {
        ['tcp']: ["TCP", "tcp"],
        ["udp"]: ['UDP', 'udp'],
        ["icmp"]: ['ICMP', 'icmp'],
        ["any"]: ['任意协议', 'any'],
      }
    }
  },
  created() {
    this.GetSecgroupList()
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
        if (val[i].total_cnt > 0) {
          this.deleteList.push(false)
        }
      }
    },
    handleSortChange({column, prop, order}) {
      switch (order) {
        case "ascending":
          this.searchInfo.numDesc = "asc";
          break;
        case "descending":
          this.searchInfo.numDesc = "desc";
          break;
        default:
          this.searchInfo.numDesc = "";
          break;
      }
      this.GetSecgroupList();
    },
    //删除规则
    handleDelete(val) {
      this.$confirm('此操作将永久删除该规则, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSecgroupRule({id: val.id}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetSecGroupRuleList();
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    },
    //删除规则
    handleDeleteByIds(val) {
      this.$confirm('此操作将永久删除该规则, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSecgroupRuleByIds({ids: this.ids}).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetSecGroupRuleList();
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    },
    //编辑规则
    handleEdit(val) {
      this.secRuleData = _.cloneDeep(val);
      this.secRuleData.cidr != undefined ? this.secRuleData.isAnyCidr = false : this.secRuleData.isAnyCidr = true;
      this.secRuleData.ports != undefined ? this.secRuleData.isAnyPorts = false : this.secRuleData.isAnyPorts = true;
      this.secRuleData.ports != undefined ? '' : this.secRuleData.ports = "";
      // this.secRuleProChange()
      this.isCreate = false;
      this.secRuleVisible = true;
    },
    //规则新增编辑确定
    handleConfirm() {
      this.$refs.secRuleData.validate((valid) => {
        if (valid) {
          if (this.isCreate) {
            console.log(this.secRuleData.description)
            if (this.secRuleData.description == undefined) {
              switch (this.secRuleData.type) {
                case "ssh":
                  this.secRuleData.description = "Linux SSH登录";
                  break;
                case "windows":
                  this.secRuleData.description = "Windows远程登录";
                  break;
                case "http":
                  this.secRuleData.description = "Web服务端口（http）";
                  break;
                case "https":
                  this.secRuleData.description = "Web服务端口（https）";
                  break;
                case "ping":
                  this.secRuleData.description = "Ping服务";
                  break;
              }
            }
            addSecgroupRule(this.secRuleData).then(res => {
              if (res.code == 0) {
                this.$message({
                  type: 'success',
                  message: '创建成功',
                })
                this.GetSecGroupRuleList();
                this.secRuleVisible = false;
              }
            })
          } else {
            updateSecgroupRule(this.secRuleData).then(res => {
              if (res.code == 0) {
                this.$message({
                  type: 'success',
                  message: '编辑成功',
                })
                this.GetSecGroupRuleList();
                this.secRuleVisible = false;
              }
            })
          }
        }
      })
    },
    //端口校验
    validatePorts(rule, value, callback) {
      console.log(this.secRuleData.isAnyPorts + "," + value)
      if (this.secRuleData.isAnyPorts && value === '') {
        // 如果允许任意端口且值为空字符串，则通过校验
        callback();
      } else if (!value) {
        // 如果值为空且不允许任意端口，则提示必填
        callback(new Error('请输入端口'));
      } else if (
          !/^(\d{1,5}|\d{1,5}-\d{1,5})(,\d{1,5}(?:-\d{1,5})?)*$/.test(value)
      ) {
        // 如果值不符合端口格式，则提示格式错误
        callback(new Error('端口格式不正确，请输入正确的格式（如80、80,443或3306-20000）'));
      } else {
        // 校验通过
        callback();
      }
    },
    //协议类型变化
    secRuleProChange() {
      switch (this.secRuleData.protocol) {
        case "tcp":
        case "udp":
          this.secRuleData.isAnyPorts = false;
          this.proDisable = false;
          break;
        case "icmp":
        case "any":
          this.secRuleData.isAnyPorts = true;
          this.proDisable = true;
          break;
      }
    },
    //规则类型变化
    secRuleTypeChange() {
      switch (this.secRuleData.type) {
        case "custom":
          this.proDisable = false;
          this.portDisable = false;
          this.secRuleData.isProtocol = false;
          break;
        case "windows":
          this.proDisable = true;
          this.secRuleData.protocol = "tcp";
          this.secRuleData.ports = "3389";
          this.secRuleData.isAnyPorts = false;
          this.portDisable = true;
          break;
        case "ssh":
          this.proDisable = true;
          this.secRuleData.protocol = "tcp";
          this.secRuleData.ports = "22";
          this.secRuleData.isAnyPorts = false;
          this.portDisable = true;
          break;
        case "http":
          this.proDisable = true;
          this.secRuleData.protocol = "tcp";
          this.secRuleData.ports = "80";
          this.secRuleData.isAnyPorts = false;
          this.portDisable = true;
          break;
        case "https":
          this.proDisable = true;
          this.secRuleData.protocol = "tcp";
          this.secRuleData.ports = "443";
          this.secRuleData.isAnyPorts = false;
          this.portDisable = true;
          break;
        case "ping":
          this.proDisable = true;
          this.secRuleData.protocol = "icmp";
          this.secRuleData.ports = "";
          this.secRuleData.isAnyPorts = true;
          this.portDisable = true;
          break;
      }
    },
    //打开新增安全组规则对话框
    OpenRuleDrawer(val) {
      this.isCreate = true;
      this.proDisable = true;
      this.portDisable = false;
      this.secRuleData = {};
      this.secRuleData.secgroup = this.secRuleInfo.id;
      this.secRuleData.isAnyCidr = true;
      this.secRuleData.direction = val;
      this.secRuleData.cidr = "";
      this.secRuleData.isAnyPorts = false;
      this.secRuleData.isProtocol = true;
      this.secRuleData.protocol = "tcp";
      this.secRuleData.action = "allow";
      this.secRuleData.priority = 1;
      this.secRuleVisible = true;
    },
    //tab切换
    tabChange() {
      this.ids = [];
      switch (this.activeName) {
        case 'in':
          this.secRuleInfo.direction = "in";
          this.GetSecGroupRuleList()
          break;
        case 'out':
          this.secRuleInfo.direction = "out";
          this.GetSecGroupRuleList()
          break;
      }
    },
    //获取安全组规则列表
    GetSecGroupRuleList() {
      getSecgroupRuleList(this.secRuleInfo).then(res => {
        if (res.code == 0) {
          this.secRuleList = JSON.parse(res.data).data;
          this.secRuleInfo.total = JSON.parse(res.data).total;
        }
      })
    },
    //关闭对话框
    closeDrawer() {
      this.ipVisible = false;
      this.ids = [];
    },
    //打开子网列表
    OpenDetailDrawer(val, type) {
      this.ipVisible = true;
      this.secgroupID = val.id;
      this.logData.id = val.id;
      this.secRuleInfo.id = val.id;
      this.activeName = type;
      this.ids = [];
      if (type == 'in') {
        this.secRuleInfo.direction = type;
        this.GetSecGroupRuleList();
      }
    },
    //删除
    HandleDelete(val) {
      deleteSecgroup({id: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetSecgroupList(this.searchInfo)
        }
      })
    },
    //删除
    HandleDeleteByIds(val) {
      deleteSecgroupByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetSecgroupList(this.searchInfo)
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
      addSecgroup(this.dataInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '创建成功',
          })
          this.createVisible = false;
          this.GetSecgroupList();
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
    GetSecgroupList() {
      getSecgroupList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.secList = response.data.data
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
      this.GetSecgroupList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetSecgroupList()
    }
  }

}
</script>
<style scoped>
.cell-content {
  display: flex;
  align-content: space-between;
  align-items: center;
}

.hint {
  font-size: 12px;
  color: #165DFF; /* 使用 Element Plus 默认的文本颜色 */
  margin-top: 5px;
}

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

</style>