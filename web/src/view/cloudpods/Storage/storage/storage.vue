<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tab-change="HandleTabChange">
      <el-tab-pane label="宿主机" name="local"></el-tab-pane>
      <el-tab-pane label="物理机" name="baremetal"></el-tab-pane>
    </el-tabs>
    <div class="gva-search-box" style="margin-top:16px">
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
            <el-button @click="GetStorageList" type="primary">
              查询
            </el-button>
            <el-button @click="searchInfo.name='';GetStorageList()" type="info">
              重置
            </el-button>
          </el-form-item>

        </el-form>

      </div>
    </div>
    <div class="gva-table-box">

      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="localList"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column type="selection" width="55"/>
        <el-table-column label="名称" width="280px">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="启用状态">
          <template #default="scope">
            <i v-if="scope.row.enabled" class="status-dot"/>
            <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
            <span>{{ scope.row.enabled ? "启用" : "禁用" }}</span>
          </template>
        </el-table-column>
        <el-table-column label="物理容量">
          <template #default="scope">
            <span style="margin-right: 4px">{{ '已使用:' + scope.row.usedInfo.actualCapacityUsed }}</span>
            <span>{{ '总量:' + scope.row.usedInfo.capacity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="虚拟容量">
          <template #default="scope">
            <span style="margin-right: 4px">{{ '已分配:' + scope.row.usedInfo.virtualUsedCapacity }}</span>
            <span>{{ '总量:' + scope.row.usedInfo.virtualCapacity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button link v-if="scope.row.enabled" type="text" @click="HandleDeleteDialog(scope.row,'delete')">禁用
            </el-button>
            <el-button link v-if="!scope.row.enabled" type="text" @click="HandleEnable(scope.row)">启用</el-button>
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
          </el-icon>禁用</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="deleteVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <span>此操作将禁用该存储,禁用后,无法使用该存储创建云主机磁盘, 是否继续?</span>
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

import {disableNasStorage, enableNasStorage, getStorageList} from '@/api/cloudpods/storage/storage'


export default {
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      timer: null,
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        storageType: 'local',
        name: '',
      },
      activeName: 'local',
      total: 0,
      localList: [],
    }
  },
  created() {
    this.GetStorageList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.HandleDisable(this.deleteId)
          break;
      }
      this.deleteVisible = false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      this.deleteVisible = true;
    },
    //tab切换
    HandleTabChange() {
      switch (this.activeName) {
        case 'local':
          this.searchInfo.storageType = "local";
          break;
        case 'baremetal':
          this.searchInfo.storageType = "baremetal";
          break;
      }
      this.GetStorageList();
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
    GetStorageList() {
      getStorageList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.localList = response.data.list;
          this.total = response.data.total
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetStorageList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetStorageList()
    },
    //禁用
    HandleDisable(val) {
      disableNasStorage({id: val}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '禁用成功',
            type: 'success'
          })
          this.GetStorageList(this.searchInfo)
        }
      })
    },
    //启用
    HandleEnable(val) {
      enableNasStorage({id: val.id}).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '启用成功',
            type: 'success'
          })
          this.GetStorageList(this.searchInfo)
        }
      })
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
</style>