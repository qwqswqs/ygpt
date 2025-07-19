<template>
  <div>
    <!-- 搜索区域 -->
    <div class="gva-search-box">
      <el-form label-position="left" ref="searchForm" :inline="true" :model="searchInfo" style="justify-content: space-between;width: 100%;display: flex">
        <el-form-item label="订单有效月份" prop="begin">
          <!--
          <el-date-picker v-model="exportTableReq.begin" type="datetime" placeholder="选择开始时间"
              :disabled-date="(date) => date > new Date()"
          />
          <el-date-picker v-model="exportTableReq.end" type="datetime" placeholder="选择结束时间"
              :disabled-date=" (date) => date < exportTableReq.begin || date > new Date()"
          />
           -->
          <el-date-picker
              v-model="value1"
              type="monthrange"
              start-placeholder="开始月份"
              end-placeholder="结束月份"
              :disabled-date="(date) => date > new Date()"
          />
        </el-form-item>
        <el-form-item label="">
          <div class="metal-types">
            <span>裸金属 ￥0.00</span>|
            <span>云主机 ￥0.00</span>|
            <span>容器 ￥0.00</span>|
            <span>存储 ￥0.00</span>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form ref="searchForm" :inline="true" :model="searchInfo">
          <!--
          <el-form-item label="">
            <el-select v-model="searchInfo.productType" placeholder="账号类型">
              <el-option label="企业用户" value="1" />
              <el-option label="个人用户" value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="">
            <el-select v-model="searchInfo.productType" placeholder="用户账号" clearable>
              <el-option label="AI" value="1" />
              <el-option label="实例" value="2" />
            </el-select>
          </el-form-item>
          -->
                    <el-form-item label="账户名称">
                      <el-select v-model="searchInfo.buyUserId" placeholder="账户名称" filterable clearable>
                        <el-option v-for="item in userList" :key="item.id" :label="item.userName" :value="item.id" />
                      </el-select>
                    </el-form-item>

          <el-form-item label="产品类型">
            <el-select v-model="searchInfo.type" placeholder="产品类型" clearable>
              <el-option label="裸金属" value="1" />
              <el-option label="云主机" value="2" />
              <el-option label="容器" value="3" />
              <el-option label="存储" value="4" />
            </el-select>
          </el-form-item>
        </el-form>

        <el-form-item>
          <el-button type="primary"  @click="onSubmit">
            查询
          </el-button>
          <el-button type="primary" @click="debounceExportTableHandle">
            导出账单
          </el-button>
        </el-form-item>
      </div>
      <el-table :data="tableData" style="width: 100%" @selection-change="handleSelectionChange" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
        <!--        <el-table-columnalign="left" type="selection" width="55" />-->
        <el-table-column align="left" prop="order.type" label="产品类型">
          <template #default="scope">
            <span v-if="scope.row.order.type === 1">裸金属</span>
            <span v-else-if="scope.row.order.type === 2">云主机</span>
            <span v-else-if="scope.row.order.type === 3">容器</span>
            <span v-else-if="scope.row.order.type === 4">存储</span>
            <span v-else>其他</span>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="order.createTime" label="账期">
          <template #default="scope">
            {{ new Date(scope.row.order.createTime).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}
          </template>
        </el-table-column>
        <el-table-column align="left" prop="order.createBy" label="账户名称" width="150" />

        <el-table-column align="left" prop="paymentCategory" label="计费账号" width="120" />
        <el-table-column align="left" prop="paymentCategory" label="计费模式" width="120" />
        <el-table-column align="left" prop="order.createTime" label="开始计费时间">
          <template #default="scope">
            {{ new Date(scope.row.order.createTime).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}
          </template>
        </el-table-column>
        <el-table-column align="left" prop="order.createTime" label="结束计费时间">
          <template #default="scope">
            {{ new Date(scope.row.order.createTime).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}
          </template>
        </el-table-column>
        <el-table-column align="left" prop="order.totalAmount" label="资源ID/编号" />
        <el-table-column align="left" prop="order.totalAmount" label="金额" />
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {useDebounceFn} from "@vueuse/core";
import {getBaseUrl} from "@/utils/format";
import {getRenterFinanceList} from "@/api/yunguan/resource/resInfo";

// 表格相关
const tableData = ref([]);
const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const responseData = ref({});
const value1 = ref([]);
const userList = ref([]);
// 搜索条件
const searchInfo = reactive({
  userId: "",
  // purchaser: "",
  status: "",
  productType: "",
  payType: "",
  startTime: "",
  endTime: "",
  sortField: "",
  sortOrder: "",
});

const exportTableReq = ref({
  begin: null,
  end: null,
  type:null
});

// 数据获取和处理
const fetchTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    keyword: "node", //表示从算力平台查询本节点所有租户的财务信息
    ...searchInfo
  };
  // 假设调用API获取数据，这里需要替换为实际的API请求
  const response = await getRenterFinanceList(params);
  if (response && response.data) {
    responseData.value = response;
    tableData.value = response.data.list;
    total.value = response.data.total;
  }
};

// 搜索功能
const onSubmit = () => {
  page.value = 1;
  fetchTableData();
};

// 重置搜索条件
const onReset = () => {
  searchInfo.userId = '';
  // searchInfo.purchaser = '';
  searchInfo.status = '';
  searchInfo.productType = '';
  searchInfo.payType = '';
  searchInfo.startTime = '';
  searchInfo.endTime = '';
  searchInfo.sortField = '';
  searchInfo.sortOrder = '';
  page.value = 1;
  fetchTableData();
};
const showExport = ref(false);
// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val;
  fetchTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  fetchTableData();
};
const exportTableHandle = () => {
  // 检查参数
  if (
      exportTableReq.value.begin === null ||
      exportTableReq.value.end === null ||
      exportTableReq.value.type === null
  ) {
    ElMessage.error("请输入完整条件");
    return;
  }
  const reqData = {
    ...exportTableReq.value,
    type: Number(exportTableReq.value.type)
  };

  fetch(`${getBaseUrl()}/order/jsManagerExport`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(reqData),
  })
      .then(async (resp) => {
        console.log(resp);
        if (resp.code === 7) {
          ElMessage.error(resp.msg);
          return;
        }
        return resp.blob();
      })
      .then((blob) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = "对账表.xlsx";
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        // 保存文件
        showExport.value = false;
        exportTableReq.value = {
          begin: null,
          end: null,
          type: null,
        };
        showExport.value = false;
      })
      .catch((err) => {
        // 处理上传失败的逻辑
        ElMessage.error(err);
      });
};

const debounceExportTableHandle = useDebounceFn(exportTableHandle, 200);

// 导出操作
const exportData = () => {
  // 执行导出操作
  ElMessage.success("导出成功");
};

// //获取选择用户
// const getSelectUserData = async () => {
//   // const params = {
//   //   page: 1,
//   //   pageSize: 20,
//   //   userName:""
//   // };
//   // const response = await getUserAdminList(params);
//   // if (response && response.data) {
//   //   userList.value = response.data.list;
//   // }
// };
// getSelectUserData();

onMounted(() => {
  fetchTableData();
});
</script>

<style scoped>
.gva-search-box {
  margin-bottom: 20px;
  display: flex;

  align-items: center;
}

.metal-types {
  margin-left: 10px;
  span{
    margin-right: 5px;
  }
}

.gva-table-box {
  margin-top: 20px;
}

.gva-btn-list {
  margin-bottom: 10px;
}

.gva-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>