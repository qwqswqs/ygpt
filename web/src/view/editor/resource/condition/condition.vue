<template>
  <div class="gva-table-box">
    <el-tabs>
      <el-tab-pane label="容器管理">
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <!-- 搜索和重置功能 -->
            <el-input
                v-model="searchQuery"
                class="filter-item"
                placeholder="请输入容器名称"
                style="width: 200px;"
            />

            <el-button class="filter-item" type="primary" @click="handleSearch">
              查询
            </el-button>
          </div>

          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="conList"
              :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column label="序号">
              <template #default="scope">
                <span>{{ scope.row.id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="容器名称">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="镜像名称">
              <template #default="scope">
                <span>{{ scope.row.imageName }}</span>
              </template>
            </el-table-column>
            <el-table-column label="主机名">
              <template #default="scope">
                <span>{{ scope.row.hostName }}</span>
              </template>
            </el-table-column>
            <el-table-column label="CPU使用率">
              <template #default="scope">
                <span>{{ scope.row.cpuUse }}</span>
              </template>
            </el-table-column>
            <el-table-column label="内存使用率">
              <template #default="scope">
                <span>{{ scope.row.memoryUse }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                <span>{{ scope.row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <el-button type="text">编辑</el-button>
              <el-button type="text">删除</el-button>
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
      </el-tab-pane>
      <el-tab-pane label="云主机管理">
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-input
                class="filter-item"
                placeholder="请输入云主机名称"
                style="width: 200px;"
            />
            <el-button class="filter-item" type="primary" @click="handleSearch">
              查询
            </el-button>
          </div>

          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="hostList"
              :header-cell-style="{'text-align':'left','color':'black'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column label="名称">
              <template #default="scope">
                <el-button type="text" @click="OpenDrawer(scope.row)">{{ scope.row.name }}</el-button>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                <i v-if="scope.row.status == '关机'" class="status-dot" style="background-color: rgb(174,185,192)" />
                <i v-else class="status-dot"/>
                <span>{{ scope.row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="IP地址">
              <template #default="scope">
                <span>{{ scope.row.ip }}</span>
              </template>
            </el-table-column>
            <el-table-column label="系统">
              <template #default="scope">
                <span>{{ scope.row.system }}</span>
              </template>
            </el-table-column>
            <el-table-column label="计费方式">
              <template #default="scope">
                <span>{{ scope.row.way }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <el-button type="text">远程连接</el-button>
              <el-button type="text">关机</el-button>
              <el-button type="text">删除</el-button>
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
      </el-tab-pane>
      <el-tab-pane label="裸金属管理">
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-input
                class="filter-item"
                placeholder="请输入裸金属名称"
                style="width: 200px;"
            />
            <el-button class="filter-item" type="primary" @click="handleSearch">
              查询
            </el-button>
          </div>

          <el-table
              v-loading="false"
              :cell-style="{'text-align':'left'}"
              :data="metalList"
              :header-cell-style="{'text-align':'left','color':'black'}"
              style="width: 100%;font-size: 15px;"
          >
            <el-table-column label="序号">
              <template #default="scope">
                <span>{{ scope.row.id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="主机名">
              <template #default="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="IP地址">
              <template #default="scope">
                <span>{{ scope.row.ip }}</span>
              </template>
            </el-table-column>
            <el-table-column label="系统">
              <template #default="scope">
                <span>{{ scope.row.system }}</span>
              </template>
            </el-table-column>
            <el-table-column label="CPU架构">
              <template #default="scope">
                <span>{{ scope.row.cpu }}</span>
              </template>
            </el-table-column>
            <el-table-column label="内存(GB)">
              <template #default="scope">
                <span>{{ scope.row.memory }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                <span>{{ scope.row.status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <el-button type="text">编辑</el-button>
              <el-button type="text">删除</el-button>
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
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
export default {
  name: 'ComplexTable',
  data() {
    return {
      searchInfo: {
        pageIndex:1,
        pageSize:10,
      },
      total: 0,
      conList: [
        { id: 1, name: '容器1', imageName: 'image1', hostName: 'host1', cpuUse: '20.14%', memoryUse: '14.52%', status: '正常' },
        { id: 2, name: '容器2', imageName: 'image2', hostName: 'host2', cpuUse: '0.00%', memoryUse: '0.00%', status: '停止' },
        { id: 4, name: '容器3', imageName: 'image3', hostName: 'host3', cpuUse: '9.45%', memoryUse: '4.25%', status: '重启' }
      ],
      vhostList: [
        { name: 'host01', status: '运行中', cpuNum: '2', memory: '16GB', time: '2024-12-31 12:00:00'}
      ],
      hostList: [
        { name: 'host1', ip: '10.10.254.14', system: 'Debian', safeGroup: 'Default', way: '按量计费', status: '运行中' },
        { name: 'host2', ip: '10.10.254.15', system: 'CentOS', safeGroup: 'Default', way: '按量计费', status: '运行中' },
        { name: 'host3', ip: '10.10.254.16', system: 'Ubuntu', safeGroup: 'Default', way: '按量计费', status: '运行中' }
      ],
      metalList: [
        { id: 1, name: 'host1', ip: '10.10.254.14', system: 'Debian', cpu: 'x86_64', memory: 12, status: '运行中' },
        { id: 2, name: 'host1', ip: '10.10.254.15', system: 'CentOS', cpu: 'x86_64', memory: 12, status: '关机' },
        { id: 3, name: 'host1', ip: '10.10.254.16', system: 'Ubuntu', cpu: 'x86_64', memory: 12, status: '运行中' }
      ]
    }
  }
}
</script>

<style>

.filter-container {
  margin-bottom: 20px;
  margin-top: 20px;
}

</style>

