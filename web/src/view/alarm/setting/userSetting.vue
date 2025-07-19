<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName"  @tabChange="handleTabChange">
      <el-tab-pane label="基础资源配置" :name="'自定义'"></el-tab-pane>
      <el-tab-pane label="租户默认配置" :name="'默认'"></el-tab-pane>
      <!--      <el-tab-pane label="投诉建议" :name="'投诉'"></el-tab-pane>-->
      <!--      <el-tab-pane label="合同发票" :name="'合同'"></el-tab-pane>-->
    </el-tabs>
    <div class="">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-input
                v-model="searchInfo.name"
                placeholder="请输入配置名称"
            />
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="handleSearch"
            >
              查询
            </el-button>
            <el-button
                type="primary"
                @click="addConfigOpenDialog"
            >
              新增
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div   style="padding: 16px 0">
        <warning-bar v-if="activeName === '自定义'" title="为基础资源（宿主机）设置告警策略" />
      </div>
      <!-- 表格展示 -->
      <el-table
          :data="configList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
<!--        <el-table-column label="配置名称">-->
<!--          <template #default="scope">-->
<!--            <span>{{ scope.row.configName }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column label="监控指标">
          <template #default="scope">
            <span>{{ showMonitorType(scope.row.monitorType) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="一般告警阈值">
          <template #default="scope">
            <span>{{ showMonitorData(scope.row.commonAlertData, scope.row.monitorType) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="严重告警阈值">
          <template #default="scope">
            <span>{{ showMonitorData(scope.row.seriesAlertData, scope.row.monitorType) }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column label="用户类别">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.alertUser.map(user => user === 1? '租户' : user === 2? '管理员' : user === 3? '算力圈租户' : '').join(', ') }}</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column label="通知方式">
          <template #default="scope">
            <span>{{ scope.row.alertWay.map(way => way === 1? '站内通知' : way === 2? '短信' : way === 3? '邮箱' : '').join(', ') }}</span>
          </template>
        </el-table-column>
<!--        <el-table-column label="告警通知模板">-->
<!--          <template #default="scope">-->
<!--            <span>{{ scope.row.noticeText }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column v-if="activeName === '自定义'" label="操作">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,  prev, pager, next, sizes,jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
      <el-drawer
          v-model="openDialog"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
          size="60%"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ IsCreated? '新增' : '编辑' }}自定义配置</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="openDialog = false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close />
          </el-icon>
        </div>

        <el-form
            ref="dataInfo"
            :model="dataInfo"
            :rules="IsCreated? rules : ''"
            label-width="120px"
        >
<!--          <el-form-item label="配置名称" prop="configName">-->
<!--            <el-input v-model="dataInfo.configName" />-->
<!--          </el-form-item>-->
          <el-form-item v-if="IsCreated" label="监控指标" prop="monitorType">
            <el-select v-model="dataInfo.monitorType" @change="monitorTypeChange">
<!--              1.CPU使用率;2.内存使用率；3.磁盘使用率；4.磁盘读速率；5.磁盘写速率；6.网络入带宽；7.网络出带宽；8.网络收包速率；9.网络发包速率；10.GPU使用率；11.GPU显存使用率；12.GPU温度-->
              <template v-for="(item, index) in monitorTypeList">
                <el-option v-if="IsShowOption(item.value)" :key="index" :label="item.label" :value="item.value" />
              </template>
            </el-select>
          </el-form-item>
          <el-form-item style="padding-left: 25px">
            <el-tag type="warning"> {{ "要为一种监控指标设置一个新的配置时，请先删除这个指标旧的配置。" }} </el-tag>
          </el-form-item>
          <el-form-item label="一般告警阈值" prop="commonAlertData">
            <el-select v-model="dataInfo.commonAlertData">
              <template v-if="dataInfo.monitorType===12">
                <el-option v-for="threshold in temperatureOptions" :key="threshold" :label="threshold + '&deg;C'" :value="threshold" />
              </template>
              <template v-else>
                <el-option v-for="threshold in thresholdOptions" :key="threshold" :label="threshold + '%'" :value="threshold" />
              </template>
            </el-select>
          </el-form-item>
          <el-form-item label="严重告警阈值" prop="seriesAlertData">
            <el-select v-model="dataInfo.seriesAlertData">
              <template v-if="dataInfo.monitorType===12">
                <el-option v-for="threshold in temperatureOptions" :key="threshold" :label="threshold + '&deg;C'" :value="threshold" />
              </template>
              <template v-else>
                <el-option v-for="threshold in thresholdOptions" :key="threshold" :label="threshold + '%'" :value="threshold" />
              </template>
            </el-select>
          </el-form-item>
<!--          <el-form-item label="监控时长">-->
<!--            <div style="display: flex; /* 垂直居中 */gap: 10px;">-->
<!--              <el-select v-model="dataInfo.monitorTime" @change="this.dataInfo.monitorInterval = this.monitorConfigData[this.dataInfo.monitorTime][0]">-->
<!--                <el-option value="1h">近1小时</el-option>-->
<!--                <el-option value="6h">近6小时</el-option>-->
<!--                <el-option value="12h">近12小时</el-option>-->
<!--                <el-option value="24h">近24小时</el-option>-->
<!--              </el-select>-->
<!--              <span style="width: 140px">取值间隔：</span>-->
<!--              <el-select style="width: 300px" v-model="dataInfo.monitorInterval">-->
<!--                <el-option v-for="(item,index) in monitorConfigData[dataInfo.monitorTime]"-->
<!--                           :label="convertTimeInterval(item)"-->
<!--                           :value="item"/>-->
<!--              </el-select>-->
<!--            </div>-->
<!--          </el-form-item>-->
          <!--          <el-form-item label="用户类别" prop="alertUser">-->
          <!--            <el-select v-model="dataInfo.alertUser" multiple>-->
          <!--              <el-option label="租户" value="1" />-->
          <!--              <el-option label="管理员" value="2" />-->
          <!--              <el-option label="算力圈租户" value="3" />-->
          <!--            </el-select>-->
          <!--          </el-form-item>-->
          <el-form-item label="通知方式" prop="alertWay" label-width="110" style="padding-left: 10px">
            <el-checkbox-group v-model="dataInfo.alertWay">
              <el-checkbox disabled label="站内通知" :value="1" />
              <el-checkbox label="短信" :value="2" />
              <el-checkbox label="邮箱" :value="3" />
            </el-checkbox-group>
          </el-form-item>
<!--          <el-form-item label="告警通知模板" prop="noticeText">-->
<!--            <el-input v-model="dataInfo.noticeText" />-->
<!--          </el-form-item>-->
        </el-form>

        <template #footer>
        <span class="dialog-footer">
              <el-button @click="closeDialog">取 消</el-button>
              <el-button type="primary" @click="submitDialog">确 定</el-button>
        </span>
        </template>
      </el-drawer>
    </div>
  </div>
</template>

<script>
import { addAlertConfig, deleteAlertConfig, queryAlertConfig, updateAlertConfig } from "@/api/yunguan/config/alertConfig";
import WarningBar from "@/components/warningBar/warningBar.vue";

export default {
  components: {WarningBar},
  data() {
    return {
      monitorTypeData: ["cpu使用率", "内存使用率", "磁盘读速率", "磁盘写速率", "网络入带宽", "网络出带宽"],
      monitorTypeList: [
        {
          label: "CPU使用率",
          value: 1,
        },
        {
          label: "内存使用率",
          value: 2,
        },
        {
          label: "磁盘使用率",
          value: 3,
        },
        {
          label: "网络入带宽",
          value: 4,
        },
        {
          label: "网络出带宽",
          value: 5,
        },
        // {
        //   label: "网络收包速率",
        //   value: 6,
        // },
        // {
        //   label: "网络发包速率",
        //   value: 7,
        // },
        {
          label: "GPU使用率",
          value: 10,
        },
        {
          label: "GPU显存使用率",
          value: 11,
        },
        {
          label: "GPU温度",
          value: 12,
        },
      ],

      monitorConfigData: {
        ['1h']: ['1m', '5m'],
        ['6h']: ['3m', '5m', '10m', '15m',],
        ['12h']: ['5m', '10m', '20m', '30m'],
        ['24h']: ['10m', '20m', '30m', '60m'],
        ['168h']: ['30m', '1h'],
        ['720h']: ['6h', '24h'],
      },
      searchInfo: {
        page: 1,
        pageSize: 10,
        name: "",
        monitorType: "",
        alertConfigType: 2,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      dataInfo: {
        ID: 0,
        configName: "",
        monitorType: "",
        commonAlertData: "",
        seriesAlertData: "",
        alertUser: [],
        alertWay: [],
        noticeText: "",
        monitorTime:'1h',
        monitorInterval:'1m',
      },
      activeName: '自定义',
      rules: {
        // configName: [
        //   { required: true, message: "请输入配置名称", trigger: "blur" },
        // ],
        monitorType: [
          { required: true, message: "请选择监控指标", trigger: "change" },
        ],
        commonAlertData: [
          { required: true, message: "请选择一般告警阈值", trigger: "change" },
        ],
        seriesAlertData: [
          { required: true, message: "请选择严重告警阈值", trigger: "change" },
        ],
        alertUser: [
          { required: true, message: "请选择用户类别", trigger: "change" },
        ],
        noticeText: [
          { required: true, message: "请输入告警通知模板", trigger: "blur" },
        ],
      },
      configList: [],
      customConfigList: [],
      thresholdOptions: Array.from({ length: 11 }, (_, i) => `${50 + i * 5}`),
      temperatureOptions: Array.from({ length: 6 }, (_, i) => `${75 + i * 5}`),

    };
  },
  created() {
    this.GetConfigList();
  },
  methods: {
    // 指标名
    showMonitorType(val) {
      for (let i = 0; i < this.monitorTypeList.length; i++) {
        if (this.monitorTypeList[i].value === val) {
          return this.monitorTypeList[i].label;
        }
      }
    },
    // 阈值单位
    showUnit(val) {
      if (val >3&&val<8) {
        return "Mbps";
      }
      else if(val === 12){
        return "°C";
      }
      return "%";
    },
    // 展示监控指标阈值和单位
    showMonitorData(val,type) {
      if (type >3&&type<8) {
        if (val < 1024) {
          return val+"bps";
        }else if (val < 1024 * 1024) {
          return (val / 1024).toFixed(2) + "Kbps";
        } else if (val < 1024 * 1024 * 1024) {
          return (val / 1024 / 1024).toFixed(2) + "Mbps";
        } else {
          return (val / 1024 / 1024 / 1024).toFixed(2) + "Gbps";
        }
      }
      else if(type === 12){
        return val+"°C";
      }
      return val+"%";
    },
    //监控类型改变
    monitorTypeChange() {
        this.dataInfo.commonAlertData = null;
        this.dataInfo.seriesAlertData = null;
    },
    //tab切换
    handleTabChange() {
      switch (this.activeName) {
        case "自定义":
          this.searchInfo.alertConfigType = 2;
          break;
        case "默认":
          this.searchInfo.alertConfigType = 1;
          break;
      }
      this.GetConfigList();
    },
    //监控项显示
    IsShowOption(val){
      const tempData=this.customConfigList.filter(item=>item.monitorType==val)
      return tempData.length>0?false:true;
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
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === "0001-01-01T00:00:00.000Z") {
        // 如果日期是不合理的，返回空字符串
        return "";
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getMinutes()).padStart(2, "0");
      const seconds = String(date.getSeconds()).padStart(2, "0");
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    handleDelete(val) {
      // 由于隐藏了删除按钮，这里的方法可以保留但实际不会被调用
      this.$confirm("此操作将永久删除该配置, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        const params = {
          id: val.ID,
        };
        deleteAlertConfig(params).then((response) => {
          if (response.code == 0) {
            this.$message({
              message: "删除成功",
              type: "success",
            });
            this.GetConfigList();
          } else {
            this.$message({
              message: "删除失败",
              type: "error",
            });
          }
        });
      }).catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除",
        });
      });
    },
    handleEdit(val) {
      // 转换监控指标值
      this.dataInfo.monitorType = String(val.monitorType);
      // 转换用户类别值
      this.dataInfo.alertUser = val.alertUser;
      // 转换通知方式值
      this.dataInfo.alertWay = val.alertWay;

      this.dataInfo.ID = val.ID;
      this.dataInfo.configName = val.configName;
      if (val.monitorType===6||val.monitorType===7) {
        this.dataInfo.commonAlertData=val.commonAlertData/1024/1024
        this.dataInfo.seriesAlertData=val.seriesAlertData/1024/1024
      } else {
        this.dataInfo.commonAlertData = String(val.commonAlertData);
        this.dataInfo.seriesAlertData = String(val.seriesAlertData);
      }
      this.dataInfo.noticeText = val.noticeText;
      this.dataInfo.monitorTime=val.monitorTime;
      this.dataInfo.monitorInterval=val.monitorInterval;

      this.IsCreated = false;
      this.openDialog = true;
    },
    handleSearch() {
      this.GetConfigList();
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val;
      this.GetConfigList();
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val;
      this.GetConfigList();
    },
    addConfigOpenDialog() {
      this.GetCustomConfigList()
      this.openDialog = true;
      this.IsCreated = true;
      this.dataInfo = {
        id: 0,
        configName: "",
        monitorType: "",
        commonAlertData: "",
        seriesAlertData: "",
        alertUser: [],
        alertWay: [1],
        noticeText: "",
        monitorTime: '1h',
        monitorInterval: '1m'
      };
    },
    closeDialog() {
      this.openDialog = false;
    },
    UpdateConfig() {
      const params = {
        id: this.dataInfo.ID,
        configName: this.dataInfo.configName,
        noticeText: this.dataInfo.noticeText,
        monitorType: parseInt(this.dataInfo.monitorType),
        commonAlertData: parseInt(this.dataInfo.commonAlertData),
        seriesAlertData: parseInt(this.dataInfo.seriesAlertData),
        alertUser: this.dataInfo.alertUser.map(Number),
        alertWay: this.dataInfo.alertWay.map(Number),
        monitorTime:this.dataInfo.monitorTime,
        monitorInterval:this.dataInfo.monitorInterval,
      };
      updateAlertConfig(params).then((response) => {
        if (response.code == 0) {
          this.$message({
            message: "编辑成功",
            type: "success",
          });
          this.openDialog = false;
          this.GetConfigList();
        } else {
          this.$message({
            message: "编辑失败",
            type: "error",
          });
        }
      });
    },
    // 新增配置
    AddConfig() {
      this.$refs.dataInfo.validate((valid) => {
        if (valid) {
          const params = {
            configName: this.dataInfo.configName,
            noticeText: this.dataInfo.noticeText,
            monitorType: parseInt(this.dataInfo.monitorType),
            commonAlertData: parseInt(this.dataInfo.commonAlertData),
            seriesAlertData: parseInt(this.dataInfo.seriesAlertData),
            alertUser: this.dataInfo.alertUser.map(Number),
            alertWay: this.dataInfo.alertWay.map(Number),
            monitorTime:this.dataInfo.monitorTime,
            monitorInterval:this.dataInfo.monitorInterval,
            alertConfigType:2,
          };
          // 在界面上输入网络带宽告警阈值为了简化操作，使用的单位是Mbps，在写入数据库前将输入的带宽值转为bps
          if (this.dataInfo.monitorType===6||this.dataInfo.monitorType===7) {
            params.commonAlertData=params.commonAlertData*1024*1024
            params.seriesAlertData=params.seriesAlertData*1024*1024
          }
          addAlertConfig(params).then((response) => {
            if (response.code == 0) {
              this.$message({
                message: "添加成功",
                type: "success",
              });
              this.openDialog = false;
              this.GetConfigList();
            }
          });
        }
      });
    },
    // 提交
    submitDialog() {
      if (this.IsCreated) this.AddConfig();
      else this.UpdateConfig();
    },
    // 获取配置列表
    GetConfigList() {
      const params = {
        pageIndex: this.searchInfo.page,
        pageSize: this.searchInfo.pageSize,
        alertConfigType:this.searchInfo.alertConfigType,
      };
      queryAlertConfig(params).then((response) => {
        if (response.code == 0) {
          this.configList = response.data.list;

          this.total = response.data.total;
        } else {
          this.$message({
            message: "获取配置列表失败",
            type: "error",
          });
        }
      });
    },
    //获取自定义配置列表
    // 获取配置列表
    GetCustomConfigList() {
      const params = {
        pageIndex: this.searchInfo.page,
        pageSize: this.searchInfo.pageSize,
        alertConfigType:2,
      };
      queryAlertConfig(params).then((response) => {
        if (response.code == 0) {
          this.customConfigList = response.data.list;

          // this.total = response.data.total;
        } else {
          this.$message({
            message: "获取配置列表失败",
            type: "error",
          });
        }
      });
    },
    // 修改状态
    changeStatusHandle(val) {
      const params = {
        id: val.id,
      };
      updateAlertConfig(params).then((response) => {
        if (response.code == 0) {
          this.$message({
            message: "状态修改成功",
            type: "success",
          });
        } else {
          this.$message({
            message: "状态修改失败",
            type: "error",
          });
        }
      });
    },
  },
};
</script>

<style scoped>
.gva-table-box {
  margin-top: 20px;
}

.gva-btn-list {
  margin-bottom: 10px;
}

.gva-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>