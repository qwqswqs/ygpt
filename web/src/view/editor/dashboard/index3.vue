<template>
  <div class="gva-container2" style="height: 100vh; width: 100%; display: flex; flex-direction: column;">

    <!-- 上方表格 -->
    <div>
      <el-table
          :data="dataList"
          row-key="ID"
      >
        <el-table-column label="资源编号">
          <template #default="scope">
            <span>{{ scope.row.ID }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="资源类别">
          <template #default="scope">
            <span>{{ resourceType[scope.row.resourceType - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="启用时间" min-width="180">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.startTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="结束时间" min-width="180">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.endTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="访问地址" min-width="120">
          <template #default="scope">
            <el-button @click="OpenSpecDialog(scope.row.specDesc)" type="text">{{ scope.row.url }}</el-button>
          </template>
        </el-table-column>
        <el-table-column  label="SSH地址"align="left" min-width="120">
          <template #default="scope">
            <span>{{ scope.row.sshHost }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SSH端口"align="left" min-width="120">
          <template #default="scope">
            <span>{{ scope.row.sshPort }}</span>
          </template>
        </el-table-column>
        <el-table-column label="SSH用户名"align="left" min-width="120">
          <template #default="scope">
            <span>{{ scope.row.sshUser }}</span>
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

    <el-dialog v-model="openDialog" :before-close="CloseDialog" :title="'规格详情'" width="500px">

      <!-- 下方内容 -->
      <div class="flex-1 flex flex-row gap-4">
        <!-- 左侧资源信息卡片 -->
        <div class="flex-1 p-6 bg-white rounded-lg shadow-md">
          <h3 class="text-xl font-bold mb-4">资源信息</h3>
          <div class="space-y-2">
            <p class="text-gray-700">类别：云主机</p>
            <p class="text-gray-700">内存：4G</p>
            <p class="text-gray-700">CPU：4 核</p>
            <p class="text-gray-700">GPU：CUDA128</p>
            <p class="text-gray-700">带宽：100M</p>
            <p class="text-gray-700">硬盘：100G</p>
            <p class="text-gray-700">ssh 账号：root  ssh 密码：12345</p>
          </div>
          <h3 class="text-xl font-bold mt-4 mb-4">环境信息</h3>
          <div class="space-y-2">
            <p class="text-gray-700">操作系统：Windows10</p>
            <p class="text-gray-700">开发环境：pytorch</p>
            <p class="text-gray-700">运行状态：运行中</p>
          </div>
        </div>
        <!-- 右侧表格 -->
        <div class="flex-1">
          <gva-card custom-class="h-full">
            <gva-table :data="paginatedData2" title="右下侧表格"/>
            <div class="gva-pagination">
              <el-pagination
                  layout="total,prev, pager, next,  sizes, jumper"
                  v-model:current-page="currentPage2"
                  v-model:page-size="pageSize2"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="tableData2.length"
                  @current-change="handleCurrentChange2"
                  @size-change="handleSizeChange2"
              />
            </div>
          </gva-card>
        </div>
      </div>
    </el-dialog>

  </div>
</template>

<script>

import {getRenterInfoList} from "@/api/yunguan/resource/resInfo";

export default {
  data() {
    return {
      resourceType:  ["裸金属", "云主机", "容器", "存储"],
      searchInfo: {
        page: 1,
        pageSize: 10,
        total:0,
      },
      openDialog:false,
      table2:[
        { ranking: 1, title: '标题3', click_num: 120 },
        { ranking: 2, title: '标题4', click_num: 90 },
      ],
      table1:[
        { ranking: 1, title: '标题1', click_num: 120 },
        { ranking: 2, title: '标题2', click_num: 90 },
      ],
      dataList:[],
    }
  },
  created() {
    this.GetRenterInfoList()
  },
  methods: {

    // 获取分区列表
    GetRenterInfoList() {
      getRenterInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.searchInfo.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    OpenSpecDialog(val) {
      this.openDialog = true
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
    handleCurrentChange(val) {
      this.searchInfo.page = val
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
    },
    CloseDialog() {
      this.openDialog = false
    },
  }
}
</script>