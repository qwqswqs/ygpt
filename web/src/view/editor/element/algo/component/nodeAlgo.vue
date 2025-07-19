<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-input
                v-model="searchInfo.name"
                placeholder="请输入算法名称"
            />
          </el-form-item>
          <el-form-item>
            <el-button

                type="primary"
                @click="handleSearch"
            >
              查询
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table
          :data="algoList"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用次数">
          <template #default="scope">
            <span>{{ scope.row.useTime }}</span>
          </template>
        </el-table-column>
        <el-table-column label="标签">
          <template #default="scope">
            <el-tag
                v-for="(tag,index) in scope.row.usage"
                :key="index"
                :disable-transitions="false"
            >
              <span> {{ tag }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="文件名称">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ GetFileName(scope.row.fileName) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="文件大小">
          <template #default="scope">
            <span>{{ GetFileSize(scope.row.fileSize) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="上传时间">
          <template #default="scope">
            <span>{{FormatDateTime( scope.row.uploadTime)}}</span>
          </template>
        </el-table-column>
<!--        <el-table-column label="操作">-->
<!--          <template #default="scope">-->
<!--            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>-->
<!--            <el-button type="text" @click="DeleteAlgoInfo(scope.row)">删除</el-button>-->
<!--          </template>-->
<!--        </el-table-column>-->
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
  </div>
</template>
<script>
import {
  getNodeElementList,
} from "@/api/yunguan/element/element";


export default {
  name: "NodeAlgo",
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
        type: 3, //1模型 3算法 2数据集
        openType: 2,//不开放1，节点内可见可用2，开放到算力节点节点内用3，开放到算力节点且可外部用4
        status: 1,//1正常、2禁用
        name: ''
      },
      total: 0,
      tagsList: ["文本", "图片"],
      algoList: []
    }
  },
  created() {
    this.GetAlgoList()
  },
  methods: {
    FormatDateTime(dateString){
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    handleSearch() {
      this.GetAlgoList()
    },
    GetFileName(val) {
      // 使用逗号分割文件名
      var parts = val.split(',');
      // 返回分割后的第一个元素，即“,”前的文件名
      return parts[0];
    },
    GetFileSize(val) {
      var B = 1024
      var KB = 1024 * 1024
      var MB = 1024 * 1024 * 1024
      var GB = 1024 * 1024 * 1024 * 1024
      if (val < KB) {
        return (val / B).toFixed(2) + "B";
      } else if (val >= KB && val <= MB) {
        return (val / KB).toFixed(2) + "KB";
      } else if (val >= MB && val <= GB) {
        return (val / MB).toFixed(2) + "MB";
      } else if (val > GB) {
        return (val / GB).toFixed(2) + "GB";
      }
    },
    GetAlgoList() {
      getNodeElementList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.algoList = response.data.list
          this.searchInfo.pageIndex = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.pageIndex = val
      this.GetAlgoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetAlgoList()
    }
  }
}
</script>
