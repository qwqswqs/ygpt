<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
            style=""
            label-position="left"
        >
          <el-form-item label="租户名">
            <el-input v-model="searchInfo.username" placeholder="请输入"/>
          </el-form-item>
          <el-form-item label="公司名">
            <el-input
                v-model="searchInfo.companyName"
                placeholder="请输入"
            />
          </el-form-item>
          <el-form-item label="租户来源">
            <el-select
                style="width: 240px;"
                v-model="searchInfo.source"
                placeholder="请选择"
            >
              <el-option label="本地注册" :value="1"/>
              <el-option label="算力注册" :value="2"/>
            </el-select>
          </el-form-item>
          <el-form-item style="float: right;">
            <el-button type="primary" @click="GetRenterList">查询</el-button>
            <el-button type="info" @click="reSetGet">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

    </div>
    <div class="gva-table-box">
      <el-button type="primary" @click="addUser" style="margin-bottom:15px">新增租户</el-button>
      <el-table
          :data="tableData"
          row-key="ID"
          @sort-change="handleSortChange"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
      >
        <el-table-column
            align="left"
            label="租户名"
            min-width="120"
            prop="username"
        />
        <el-table-column
            align="left"
            label="公司名"
            min-width="120"
            prop="companyName"
        >
          <template #default="scope">
            {{ scope.row.companyName!=''?scope.row.companyName:'-' }}
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="租户来源"
        >
          <template #default="scope">
            {{ Source[scope.row.source - 1] }}
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <!--            <el-tag :type="handleTagType(scope.row.status)">-->
            <!--              -->
            <!--            </el-tag>-->
            <div style="display: flex;align-items: center">
              <div :class="handleTagType(scope.row.status)"> </div>
              <div>{{ scope.row.status == 1 ? '正常' : '过期' }}</div>
            </div>
          </template>
        </el-table-column>
<!--        <el-table-column-->
<!--            align="left"-->
<!--            label="状态"-->
<!--        >-->
<!--          <template #default="scope">-->
<!--            <span style="color: green">{{ scope.row.status == 1 ? '正常' : '过期' }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column
            align="left"
            label="创建时间"
            sortable="custom"
        >
          <template #default="scope">
            {{ FormatDateTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="text" link @click="ResetPassword(scope.row)">重置密码</el-button>
            <el-button type="text" link @click="OpenResDrawer(scope.row)">资源列表</el-button>
          </template>
        </el-table-column>

      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="searchInfo.total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!--    资源列表-->
    <el-drawer  :show-close="false" v-model="resDialog" size="1200px">
      <div class="flex justify-between items-center" style="margin-bottom: 16px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">资源列表</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="resDialog=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-tabs v-model="activeName" @tab-change="handleTabsClick">
        <el-tab-pane :label="'裸金属'" name="baremetal">
          <div class="gva-table-box" style="margin-top: 16px">

            <!-- 表格展示 -->
            <el-table
                :data="renterResList"
                row-key="ID"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
            >
              <el-table-column label="名称"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.name }}</span>
                </template>
              </el-table-column>
              <el-table-column label="密码"align="left" min-width="120">
                <template #default="scope">
                  <el-tooltip>
                    <template #content>
                      <span>{{ '用户名:' + scope.row.username + '\n' + '密码:' + scope.row.password }}</span>
                    </template>
                    <span>查看密码</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column label="私有ip地址"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.privateIp }}</span>
                </template>
              </el-table-column>
              <el-table-column label="运行状态"align="left" min-width="120">
                <template #default="scope">
                  <i v-if="scope.row.status == 'running'" class="status-dot"/>
                  <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                  <span>{{ hostStatus[scope.row.status] }}</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="gva-pagination">
              <el-pagination
                  :current-page="tabSearchData.page"
                  :page-size="tabSearchData.pageSize"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="tabSearchData.total"
                  layout="total,prev, pager, next,  sizes, jumper"
                  @current-change="handleTabCurrentChange"
                  @size-change="handleTabSizeChange"
              />
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="'云主机'" name="virtual">
          <div class="gva-table-box" style="margin-top: 16px">

            <!-- 表格展示 -->
            <el-table
                :data="renterResList"
                row-key="ID"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
            >
              <el-table-column label="名称"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.name }}</span>
                </template>
              </el-table-column>
              <el-table-column label="密码"align="left" min-width="120">
                <template #default="scope">
                  <el-popover placement="top" :width="180" trigger="click">
                    <template #reference>
                      <el-button type="text" style="margin-right: 16px">查看密码
                      </el-button>
                    </template>
                    <span>用户名:{{scope.row.username }}</span><br/>
                    <span>密码:{{ scope.row.password }}</span>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="私有ip地址"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.privateIp }}</span>
                </template>
              </el-table-column>
              <el-table-column label="运行状态"align="left" min-width="120">
                <template #default="scope">
                  <i v-if="scope.row.status == 'running'" class="status-dot"/>
                  <i v-else class="status-dot" style="background-color: rgb(174,185,192)"/>
                  <span>{{ hostStatus[scope.row.status] }}</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="gva-pagination">
              <el-pagination
                  :current-page="tabSearchData.page"
                  :page-size="tabSearchData.pageSize"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="tabSearchData.total"
                  layout="total,prev, pager, next,  sizes, jumper"
                  @current-change="handleTabCurrentChange"
                  @size-change="handleTabSizeChange"
              />
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="'容器'" name="container">
          <div class="gva-table-box" style="margin-top: 16px">

            <!-- 表格展示 -->
            <el-table
                :data="renterResList"
                row-key="ID"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"

            >
              <el-table-column label="名称"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.name }}</span>
                </template>
              </el-table-column>
              <el-table-column label="私有ip地址"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.privateIp }}</span>
                </template>
              </el-table-column>
              <el-table-column label="镜像"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.image }}</span>
                </template>
              </el-table-column>
              <el-table-column label="运行状态"align="left" min-width="120">
                <template #default="scope">
                  <span>{{ scope.row.status }}</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="gva-pagination">
              <el-pagination
                  :current-page="tabSearchData.page"
                  :page-size="tabSearchData.pageSize"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="tabSearchData.total"
                  layout="total,prev, pager, next,  sizes, jumper"
                  @current-change="handleTabCurrentChange"
                  @size-change="handleTabSizeChange"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
    <!--    新增租户-->
    <el-drawer
        v-model="addUserDialog"
        size="520px"
        :show-close="false"
        :close-on-press-escape="false"
        :close-on-click-modal="false"
    >
      <div class="flex justify-between items-center" style="margin-bottom: 16px">
        <span class="text-lg" style="font-size: 24px;font-weight: 500;color: rgba(29, 33, 41, 1);">新增租户</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="closeAddUserDialog"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form
          ref="userForm"
          :rules="rules"
          :model="userInfo"
          label-width="80px"
      >
        <el-form-item
            label="手机号"
            prop="phone"
        >
          <el-input placeholder="请输入" style="width: 320px" v-model="userInfo.phone"/>
        </el-form-item>
        <el-form-item
            label="密码"
            prop="password"
        >
          <el-input  placeholder="请输入" style="width: 320px" v-model="userInfo.password"/>
        </el-form-item>
        <el-form-item
            label="邮箱"
            prop="email"
            label-width="70px" style="padding-left: 10px"
        >
          <el-input  placeholder="请输入" style="width: 320px" v-model="userInfo.email"/>
        </el-form-item>
        <el-form-item
            label="公司名"
            prop="companyName"
            label-width="70px" style="padding-left: 10px"
        >
          <el-input  placeholder="请输入" style="width: 320px" v-model="userInfo.companyName"/>
        </el-form-item>
      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="closeAddUserDialog">取 消</el-button>
        <el-tooltip :disabled="licenseConfig.renterNum>searchInfo.total" content="当前系统租户数量已达上限"
                    placement="top">
          <el-button
              :disabled="licenseConfig.renterNum<=searchInfo.total"
              type="primary"
              @click="enterAddUserDialog"
          >确 定
          </el-button>
        </el-tooltip>
      </div>


    </el-drawer>
  </div>
</template>
<script>
import {addRenter, getRenterList} from "@/api/yunguan/run/renter";
import {resetPassword} from "@/api/system/user";
import {getRenterResList} from "@/api/yunguan/resource/resInfo";
import SelectImage from "@/components/selectImage/selectImage.vue";
import statusType from "@/locales/zh-CN.json";
import {queryLicense} from "@/api/yunguan/config/licenseConfig";

export default {
  components: {SelectImage},
  data() {
    return {
      licenseConfig:{},
      addUserDialog: false,
      activeName: 'baremetal',
      renterID: '',
      searchInfo: {
        page: 1,
        pageSize: 10,
        timeDesc:"",
        total:0,
      },
      userInfo: {
        username: '',
        phone: '',
        password: '',
        companyName: '',
        email: '',
        source: 1,
        enable: 1,
        authorityId: 777,
        authorityIds: [777],
      },
      rules: {
        phone: [
          {required: true, message: '请输入手机号码', trigger: 'blur'},
          {
            pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
            message: '请输入合法手机号',
            trigger: 'blur'
          },
        ],
        password: [
          {required: true, message: '请输入用户密码', trigger: 'blur'},
          {min: 6, message: '最低6位字符', trigger: 'blur'}
        ],
        email: [
          {
            pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
            message: '请输入正确的邮箱',
            trigger: 'blur'
          },
        ],
      },
      tabSearchData: {
        page: 1,
        pageSize: 10,
        total: 0,
        renterID: '',
        type: 1,
      },
      tableData: [],
      Source: ["本地注册", "算力注册"],
      renterResList: [],
      resDialog: false,
      hostStatus: statusType.status.server,
    }
  },
  created() {
    this.GetRenterList();
  },
  methods: {
    handleSortChange({ column, prop, order }) {
      switch (order) {
        case "ascending":
          this.searchInfo.timeDesc="asc";
          break;
        case "descending":
          this.searchInfo.timeDesc="desc";
          break;
        default:
          this.searchInfo.timeDesc="";
          break;
      }
      this.GetRenterList();
    },
    handleTagType(status) {
      switch (status) {
        case 1:
          return 'success'
        case 2:
          return 'danger'
        default:
          return 'info'
      }
    },
    //获取License
    GetLicenseConfig() {
      queryLicense().then(res => {
        if (res.code == 0) {
          this.licenseConfig = res.data;
        }
      })
    },
    //时间格式化
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        return '';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    //打开资源列表
    OpenResDrawer(val) {
      this.tabSearchData.renterID = val.renterID;
      this.tabSearchData.type = 1;
      this.renterResList = [];
      this.resDialog = true;
      this.activeName = 'baremetal'
      this.GetRenterRes();
    },
    //资源列表tab切换
    handleTabsClick() {
      switch (this.activeName) {
        case "baremetal":
          this.tabSearchData.type = 1
          break;
        case "virtual":
          this.tabSearchData.type = 2
          break;
        case "container":
          this.tabSearchData.type = 3
          break;
      }
      this.GetRenterRes();
    },
    //资源列表
    GetRenterRes() {
      this.renterResList = []
      getRenterResList(this.tabSearchData).then(res => {
        if (res.code == 0) {
          this.renterResList = res.data.list
          this.tabSearchData.total = res.data.total
        }
      })
    },
    //重置密码
    ResetPassword(val) {
      this.$confirm('是否将此用户密码重置为123456?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let obj = {}
        obj.id = val.renterID
        resetPassword(obj).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '重置成功',
              type: 'success'
            })
            this.GetRenterList(this.searchInfo)
          } else {
            this.$message({
              message: '重置失败',
              type: 'error'
            });
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消'
        });
      });
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetRenterList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetRenterList()
    },
    handleTabCurrentChange(val) {
      this.tabSearchData.page = val
      this.GetRenterRes()
    },
    handleTabSizeChange(val) {
      this.tabSearchData.pageSize = val
      this.GetRenterRes()
    },
    //获取租户列表
    GetRenterList() {
      getRenterList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.tableData = response.data.list
          this.searchInfo.total = response.data.total
        }
      })
    },
    reSetGet() {
      this.searchInfo= {
        page: 1,
        pageSize: 10,
        timeDesc:"",
      }
      this.getRenterList();
    },
    //打开新增对话框
    addUser() {
      this.GetLicenseConfig();
      this.addUserDialog = true;
      this.userInfo.phone = '';
      this.userInfo.companyName = '';
      this.userInfo.email = '';
      this.userInfo.password = '';
    },
    //取消
    closeAddUserDialog() {
      this.addUserDialog = false;
    },
    //确定
    enterAddUserDialog() {
      this.userInfo.username = this.userInfo.phone;
      this.$refs.userForm.validate(valid => {
        if (valid) {
          addRenter(this.userInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '创建成功',
              })
            }
            this.GetRenterList();
            this.addUserDialog = false;
          })
        }
      })
    },
  }
}
</script>

<style scoped lang="scss">
//文本溢出截断
.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.danger {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(245, 63, 63, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}

.success {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(0, 180, 42, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}

</style>
