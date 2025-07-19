<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tabChange="HandleTabChange">
      <el-tab-pane name="group" label="弹性伸缩组">
        <div class="gva-search-box">
          <div class="gva-btn-list">
            <el-form :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
              <el-form-item label="名称">
                <el-input
                    v-model="searchInfo.name"
                    class="filter-item"
                    placeholder="请输入"
                    style="width: 200px;"
                />
              </el-form-item>
              <el-form-item>
                <el-button @click="GetGroupList" type="primary">
                  查询
                </el-button>
                <el-button type="info" @click="searchInfo.name='';GetGroupList()">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="gva-table-box">

          <div style="margin-bottom: 16px;">
            <el-button @click="OpenCreateGroupDrawer"
                       type="primary">
              新增
            </el-button>
            <el-tooltip :disabled="!deleteList.includes(false)" content="请先禁用弹性伸缩组"
                        placement="top">
              <el-button :disabled="!ids.length||deleteList.includes(false)"
                         type="danger" @click="HandleDeleteDialog('','deleteGroupIds')">
                删除
              </el-button>
            </el-tooltip>
          </div>
          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="groupList"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column type="selection" width="55"/>
            <el-table-column label="名称">
              <template #default="scope">
                <el-button type="text" @click="OpenGroupInfoDrawer(scope.row)">{{ scope.row.name }}</el-button>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <span>{{ scope.row.enabled ? "启用" : '禁用' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="伸缩配置">
              <template #default="scope">
                <span>{{ scope.row.guest_template }}</span>
              </template>
            </el-table-column>
            <el-table-column label="当前实例数">
              <template #default="scope">
                <span>{{ scope.row.instance_number }}</span>
              </template>
            </el-table-column>
            <el-table-column label="最大实例数">
              <template #default="scope">
                <span>{{ scope.row.max_instance_number }}</span>
              </template>
            </el-table-column>
            <el-table-column label="起始实例数">
              <template #default="scope">
                <span>{{ scope.row.desire_instance_number }}</span>
              </template>
            </el-table-column>
            <el-table-column label="最小实例数">
              <template #default="scope">
                <span>{{ scope.row.min_instance_number }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text" link v-if="!scope.row.enabled" @click="HandleEnableGroup(scope.row)">启用
                </el-button>
                <el-button type="text" link v-if="scope.row.enabled" @click="HandleDisableGroup(scope.row)">禁用
                </el-button>
                <el-button type="text" link :disabled="scope.row.enabled"
                           :title="scope.row.enabled?'请先禁用弹性伸缩组':''"
                           @click="HandleDeleteDialog(scope.row,'deleteGroup')">删除
                </el-button>
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
      </el-tab-pane>
      <el-tab-pane name="config" label="伸缩配置">
        <div class="gva-search-box">
          <div class="gva-btn-list">
            <el-form :inline="true" :model="searchInfo" style="display: flex;justify-content: space-between">
              <el-form-item label="名称">
                <el-input
                    v-model="searchInfo.name"
                    class="filter-item"
                    placeholder="请输入"
                    style="width: 200px;"
                />
              </el-form-item>
              <el-form-item>
                <el-button @click="GetConfigList" type="primary">
                  查询
                </el-button>
                <el-button type="info" @click="searchInfo.name='';GetConfigList()">
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
            <el-tooltip :disabled="!deleteList.includes(false)" content="存在不能删除的项"
                        placement="top">
              <el-button :disabled="!ids.length||deleteList.includes(false)"
                         type="danger" @click="HandleDeleteDialog('','deleteConfigIds')">
                删除
              </el-button>
            </el-tooltip>
          </div>
          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="configList"
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
                <span>{{ scope.row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="CPU架构">
              <template #default="scope">
                <span>{{ scope.row.content.os_arch }}</span>
              </template>
            </el-table-column>
            <el-table-column label="系统镜像">
              <template #default="scope">
                <span>{{ scope.row.config_info.image }}</span>
              </template>
            </el-table-column>
            <el-table-column label="VPC">
              <template #default="scope">
                <span>{{ scope.row.vpc }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button :disabled="!scope.row.can_delete"
                           :title="scope.row.can_delete?'':scope.row.delete_fail_reason.details" type="text"
                           @click="HandleDeleteDialog(scope.row,'deleteConfig')">删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchInfo.pageIndex"
                :page-size="searchInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="searchInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!--    新增伸缩配置-->
    <el-drawer
        v-model="createVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="60%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增伸缩配置</span>
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
          :rules="rules"
          label-width="100px"
          label-position="left"
      >
        <el-form-item label="名称" prop="generate_name">
          <el-input v-model="serverDataInfo.generate_name"/>
        </el-form-item>
        <div style="display: flex;align-items: center; gap: 10px">
          <el-form-item label="CPU" prop="generate_name">
            <el-input-number style="width: 200px;" controls-position="right" :precision="0" :min="1"
                             v-model="skuData.cpu"/>
            <span class="unit">核</span>
          </el-form-item>
          <el-form-item label="内存" prop="generate_name">
            <el-input-number controls-position="right" :precision="0" :min="1" v-model="skuData.mem"/>
            <span class="unit">GB</span>
          </el-form-item>
        </div>
        <el-form-item label="GPU" style="padding-left: 10px" label-width="90">
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
          <el-form-item v-if="openDevice" label="GPU类型" style="flex: 1;padding-left: 10px" label-width="90">
            <el-select placeholder="请选择GPU类型" style="width: 200px;" v-model="item.index">
              <el-option v-for="(item1,index1) in deviceInfoArray"
                         :label="'类型: '+item1.model"
                         :value="index1"/>
            </el-select>
          </el-form-item>
          <el-form-item v-if="openDevice" label="数量" style="flex: 1;padding-left: 10px" label-width="90">
            <el-radio-group v-model="item.num">
              <div v-for="n in deviceInfoArray[item.index].count" :key="n">
                <el-radio-button :label="n" :value="n"/>
              </div>
            </el-radio-group>
          </el-form-item>
        </div>

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
        <el-form-item label="VPC" prop="vpcId" style="padding-left: 10px" label-width="90">
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
        <el-form-item label="系统盘(GB)" prop="disks"
                      style="display: flex; align-items: center; gap: 0px;padding-left: 10px" label-width="90">
          <el-select style="width: 200px;" v-model="serverDataInfo.disks[0].backend"
                     @change="HandleDiskChange('local',0)">
            <el-option :value="'local'" label="本地硬盘"/>
            <el-option :value="'nfs'" v-if="nfsStoList.length>0" label="云硬盘(NFS)"/>
            <el-option :value="'rbd'" v-if="cephStoList.length>0" label="云硬盘(Ceph)"/>
          </el-select>
          <div class="input-with-unit">
            <el-tooltip
                v-if="serverDataInfo.disks[0].size === minDiskSize"
                :content="'该镜像系统盘最小容量为'+ minDiskSize +'GB'"
                placement="top">
              <el-input-number

                  controls-position="right"
                  :precision="0"
                  v-model="serverDataInfo.disks[0].size"
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
                v-model="serverDataInfo.disks[0].size"
                type="number"
                style="width: 150px;"
                :min="minDiskSize"
                :step="10"
            />
            <span class="unit">GB</span>
          </div>
          <el-button style="margin-left: 10px" v-if="!('storage_id' in serverDataInfo.disks[0])"
                     @click="SetSystemStorage('local',0)">指定存储
          </el-button>
          <el-select style="width: 200px;margin-left: 10px" v-if="'storage_id' in serverDataInfo.disks[0]"
                     v-model="serverDataInfo.disks[0].storage_id">
            <el-option
                v-for="(item,index) in (serverDataInfo.disks[0].backend=='local'?localStoList:serverDataInfo.disks[0].backend=='nfs'?nfsStoList:cephStoList)"
                :label="item.name+' 总量:'+item.usedInfo.capacity+'/已分配:'+item.usedInfo.virtualUsedCapacity+'/已使用:'+item.usedInfo.actualCapacityUsed"
                :value="item.id"/>
          </el-select>
          <el-button style="margin-left: 10px" v-if="'storage_id' in serverDataInfo.disks[0]"
                     @click="RemoveSystemStorage('local',0)">取消
          </el-button>
        </el-form-item>
        <el-form-item label="数据盘(GB)" prop="disks"
                      style="display: flex; align-items: center; gap: 0px;padding-left: 10px" label-width="90">
          <div v-for="(item,index) in dataDiskList" :key="index" style="display: block;width: 100%;">
            <div v-if="index==0">
              <el-select style="width: 200px;" v-model="item.backend" @change="HandleDiskChange('data',index)">
                <el-option :value="'local'" label="本地硬盘"/>
                <el-option :value="'nfs'" v-if="nfsStoList.length>0" label="云硬盘(NFS)"/>
                <el-option :value="'rbd'" v-if="cephStoList.length>0" label="云硬盘(Ceph)"/>
              </el-select>
              <div class="input-with-unit">
                <el-input-number controls-position="right" :precision="0" v-model="item.size" style="width: 150px"
                                 :min="10" :step="10"></el-input-number>
                <span class="unit">GB</span>
              </div>
              <el-button style="margin-left: 10px" v-if="!('storage_id' in item)"
                         @click="SetSystemStorage('data',index)">指定存储
              </el-button>
              <el-select @change="SetDataStorage" style="width: 200px;margin-left: 10px" v-if="'storage_id' in item"
                         v-model="item.storage_id">
                <el-option
                    v-for="(item1,index1) in (item.backend=='local'?localStoList:item.backend=='nfs'?nfsStoList:cephStoList)"
                    :label="item1.name+' 总量:'+item1.usedInfo.capacity+'/已分配:'+item1.usedInfo.virtualUsedCapacity+'/已使用:'+item1.usedInfo.actualCapacityUsed"
                    :value="item1.id"/>
              </el-select>
              <el-button style="margin-left: 10px" v-if="'storage_id' in item"
                         @click="RemoveSystemStorage('data',index)">取消
              </el-button>
            </div>
            <div v-if="index>0" style="margin-top: 5px;">
              <el-tag style="width: 200px;height: 30px" type="primary">
                {{ item.backend == 'local' ? '本地硬盘' : item.backend == 'nfs' ? '云硬盘(NFS)' : '云硬盘(Ceph)' }}
              </el-tag>
              <div class="input-with-unit">
                <el-input-number controls-position="right" :precision="0" v-model="item.size" style="width: 150px"
                                 :min="10" :step="10"></el-input-number>
                <span class="unit">GB</span>
              </div>
              <el-tag style="width: 200px;height: 30px;margin-left: 10px" v-if="'storage_id' in item">
                {{ dataStorageBlock.name }}
              </el-tag>
            </div>
          </div>
          <el-button style="margin-left: 10px" type="text" @click="AddNewDisk">添加新磁盘</el-button>
          <el-button v-if="dataDiskList.length != 0" type="text" @click="RemoveNewDisk">移除磁盘</el-button>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="createVisible=false">取 消</el-button>
        <el-button type="primary" @click="SubmitDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    新增伸缩策略-->
    <el-drawer
        v-model="createPolicyVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="40%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增伸缩策略</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createPolicyVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>

      <el-form
          ref="addPolicyInfo"
          :model="addPolicyInfo"
          :rules="rules"
          label-width="110px"
          label-position="left"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="addPolicyInfo.name"/>
        </el-form-item>
        <div style="display: flex;">
          <el-form-item label="触发时间" style="padding-left: 10px" label-width="100">
            <el-input-number
                controls-position="right"
                :precision="0"
                v-model="addPolicyInfo.cycleTimer.hour"
                type="number"
                :min="1"
                :max="24"
            />
            <span class="unit">点</span>
          </el-form-item>
        </div>
        <el-form-item label="有效时间" prop="dateTime">
          <el-date-picker
              style="width: 100%"
              v-model="addPolicyInfo.dateTime"
              type="daterange"
              placeholder="选择时间"
              :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item label="执行动作" style="padding-left: 10px" label-width="100">
          <div style="display: flex;align-items: center; gap: 10px">
            <el-select style="flex: 1;" v-model="addPolicyInfo.action">
              <el-option :value="'add'" label="增加"/>
              <el-option :value="'remove'" label="移除"/>
              <el-option :value="'set'" label="调整为"/>
            </el-select>
            <el-input-number
                style="flex: 1;"
                controls-position="right"
                :precision="0"
                v-model="addPolicyInfo.number"
                type="number"
                :title="'范围在0~最大实例数('+groupData.max_instance_number+')'"
                :min="0"
                :max="groupData.max_instance_number"
            />
            个实例
          </div>
        </el-form-item>

        <div style="display: flex;align-items: center; gap: 10px">
          <el-form-item label="默认冷却时间" style="padding-left: 10px" label-width="100">
            <el-input-number
                controls-position="right"
                :precision="0"
                v-model="addPolicyInfo.cooling_time"
                type="number"
                title="范围在0~1000"
                :min="0"
                :max="1000"
            />
            <span class="unit">秒</span>
          </el-form-item>
        </div>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="createPolicyVisible=false">取 消</el-button>
        <el-button type="primary" @click="SubmitPolicyDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--    新增弹性伸缩-->
    <el-drawer
        v-model="createGroupVisible"
        :close-on-click-modal="true"
        :close-on-press-escape="true"
        :show-close="false"
        size="40%"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增弹性伸缩</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="createGroupVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>

      <el-form
          ref="addGroupInfo"
          :model="addGroupInfo"
          :rules="rules"
          label-position="left"
          label-width="110px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="addGroupInfo.name"/>
        </el-form-item>
        <el-form-item label="伸缩配置" prop="scalingConfigId" style="margin-top: 30px">
          <el-select
              v-model="addGroupInfo.scalingConfigId"
              placeholder="请选择伸缩配置"
              style="width: 100%"
              @change="HandleSelectConfig">
            <el-option
                v-for="(optionItem, optionIndex) in configList"
                :key="optionIndex"
                :label="'名称:'+optionItem.name"
                :value="optionItem.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="最大实例数" style="padding-left: 10px" label-width="100">
          <el-input-number
              style="width: 100%"
              controls-position="right"
              :precision="0"
              v-model="addGroupInfo.maxNumber"
              type="number"
              title="范围在1~1000"
              :min="1"
              :max="1000"
          />
        </el-form-item>
        <el-form-item label="起始实例数" style="padding-left: 10px" label-width="100">
          <el-input-number
              style="width: 100%"
              controls-position="right"
              :precision="0"
              v-model="addGroupInfo.beginNumber"
              type="number"
              title="范围在最小实例数~最大实例数"
              :min="addGroupInfo.minNumber"
              :max="addGroupInfo.maxNumber"
          />
        </el-form-item>
        <el-form-item label="最小实例数" style="padding-left: 10px" label-width="100">
          <el-input-number
              style="width: 100%"
              controls-position="right"
              :precision="0"
              v-model="addGroupInfo.minNumber"
              type="number"
              title="范围在0~起始实例数"
              :min="0"
              :max="addGroupInfo.beginNumber"
          />
        </el-form-item>
        <el-form-item label="实例移除策略" style="padding-left: 10px" label-width="100">
          <el-select
              v-model="addGroupInfo.shrinkPrinciple"
              placeholder="请选择实例移除策略"
              style="width: 100%">
            <el-option label="优先移除最早创建的实例" :value="'earliest'"/>
            <el-option label="优先移除最晚创建的实例" :value="'latest'"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button @click="createGroupVisible=false">取 消</el-button>
        <el-button type="primary" @click="SubmitGroupDrawer">确 定</el-button>
      </div>
    </el-drawer>
    <!--弹性伸缩组详情-->
    <el-drawer v-model="visible"
               :close-on-click-modal="true"
               :close-on-press-escape="true"
               :show-close="false" size="50%">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">弹性伸缩组详情</span>
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
          <el-tab-pane name="policy" label="伸缩策略">
            <div class="gva-table-box">
              <div class="gva-btn-list">
                <el-form :inline="true" :model="searchPolicyInfo">
                  <el-form-item>
                    <el-input
                        v-model="searchPolicyInfo.name"
                        class="filter-item"
                        placeholder="请输入名称"
                        style="width: 200px;"
                    />
                  </el-form-item>
                  <el-form-item>
                    <el-button @click="GetPolicyList" type="primary">
                      搜 索
                    </el-button>
                    <el-button type="primary" @click="OpenPolicyDrawer">
                      新 建
                    </el-button>
                  </el-form-item>
                </el-form>
              </div>
              <el-table
                  v-loading="false"
                  :cell-style="{'text-align':'left'}"
                  :data="policyList"
                  :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
                  style="width: 100%;font-size: 15px;"
              >
                <el-table-column label="名称">
                  <template #default="scope">
                    <span>{{ scope.row.name }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="启用状态">
                  <template #default="scope">
                    <span>{{ scope.row.enabled ? '启用' : '禁用' }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="触发条件">
                  <template #default="scope">
                    <span>{{ '每天' + GetStartTime(scope.row) + '点触发，有效时间为' + FormatDateTime(scope.row.cycle_timer.start_time) + '至' + FormatDateTime(scope.row.cycle_timer.end_time) }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="执行动作">
                  <template #default="scope">
                    <span>{{ actionType[scope.row.action] + scope.row.number + '个实例' }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="冷却时间">
                  <template #default="scope">
                    <span>{{ scope.row.cooling_time + '秒' }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作">
                  <template #default="scope">
                    <el-button type="text" v-if="!scope.row.enabled" @click="HandleEnablePolicy(scope.row)">启用
                    </el-button>
                    <el-button type="text" v-if="scope.row.enabled" @click="HandleDisablePolicy(scope.row)">禁用
                    </el-button>
                    <el-button type="text" :disabled="scope.row.enabled"
                               :title="scope.row.enabled?'仅在禁用状态下支持该操作':''"
                               @click="HandleDeleteDialog(scope.row,'deletePolicy')">删除
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
              <div class="gva-pagination">
                <el-pagination
                    :current-page="searchPolicyInfo.pageIndex"
                    :page-size="searchPolicyInfo.pageSize"
                    :page-sizes="[10, 30, 50, 100]"
                    :total="searchPolicyInfo.total"
                    layout="total,  prev, pager, next, sizes,jumper"
                    @current-change="handleCurrentChange"
                    @size-change="handleSizeChange"
                />
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane name="host" label="当前实例">
            <Server :scaling-group="scalingGroup"></Server>
          </el-tab-pane>
          <el-tab-pane name="log" label="操作日志">
            <div v-if="tabActiveName=='log'">
              <Logs :log-data="logData"/>
            </div>
          </el-tab-pane>
        </el-tabs>
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
import {getStorageList} from "@/api/cloudpods/storage/storage";
import {
  addScalingConfig,
  addScalingGroup,
  addScalingPolicy,
  deleteScalingConfig,
  deleteScalingConfigByIds,
  deleteScalingGroup,
  deleteScalingGroupByIds,
  deleteScalingPolicy,
  disableScalingGroup,
  disableScalingPolicy,
  enableScalingGroup,
  enableScalingPolicy,
  getScalingConfigList,
  getScalingGroupList,
  getScalingPolicyList
} from "@/api/cloudpods/compute/spring";
import {getNetIpList} from "@/api/cloudpods/network/ip";
import {getVpcList} from "@/api/cloudpods/network/vpc";
import _ from "lodash";
import {getGpuList} from "@/api/cloudpods/compute/servers";
import {getImageList} from "@/api/cloudpods/image/image";
import WarningBar from "@/components/warningBar/warningBar.vue";
import Logs from "@/view/cloudpods/component/log.vue";
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import Snapshot from "@/view/cloudpods/Storage/snapshot/snapshot.vue";
import Server from "@/view/cloudpods/compute/servers/servers.vue";


export default {
  components: {Server, GvaChart, Snapshot, GvaCard, Logs, WarningBar},
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
      deleteType: '',
      ids: [],
      deleteList: [],
      groupData: null,
      actionType: {
        'add': '增加',
        'remove': '移除',
        'set': '调整为',
      },
      createPolicyVisible: false,
      addPolicyInfo: null,
      policyInfo: {
        name: '',
        unit: 's',
        scaling_group: '',
        trigger_type: 'cycle',
        action: 'add',//add remove set
        cooling_time: 180,
        number: '',
        dateTime: [],
        cycleTimer: {
          cycle_type: 'day',
          hour: 8,
          minute: 0,
          startTime: '',
          endTime: '',
        }
      },
      searchPolicyInfo: {
        id: '',
        name: '',
        pageIndex: 1,
        pageSize: 10,
        total: 0,
      },
      policyList: [],
      scalingGroup: '',
      visible: false,
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      tabActiveName: 'policy',
      activeName: 'group',
      searchInfo: {
        name: '',
        pageIndex: 1,
        pageSize: 10,
        total: 0,
      },
      groupList: [],
      configList: [],
      createVisible: false,
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
        isolated_devices: [],
        vpcId: 'default',
      },
      netDataPrototype: {
        network: '',
      },
      rules: {
        generate_name: [
          {required: true, message: '请输入名称', trigger: 'blur'},
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
            trigger: 'blur'
          }
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
        scalingConfigId: [
          {required: true, message: '请选择伸缩配置', trigger: 'blur'}
        ],
        dateTime: [
          {required: true, message: '请选择有效时间', trigger: 'blur'}
        ],
      },
      skuData: {
        cpu: 2,
        mem: 2,
      },
      deviceInfoArray: [],
      openDevice: false,
      selectDevice: [
        {
          index: 0,
          num: 1,
        }
      ],
      imageList: [],
      IpList: [],
      vpcList: [],
      localStoList: [],
      nfsStoList: [],
      cephStoList: [],
      dataDiskList: [],
      dataStorageBlock: null,
      minDiskSize: 10,
      addConfigInfo: null,
      groupInfo: {
        name: '',
        scalingConfigId: '',
        vpc: '',
        maxNumber: 2,
        beginNumber: 0,
        minNumber: 0,
        shrinkPrinciple: 'earliest', //earliest latest
        networks: [''],
      },
      addGroupInfo: null,
      createGroupVisible: false,
    }
  },
  created() {
    this.GetGroupList();
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'deletePolicy':
          this.HandleDeletePolicy(this.deleteId)
          break;
        case 'deleteGroup':
          this.HandleDeleteGroup(this.deleteId)
          break;
        case 'deleteConfig':
          this.handleDeleteConfig(this.deleteId)
          break;
        case 'deleteGroupIds':
          this.HandleDeleteGroupByIds()
          break;
        case 'deleteConfigIds':
          this.handleDeleteConfigByIds()
          break;
      }
      this.deleteVisible = false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      switch (type) {
        case 'deletePolicy':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该伸缩策略，是否继续?';
          break;
        case 'deleteGroup':
        case 'deleteGroupIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该弹性伸缩，是否继续?';
          break;
        case 'deleteConfig':
        case 'deleteConfigIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久删除该伸缩配置，是否继续?';
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
        if (this.activeName == 'group') {
          if (val[i].enabled) {
            this.deleteList.push(false)
          }
        } else {
          if (!val[i].can_delete) {
            this.deleteList.push(false)
          }
        }
      }
    },
    //禁用伸缩策略
    HandleDisablePolicy(val) {
      disableScalingPolicy({id: val.id}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '禁用成功',
          })
          this.GetPolicyList();
        }
      })
    },
    //启用伸缩策略
    HandleEnablePolicy(val) {
      enableScalingPolicy({id: val.id}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '启用成功',
          })
          this.GetPolicyList();
        }
      })
    },
    //删除伸缩策略
    HandleDeletePolicy(val) {
      deleteScalingPolicy({id: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功',
          })
          this.GetPolicyList();
        }
      })
    },
    //格式化时间
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        // 如果日期是不合理的，返回空字符串
        return '';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    },
    //获取触发时间
    GetStartTime(val) {
      let tempData = val.cycle_timer.hour + 8;
      if (tempData > 24) {
        return tempData - 24;
      }
      return tempData;
    },
    //提交伸缩策略
    SubmitPolicyDrawer() {
      this.$refs.addPolicyInfo.validate(valid => {
        if (valid) {
          console.log("开始创建....")
          this.addPolicyInfo.scaling_group = this.groupData.id;
          console.log(this.addPolicyInfo)
          let tempData = this.addPolicyInfo.dateTime;
          this.addPolicyInfo.cycleTimer.startTime = tempData[0];
          this.addPolicyInfo.cycleTimer.endTime = tempData[1];
          let tempHour = this.addPolicyInfo.cycleTimer.hour - 8;
          if (tempHour < 0) {
            this.addPolicyInfo.cycleTimer.hour = 24 + tempHour;
          } else {
            this.addPolicyInfo.cycleTimer.hour = tempHour;
          }
          addScalingPolicy(this.addPolicyInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createPolicyVisible = false;
              this.GetPolicyList();
            }
          })
        }
      })
    },
    //打开创建伸缩策略对话框
    OpenPolicyDrawer() {
      this.createPolicyVisible = true;
      this.addPolicyInfo = _.cloneDeep(this.policyInfo)
    },
    disabledDate(time) {
      return time.getTime() < Date.now() - 8.64e7;
    },
    //获取伸缩策略列表
    GetPolicyList() {
      getScalingPolicyList(this.searchPolicyInfo).then(res => {
        if (res.code == 0) {
          this.policyList = res.data.list;
          this.searchPolicyInfo.total = res.data.total
        }
      })
    },
    //打开弹性伸缩详情对话框
    OpenGroupInfoDrawer(val) {
      this.visible = true;
      this.scalingGroup = val.id;
      this.logData.id = val.id;
      this.searchPolicyInfo.id = val.id;
      this.groupData = val;
      this.GetPolicyList();
    },
    //禁用弹性伸缩
    HandleDisableGroup(val) {
      disableScalingGroup({id: val.id}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '禁用成功',
          })
          this.GetGroupList();
        }
      })
    },
    //启用弹性伸缩
    HandleEnableGroup(val) {
      enableScalingGroup({id: val.id}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '启用成功',
          })
          this.GetGroupList();
        }
      })
    },
    //删除弹性伸缩
    HandleDeleteGroup(val) {
      deleteScalingGroup({id: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功',
          })
          this.GetGroupList();
        }
      })
    },
    //批量删除弹性伸缩
    HandleDeleteGroupByIds(val) {
      deleteScalingGroupByIds({ids: this.ids}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功',
          })
          this.GetGroupList();
        }
      })
    },
    //提交创建弹性伸缩对话框
    SubmitGroupDrawer() {
      this.$refs.addGroupInfo.validate(valid => {
        if (valid) {
          addScalingGroup(this.addGroupInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
              this.createGroupVisible = false;
              this.GetGroupList();
            }
          })
        }
      })
    },
    //选择伸缩配置
    HandleSelectConfig() {
      let tempData = this.configList.filter(item => item.id == this.addGroupInfo.scalingConfigId);
      this.addGroupInfo.vpc = tempData[0].vpc;
      this.addGroupInfo.networks[0] = tempData[0].content.nets[0].network;
    },
    //打开创建弹性伸缩对话框
    OpenCreateGroupDrawer() {
      this.addGroupInfo = _.cloneDeep(this.groupInfo)
      this.createGroupVisible = true;
      this.GetConfigList();
    },
    //删除伸缩配置
    handleDeleteConfig(val) {
      deleteScalingConfig({id: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功',
          })
          this.GetConfigList();
        }
      })
    },
    //批量删除伸缩配置
    handleDeleteConfigByIds(val) {
      deleteScalingConfigByIds({ids: this.ids}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功',
          })
          this.GetConfigList();
        }
      })
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
    //打开创建伸缩配置对话框
    OpenCreateDrawer() {
      this.GetImageList();
      this.GetVpcList();
      this.GetStorageList();
      this.GetGpuList();
      this.serverDataInfo = {};
      this.serverDataInfo.isolated_devices = [];
      this.openDevice = false;
      this.serverList = [{
        index: 0,
        num: 1,
      }]
      this.serverDataInfo.vpcId = "default";
      this.netDataPrototype.network = "";
      this.serverDataInfo.imageId = "";
      this.skuData.cpu = 1;
      this.skuData.mem = 1;
      this.dataDiskList = [];
      this.serverDataInfo.disks = [];
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
    //创建伸缩配置
    CreateConfig() {
      addScalingConfig(this.addConfigInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '创建成功',
          })
          this.createVisible = false;
          this.GetConfigList();
        }
      })
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
      if (this.dataDiskList.length != 0) {
        this.serverDataInfo.disks = this.serverDataInfo.disks.concat(this.dataDiskList)
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
      this.netDataPrototype.network = this.serverDataInfo.network;
      this.serverDataInfo.nets = [];
      this.serverDataInfo.nets.push(this.netDataPrototype)
      this.serverDataInfo.disks[0].image_id = this.serverDataInfo.imageId;
      this.addConfigInfo = _.cloneDeep(this.serverDataInfo)
      for (let i = 0; i < this.addConfigInfo.disks.length; i++) {
        this.addConfigInfo.disks[i].size = this.addConfigInfo.disks[i].size * 1024
      }
      this.CreateConfig()
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
    //选择套餐后
    skuChange(value) {
      const minDiskSizeTemp = Number(this.imageList.find(item => item.id == this.serverDataInfo.imageId).min_disk / 1024).toFixed(0)
      this.minDiskSize = minDiskSizeTemp < 10 ? 10 : minDiskSizeTemp;
      this.serverDataInfo.disks[0].size = this.minDiskSize;
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
    //获取弹性伸缩列表
    GetGroupList() {
      getScalingGroupList(this.searchInfo).then(res => {
        if (res.code == 0) {
          this.groupList = res.data.list;
          this.searchInfo.total = res.data.total;
        }
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
    //获取伸缩配置列表
    GetConfigList() {
      getScalingConfigList(this.searchInfo).then(res => {
        if (res.code == 0) {
          this.configList = res.data.list;
          this.searchInfo.total = res.data.total;
        }
      })
    },
    //tab切换
    HandleTabChange() {
      this.ids = [];
      this.deleteList = [];
      this.searchInfo.name = '';
      this.searchInfo.pageIndex = 1;
      this.searchInfo.pageSize = 10;
      switch (this.activeName) {
        case 'policy':
          this.GetPolicyList();
          break;
        case 'group':
          this.GetGroupList();
          break;
        case 'config':
          this.GetConfigList();
          break;
      }
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      if (this.visible) {
        this.searchPolicyInfo.pageIndex = val
        this.GetPolicyList()
      }
      if (this.activeName == 'group') this.GetGroupList()
      else this.GetConfigList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      if (this.visible) {
        this.searchPolicyInfo.pageSize = val
        this.GetPolicyList()
      }
      if (this.activeName == 'group') this.GetGroupList()
      else this.GetConfigList()
    },
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