<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
          class="filter-item"
          placeholder="请输入名称"
          style="width: 200px;"
      />

      <el-button v-waves class="filter-item"  type="primary">
        搜索
      </el-button>
      <el-button class="filter-item"  style="margin-left: 10px;" type="primary">
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
      <el-table-column label="格式">
        <template #default="scope">
          <span>{{ scope.row.format }}</span>
        </template>
      </el-table-column>
      <el-table-column label="CPU架构">
        <template #default="scope">
          <span>{{ scope.row.CPUType }}</span>
        </template>
      </el-table-column>
      <el-table-column label="系统">
        <template #default="scope">
          <span>{{ scope.row.system }}</span>
        </template>
      </el-table-column>
      <el-table-column label="镜像大小">
        <template #default="scope">
          <span>{{ scope.row.imageSize }}</span>
        </template>
      </el-table-column>
      <el-table-column label="镜像类型">
        <template #default="scope">
          <span>{{ scope.row.imageType }}</span>
        </template>
      </el-table-column>
      <el-table-column label="共享范围">
        <template #default="scope">
          <span>{{ scope.row.shareRange }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <el-button type="text">禁用</el-button>
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
    <el-drawer v-model="visible" :before-close="CloseDrawer" :title="drawerTitle" size="50%">
      <div style="margin: 0 10px">
        <el-tabs>
          <el-tab-pane label="详情">
            <div class="boxer">
              <div class="brand-info">
                <el-card>
                  <div slot="header" class="clearfix">
                    <span>基本信息</span>
                  </div>
                  <el-table :data="baseInfoData" style="width: 100%">
                    <el-table-column label="名称" prop="label" width="150"></el-table-column>
                    <el-table-column label="信息" prop="value"></el-table-column>
                  </el-table>
                </el-card>
              </div>
              <div class="brand-info">
                <el-card>
                  <div slot="header" class="clearfix">
                    <span>镜像属性</span>
                  </div>
                  <el-table :data="imageAttributesData" style="width: 100%">
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
          name: 'debian10',
          status: '可用',
          format: 'QCOW2',
          CPUType: 'x86_64',
          system: 'Debian',
          imageSize: '1.5G',
          imageType: '公共镜像',
          shareRange: '全局共享'
        },
        {
          name: 'host',
          status: '可用',
          format: 'ISO',
          CPUType: 'x86_64',
          system: 'linux',
          imageSize: '1.5G',
          imageType: '自定义镜像',
          shareRange: '不共享'
        },
        {
          name: 'linux',
          status: '可用',
          format: 'QCOW2',
          CPUType: 'x86_64',
          system: 'CentOS',
          imageSize: '1.5G',
          imageType: '公共镜像',
          shareRange: '全局共享'
        }
      ],
      visible: false,
      drawerTitle: '',
      baseInfoData: [
        {label: 'ID', value: '12de24c0-66c5-485c-81c1-1db7c91b1704'},
        {label: '状态', value: '可用'},
        {label: '名称', value: 'debian10'},
        {label: '域', value: 'Default'},
        {label: '项目', value: 'system'},
        {label: '标签', value: '编辑标签'},
        {label: '共享范围', value: '不共享'},
        {label: '类型', value: '自定义镜像'},
        {label: '大小', value: '945M'},
        {label: '创建时间', value: '2024-09-27 15:18:31'},
        {label: '更新时间', value: '2024-10-08 16:50:09'},
        {label: '备注', value: ''}
      ],
      imageAttributesData: [
        {label: 'CPU架构', value: 'x86_64'},
        {label: '磁盘格式', value: 'iso'},
        {label: '系统语言', value: '-'},
        {label: '最小内存要求', value: '0'},
        {label: '最小磁盘要求', value: '945M'},
        {label: '磁盘驱动', value: '-'},
        {label: '网卡驱动', value: '-'},
        {label: '引导方式', value: 'BIOS'},
        {label: '远程终端协议', value: '-'},
        {label: '应用平台', value: '-'}
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

.brand-info {
  margin: 20px;
}

.boxer {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.el-table-column {
  text-align: left;
}

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50; /* Green color */
  margin-right: 5px;
}
</style>
