<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入名称"
            style="width: 200px;"
        />

        <el-button v-waves class="filter-item" type="primary">
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
        <el-table-column label="磁盘容量">
          <template #default="scope">
            <span>{{ scope.row.size }}</span>
          </template>
        </el-table-column>
        <el-table-column label="是否挂载">
          <template #default="scope">
            <span>{{ scope.row.isMount }}</span>
          </template>
        </el-table-column>
        <el-table-column label="云主机">
          <template #default="scope">
            <span>{{ scope.row.host }}</span>
          </template>
        </el-table-column>
        <el-table-column label="主存储">
          <template #default="scope">
            <span>{{ scope.row.mainStorage }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间">
          <template #default="scope">
            <span>{{ scope.row.createTime }}</span>
          </template>
        </el-table-column>
        <el-table-column label="计费方式">
          <template #default="scope">
            <span>{{ scope.row.model }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span>{{ scope.row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column label="介质类型">
          <template #default="scope">
            <span>{{ scope.row.mediumType }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <el-button type="text" link>扩容</el-button>
          <el-button type="text" link>挂载</el-button>
          <el-button type="text" link>删除</el-button>
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
                  <el-table :data="basicInfoData" style="width: 100%">
                    <el-table-column label="名称" prop="label" width="150"></el-table-column>
                    <el-table-column label="信息" prop="value"></el-table-column>
                  </el-table>
                </el-card>
              </div>
            </div>
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
          name: 'CS',
          size: '10G',
          isMount: '已挂载',
          mainStorage: 'host_10.10.1.20',
          host: 'chen3',
          createTime: '2024-09-29 17:32:051',
          model: '计量计费',
          status: '可用',
          mediumType: '机械盘'
        },
        {
          name: 'test2',
          size: '30G',
          isMount: '已挂载',
          mainStorage: 'host_10.10.1.20',
          host: 'chen3',
          createTime: '2024-09-23 18:24:31',
          model: '计量计费',
          status: '可用',
          mediumType: '机械盘'
        },
        {
          name: 'test3',
          size: '20G',
          isMount: '未挂载',
          host: '',
          mainStorage: 'host_10.10.1.20',
          createTime: '2024-09-15 15:56:31',
          model: '计量计费',
          status: '可用',
          mediumType: '固态硬盘'
        }
      ],
      visible: false,
      drawerTitle: '',
      basicInfoData: [
        {label: 'ID', value: '8be448c8-af6e-4ca7-8da1-9c2868e3b51c'},
        {label: '状态', value: '分配失败'},
        {label: '名称', value: 'CS'},
        {label: '域', value: 'Default'},
        {label: '计费方式', value: '按量付费'},
        {label: '容量', value: '30G'},
        {label: '磁盘制备类型', value: '关闭'},
        {label: '磁盘类型', value: '系统盘'},
        {label: '是否挂载', value: '未挂载'},
        {label: '介质类型', value: '机械盘'},
        {label: '云主机', value: 'chen3'},
        {label: '存储', value: 'host_10.10.1.20'},
        {label: '快照策略', value: '-'},
        {label: '区域', value: 'Default'},
        {label: '可用区', value: 'zone0'},
        {label: '创建时间', value: '2024-09-29 17:32:05'},
        {label: '更新时间', value: '2024-09-29 17:32:09'},
        {label: '备注', value: '-'}
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
.el-drawer.rtl {
  overflow: scroll
}
</style>
