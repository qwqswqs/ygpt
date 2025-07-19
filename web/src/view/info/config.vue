<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <!-- 条件筛选栏 -->
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :model="searchInfo" :inline="true" style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="配置项">
              <el-select v-model="searchInfo.type" style="width: 200px" clearable filterable placeholder="请选择">
                <el-option v-for="item in typeList" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="搜索">
              <el-input v-model="searchInfo.keyword" placeholder="关键词" style="width: 200px" clearable />
            </el-form-item>
          </div>
          <el-form-item>
            <el-button class="button-gap" type="primary"  @click="onSearch">
              查询
            </el-button>
            <el-button class="button-gap" type="info"  @click="onReset">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <!-- 条件选择 -->

    </div>

    <!-- 表格栏 -->
    <div class="gva-table-box">
      <!-- 表格操作按钮 -->
      <div style="margin-bottom: 16px">
        <el-button type="primary"  @click="addFunc('add')">
          新增
        </el-button>
        <el-button type="grey" :disabled="!selectionList.length" @click="openDeleteDialog('batch')">
          删除
        </el-button>
      </div>
      <!-- 表格 -->
      <div v-if="tableData.length > 0">
        <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
          <el-table-column type="selection" width="55" />
          <el-table-column align="left" label="序号" type="index" min-width="100" />
          <el-table-column align="left" label="配置项" min-width="200" prop="type" sortable="custom">
            <template #default="{ row }">
              <div class="text-overflow">{{ row.type }}</div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="内容" min-width="200" prop="name">
            <template #default="{ row }">
              <div class="text-overflow">{{ row.name }}</div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="单位" min-width="200" prop="name">
            <template #default="{ row }">
              <div class="text-overflow">{{ row.unit ? row.unit : '--' }}</div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="操作" fixed="right" width="300">
            <template #default="scope">
              <el-button  class="text-blue-500"  type="text" link @click="editFunc(scope.row)">
                编辑
              </el-button>
              <el-button class="text-blue-500" type="text" link @click="openDeleteDialog('single',scope.row)">
                删除
              </el-button>
            </template>
          </el-table-column>


        </el-table>
      </div>
      <div v-else class="text-center">
        <h3>暂无数据</h3>
      </div>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageNum" :page-size="searchInfo.pageSize" :page-sizes="[10, 20, 30, 50, 100]" :total="searchInfo.total"
            layout="total,prev, pager, next,  sizes, jumper" @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 抽屉：增加/编辑 -->
    <el-dialog :show-close="false" v-model="drawerFormVisible" width="30%" :before-close="closeDrawer">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ drawerTitle }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="drawerFormVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form
          label-position="left" ref="drawerForm" :model="row" :rules="rules" label-width="80px">
        <el-form-item label="配置项" prop="type">
          <el-select v-model="row.type" placeholder="请选择" clearable>
            <el-option v-for="item in typeList" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="name">
          <el-input v-model="row.name" autocomplete="off" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item v-if="row.type === 17" label="单位" prop="unit">
          <el-input v-model="row.unit" autocomplete="off" placeholder="请输入（可选）" clearable />
        </el-form-item>
      </el-form>
      <div class="flex justify-end mt-4">
        <el-button type="info" @click="closeDrawer">
          取 消
        </el-button>
        <el-button type="primary" @click="confirmForm">
          确 定
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
        v-model="deleteVisible"
        :show-close="false"
        width="520px"
        @close="handleDialogClose"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);display: flex;align-items: center;">
          <el-icon :style="{ color: 'rgba(255, 125, 0, 1)', fontSize: '1.25em',marginRight: '8px' }">
            <WarningFilled />
          </el-icon>删除</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="deleteVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <span>此操作将永远删除该项，是否继续？</span>
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

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted } from 'vue'
import {
  getCategoryList,
  addCategory,
  updateCategory, deleteCategory, deleteCategoryByIds
} from "@/api/yunguan/system/type"

const dialogVisible = ref(false)
const deleteVisible =ref(false)
const deleteType = ref(''); // 删除类型：single(单删) / batch(批量删)
const currentProductId = ref(null); // 单删时的产品ID
const selectedIds = ref([]); // 批量删除时的ID数组

const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};

// 打开删除弹窗（统一入口）
const openDeleteDialog = (type, row = null) => {
  deleteType.value = type;

  if (type === 'single' && row) {
    currentProductId.value = row.id;
    console.log(row.id,'sfa')
  } else if (type === 'batch') {
    selectedIds.value = selectionList.value.map(item => item.id);
  }

  deleteVisible.value = true;
};

const confirmDelete = async () => {
  try {
    let resp;

    if (deleteType.value === 'single') {
      resp = await deleteCategory({ id: currentProductId.value });
    } else if (deleteType.value === 'batch') {
      resp = await deleteCategoryByIds({ ids: selectedIds.value });
    }

    // 严格判断接口成功状态
    if (resp.code === 0) {
      // 显示对应成功消息
      ElMessage({
        type: 'success',
        message: deleteType.value === 'single'
            ? '删除成功'
            : `批量删除 ${selectedIds.value.length} 条记录成功`
      });

      // 执行成功后的操作
      await getTableData(); // 刷新表格
      emits('confirm');     // 通知父组件
    } else {
      // 接口返回失败状态时显示错误
      ElMessage({
        showClose: true,
        message: resp.msg || '删除失败',
        type: 'error',
      });
      return; // 失败时直接返回，不执行后续操作
    }
  } catch (error) {

    console.error('删除过程中出现异常:', error);
  } finally {
    // 清理状态（无论成功失败都执行）
    resetDeleteState();
  }
};
const resetDeleteState = () => {
  deleteVisible.value = false;
  deleteType.value = '';
  currentProductId.value = null;
  selectedIds.value = [];
};

/*********
 * 条件筛选
 *********/
// 搜索条件
const searchInfo = ref({
  pageNum: 1,
  pageSize: 10,
  total: 0,
  orderBy: null,
  isDesc: null,
  keyword: null,
  type: null,
})

// 类型列表
const typeList = ref([
  { label: '计算类型', value: 10 },
  { label: '卡类型', value: 14 },
  { label: '品牌', value: 11 },
  { label: '型号', value: 15 },
  { label: '供电方式', value: 12 },
  { label: '运营商', value: 13 },
  { label: '路由类型', value: 16 },
  { label: '供需规格', value: 17}
])


// 搜索按钮
const onSearch = () => {
  getTableData()
}

// 重置按钮
const onReset = () => {
  resetSearchInfo()
  getTableData()
}

// 重置搜索条件
const resetSearchInfo = () => {
  searchInfo.value = {
    pageNum: 1,
    pageSize: 10,
    total: 0,
    orderBy: null,
    isDesc: null,
    keyword: null,
  }
}

// 搜索获取表格数据
const getTableData = async () => {
  await getCategoryList(searchInfo.value).then((response) => {
    if (response.code === 0) {
      tableData.value = response.data.list ? response.data.list : []
      searchInfo.value.total = response.data.total
    }
  })
}

/********
 * 表格栏
 ********/
// 表格数据
const tableData = ref([])
// 多选
const selectionList = ref([])

// 多选改变
const handleSelectionChange = (val) => {
  selectionList.value = val
}

// 分页
const handleCurrentChange = (val) => {
  searchInfo.value.pageNum = val
  getTableData()
}
const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val
  getTableData()
}

// 排序
const sortChange = ({ prop, order }) => {
  if (order) {
    searchInfo.value.orderBy = prop
    searchInfo.value.isDesc = order === 'descending'
  } else {
    searchInfo.value.orderBy = null
    searchInfo.value.isDesc = null
  }
  getTableData()
}

// 单个删除
const deleteFunc = async (row) => {
  ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const resp = await deleteCategory({id:row.id})
    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
    } else {
      ElMessage({
        showClose: true,
        message: '删除失败',
        type: 'error',
      })
    }
    await getTableData()
  })
}

// 批量删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = selectionList.value.map(item => item.id)
    const resp = await deleteCategoryByIds({ ids: ids })
    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '批量删除成功'
      })
      if (tableData.value.length === ids.length && searchInfo.value.pageNum > 1) {
        searchInfo.value.pageNum.value--
      }
    } else {
      ElMessage({
        showClose: true,
        message: '批量删除失败',
        type: 'error',
      })
    }
    await getTableData()
  })
}

/******
 * 抽屉
 ******/
const rules = ref({
  type: [{ required: true, message: "请选择配置项", trigger: "change" }],
  name: [{ required: true, message: "请输入内容", trigger: "blur" }],
  // unit: [{ required: true, message: "请输入单位", trigger: "blur" }],
})

const drawerTitle = ref(null)
const drawerForm = ref(null)
const drawerFormVisible = ref(false)
const drawerAuditVisible = ref(false)
const row = ref(null)

const addFunc = function () {
  row.value = {}
  openDrawer('add')
}

const editFunc = function (r) {
  row.value = { ...r }
  row.value.type = typeList.value
      .filter(item => item.label === r.type)
      .map(item => item.value)
  row.value.type = row.value.type.length > 0 ? row.value.type[0] : null
  openDrawer('edit')
}

const openDrawer = (key) => {
  switch (key) {
    case 'add':
      drawerTitle.value = '新增'
      drawerFormVisible.value = true
      break
    case 'edit':
      drawerTitle.value = '编辑'
      drawerFormVisible.value = true
      break
    default:
      break
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  drawerAuditVisible.value = false
  drawerForm.value.resetFields()
}


// 新增/编辑确认按钮
const confirmForm = async () => {
  drawerForm.value.validate(async (valid) => {
    if (valid) {
      if (drawerTitle.value === '新增') {
        const resp = await addCategory({
          type: row.value.type,
          name: row.value.name,
          unit: row.value.unit
        })
        if (resp.code === 0) {
          ElMessage({
            type: 'success',
            message: '添加成功',
            showClose: true
          })
          closeDrawer()
          await getTableData()
        }
      }
      if (drawerTitle.value === '编辑') {
        const resp = await updateCategory({
          id: row.value.id,
          type: row.value.type,
          name: row.value.name,
          unit: row.value.unit
        })
        if (resp.code === 0) {
          ElMessage({
            type: 'success',
            message: '编辑成功',
            showClose: true
          })
          closeDrawer()
          await getTableData()
        }
      }
    }
  })
}

/********
 * 初始化
 ********/
onMounted(() => {
  getTableData()
})
</script>


<style scoped lang="scss">
.warning {
  color: #dc143c;
}

// 设置筛选标签的下边距
.tabs {
  margin-bottom: 30px;
}

// 设置按钮间隔且保证挤压时上下对齐
.button-gap {
  margin-left: 0;

}

// 文本溢出表格后截断
.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
