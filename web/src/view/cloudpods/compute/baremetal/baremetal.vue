<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo"
                 style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="名称">
              <el-input
                  class="filter-item"
                  placeholder="请输入"
                  style="width: 240px;"
                  v-model="searchInfo.name"
              />
            </el-form-item>
            <el-form-item label="状态">
              <el-select style="width: 240px;" v-model="searchInfo.status" placeholder="请选择" clearable>
                <el-option label="运行中" value="running"/>
                <el-option label="关机" value="ready"/>
                <el-option label="挂起" value="suspend"/>
              </el-select>
            </el-form-item>
            <el-form-item label="系统">
              <el-select style="width: 240px;" v-model="searchInfo.osDist" placeholder="请选择" clearable>
                <el-option label="CentOS" value="CentOS"/>
                <el-option label="Ubuntu" value="Ubuntu"/>
                <el-option label="Debian" value="Debian"/>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item>
            <el-button type="primary" @click="GetBaremetalList">
              查询
            </el-button>
            <el-button @click="searchInfo.name='';searchInfo.status='';searchInfo.osDist='';GetBaremetalList()"
                       type="info">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button @click="OpenCreateDrawer"
                   type="primary">
          新增
        </el-button>
        <el-tooltip :disabled="!deleteList.includes(false)" content="VPC不为空,请先删除其中的网络"
                    placement="top">
          <el-button :disabled="!ids.length||deleteList.includes(false)"
                     type="grey" @click="HandleDeleteDialog('','deleteIds')">
            删除
          </el-button>
        </el-tooltip>
      </div>
      <el-table v-loading="false" :cell-style="{ 'text-align': 'left' }" :data="hostDataList"
                @selection-change="handleSelectionChangeData"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
                style="width: 100%;font-size: 15px;">
        <el-table-column type="selection" width="55"/>
        <!--        <el-table-column label="租户名称">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.description }}</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="OpenDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <el-button type="text">
              {{ GetRenterInfo(scope.row.renterId) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <i v-if="scope.row.status == 'running'" class="status-dot"/>
            <i v-else-if="scope.row.status == 'sched_fail'" style="background-color: red" class="status-dot"/>
            <i v-else-if="scope.row.status == 'ready'" style="background-color: rgb(174,185,192)" class="status-dot"/>
            <el-icon v-else class="status-dot-other">
              <Loading/>
            </el-icon>
            <span>{{ hostStatusType[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP地址">
          <template #default="scope">
            <span>{{ scope.row.ips }}</span>
          </template>
        </el-table-column>
        <el-table-column label="初始账号">
          <template #default="scope">
            <el-popover placement="top" :width="popoverWidth" trigger="click">
              <template #reference>
                <el-button @click="HandleGetLoginInfo(scope.row)" type="text"
                           style="margin-right: 16px">查看
                </el-button>
              </template>
              <span>用户名:{{ impiData.username }}</span><br/>
              <span>密码:{{ impiData.password }}</span>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="系统">
          <template #default="scope">
            <span>{{
                scope.row.metadata.os_distribution != "" ? scope.row.metadata.os_distribution : scope.row.os_type
              }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="250px">
          <template #default="scope">
            <el-dropdown style="margin-top: 3px;"
                         :title="scope.row.status=='running'?'':'云主机必须在开机状态下才可进行此操作'">
              <el-button type="text" link>
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
                  <el-dropdown-item :disabled="scope.row.status!='running'" @click="RemoteConnect(scope.row,'ssh')">
                    SSH 连接
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-dropdown style="margin-left: 10px;margin-top: 3px;" >
              <el-button type="text" link>
                实例状态
                <el-icon class="el-icon--right">
                  <arrow-down/>
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :disabled="DisableType(scope.row, '关机')" @click="HostShutDown(scope.row)">关机
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="DisableType(scope.row, '开机')" @click="HostBoot(scope.row)">开机
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="DisableType(scope.row, '重启')" @click="HostReboot(scope.row)">重启
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.status != 'ready'"
                                    :title="scope.row.status == 'ready'?'':'裸金属必须在关机状态下才可进行此操作'"
                                    @click="OpenResetDrawer(scope.row)">重置密码
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.status!='ready'&&!scope.row.status.includes('fail')"
                                    :title="scope.row.status == 'ready'?'':'裸金属必须在关机状态下才可进行此操作'"
                                    @click="HandleDeleteDialog(scope.row,'delete')">删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="searchInfo.pageIndex" :page-size="searchInfo.pageSize"
                       :page-sizes="[10, 30, 50, 100]" :total="total" layout="total,  prev, pager, next, sizes,jumper"
                       @current-change="handleCurrentChange" @size-change="handleSizeChange"/>
      </div>
    </div>
    <!--云主机详情-->
    <el-drawer v-model="visible" :before-close="CloseDrawer" :title="drawerTitle" size="1200px">
      <div style="margin: 0 10px">
        <el-tabs v-model="tabActiveName" @tab-change="HandleTabChange">
          <el-tab-pane name="详情" label="详情">
            <div style="margin-top: 20px">
              <el-descriptions :column="1" label-width="100px">
                <!--                    <el-descriptions-item label="ID">{{ configData.id }}</el-descriptions-item>-->
                <el-descriptions-item label="名称">{{ configData.name }}</el-descriptions-item>
                <el-descriptions-item label="操作系统">
                  {{
                    configData.metadata.os_distribution != undefined ? configData.metadata.os_distribution : '-' + ' ' + configData.metadata.os_version != undefined ? configData.metadata.os_version : ''
                  }}
                </el-descriptions-item>
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
                  {{ configData.metadata.os_arch != undefined ? configData.metadata.os_arch : '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="CPU">{{ configData.vcpu_count }}核</el-descriptions-item>
                <el-descriptions-item label="内存">{{ FormatMem(configData.vmem_size) }}GB</el-descriptions-item>
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
                      <span>取值间隔（分钟）</span>
                      <el-select @change="handleSelectChange" style="width: 200px" v-model="monitorData.interval">
                        <el-option v-for="(item,index) in monitorConfigData[monitorData.time]"
                                   :label="convertTimeInterval(item)"
                                   :value="item"/>
                      </el-select>
                    </div>
                    <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
                      <el-radio-group v-model="monitorData.time" @change="handleRadioChange">
                        <el-radio-button type="key" value="1h">近1小时</el-radio-button>
                        <el-radio-button type="key" value="6h">近6小时</el-radio-button>
                        <el-radio-button type="key" value="12h">近12小时</el-radio-button>
                        <el-radio-button type="key" value="24h">近24小时</el-radio-button>
                        <!--                      <el-radio-button value="168h">近1周</el-radio-button>-->
                        <!--                      <el-radio-button value="720h">近1月</el-radio-button>-->
                      </el-radio-group>
                      <el-button type="key" @click="handleRefresh" icon="refresh" value="refresh"></el-button>
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
    </el-drawer>
    <!--    新增和编辑-->
    <el-drawer v-model="createVisible" :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false"
               size="1200px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{
            IsCreated ? '新增' : '编辑'
          }}裸金属</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>

      <el-form
          label-position="left"
          ref="serverDataInfo"
          :model="serverDataInfo"
          :rules="IsCreated ? rules : ''"
          label-width="90px">
        <el-form-item label="名称" prop="generate_name">
          <el-input style="width: 320px;" :disabled="true" v-model="serverDataInfo.generate_name"/>
        </el-form-item>
        <el-form-item label="镜像" prop="imageId">
          <el-select style="width: 320px;" v-model="serverDataInfo.imageId" placeholder="请选择镜像" @change="skuChange">
            <el-option v-for="(optionItem, optionIndex) in imageList" :key="optionIndex" :label="optionItem.name"
                       :value="optionItem.id"/>
          </el-select>
        </el-form-item>
        <el-form-item style="padding-left: 10px;" label-width="80" label="绑定租户">
          <el-switch
              v-model="isBindRenter"
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
        <el-form-item prop="renterID" v-if="isBindRenter" label="租户">
          <el-select style="width: 320px;" placeholder="请选择租户" @change="generateUniqueName(serverDataInfo.renterID)"
                     v-model="serverDataInfo.renterID">
            <el-option v-for="item in renterList"
                       :label="item.username"
                       :value="item.renterID"/>
          </el-select>
        </el-form-item>
        <el-form-item style="padding-left: 10px" label-width="80" label="VPC" prop="vpcId">
          <el-select
              style="width: 320px;"
              v-model="serverDataInfo.vpcId"
              placeholder="请选择vpc"
              @change="GetIpList();GetExitIpList();serverDataInfo.network='';serverDataInfo.exitIP='';">
            <el-option
                v-for="(optionItem, optionIndex) in vpcList"
                :key="optionIndex"
                :label="'名称:'+optionItem.name"
                :value="optionItem.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item style="padding-left: 10px" label-width="80" label="公网IP">
          <el-select style="width: 320px;" v-model="serverDataInfo.exitIP" placeholder="请选择公网IP">
            <el-option v-for="(optionItem, optionIndex) in exitIpList" :key="optionIndex"
                       :disabled="(optionItem.ports - optionItem.ports_used) == 0"
                       :label="optionItem.name + '(' + optionItem.guest_ip_start + '-' + optionItem.guest_ip_end + ')    ' + '可用：' + (optionItem.ports - optionItem.ports_used)"
                       :value="optionItem.id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="IP子网" prop="network">
          <el-select style="width: 320px;" v-model="serverDataInfo.network" placeholder="请选择IP子网" >
            <el-option v-for="(optionItem, optionIndex) in IpList" :key="optionIndex"
                       :disabled="(optionItem.ports - optionItem.ports_used) == 0"
                       :label="optionItem.name + '(' + optionItem.guest_ip_start + '-' + optionItem.guest_ip_end + ')    ' + '可用：' + (optionItem.ports - optionItem.ports_used)"
                       :value="optionItem.id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="物理机器" prop="host">
          <el-select style="width: 320px;" @change="handleChange" v-model="serverDataInfo.host">
            <el-option v-for="baseHost in baseHostList" :disabled="baseHost.server!=undefined"
                       :label="baseHost.name+(baseHost.server!=undefined?' 已分配':' 未分配')" :value="baseHost">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="磁盘配置" style="padding-left: 10px" label-width="80">
          <el-text v-show="baseHostDiskInfo === null" type="warning">
            请选择物理机
          </el-text>
          <!-- 渲染物理机的磁盘配置 -->
          <div class=" left-64" v-show="baseHostDiskInfo !== null" id="diskChart">
            <el-table border style="font-size: 10px" :data="hostDiskList">
              <el-table-column label="驱动" width="100px">
                <template #default="scope">
                  {{ scope.row.driver }}
                </template>
              </el-table-column>
              <el-table-column width="300px" label="型号">
                <template #default="scope">
                  {{ scope.row.model }}
                </template>
              </el-table-column>
              <el-table-column label="容量">
                <template #default="scope">
                  {{ GetImageSize(scope.row.size) }}
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-form-item>
        <el-form-item label="root密码" style="padding-left: 10px" label-width="80">
          <el-input style="width: 320px;" type="password" v-model="serverDataInfo.password"></el-input>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    重置密码-->
    <el-dialog v-model="passVisible" title="重置密码" width="40%">
      <el-form ref="resetData" :model="resetData" :rules="rules" label-width="120px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="resetData.username" placeholder="请输入用户名"/>
        </el-form-item>
        <el-form-item label="密码类型">
          <el-radio-group v-model="resetData.type">
            <el-radio-button :value="false">随机生成</el-radio-button>
            <el-radio-button :value="true">手工输入</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="resetData.type" label="密码" prop="password">
          <el-input v-model="resetData.password" placeholder="请输入密码"/>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="HandleResetConfirm">确 定</el-button>
        <el-button @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <!--    监控-->
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
    <!--    远程连接-->
    <el-dialog
        v-model="dialogVisible"
        title="请输入SSH账户密码"
        width="40%"
    >
      <el-form
          ref="sshData"
          :model="sshData"
          :rules="sshRules"
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
        <el-button @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <!--云主机详情-->
    <el-drawer v-model="renterVisible" :before-close="CloseDrawer" :title="'详情'" size="50%">
      <div style="margin: 0 10px" v-if="resInfo.renterID!=undefined">
        <el-descriptions :column="2" border title="基本信息">
          <el-descriptions-item label="所属租户">{{ GetRenterInfo(resInfo.renterID) }}</el-descriptions-item>
        </el-descriptions>
        <div v-if="orderInfo.code!=undefined">
          <el-descriptions :column="2" border>
            <el-descriptions-item :span="3" label="订单号">
              {{ orderInfo.code }}
            </el-descriptions-item>
            <el-descriptions-item label="产品名称">
              {{ orderInfo.productName }}
            </el-descriptions-item>
            <el-descriptions-item label="购买数量">
              {{ orderInfo.productNum }}
            </el-descriptions-item>
            <el-descriptions-item label="购买时长">
              {{ orderInfo.duration }}
              【{{
                orderInfo.billingType === 1 ? '卡时' :
                    orderInfo.billingType === 2 ? '日' :
                        orderInfo.billingType === 4 ? '周' :
                            orderInfo.billingType === 8 ? '月' :
                                orderInfo.billingType === 32 ? '年' : '无限制'
              }}】
            </el-descriptions-item>
            <el-descriptions-item
                :rowspan="1"
                label="已选规格"
            >
              <el-descriptions size="small" border>
                <el-descriptions-item
                    v-for="(item, index) in orderInfo.specJson"
                    :key="index"
                    :label="item.name"
                    :span="orderInfo.specJson.length"
                    label-align="right"
                    align="left"
                >
                  {{ item.value[item.valueIndex != undefined ? item.valueIndex : 0] }}
                </el-descriptions-item>
              </el-descriptions>
            </el-descriptions-item>
          </el-descriptions>
        </div>
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
  addHost,
  createBakemeats,
  deleteHost,
  deleteHostByIds,
  forecastAgent,
  forecastAgentCanInstall,
  getBaremetalList,
  getLoginInfo,
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
  installAgent,
  resetPassword,
  setAgentCanInstall,
  setHostDelete,
  setHostRestart,
  setHostShutDown,
  setHostStart,
} from "@/api/cloudpods/compute/servers";
import Logs from "@/view/cloudpods/component/log.vue"
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import {getNetIpList} from "@/api/cloudpods/network/ip";
import {getImageList} from "@/api/cloudpods/image/image";
import {getBaseBareHostList, getWebConsoleSSH} from "@/api/cloudpods/baseRes/host";
import {monitorInstance} from "@/api/yunguan/cloudpods/instance/instance";
import {
  addRenterRes,
  deleteRenterRes,
  getRenterList,
  queryRenterResInfo,
  updateRenterRes
} from "@/api/yunguan/run/renter";
import {GetVirtualSSH} from "@/api/yunguan/cloudpods/virtual/virtual";
import statusType from "@/locales/zh-CN.json";
import {getVpcList} from "@/api/cloudpods/network/vpc";
import WarningBar from "@/components/warningBar/warningBar.vue";
import {format} from "date-fns";
import _ from "lodash";
import {graphic} from "echarts";

export default {
  name: 'ComplexTable',
  components: {WarningBar, GvaCard, GvaChart, Logs},
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
      deleteType: '',
      ids: [],
      deleteList: [],
      isCanInstallAgent: true,
      isSetAgent: false,
      agentTimer: null,
      isInstallAgent: {
        status: false,
        text: '未安装',
      },
      isUpdateRenterRes: false,
      modRenterRes: {
        id: '',
        resourceID: '',
        startTime: '',
      },
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      timer: null,
      orderInfo: {},
      impiData: {},
      popoverWidth: 180,
      tabActiveName: "详情",
      searchInfo: {
        description: '裸金属',
        pageIndex: 1,
        pageSize: 10,
      },
      total: 0,
      visible: false,
      createVisible: false,
      IsCreated: true,
      drawerTitle: '',
      hostDataList: [],
      configData: {},
      baseHostDiskInfo: null,
      hostDiskList: [],
      baseHostList: [], //物理机器列表
      webUrl: '',
      hostStatusType: statusType.status.server,
      serverID: '',
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
        vpcId: '',
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
        exit: false,
      },
      IpList: [],
      exitIpList: [],
      rules: {
        generate_name: [
          {required: true, message: '请输入云主机名称', trigger: 'blur'}
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
        imageId: [
          {required: true, message: '请选择镜像', trigger: 'blur'}
        ],
        network: [
          {required: true, message: '请选择IP子网', trigger: 'blur'}
        ],
        host: [
          {required: true, message: '请选择物理机', trigger: 'blur'}
        ],
        name: [
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
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
      },
      skuList: [],
      imageList: [],
      serverList: [],
      dataDiskList: [],
      passVisible: false,
      sshData: {
        id: '',
        port: 22,
        username: '',
        password: '',
        consoleType: 'vnc',
      },
      sshRules: {
        username: [
          {required: true, message: '请输入账号', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'}
        ],
        port: [
          {required: true, message: '请输入端口', trigger: 'blur'}
        ],
        renterID: [
          {required: true, message: '请选择租户', trigger: 'blur'}
        ],
      },
      dialogVisible: false,
      resetData: {
        id: '',
        username: '',
        password: '',
        type: false,
      },
      openDialog: false,
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
      hostBaseData: {},
      resInfo: {},
      specDetail: [],
      renterList: [],
      renterVisible: false,
      isBindRenter: false,
      vpcList: [],
      addRenterResInfo: {
        renterID: '',
        type: 1,
        sshHost: '',
        sshUser: '',
        sshPwd: '',
        resourceID: '',
        privateIp: '',
        aiList: '[]',
      },
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
    this.GetVpcList()
    this.GetAllRenterList()
    this.GetImageList()
    this.GetBaremetalList()
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
      switch (type) {
        case 'delete':
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该裸金属，是否继续?';
          break;
      }
      this.deleteVisible = true;
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.deleteList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        if (val[i].status != 'ready' && !val[i].status.includes('fail')) {
          this.deleteList.push(false)
        }
      }
    },
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
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetBaremetalList();
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
    //打开绑定租户对话框
    GetResourceInfo(val) {
      this.resInfo = {};
      this.orderInfo = {};
      this.specDetail = []
      queryRenterResInfo({resourceID: val.id}).then(res => {
        if (res.code == 0) {
          this.resInfo = res.data
          if (res.data.description != "") {
            this.orderInfo = JSON.parse(res.data.description)
            this.orderInfo.specJson = JSON.parse(this.orderInfo.specJson)
            this.orderInfo.storage = JSON.parse(this.orderInfo.storage)
            this.orderInfo.specJson.push({name: '存储类型', value: [this.orderInfo.storage[0].name], valueIndex: 0})
            this.orderInfo.specJson.push({name: '操作系统', value: [this.orderInfo.os], valueIndex: 0})
          }
        }
      })
      this.renterVisible = true;
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
    //获取根据资源ID获取租户信息
    GetRenterInfo(val) {
      // 遍历租户资源列表
      for (const renter of this.renterList) {
        if (renter.renterID === val) {
          return renter.username;
        }
      }
      // 如果没有找到匹配的租户，返回“无”
      return "无";
    },
    //监控
    handleMonitor(val) {
      let obj = {}
      obj.privateIp = val.ips
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
    // 切换物理机选择处理
    handleChange() {
      console.log(this.serverDataInfo.host)
      this.baseHostDiskInfo = this.serverDataInfo.host.spec.disk
      this.hostDiskList = this.serverDataInfo.host.spec.storage_info
      console.log(this.hostDiskList)
      // this.drawDiskChart()
    },
    //打开重置密码对话框
    OpenResetDrawer(val) {
      this.resetData = {}
      this.passVisible = true;
      this.resetData.id = val.id;
      this.resetData.username = val.metadata.login_account;
      this.resetData.type = false;
    },
    //重置密码确认
    HandleResetConfirm() {
      this.$refs.resetData.validate((valid) => {
        if (valid) {
          resetPassword(this.resetData).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '重置成功',
              })
              this.passVisible = false;
              this.GetBaremetalList()
            }
          })
        }
      })
    },
    //获取登录信息
    HandleGetLoginInfo(val) {
      getLoginInfo({id: val.id}).then(res => {
        if (res.code == 0) {
          this.impiData = JSON.parse(res.data)
          const usernameLength = this.impiData.username.length;
          const passwordLength = this.impiData.password.length;

          // 根据用户名和密码的长度计算宽度
          const baseWidth = 14; // 基础宽度
          const additionalWidthPerChar = 10; // 每个字符增加的宽度

          const totalLength = usernameLength + passwordLength;
          const dynamicWidth = baseWidth + (totalLength * additionalWidthPerChar);
          this.popoverWidth = dynamicWidth;
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetBaremetalList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetBaremetalList()
    },
    //删除
    HandleDelete(val) {
      setHostDelete({id: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetBaremetalList(this.searchInfo)
        }
      })
    },
    //批量删除
    HandleDeleteByIds(val) {
      deleteHostByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetBaremetalList(this.searchInfo)
        }
      })
    },
    OpenDrawer(val) {
      this.visible = true
      this.drawerTitle = val.name
      this.configData = val
      this.serverID = val.id
      this.logData.id = val.id;
      if (this.configData.status == 'running') this.GetHostIsInstallAgent();
    },
    OpenCreateDrawer() {
      this.generateUniqueName(0)
      this.serverDataInfo.host = null;
      this.serverDataInfo.vpcId = 'default';
      this.GetIpList();
      this.GetExitIpList();
      this.baseHostDiskInfo = null;
      // 请求物理机信息
      getBaseBareHostList({
        pageIndex: 1,
        pageSize: 10000,
      }).then((resp) => {
        this.baseHostList = resp.data.data
      })
      this.createVisible = true;
    },
    //tab按钮切换
    HandleTabChange() {
      switch (this.tabActiveName) {
        case "详情":
          break;
        case "监控":
          if (this.isInstallAgent.status) {
            this.monitorData.id = this.configData.id
            this.GetMonitorData(this.monitorData)
          }
          break;
      }
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
    // 获取镜像列表
    GetImageList() {
      let obj = {}
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
          getLoginInfo({id: this.sshData.id}).then(res => {
            if (res.code == 0) {
              this.impiData = JSON.parse(res.data)
              if (this.impiData.username == this.sshData.username && this.impiData.password == this.sshData.password) {
                this.HandleVirtualRemoteConnect()
              } else {
                this.$message({
                  message: '用户名或密码错误',
                  type: 'error',
                })
              }
            }
          })

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
          // const routeData = this.$router.resolve({
          //   name: 'virTerminal',
          //   query: {id: data},
          // });
          window.open(`${data.access_url}&session_id=${data.session}`, '_blank');
          // window.open(routeData.href, '_blank');
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
    ParserBaremetalDisk(bastHost) {
      //解析格式BaremetalDiskConfigs
      let range = []
      console.log(bastHost)
      bastHost.spec.storage_info.forEach((item) => {
        range.push(item.slot)
      })
      let diskConfig = [{
        conf: "raid0",
        driver: bastHost.spec.driver,
        count: bastHost.spec.storage_info.length,
        type: bastHost.storage_type,
        adapter: bastHost.spec.storage_info[0].adapter,
        range: range
      }]
      console.log(diskConfig)
      return diskConfig
    },
    //提交对话框
    SubmitDrawer() {
      this.$refs.serverDataInfo.validate(valid => {
        if (valid) {
          this.serverDataInfo.vmem_size = this.serverDataInfo.host.mem_size
          this.serverDataInfo.vcpu_count = this.serverDataInfo.host.spec.cpu
          this.serverDataInfo.disks = [
            {
              size: this.serverDataInfo.host.storage_size - 30 * 1024,
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
          createBakemeats(this.serverDataInfo).then((resp) => {
            if (resp.code === 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              if (this.isBindRenter) {
                this.CreateRenterRes(JSON.parse(resp.data))
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

    handleBindRenterChange(value) {
      this.serverDataInfo.renterID = '';
      this.generateUniqueName(0);
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
    //创建租户资源
    CreateRenterRes(val) {
      this.addRenterResInfo.resourceID = val.id;
      this.addRenterResInfo.renterID = this.serverDataInfo.renterID;
      addRenterRes(this.addRenterResInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '绑定租户成功',
          })
          this.isUpdateRenterRes = true;
          this.modRenterRes.id = res.data.ID;
          this.modRenterRes.resourceID = val.id;
        }
      })
    },
    //获取主机列表
    GetBaremetalList() {
      getBaremetalList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData = response.data.data;
          if (JSON.stringify(tempData) !== JSON.stringify(this.hostDataList)) {
            this.hostDataList = response.data.data
            this.total = response.data.total
            this.startInterval()
            if (this.isUpdateRenterRes) {
              const tempData = this.hostDataList.filter(item => item.id == this.modRenterRes.resourceID);
              if (tempData.length > 0) {
                if (tempData[0].status.includes('fail')) {
                  deleteRenterRes(this.modRenterRes).then(res => {
                    if (res.code == 0) {
                      this.isUpdateRenterRes = false;
                    }
                  })
                }
                if (tempData[0].status == "running") {
                  const now = new Date();
                  const cTime = new Date(now.getTime() + 8 * 60 * 60 * 1000);
                  const isoString = cTime.toISOString().replace("Z", "+08:00");
                  this.modRenterRes.startTime = isoString;
                  updateRenterRes(this.modRenterRes).then(res => {
                    if (res.code == 0) {
                      this.isUpdateRenterRes = false;
                    }
                  })
                }
              }
            }
          } else {
            this.stopInterval()
          }
        }
      })
    },
    //创建云主机
    AddHost() {
      addHost(this.serverDataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            type: 'success',
            message: '创建成功',
          })
          this.GetBaremetalList()
          this.createVisible = false
        }
      })
      this.serverDataInfo = []
      this.serverDataInfo.disks = [],
          this.serverDataInfo.nets = [],
          this.serverDataInfo.disks.push({
            backend: 'local',
            disk_type: 'sys',
            index: 0,
            medium: 'rotate',
            size: 30,
            image_id: '',
          })
      this.imageId = ""
      this.dataDiskList = []
    },
    FormatMem(val) {
      return val / 1024
    },
    GetImageSize(val) {
      return (val / 1024).toFixed(0) + "GB";
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
    //远程连接
    RemoteConnect(val, data) {
      this.isSetAgent = false;
      this.sshData.id = val.id
      this.sshData.port = 22
      this.sshData.consoleType = data
      this.dialogVisible = true
      this.tempSSH = val
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
    //删除
    HostDelete(val) {
      let obj = {}
      obj.id = val.id
      deleteHost(obj).then(response => {
        if (response.code == 0) {
          this.GetBaremetalList()
          setTimeout(() => {
            this.GetBaremetalList()
          }, 3000)
          setTimeout(() => {
            this.GetBaremetalList()
          }, 5000)
        }
      })
    },
    //关机
    HostShutDown(val) {
      let obj = {}
      obj.id = val.id
      setHostShutDown(obj).then(response => {
        if (response.code == 0) {
          this.GetBaremetalList()
          setTimeout(() => {
            this.GetBaremetalList()
          }, 3000)
          setTimeout(() => {
            this.GetBaremetalList()
          }, 5000)
        }
      })
    },
    //开机
    HostBoot(val) {
      let obj = {}
      obj.id = val.id
      setHostStart(obj).then(response => {
        if (response.code == 0) {
          this.GetBaremetalList()
          setTimeout(() => {
            this.GetBaremetalList()
          }, 3000)
          setTimeout(() => {
            this.GetBaremetalList()
          }, 5000)
        }
      })
    },
    //重启
    HostReboot(val) {
      let obj = {}
      obj.id = val.id
      setHostRestart(obj).then(response => {
        if (response.code == 0) {
          this.GetBaremetalList()
          setTimeout(() => {
            this.GetBaremetalList()
          }, 5000)
          setTimeout(() => {
            this.GetBaremetalList()
          }, 15000)
          setTimeout(() => {
            this.GetBaremetalList()
          }, 30000)
        }
      })
    },

    //获取VPC列表
    GetVpcList() {
      let obj = {};
      obj.pageSize = 0;
      obj.id = 'default'
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
  },
  mounted() {
    this.startInterval()
  },
  unmounted() {
    this.stopInterval()
  },
}
</script>
<style lang="scss">
.el-popper[x-placement^=bottom] .popper__arrow {
  border: none;
}

.el-popper[x-placement^=bottom] .popper__arrow::after {
  border: none;
}

.el-drawer.rtl {
  overflow: scroll
}
</style>
<style scoped>
.brand-info {
  margin: 20px;
}

.boxer {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.el-table-column {
  text-align: left;
}

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50;
  /* Green color */
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
<style scoped>
.example-showcase .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
</style>
