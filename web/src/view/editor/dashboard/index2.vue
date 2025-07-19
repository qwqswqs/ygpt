<template>
  <div class="gva-container2" style="height: 100vh; width: 100%; display: flex; flex-direction: column;">

    <!-- 上方表格 -->
    <div class="flex-1 mb-4">
      <gva-card custom-class="h-full">
        <gva-table :data="paginatedData1" title="上方表格"/>
        <div class="gva-pagination">
          <el-pagination
              layout="total,prev, pager, next,  sizes, jumper"
              v-model:current-page="currentPage1"
              v-model:page-size="pageSize1"
              :page-sizes="[10, 30, 50, 100]"
              :total="tableData.length"
              @current-change="handleCurrentChange1"
              @size-change="handleSizeChange1"
          />
        </div>
      </gva-card>
    </div>

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
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { GvaTable, GvaCard } from "./componenst"; // 确保路径正确
import { ElPagination } from 'element-plus';

// 示例数据
const tableData = ref([
  { ranking: 1, title: '标题1', click_num: 100 },
  { ranking: 2, title: '标题2', click_num: 80 },
  // 更多数据...
]);

const tableData2 = ref([
  { ranking: 1, title: '标题3', click_num: 120 },
  { ranking: 2, title: '标题4', click_num: 90 },
  // 更多数据...
]);

const currentPage1 = ref(1);
const currentPage2 = ref(1);
const pageSize1 = ref(10);
const pageSize2 = ref(10);

const paginatedData1 = computed(() => {
  const start = (currentPage1.value - 1) * pageSize1.value;
  const end = currentPage1.value * pageSize1.value;
  return tableData.value.slice(start, end);
});

const paginatedData2 = computed(() => {
  const start = (currentPage2.value - 1) * pageSize2.value;
  const end = currentPage2.value * pageSize2.value;
  return tableData2.value.slice(start, end);
});

const handleCurrentChange1 = (newPage) => {
  currentPage1.value = newPage;
};

const handleSizeChange1 = (newSize) => {
  pageSize1.value = newSize;
  currentPage1.value = 1; // 重置为第一页
};

const handleCurrentChange2 = (newPage) => {
  currentPage2.value = newPage;
};

const handleSizeChange2 = (newSize) => {
  pageSize2.value = newSize;
  currentPage2.value = 1; // 重置为第一页
};
</script>