<template>
  <div>
    <div class="gva-table-box">
<!--      <div class="gva-btn-list">-->
<!--        <el-button-->
<!--          type="primary"-->
<!--          -->
<!--          @click="addMenu(0)"-->
<!--        >-->
<!--          新增根菜单-->
<!--        </el-button>-->
<!--      </div>-->

      <!-- 由于此处菜单跟左侧列表一一对应所以不需要分页 pageSize默认999 -->
      <el-table
        :data="tableData"
        row-key="ID"
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
<!--            <el-button-->
<!--              type="primary"-->
<!--              link-->
<!--              -->
<!--              @click="addMenu(scope.row.ID)"-->
<!--            >-->
<!--              添加子菜单-->
<!--            </el-button>-->
            <el-button
              type="primary"
              link

              @click="editMenu(scope.row.ID)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
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
      size="60%"
      :before-close="handleClose"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">
              取 消
            </el-button>
            <el-button
              type="primary"
              @click="enterDialog"
            >
              确 定
            </el-button>
          </div>
        </div>
      </template>

      <warning-bar title="新增菜单，需要在角色管理内配置权限才可使用" />
      <el-form
        v-if="dialogFormVisible"
        ref="menuForm"
        :inline="true"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <el-row class="w-full">
<!--          <el-col :span="16">-->
<!--            <el-form-item-->
<!--              label="文件路径"-->
<!--              prop="component"-->
<!--            >-->
<!--              <el-input-->
<!--                v-model="form.component"-->
<!--                filterable-->
<!--                allow-create-->
<!--                autocomplete="off"-->
<!--                style="width: 100%"-->
<!--                placeholder="页面:view/xxx/xx.vue 插件:plugin/xx/xx.vue"-->
<!--                default-first-option-->
<!--                @change="fmtComponent"-->
<!--              >-->
<!--              </el-input>-->
<!--              <span style="font-size: 12px; margin-right: 12px">如果菜单包含子菜单，请创建router-view二级路由页面或者</span><el-button-->
<!--                style="margin-top: 4px"-->
<!--                @click="form.component = 'view/routerHolder.vue'"-->
<!--              >-->
<!--                点我设置-->
<!--              </el-button>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
          <el-col :span="8">
            <el-form-item
              label="展示名称"
              prop="meta.title"
            >
              <el-input
                v-model="form.meta.title"
                autocomplete="off"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
                label="是否隐藏"
            >
              <el-select
                  v-model="form.hidden"
                  style="width: 100%"
                  placeholder="是否在列表隐藏"
              >
                <el-option
                    :value="false"
                    label="否"
                />
                <el-option
                    :value="true"
                    label="是"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
<!--        <el-row class="w-full">-->
<!--&lt;!&ndash;          <el-col :span="8">&ndash;&gt;-->
<!--&lt;!&ndash;            <el-form-item&ndash;&gt;-->
<!--&lt;!&ndash;              label="路由Name"&ndash;&gt;-->
<!--&lt;!&ndash;              prop="path"&ndash;&gt;-->
<!--&lt;!&ndash;            >&ndash;&gt;-->
<!--&lt;!&ndash;              <el-input&ndash;&gt;-->
<!--&lt;!&ndash;                v-model="form.name"&ndash;&gt;-->
<!--&lt;!&ndash;                autocomplete="off"&ndash;&gt;-->
<!--&lt;!&ndash;                placeholder="唯一英文字符串"&ndash;&gt;-->
<!--&lt;!&ndash;                @change="changeName"&ndash;&gt;-->
<!--&lt;!&ndash;              />&ndash;&gt;-->
<!--&lt;!&ndash;            </el-form-item>&ndash;&gt;-->
<!--&lt;!&ndash;          </el-col>&ndash;&gt;-->
<!--&lt;!&ndash;          <el-col :span="8">&ndash;&gt;-->
<!--&lt;!&ndash;            <el-form-item&ndash;&gt;-->
<!--&lt;!&ndash;              prop="path"&ndash;&gt;-->
<!--&lt;!&ndash;            >&ndash;&gt;-->
<!--&lt;!&ndash;              <template #label>&ndash;&gt;-->
<!--&lt;!&ndash;                <span style="display: inline-flex; align-items: center">&ndash;&gt;-->
<!--&lt;!&ndash;                  <span>路由Path</span>&ndash;&gt;-->
<!--&lt;!&ndash;                  <el-checkbox&ndash;&gt;-->
<!--&lt;!&ndash;                    v-model="checkFlag"&ndash;&gt;-->
<!--&lt;!&ndash;                    style="margin-left: 12px; height: auto"&ndash;&gt;-->
<!--&lt;!&ndash;                  >添加参数</el-checkbox>&ndash;&gt;-->
<!--&lt;!&ndash;                </span>&ndash;&gt;-->
<!--&lt;!&ndash;              </template>&ndash;&gt;-->

<!--&lt;!&ndash;              <el-input&ndash;&gt;-->
<!--&lt;!&ndash;                v-model="form.path"&ndash;&gt;-->
<!--&lt;!&ndash;                :disabled="!checkFlag"&ndash;&gt;-->
<!--&lt;!&ndash;                autocomplete="off"&ndash;&gt;-->
<!--&lt;!&ndash;                placeholder="建议只在后方拼接参数"&ndash;&gt;-->
<!--&lt;!&ndash;              />&ndash;&gt;-->
<!--&lt;!&ndash;            </el-form-item>&ndash;&gt;-->
<!--&lt;!&ndash;          </el-col>&ndash;&gt;-->
<!--         -->
<!--        </el-row>-->
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item
              label="父节点ID"
            >
              <el-cascader
                v-model="form.parentId"
                style="width: 100%"
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
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="图标"
              prop="meta.icon"
            >
              <icon
                v-model="form.meta.icon"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="排序标记"
              prop="sort"
            >
              <el-input
                v-model.number="form.sort"
                autocomplete="off"
              />
            </el-form-item>
          </el-col>
        </el-row>
<!--        <el-row class="w-full">-->
<!--          <el-col :span="8">-->
<!--            <el-form-item-->
<!--              prop="meta.activeName"-->
<!--            >-->
<!--              <template #label>-->
<!--                <div>-->
<!--                  <span> 高亮菜单 </span>-->
<!--                  <el-tooltip-->
<!--                    content="注：当到达此路由时候，指定左侧菜单指定name会处于活跃状态（亮起），可为空，为空则为本路由Name。"-->
<!--                    placement="top"-->
<!--                    effect="light"-->
<!--                  >-->
<!--                    <el-icon><QuestionFilled /></el-icon>-->
<!--                  </el-tooltip>-->
<!--                </div>-->
<!--              </template>-->
<!--              <el-input-->
<!--                v-model="form.meta.activeName"-->
<!--                :placeholder="form.name"-->
<!--                autocomplete="off"-->
<!--              />-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--          <el-col :span="8">-->
<!--            <el-form-item-->
<!--              label="KeepAlive"-->
<!--              prop="meta.keepAlive"-->
<!--            >-->
<!--              <el-select-->
<!--                v-model="form.meta.keepAlive"-->
<!--                style="width: 100%"-->
<!--                placeholder="是否keepAlive缓存页面"-->
<!--              >-->
<!--                <el-option-->
<!--                  :value="false"-->
<!--                  label="否"-->
<!--                />-->
<!--                <el-option-->
<!--                  :value="true"-->
<!--                  label="是"-->
<!--                />-->
<!--              </el-select>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--          <el-col :span="8">-->
<!--            <el-form-item-->
<!--              label="CloseTab"-->
<!--              prop="meta.closeTab"-->
<!--            >-->
<!--              <el-select-->
<!--                v-model="form.meta.closeTab"-->
<!--                style="width: 100%"-->
<!--                placeholder="是否自动关闭tab"-->
<!--              >-->
<!--                <el-option-->
<!--                  :value="false"-->
<!--                  label="否"-->
<!--                />-->
<!--                <el-option-->
<!--                  :value="true"-->
<!--                  label="是"-->
<!--                />-->
<!--              </el-select>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--        </el-row>-->
<!--        <el-row class="w-full">-->
<!--          <el-col :span="8">-->
<!--            <el-form-item>-->
<!--              <template #label>-->
<!--                <div>-->
<!--                  <span> 是否为基础页面 </span>-->
<!--                  <el-tooltip-->
<!--                    content="此项选择为是，则不会展示左侧菜单以及顶部信息。"-->
<!--                    placement="top"-->
<!--                    effect="light"-->
<!--                  >-->
<!--                    <el-icon><QuestionFilled /></el-icon>-->
<!--                  </el-tooltip>-->
<!--                </div>-->
<!--              </template>-->

<!--              <el-select-->
<!--                v-model="form.meta.defaultMenu"-->
<!--                style="width: 100%"-->
<!--                placeholder="是否为基础页面"-->
<!--              >-->
<!--                <el-option-->
<!--                  :value="false"-->
<!--                  label="否"-->
<!--                />-->
<!--                <el-option-->
<!--                  :value="true"-->
<!--                  label="是"-->
<!--                />-->
<!--              </el-select>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--        </el-row>-->
      </el-form>
      <div>
<!--        <div class="flex items-center gap-2">-->
<!--          <el-button-->
<!--            type="primary"-->
<!--            -->
<!--            @click="addParameter(form)"-->
<!--          >-->
<!--            新增菜单参数-->
<!--          </el-button>-->
<!--        </div>-->
<!--        <el-table-->
<!--          :data="form.parameters"-->
<!--          style="width: 100%; margin-top: 12px"-->
<!--        >-->
<!--          <el-table-column-->
<!--            align="left"-->
<!--            prop="type"-->
<!--            label="参数类型"-->
<!--            width="180"-->
<!--          >-->
<!--            <template #default="scope">-->
<!--              <el-select-->
<!--                v-model="scope.row.type"-->
<!--                placeholder="请选择"-->
<!--              >-->
<!--                <el-option-->
<!--                  key="query"-->
<!--                  value="query"-->
<!--                  label="query"-->
<!--                />-->
<!--                <el-option-->
<!--                  key="params"-->
<!--                  value="params"-->
<!--                  label="params"-->
<!--                />-->
<!--              </el-select>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--          <el-table-column-->
<!--            align="left"-->
<!--            prop="key"-->
<!--            label="参数key"-->
<!--            width="180"-->
<!--          >-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-input v-model="scope.row.key" />-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--          <el-table-column-->
<!--            align="left"-->
<!--            prop="value"-->
<!--            label="参数值"-->
<!--          >-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-input v-model="scope.row.value" />-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--          <el-table-column align="left">-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-button-->
<!--                  type="primary"-->
<!--                  -->
<!--                  @click="deleteParameter(form.parameters, scope.$index)"-->
<!--                >-->
<!--                  删除-->
<!--                </el-button>-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--        </el-table>-->

<!--        <div class="flex items-center gap-2 mt-3">-->
<!--          <el-button-->
<!--            type="primary"-->
<!--            -->
<!--            @click="addBtn(form)"-->
<!--          >-->
<!--            新增可控按钮-->
<!--          </el-button>-->
<!--          <el-icon-->
<!--            class="cursor-pointer"-->
<!--            @click="-->
<!--              toDoc('https://www.gin-vue-admin.com/guide/web/button-auth.html')-->
<!--            "-->
<!--          >-->
<!--            <QuestionFilled />-->
<!--          </el-icon>-->
<!--        </div>-->

<!--        <el-table-->
<!--          :data="form.menuBtn"-->
<!--          style="width: 100%; margin-top: 12px"-->
<!--        >-->
<!--          <el-table-column-->
<!--            align="left"-->
<!--            prop="name"-->
<!--            label="按钮名称"-->
<!--            width="180"-->
<!--          >-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-input v-model="scope.row.name" />-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--          <el-table-column-->
<!--            align="left"-->
<!--            prop="name"-->
<!--            label="备注"-->
<!--            width="180"-->
<!--          >-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-input v-model="scope.row.desc" />-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--          <el-table-column align="left">-->
<!--            <template #default="scope">-->
<!--              <div>-->
<!--                <el-button-->
<!--                  type="primary"-->
<!--                  -->
<!--                  @click="deleteBtn(form.menuBtn, scope.$index)"-->
<!--                >-->
<!--                  删除-->
<!--                </el-button>-->
<!--              </div>-->
<!--            </template>-->
<!--          </el-table-column>-->
<!--        </el-table>-->
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
