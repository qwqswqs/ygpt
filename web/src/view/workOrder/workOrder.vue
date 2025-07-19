<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="概要">
          <el-input v-model="searchInfo.title" placeholder="概要"></el-input>
        </el-form-item>
        <!-- <el-form-item label="操作人员">
          <el-input
            v-model="searchInfo.operator"
            placeholder="操作人员"
          ></el-input>
        </el-form-item> -->
        <el-form-item label="类型">
          <el-select
              v-model="searchInfo.type"
              placeholder="请选择类型"
              clearable
          >
            <el-option label="其他" value="其他"></el-option>
            <el-option label="意见反馈" value="意见反馈"></el-option>
            <el-option label="投诉" value="投诉"></el-option>
            <el-option label="解决方案" value="解决方案"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary"  @click="onSubmit">
            查询
          </el-button>
          <el-button  @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary"  @click="addApiFunc('addApi')">
          新增
        </el-button>
        <el-button :disabled="!apis.length" @click="onDelete">
          删除
        </el-button>
      </div>
      <el-table
          :data="tableData"
          @sort-change="sortChange"
          @selection-change="handleSelectionChange"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column type="index" label="序号" width="55"></el-table-column>
        <el-table-column
            prop="title"
           align="left"
            label="工单类别"
            min-width="150"
        ></el-table-column>
        <el-table-column
            prop="type"
           align="left"
            label="所属租户id"
            min-width="100"
        ></el-table-column>

        <el-table-column
            prop="userName"
           align="left"
            label="工单描述"
            min-width="120"
        ></el-table-column>
        <el-table-column
            prop="phoneNumber"
           align="left"
            label="生成时间"
            min-width="120"
        ></el-table-column>
        <el-table-column
            prop="email"
           align="left"
            label="创建用户id"
            min-width="180"
        ></el-table-column>
        <el-table-column
            prop="workOrderStatus"
           align="left"
            label="处理用户id"
            min-width="100"
        >
          <!-- 修改为工单状态 -->
          <template #default="scope">
            <span>
              {{
                scope.row.workOrderStatus === "normal"
                    ? "工单完成"
                    : scope.row.workOrderStatus === "processing"
                        ? "处理中"
                        : scope.row.workOrderStatus === "abnormal"
                            ? "工单关闭"
                            : "未知状态"
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
            prop="workOrderStatus"
           align="left"
            label="完成情况"
            min-width="100"
        >
          <!-- 修改为工单状态 -->
          <template #default="scope">
            <span>
              {{
                scope.row.workOrderStatus === "normal"
                    ? "工单完成"
                    : scope.row.workOrderStatus === "processing"
                        ? "处理中"
                        : scope.row.workOrderStatus === "abnormal"
                            ? "工单关闭"
                            : "未知状态"
              }}
            </span>
          </template>
        </el-table-column>
‘        <el-table-column
          prop="workOrderStatus"
         align="left"
          label="状态"
          min-width="100"
      >
        <!-- 修改为工单状态 -->
        <template #default="scope">
            <span>
              {{
                scope.row.workOrderStatus === "normal"
                    ? "工单完成"
                    : scope.row.workOrderStatus === "processing"
                        ? "处理中"
                        : scope.row.workOrderStatus === "abnormal"
                            ? "工单关闭"
                            : "未知状态"
              }}
            </span>
        </template>
      </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" width="250">
          <template #default="scope">
            <el-button  type="primary" link @click="detailApiFunc(scope.row)">
              详情
            </el-button>
            <template v-if="scope.row.workOrderStatus !== 'normal'">
              <el-button

                  type="primary"
                  link
                  @click="editApiFunc(scope.row)"
              >
                处理
              </el-button>
            </template>
            <el-button
                type="primary"
                link
                @click="deleteApiFunc(scope.row)"
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

    <el-dialog
        v-model="dialogFormVisible"
        size="60%"
        :before-close="closeDialog"
        :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-between items-center">
          <span class="text-lg"></span>
          <div>
            <el-button @click="closeDialog"> 取 消 </el-button>
            <el-button  :disabled="!isChangeStatus" type="primary" @click="enterDialog"> 确 定 </el-button>
          </div>
        </div>
      </template>

      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="概要" v-if="isEditing">
          <el-input v-model="form.title" :disabled="!isEditing"></el-input>
        </el-form-item>
        <!-- <el-form-item label="概要">
          <el-select v-model="form.type">
            <el-option label="公告" value="公告"></el-option>
            <el-option label="通知" value="通知"></el-option>
          </el-select>
        </el-form-item> -->
        <el-form-item label="类型" v-if="isEditing" >
          <el-select v-model="form.type" :disabled="!isEditing">
            <el-option label="其他" value="其他"></el-option>
            <el-option label="意见反馈" value="意见反馈"></el-option>
            <el-option label="投诉" value="投诉"></el-option>
            <el-option label="解决方案" value="解决方案"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述" v-if="isEditing">
          <el-input v-model="form.des" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="诉求" v-if="isEditing">
          <el-input v-model="form.need" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="姓名" v-if="isEditing">
          <el-input v-model="form.userName" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="电话" v-if="isEditing">
          <el-input v-model="form.phoneNumber" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" v-if="isEditing">
          <el-input v-model="form.email" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="内容" v-if="isEditing">
          <el-input v-model="form.content" :disabled="!isEditing"></el-input>
        </el-form-item>
        <el-form-item label="状态" v-if="isEditing">
          <el-radio-group v-model="form.workOrderStatus" :disabled="!isEditing">
            <el-radio label="normal">工单完成</el-radio>
            <el-radio label="abnormal">工单关闭</el-radio>
            <el-radio label="processing">处理中</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="处理意见" prop="advice">
          <el-input v-model="form.advice"></el-input>
        </el-form-item>
        <!-- 处理结果 -->
        <el-form-item
            label="处理结果"
            prop="result"
            v-if="form.workOrderStatus === 'normal'"
        >
          <el-input v-model="form.result"></el-input>
        </el-form-item>
        <!-- 处理评价 -->
        <el-form-item
            label="处理评价"
            prop="evaluation"
            v-if="form.workOrderStatus === 'normal'"
        >
          <el-input v-model="form.evaluation"></el-input>
        </el-form-item>
        <!-- 评价内容 -->
        <el-form-item
            label="评价内容"
            prop="evaluationContent"
            v-if="form.workOrderStatus === 'normal'"
        >
          <el-input type="textarea" v-model="form.evaluationContent"></el-input>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

//编辑页面，详细页面，新增页面的弹窗应该加以修改

const apis = ref([]);
//这里应该是编辑或者新增页面的双向绑定的数据
const form = ref({
  title: "", // 概要
  type: "", // 类型
  des: "", // 描述
  need: "", // 诉求
  userName: "", // 姓名
  phoneNumber: "", // 电话
  email: "", // 邮箱
  content: "", // 内容
  workOrderStatus: "normal", // 状态，默认值为工单完成
  advice: "", // 处理意见
  result: "", // 处理结果
  evaluation: "", // 处理评价
  evaluationContent: "", // 评价内容
});

const type = ref("");
const rules = {
  title: [{ required: true, message: "请输入概要", trigger: "blur" }],
  type: [{ required: true, message: "请选择类型", trigger: "change" }],
  userName: [{ required: true, message: "请输入姓名", trigger: "blur" }],
  phoneNumber: [{ required: true, message: "请输入电话", trigger: "blur" }],
  email: [{ required: true, message: "请输入邮箱", trigger: "blur" }],
  content: [{ required: true, message: "请输入内容", trigger: "blur" }],
  workOrderStatus: [
    { required: true, message: "请选择状态", trigger: "change" },
  ],
  result: [
    {
      required: true,
      message: "请填写处理结果",
      trigger: "blur",
    },
  ],

};

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([
  {
    ID: 1,
    title: "系统维护通知", // 修改为概述
    type: "其他", // 修改为类别
    workOrderStatus: "processing", // 修改为工单状态
    userName: "张三", // 修改为用户名
    phoneNumber: "13800138000", // 修改为电话
    email: "13800138000@qq.com", // 修改为邮箱
    createdAt: "2024-09-01 10:00:00",
    content: "454",
  },
  {
    ID: 2,
    title: "版本更新公告", // 修改为概述
    type: "意见反馈", // 修改为类别
    workOrderStatus: "normal", // 修改为工单状态
    userName: "李四", // 修改为用户名
    phoneNumber: "13800138000", // 修改为电话
    email: "13800138000@qq.com", // 修改为邮箱
    createdAt: "2024-09-02 14:30:00",
    content: "454",
  },
  {
    ID: 3,
    title: "安全警告", // 修改为概述
    type: "投诉", // 修改为类别
    workOrderStatus: "abnormal", // 修改为工单状态
    userName: "王五", // 修改为用户名
    phoneNumber: "13800138000", // 修改为电话
    email: "13800138000@qq.com", // 修改为邮箱
    createdAt: "2024-09-03 09:15:00",
    content: "454",
  },
  {
    ID: 4,
    title: "年度总结报告", // 修改为概述
    type: "解决方案", // 修改为类别
    workOrderStatus: "normal", // 修改为工单状态
    userName: "赵六", // 修改为用户名
    phoneNumber: "13800138000", // 修改为电话
    email: "13800138000@qq.com", // 修改为邮箱
    createdAt: "2024-09-04 11:45:00",
    content: "454",
  },
  {
    ID: 5,
    title: "节假日放假通知", // 修改为概述
    type: "解决方案", // 修改为类别
    workOrderStatus: "processing", // 修改为工单状态
    userName: "钱七", // 修改为用户名
    phoneNumber: "13800138000", // 修改为电话
    email: "13800138000@qq.com", // 修改为邮箱
    createdAt: "2024-09-05 08:00:00",
    content: "454",
  },
]);

const searchInfo = ref({
  title: "", // 公告标题
  operator: "", // 操作人
  type: "", // 公告类型
});

const onReset = () => {
  searchInfo.value.title = "";
  searchInfo.value.operator = "";
  searchInfo.value.type = "";

  // 恢复为原始的表格数据
  tableData.value = [...originalTableData.value];
};

// 搜索

//列表数据和获取
const getTableData = async () => {
  //定义查询参数
  const params = {
    page: currentPage.value,
    pageSize: pageSize.value,
    // title: searchInfo.title,
    // type: searchInfo.type
  };
  //打印查询参数
  console.log(
      currentPage.value,
      pageSize.value,
      searchInfo.title,
      searchInfo.type
  );
  try {
    // 调用获取数据的接口（需要更改此接口）
    const table = await getNodeList(params);
    if (table.code === 0) {
      console.log(table.data);
      tableData.value = table.data.list || [];
      console.log(tableData.value, "64");
      total.value = table.data.total || 0;
      // 如果后端返回当前页码和每页大小，可以在这里更新
      // currentPage.value = table.data.page
      // pageSize.value = table.data.pageSize
    }
  } catch (error) {
    console.error("获取表格数据失败:", error);
  }
};

const originalTableData = ref([...tableData.value]); // 保存原始数据
//搜索功能
const onSubmit = async () => {
  // 重置分页
  currentPage.value = 1; // 假设 currentPage 是一个响应式引用
  pageSize.value = 10; // 假设 pageSize 是一个响应式引用

  // 定义查询参数
  const params = {
    page: currentPage.value,
    pageSize: pageSize.value,
    title: searchInfo.value.title, // 使用搜索标题
    // operator: searchInfo.value.operator, // 使用操作人
    type: searchInfo.value.type, // 使用公告类型
  };

  try {
    // 调用获取数据的接口（需要更改此接口）
    const table = await getNodeList(params);
    if (table.code === 0) {
      console.log(table.data);
      tableData.value = table.data.list || [];
      console.log(tableData.value, "64");
      total.value = table.data.total || 0;
    }
  } catch (error) {
    console.error("搜索数据失败:", error);
  }
};

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
};

const handleCurrentChange = (val) => {
  page.value = val;
};

// 排序
const sortChange = ({ prop, order }) => {
  if (prop) {
    if (prop === "ID") {
      prop = "id";
    }
    searchInfo.value.orderKey = toSQLLine(prop);
    searchInfo.value.desc = order === "descending";
  }
};

// 批量操作
const handleSelectionChange = (val) => {
  apis.value = val;
  // 提取选中项的 ID
  const selectedIds = val.map((item) => item.ID); // 假设 ID 是你需要的字段名
  console.log("选中的 ID:", selectedIds);
};

//批量删除
const onDelete = async () => {
  if (apis.value.length === 0) {
    ElMessage({
      type: "warning",
      message: "请先选择要删除的项",
    });
    return;
  }

  ElMessageBox.confirm("确定要删除所选项吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
      .then(async () => {
        try {
          // 假设你有一个删除接口 deleteNode
          const response = await deleteNode({ ids: apis.value }); // 传入选中的 ID 数组
          if (response.code === 0) {
            ElMessage({
              type: "success",
              message: "删除成功",
            });
            // 刷新数据或更新表格
            await getTableData(); // 重新获取数据
          } else {
            ElMessage({
              type: "error",
              message: response.message || "删除失败",
            });
          }
        } catch (error) {
          console.error("删除失败:", error);
          ElMessage({
            type: "error",
            message: "删除过程中发生错误",
          });
        }
      })
      .catch(() => {
        // 用户取消删除
        ElMessage({
          type: "info",
          message: "删除已取消",
        });
      });
};

//单个删除
const deleteApiFunc = async (row) => {
  ElMessageBox.confirm("此操作将永久删除, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
      .then(async () => {
        try {
          // 假设 row.ID 是要删除的项的 ID
          const response = await deleteNode(row.id); // 调用删除接口

          if (response.code === 0) {
            ElMessage({
              type: "success",
              message: "删除成功",
            });
            // 刷新数据或更新表格
            await getTableData(); // 重新获取数据
          } else {
            ElMessage({
              type: "error",
              message: response.message || "删除失败",
            });
          }
        } catch (error) {
          console.error("删除失败:", error);
          ElMessage({
            type: "error",
            message: "删除过程中发生错误",
          });
        }
      })
      .catch(() => {
        // 用户取消删除
        ElMessage({
          type: "info",
          message: "删除已取消",
        });
      });
};

// 弹窗相关
const apiForm = ref(null);
const initForm = () => {
  apiForm.value.resetFields();
  form.value = {
    path: "",
    apiGroup: "",
    method: "",
    description: "",
  };
};

const dialogTitle = ref("新增工单");
const dialogFormVisible = ref(false);
const openDialog = (key) => {
  switch (key) {
    case "addApi":
      dialogTitle.value = "新增工单";
      break;
    case "edit":
      dialogTitle.value = "编辑工单";
      break;
    case "detail":
      dialogTitle.value = "工单详情";
      break;

    default:
      break;
  }
  type.value = key;
  dialogFormVisible.value = true;
};

const isEditing = ref(false); // 是否为编辑模式
const isChangeStatus = ref(false); // 是否为修改状态
const closeDialog = () => {
  initForm();
  dialogFormVisible.value = false;
};

const editApiFunc = async (row) => {
  openDialog("edit");
  //绑定数据列表
  form.value = { ...row };
  isEditing.value = false; // 设置为编辑模式
  isChangeStatus.value = true; // 设置为编辑模式
};
const detailApiFunc = async (row) => {
  openDialog("detail");
  //绑定数据列表
  form.value = { ...row };
  isEditing.value = true;
  isChangeStatus.value = false; // 设置为编辑模式
};
const addApiFunc = async (row) => {
  openDialog("addApi");
  //绑定数据列表
  form.value = { ...row };
  isEditing.value = true;
  isChangeStatus.value = true; // 设置为编辑模式
};
// 新增和编辑提交表单
const enterDialog = async () => {
  apiForm.value.validate(async (valid) => {
    if (valid) {
      const formData = { ...form.value }; // 获取表单数据

      try {
        const response = await submitForm(formData); // 调用接口
        if (response.data.code === 0) {
          ElMessage({
            type: "success",
            message: type.value === "addApi" ? "新增成功" : "编辑成功",
            showClose: true,
          });
          closeDialog(); // 关闭对话框
          await getTableData(); // 刷新表格数据
        } else {
          ElMessage({
            type: "error",
            message: response.data.message || "操作失败",
            showClose: true,
          });
        }
      } catch (error) {
        console.error("操作失败:", error);
        ElMessage({
          type: "error",
          message: "操作过程中发生错误",
        });
      }
    }
  });
};
</script>

<style scoped lang="scss">
.warning {
  color: #dc143c;
}
</style>