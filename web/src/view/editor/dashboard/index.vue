<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tab-change="handleTabsClick">
      <el-tab-pane :label="'裸金属'" name="baremetal">
        <div class="gva-table-box">
          <!-- 表格展示 -->
          <el-table
              :data="dataList"
              row-key="ID"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column label="名称" align="left" min-width="120">
              <template #default="scope">
                <el-button type="text" @click="OpenDrawer(scope.row)">{{ scope.row.name }}</el-button>
              </template>
            </el-table-column>
            <el-table-column label="密码" align="left" min-width="120">
              <template #default="scope">
                <el-popover placement="top" :width="150" trigger="click">
                  <template #reference>
                    <el-button type="text" style="margin-right: 16px">查看密码
                    </el-button>
                  </template>
                  <span>用户名:{{ scope.row.username }}</span><br/>
                  <span>密码:{{ scope.row.password }}</span>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column label="私有ip地址" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.privateIp }}</span>
              </template>
            </el-table-column>
            <el-table-column label="运行状态" align="left" min-width="120">
              <template #default="scope">
                <i v-if="scope.row.status == 'running'" class="status-dot"/>
                <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                <span>{{ hostStatus[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" align="left" :min-width="150" :fixed="'right'">
              <template #default="scope">
                <el-tooltip :disabled="scope.row.status=='running'" content="云主机必须在开机状态下才可进行此操作"
                            placement="top">
                  <el-dropdown style="margin-left: 10px">
                    <el-button type="text">
                      远程连接
                      <el-icon class="el-icon--right">
                        <arrow-down/>
                      </el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :disabled="scope.row.status!='running'" @click="HandleRemoteConn(scope.row)">
                          SOL 远程终端
                        </el-dropdown-item>
                        <el-dropdown-item :disabled="scope.row.status!='running'"
                                          @click="RemoteConnect(scope.row,'ssh')">
                          SSH 连接
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </el-tooltip>
                <el-dropdown style="margin-left: 10px">
                  <el-button type="text">
                    实例操作
                    <el-icon class="el-icon--right">
                      <arrow-down/>
                    </el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :disabled="DisableType(scope.row,'关机')" @click="HostShutDown(scope.row)">
                        关机
                      </el-dropdown-item>
                      <el-dropdown-item :disabled="DisableType(scope.row,'开机')" @click="HostBoot(scope.row)">开机
                      </el-dropdown-item>
                      <el-dropdown-item :disabled="DisableType(scope.row,'重启')" @click="HostReboot(scope.row)">重启
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchInfo.page"
                :page-size="searchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'云主机'" name="virtual">
        <div class="gva-table-box">
          <!-- 表格展示 -->
          <el-table
              :data="dataList"
              row-key="ID"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column label="名称" align="left" min-width="120">
              <template #default="scope">
                <el-button type="text" @click="OpenDrawer(scope.row)">{{ scope.row.name }}</el-button>
              </template>
            </el-table-column>
            <el-table-column label="密码" align="left" min-width="120">
              <template #default="scope">
                <el-popover placement="top" :width="150" trigger="click">
                  <template #reference>
                    <el-button type="text" style="margin-right: 16px">查看密码
                    </el-button>
                  </template>
                  <span>用户名:{{ scope.row.username }}</span><br/>
                  <span>密码:{{ scope.row.password }}</span>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column label="私有ip地址" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.privateIp }}</span>
              </template>
            </el-table-column>
            <el-table-column label="运行状态" align="left" min-width="120">
              <template #default="scope">
                <i v-if="scope.row.status == 'running'" class="status-dot"/>
                <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                <span>{{ hostStatus[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" align="left" :min-width="150" :fixed="'right'">
              <template #default="scope">
                <el-dropdown style="margin-left: 10px">
                  <el-button type="text">
                    远程连接
                    <el-icon class="el-icon--right">
                      <arrow-down/>
                    </el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :disabled="scope.row.status!='running'" @click="RemoteConnect(scope.row,'vnc')">
                        VNC 远程终端
                      </el-dropdown-item>
                      <el-dropdown-item :disabled="scope.row.status!='running'" @click="RemoteConnect(scope.row,'ssh')">
                        SSH 连接
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <el-dropdown style="margin-left: 10px">
                  <el-button type="text">
                    实例操作
                    <el-icon class="el-icon--right">
                      <arrow-down/>
                    </el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :disabled="DisableType(scope.row,'关机')" @click="HostShutDown(scope.row)">
                        关机
                      </el-dropdown-item>
                      <el-dropdown-item :disabled="DisableType(scope.row,'开机')" @click="HostBoot(scope.row)">开机
                      </el-dropdown-item>
                      <el-dropdown-item :disabled="DisableType(scope.row,'重启')" @click="HostReboot(scope.row)">重启
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchInfo.page"
                :page-size="searchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'容器'" name="container">
        <div class="gva-table-box">

          <!-- 表格展示 -->
          <el-table
              :data="dataList"
              row-key="ID"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column label="名称" align="left" min-width="120">
              <template #default="scope">
                <span type="text">{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="私有ip地址" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.privateIp }}</span>
              </template>
            </el-table-column>
            <el-table-column label="运行状态" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" align="left" :min-width="150" :fixed="'right'">
              <template #default="scope">
                <!--                <el-button type="text" @click="handleMonitor(scope.row)">监控</el-button>-->
                <el-button type="text" @click="HandleContainerRemoteConnect(scope.row)">远程连接</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchInfo.page"
                :page-size="searchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <!--    SSH连接-->
    <el-dialog
        v-model="dialogVisible"
        title="请输入SSH账户密码"
        width="40%"
    >
      <el-form
          ref="sshData"
          :model="sshData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="账号" prop="username">
          <el-input v-model="sshData.username" placeholder="请输入账号"/>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="sshData.password" placeholder="请输入密码"/>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input type="number" v-model="sshData.port" placeholder="请输入端口"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="handleConnectConfirm">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
    <!--    创建快照-->
    <el-dialog
        v-model="snapVisible"
        title="创建快照"
        width="40%"
    >
      <el-form
          ref="snapData"
          :model="snapData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="快照名称" prop="name">
          <el-input v-model="snapData.name" placeholder="请输入快照名称"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="handleSnapConfirm">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
    <!--    修改名称-->
    <el-dialog
        v-model="nameVisible"
        title="修改名称"
        width="40%"
    >
      <el-form
          ref="snapData"
          :model="snapData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="snapData.name" placeholder="请输入名称"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="handleNameConfirm">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
    <!--    监控面板-->
    <el-drawer
        v-model="openDialog"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="true"
        size="50%"
        title="监控信息">
      <el-descriptions :column="2" border title="基础信息">
        <el-descriptions-item label="名称">{{ hostBaseData.host.hostname }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ hostBaseData.host.os }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ hostBaseData.host.platform }}</el-descriptions-item>
        <el-descriptions-item label="cpu架构">{{ hostBaseData.host.kernel_arch }}</el-descriptions-item>
        <el-descriptions-item label="cpu型号">{{ hostBaseData.cpu.model }}</el-descriptions-item>
        <el-descriptions-item label="cpu核数">{{ hostBaseData.cpu.cores }}</el-descriptions-item>
        <el-descriptions-item label="内存容量">{{ hostBaseData.mem.total }}GB</el-descriptions-item>
        <el-descriptions-item label="内存已使用">{{ hostBaseData.mem.used }}GB</el-descriptions-item>
        <el-descriptions-item label="磁盘容量">{{ hostBaseData.disk.total }}GB</el-descriptions-item>
        <el-descriptions-item label="磁盘已使用">{{ hostBaseData.disk.used }}GB</el-descriptions-item>
      </el-descriptions>
      <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
        <GvaCard custom-class="h-42" title="服务器CPU利用率">
          <GvaChart :dash-data="cpuUse" :type="6"/>
        </GvaCard>
        <GvaCard custom-class="h-42" title="服务器磁盘使用率">
          <GvaChart :pie-data="diskUse" :type="5"/>
        </GvaCard>
        <GvaCard custom-class="h-42" title="服务器内存占用率">
          <GvaChart :pie-data="memoryUse" :type="5"/>
        </GvaCard>
      </div>
    </el-drawer>
    <!--云主机详情-->
    <el-drawer v-model="visible" :show-close="false" size="50%">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">详情</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="visible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <div style="margin: 0 10px">
        <el-tabs v-model="tabActiveName" @tab-change="HandleTabChange">
          <el-tab-pane name="详情" label="详情">
            <div style="margin-top: 20px">
              <el-descriptions :column="1" label-width="100px">
                <!--                    <el-descriptions-item label="ID">{{ configData.id }}</el-descriptions-item>-->
                <el-descriptions-item label="名称">{{ configData.name }}</el-descriptions-item>
                <el-descriptions-item label="IP">{{
                    configData.ips != undefined ? configData.ips : '-'
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="MAC">{{
                    configData.macs != undefined ? configData.macs : '-'
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="所属宿主机">{{
                    configData.host != undefined ? configData.host : '-'
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="所属安全组">
                  {{ configData.secgroup != undefined ? configData.secgroup : '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="所属VPC">{{
                    configData.vpc != undefined ? configData.vpc : '-'
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="系统镜像">
                  {{ configData.disks_info != undefined ? configData.disks_info[0].image : '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="CPU架构">
                  {{ configData.os_arch != undefined ? configData.os_arch : '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="CPU">{{ configData.vcpu_count }}核</el-descriptions-item>
                <el-descriptions-item label="内存">{{ configData.vmem_size / 1024 }}GB</el-descriptions-item>
                <el-descriptions-item label="GPU">
                  {{
                    configData.isolated_devices && configData.isolated_devices.length > 0
                        ? getGpuInfo(configData.isolated_devices)
                        : '-'
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="系统盘">{{ GetDiskData(configData.disks_info, 'sys') }}
                </el-descriptions-item>
                <el-descriptions-item label="数据盘">{{ GetDiskData(configData.disks_info, 'data') }}
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">{{ FormatDateTime(configData.created_at) }}
                </el-descriptions-item>
                <!--                    <el-descriptions-item label="更新时间">{{ FormatDateTime(configData.updated_at) }}-->
                <!--                    </el-descriptions-item>-->
              </el-descriptions>
            </div>
          </el-tab-pane>
          <el-tab-pane name="监控" label="监控">
            <div v-if="configData.status=='running'">
              <div v-if="!isCanInstallAgent">
                <warning-bar :title="'该实例不满足安装agent的条件'"/>
                <el-button type="text" @click="OpenSSHDialog">前往设置agent</el-button>
              </div>
              <div v-if="isCanInstallAgent">
                <warning-bar v-if="!isInstallAgent.status" :title="'该实例'+isInstallAgent.text"/>
                <el-button v-if="!isInstallAgent.status&&!isInstallAgent.text.includes('正在安装')" type="text"
                           @click="HandleInstallAgent">前往安装
                </el-button>
                <div v-if="isInstallAgent.status">
                  <div
                      style="margin-top: 10px;margin-bottom: 10px;display: flex;justify-content: space-between;align-items: center;/* 垂直居中 */gap: 10px;">
                    <div style="display: flex;align-items: center;">
                      <span>取值间隔：</span>
                      <el-select @change="handleSelectChange" style="width: 150px" v-model="monitorData.interval">
                        <el-option v-for="(item,index) in monitorConfigData[monitorData.time]"
                                   :label="convertTimeInterval(item)"
                                   :value="item"/>
                      </el-select>
                    </div>
                    <div style="display: flex;align-items: center;">
                      <el-radio-group v-model="monitorData.time" @change="handleRadioChange">
                        <el-radio-button value="1h">近1小时</el-radio-button>
                        <el-radio-button value="6h">近6小时</el-radio-button>
                        <el-radio-button value="12h">近12小时</el-radio-button>
                        <el-radio-button value="24h">近24小时</el-radio-button>
                        <!--                      <el-radio-button value="168h">近1周</el-radio-button>-->
                        <!--                      <el-radio-button value="720h">近1月</el-radio-button>-->
                      </el-radio-group>
                      <el-button @click="handleRefresh" icon="refresh" value="refresh"></el-button>
                    </div>
                  </div>
                  <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
                    <GvaCard v-if="monitorDataset.cpuData.xData.length>0" custom-class="h-42" title="CPU使用率">
                      <GvaChart :line-title="'CPU使用率'" :x-line-data="monitorDataset.cpuData.xData"
                                :y-line-data="monitorDataset.cpuData.yData"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.memData.xData.length>0" custom-class="h-42" title="内存使用率">
                      <GvaChart :line-title="'内存使用率'" :x-line-data="monitorDataset.memData.xData"
                                :y-line-data="monitorDataset.memData.yData"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.diskData.xData.length>0" custom-class="h-42" title="磁盘使用率">
                      <GvaChart :line-title="'磁盘使用率'" :x-line-data="monitorDataset.diskData.xData"
                                :y-line-data="monitorDataset.diskData.yData"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.diskReadData.xData.length>0" custom-class="h-42" title="磁盘读速率">
                      <GvaChart :line-title="'磁盘读速率'" :x-line-data="monitorDataset.diskReadData.xData"
                                :y-line-data="monitorDataset.diskReadData.yData"
                                :chart-value="monitorDataset.diskReadData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.diskWriteData.xData.length>0" custom-class="h-42" title="磁盘写速率">
                      <GvaChart :line-title="'磁盘写速率'" :x-line-data="monitorDataset.diskWriteData.xData"
                                :y-line-data="monitorDataset.diskWriteData.yData"
                                :chart-value="monitorDataset.diskWriteData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.bpsRecvData.xData.length>0" custom-class="h-42" title="网络入带宽">
                      <GvaChart :line-title="'网络入带宽'" :x-line-data="monitorDataset.bpsRecvData.xData"
                                :y-line-data="monitorDataset.bpsRecvData.yData"
                                :chart-value="monitorDataset.bpsRecvData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.bpsSentData.xData.length>0" custom-class="h-42" title="网络出带宽">
                      <GvaChart :line-title="'网络出带宽'" :x-line-data="monitorDataset.bpsSentData.xData"
                                :y-line-data="monitorDataset.bpsSentData.yData"
                                :chart-value="monitorDataset.bpsSentData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.ppsRecvData.xData.length>0" custom-class="h-42" title="网络收包数">
                      <GvaChart :line-title="'网络收包数'" :x-line-data="monitorDataset.ppsRecvData.xData"
                                :y-line-data="monitorDataset.ppsRecvData.yData"
                                :chart-value="monitorDataset.ppsRecvData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.ppsSentData.xData.length>0" custom-class="h-42" title="网络发包数">
                      <GvaChart :line-title="'网络发包数'" :x-line-data="monitorDataset.ppsSentData.xData"
                                :y-line-data="monitorDataset.ppsSentData.yData"
                                :chart-value="monitorDataset.ppsSentData.dataValue"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.gpuData.xData.length>0" custom-class="h-42" title="GPU使用率">
                      <GvaChart :line-title="'GPU使用率'" :x-line-data="monitorDataset.gpuData.xData"
                                :y-line-data="monitorDataset.gpuData.yData"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.gpuMemData.xData.length>0" custom-class="h-42" title="GPU显存使用率">
                      <GvaChart :line-title="'GPU显存使用率'" :x-line-data="monitorDataset.gpuMemData.xData"
                                :y-line-data="monitorDataset.gpuMemData.yData"
                                :type="8"/>
                    </GvaCard>
                    <GvaCard v-if="monitorDataset.gpuTempData.xData.length>0" custom-class="h-42" title="GPU温度">
                      <GvaChart :line-title="'GPU温度'" :x-line-data="monitorDataset.gpuTempData.xData"
                                :y-line-data="monitorDataset.gpuTempData.yData"
                                :chart-value="monitorDataset.gpuTempData.dataValue"
                                :type="8"/>
                    </GvaCard>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="configData.status!='running'">
              <warning-bar :title="'必须处于开机状态下才能进行监控'"/>
            </div>
          </el-tab-pane>
          <el-tab-pane name="操作日志" label="操作日志">
            <div v-if="tabActiveName=='操作日志'">
              <Logs :log-data="logData"/>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="visible=false" type="info">关闭</el-button>
      </div>
    </el-drawer>
  </div>
</template>
<script>
import {getRenterResList} from "@/api/yunguan/resource/resInfo";
import {GetContainSSH} from "@/api/yunguan/cloudpods/container/container";
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import {
  CreateServerSnap,
  GetVirtualSSH,
  setHostRestart,
  setHostShutDown,
  setHostStart,
} from "@/api/yunguan/cloudpods/virtual/virtual";
import {monitorInstance} from "@/api/yunguan/cloudpods/instance/instance";
import {
  forecastAgent,
  forecastAgentCanInstall,
  getMonitorBpsRecvData,
  getMonitorBpsSentData,
  getMonitorCpuData,
  getMonitorDiskData,
  getMonitorDiskReadData,
  getMonitorDiskWriteData,
  getMonitorGpuData,
  getMonitorGpuMemData,
  getMonitorGpuTemperatureData,
  getMonitorMemData,
  getMonitorPpsRecvData,
  getMonitorPpsSentData,
  hostDetails,
  installAgent,
  setAgentCanInstall,
  updateName
} from "@/api/cloudpods/compute/servers";
import {getWebConsoleSSH} from "@/api/cloudpods/baseRes/host";
import Snapshot from "@/view/cloudpods/Storage/snapshot/snapshot.vue";
import Logs from "@/view/cloudpods/component/log.vue";
import {format} from "date-fns";
import WarningBar from "@/components/warningBar/warningBar.vue";
import _ from "lodash";
import {graphic} from "echarts";

export default {
  components: {WarningBar, Logs, Snapshot, GvaCard, GvaChart},
  data() {
    return {
      isCanInstallAgent: true,
      isSetAgent: false,
      agentTimer: null,
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      serverID: '',
      monitorDataset: {
        cpuData: {
          xData: [],
          yData: [],
        },
        memData: {
          xData: [],
          yData: [],
        },
        diskData: {
          xData: [],
          yData: [],
        },
        diskReadData: {
          xData: [],
          yData: [],
          dataValue: 'Mb/s',
        },
        diskWriteData: {
          xData: [],
          yData: [],
          dataValue: 'Mb/s',
        },
        bpsRecvData: {
          xData: [],
          yData: [],
          dataValue: 'b/s',
        },
        bpsSentData: {
          xData: [],
          yData: [],
          dataValue: 'b/s',
        },
        ppsRecvData: {
          xData: [],
          yData: [],
          dataValue: 'pps',
        },
        ppsSentData: {
          xData: [],
          yData: [],
          dataValue: 'pps',
        },
        gpuData: {
          xData: [],
          yData: [],
        },
        gpuMemData: {
          xData: [],
          yData: [],
        },
        gpuTempData: {
          xData: [],
          yData: [],
          dataValue: '℃',
        }
      },
      monitorData: {
        time: '1h',
        interval: '1m',
        id: '',
      },
      monitorConfigData: {
        ['1h']: ['1m', '5m'],
        ['6h']: ['3m', '5m', '10m', '15m',],
        ['12h']: ['5m', '10m', '20m', '30m'],
        ['24h']: ['10m', '20m', '30m', '60m'],
        ['168h']: ['30m', '1h'],
        ['720h']: ['6h', '24h'],
      },
      tabActiveName: "详情",
      visible: false,
      nameVisible: false,
      orderInfo: {},
      resInfo: {},
      hostStatus: {
        ['running']: '开机',
        ['disk']: '分配磁盘中',
        ['stopping']: '关机中',
        ['ready']: '关机',
        ['suspend']: '挂起',
        ['suspending']: '挂起中',
        ['resuming']: '恢复中',
        ['starting']: '启动中'
      },
      dialogVisible: false,
      searchInfo: {
        page: 1,
        pageSize: 10,
        type: 1,
      },
      //CPU利用率
      cpuUse: [{value: 13, name: '利用率'},],
      //磁盘使用率
      diskUse: [
        {value: 2, name: '已使用'},
        {value: 5, name: '未使用'}
      ],
      //内存占用率
      memoryUse: [
        {value: 2, name: '已使用'},
        {value: 5, name: '未使用'}
      ],
      activeName: 'baremetal',
      total: 0,
      resourceType: ["裸金属", "云主机", "容器", "存储"],
      dataList: [],
      openDialog: false,
      sshData: {
        id: '',
        port: 22,
        username: '',
        password: '',
        consoleType: 'vnc',
      },
      rules: {
        username: [
          {required: true, message: '请输入账号', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'}
        ],
        port: [
          {required: true, message: '请输入端口', trigger: 'blur'}
        ],
        name: [
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]*$/, // 正则表达式，匹配以字母开头并仅包含字母和数字的字符串
            message: '名称必须以字母开头，且只能包含英文字符和数字',
            trigger: 'blur'
          }
        ],
      },
      tempSSH: {
        username: '',
        password: '',
      },
      snapVisible: false,
      snapData: {
        id: '',
        name: '',
      },
      hostBaseData: {},
      getDetails: {
        id: '',
        serverType: 'vm',
      },
      configData: {},
      isInstallAgent: {
        status: false,
        text: '未安装',
      },
      chartData: {
        type: 'line',
        smooth: true,
        showSymbol: false,
        symbol: 'circle',
        symbolSize: 6,
        name: '',
        data: [],
        lineStyle: {
          normal: {
            shadowBlur: 3,
            shadowColor: 'rgba(0, 0, 0, .12)',
            shadowOffsetX: 4,
            shadowOffsetY: 4,
          },
        },
        emphasis: {
          focus: 'series',
          itemStyle: {
            borderColor: 'blue',
            borderWidth: 2,
          }, areaStyle: {
            color: new graphic.LinearGradient(0, 0, 0, 1, [
              {offset: 0, color: 'rgba(70, 130, 180, 0.7)'},
              {offset: 1, color: 'rgba(70, 130, 180, 0)'}
            ]),
          }
        },
        areaStyle: {
          color: new graphic.LinearGradient(0, 0, 0, 1, [
            {offset: 0, color: 'rgba(70, 130, 180, 0.5)'},
            {offset: 1, color: 'rgba(70, 130, 180, 0)'}
          ]),
        },
      },
    }
  },
  created() {
    this.GetRenterInfoList()
  },
  methods: {
    OpenSSHDialog() {
      this.isSetAgent = true;
      this.sshData.id = this.configData.id;
      this.dialogVisible = true;
    },
    HandleForecastAgentCanInstall() {
      forecastAgentCanInstall({id: this.configData.id}).then(res => {
        if (res.code == 0) {
          this.isCanInstallAgent = res.data;
        }
      })
    },
    //安装agent
    HandleInstallAgent() {
      installAgent({id: this.configData.id}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '正在安装中，请稍等'
          })
          this.startAgentInterval()
        }
      })
    },
    //判断是否安装agent
    GetHostIsInstallAgent() {
      forecastAgent({id: this.configData.id}).then(res => {
        if (res.code == 0) {
          if (res.data.length == 0) {
            this.isInstallAgent.status = false;
            this.isInstallAgent.text = '未安装agent';
            this.HandleForecastAgentCanInstall();
          } else {
            if (res.data[0].status == 'succeed') {
              this.isInstallAgent.status = true;
              this.isInstallAgent.text = '已安装';
              this.stopAgentInterval()
            } else if (res.data[0].status == 'applying') {
              this.isInstallAgent.status = false;
              this.isInstallAgent.text = '正在安装agent中';
            } else {
              this.isInstallAgent.status = false;
              this.isInstallAgent.text = '安装失败';
              this.stopAgentInterval()
              this.HandleForecastAgentCanInstall();
            }
          }
        }
      })
    },
    //开启定时器
    startAgentInterval() {
      if (this.agentTimer == null) {
        this.agentTimer = setInterval(() => {
          this.GetHostIsInstallAgent();
        }, 5000);
      } else {
      }
    },
    //关闭定时器
    stopAgentInterval() {
      if (this.agentTimer) {
        clearInterval(this.agentTimer);
        this.agentTimer = null;
      }
    },
    OpenDrawer(val) {
      this.visible = true
      this.getDetails.id = val.resourceID;;
      this.serverID = val.resourceID;
      this.logData.id = val.resourceID;
      this.monitorData.id = this.serverID;
      this.GetHostDetail();
      if (val.status == 'running') this.GetHostIsInstallAgent();
    },
    GetHostDetail() {
      if (this.activeName == 'virtual') this.getDetails.serverType = "vm";
      else if (this.activeName == 'baremetal') this.getDetails.serverType = "baremetal";
      hostDetails(this.getDetails).then(res => {
        if (res.code == 0) {
          this.configData = res.data;
          if (this.configData.metadata.os_distribution != undefined) {
            this.visible = true
          }
        }
      })
    },
    getGpuInfo(devices) {
      const modelCount = {};
      devices.forEach(device => {
        const model = device.model || '未知型号';
        modelCount[model] = (modelCount[model] || 0) + 1;
      });

      return Object.keys(modelCount)
          .map(model => `${model}*${modelCount[model]}卡`)
          .join('，');
    },
    GetDiskData(val, type) {
      if (val == undefined) return '';
      let num = 0;
      let stoType;

      for (let i = 0; i < val.length; i++) {
        if (val[i].disk_type == type) {
          num += val[i].size;
          stoType = val[i].storage_type;
        }
      }

      // 根据存储类型返回不同的描述
      let storageDesc = '';
      if (stoType === 'local') {
        storageDesc = '(本地硬盘)';
      } else if (stoType === 'nfs') {
        storageDesc = '(云硬盘-NFS)';
      } else if (stoType === 'rbd') {
        storageDesc = '(云硬盘-Ceph)';
      }

      // 使用 Math.round 确保四舍五入取整
      return Math.round(num / 1024) + 'GB' + storageDesc;
    },
    //刷新按钮
    handleRefresh() {
      this.GetMonitorData(this.monitorData)
    },
    //下拉框变动
    handleSelectChange() {
      this.GetMonitorData(this.monitorData)
    },
    //radio按钮监控时间变化
    handleRadioChange() {
      this.monitorData.interval = this.monitorConfigData[this.monitorData.time][0]
      this.GetMonitorData(this.monitorData)
    },
    clearData() {
      this.monitorDataset.cpuData.yData = []
      this.monitorDataset.cpuData.xData = []
      this.monitorDataset.memData.yData = []
      this.monitorDataset.memData.xData = []
      this.monitorDataset.diskData.yData = []
      this.monitorDataset.diskData.xData = []
      this.monitorDataset.diskReadData.yData = []
      this.monitorDataset.diskReadData.xData = []
      this.monitorDataset.diskWriteData.yData = []
      this.monitorDataset.diskWriteData.xData = []
      this.monitorDataset.bpsRecvData.yData = []
      this.monitorDataset.bpsRecvData.xData = []
      this.monitorDataset.bpsSentData.yData = []
      this.monitorDataset.bpsSentData.xData = []
      this.monitorDataset.ppsRecvData.yData = []
      this.monitorDataset.ppsRecvData.xData = []
      this.monitorDataset.ppsSentData.yData = []
      this.monitorDataset.ppsSentData.xData = []
      this.monitorDataset.gpuData.yData = []
      this.monitorDataset.gpuData.xData = []
      this.monitorDataset.gpuMemData.yData = []
      this.monitorDataset.gpuMemData.xData = []
      this.monitorDataset.gpuTempData.yData = []
      this.monitorDataset.gpuTempData.xData = []
    },

    //获取最大值
    GetMaxData(val) {
      let maxData = 0;
      for (let j = 0; j < val.length; j++) {
        let memData = val[j].points;
        for (let i = 0; i < memData.length; i++) {
          let tempData = memData[i]
          if (tempData[0] >= maxData) {
            maxData = tempData[0];
          }
        }
      }

      if (maxData < 1024) return 1;
      else if (maxData < 1024 * 1024) return 1024;
      else if (maxData < 1024 * 1024 * 1024) return (1024 * 1024);
      return (1024 * 1024 * 1024);
    },
    //获取监控数据
    GetMonitorData(obj) {
      this.clearData();
      getMonitorCpuData(obj).then(res => {
        if (res.code == 0 && res.data.cpu.series_total > 0) {
          for (let j = 0; j < res.data.cpu.series.length; j++) {
            let memData = res.data.cpu.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.cpu.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.cpuData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.cpuData.yData.push(tempChartData);
          }
        }
      })
      getMonitorMemData(obj).then(res => {
        if (res.code == 0 && res.data.mem.series_total > 0) {
          for (let j = 0; j < res.data.mem.series.length; j++) {
            let memData = res.data.mem.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.mem.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.memData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.memData.yData.push(tempChartData);
          }
        }
      })
      getMonitorDiskData(obj).then(res => {
        if (res.code == 0 && res.data.disk.series_total > 0) {
          for (let j = 0; j < res.data.disk.series.length; j++) {
            let memData = res.data.disk.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.disk.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.diskData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.diskData.yData.push(tempChartData);
          }
        }
      })
      getMonitorDiskReadData(obj).then(res => {
        if (res.code == 0 && res.data.diskRead.series_total > 0) {
          let maxData = this.GetMaxData(res.data.diskRead.series)
          if (maxData < 1024) this.monitorDataset.diskReadData.dataValue = 'b/s'
          else if (maxData < 1024 * 1024) this.monitorDataset.diskReadData.dataValue = 'Kb/s'
          else if (maxData < 1024 * 1024 * 1024) this.monitorDataset.diskReadData.dataValue = 'Mb/s'
          else this.monitorDataset.diskReadData.dataValue = 'Gb/s'
          for (let j = 0; j < res.data.diskRead.series.length; j++) {
            let memData = res.data.diskRead.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.diskRead.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push((tempData[0] / maxData).toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.diskReadData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.diskReadData.yData.push(tempChartData);
          }
        }
      })
      getMonitorDiskWriteData(obj).then(res => {
        if (res.code == 0 && res.data.diskWrite.series_total > 0) {
          let maxData = this.GetMaxData(res.data.diskWrite.series)
          if (maxData < 1024) this.monitorDataset.diskWriteData.dataValue = 'b/s'
          else if (maxData < 1024 * 1024) this.monitorDataset.diskWriteData.dataValue = 'Kb/s'
          else if (maxData < 1024 * 1024 * 1024) this.monitorDataset.diskWriteData.dataValue = 'Mb/s'
          else this.monitorDataset.diskWriteData.dataValue = 'Gb/s'
          for (let j = 0; j < res.data.diskWrite.series.length; j++) {
            let memData = res.data.diskWrite.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.diskWrite.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push((tempData[0] / maxData).toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.diskWriteData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.diskWriteData.yData.push(tempChartData);
          }
        }
      })
      getMonitorBpsRecvData(obj).then(res => {
        if (res.code == 0 && res.data.bpsRecv.series_total > 0) {
          let maxData = this.GetMaxData(res.data.bpsRecv.series)
          if (maxData < 1024) this.monitorDataset.bpsRecvData.dataValue = 'b/s'
          else if (maxData < 1024 * 1024) this.monitorDataset.bpsRecvData.dataValue = 'Kb/s'
          else if (maxData < 1024 * 1024 * 1024) this.monitorDataset.bpsRecvData.dataValue = 'Mb/s'
          else this.monitorDataset.bpsRecvData.dataValue = 'Gb/s'
          for (let j = 0; j < res.data.bpsRecv.series.length; j++) {
            let memData = res.data.bpsRecv.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.bpsRecv.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push((tempData[0] / maxData).toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.bpsRecvData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.bpsRecvData.yData.push(tempChartData);
          }
        }
      })
      getMonitorBpsSentData(obj).then(res => {
        if (res.code == 0 && res.data.bpsSent.series_total > 0) {
          let maxData = this.GetMaxData(res.data.bpsSent.series)
          if (maxData < 1024) this.monitorDataset.bpsSentData.dataValue = 'b/s'
          else if (maxData < 1024 * 1024) this.monitorDataset.bpsSentData.dataValue = 'Kb/s'
          else if (maxData < 1024 * 1024 * 1024) this.monitorDataset.bpsSentData.dataValue = 'Mb/s'
          else this.monitorDataset.bpsSentData.dataValue = 'Gb/s'
          for (let j = 0; j < res.data.bpsSent.series.length; j++) {
            let memData = res.data.bpsSent.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.bpsSent.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push((tempData[0] / maxData).toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.bpsSentData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.bpsSentData.yData.push(tempChartData);
          }
        }
      })
      getMonitorPpsRecvData(obj).then(res => {
        if (res.code == 0 && res.data.ppsRecv.series_total > 0) {
          for (let j = 0; j < res.data.ppsRecv.series.length; j++) {
            let memData = res.data.ppsRecv.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.ppsRecv.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.ppsRecvData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.ppsRecvData.yData.push(tempChartData);
          }
        }
      })
      getMonitorPpsSentData(obj).then(res => {
        if (res.code == 0 && res.data.ppsSent.series_total > 0) {
          for (let j = 0; j < res.data.ppsSent.series.length; j++) {
            let memData = res.data.ppsSent.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.ppsSent.series[j].name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.ppsSentData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.ppsSentData.yData.push(tempChartData);
          }
        }
      })
      getMonitorGpuData(obj).then(res => {
        if (res.code == 0 && res.data.utilization_gpu.series_total > 0) {
          for (let j = 0; j < res.data.utilization_gpu.series.length; j++) {
            let memData = res.data.utilization_gpu.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.utilization_gpu.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.gpuData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.gpuData.yData.push(tempChartData);
          }
        }
      })
      getMonitorGpuMemData(obj).then(res => {
        if (res.code == 0 && res.data.utilization_memory.series_total > 0) {
          for (let j = 0; j < res.data.utilization_memory.series.length; j++) {
            let memData = res.data.utilization_memory.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.utilization_memory.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.gpuMemData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.gpuMemData.yData.push(tempChartData);
          }
        }
      })
      getMonitorGpuTemperatureData(obj).then(res => {
        if (res.code == 0 && res.data.temperature_gpu.series_total > 0) {
          for (let j = 0; j < res.data.temperature_gpu.series.length; j++) {
            let memData = res.data.temperature_gpu.series[j].points;
            let tempChartData = _.cloneDeep(this.chartData)
            tempChartData.name = res.data.temperature_gpu.series[j].raw_name;
            for (let i = 0; i < memData.length; i++) {
              let tempData = memData[i]
              tempChartData.data.push(tempData[0].toFixed(2));
              if (j == 0) {
                const formattedDate = format(new Date(tempData[1]), 'MM-dd HH:mm');
                this.monitorDataset.gpuTempData.xData.push(formattedDate);
              }
            }
            this.monitorDataset.gpuTempData.yData.push(tempChartData);
          }
        }
      })
    },
    //时间间隔转化
    convertTimeInterval(interval) {
      const timeMap = {
        '1m': '1分钟',
        '3m': '3分钟',
        '5m': '5分钟',
        '10m': '10分钟',
        '15m': '15分钟',
        '20m': '20分钟',
        '30m': '30分钟',
        '60m': '1小时',
        '1h': '1小时',
        '6h': '6小时',
        '24h': '1天',
        // 可以根据需要添加更多映射
      };
      return timeMap[interval] || interval; // 如果没有找到匹配项，则返回原值
    },
    //tab按钮切换
    HandleTabChange() {
      switch (this.tabActiveName) {
        case "详情":
          this.GetHostDetail();
          break;
        case "监控":
          if (this.isInstallAgent.status) {
            this.monitorData.id = this.configData.id
            this.GetMonitorData(this.monitorData)
          }
          break;
      }
    },
    DisableType(val, type) {
      let isDisable = false
      if (type == "开机") {
        switch (val.status) {
          case "running":
          case "stopping":
          case "suspending":
          case "resuming":
          case "starting":
          case "suspend":
            isDisable = true;
            break;
          case "ready":
            isDisable = false;
            break;
        }
      } else if (type == "关机" || type == "重启" || type == "挂起") {
        switch (val.status) {
          case "ready":
          case "stopping":
          case "suspending":
          case "resuming":
          case "starting":
          case "suspend":
            isDisable = true;
            break;
          case "running":
            isDisable = false;
            break;
        }
      } else if (type == "恢复") {
        if (val.status == "suspend") isDisable = false
        else isDisable = true
      }
      return isDisable
    },
    //VNC远程连接
    HandleRemoteConn(val) {
      getWebConsoleSSH({id: val.host_id}).then(res => {
        if (res.code == 0) {
          const data = JSON.parse(res.data)
          // const params = new URLSearchParams(data).toString();
          // const routeData = this.$router.resolve({
          //   name: 'virTerminal',
          //   query: { id: data },
          // });
          window.open(`${data.access_url}&session_id=${data.session}`, '_blank');
        }
      })
    },
    RemoteConnect(val, data) {
      this.sshData.id = val.resourceID
      this.sshData.port = 22
      this.sshData.consoleType = data
      if (data == "vnc") {
        this.HandleVirtualRemoteConnect()
      } else {
        this.dialogVisible = true
        this.tempSSH = val
      }
    },
    //关闭对话框
    CloseDialog() {
      this.dialogVisible = false;
      this.snapVisible = false;
      this.nameVisible = false;
    },
    //打开创建快照对话框
    OpenCreateSnapDialog(val) {
      this.snapData.id = val.resourceID
      this.snapVisible = true
    },
    //打开修改名称对话框
    OpenNameDialog(val) {
      this.snapData.id = val.resourceID
      this.snapData.name = val.name
      this.nameVisible = true
    },
    //创建快照确定
    handleSnapConfirm() {
      this.$refs.snapData.validate(valid => {
        if (valid) {
          CreateServerSnap(this.snapData).then(res => {
            if (res.code == 0) {
              this.snapVisible = false;
              this.$message({
                message: '创建成功',
                type: 'success'
              })
            }
          })
        }
      })
    },
    //修改名称确定
    handleNameConfirm() {
      this.$refs.snapData.validate(valid => {
        if (valid) {
          updateName(this.snapData).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '修改成功'
              })
              this.nameVisible = false;
              this.GetRenterInfoList()
            }
          })
        }
      })
    },
    //ssh连接确定
    handleConnectConfirm() {
      if (this.isSetAgent) {
        setAgentCanInstall(this.sshData).then(res => {
          if (res.code == 0) {
            this.dialogVisible = false
            this.HandleForecastAgentCanInstall();
          }
        })
        return;
      }
      this.$refs.sshData.validate(valid => {
        if (valid) {
          if (this.tempSSH.username == this.sshData.username && this.tempSSH.password == this.sshData.password) {
            this.HandleVirtualRemoteConnect()
          } else {
            this.$message({
              message: '用户名或密码错误',
              type: 'error',
            })
          }
        }
      })
    },
    //关机
    HostShutDown(val) {
      this.$confirm('此操作将关机该云主机, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let obj = {}
        obj.id = val.resourceID
        setHostShutDown(obj).then(response => {
          if (response.code == 0) {
            this.GetRenterInfoList()
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消关机'
        });
      });
    },
    //开机
    HostBoot(val) {
      let obj = {}
      obj.id = val.resourceID
      setHostStart(obj).then(response => {
        if (response.code == 0) {
          this.GetRenterInfoList()
        }
      })
    },
    //重启
    HostReboot(val) {
      this.$confirm('强制重启模式会导致服务器实例当前未保存的数据丢失, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let obj = {}
        obj.id = val.resourceID
        setHostRestart(obj).then(response => {
          if (response.code == 0) {
            this.GetRenterInfoList()
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消重启'
        });
      });
    },
    //tab切换
    handleTabsClick() {
      console.log(this.activeName)
      switch (this.activeName) {
        case "baremetal":
          this.searchInfo.type = 1
          break;
        case "virtual":
          this.searchInfo.type = 2
          break;
        case "container":
          this.searchInfo.type = 3
          break;
      }
      this.GetRenterInfoList()
    },
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date === '0001-01-01T00:00:00.000Z') {
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
    //容器远程连接
    HandleContainerRemoteConnect(val) {
      let obj = {}
      obj.id = val.resourceID
      GetContainSSH(obj).then(res => {
        if (res.code == 0) {
          const data = JSON.parse(res.data).connect_params
          // const params = new URLSearchParams(data).toString();
          const routeData = this.$router.resolve({
            name: 'terminal',
            query: {id: data},
          });
          window.open(routeData.href, '_blank');
        }
      })
    },
    //云主机远程连接
    HandleVirtualRemoteConnect() {
      GetVirtualSSH(this.sshData).then(res => {
        if (res.code == 0) {
          this.dialogVisible = false
          const data = JSON.parse(res.data)
          // const params = new URLSearchParams(data).toString();
          const routeData = this.$router.resolve({
            name: 'virTerminal',
            query: {id: data},
          });
          window.open(`${data.access_url}&session_id=${data.session}`, '_blank');
          // window.open(routeData.href, '_blank');
        }
      })
    },
    //监控
    handleMonitor(val) {
      let obj = {}
      obj.privateIp = val.privateIp
      monitorInstance(obj).then(res => {
        if (res.code == 0) {
          this.hostBaseData = res.data
          this.cpuUse[0].value = res.data.cpu.used_percent
          this.diskUse[0].value = Number(res.data.disk.used)
          this.diskUse[1].value = Number(res.data.disk.free)
          this.memoryUse[0].value = Number(res.data.mem.used)
          this.memoryUse[1].value = Number(res.data.mem.available)
          this.openDialog = true
        }
      })
    },
    handleEdit(val) {
      this.dataInfo = val
      this.dataInfo.specDesc = JSON.parse(val.specDesc)
      this.IsCreated = false
      this.openDialog = true
      this.GetAllConfigInfoList()
    },
    handleSearch() {
      this.GetRenterInfoList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetRenterInfoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetRenterInfoList()
    },
    // 获取分区列表
    GetRenterInfoList() {
      getRenterResList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
  }
}
</script>

<style scoped>
.brand-info {
  margin: 20px;
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
</style>