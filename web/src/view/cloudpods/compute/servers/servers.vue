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
            <el-form-item label="系统">
              <el-select style="width: 240px;" v-model="searchInfo.osDist" placeholder="请选择" clearable>
                <el-option label="CentOS" value="CentOS"/>
                <el-option label="Ubuntu" value="Ubuntu"/>
                <el-option label="Debian" value="Debian"/>
              </el-select>
            </el-form-item>
            <el-form-item label="状态">
              <el-select style="width: 240px;" v-model="searchInfo.status" placeholder="请选择" clearable>
                <el-option label="运行中" value="running"/>
                <el-option label="关机" value="ready"/>
                <el-option label="挂起" value="suspend"/>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item>
            <el-button type="primary" @click="GetHostList">
              查询
            </el-button>
            <el-button type="info" @click="searchInfo.name='';searchInfo.status='';searchInfo.osDist='';GetHostList()">
              重置
            </el-button>
          </el-form-item>
        </el-form>

      </div>
    </div>
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button v-if="secGroupID==''&&scalingGroup==''" @click="OpenCreateDrawer"
                   type="primary">
          新增
        </el-button>
        <el-tooltip :disabled="!statusList.includes('running')" content="云主机必须在关机状态下才可进行此操作"
                    placement="top">
          <el-button v-if="secGroupID==''&&scalingGroup==''"
                     :disabled="!ids.length||statusList.includes('start_delete')||statusList.includes('running')"
                     type="grey" @click="HandleDeleteDialog('','deleteIds')">
            删除
          </el-button>
        </el-tooltip>
        <el-button v-if="secGroupID!=''" @click="OpenHostDrawer"
                   type="info">
          关联云主机
        </el-button>
      </div>

      <el-table
          :cell-style="{'text-align':'left'}"
          :data="hostDataList"
          @sort-change="handleSortChange"
          @selection-change="handleSelectionChangeData"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
<!--          style="width: 100%;font-size: 15px;"-->
      >
        <el-table-column type="selection" width="30"/>
        <el-table-column align="left" label="名称" min-width="100">
          <template #default="{ row }">
            <el-tooltip
                effect="light"
                :content="row.name"
            placement="top"
            :open-delay="300"
            >

            <el-button
                type="text"
                @click="OpenDrawer(row)"
                class="name-button"
            >
            {{ row.name }}
            </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column :show-overflow-tooltip='true' label="描述">
          <template #default="scope">
            {{ scope.row.description != undefined ? scope.row.description : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <el-button type="text">
              {{ GetRenterInfo(scope.row.renterId) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="140px">
          <template #default="scope">
            <i v-if="scope.row.status == 'running'" class="status-dot"/>
            <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
            <el-icon v-else-if="scope.row.status.includes('ing')||scope.row.status == 'schedule'"
                     class="status-dot-other">
              <Loading/>
            </el-icon>
            <i v-else style="background-color: rgb(174,185,192)" class="status-dot"/>
            <span>{{ hostStatusType[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP地址" sortable="custom" width="180px">
          <template #default="scope">
            <span>{{ scope.row.ips }}</span>
          </template>
        </el-table-column>
        <el-table-column label="初始账号">
          <template #default="scope">
            <el-popover placement="top" :width="popoverWidth" trigger="click">
              <template #reference>
                <el-button @click="HandleGetLoginInfo(scope.row)" type="text" style="margin-right: 16px">查看
                </el-button>
              </template>
              <span>用户名:{{ impiData.username }}</span><br/>
              <span>密码:{{ impiData.password }}</span>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="系统" width="100px">
          <template #default="scope">
            <span>{{ scope.row.metadata.os_distribution }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="secGroupID==''&&scalingGroup==''" label="操作" fixed="right" width="250px">
          <template #default="scope">
            <el-dropdown style="margin-top: 3px;">
              <el-button type="text" link>
                远程连接
                <el-icon class="el-icon--right">
                  <arrow-down/>
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :title="scope.row.status!='running'?'云主机必须在开机状态下才可进行此操作':''"
                                    :disabled="scope.row.status!='running'" @click="RemoteConnect(scope.row,'vnc')">
                    VNC 远程终端
                  </el-dropdown-item>
                  <el-dropdown-item :title="scope.row.status!='running'?'云主机必须在开机状态下才可进行此操作':''"
                                    :disabled="scope.row.status!='running'" @click="RemoteConnect(scope.row,'ssh')">
                    SSH 连接
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-dropdown style="margin-left: 10px;margin-top: 3px;">
              <el-button type="text" link>
                实例状态
                <el-icon class="el-icon--right">
                  <arrow-down/>
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>

                  <el-dropdown-item @click="HandleSetGpu(scope.row)">GPU管理</el-dropdown-item>


                  <el-dropdown-item type="text" @click="OpenCreateSnapDialog(scope.row)">创建快照
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.is_gpu" :title="scope.row.is_gpu?'请先卸载GPU后进行迁移':''"
                                    v-if="scope.row.status=='running'" type="text"
                                    @click="OpenMigrateDialog(scope.row,'hot')">热迁移
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.is_gpu" :title="scope.row.is_gpu?'请先卸载GPU后进行迁移':''"
                                    v-if="scope.row.status=='ready'" type="text"
                                    @click="OpenMigrateDialog(scope.row,'cold')">冷迁移
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.is_gpu" :title="scope.row.is_gpu?'请先卸载GPU后进行迁移':''"
                                    v-if="scope.row.status=='unknown'" type="text"
                                    @click="OpenMigrateDialog(scope.row,'down')">宕机迁移
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-dropdown style="margin-left: 10px;margin-top: 3px;">
              <el-button type="text" link>
                更多
                <el-icon class="el-icon--right">
                  <arrow-down/>
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="OpenResetDrawer(scope.row)">重置密码</el-dropdown-item>
                  <el-dropdown-item :disabled="scope.row.status!='ready'&&!scope.row.status.includes('fail')"
                                    :title="scope.row.status!='ready'&&!scope.row.status.includes('fail')?'云主机必须在关机状态下才可进行此操作':''"
                                    @click="HandleDeleteDialog(scope.row,'delete')">删除
                  </el-dropdown-item>
                  <el-dropdown>
                    <el-button style="margin-left: 15px;margin-top: 5px;margin-bottom: 5px;color: #000!important;"  link>
                      实例状态
                      <el-icon class="el-icon--right">
                        <arrow-down/>
                      </el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :disabled="scope.row.status!='running'||scope.row.status.includes('fail')"
                                          @click="HandleDeleteDialog(scope.row,'shut')">关机
                        </el-dropdown-item>
                        <el-dropdown-item :disabled="scope.row.status!='ready'||scope.row.status.includes('fail')"
                                          @click="HostBoot(scope.row)">开机
                        </el-dropdown-item>
                        <el-dropdown-item :disabled="DisableType(scope.row,'重启')||scope.row.status.includes('fail')"
                                          @click="HandleDeleteDialog(scope.row,'restart')">重启
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
        <el-table-column v-if="secGroupID!=''" label="操作">
          <template #default="scope">
            <el-tooltip :disabled="scope.row.secgroups.length!=1" content="该云主机只有一个安全组，不支持该操作"
                        placement="top">
              <el-button type="text" :disabled="scope.row.secgroups.length==1" @click="HandleDeleteDialog(scope.row,'revoke')">解绑
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column v-if="scalingGroup!=''" label="操作">
          <template #default="scope">
            <el-button type="text" @click="HandleDetachGroup(scope.row)">移除</el-button>
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
    <!--云主机详情-->
    <el-drawer v-model="visible" :show-close="false" size="1200">

      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{
            drawerTitle
          }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="CloseDrawer"
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
                <el-descriptions-item label="ID">{{ configData.id }}</el-descriptions-item>
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
                <el-descriptions-item label="更新时间">{{
                    FormatDateTime(configData.updated_at)
                  }}
                </el-descriptions-item>
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
          <el-tab-pane name="快照" label="快照">
            <div v-if="tabActiveName=='快照'">
              <Snapshot :serverID="serverID"/>
            </div>
          </el-tab-pane>
          <el-tab-pane name="GPU" label="GPU">
            <div class="gva-table-box">
              <div class="gva-btn-list">
                <el-form-item label="GPU类型" >
                  <el-select label="GPU型号" v-model="unBindGpuIds" style="width: 320px" placeholder="请选择gpu" multiple>
                    <el-option v-for="item in unBindGpuList"
                               :value="item.id"
                               :label="item.model+' '+item.addr"
                    />
                  </el-select>
                </el-form-item>
              </div>
              <div style="margin: 20px 0;">
                <el-button style="margin-left: 5px" :disabled="!unBindGpuIds.length"
                           type="primary" @click="HandleManageGPU('bind')">
                  绑定GPU
                </el-button>
                <el-button :disabled="!gpuIds.length" icon="delete" type="grey" @click="HandleManageGPU('unbind')">
                  解绑GPU
                </el-button>
              </div>
              <el-table
                  v-loading="false"
                  :cell-style="{'text-align':'left'}"
                  :data="gpuList"
                  @selection-change="handleSelectionGpuData"
                  :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
                  style="width: 100%;font-size: 15px;"
              >
                <el-table-column label="GPU型号">
                  <template #default="scope">
                    <span>{{ scope.row.model }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="供应商">
                  <template #default="scope">
                    <span>{{ GetSupplierName(scope.row.vendorDeviceId) }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="PCI总线地址">
                  <template #default="scope">
                    <span>{{ scope.row.addr }}</span>
                  </template>
                </el-table-column>
              </el-table>
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
        <el-button type="info" @click="visible=false">关闭</el-button>
      </div>
    </el-drawer>
    <!--    新增和编辑-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="1200px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{
            IsCreated ? '新增' : '编辑'
          }}云主机</span>
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
          :model="serverDataInfo"
          :rules="IsCreated?rules:''"
          label-width="120px"
          label-position="left"
      >
        <el-form-item label="名称" prop="generate_name">
          <el-input style="width: 320px;" class="formStyle" :disabled="true" v-model="serverDataInfo.generate_name"/>
        </el-form-item>
        <el-form-item label="CPU(核)" prop="generate_name">
          <el-input-number style="width: 200px;" controls-position="right" :precision="0" :min="1"
                           v-model="skuData.cpu"/>
        </el-form-item>
        <el-form-item label="内存(GB)" prop="generate_name">
          <el-input-number style="width: 200px;" controls-position="right" :precision="0" :min="1"
                           v-model="skuData.mem"/>
        </el-form-item>
        <!--        <el-form-item label="到期时间">-->
        <!--          <div style="display: flex">-->
        <!--            <el-input-number style="flex: 2;width: 100%" controls-position="right" :precision="0" :min="1" v-model="serverDataInfo.timeNum"/>-->
        <!--            <el-select style="flex: 1;width: 100%" v-model="serverDataInfo.timeType">-->
        <!--              <el-option label="小时" :value="'h'"/>-->
        <!--              <el-option label="天" :value="'d'"/>-->
        <!--              <el-option label="月" :value="'m'"/>-->
        <!--            </el-select>-->
        <!--          </div>-->
        <!--        </el-form-item>-->
        <el-form-item label="GPU" style="padding-left: 10px" label-width="110">
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
            <el-select style="width: 320px;" placeholder="请选择GPU类型"  v-model="item.index">
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
        <!--          <el-button type="text" @click="AddNewGPU">添加透传设备</el-button>-->
        <!--          <el-button v-if="selectDevice.length != 1" type="text" @click="RemoveNewGPU">移除</el-button>-->

        <el-form-item label="绑定租户" style="padding-left: 10px" label-width="110">
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
          <el-select style="width: 320px;" class="formStyle" placeholder="请选择租户" @change="generateUniqueName(serverDataInfo.renterID)"
                     v-model="serverDataInfo.renterID">
            <el-option v-for="item in renterList"
                       :label="item.username"
                       :value="item.renterID"/>
          </el-select>
        </el-form-item>

        <el-form-item label="镜像" prop="imageId">
          <el-select
              style="width: 320px;"
              v-model="serverDataInfo.imageId"
              placeholder="请选择镜像"
              class="formStyle"
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
              style="width: 320px;"
              v-model="serverDataInfo.vpcId"
              placeholder="请选择vpc"
              class="formStyle"
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
              style="width: 320px;"
              v-model="serverDataInfo.network"
              placeholder="请选择IP子网"
              class="formStyle">
            <el-option
                v-for="(optionItem, optionIndex) in IpList"
                :key="optionIndex"
                :disabled="(optionItem.ports-optionItem.ports_used)==0"
                :label="optionItem.name+'('+optionItem.guest_ip_start+'-'+optionItem.guest_ip_end+')    '+'可用：'+(optionItem.ports-optionItem.ports_used)"
                :value="optionItem.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item style="padding-left: 10px" label-width="110" label="绑定弹性公网IP" v-if="serverDataInfo.vpcId!=''&&serverDataInfo.vpcId!='default'">
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
          <el-form-item style="padding-left: 10px" label-width="110" label="网络计费方式">
            <el-tag>按带宽计费</el-tag>
          </el-form-item>
          <el-form-item style="padding-left: 10px" label-width="110" label="带宽(Mbps)" prop="bandwidth">
            <el-input style="width: 200px;" class="formStyle" type="number" placeholder="请输入带宽峰值"
                      v-model.number="serverDataInfo.eip_bw"/>
          </el-form-item>
        </div>
        <el-form-item label="系统盘(GB)" prop="disks"
                      style="display: flex; align-items: center; gap: 0px;padding-left: 10px;;" label-width="110">
          <el-table :data="serverDataInfo.disks"
                    :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
            <el-table-column label="硬盘类别">
              <template #default="scope">
                <el-select style="width: 200px" v-model="scope.row.backend"
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
        <el-form-item label="数据盘(GB)" prop="disks" style="padding-left: 10px" label-width="110">
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
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    创建快照-->
    <el-dialog
        v-model="snapVisible"
        :show-close="false"
        width="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">创建快照</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="snapVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="snapData"
          :model="snapData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="快照名称" prop="name">
          <el-input style="width: 320px;" v-model="snapData.name" placeholder="请输入快照名称"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer" style="justify-content: right;margin-top: 24px;display: flex">
        <el-button type="primary" @click="handleSnapConfirm">确 定</el-button>
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <!--    重置密码-->
    <el-dialog
        v-model="passVisible"
        :show-close="false"
        width="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">重置密码</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="passVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="resetData"
          :model="resetData"
          :rules="rules"
          label-width="120px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input style="width: 320px;" v-model="resetData.username" placeholder="请输入用户名"/>
        </el-form-item>
        <el-form-item  label-width="110px" label="密码类型" style="padding-left: 10px">
          <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
            <el-radio-group v-model="resetData.type">
              <el-radio-button :value="false">随机生成</el-radio-button>
              <el-radio-button :value="true">手工输入</el-radio-button>
            </el-radio-group>
          </div>
        </el-form-item>
        <el-form-item v-if="resetData.type" label="密码" prop="password">
          <el-input style="width: 320px;" v-model="resetData.password" placeholder="请输入密码"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer" style="justify-content: right;margin-top: 24px;display: flex">
        <el-button type="primary" @click="HandleResetConfirm">确 定</el-button>
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <!--    可迁移宿主机列表-->
    <el-drawer
        v-model="migrateDialog"
        :show-close="true"
        size="50%"
        title="可迁移宿主机列表"
    >
      <el-table
          :data="migrateHostList"
      >
        <el-table-column label="名称">
          <template #default="scope">
            {{ scope.row.name }}
          </template>
        </el-table-column>
        <el-table-column label="CPU">
          <template #default="scope">
            <span>{{ '已分配:' + scope.row.cpu_commit }}</span><br/>
            <span>{{ '虚拟总量:' + (scope.row.cpu_commit_bound * scope.row.cpu_count - scope.row.cpu_reserved) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="内存" min-width="120">
          <template #default="scope">
            <span>{{ '已分配:' + FormatCapacitySize(scope.row.mem_commit) }}</span><br/>
            <span>{{ '虚拟总量:' + FormatCapacitySize(scope.row.mem_size - scope.row.mem_reserved) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="存储" min-width="120">
          <template #default="scope">
            <span>{{ '已分配:' + FormatCapacitySize(scope.row.storage_used) }}</span><br/>
            <span>{{ '虚拟总量:' + FormatCapacitySize(scope.row.storage_virtual) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" @click="HandleDeleteDialog(scope.row,'migrate')">迁移</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="getMigrateHostList.pageIndex"
            :page-size="getMigrateHostList.pageSize"
            :page-sizes="[1,2,10, 30, 50, 100]"
            :total="getMigrateHostList.total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="GetMigrateHostsList"
            @size-change="GetMigrateHostsList"
        />
      </div>
    </el-drawer>
    <!--    远程连接-->
    <el-dialog
        v-model="dialogVisible"
        :show-close="false"
        width="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">请输入SSH账户密码</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="dialogVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="sshData"
          :model="sshData"
          :rules="sshRules"
          label-width="120px"
      >
        <el-form-item label="账号" prop="username">
          <el-input style="width: 320px" v-model="sshData.username" placeholder="请输入账号"/>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input style="width: 320px" v-model="sshData.password" placeholder="请输入密码"/>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input style="width: 320px;" type="number" v-model="sshData.port" placeholder="请输入端口"/>
        </el-form-item>

      </el-form>
      <div class="dialog-footer"  style="justify-content: right;margin-top: 24px;display: flex">
        <el-button type="primary" @click="handleConnectConfirm">确 定</el-button>
        <el-button type="info" @click="CloseDrawer">取 消</el-button>
      </div>
    </el-dialog>
    <!--    绑定/修改租户-->
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
    <!--    关联宿主机-->
    <el-drawer v-model="secGroupVisible"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :show-close="false"
               size="60%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">关联云主机</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="secGroupVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <div style="margin-top: 10px;">
        <span>已选：</span>
        <el-tag
            style="margin-left: 2px;margin-top: 2px"
            v-for="(tag, index) in selectData"
            :key="index"
        >
          {{ tag.name }}
        </el-tag>
      </div>
      <div class="gva-btn-list" style="margin-top: 10px;">
        <el-input
            class="filter-item"
            placeholder="请输入名称"
            style="width: 200px;"
            v-model="secSearchInfo.name"
        />

        <el-button type="primary" @click="GetSecHostList">
          搜索
        </el-button>
      </div>
      <el-table
          :data="secHostData"
          @selection-change="handleSelectionChange"
          style="width: 100%"
      >

        <el-table-column type="selection" width="55" :selectable="isRowSelectable"/>
        <el-table-column prop="name" label="名称"/>
        <el-table-column prop="ips" label="IP地址"/>
        <el-table-column label="状态">
          <template #default="scope">
            <i v-if="scope.row.status == 'running'" class="status-dot"/>
            <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
            <el-icon v-else-if="scope.row.status.includes('ing')" class="status-dot-other">
              <Loading/>
            </el-icon>
            <i v-else style="background-color: rgb(174,185,192)" class="status-dot"/>
            <span>{{ hostStatusType[scope.row.status] }}</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="CloseDrawer">取 消</el-button>
        <el-button type="primary" @click="SubmitSecDrawer">确 定</el-button>
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
import Snapshot from "@/view/cloudpods/Storage/snapshot/snapshot1.vue"
import Logs from "@/view/cloudpods/component/log.vue"
import {
  addHost,
  deleteHost,
  deleteHostByIds,
  doMigrate,
  forecastAgent,
  forecastAgentCanInstall,
  getGpuInfo,
  getGpuList,
  getHostList,
  getLoginInfo,
  getMigrateForecastHost,
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
  manageGpu,
  resetPassword,
  setAgentCanInstall,
  setHostRestart,
  setHostShutDown,
  setHostStart
} from "@/api/cloudpods/compute/servers";
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import {getNetIpList} from "@/api/cloudpods/network/ip";
import {getImageList} from "@/api/cloudpods/image/image";
import {format} from 'date-fns'
import {CreateServerSnap, GetVirtualSSH} from "@/api/yunguan/cloudpods/virtual/virtual";
import {getStorageList} from "@/api/cloudpods/storage/storage";
import {
  addRenterRes,
  deleteRenterRes,
  getRenterList,
  queryAllRenterRes,
  queryRenterResInfo,
  updateRenterRes
} from "@/api/yunguan/run/renter";
import statusType from '@/locales/zh-CN.json';
import {addSeverSecgroupRule, revokeSeverSecgroupRule} from "@/api/cloudpods/network/secgroup";
import {getVpcList} from "@/api/cloudpods/network/vpc";
import _ from "lodash";
import WarningBar from "@/components/warningBar/warningBar.vue";
import {detachScalingGroup} from "@/api/cloudpods/compute/spring";
import {graphic} from "echarts";

export default {
  name: 'ComplexTable',
  props: {
    secGroupID: {
      type: String,
      default: ''
    },
    scalingGroup: {
      type: String,
      default: ''
    }
  },
  components: {WarningBar, GvaCard, GvaChart, Snapshot, Logs},
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
      deleteType: '',
      isUseBlock: false,
      isUseLocalBlock: false,
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
          },
          areaStyle: {
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
      supplier: {
        '10de': '英伟达',
      },
      unBindGpuIds: [],
      gpuIds: [],
      manageGpu: {
        vmId: '',
        attachIds: [],
        detachIds: [],
      },
      unBindGpuList: [],
      gpuList: [],
      isCanInstallAgent: true,
      isSetAgent: false,
      isInstallAgent: {
        status: false,
        text: '未安装',
      },
      migrateDialog: false,
      migrateHostList: [],
      getMigrateHostList: {
        vmId: '',
        migrateType: '',
        pageIndex: 1,
        pageSize: 10,
        total: 0,
      },
      migrateInfo: {
        vmId: '',
        hostId: '',
        migrateType: '',
      },
      isUpdateRenterRes: false,
      modRenterRes: {
        id: '',
        resourceID: '',
        startTime: '',
        endTime: '',
      },
      minDiskSize: 10,
      secGroupVisible: false,
      secHostData: [],
      selectData: [],
      secSeverInfo: {
        id: [],
        secGroup: this.secGroupID,
      },
      addSelectData: [],
      secSearchInfo: {
        pageIndex: 1,
        pageSize: 100000,
        total: 0,
        name: '',
      },
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      renterResList: [],
      renterList: [],
      bindRenterData: {
        resourceID: '',
        renterID: '',
        type: 2,
        privateIp: '',
        sshHost: '',
        sshUser: '',
        sshPwd: '',
        aiList: "[]",
      },
      skuData: {
        cpu: 2,
        mem: 2,
      },
      tempSSH: {
        username: '',
        password: '',
      },
      dialogVisible: false,
      sshData: {
        id: '',
        port: 22,
        username: '',
        password: '',
        consoleType: 'vnc',
      },
      tabActiveName: "详情",
      hostStatusType: statusType.status.server,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
        secGroupID: this.secGroupID,
        scalingGroup: this.scalingGroup,
      },
      total: 0,
      visible: false,
      createVisible: false,
      IsCreated: true,
      drawerTitle: '',
      hostDataList: [],
      configData: {},
      webUrl: '',
      serverID: '',
      serverDataInfo: {
        disks: [{
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: 10,
        }],
        network: '',
        imageId: '',
        renterID: '',
        generate_name: '',
        nets: [],
        sku: '',
        vcpu_count: 0,
        vmem_size: 0,
        eip_bw: 30,
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
        vcpu_count: 0,
        vmem_size: 0,
        eip_bw: 30,
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
      imageList: [],
      serverList: [],
      dataDiskList: [],
      dataStorageBlock: null,
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
      snapData: {
        id: '',
        name: '',
      },
      snapVisible: false,
      impiData: {},
      popoverWidth: 180,
      passVisible: false,
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
      orderInfo: {},
      resInfo: {},
      timer: null,
      agentTimer: null,
      ids: [],
      statusList: [],
      isBindRenter: false,
      isEip: false,
      addRenterResInfo: {
        renterID: '',
        type: 2,
        sshHost: '',
        sshUser: '',
        sshPwd: '',
        resourceID: '',
        privateIp: '',
        aiList: '[]',
      },
      vpcList: [],
    }
  },
  created() {
    this.GetAllRenterList()
    this.GetVpcList()
    this.GetImageList()
    this.GetHostList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HostDelete(this.deleteId)
          break;
        case 'restart':
          this.HostReboot(this.deleteId)
          break;
        case 'shut':
          this.HostShutDown(this.deleteId)
          break;
        case 'deleteIds':
          this.onDelete()
          break;
        case 'revoke':
          this.RevokeHost(this.deleteId)
          break;
        case 'migrate':
          this.HandleMigrate(this.deleteId)
          break;
      }
      this.deleteVisible=false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      switch (type) {
        case 'delete':
        case 'deleteIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该云主机，是否继续?';
          break;
        case 'restart':
          this.dialogTitle = '重启';
          this.dialogInfo = '强制重启模式会导致服务器实例当前未保存的数据丢失, 是否继续?';
          break;
        case 'shut':
          this.dialogTitle = '关机';
          this.dialogInfo = '此操作将关机该云主机, 是否继续?';
          break;
        case 'migrate':
          this.dialogTitle = '迁移';
          this.dialogInfo = '此操作将迁移该云主机, 是否继续?';
          break;
        case 'revoke':
          this.dialogTitle = '解绑';
          this.dialogInfo = '此操作将该安全组与该云主机解绑, 是否继续?';
          break;
      }
      this.deleteVisible = true;
    },
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
    handleSortChange({column, prop, order}) {
      switch (order) {
        case "ascending":
          this.searchInfo.ipDesc = "asc";
          break;
        case "descending":
          this.searchInfo.ipDesc = "desc";
          break;
        default:
          this.searchInfo.ipDesc = "";
          break;
      }
      this.GetHostList();
    },
    //移除主机模板实例
    HandleDetachGroup(val) {
      this.$confirm('此操作将永久移除该实例，是否继续', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        detachScalingGroup({vmIds: [val.id], scalingGroupId: this.scalingGroup}).then(res => {
          if (res.code == 0) {
            this.$message({
              type: 'success',
              message: '移除成功',
            })
            this.GetHostList();
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消移除',
        })
      })
    },
    //获取供应商名称
    GetSupplierName(val) {
      let tempData = val.split(':');
      return this.supplier[tempData[0]];
    },
    //GPU多选
    handleSelectionGpuData(val) {
      this.gpuIds = [];
      for (let i = 0; i < val.length; i++) {
        this.gpuIds.push(val[i].id)
      }
    },
    //管理GPU
    HandleManageGPU(type) {
      this.manageGpu = {};
      this.manageGpu.vmId = this.configData.id;
      if (type == 'bind') {
        this.manageGpu.attachIds = _.cloneDeep(this.unBindGpuIds);
      } else {
        this.manageGpu.detachIds = this.gpuIds;
      }
      manageGpu(this.manageGpu).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '操作成功'
          })
          this.GetGPUInfoList();
        }
      })
    },
    //管理GPU
    HandleSetGpu(val) {
      this.visible = true
      this.isCanInstallAgent = true
      this.drawerTitle = val.name
      this.configData = val
      this.serverID = val.id
      this.logData.id = val.id;
      this.tabActiveName = 'GPU'
      if (this.configData.status == 'running') this.GetHostIsInstallAgent();
      this.GetGPUInfoList();
    },
    //获取GPU信息
    GetGPUInfoList() {
      this.unBindGpuIds = [];
      getGpuInfo({'hostId': this.configData.host_id}).then(res => {
        if (res.code == 0) {
          let tempData = res.data;
          this.gpuList = tempData.filter(item => item.vmId == this.configData.id)
          this.unBindGpuList = tempData.filter(item => item.vmId == "")
        }
      })
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
    //获取可迁移宿主机列表
    OpenMigrateDialog(val, type) {
      this.getMigrateHostList.vmId = val.id;
      this.getMigrateHostList.migrateType = type;
      this.migrateInfo.vmId = val.id;
      this.migrateInfo.migrateType = type;
      this.migrateDialog = true;
      getMigrateForecastHost(this.getMigrateHostList).then(res => {
        if (res.code == 0) {
          this.migrateHostList = res.data.list;
          this.getMigrateHostList.total = res.data.total;
        }
      })
    },
    GetMigrateHostsList() {
      getMigrateForecastHost(this.getMigrateHostList).then(res => {
        if (res.code == 0) {
          this.migrateHostList = res.data.list;
          this.getMigrateHostList.total = res.data.total;
        }
      })
    },
    //开始迁移
    HandleMigrate(val) {
      this.migrateInfo.hostId = val;
      doMigrate(this.migrateInfo).then(response => {
        if (response.code == 0) {
          this.migrateDialog = false;
          this.GetHostList()
        }
      })
    },
    //批量删除
    onDelete() {
      deleteHostByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.GetHostList()
        }
      })
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.statusList = [];
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        this.statusList.push(val[i].status)
      }
    },
    //解绑云主机
    RevokeHost(val) {
      this.$confirm('此操作将该安全组与该云主机解绑, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.secSeverInfo.id = [val.id]
        revokeSeverSecgroupRule(this.secSeverInfo).then(response => {
          if (response.code == 0) {
            this.$message({
              type: 'success',
              message: '解绑成功',
            })
            this.GetHostList()
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    },
    //关联云主机确定
    SubmitSecDrawer() {
      this.secSeverInfo.id = this.addSelectData;
      addSeverSecgroupRule(this.secSeverInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '关联成功',
          })
          this.GetHostList()
          this.secGroupVisible = false;
        }
      })
    },
    isRowSelectable(row) {
      return !row.secgroups.some(item => item.id === this.secGroupID);
    },
    // 处理表格中的选择变化
    handleSelectionChange(selection) {
      this.selectData = selection; // 更新完整的选中数据
      this.addSelectData = selection
          .map(item => item.id) // 提取选中行的 id
          .filter(id => !this.secHostData.includes(id)); // 移除 secHostData 中已有的值

    },
    GetSecHostList() {
      getHostList(this.secSearchInfo).then(res => {
        if (res.code == 0) {
          this.secHostData = res.data.data
          this.secSearchInfo.total = res.data.total
        }
      })
    },
    //打开管理云主机对话框
    OpenHostDrawer() {
      this.GetSecHostList()
      getHostList(this.secSearchInfo).then(res => {
        if (res.code == 0) {
          this.secHostData = res.data.data
          this.secSearchInfo.total = res.data.total
        }
      })
      this.selectData = [];
      this.secGroupVisible = true;
    },
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        this.timer = setInterval(() => {
          this.GetHostList();
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
    //打开绑定租户对话框
    GetResourceInfo(val) {
      this.resInfo = {};
      this.orderInfo = {};
      queryRenterResInfo({resourceID: val.id}).then(res => {
        if (res.code == 0) {
          this.resInfo = res.data
          if (res.data.description != "") {
            this.orderInfo = JSON.parse(res.data.description)
            this.orderInfo.specJson = JSON.parse(this.orderInfo.specJson)
            this.orderInfo.storage = JSON.parse(this.orderInfo.storage)
            this.orderInfo.specJson.push({name: '存储类型', value: [this.orderInfo.storage.name], valueIndex: 0})
            this.orderInfo.specJson.push({name: '操作系统', value: [this.orderInfo.os], valueIndex: 0})
          }
        }
        this.renterVisible = true;
      })
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
    //获取所有租户资源列表
    GetAllRenterResList() {
      queryAllRenterRes().then(res => {
        if (res.code == 0) {
          this.renterResList = res.data.list;
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
    //绑定租户确定
    handleBindConfirm() {
      this.$refs.bindRenterData.validate(valid => {
        if (valid) {
          addRenterRes(this.bindRenterData).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '绑定成功',
              })
              this.renterVisible = false;
              this.GetHostList()
              this.bindRenterData = {}
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
    //磁盘类型切换
    HandleDiskChange(type) {
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
    SetSystemStorage(type) {
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
    RemoveSystemStorage(type) {
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
              this.GetHostList()
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
          this.bindRenterData.sshUser = this.impiData.username
          this.bindRenterData.sshPwd = this.impiData.password

          // 根据用户名和密码的长度计算宽度
          const baseWidth = 14; // 基础宽度
          const additionalWidthPerChar = 10; // 每个字符增加的宽度

          const totalLength = usernameLength + passwordLength;
          const dynamicWidth = baseWidth + (totalLength * additionalWidthPerChar);
          this.popoverWidth = dynamicWidth;
        }
      })
    },
    //分页查询
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetHostList()
    },
    //分页查询
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetHostList()
    },
    //打开创建快照对话框
    OpenCreateSnapDialog(val) {
      this.snapData = {}
      this.snapData.id = val.id
      this.snapVisible = true
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
    OpenDrawer(val) {
      this.visible = true
      this.isCanInstallAgent = true
      this.drawerTitle = val.name
      this.configData = val
      this.serverID = val.id
      this.logData.id = val.id
      if (this.configData.status == 'running') this.GetHostIsInstallAgent();
    },
    OpenCreateDrawer() {
      this.isUseBlock = false;
      this.isUseLocalBlock = false;
      this.GetStorageList()
      this.GetGpuList()
      this.serverDataInfo = {};
      this.serverDataInfo.isolated_devices = [];
      this.openDevice = false;
      this.serverList = [{
        index: 0,
        num: 1,
      }]
      this.generateUniqueName(0)
      this.serverDataInfo.vpcId = "default";
      this.netDataPrototype.network = "";
      this.serverDataInfo.imageId = "";
      this.skuData.cpu = 1;
      this.skuData.mem = 1;
      this.dataDiskList = [];
      this.serverDataInfo.disks = [];
      this.serverDataInfo.eip_bw = 30;
      this.serverDataInfo.nets = [];
      this.serverDataInfo.disks = [
        {
          backend: 'local',
          disk_type: 'sys',
          index: 0,
          medium: 'rotate',
          size: 10,
        }
      ]
      this.minDiskSize = 10;
      this.GetIpList();
      this.createVisible = true;
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
        case "GPU":
          this.GetGPUInfoList();
          break;
      }
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
    //添加新磁盘
    AddNewDisk() {
      if (this.dataDiskList.length == 0) {
        this.dataDiskList.push({
          backend: 'local',
          disk_type: 'data',
          index: this.dataDiskList.length + 1,
          medium: 'rotate',
          size: 10,
        })
      } else {
        if (this.dataDiskList[0].storage_id != undefined) {
          this.dataDiskList.push({
            backend: this.dataDiskList[0].backend,
            disk_type: 'data',
            index: this.dataDiskList.length + 1,
            medium: 'rotate',
            storage_id: this.dataDiskList[0].storage_id,
            size: 10,
          })
        } else {
          this.dataDiskList.push({
            backend: this.dataDiskList[0].backend,
            disk_type: 'data',
            index: this.dataDiskList.length + 1,
            medium: 'rotate',
            size: 10,
          })
        }
      }
    },
    //移除磁盘
    RemoveNewDisk() {
      this.dataDiskList.splice(this.dataDiskList.length - 1, 1)
    },
    //添加新GPU
    AddNewGPU() {
      this.selectDevice.push({
        index: 0,
        num: 1,
      })
    },
    //移除GPU
    RemoveNewGPU() {
      this.selectDevice.splice(this.selectDevice.length - 1, 1)
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
      this.tabActiveName = "详情";
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
      if (!this.isEip) {
        this.serverDataInfo.eip_bw = 0;
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
    //获取主机列表
    GetHostList() {
      getHostList(this.searchInfo).then(response => {
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
                  const expiredAt = tempData[0].expired_at;
                  const utcDate = new Date(expiredAt);
                  const chinaTime = new Date(utcDate.getTime() + 8 * 60 * 60 * 1000);
                  this.modRenterRes.endTime = chinaTime.toISOString().replace("Z", "+08:00");
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
      this.$refs.serverDataInfo.validate(valid => {
        if (valid) {

          addHost(this.addServerDataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              if (this.isBindRenter) {
                this.CreateRenterRes(JSON.parse(response.data))
              }
              this.GetHostList()
              this.createVisible = false
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
    FormatMem(val) {
      return val / 1024
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
    //选择套餐后
    skuChange(value) {
      const minDiskSizeTemp = Number(this.imageList.find(item => item.id == this.serverDataInfo.imageId).min_disk / 1024).toFixed(0)
      this.minDiskSize = minDiskSizeTemp < 10 ? 10 : minDiskSizeTemp;
      this.serverDataInfo.disks[0].size = this.minDiskSize;
    },
    //删除
    HostDelete(val) {
      let obj = {}
      obj.id = val
      deleteHost(obj).then(response => {
        if (response.code == 0) {
          this.GetHostList()
        }
      })
    },
    //关机
    HostShutDown(val) {
      let obj = {}
      obj.id = val
      setHostShutDown(obj).then(response => {
        if (response.code == 0) {
          this.GetHostList()
        }
      })
    },
    //开机
    HostBoot(val) {
      let obj = {}
      obj.id = val.id
      setHostStart(obj).then(response => {
        if (response.code == 0) {
          this.GetHostList()
        }
      })
    },
    //重启
    HostReboot(val) {
      let obj = {}
      obj.id = val
      setHostRestart(obj).then(response => {
        if (response.code == 0) {
          this.GetHostList()
        }
      })
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
    //远程连接
    RemoteConnect(val, data) {
      this.isSetAgent = false;
      this.sshData.id = val.id
      this.sshData.port = 22
      this.sshData.consoleType = data
      this.HandleGetLoginInfo(this.sshData)
      if (data == "vnc") {
        this.HandleVirtualRemoteConnect()
      } else {
        this.dialogVisible = true
        this.tempSSH = val
      }
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
.formStyle {
  width: 525px;
}

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
<style scoped>
.example-showcase .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 自定义提示框样式 */
.custom-tooltip-content {
  max-width: 200px;
  max-height: 102px;
  overflow: hidden;
  white-space: normal;
  word-break: break-all;
}

</style>
