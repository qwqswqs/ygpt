<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div >
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-button
                type="primary"

                @click="addMenu(0)"
            >
              新增根菜单
            </el-button>
          </el-form-item>
        </el-form>

      </div>

      <!-- 由于此处菜单跟左侧列表一一对应所以不需要分页 pageSize默认999 -->
      <el-table
        :data="tableData"
        row-key="ID"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
<!--        <el-table-column-->
<!--          align="left"-->
<!--          label="ID"-->
<!--          min-width="100"-->
<!--          prop="ID"-->
<!--        />-->
        <el-table-column
          align="left"
          label="展示名称"
          min-width="120"
          prop="authorityName"
        >
          <template #default="scope">
            <span>{{ scope.row.meta.title }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="图标"
          min-width="140"
          prop="authorityName"
        >
          <template #default="scope">
            <div
              v-if="scope.row.meta.icon"
              class="icon-column"
            >
              <el-icon>
                <component :is="scope.row.meta.icon" />
              </el-icon>
              <span>{{ scope.row.meta.icon }}</span>
            </div>
          </template>
        </el-table-column>
<!--        <el-table-column-->
<!--          align="left"-->
<!--          label="路由Name"-->
<!--          show-overflow-tooltip-->
<!--          min-width="160"-->
<!--          prop="name"-->
<!--        />-->
<!--        <el-table-column-->
<!--          align="left"-->
<!--          label="路由Path"-->
<!--          show-overflow-tooltip-->
<!--          min-width="160"-->
<!--          prop="path"-->
<!--        />-->
        <el-table-column
          align="left"
          label="是否隐藏"
          min-width="100"
          prop="hidden"
        >
          <template #default="scope">
            <span>{{ scope.row.hidden ? "隐藏" : "显示" }}</span>
          </template>
        </el-table-column>
<!--        <el-table-column-->
<!--          align="left"-->
<!--          label="父节点"-->
<!--          min-width="90"-->
<!--          prop="parentId"-->
<!--        />-->
        <el-table-column
          align="left"
          label="排序"
          min-width="70"
          prop="sort"
        />
<!--        <el-table-column-->
<!--          align="left"-->
<!--          label="文件路径"-->
<!--          min-width="360"-->
<!--          prop="component"-->
<!--        />-->
        <el-table-column
          align="left"
          fixed="right"
          label="操作"
          width="300"
        >
          <template #default="scope">
            <el-button
                type="text" class="text-blue-500"
              link
              @click="addMenu(scope.row.ID)"
            >
              添加子菜单
            </el-button>
            <el-button
                type="text" class="text-blue-500"
              link
              @click="editMenu(scope.row.ID)"
            >
              编辑
            </el-button>
            <el-button
                type="text" class="text-blue-500"
              link
              @click="openDeleteDialog(scope.row.ID)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      size="1200px"
      :before-close="handleClose"
      :show-close="false"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ dialogTitle }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="closeDialog"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <warning-bar title="新增菜单，需要在角色管理内配置权限才可使用" />
      <el-form
          v-if="dialogFormVisible"
          ref="menuForm"
          :inline="true"
          :model="form"
          :rules="rules"
          label-position="top"
      >
        <div style="display: flex">
          <el-form-item
              label="展示名称"
              prop="meta.title"
              style="width: 320px"
          >
            <el-input
                v-model="form.meta.title"
                autocomplete="off"
                placeholder="请输入"
                style="width: 320px"
            />
          </el-form-item>
          <el-form-item label="是否隐藏" style="width: 320px">
            <el-radio-group v-model="form.hidden">
              <el-radio :label="true">是</el-radio>
              <el-radio :label="false">否</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item
              label="父节点ID"
              style="width: 320px;margin-right: 40px"
          >
            <el-cascader
                v-model="form.parentId"
                style="width: 320px"
                :disabled="!isEdit"
                :options="menuOption"
                :props="{
                  checkStrictly: true,
                  label: 'title',
                  value: 'ID',
                  disabled: 'disabled',
                  emitPath: false,
                }"
                :show-all-levels="false"
                filterable
            />
          </el-form-item>
        </div>

        <div style="display: flex">

          <el-form-item
              label="图标"
              prop="meta.icon"
              style="width: 320px;margin-right: 40px"
          >
            <icon
                v-model="form.meta.icon"
            />
          </el-form-item>
          <el-form-item
              label="排序标记"
              prop="sort"
              style="width: 320px"
          >
            <el-input
                v-model.number="form.sort"
                autocomplete="off"
                placeholder="请输入"
                style="width: 320px"
            />
          </el-form-item>
        </div>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="closeDialog">取 消</el-button>
        <el-button
            type="primary"
            @click="enterDialog"
        >确 定</el-button>
      </div>
    </el-drawer>
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
import {
  updateBaseMenu,
  getMenuList,
  addBaseMenu,
  deleteBaseMenu,
  getBaseMenuById,
} from '@/api/system/menu'
import icon from '@/view/superAdmin/menu/icon.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { canRemoveAuthorityBtnApi } from '@/api/system/authorityBtn'
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import pathInfo from '@/pathInfo.json'

import { toDoc } from '@/utils/doc'
import { toLowerCase } from '@/utils/stringFun'
import {deleteAuthority} from "@/api/system/authority";

defineOptions({
  name: 'Menus',
})
const deleteVisible =ref(false)
const currentProductId = ref(null);
const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};
// 打开删除弹窗（统一入口）
const openDeleteDialog = ( ID) => {
  currentProductId.value = ID;
  console.log(ID,'id')
  deleteVisible.value = true;
};
const confirmDelete = async () => {
  try {
    let resp;
    resp = await deleteBaseMenu({ ID: currentProductId.value });

    // 严格判断接口成功状态
    if (resp.code === 0) {
      // 显示对应成功消息
      ElMessage({
        type: 'success',
        message: '删除成功'
      });

      // 执行成功后的操作
      await getTableData(); // 刷新表格
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
  }
};
const pathOptions = reactive({})

onMounted(()=>{
  for (let pathInfoKey in pathInfo) {
    // 取消掉最前面的 /src/
    pathOptions[pathInfoKey.replace(/^\/src\//, '')] = pathInfo[pathInfoKey]
  }
})

const rules = reactive({
  path: [{ required: true, message: '请输入菜单name', trigger: 'blur' }],
  component: [{ required: true, message: '请输入文件路径', trigger: 'blur' }],
  'meta.title': [
    { required: true, message: '请输入菜单展示名称', trigger: 'blur' },
  ],
})

const tableData = ref([])
// 查询
const getTableData = async() => {
  const table = await getMenuList()
  if (table.code === 0) {
    tableData.value = table.data
  }
}

getTableData()

// 新增参数
const addParameter = (form) => {
  if (!form.parameters) {
    form.parameters = []
  }
  form.parameters.push({
    type: 'query',
    key: '',
    value: '',
  })
}

const fmtComponent = () => {
  form.value.component = form.value.component.replace(/\\/g, '/')
  form.value.name = toLowerCase(pathOptions[form.value.component])
  form.value.path = form.value.name
}

// 删除参数
const deleteParameter = (parameters, index) => {
  parameters.splice(index, 1)
}

// 新增可控按钮
const addBtn = (form) => {
  if (!form.menuBtn) {
    form.menuBtn = []
  }
  form.menuBtn.push({
    name: '',
    desc: '',
  })
}
// 删除可控按钮
const deleteBtn = async(btns, index) => {
  const btn = btns[index]
  if (btn.ID === 0) {
    btns.splice(index, 1)
    return
  }
  const res = await canRemoveAuthorityBtnApi({ id: btn.ID })
  if (res.code === 0) {
    btns.splice(index, 1)
  }
}

const form = ref({
  ID: 0,
  path: '',
  name: '',
  hidden: false,
  parentId: 0,
  component: '',
  meta: {
    activeName: '',
    title: '',
    icon: '',
    defaultMenu: false,
    closeTab: false,
    keepAlive: false,
  },
  parameters: [],
  menuBtn: [],
})
const changeName = () => {
  form.value.path = form.value.name
}

const handleClose = (done) => {
  initForm()
  done()
}
// 删除菜单
const deleteMenu = (ID) => {
  ElMessageBox.confirm('此操作将永久删除所有角色下该菜单, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async() => {
      const res = await deleteBaseMenu({ ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!',
        })

        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除',
      })
    })
}
// 初始化弹窗内表格方法
const menuForm = ref(null)
const checkFlag = ref(false)
const initForm = () => {
  checkFlag.value = false
  menuForm.value.resetFields()
  form.value = {
    ID: 0,
    path: '',
    name: '',
    hidden: false,
    parentId: 0,
    component: '',
    meta: {
      title: '',
      icon: '',
      defaultMenu: false,
      closeTab: false,
      keepAlive: false,
    },
  }
}
// 关闭弹窗

const dialogFormVisible = ref(false)
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}
// 添加menu
const enterDialog = async() => {
  menuForm.value.validate(async(valid) => {
    if (valid) {
      let res
      if (isEdit.value) {
        res = await updateBaseMenu(form.value)
      } else {
        res = await addBaseMenu(form.value)
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: isEdit.value ? '编辑成功' : '添加成功!',
        })
        getTableData()
      }
      initForm()
      dialogFormVisible.value = false
    }
  })
}

const menuOption = ref([
  {
    ID: '0',
    title: '根菜单',
  },
])
const setOptions = () => {
  menuOption.value = [
    {
      ID: 0,
      title: '根目录',
    },
  ]
  setMenuOptions(tableData.value, menuOption.value, false)
}
const setMenuOptions = (menuData, optionsData, disabled) => {
  menuData &&
    menuData.forEach((item) => {
      if (item.children && item.children.length) {
        const option = {
          title: item.meta.title,
          ID: item.ID,
          disabled: disabled || item.ID === form.value.ID,
          children: [],
        }
        setMenuOptions(
          item.children,
          option.children,
          disabled || item.ID === form.value.ID
        )
        optionsData.push(option)
      } else {
        const option = {
          title: item.meta.title,
          ID: item.ID,
          disabled: disabled || item.ID === form.value.ID,
        }
        optionsData.push(option)
      }
    })
}

// 添加菜单方法，id为 0则为添加根菜单
const isEdit = ref(false)
const dialogTitle = ref('新增菜单')
const addMenu = (id) => {
  dialogTitle.value = '新增菜单'
  form.value.parentId = id
  isEdit.value = false
  setOptions()
  dialogFormVisible.value = true
}
// 修改菜单方法
const editMenu = async(id) => {
  dialogTitle.value = '编辑菜单'
  const res = await getBaseMenuById({ id })
  form.value = res.data.menu
  isEdit.value = true
  setOptions()
  dialogFormVisible.value = true
}
</script>

<style scoped lang="scss">
.warning {
  color: #dc143c;
}
.icon-column {
  display: flex;
  align-items: center;
  .el-icon {
    margin-right: 8px;
  }
}
</style>
