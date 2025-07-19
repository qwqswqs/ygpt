<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
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



      <!-- 由于此处菜单跟左侧列表一一对应所以不需要分页 pageSize默认999 -->
      <el-table
        :data="tableData"
        row-key="ID"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
        <el-table-column
          align="left"
          label="ID"
          min-width="100"
          prop="ID"
        />
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
        <el-table-column
          align="left"
          label="路由Name"
          show-overflow-tooltip
          min-width="160"
          prop="name"
        />
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
        <el-table-column
          align="left"
          label="父节点"
          min-width="90"
          prop="parentId"
        />
        <el-table-column
          align="left"
          label="排序"
          min-width="70"
          prop="sort"
        />
        <el-table-column
          align="left"
          label="文件路径"
          min-width="360"
          prop="component"
        />
        <el-table-column
          align="left"
          fixed="right"
          label="操作"
          width="250"
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
              @click="deleteMenu(scope.row.ID)"
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
              label="文件路径"
              prop="component"
              style="width: 680px;margin-right: 40px"
          >
            <el-input
                v-model="form.component"
                filterable
                allow-create
                autocomplete="off"
                style="width: 680px"
                placeholder="请输入"
                default-first-option
                @change="fmtComponent"
            >
            </el-input>
            <span style="font-size: 14px; ">如果菜单包含子菜单，请创建router-view二级路由页面或者 </span><el-button
              style="font-size: 14px"
              @click="form.component = 'view/routerHolder.vue'"
              type="text"
              link
          >
            点我设置
          </el-button>
          </el-form-item>
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
        </div>

        <div style="display: flex">
          <el-form-item
              label="路由Name"
              prop="path"
              style="width: 320px;margin-right: 40px"
          >
            <el-input
                v-model="form.name"
                autocomplete="off"
                placeholder="请输入"
                @change="changeName"
                style="width: 320px"
            />
          </el-form-item>
          <el-form-item
              prop="path"
              style="width: 320px;margin-right: 40px"
          >
            <template #label>
                <span style="display: inline-flex; align-items: center">
                  <span>路由Path</span>
                  <!--                  <el-checkbox-->
                  <!--                    v-model="checkFlag"-->
                  <!--                    style="margin-left: 12px; height: auto"-->
                  <!--                  >添加参数</el-checkbox>-->
                </span>
            </template>

            <el-input
                v-model="form.path"
                :disabled="!checkFlag"
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
        </div>

        <div style="display: flex">
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

        <div style="display: flex">
          <el-form-item
              prop="meta.activeName"
              style="width: 320px;margin-right: 40px"

          >
            <template #label>
              <div>
                <span> 高亮菜单 </span>
                <el-tooltip
                    content="注：当到达此路由时候，指定左侧菜单指定name会处于活跃状态（亮起），可为空，为空则为本路由Name。"
                    placement="top"
                    effect="light"

                >
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input
                v-model="form.meta.activeName"
                :placeholder="form.name || '请输入'"
                autocomplete="off"
                style="width: 320px;"
            />
          </el-form-item>
          <el-form-item style="width: 320px;margin-right: 40px" label="KeepAlive" prop="meta.keepAlive">
            <el-radio-group v-model="form.meta.keepAlive">
              <el-radio :label="true">是</el-radio>
              <el-radio :label="false">否</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item style="width: 320px" label="CloseTab" prop="meta.closeTab">
            <!-- 使用el-radio-group包裹el-radio，并通过v-model绑定值 -->
            <el-radio-group v-model="form.meta.closeTab">
              <el-radio :label="true">是</el-radio>
              <el-radio :label="false">否</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>

        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item>
              <template #label>
                <div>
                  <span> 是否为基础页面 </span>
                  <el-tooltip
                      content="此项选择为是，则不会展示左侧菜单以及顶部信息。"
                      placement="top"
                      effect="light"
                  >
                    <el-icon><QuestionFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>

              <!-- 使用 el-radio-group 实现单选 -->
              <el-radio-group v-model="form.meta.defaultMenu">
                <el-radio :label="true">是</el-radio>
                <el-radio :label="false">否</el-radio>
              </el-radio-group>
            </el-form-item>

          </el-col>
        </el-row>
      </el-form>
      <div>
        <div class="flex items-center gap-2">
          <el-button
            type="primary"

            @click="addParameter(form)"
          >
            新增菜单参数
          </el-button>
        </div>
        <el-table
          :data="form.parameters"
          style="width: 100%; margin-top: 12px"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
        >
          <el-table-column
            align="left"
            prop="type"
            label="参数类型"
            width="336"
          >
            <template #default="scope">
              <el-select
                v-model="scope.row.type"
                placeholder="请选择"
                style="width: 200px"
              >
                <el-option
                  key="query"
                  value="query"
                  label="query"
                />
                <el-option
                  key="params"
                  value="params"
                  label="params"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="key"
            label="参数key"
            width="332"
          >
            <template #default="scope">
              <div>
                <el-input placeholder="请输入" style="width: 200px" v-model="scope.row.key" />
              </div>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="value"
            label="参数值"
            width="322"
          >
            <template #default="scope">
              <div>
                <el-input placeholder="请输入" style="width: 280px" v-model="scope.row.value" />
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" align="left" width="130">
            <template #default="scope">
              <div>
                <el-button
                  type="text"
                  link
                  @click="deleteParameter(form.parameters, scope.$index)"
                >
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <div class="flex items-center gap-2 mt-3">
          <el-button
            type="primary"
            @click="addBtn(form)"
          >
            新增可控按钮
          </el-button>
          <el-icon
            class="cursor-pointer"
            @click="
              toDoc('https://www.gin-vue-admin.com/guide/web/button-auth.html')
            "
          >
            <QuestionFilled />
          </el-icon>
        </div>

        <el-table
          :data="form.menuBtn"
          style="width: 100%; margin-top: 12px"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
        >
          <el-table-column
            align="left"
            prop="name"
            label="按钮名称"
            width="336"
          >
            <template #default="scope">
              <div>
                <el-input placeholder="请输入" style="width: 200px" v-model="scope.row.name" />
              </div>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            prop="name"
            label="备注"
            width="308"
          >
            <template #default="scope">
              <div>
                <el-input placeholder="请输入" style="width: 200px" v-model="scope.row.desc" />
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left">
            <template #default="scope">
              <div>
                <el-button
                  type="text"
                  link
                  @click="deleteBtn(form.menuBtn, scope.$index)"
                >
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="flex" style="justify-content: right;margin-top: 108px">
        <el-button type="info" @click="closeDialog">取 消</el-button>
        <el-button
            type="primary"
            @click="enterDialog"
        >确 定</el-button>
      </div>
    </el-drawer>
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

defineOptions({
  name: 'Menus',
})

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
