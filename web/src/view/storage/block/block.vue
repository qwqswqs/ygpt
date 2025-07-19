<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入名称"
            style="width: 200px;"
        />

        <el-button v-waves class="filter-item"  type="primary">
          搜索
        </el-button>
        <el-button class="filter-item" style="margin-left: 10px;" type="primary">
          添加
        </el-button>
      </div>

      <el-table
          v-loading="false"
          :cell-style="{'text-align':'left'}"
          :data="list"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%;font-size: 15px;"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="OpenDrawer(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <i class="status-dot"/>
            <span>{{ scope.row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="启用状态">
          <template #default="scope">
            <i class="status-dot"/>
            <span>{{ scope.row.enableStatus }}</span>
          </template>
        </el-table-column>
        <el-table-column label="物理容量">
          <template #default="scope">
            <span>{{ scope.row.physicalSize }}</span>
          </template>
        </el-table-column>
        <el-table-column label="虚拟容量">
          <template #default="scope">
            <span>{{ scope.row.virtualSize }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <el-button link type="text">修改属性</el-button>
          <el-dropdown split-button style="margin-left: 10px" type="text">
            更多
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>启用</el-dropdown-item>
              <el-dropdown-item disabled>禁用</el-dropdown-item>
              <el-dropdown-item>删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
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
    <el-drawer v-model="visible" :before-close="CloseDrawer" :title="drawerTitle" size="50%">
      <div style="margin: 0 10px">
        <el-tabs>
          <el-tab-pane label="详情">
            <div class="boxer">
              <div class="brand-info">
                <el-card>
                  <div slot="header" class="clearfix">
                    <span>用量统计</span>
                  </div>
                  <el-table :data="usageStatsData" style="width: 100%">
                    <el-table-column label="名称" prop="label" width="150"></el-table-column>
                    <el-table-column label="信息" prop="value"></el-table-column>
                  </el-table>
                </el-card>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="宿主机">
            <el-table
                v-loading="false"
                :cell-style="{'text-align':'left'}"
                :data="hardwareList"
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
                  <i class="status-dot"/>
                  <span>{{ scope.row.status }}</span>
                </template>
              </el-table-column>
              <el-table-column label="启用状态">
                <template #default="scope">
                  <i class="status-dot"/>
                  <span>{{ scope.row.enableStatus }}</span>
                </template>
              </el-table-column>
              <el-table-column label="IP">
                <template #default="scope">
                  <span>{{ scope.row.IP }}</span>
                </template>
              </el-table-column>
              <el-table-column label="CPU架构">
                <template #default="scope">
                  <span>{{ scope.row.CPUSchema }}</span>
                </template>
              </el-table-column>
              <el-table-column label="物理内存">
                <template #default="scope">
                  <span>{{ scope.row.memory }}</span>
                </template>
              </el-table-column>
              <el-table-column label="物理磁盘">
                <template #default="scope">
                  <span>{{ scope.row.disk }}</span>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-drawer>
  </div>
</template>

<script>

export default {
  name: 'ComplexTable',
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
      },
      total: 0,
      list: [
        {
          name: 'host_10.10.254.44_local',
          status: '在线',
          enableStatus: '启用',
          physicalSize: '已使用：89G' + '\n' + '总量:198G',
          virtualSize: '已分配：0G' + '\n' + '总量:195G'
        },
        {
          name: 'host_10.10.1.31_local',
          status: '在线',
          enableStatus: '启用',
          physicalSize: '已使用：89G' + '\n' + '总量:198G',
          virtualSize: '已分配：0G' + '\n' + '总量:195G'
        },
        {
          name: 'host_10.10.254.43_local',
          status: '在线',
          enableStatus: '启用',
          physicalSize: '已使用：89G' + '\n' + '总量:198G',
          virtualSize: '已分配：0G' + '\n' + '总量:195G'
        }
      ],
      visible: false,
      drawerTitle: '',
      usageStatsData: [
        {label: '物理容量', value: '已使用:89G/总量:195G'},
        // {label: '超售比', value: '1'},
        {label: '虚拟容量', value: '195G'},
        {label: '分配率', value: '0'},
        {label: '分配', value: '-'},
        {label: '预留', value: '-'},
        {label: '浪费', value: '-'}
      ],
      hardwareList: [
        {
          name: 'debian-test-host',
          status: '运行中',
          enableStatus: '启用',
          IP: '10.10.254.44',
          CPUSchema: 'x86_64',
          memory: '20G/11%',
          disk: '195G/15%'
        },
        {
          name: 'debian-test-host',
          status: '运行中',
          enableStatus: '启用',
          IP: '10.10.254.44',
          CPUSchema: 'x86_64',
          memory: '20G/11%',
          disk: '195G/15%'
        },
        {
          name: 'debian-10-host',
          status: '运行中',
          enableStatus: '禁用',
          IP: '10.10.1.31',
          CPUSchema: 'x86_64',
          memory: '252G/0%',
          disk: '2.7T/0%'
        }
      ]
    }
  },
  methods: {
    OpenDrawer(val) {
      this.visible = true
      this.drawerTitle = val.name
    },
    CloseDrawer() {
      this.visible = false
    }
  }
}
</script>
<style lang="scss">
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}
</style>
