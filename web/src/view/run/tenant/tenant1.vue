<template>
  <div>
    <div class="gva-search-box">
      <el-form
          ref="searchForm"
          :inline="true"
          :model="searchInfo"
          label-position="left"
      >
        <el-form-item label="用户名">
          <el-input
              v-model="searchInfo.username"
              placeholder="用户名"
          />
        </el-form-item>
        <el-form-item>
          <el-button

              type="primary"
              @click="onSubmit"
          >
            查询
          </el-button>
          <!-- <el-button

            type="primary"
            @click="addUser"
        >新增租户
        </el-button> -->
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
          :data="tableData"
          row-key="ID"
      >
        <el-table-column
          type="selection"
          width="55">
        </el-table-column>
        <el-table-column
            align="left"
            type="index"
            label="序号"
            min-width="60"
        />
        <el-table-column
            align="left"
            label="租户编号"
            min-width="150"
            prop="tenantId"
        />
        <el-table-column
            align="left"
            label="租户名"
            min-width="150"
        >
          <template #default="scope">
            {{ scope.row.username || '-' }}
          </template>
        </el-table-column>>
        <el-table-column
            align="left"
            label="资源类别"
            min-width="180"
            prop="resourceType"
        />
        <el-table-column
            align="left"
            label="关联订单"
            min-width="180"
            prop="orderId"
        />
        <el-table-column
            align="left"
            label="状态"
            min-width="150">
          <template #default="scope">
            <span>{{scope.row.status == 1 ?'正常':'禁用'}}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="生成日期"
            min-width="180"
            prop="created_at">
        </el-table-column>
        <!-- <el-table-column
            align="left"
            label="启用"
            min-width="150"
        >
          <template #default="scope">
            <span>{{scope.row.status == 1 ?'正常':'禁用'}}</span>
          </template>
        </el-table-column> -->
        <!-- <el-table-column
            align="left"
            label="租户来源"
            min-width="180"
            prop="source"
        >
          <template #default="scope">
            <span>{{ (scope.row.source)==1?'本地注册': '算力注册' }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="上次登录时间"
            min-width="180"
            prop="loginTime"
        >
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.loginTime) }}</span>
          </template>
        </el-table-column> -->
        <el-table-column
            fixed="right"
            label="操作"
            min-width="250"
        >
          <template #default="scope">
            <el-button

                link
                type="primary"
                @click="deleteRenterFunc(scope.row)"
            >删除
            </el-button>
            <el-button
                icon="magic-stick"
                link
                type="primary"
                @click="resetPasswordFunc(scope.row)"
            >重置密码
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
        v-model="addUserDialog"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="false"
        size="60%"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button
                type="primary"
                @click="enterAddUserDialog"
            >确 定
            </el-button>
          </div>
        </div>
      </template>

      <el-form
          ref="userForm"
          :model="userInfo"
          :rules="rules"
          label-width="80px"
          label-position="left"
      >
        <el-form-item
            label="用户名"
            prop="userName"
        >
          <el-input v-model="userInfo.userName"/>
        </el-form-item>
        <el-form-item
            label="密码"
            prop="password"
        >
          <el-input v-model="userInfo.password"/>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>

import {
  deleteRenter, getRenterList, addRenter
} from '@/api/yunguan/run/renter'

import {setUserInfo, resetPassword} from '@/api/system/user.js'

import {nextTick, ref, watch} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'

defineOptions({
  name: 'User',
})

const searchInfo = ref({
  username: '',
  nickname: '',
  phone: '',
  email: ''
})

const onSubmit = () => {
  page.value = 1
  getTableData()
}
const FormatDateTime = (dateString) => {
  const date = new Date(dateString);
  if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
    // 如果日期是不合理的，返回空字符串
    return '';
  }
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

const onReset = () => {
  searchInfo.value = {
    username: '',
    nickname: '',
    phone: '',
    email: ''
  }
  getTableData()
}

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
const getTableData = async () => {
  const table = await getRenterList({page: page.value, pageSize: pageSize.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

watch(() => tableData.value, () => {
  setAuthorityIds()
})

const initPage = async () => {
  getTableData()
  // const res = await getAuthorityList()
  // setOptions(res.data)
}

initPage()

const resetPasswordFunc = (row) => {
  ElMessageBox.confirm(
      '是否将此用户密码重置为123456?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(async () => {
    const res = await resetPassword({
      ID: row.ID,
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}
const setAuthorityIds = () => {
  tableData.value && tableData.value.forEach((user) => {
    user.authorityIds = user.authorities && user.authorities.map(i => {
      return i.authorityId
    })
  })
}

const deleteRenterFunc = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteRenter({id: row.ID})
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

// 弹窗相关
const userInfo = ref({
  userName: '',
  password: '',
  authorityId: 777,
  authorityIds: [777],
  enable: 1,
  source:1,
})

const rules = ref({
  userName: [
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 5, message: '最低5位字符', trigger: 'blur'}
  ],
  password: [
    {required: true, message: '请输入用户密码', trigger: 'blur'},
    {min: 6, message: '最低6位字符', trigger: 'blur'}
  ],
  phone: [
    {pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/, message: '请输入合法手机号', trigger: 'blur'},
  ],
  email: [
    {
      pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
      message: '请输入正确的邮箱',
      trigger: 'blur'
    },
  ],
  authorityId: [
    {required: true, message: '请选择用户角色', trigger: 'blur'}
  ]
})
const userForm = ref(null)
const enterAddUserDialog = async () => {
  userForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...userInfo.value
      }
      if (dialogFlag.value === 'add') {
        const res = await addRenter(req)
        if (res.code === 0) {
          ElMessage({type: 'success', message: '创建成功'})
          await getTableData()
          closeAddUserDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await setUserInfo(req)
        if (res.code === 0) {
          ElMessage({type: 'success', message: '编辑成功'})
          await getTableData()
          closeAddUserDialog()
        }
      }
    }
  })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  userForm.value.resetFields()
  userInfo.value.headerImg = ''
  userInfo.value.authorityIds = []
  addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
  dialogFlag.value = 'add'
  addUserDialog.value = true
}

</script>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}
</style>
