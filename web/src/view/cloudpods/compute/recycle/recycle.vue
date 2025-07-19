<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName">
      <el-tab-pane :label="'主机'" name="host">
        <div class="gva-search-box" >
          <div class="gva-btn-list">
            <el-form
                ref="searchForm"
                :inline="true"
                :model="searchHostInfo"
                style="display: flex;justify-content: space-between"
            >
              <el-form-item label="名称">
                <el-input
                    style="width: 240px;"
                    v-model="searchHostInfo.name"
                    placeholder="请输入"
                />
              </el-form-item>
              <el-form-item>
                <el-button

                    type="primary"
                    @click="GetRecycleServerList"
                >
                  查询
                </el-button>
                <el-button type="info" @click="searchHostInfo.name='';GetRecycleServerList()">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <div style="margin-bottom: 16px;">
            <el-button :disabled="!ids.length" type="grey"
                       @click="HandleDeleteDialog('','deleteHostIds')">
              删除
            </el-button>
          </div>
          <!-- 表格展示 -->
          <el-table
              :data="hostList"
              @selection-change="handleSelectionChangeData"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column type="selection" width="55"/>
            <el-table-column label="名称" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="scope">
                <i v-if="scope.row.status == 'running'" class="status-dot"/>
                <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
                <i v-else-if="scope.row.status == 'ready'" style="background-color: rgb(174,185,192)"
                   class="status-dot"/>
                <el-icon v-else class="status-dot-other">
                  <Loading/>
                </el-icon>
                <span>{{ hostStatus[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="IP" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.ips }}</span>
              </template>
            </el-table-column>
            <el-table-column label="所属主机" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.host }}</span>
              </template>
            </el-table-column>
            <el-table-column label="自动清除时间" min-width="120">
              <template #default="scope">
                <span>{{ FormatDateTime(scope.row.auto_delete_at) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" :min-width="150" :fixed="'right'">
              <template #default="scope">
                <el-button
                    :disabled="scope.row.status=='deleting'||scope.row.status=='start_delete'||scope.row.status=='restore_state'"
                    type="text" @click="HandleDeleteDialog(scope.row,'deleteHost')">清除
                </el-button>
                <el-button
                    :disabled="scope.row.status=='deleting'||scope.row.status=='start_delete'||scope.row.status=='restore_state'"
                    type="text" @click="HandleDeleteDialog(scope.row,'recoverHost')">
                  恢复
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchHostInfo.pageIndex"
                :page-size="searchHostInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="searchHostInfo.total"
                layout="total,  prev, pager, next, sizes,jumper"
                @current-change="handleHostCurrentChange"
                @size-change="handleHostSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'硬盘'" name="disk">
        <div class="gva-search-box">
          <div class="gva-btn-list">
            <el-form
                ref="searchForm"
                :inline="true"
                :model="searchDiskInfo"
                style="display: flex;justify-content: space-between"
            >
              <el-form-item label="名称">
                <el-input
                    style="width: 240px;"
                    v-model="searchDiskInfo.name"
                    placeholder="请输入"
                />
              </el-form-item>
              <el-form-item>
                <el-button

                    type="primary"
                    @click="GetRecycleDiskList"
                >
                  查询
                </el-button>
                <el-button type="info" @click="searchDiskInfo.name='';GetRecycleDiskList()">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <!-- 表格展示 -->
          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="diskList"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column type="selection" width="55"/>
            <el-table-column label="名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘容量">
              <template #default="scope">
                <span>{{ scope.row.storage }}</span>
              </template>
            </el-table-column>
            <el-table-column label="启用状态">
              <template #default="scope">
                <i v-if="scope.row.status == 'ready'" class="status-dot"/>
                <i v-else-if="scope.row.status.includes('fail')" style="background-color: red" class="status-dot"/>
                <i v-else-if="scope.row.status == 'unknown'" style="background-color: rgb(174,185,192)"
                   class="status-dot"/>
                <el-icon v-else class="status-dot-other">
                  <Loading/>
                </el-icon>
                <span>{{ diskStatusType[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="磁盘类型">
              <template #default="scope">
                <span>{{ scope.row.type == 'sys' ? '系统盘' : '数据盘' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="主机">
              <template #default="scope">
                <span>{{ scope.row.mounted_by != '' ? scope.row.mounted_by : '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="存储">
              <template #default="scope">
                <span>{{ scope.row.from_vm_nas }}</span>
              </template>
            </el-table-column>
            <el-table-column label="创建时间">
              <template #default="scope">
                <span>{{ FormatDateTime(scope.row.age) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
<!--                <el-button type="text" @click="HandleDeleteDialog(scope.row,'deleteDisk')">清除</el-button>-->
                <el-button type="text" @click="HandleDeleteDialog(scope.row,'recoverDisk')">恢复</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchDiskInfo.pageIndex"
                :page-size="searchDiskInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="searchDiskInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleDiskCurrentChange"
                @size-change="handleDiskSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'镜像'" name="image">
        <div class="gva-search-box" >
          <div class="gva-btn-list">
            <el-form
                ref="searchForm"
                :inline="true"
                :model="searchDiskInfo"
                style="display: flex;justify-content: space-between"
            >
              <el-form-item label="区域名称">
                <el-input
                    style="width: 240px;"
                    v-model="searchImageInfo.name"
                    placeholder="请输入"
                />
              </el-form-item>
              <el-form-item>
                <el-button

                    type="primary"
                    @click="GetRecycleImageList"
                >
                  查询
                </el-button>
                <el-button type="info" @click="searchImageInfo.name='';GetRecycleImageList()">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="gva-table-box">
          <!-- 表格展示 -->
          <el-table
              :data="imageList"
              row-key="ID"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          >
            <el-table-column type="selection" width="55"/>
            <el-table-column label="名称" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" align="left" min-width="120">
              <template #default="scope">
                <span>{{ hostStatus[scope.row.status] }}</span>
              </template>
            </el-table-column>
            <el-table-column label="格式" align="left" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.disk_format }}</span>
              </template>
            </el-table-column>
            <el-table-column label="镜像大小" align="left" min-width="120">
              <template #default="scope">
                <span>{{ FormatImageSize(scope.row.size) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="自动清除时间" align="left" min-width="120">
              <template #default="scope">
                <span>{{ FormatDateTime(scope.row.auto_delete_at) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" align="left" :min-width="150" :fixed="'right'">
              <template #default="scope">
                <el-button type="text" link @click="HandleDeleteDialog(scope.row,'deleteImage')">清除</el-button>
                <el-button type="text" link @click="HandleDeleteDialog(scope.row,'recoverImage')">恢复</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                :current-page="searchImageInfo.pageIndex"
                :page-size="searchImageInfo.pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="searchImageInfo.total"
                layout="total,prev, pager, next,  sizes, jumper"
                @current-change="handleImageCurrentChange"
                @size-change="handleImageSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

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
  clearRecycleHost,
  clearRecycleHostByIds,
  getRecycleServerList,
  recoverRecycleHost
} from "@/api/cloudpods/compute/servers";
import {clearRecycleImage, getRecycleImageList, recoverRecycleImage} from "@/api/cloudpods/image/image";
import {clearRecycleVmEvs, getRecycleVmEvsList, recoverRecycleVmEvs} from "@/api/cloudpods/storage/evs";
import statusType from "@/locales/zh-CN.json";

export default {
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      dialogTitle: '',
      dialogInfo: '',
      deleteType: '',
      timer: null,
      diskStatusType: statusType.status.disk,
      hostStatus: statusType.status.server,
      searchHostInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
        total: 0,
      },
      searchDiskInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
        total: 0,
      },
      searchImageInfo: {
        pageIndex: 1,
        pageSize: 10,
        name: '',
        total: 0,
      },
      activeName: 'host',
      hostList: [],
      imageList: [],
      ids: [],
      statusList: [],
      diskList: [],
    }
  },
  created() {
    this.GetRecycleServerList()
    this.GetRecycleDiskList()
    this.GetRecycleImageList()
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'deleteHost':
          this.ClearRecycleHost(this.deleteId)
          break;
        case 'deleteHostIds':
          this.onDelete();
          break;
        case 'deleteDisk':
          this.ClearRecycleDisk(this.deleteId)
          break;
        case 'deleteImage':
          this.ClearRecycleImage(this.deleteId)
          break;
        case 'recoverHost':
          this.RecoverRecycleHost(this.deleteId)
          break;
        case 'recoverImage':
          this.RecoverRecycleImage(this.deleteId)
          break;
        case 'recoverDisk':
          this.RecoverRecycleDisk(this.deleteId)
          break;
      }
      this.deleteVisible = false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.id;
      this.deleteType = type;
      switch (type) {
        case 'deleteHost':
        case 'deleteHostIds':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久清除该云主机，是否继续?';
          break;
        case 'deleteDisk':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久清除该硬盘，是否继续?';
          break;
        case 'deleteImage':
          this.dialogTitle = '删除';
          this.dialogInfo = '此操作将永久清除该镜像，是否继续?';
          break;
        case 'recoverHost':
          this.dialogTitle = '恢复';
          this.dialogInfo = '此操作将恢复该主机，是否继续?';
          break;
        case 'recoverImage':
          this.dialogTitle = '恢复';
          this.dialogInfo = '此操作将恢复该镜像，是否继续?';
          break;
        case 'recoverDisk':
          this.dialogTitle = '恢复';
          this.dialogInfo = '此操作将恢复该硬盘，是否继续?';
          break;
      }
      this.deleteVisible = true;
    },
    //批量删除
    onDelete() {
      clearRecycleHostByIds({ids: this.ids}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '清除成功',
          })
          this.GetRecycleServerList()
        }
      })
    },
    //批量选中
    handleSelectionChangeData(val) {
      this.ids = [];
      this.statusList = [];
      console.log(val)
      for (let i = 0; i < val.length; i++) {
        this.ids.push(val[i].id)
        this.statusList.push(val[i].status)
      }
    },
    //开启定时器
    startInterval() {
      if (this.timer == null) {
        console.log('定时器启动');
        this.timer = setInterval(() => {
          this.GetRecycleServerList()
          this.GetRecycleImageList()
        }, 10000);
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
    //格式化镜像大小
    FormatImageSize(val) {
      var B = 1
      var KB = 1024
      var MB = 1024 * 1024
      var GB = 1024 * 1024 * 1024
      if (val < KB) {
        return (val / B).toFixed(1) + "B";
      } else if (val >= KB && val <= MB) {
        return (val / KB).toFixed(1) + "KB";
      } else if (val >= MB && val <= GB) {
        return (val / MB).toFixed(1) + "MB";
      } else if (val > GB) {
        return (val / GB).toFixed(1) + "GB";
      }
    },
    //清除回收主机
    ClearRecycleHost(val) {
      clearRecycleHost({ID: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '清除成功',
          })
          this.GetRecycleServerList()
        }
      })
    },
    //恢复回收主机
    RecoverRecycleHost(val) {
      recoverRecycleHost({ID: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '恢复成功',
          })
          this.GetRecycleServerList()
        }
      })
    },
    //清除回收镜像
    ClearRecycleDisk(val) {
      clearRecycleVmEvs({id: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '清除成功',
          })
          this.GetRecycleDiskList()
        }
      })
    },
    //恢复回收镜像
    RecoverRecycleDisk(val) {
      recoverRecycleVmEvs({id: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '恢复成功',
          })
          this.GetRecycleDiskList();
        }
      })
    },
    //清除回收镜像
    ClearRecycleImage(val) {
      clearRecycleImage({ID: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '清除成功',
          })
          this.GetRecycleImageList()
        }
      })
    },
    //恢复回收镜像
    RecoverRecycleImage(val) {
      recoverRecycleImage({ID: val}).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '恢复成功',
          })
          this.GetRecycleImageList()
        }
      })
    },
    //获取主机回收站列表
    GetRecycleServerList() {
      getRecycleServerList(this.searchHostInfo).then(res => {
        if (res.code == 0) {
          const tempData = res.data.data;
          if (tempData !== this.hostList) {
            this.hostList = res.data.data
            this.total = res.data.total
            this.startInterval()
          } else {
            this.stopInterval()
          }
        }
      })
    },
    //获取主机回收站列表
    GetRecycleDiskList() {
      getRecycleVmEvsList(this.searchDiskInfo).then(res => {
        if (res.code == 0) {
          const tempData = res.data.list;
          if (tempData !== this.diskList) {
            this.diskList = res.data.list;
            this.searchDiskInfo.total = res.data.total;
          }
        }
      })
    },
    //获取镜像回收站列表
    GetRecycleImageList() {
      getRecycleImageList(this.searchImageInfo).then(res => {
        if (res.code == 0) {
          const tempData = res.data.data;
          if (JSON.stringify(tempData) !== JSON.stringify(this.imageList)) {
            this.imageList = res.data.data
            this.total = res.data.total
            this.startInterval()
          } else {
            this.stopInterval()
          }
        }
      })
    },
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
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
    handleHostCurrentChange(val) {
      this.searchHostInfo.pageIndex = val
      this.GetRecycleServerList()
    },
    handleHostSizeChange(val) {
      this.searchHostInfo.pageSize = val
      this.GetRecycleServerList()
    },
    handleDiskCurrentChange(val) {
      this.searchDiskInfo.pageIndex = val
      this.GetRecycleDiskList()
    },
    handleDiskSizeChange(val) {
      this.searchDiskInfo.pageSize = val
      this.GetRecycleDiskList()
    },
    handleImageCurrentChange(val) {
      this.searchImageInfo.pageIndex = val
      this.GetRecycleImageList()
    },
    handleImageSizeChange(val) {
      this.searchImageInfo.pageSize = val
      this.GetRecycleImageList()
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