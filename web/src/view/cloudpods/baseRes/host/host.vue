<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo"
                 label-width="70"
                 style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="名称">
              <el-input style="width: 240px" v-model="searchInfo.name" placeholder="请输入" clearable/>
            </el-form-item>
            <el-form-item label="启用状态">
              <el-select style="width: 240px" v-model="searchInfo.enabled" placeholder="请选择" clearable>
                <el-option label="启用" :value="true"/>
                <el-option label="禁用" :value="false"/>
              </el-select>
            </el-form-item>
            <el-form-item label="CPU架构">
              <el-select style="width: 240px" v-model="searchInfo.osArch" placeholder="请选择" clearable>
                <el-option label="x86_64" value="x86"/>
                <el-option label="aarch64" value="arm"/>
              </el-select>
            </el-form-item>
          </div>

          <el-form-item>
            <el-button type="primary" @click="GetBaseHostList">查询</el-button>
            <el-button  class="button-gap" type="info" @click="onReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div class="gva-table-box">
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="hostList"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55"/>
        <el-table-column label="名称" width="120" :show-overflow-tooltip="true">
          <template #default="scope">
            <el-button type="text" @click="OpenMoreDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <i v-if="scope.row.host_status=='online'" class="status-dot"/>
            <i v-else
               style="background-color: rgb(174,185,192)" class="status-dot"/>
            <span>{{ scope.row.host_status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="启用状态">
          <template #default="scope">
            <i v-if="scope.row.enabled" class="status-dot"/>
            <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
            <span>{{ scope.row.enabled ? '启用' : '禁用' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="IP" :show-overflow-tooltip="true" width="130">
          <template #default="scope">
            <span>{{ scope.row.access_ip }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU架构">
          <template #default="scope">
            <span>{{ scope.row.cpu_architecture }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="140">
          <template #default="scope">
            <span>{{ '已分配:' + scope.row.cpu_commit }}</span><br/>
            <span>{{ '物理总量:' + scope.row.cpu_count }}</span><br/>
            <span>{{ '虚拟总量:' + (scope.row.cpu_commit_bound * scope.row.cpu_count - scope.row.cpu_reserved) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="内存" width="140">
          <template #default="scope">
            <span>{{ '已分配:' + FormatCapacitySize(scope.row.mem_commit) }}</span><br/>
            <span>{{ '物理总量:' + FormatCapacitySize(scope.row.mem_size) }}</span><br/>
            <span>{{ '虚拟总量:' + FormatCapacitySize(scope.row.mem_size - scope.row.mem_reserved) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="存储" width="140">
          <template #default="scope">
            <span>{{ '已分配:' + FormatCapacitySize(scope.row.storage_used) }}</span><br/>
            <span>{{ '物理总量:' + FormatCapacitySize(scope.row.storage_size) }}</span><br/>
            <span>{{ '虚拟总量:' + FormatCapacitySize(scope.row.storage_virtual) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button link v-if="scope.row.enabled" @click="HandleModStatus(scope.row,'disable','禁用')" type="text">禁用
            </el-button>
            <el-button link v-if="!scope.row.enabled" @click="HandleModStatus(scope.row,'enable','启用')" type="text">启用
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!--      详情-->
    <el-drawer v-model="visible" :close-on-press-escape="true" :show-close="true" :before-close="CloseDrawer"
               :title="'详情'" size="50%">
      <div style="margin: 0 10px">
        <el-tabs v-model="activeName" @tab-change="HandleTabChange">
          <el-tab-pane name="detail" label="详情">
            <el-collapse style="margin-top: 10px" v-model="activeNames">
              <el-collapse-item isActive :name="1">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">基本信息</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="名称">{{ detailData.name }}</el-descriptions-item>
                    <el-descriptions-item label="IP">{{ detailData.access_ip }}</el-descriptions-item>
                    <el-descriptions-item label="MAC地址">{{ detailData.access_mac }}</el-descriptions-item>
                    <el-descriptions-item label="云主机数量">{{ detailData.guests }}</el-descriptions-item>
                    <el-descriptions-item label="内核版本">{{ detailData.metadata.kernel_version }}
                    </el-descriptions-item>
                    <el-descriptions-item label="操作系统">
                      {{ detailData.metadata.os_distribution + '(' + detailData.metadata.os_version + ')' }}
                    </el-descriptions-item>
                    <!--                      <el-descriptions-item label="OVS版本">{{ detailData.metadata.ovs_version }}</el-descriptions-item>-->
                    <!--                      <el-descriptions-item label="内核OVS版本">{{ detailData.metadata.ovs_kmod_version }}-->
                    <!--                      </el-descriptions-item>-->
                    <el-descriptions-item label="创建时间">{{ FormatDateTime(detailData.created_at) }}
                    </el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>
              <el-collapse-item :name="2">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">品牌名称</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="名称">{{ detailData.metadata.manufacture }}</el-descriptions-item>
                    <el-descriptions-item label="型号">{{ detailData.metadata.model }}</el-descriptions-item>
                    <el-descriptions-item label="序列号">{{ detailData.metadata.sn }}</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>
              <el-collapse-item :name="3">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">CPU</div>
                </template>
                <div class="form-section">
                  <el-descriptions :column="1" label-width="120px">
                    <el-descriptions-item label="物理CPU">{{ detailData.cpu_count }}核</el-descriptions-item>
                    <!--                      <el-descriptions-item label="超售比">{{ detailData.cpu_commit_bound }}</el-descriptions-item>-->
                    <el-descriptions-item label="虚拟CPU">{{ detailData.cpu_count * detailData.cpu_commit_bound }}核
                    </el-descriptions-item>
                    <el-descriptions-item label="描述">{{ detailData.sys_info.cpu_info.processors[0].model }}
                    </el-descriptions-item>
                    <el-descriptions-item label="CPU架构">{{ detailData.cpu_architecture }}</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>
              <el-collapse-item :name="4">
                <template #title>
                  <div :style="{ marginLeft: '10px' }">内存</div>
                </template>
                <div class="form-section">
                  <el-descriptions label-width="120px" :column="1">
                    <el-descriptions-item label="物理容量">
                      {{
                        '已使用：' + FormatCapacitySize(detailData.mem_commit) + '/总量：' + FormatCapacitySize(detailData.mem_size)
                      }}
                    </el-descriptions-item>
                    <!--                      <el-descriptions-item label="超售比">{{ detailData.mem_cmtbound }}</el-descriptions-item>-->
                    <el-descriptions-item label="系统预留">{{ FormatCapacitySize(detailData.mem_reserved) }}
                    </el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-collapse-item>
            </el-collapse>
          </el-tab-pane>
          <el-tab-pane name="监控" label="监控">
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
            </div>
          </el-tab-pane>
          <el-tab-pane name="log" label="操作日志">
            <Logs :log-data="logData"/>
          </el-tab-pane>
        </el-tabs>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="visible=false">关闭</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import {
  getBaseHostList,
  getMonitorHostBpsRecvData,
  getMonitorHostBpsSentData,
  getMonitorHostCpuData,
  getMonitorHostDiskData,
  getMonitorHostDiskReadData,
  getMonitorHostDiskWriteData,
  getMonitorHostMemData,
  updateBareHostStatus,
} from "@/api/cloudpods/baseRes/host"
import Logs from "@/view/cloudpods/component/log.vue";
import statusType from "@/locales/zh-CN.json";
import WarningBar from "@/components/warningBar/warningBar.vue";
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import {format} from "date-fns";
import _ from "lodash";
import {graphic} from "echarts";


export default {
  components: {GvaChart, GvaCard, WarningBar, Logs},
  data() {
    return {
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
      activeName: 'detail',
      activeNames:[1,2,3,4],
      logData: {
        id: '',
        isContainer: false,
        ownerKind: '',
        ownerName: '',
        namespace: '',
        clusterId: '',
      },
      visible: false,
      timer: null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
      },
      total: 0,
      hostList: [],
      hostStatus: statusType.status.host,
      updateData: {
        id: '',
        status: '',
      },
      detailData: {},
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
    this.GetBaseHostList()
  },
  methods: {
    //tab按钮切换
    HandleTabChange() {
      switch (this.activeName) {
        case "详情":
          break;
        case "监控":
          this.monitorData.id = this.configData.id
          this.GetMonitorData(this.monitorData)
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
    GetMonitorData(obj) {
      this.clearData();
      getMonitorHostCpuData(obj).then(res => {
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
      getMonitorHostMemData(obj).then(res => {
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
      getMonitorHostDiskData(obj).then(res => {
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
      getMonitorHostDiskReadData(obj).then(res => {
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
      getMonitorHostDiskWriteData(obj).then(res => {
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
      getMonitorHostBpsRecvData(obj).then(res => {
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
      getMonitorHostBpsSentData(obj).then(res => {
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
    //打开详情对话框
    OpenMoreDrawer(val) {
      this.logData.id = val.id;
      this.detailData = val;
      this.monitorData.id = val.id;
      this.visible = true;
      this.GetMonitorData(this.monitorData);
    },
    CloseDrawer() {
      this.visible = false;
    },
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetBaseHostList();
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
    GetBaseHostList() {
      getBaseHostList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData = response.data.data;
          if (JSON.stringify(tempData) !== JSON.stringify(this.hostList)) {
            this.hostList = response.data.data
            this.total = response.data.total
          }
        }
      })
    },
    onReset() {
      this.searchInfo.name = '';
      this.searchInfo.enabled = '';
      this.searchInfo.osArch = '';
      this.GetBaseHostList()
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetBaseHostList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetBaseHostList()
    },
    //修改状态
    HandleModStatus(val, status, text) {
      this.updateData.id = val.id;
      this.updateData.status = status;
      updateBareHostStatus(this.updateData).then(res => {
        if (res.code == 0) {
          this.GetBaseHostList()
          this.$message({
            type: 'success',
            message: text + '成功',
          })
        }
      })
    }
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