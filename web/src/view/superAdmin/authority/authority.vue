<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-table-box">
      <div style="margin-bottom: 16px;">
        <el-button
          type="primary"
          @click="addAuthority(0)"
        >新增角色</el-button>
      </div>
      <el-table
        :data="tableData"
        :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
        :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
        row-key="authorityId"
        style="width: 100%"
      >
<!--        <el-table-column-->
<!--          label="角色ID"-->
<!--          min-width="180"-->
<!--          prop="authorityId"-->
<!--        />-->
        <el-table-column
          align="left"
          label="角色名称"
          min-width="180"
          prop="authorityName"
        />
        <el-table-column
          align="left"
          label="操作"
          width="460"
        >
          <template #default="scope">
            <el-button
              :disabled="scope.row.authorityId==1&&(userStore.userInfo.authority.authorityId!=1)"
              type="text" class="text-blue-500"
              link
              @click="openDrawer(scope.row)"
            >设置权限</el-button>
<!--            <el-button-->
<!--              -->

<!--              type="primary"-->
<!--              link-->
<!--              @click="addAuthority(scope.row.authorityId)"-->
<!--            >新增子角色</el-button>-->
<!--            <el-button-->
<!--              icon="copy-document"-->

<!--              type="primary"-->
<!--              link-->
<!--              @click="copyAuthorityFunc(scope.row)"-->
<!--            >拷贝</el-button>-->
            <el-button
                v-if="scope.row.authorityId != 777"
                :class="{'custom-disabled': scope.row.authorityId === 777}"
                type="text" class="text-blue-500"
                link
                @click="editAuthority(scope.row)"
            >编辑</el-button>
            <el-button
                v-if="scope.row.authorityId != 777"
                :class="{'custom-disabled': scope.row.authorityId === 777}"
                type="text" class="text-blue-500"
                link
                @click="openDeleteDialog(scope.row)"
            >删除</el-button>

          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[2,10, 30, 50, 100]"
            :total="searchInfo.total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!-- 新增角色弹窗 -->
    <el-drawer
      v-model="authorityFormVisible"
      :show-close="false"
      size="520px"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ authorityTitleForm }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="authorityFormVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
        <Close />
        </el-icon>
      </div>
      <el-form
        ref="authorityForm"
        :model="form"
        :rules="rules"
        label-width="80px"
        :style="{ marginBottom: '24px' }"
      >
<!--        <el-form-item-->
<!--          label="父级角色"-->
<!--          prop="parentId"-->
<!--        >-->
<!--          <el-cascader-->
<!--            v-model="form.parentId"-->
<!--            style="width:100%"-->
<!--            :disabled="dialogType==='add'"-->
<!--            :options="AuthorityOption"-->
<!--            :props="{ checkStrictly: true,label:'authorityName',value:'authorityId',disabled:'disabled',emitPath:false}"-->
<!--            :show-all-levels="false"-->
<!--            filterable-->
<!--          />-->
<!--        </el-form-item>-->
        <el-form-item
          label="角色姓名"
          prop="authorityName"
        >
          <el-input
            v-model="form.authorityName"
            autocomplete="off"
            style="width: 320px"
          />
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right">
        <el-button type="info" @click="closeAuthorityForm">取 消</el-button>
        <el-button
            type="primary"
            @click="submitAuthorityForm"
        >确 定</el-button>
      </div>
    </el-drawer>

    <el-drawer
      v-if="drawer"
      v-model="drawer"
      :with-header="false"
      size="520px"
      title="角色配置"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">角色菜单</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="drawer = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>

      <Menus
          ref="menus"
          :row="activeRow"
          @changeRow="changeRow"
          v-model="drawer"
      />
<!--      <el-tabs-->
<!--        :before-leave="autoEnter"-->
<!--        type="border-card"-->
<!--      >-->
<!--        <el-tab-pane label="角色菜单">-->

<!--        </el-tab-pane>-->
<!--&lt;!&ndash;        <el-tab-pane label="角色api">&ndash;&gt;-->
<!--&lt;!&ndash;          <Apis&ndash;&gt;-->
<!--&lt;!&ndash;            ref="apis"&ndash;&gt;-->
<!--&lt;!&ndash;            :row="activeRow"&ndash;&gt;-->
<!--&lt;!&ndash;            @changeRow="changeRow"&ndash;&gt;-->
<!--&lt;!&ndash;          />&ndash;&gt;-->
<!--&lt;!&ndash;        </el-tab-pane>&ndash;&gt;-->
<!--&lt;!&ndash;        <el-tab-pane label="资源权限">&ndash;&gt;-->
<!--&lt;!&ndash;          <Datas&ndash;&gt;-->
<!--&lt;!&ndash;            ref="datas"&ndash;&gt;-->
<!--&lt;!&ndash;            :authority="tableData"&ndash;&gt;-->
<!--&lt;!&ndash;            :row="activeRow"&ndash;&gt;-->
<!--&lt;!&ndash;            @changeRow="changeRow"&ndash;&gt;-->
<!--&lt;!&ndash;          />&ndash;&gt;-->
<!--&lt;!&ndash;        </el-tab-pane>&ndash;&gt;-->
<!--      </el-tabs>-->
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
  getAuthorityList,
  deleteAuthority,
  createAuthority,
  updateAuthority,
  copyAuthority, getAllDataAuthority
} from '@/api/system/authority'

import Menus from '@/view/superAdmin/authority/components/menus.vue'
import Apis from '@/view/superAdmin/authority/components/apis.vue'
import Datas from '@/view/superAdmin/authority/components/datas.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'

import { ref,  provide} from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {useUserStore} from "@/pinia";
import {deleteCategory, deleteCategoryByIds} from "@/api/yunguan/system/type";

const deleteVisible =ref(false)
const currentProductId = ref(null);
const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};
// 打开删除弹窗（统一入口）
const openDeleteDialog = ( row = null) => {
  currentProductId.value = row.authorityId;
  deleteVisible.value = true;
};
const confirmDelete = async () => {
  try {
    let resp;
    resp = await deleteAuthority({ authorityId: currentProductId.value });

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
defineOptions({
  name: 'Authority'
})
const userStore = useUserStore();
const mustUint = (rule, value, callback) => {
  if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
    return callback(new Error('请输入正整数'))
  }
  return callback()
}

const AuthorityOption = ref([
  {
    authorityId: 0,
    authorityName: '根角色/严格模式下为当前角色'
  }
])
const drawer = ref(false)
const dialogType = ref('add')
const activeRow = ref({})

const authorityTitleForm = ref('新增角色')
const authorityFormVisible = ref(false)
const apiDialogFlag = ref(false)
const copyForm = ref({})

const form = ref({
  authorityId: 22,
  authorityName: '',
  parentId: 1,
})
const rules = ref({
  authorityId: [
    { required: true, message: '请输入角色ID', trigger: 'blur' },
    { validator: mustUint, trigger: 'blur', message: '必须为正整数' }
  ],
  authorityName: [
    { required: true, message: '请输入角色名', trigger: 'blur' }
  ],
  parentId: [
    { required: true, message: '请选择父角色', trigger: 'blur' },
  ]
})
// 分页
const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val;
  getTableData();
};

const handleCurrentChange = (val) => {
  searchInfo.value.page = val;
  getTableData();
};
const tableData = ref([])
const allDataList = ref([])
const searchInfo=ref({
  page:1,
  pageSize:10,
  total:0,
})
// 查询
const getTableData = async() => {
  const table = await getAuthorityList(searchInfo.value)
  if (table.code === 0) {
    tableData.value = table.data.list;
    searchInfo.value.total=table.data.total;
  }
}
// 查询
const getAllDataList = async() => {
  const table = await getAllDataAuthority()
  if (table.code === 0) {
    allDataList.value = table.data
  }
}

getAllDataList()
getTableData()

const changeRow = (key, value) => {
  activeRow.value[key] = value
}
const menus = ref(null)
const apis = ref(null)
const datas = ref(null)
const autoEnter = (activeName, oldActiveName) => {
  const paneArr = [menus, apis, datas]
  if (oldActiveName) {
    if (paneArr[oldActiveName].value.needConfirm) {
      paneArr[oldActiveName].value.enterAndNext()
      paneArr[oldActiveName].value.needConfirm = false
    }
  }
}
// 拷贝角色
const copyAuthorityFunc = (row) => {
  setOptions()
  authorityTitleForm.value = '拷贝角色'
  dialogType.value = 'copy'
  for (const k in form.value) {
    form.value[k] = row[k]
  }
  copyForm.value = row
  authorityFormVisible.value = true
}
const openDrawer = (row) => {
  drawer.value = true
  activeRow.value = row
}
// 删除角色
const deleteAuth = (row) => {
  ElMessageBox.confirm('此操作将永久删除该角色, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await deleteAuthority({ authorityId: row.authorityId })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })

        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除'
      })
    })
}
// 初始化表单
const authorityForm = ref(null)
const initForm = () => {
  if (authorityForm.value) {
    authorityForm.value.resetFields()
  }
  form.value = {
    authorityId: 22,
    authorityName: '',
    parentId: 1
  }
}
// 关闭窗口
const closeAuthorityForm = () => {
  initForm()
  authorityFormVisible.value = false
  apiDialogFlag.value = false
}

provide('closeDrawer', () => {
  drawer.value = false;
});
// 确定弹窗

const submitAuthorityForm = () => {
  authorityForm.value.validate(async valid => {
    if (valid) {
      form.value.authorityId = Number(form.value.authorityId)
      switch (dialogType.value) {
        case 'add':
          {
            form.value.authorityId=22
            form.value.parentId=0
            const res = await createAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '添加成功!'
              })
              getTableData()
              closeAuthorityForm()
            }
          }
          break
        case 'edit':
          {
            const res = await updateAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '添加成功!'
              })
              getTableData()
              closeAuthorityForm()
            }
          }
          break
        case 'copy': {
          const data = {
            authority: {
              authorityId: 0,
              authorityName: '',
              datauthorityId: [],
              parentId: 0
            },
            oldAuthorityId: 0
          }
          data.authority.authorityId = form.value.authorityId
          data.authority.authorityName = form.value.authorityName
          data.authority.parentId = form.value.parentId
          data.authority.dataAuthorityId = copyForm.value.dataAuthorityId
          data.oldAuthorityId = copyForm.value.authorityId
          const res = await copyAuthority(data)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '复制成功！'
            })
            getTableData()
          }
        }
      }

      initForm()
      authorityFormVisible.value = false
    }
  })
}
const setOptions = () => {
  AuthorityOption.value = [
    {
      authorityId: 22,
      authorityName: '根角色(严格模式下为当前用户角色)'
    }
  ]
  setAuthorityOptions(tableData.value, AuthorityOption.value, false)
}
const setAuthorityOptions = (AuthorityData, optionsData, disabled) => {
  AuthorityData &&
        AuthorityData.forEach(item => {
          if (item.children && item.children.length) {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId,
              children: []
            }
            setAuthorityOptions(
              item.children,
              option.children,
              disabled || item.authorityId === form.value.authorityId
            )
            optionsData.push(option)
          } else {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId
            }
            optionsData.push(option)
          }
        })
}
// 增加角色
const addAuthority = (parentId) => {
  initForm()
  authorityTitleForm.value = '新增角色'
  dialogType.value = 'add'
  form.value.parentId = 1
  setOptions()
  authorityFormVisible.value = true
}
// 编辑角色
const editAuthority = (row) => {
  setOptions()
  authorityTitleForm.value = '编辑角色'
  dialogType.value = 'edit'
  for (const key in form.value) {
    form.value[key] = row[key]
  }
  setOptions()
  authorityForm.value && authorityForm.value.clearValidate()
  authorityFormVisible.value = true
}

</script>

<style lang="scss">
.authority {
  .el-input-number {
    margin-left: 15px;
    span {
      display: none;
    }
  }
}
.tree-content{
  margin-top: 10px;
  height: calc(100vh - 158px);
  overflow: auto;
}

.custom-disabled {
  color: rgba(201, 205, 212, 1) !important;
  cursor: not-allowed;
}
.custom-disabled:hover {
  color: rgba(201, 205, 212, 1) !important;
}

</style>
