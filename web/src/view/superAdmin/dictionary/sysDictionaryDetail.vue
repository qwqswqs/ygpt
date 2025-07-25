<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list justify-between">
        <span class="text font-bold">字典详细内容</span>
        <el-button
          type="primary"

          @click="openDrawer"
        >
          新增字典项
        </el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          align="left"
          label="日期"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="展示值"
          prop="label"
        />

        <el-table-column
          align="left"
          label="字典值"
          prop="value"
        />

        <el-table-column
          align="left"
          label="扩展值"
          prop="extend"
        />

        <el-table-column
          align="left"
          label="启用状态"
          prop="status"
          width="120"
        >
          <template #default="scope">
            {{ formatBoolean(scope.row.status) }}
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="排序标记"
          prop="sort"
          width="120"
        />

        <el-table-column
          align="left"
          label="操作"
          width="180"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link

              @click="updateSysDictionaryDetailFunc(scope.row)"
            >
              变更
            </el-button>
            <el-button
              type="primary"
              link

              @click="deleteSysDictionaryDetailFunc(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total,prev, pager, next,  sizes, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      size="30%"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create' ? '添加字典项' : '修改字典项' }}</span>
          <div>
            <el-button @click="closeDrawer">
              取 消
            </el-button>
            <el-button type="primary" @click="enterDrawer">
              确 定
            </el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="drawerForm"
        :model="formData"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item
          label="展示值"
          prop="label"
        >
          <el-input
            v-model="formData.label"
            placeholder="请输入展示值"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="字典值"
          prop="value"
        >
          <el-input
            v-model="formData.value"
            placeholder="请输入字典值"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="扩展值"
          prop="extend"
        >
          <el-input
            v-model="formData.extend"
            placeholder="请输入扩展值"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="启用状态"
          prop="status"
          required
        >
          <el-switch
            v-model="formData.status"
            active-text="开启"
            inactive-text="停用"
          />
        </el-form-item>
        <el-form-item
          label="排序标记"
          prop="sort"
        >
          <el-input-number
            v-model.number="formData.sort"
            placeholder="排序标记"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysDictionaryDetail,
  deleteSysDictionaryDetail,
  updateSysDictionaryDetail,
  findSysDictionaryDetail,
  getSysDictionaryDetailList
} from '@/api/system/sysDictionaryDetail' // 此处请自行替换地址
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatBoolean, formatDate } from '@/utils/format'

defineOptions({
  name: 'SysDictionaryDetail'
})

const props = defineProps({
  sysDictionaryID: {
    type: Number,
    default: 0
  }
})

const formData = ref({
  label: null,
  value: null,
  status: true,
  sort: null
})
const rules = ref({
  label: [
    {
      required: true,
      message: '请输入展示值',
      trigger: 'blur'
    }
  ],
  value: [
    {
      required: true,
      message: '请输入字典值',
      trigger: 'blur'
    }
  ],
  sort: [
    {
      required: true,
      message: '排序标记',
      trigger: 'blur'
    }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  if(!props.sysDictionaryID) return
  const table = await getSysDictionaryDetailList({
    page: page.value,
    pageSize: pageSize.value,
    sysDictionaryID: props.sysDictionaryID
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const type = ref('')
const drawerFormVisible = ref(false)
const updateSysDictionaryDetailFunc = async(row) => {
  drawerForm.value && drawerForm.value.clearValidate()
  const res = await findSysDictionaryDetail({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reSysDictionaryDetail
    drawerFormVisible.value = true
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  formData.value = {
    label: null,
    value: null,
    status: true,
    sort: null,
    sysDictionaryID: props.sysDictionaryID
  }
}
const deleteSysDictionaryDetailFunc = async(row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await deleteSysDictionaryDetail({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const drawerForm = ref(null)
const enterDrawer = async() => {
  drawerForm.value.validate(async valid => {
    formData.value.sysDictionaryID = props.sysDictionaryID
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createSysDictionaryDetail(formData.value)
        break
      case 'update':
        res = await updateSysDictionaryDetail(formData.value)
        break
      default:
        res = await createSysDictionaryDetail(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDrawer()
      getTableData()
    }
  })
}
const openDrawer = () => {
  type.value = 'create'
  drawerForm.value && drawerForm.value.clearValidate()
  drawerFormVisible.value = true
}

watch(() => props.sysDictionaryID, () => {
  getTableData()
})

</script>

<style>
</style>
