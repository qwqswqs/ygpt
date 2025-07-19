<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:10px">
      <div class="gva-btn-list">
        <el-form label-width="70px" :inline="true" :model="searchInfo"
                 style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="名称">
              <el-input
                  v-model="searchInfo.name"
                  class="filter-item"
                  placeholder="请输入"
                  style="width: 240px;"
              />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="searchInfo.status" style="width: 240px;" placeholder="请选择" clearable>
                <el-option label="可用" value="ready"/>
                <el-option label="故障" value="fail"/>
              </el-select>
            </el-form-item>
            <el-form-item label="CPU架构">
              <el-select style="width: 240px;" v-model="searchInfo.osArch" placeholder="请选择" clearable>
                <el-option label="x86_64" value="x86"/>
                <el-option label="aarch64" value="arm"/>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item>
            <el-button @click="GetSnapshotList" type="primary">
              查询
            </el-button>
            <el-button class="button-gap" type="info"
                       @click="searchInfo.name='';searchInfo.osArch='';GetSnapshotList()">重置
            </el-button>
          </el-form-item>

        </el-form>

      </div>
    </div>
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button :disabled="!ids.length" type="grey" @click="HandleDeleteDialog('','deleteIds')">
          删除
        </el-button>
      </div>
      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="snapshotList"
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
            <i v-if="scope.row.status=='ready'" class="status-dot"/>
            <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
            <span>{{ snapStatus[scope.row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU架构">
          <template #default="scope">
            <span>{{ scope.row.server_metadata.os_arch }}</span>
          </template>
        </el-table-column>
        <el-table-column label="快照大小" sortable="custom">
          <template #default="scope">
            <span>{{ FormatCapacitySize(scope.row.size) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="云主机" width="200">
          <template #default="scope">
            <span>{{ scope.row.guest }}</span>
            <i v-if="scope.row.guest_status=='running'" class="status-dot"/>
            <i v-else-if="scope.row.guest_status=='ready'" class="status-dot"
               style="background-color: rgb(174,185,192)"/>
            <span>{{ hostStatus[scope.row.guest_status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-tooltip :disabled="scope.row.guest_status=='ready'" content="云主机必须在关机状态下才可进行此操作"
                        placement="top">
              <el-button type="text" @click="HandleSnapshotReset(scope.row)"
                         :disabled="scope.row.guest_status!='ready'">回滚主机
              </el-button>
            </el-tooltip>
            <el-button link type="text" @click="HandleDeleteDialog(scope.row,'delete')">删除</el-button>
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

import {deleteSnapshot, deleteSnapshotByIds, getSnapshotList, snapshotReset} from '@/api/cloudpods/storage/snapshot'


export default {
  props: {
    serverID: {
      type: String,
      default: ''
    },
  },
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      ids: [],
      timer: null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        serverID: this.serverID,
        name: '',
      },
      total: 0,
      snapshotList: [],
      hostStatus: {
        ['running']: '开机',
        ['disk']: '分配磁盘中',
        ['stopping']: '关机中',
        ['ready']: '关机',
        ['suspend']: '挂起',
        ['suspending']: '挂起中',
        ['resuming']: '恢复中',
        ['starting']: '启动中',
        ['disk_reset']: '回滚中',
      },
      snapStatus: {
        ['ready']: '可用',
        ['instance_snapshot_start_delete']: '主机快照开始删除',
      },
      snapResetData: {
        hostId: '',
        snapshotId: '',
      }
    }
  },
  created() {
    this.GetSnapshotList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDeleteSnapshot(this.deleteId)
          break;
        case 'deleteIds':
          this.HandleDeleteSnapshotByIds()
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
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
      }
    },
    handleSortChange({column, prop, order}) {
      switch (order) {
        case "ascending":
          this.searchInfo.sizeDesc = "asc";
          break;
        case "descending":
          this.searchInfo.sizeDesc = "desc";
          break;
        default:
          this.searchInfo.sizeDesc = "";
          break;
      }
      this.GetSnapshotList();
    },
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetSnapshotList();
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
    //删除主机快照
    HandleDeleteSnapshot(val) {
      deleteSnapshot({snapshotId: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetSnapshotList(this.searchInfo)
        }
      })
    },
    //批量删除主机快照
    HandleDeleteSnapshotByIds(val) {
      deleteSnapshotByIds({ids: this.ids}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.GetSnapshotList(this.searchInfo)
        }
      })
    },
    //回滚主机
    HandleSnapshotReset(val) {
      this.$confirm('此操作将回滚主机状态, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.snapResetData.hostId = val.guest_id
        this.snapResetData.snapshotId = val.id
        snapshotReset(this.snapResetData).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '回滚成功',
              type: 'success'
            })
            this.GetSnapshotList(this.searchInfo)
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消回滚'
        });
      });
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
    GetSnapshotList() {
      getSnapshotList(this.searchInfo).then(response => {
        if (response.code == 0) {
          const tempData = response.data.data;
          if (JSON.stringify(tempData) !== JSON.stringify(this.snapshotList)) {
            this.snapshotList = response.data.data
            this.total = response.data.total
            this.startInterval()
          } else {
            this.stopInterval()
          }
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetSnapshotList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetSnapshotList()
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
.status-dot {
  margin-left: 5px;
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}
</style>