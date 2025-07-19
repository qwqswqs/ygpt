<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <warning-bar title="注：右上角头像下拉可切换角色"/>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
            style="display: flex;justify-content: space-between"
        >
          <div>
            <el-form-item label="用户名">
              <el-input
                  v-model="searchInfo.username"
                  placeholder="请输入"
              />
            </el-form-item>
            <el-form-item label="昵称">
              <el-input
                  v-model="searchInfo.nickname"
                  placeholder="请输入"
              />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input
                  v-model="searchInfo.phone"
                  placeholder="请输入"
              />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input
                  v-model="searchInfo.email"
                  placeholder="请输入"
              />
            </el-form-item>
          </div>
          <el-form-item style="margin-right: 0;">
            <el-button
                type="primary"
                @click="onSubmit"
            >
              查询
            </el-button>
            <el-button
                class="button-gap" type="info"
                @click="onReset"
            >
              重置
            </el-button>
            
          </el-form-item>
        </el-form>
      </div>

    </div>
    <div class="gva-table-box">
      <el-button  type="primary" @click="addUser" style="margin-bottom:15px">新增用户</el-button>
      <el-table
          :data="tableData"
          row-key="ID"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
        <el-table-column
            align="left"
            label="头像"
            min-width="75"

        >
          <template #default="scope">
            <CustomPic
                :pic-src="scope.row.headerImg"
            />
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="用户名"
            min-width="150"
            prop="userName"
        />
        <el-table-column
            align="left"
            label="昵称"
            min-width="150"
            prop="nickName"
        />
        <el-table-column
            align="left"
            label="手机号"
            min-width="130"
            prop="phone"
        />
        <el-table-column
            align="left"
            label="邮箱"
            min-width="180"
            prop="email"
        />
        <el-table-column
            align="left"
            label="用户角色"
            min-width="300"
        >
          <template #default="scope">
            <span v-for="item in scope.row.authorities">
              <el-tag style="margin-right: 5px">{{ item.authorityName }}</el-tag>
            </span>
            <!--            <el-cascader-->
            <!--              v-model="scope.row.authorityIds"-->
            <!--              :options="authOptions"-->
            <!--              :show-all-levels="false"-->
            <!--              collapse-tags-->
            <!--              :props="{ multiple:true,checkStrictly: true,label:'authorityName',value:'authorityId',disabled:'disabled',emitPath:false}"-->
            <!--              :clearable="false"-->
            <!--              @visible-change="(flag)=>{changeAuthority(scope.row,flag,0)}"-->
            <!--              @remove-tag="(removeAuth)=>{changeAuthority(scope.row,false,removeAuth)}"-->
            <!--            />-->
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="启用"
            min-width="80"
        >
          <template #default="scope">
            <el-switch
                v-model="scope.row.enable"
                inline-prompt
                :active-value="1"
                :inactive-value="2"
                @change="()=>{switchEnable(scope.row)}"
            />
          </template>
        </el-table-column>

        <el-table-column
            label="操作"
            min-width="250"
            fixed="right"
        >
          <template #default="scope">
            <el-button
                type="text"
                class="text-blue-500"
                link
                @click="openDeleteDialog(scope.row)"
            >删除
            </el-button>
            <el-button
                type="text"
                class="text-blue-500"
                link
                @click="openEdit(scope.row)"
            >编辑
            </el-button>
            <el-button
                type="text"
                class="text-blue-500"
                link
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
        size="520px"
        :show-close="false"
        :close-on-press-escape="true"
        :close-on-click-modal="true"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 24px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增用户</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="closeAddUserDialog"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form
          ref="userForm"
          :rules="rules"
          :model="userInfo"
          label-width="80px"
      >
        <el-form-item
            label="用户名"
            prop="userName"

        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.userName"/>
        </el-form-item>
        <el-form-item
            label="密码"
            prop="password"
        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.password"/>
        </el-form-item>
        <el-form-item
            label="昵称"
            prop="nickName"
        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.nickName"/>
        </el-form-item>
        <el-form-item
            label="手机号"
            prop="phone"
            label-width="70px" style="padding-left: 10px"
        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.phone"/>
        </el-form-item>
        <el-form-item
            label="邮箱"
            prop="email"
            label-width="70px" style="padding-left: 10px"
        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.email"/>
        </el-form-item>
        <el-form-item
            label="用户角色"
            prop="authorityId"
        >
          <el-cascader
              v-model="userInfo.authorityIds"
              style="width:320px"
              :options="authOptions"
              :show-all-levels="false"
              :props="{ multiple:true,checkStrictly: true,label:'authorityName',value:'authorityId',disabled:'disabled',emitPath:false}"
              :clearable="false"
          />
        </el-form-item>
        <el-form-item
            label="启用"
            prop="disabled"
            label-width="70px" style="padding-left: 10px"
        >
          <el-switch
              v-model="userInfo.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
          />
        </el-form-item>
        <el-form-item
            label="头像"
            label-width="70px" style="padding-left: 10px"
        >
          <SelectImage
              style="width: 100px;height: 100px"
              v-model="userInfo.headerImg"
          />
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="closeAddUserDialog">取 消</el-button>

        <el-tooltip v-if="dialogFlag === 'add'" :disabled="licenseConfig.adminNum>total"
                    content="当前系统管理员数量已达上限"
                    placement="top">
          <el-button
              :disabled="licenseConfig.adminNum<=total"
              type="primary"
              @click="enterAddUserDialog"
          >确 定
          </el-button>
        </el-tooltip>
        <el-button
            v-if="dialogFlag === 'edit'"
            type="primary"
            @click="enterAddUserDialog"
        >确 定
        </el-button>
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

import {deleteUser, getUserList, register, setUserAuthorities} from '@/api/system/user'

import {deleteAuthority, getAuthorityList} from '@/api/system/authority'
import CustomPic from '@/components/customPic/index.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import {resetPassword, setUserInfo} from '@/api/system/user.js'

import {nextTick, ref, watch} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'
import {queryLicense} from "@/api/yunguan/config/licenseConfig";

defineOptions({
  name: 'User',
})

const deleteVisible =ref(false)
const currentProductId = ref(null);
const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};
// 打开删除弹窗（统一入口）
const openDeleteDialog = ( row = null) => {
  currentProductId.value = row.ID;
  deleteVisible.value = true;
};
const confirmDelete = async () => {
  try {
    let resp;
    resp = await deleteUser({ id: currentProductId.value });

    // 严格判断接口成功状态
    if (resp.code === 0) {
      // 显示对应成功消息
      ElMessage({
        type: 'success',
        message:'删除成功'
      });
      deleteVisible.value = false;

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
const licenseConfig = ref({})
const GetLicenseConfig = () => {
  queryLicense().then(res => {
    if (res.code == 0) {
      licenseConfig.value = res.data;
    }
  })
};

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

const onReset = () => {
  searchInfo.value = {
    username: '',
    nickname: '',
    phone: '',
    email: ''
  }
  getTableData()
}
// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
  AuthorityData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName,
        children: []
      }
      setAuthorityOptions(item.children, option.children)
      optionsData.push(option)
    } else {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName
      }
      optionsData.push(option)
    }
  })
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
  const table = await getUserList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
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
  const res = await getAuthorityList({page: 1, pageSize: 1000})
  console.log(res.data)
  setOptions(res.data)
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

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  const tempData = authData.list.filter(item => item.authorityId != 777)
  setAuthorityOptions(tempData, authOptions.value)
}

const deleteUserFunc = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteUser({id: row.ID})
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
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
  enable: 1,
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
  nickName: [
    {required: true, message: '请输入用户昵称', trigger: 'blur'}
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
  userInfo.value.authorityId = userInfo.value.authorityIds[0]
  userForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...userInfo.value
      }
      if (dialogFlag.value === 'add') {
        const res = await register(req)
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
  GetLicenseConfig();
}

const tempAuth = {}
const changeAuthority = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds]
    }
    return
  }
  await nextTick()
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds
  })
  if (res.code === 0) {
    ElMessage({type: 'success', message: '角色设置成功'})
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds]
    }
  }
}

const openEdit = (row) => {
  dialogFlag.value = 'edit'
  userInfo.value = JSON.parse(JSON.stringify(row))
  addUserDialog.value = true
}

const switchEnable = async (row) => {
  userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...userInfo.value
  }
  const res = await setUserInfo(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `${req.enable === 2 ? '禁用' : '启用'}成功`})
    await getTableData()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
  }
}

</script>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}


</style>
