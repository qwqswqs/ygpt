<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <div class="gva-search-box" style="margin-top: 16px">
        <div class="gva-btn-list">
          <el-form ref="searchForm" :model="searchInfo" label-width="70px"
                   :inline="true" style="display: flex;justify-content: space-between">
            <div>
              <el-form-item label="产品名称">
                <el-input v-model="searchInfo.keyword" clearable placeholder="请输入" style="width: 180px"/>
              </el-form-item>
              <el-form-item label="产品类型">
                <el-select v-model="searchInfo.subType" clearable placeholder="请选择" style="width: 180px">
                  <el-option value="5" label="模型"/>
                  <el-option value="6" label="数据集"/>
                </el-select>
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="searchInfo.status" clearable placeholder="请选择" style="width: 180px">
                  <el-option v-for="item in searchOptions.status" :value="item.value" :label="item.label"/>
                </el-select>
              </el-form-item>
              <el-form-item label="资源类型">
                <el-select v-model="searchInfo.eleType" clearable placeholder="请选择" style="width: 180px">
                  <el-option v-for="item in subTypeList" :value="item.value" :label="item.label"/>
                </el-select>
              </el-form-item>
            </div>
            <el-form-item>
              <el-button type="primary" @click="onSearch">
                查询
              </el-button>
              <el-button type="info" @click="onReset">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <!-- 表格选择 -->

        <!-- 搜索条件 -->

      </div>
      <!-- 表格 -->
      <div class="gva-table-box">
        <div style="margin-bottom: 16px;">
          <el-button type="primary" @click="addFunc('add')">
            新增
          </el-button>
          <el-button :disabled="!selections.length" type="grey" @click="onDelete">
            删除
          </el-button>
        </div>
        <!-- 表格具体内容 -->
        <div v-if="tableData && tableData.length > 0">
          <el-table :data="tableData" @selection-change="handleSelectionChange"
                    :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
            <el-table-column type="selection" width="55"/>
            <el-table-column align="left" class="text-overflow" label="产品名称" min-width="120" prop="product.name">
              <template #default="scope">
                <div class="text-overflow text-details" @click="handleDetailsClick(scope.row)">
                  {{ scope.row.product.name }}
                </div>
              </template>
            </el-table-column>
            <el-table-column align="left" label="产品类型" width="100">
              <template #default="scope">
                {{ scope.row.product.subtype == 5 ? '模型' : '数据集' }}
              </template>
            </el-table-column>
            <el-table-column align="left" label="描述" min-width="200" prop="product.description">
              <template #default="scope">
                <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                    scope.row.product.description
                  }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="价格（元）" min-width="120" prop="product.price"/>
            <el-table-column align="left" label="状态" min-width="120" sortable="custom">
              <template #default="{ row }">
                <div style="display: flex;align-items: center">
                  <div :class="handleTagType(row.product.status)"></div>
                  {{
                    searchOptions.status.find(item => item.value === row.product.status)?.label || '-'
                  }}
                </div>

              </template>
            </el-table-column>

            <el-table-column align="left" label="资源类型" min-width="120">
              <template #default="scope">
                {{ subTypeList.find(item => item.value == scope.row.sysType.parent_id)?.label }}
              </template>
            </el-table-column>
            <el-table-column align="left" label="二级类型" min-width="120" prop="sysType.name"/>
            <el-table-column align="left" fixed="right" label="操作" width="300">
              <template #default="scope">
                <el-button link type="text" class="text-blue-500" @click="editFunc(scope.row)">
                  编辑
                </el-button>
                <el-button
                    v-if="scope.row.product.isSynced === 0&&scope.row.product.auditStatus==1" type="text"
                    class="text-blue-500" link
                    @click="syncFunc(scope.row)"
                    :disabled="scope.row.product.isSynced===1"
                >
                  同步
                </el-button>
                <el-button link type="text" class="text-blue-500" @click="openDeleteDialog(scope.row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div v-else class="text-center">
          <h3>暂无数据</h3>
        </div>
        <div class="gva-pagination">
          <el-pagination
              :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
              layout="total,prev, pager, next,  sizes, jumper" @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
          />
        </div>
      </div>
    </div>


    <!-- 抽屉：增加/编辑 -->
    <el-drawer v-model="dialogFormVisible" :before-close="closeDialog" :show-close="false" size="1200px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span v-if="activeName==='all'" class="text-lg">{{ dialogTitle }}AI产品</span>
        <span v-else class="text-lg">{{ dialogTitle }}{{ activeName === "model" ? "模型" : "数据集" }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="dialogFormVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form ref="form" :model="editForm" :rules="editRules" label-width="100px">
        <!-- 基本信息 -->
        <el-form-item v-if="activeName==='all'" label="产品类型" label-width="90px" prop="productType">
          <div style="display: flex;background: rgba(242, 243, 245, 1);">
            <el-button
                v-for="type in productTypeList.subType"
                :key="type.value"
                :type="editForm.productType === type.value ? 'white' : 'key'"
                @click="editForm.productType = type.value"
            >
              {{ type.label }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="产品名称" label-width="90px" prop="product.name">
          <el-input style="width: 320px" v-model="editForm.product.name" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="资源类型" label-width="90px" prop="subType">
          <el-radio-group v-model="editForm.subType" @change="selectSubType(editForm.subType)">
            <el-radio
                v-for="subType in subTypeList"
                :key="subType.value"
                :label="subType.value"
                :border="false"

            >
              {{ subType.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="secondConfigList.length > 0" label="二级类型" label-width="90px" prop="secondType">
          <el-select style="width: 320px" v-model="editForm.secondType" placeholder="选择二级类型" @change="handleSecondClick()">
            <el-option v-for="second in secondConfigList" :key="second.id" :label="second.name" :value="second.id" />
          </el-select>
        </el-form-item>

        <el-form-item v-if="thirdConfigList.length > 0" label="三级类型" label-width="90px" prop="thirdType">
          <el-select style="width: 320px" v-model="editForm.thirdType" placeholder="选择三级类型">
            <el-option v-for="third in thirdConfigList" :key="third.id" :label="third.name" :value="third.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="来源" label-width="90px" prop="product.source">
          <el-radio-group v-model="editForm.product.source">
            <el-radio
                v-for="source in sourceList"
                :key="source.value"
                :label="source.value"
                :border="false"
            >
              {{ source.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="开放类别" label-width="90px" prop="product.openType">
          <el-select style="width: 320px" v-model="editForm.product.openType" placeholder="选择开放类别">
            <el-option v-for="openType in openTypeList" :key="openType.value" :label="openType.label" :value="openType.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="定制价格" label-width="90px" prop="product.price">
          <el-input-number
              v-model="editForm.product.price"
              :min="0"
              controls-position="horizontal"
              style="width: 120px"
          />
        </el-form-item>
        <el-form-item label="可否下载" label-width="90px" prop="product.isDownload">
          <el-radio-group v-model="editForm.product.isDownload">
            <el-radio :label="0">
              可以下载
            </el-radio>
            <el-radio :label="1">
              不可以下载
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <div v-if="dialogTitle==='新增'">
          <el-form-item label="上传类型"  label-width="90px" prop="uploadType">
            <div style="display: flex;background: rgba(242, 243, 245, 1);">
              <el-button
                  v-for="type in productTypeList.uploadType"
                  :key="type.value"
                  :type="editForm.uploadType === type.value ? 'white' : 'key'"
                  @click="editForm.uploadType = type.value"
              >
                {{ type.label }}
              </el-button>
            </div>
          </el-form-item>
          <el-form-item v-if="editForm.uploadType===1" label-width="90px"  label="附件" prop="attachment">
            <el-upload
                :action="uploadUrl"
                :file-list="editForm.files"
                @success="handleUploadSuccess"
                @change="handleChanged"

            >
              <el-button type="primary">
                上传文件
              </el-button>
            </el-upload>
          </el-form-item>
          <el-form-item v-if="editForm.uploadType===2"  label-width="90px" label="文件链接" prop="url">
            <el-input style="width: 320px" v-model="editForm.url" placeholder="请输入文件链接" />
          </el-form-item>
        </div>
        <div v-else-if="dialogTitle==='编辑'||dialogTitle==='详情'">
          <el-form-item v-if="editForm.productFiles!==''&&editForm.productFiles!=null">
            <!-- 基本信息 -->
            <div v-for="(product, productIndex) in editForm.productFiles" :key="productIndex" style="margin: 20px 0 ;width: 100%">
              <el-table :data="[product]" border style="width: 100%">
                <el-table-column v-if="product.name" prop="name" label="文件名称" />
                <el-table-column v-if="product.url" prop="url" label="文件链接">
                  <template #default="scope">
                    <el-link :href="scope.row.url" target="_blank" style="color: blue;">
                      {{ scope.row.url }}
                    </el-link>
                  </template>
                </el-table-column>
                <el-table-column v-if="product.size&&product.size!==0"align="left" label="文件大小" min-width="100">
                  <template #default="scope">
                    {{ formatFileSize(scope.row.size) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-form-item>
          <el-form-item v-else label="文件">
            <el-tag type="primary">
              无
            </el-tag>
          </el-form-item>
        </div>
        <el-form-item label="详细描述" label-width="90px"  prop="product.description">
          <el-input
              v-model="editForm.product.description" maxlength="300" placeholder="请输入详细描述" show-word-limit
              style="width: 1044px;" type="textarea"
          />
        </el-form-item>


      </el-form>
      <div v-if="dialogTitle !== '详情'" class="flex justify-end " style="margin-top: 24px;">
        <el-button type="info" @click="closeDialog">
          取 消
        </el-button>
        <el-button type="primary" :disabled="isConfirmButtonDisabled" @click="confirmProduct">
          确 定
        </el-button>
      </div>
    </el-drawer>

    <el-drawer v-model="dialogDetialVisible" :before-close="closeDialog" :show-close="false" size="720px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span  class="text-lg">产品详情</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="dialogDetialVisible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form ref="form" :model="editForm" :rules="editRules" label-width="100px">
        <!-- 基本信息 -->
        <el-form-item v-if="activeName==='all'" label="产品类型">
          <span>{{ getProductTypeLabel(editForm.productType) }}</span>
        </el-form-item>
        <el-form-item label="产品名称">
          <span>{{ editForm.product.name }}</span>
        </el-form-item>
        <el-form-item label="资源类型">
          <span v-if="editForm.sysType.parent_id===29">图像</span>
          <span v-if="editForm.sysType.parent_id===1">文本</span>
          <span v-if="editForm.sysType.parent_id===16">语音</span>
          <span v-if="editForm.sysType.parent_id===39">视频</span>

        </el-form-item>
        <el-form-item v-if="secondConfigList.length > 0" label="二级类型">
          <span>{{ editForm.sysType.name }}</span>
        </el-form-item>

        <el-form-item v-if="thirdConfigList.length > 0" label="三级类型" >
          <el-select style="width: 320px" v-model="editForm.thirdType" placeholder="选择三级类型">
            <el-option v-for="third in thirdConfigList" :key="third.id" :label="third.name" :value="third.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="来源">
          <span>{{ getSourceLabel(editForm.product.source) }}</span>
        </el-form-item>
        <el-form-item label="开放类别">
          <span>{{ getOpenTypeLabel(editForm.product.openType) }}</span>
        </el-form-item>
        <el-form-item label="定制价格" >
          <span>{{ editForm.product.price }}</span>
        </el-form-item>
        <el-form-item label="可否下载">
          <span>{{ editForm.product.isDownload === 0 ? '可以下载' : '不可以下载' }}</span>
        </el-form-item>
        <el-form-item label="详细描述">
          <span>{{editForm.product.description }}</span>
        </el-form-item>
        <el-form-item v-if="editForm.productFiles!==''&&editForm.productFiles!=null">
          <!-- 基本信息 -->
          <div v-for="(product, productIndex) in editForm.productFiles" :key="productIndex" style="width: 100%">
            <el-table :data="[product]"  style="width: 100%" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
              <el-table-column v-if="product.name" prop="name" label="文件名称" />
              <el-table-column v-if="product.url" prop="url" label="文件链接">
                <template #default="scope">
                  <el-link :href="scope.row.url" target="_blank" style="color: blue;">
                    {{ scope.row.url }}
                  </el-link>
                </template>
              </el-table-column>
              <el-table-column v-if="product.size&&product.size!==0"align="left" label="文件大小" min-width="100">
                <template #default="scope">
                  {{ formatFileSize(scope.row.size) }}
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-form-item>
        <el-form-item v-else label="文件">
          <el-tag type="primary">
            无
          </el-tag>
        </el-form-item>



      </el-form>
      <div v-if="dialogTitle !== '详情'" class="flex justify-end mt-4">
        <el-button type="info" @click="closeDialog">
          取 消
        </el-button>
        <el-button type="primary" :disabled="isConfirmButtonDisabled" @click="confirmProduct">
          确 定
        </el-button>
      </div>
    </el-drawer>

    <!-- 抽屉：审核 -->
    <el-dialog v-model="dialogAuditVisible" :before-close="closeDialog" :show-close="false" size="30%">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
        </div>
      </template>
      <el-form ref="form" :model="formAudit">
        <el-form-item label="审核状态" prop="status">
          <el-radio-group v-model="formAudit.status">
            <el-radio border :value="1">
              通过
            </el-radio>
            <el-radio border :value="-1">
              驳回
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div class="flex justify-end mt-4">
        <el-button @click="closeDialog">
          取 消
        </el-button>
        <el-button type="primary" @click="confirmAudit()">
          确 定
        </el-button>
      </div>
    </el-dialog>
    <el-drawer v-model="dialogFilesVisible" :before-close="closeDialog" :show-close="false" size="60%">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle
            }}{{ activeName === "model" ? "模型" : "数据集" }}</span>
        </div>
      </template>
      <!-- 这里是表单内容 -->
      <!-- ref表示表单的引用，model表示表单的数据，rules表示表单的校验规则，label-width表示表单项的label宽度 -->
      <el-form ref="form" :model="files" :rules="editRules" label-width="100px">
        <!-- 基本信息 -->
        <div v-for="(product, productIndex) in editForm.ProductFiles" :key="productIndex" style="margin: 20px 0">
          <el-table :data="[product]" border style="width: 100%">
            <el-table-column prop="name" label="文件名称" />
            <el-table-column prop="description" label="文件描述" />
            <el-table-column prop="url" label="文件地址" />
          </el-table>
        </div>
      </el-form>

      <div v-if="dialogTitle !== '详情'" class="flex justify-end mt-4">
        <el-button class="ml-2" @click="closeDialog">
          取 消
        </el-button>
        <el-button type="primary" @click="confirmProduct">
          确 定
        </el-button>
      </div>
    </el-drawer>
    <el-dialog v-model="forceDeleteDialogVisible" size="50%" title="强制删除" width="620px">
      <div>
        <el-form-item>
          <span class="text">有所选产品所在云管平台未对接算力平台,删除可能导致算力平台与云管平台产品信息不一致。</span>
          <span class="text-red-500" style="font-size: 16px; font-weight: bold;">
            是否强制删除该产品？
          </span>
        </el-form-item>
      </div>
      <div class="flex justify-end mt-4">
        <el-button class="ml-2" @click="forceDeleteDialogVisible = false">
          取 消
        </el-button>
        <el-button type="primary" @click="onForceDelete">
          确 定
        </el-button>
      </div>
    </el-dialog>
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
import {ref,computed } from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
//旧接口
import {
  getElementProductList,
  createProduct,
  updateProduct, deleteProductByIds, deleteProduct, syncElement,
} from "@/api/yunguan/element/element";

import {onMounted} from "vue";
import {getBaseUrl} from "@/utils/format";
import {getElemSysTypeList} from "@/api/yunguan/system/type";
import {deleteUser} from "@/api/system/user";

defineOptions({
  name: "AiProductViewPage",
});

const deleteVisible =ref(false)
const currentProductId = ref(null);
const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};
// 打开删除弹窗（统一入口）
const openDeleteDialog = ( row = null) => {
  currentProductId.value = row.product.id;
  deleteVisible.value = true;
};
const confirmDelete = async () => {
  try {
    let resp;
    resp = await deleteProduct({ id: currentProductId.value });

    // 严格判断接口成功状态
    if (resp.code === 0) {
      // 显示对应成功消息
      ElMessage({
        type: 'success',
        message:'删除成功'
      });
      deleteVisible.value = false;

      // 执行成功后的操作
      await getAllTableData(); // 刷新表格
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
const productTypeList = ref({
  subType: [
    {label: '模型', value: 5},
    {label: '数据集', value: 6},
  ],
  uploadType:[
    {label: '本地上传', value: 1},
    {label: '文件链接', value: 2},
  ]
});
// 使用 computed 缓存映射结果
const getProductTypeLabel = computed(() => {
  return (value) => {
    const item = productTypeList.value.subType.find(item => item.value === value);
    return item ? item.label : '未知';
  };
});
//文件上传
const visible = ref(false);
const productId = ref();
const files = ref([]);
// 上传相关
const uploadUrl = `${getBaseUrl()}/product/elem/upload`; // 替换为实际的上传接口
const searchOptions = {

  status: [
    {label: '未通过', value: -1},
    {label: '待审核', value: 0},
    {label: '未同步', value: 1},
    {label: '已同步', value: 2}
  ],
}
//定义数据区域

const handleTagType = (status) => {
  switch (status) {
    case 2:
      return 'blue'
    case -1:
      return 'danger'
    case 0:
      return 'orange'
    case 1:
      return 'grey'
    default:
      return 'info'
  }
}

//批量选中的数据
const selections = ref([]);
//上下架
const upAndDownCloudProductsFun = async (data) => {
  let params = {...data.product};
  params.status = data.product.status;
  params.id = data.product.id;
  const response = await upAndDownCloudProducts({id: params.id, status: params.status});
  if (response.code === 0) {
    ElMessage.success("操作成功");
    getAllTableData();
    const param = {
      id: params.id,
      status: params.status,
    };
    await notifyCloudProductStatus(param);
  } else {
    ElMessage.error("操作失败: " + (response.msg || "未知错误"));
  }
};
/*
 * 表格
 */
const tableData = ref([]);

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);

/*
 * 定义抽屉，可能
 */
//新增编辑
// const editForm = ref({})

const editForm = ref({
  product: {
    // price: '',
    // 其他属性可以在这里定义
  },
  sysType: {},
  user: {},
  files: [],
  uploadType: 1,
  subType: 29,
  url: "",
  size:0,
  openType: 0,
  // 其他字段...
});
//审核
const formAudit = ref({
  status:1,
});
const form = ref(null); // 确保 form 的引用正确初始化
//抽屉是否展示
const dialogFormVisible = ref(false);
const dialogDetialVisible = ref(false);
const dialogFilesVisible = ref(false);

const dialogAuditVisible = ref(false);
const dialogAddActivityVisible = ref(false);
const forceDeleteDialogVisible = ref(false);
const singleForceDeleteID = ref(0);
const dialogTitle = ref("model");
// 定义按钮禁用状态
const isConfirmButtonDisabled = ref(false);
const isUploaded = ref(false);
//审核操作值
const auditStatus = ref(0);

//定义筛选逻辑
const searchInfo = ref({
  //是否过期
  isExpired: "",
  status: null,
  //是否上架
  isActive: "",
  // 审核状态
  auditStatus: "",
  //子资源类型
  subType: "",
  //发布者
  publisher: "",
  // 关键字
  keyword: "",
});
//筛选框的具体选择项映射，统一在命名后加List
//是否过期
const isExpiredList = [
  {value: "", label: "全部"},
  {value: 0, label: "未过期"},
  {value: 1, label: "已过期"},
];
//是否上架
// const isActiveList = [
//     { value: '', label: '全部' },
//     { value: 0, label: '未上架' },
//     { value: 1, label: '已上架' }
// ]
//审核状态
const auditStatusList = [
  {value: "", label: "全部"},
  {value: 0, label: "待审核"},
  {value: 1, label: "未上架"},
  {value: 2, label: "已上架"},
  {value: -1, label: "已驳回"},
];
//子资源类型
const subTypeList = [
  // { value: "", label: "全部" },
  {value: 29, label: "图像"},
  {value: 1, label: "文本"},
  {value: 16, label: "语音"},
  {value: 39, label: "视频"},
];

const sourceList = [
  {value: 1, label: "节点"},
];
// 定义计算函数，根据值获取对应的标签
const getSourceLabel = computed(() => {
  return (value) => {
    const target = sourceList.find(item => item.value === value);
    return target ? target.label : '未知';
  };
});

const openTypeList = [
  {label: '不开放', value: 1},
  {label: '节点内可见可用', value: 2},
  {label: '开放到算力节点节点内用', value: 3},
  {label: '开放到算力节点且可外部用', value: 4},
];
// 定义计算函数，根据值获取对应的标签
const getOpenTypeLabel = computed(() => {
  return (value) => {
    const target = openTypeList.find(item => item.value === value);
    return target ? target.label : '未知';
  };
});
//发布者
const publisherList = [
  {value: "", label: "全部"},
  {value: 0, label: "平台"},
  {value: 1, label: "用户"},
  {value: 2, label: "节点"},
  {value: 3, label: "运营商"},
  {value: 4, label: "第三方"},
  {value: 5, label: "智算中心"},
];



// 定义更新列表表单
const auditFormData = ref({});

//定义获取活动设置表单
const activityForm = ref({});

// 城市列表的接收
const cities = ref([]);
const cityList = [];
const provinceList = [];

const getCityList = async () => {
  const res = await getSupplyOrDemandCityList();
  if (res.code === 0) {
    // 清空之前的列表
    cityList.length = 0;
    provinceList.length = 0;

    // 将返回的城市数据赋值
    cities.value = res.data.list;

    // 遍历数据
    cities.value.forEach((item) => {
      // 提取城市
      if (item.cities && Array.isArray(item.cities)) {
        cityList.push(...item.cities);
      }

      // 提取省份
      if (item.province) {
        provinceList.push(item.province);
      }
    });
  }
};

const addActivity = async () => {
  const params = {
    activeId: activityForm.value.activityId,
    productId: activityForm.value.id,
  };
  const res = await addProductActivity(params);
  if (res.code === 0) {
    ElMessage.success("添加活动成功");
    dialogAddActivityVisible.value = false;
    getAllTableData();
  } else {
    ElMessage.error("添加活动失败");
  }
};

//获取子类型配置
const productConfigList = ref([]);
const secondConfigList = ref([]);
const thirdConfigList = ref([]);
// 根据 ID 获取二级类型名称（添加空值检查）
const getSecondTypeName = computed(() => {
  return (id) => {
    // 检查数组是否存在且包含元素
    if (!secondConfigList.value || secondConfigList.value.length === 0) {
      return '加载中...';
    }

    const item = secondConfigList.value.find(item => item.id === id);
    return item ? item.name : '未知';
  };
});

const getProductConfig = async () => {
  const res = await getElemSysTypeList({type: 1});
  if (res.code === 0) {
    productConfigList.value = res.data.list;
  }
};
// 选择资源类型
const selectSubType = (subTypeName) => {
  editForm.value.subType = subTypeName;
  //打印subTypeName
  // 映射配置
  //子资源类型

  const selectedSubType = subTypeList.find(
      (subType) => subType.value === subTypeName
  );
  const label = selectedSubType ? selectedSubType.label : "未知类型"; // 默认值为 '未知类型

  // 查找对应的配置
  const selectedConfig = productConfigList.value.find(
      (config) => config.root.name === label
  );
  if (selectedConfig) {
    secondConfigList.value = selectedConfig.second || [];
    thirdConfigList.value = selectedConfig.third || [];
  } else {
    secondConfigList.value = [];
    thirdConfigList.value = [];
  }
};

const selectSource = (sourceName) => {
  editForm.value.source = sourceName;
  //打印subTypeName
};

const selectSecondType = (secondConfigId) => {


  // 根据二级类型ID查找对应的三级类型配置
  const selectedSecondConfig = thirdConfigList.value.filter(
      (config) => config.parent_id === secondConfigId.value
  );

  // 如果找到对应的三级类型配置，则更新三级类型列表
  if (selectedSecondConfig.length > 0) {
    thirdConfigList.value = selectedSecondConfig; // 直接赋值为找到的配置
  } else {
    thirdConfigList.value = []; // 如果未找到，清空列表
  }

  // 这里可以添加映射配置的逻辑
  // 例如：映射到其他变量或处理其他状态
};
//
const handleSubTypeClick = async (subType) => {
  editForm.value.subType = subType.value;
  //清空二级类型
  editForm.value.secondType = "";
  //清空三级类型
  editForm.value.thirdType = "";
  selectSubType(subType.value);
};
const handleSourceClick = async (source) => {
  editForm.value.product.source = source.value;
  selectSource(source.value);
};

const secondConfigId = ref();
const handleSecondClick = async () => {
  secondConfigId.value = editForm.value.secondType;
  //打印secondConfigId
  selectSecondType(secondConfigId);
};

/*
 * 产品分类
 */

const editRules = ref({
  product : {
    isDownload: [{ required: true, message: "请选择是否可下载", trigger: "blur" }],
    price: [{ required: true, message: "价格不能为空", trigger: "blur" }],
    source: [{ required: true, message: "请选择价格", trigger: "blur" }],
    openType : [{ required: true, message: "请选择开放类别", trigger: "blur" }],
    name: [{ required: true, message: "请输入标题", trigger: "blur" }],
    description: [{ required: true, message: "请输入详细描述", trigger: "blur" }],
  },
  uploadType: [{ required: true, message: "请选择附件类型", trigger: "blur" }],
  productType : [{ required: true, message: "请选择产品类型", trigger: "blur" }],
  subType: [{ required: true, message: "请选择产品类型", trigger: "blur" }],
  secondType: [{ required: true, message: "请选择标签", trigger: "blur" }],
  thirdType: [{ required: true, message: "请选择标签", trigger: "blur" }],
  startDate: [{ required: true, message: "请选择", trigger: "blur" }],
  attachment: [{required: true, message: "请上传文件", trigger: "blur"}],
  url: [{required: true, message: "请输入文件链接", trigger: "blur"}],
  endDate: [{ required: true, message: "请选择", trigger: "blur" }],
});

//最高级页面切换（供给需求）
const activeName = ref("all"); // 默认选中资源

// 搜索
const onSearch = () => {
  page.value = 1;
  pageSize.value = 10;
  getAllTableData();
};

// 重置搜索条件
const onReset = () => {
  searchInfo.value = {};
  getAllTableData();
};

// 批量选中操作
const handleSelectionChange = (val) => {
  selections.value = val;
  //打印selections
};
// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getAllTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getAllTableData();
};
// 定义 editForm 的初始结构
const getEditFormDefault = () => ({
  product: {},
  sysType: {},
  user: {},
  files: [],
  uploadType: 1,
  url: "",
  size: 0, // 确保 size 类型为 number
});
const addFunc = function (r) {
  auditFormData.value = {...r};
  IsCreated.value = true;
  // 使用 getEditFormDefault 方法确保类型一致
  editForm.value = getEditFormDefault();
  openDialog("add");
};

const handleDetailsClick = (r) => {
  IsCreated.value = false;
  editForm.value = {...r};
  //后端返回的数据与 editForm.value 的数据类型可能不一致，所以需要手动对应
  editForm.value.subType = r.sysType.parent_id
  selectSubType(editForm.value.subType);
  editForm.value.secondType = r.sysType.id
  editForm.value.product.source = r.product.source;
  editForm.value.startDate = r.product.startTime;
  editForm.value.endDate = r.product.endTime;
  editForm.value.productType = r.product.subtype;
  openDialog("show");
};

const handleFilesClick = (r) => {
  IsCreated.value = false;
  editForm.value = {...r};
  openDialog("file");
};

const editFunc = function (r) {
  IsCreated.value = false;
  editForm.value = {...r};
  //后端返回的数据与 editForm.value 的数据类型可能不一致，所以需要手动对应
  editForm.value.subType = r.sysType.parent_id
  selectSubType(editForm.value.subType);
  editForm.value.secondType = r.sysType.id
  editForm.value.product.source = r.product.source;
  editForm.value.startDate = r.product.startTime;
  editForm.value.endDate = r.product.endTime;
  editForm.value.productType = r.product.subtype;
  openDialog("edit");
};

const auditFunc = function (r) {
  auditFormData.value = r.product;
  openDialog("audit");
};
const syncFunc = function (r) {
  syncById(r.product.id);
};
const addActivityFunc = function (r) {
  activityForm.value = r.product;
  openDialog("addActivity");
};
// 单个删除
const deleteFunc = (row) => {
  auditFormData.value = {...row}; // 使用展开运算符复制行数据
  deleteById();
};

const openDialog = (key) => {
  switch (key) {
    case "add":
      editForm.value.product.source = 1;
      dialogTitle.value = "新增";
      dialogFormVisible.value = true;

      break;
    case "edit":
      dialogTitle.value = "编辑";
      dialogFormVisible.value = true;
      break;
    case "show":
      dialogTitle.value = "详情";
      dialogDetialVisible.value = true;
      break;
    case "audit":
      dialogTitle.value = "审核";
      dialogAuditVisible.value = true;
      break;
    case "addActivity":
      dialogTitle.value = "添加活动";
      dialogAddActivityVisible.value = true;
      break;
    case "file":
      dialogTitle.value = "文件列表";
      dialogFilesVisible.value = true;
      break;
    default:
      break;
  }
};

const closeDialog = () => {
  // 检查表单引用是否存在并重置
  // if (form.value) {
  //   form.value.resetFields(); // 重置主表单
  // }
  // 重置其他状态
  isUploaded.value = false;
  isConfirmButtonDisabled.value = false;
  auditStatus.value = 0; // 将审核状态重置为初始值
  // specs.value = []; // 清空规格数据
  dialogFormVisible.value = false; // 关闭主表单抽屉
  dialogAuditVisible.value = false; // 关闭审核抽屉
  dialogAddActivityVisible.value = false; // 关闭添加活动抽屉
  dialogFilesVisible.value = false; // 关闭规格抽屉
};

const activeNameTypeMap = {
  'all': 0,
  'model': 5,
  'dataset': 6,
}

// 初始化
const getAllTableData = async () => {
  // if (searchInfo.value.type === undefined) {
  //   searchInfo.value.type = 1; // 设置默认值
  // }
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    subType: searchInfo.value.subType,
  };
  if (searchInfo.value.status !== "") {
    params.status = searchInfo.value.status;
  }
  if (searchInfo.value.keyword !== "") {
    params.keyword = searchInfo.value.keyword;
  }
  if (searchInfo.value.eleType !== null) {
    params.eleType = searchInfo.value.eleType;
  }


  // if (searchInfo.value.isActive !== undefined) {
  //   params.status = searchInfo.value.isActive
  // }

  const table = await getElementProductList(params);
  // const table = await getProductListByCondition({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    // 处理数据

    // 确保 table.data 和 table.data.list 存在
    if (table.data && Array.isArray(table.data.list)) {
      tableData.value = table.data.list;
      total.value = table.data.total;
      page.value = table.data.page;
      pageSize.value = table.data.pageSize;
    } else {
      tableData.value = []; // 设置为空数组
      total.value = 0; // 总数为 0
    }
  } else {
    ElMessage.error(table.msg);
    // tableData.value = [];
  }
};


// 排序

// 批量删除
const onDelete = async () => {
  //打印selections
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    //打印selections
    const ids = selections.value.map((item) => item.product.id);
    //打印ids
    const res = await deleteProductByIds({ids: ids});
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: res.msg,
      });
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--;
      }
      getAllTableData();
    }
    if (res.code === -1) {
      forceDeleteDialogVisible.value = true;
    }
  });
};
const onForceDelete = async () => {
  if (singleForceDeleteID.value !== 0) {
    forceDeleteDialogVisible.value = false;
    const res = await deleteProduct({id: singleForceDeleteID.value, isForce: true});
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "删除成功!",
      });
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--;
      }
    }
    singleForceDeleteID.value = 0;
  } else {
    forceDeleteDialogVisible.value = false;
    const ids = selections.value.map((item) => item.product.id);
    //打印ids
    await deleteProductByIds({ids: ids, isForce: true});
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--;
    }
  }
  getAllTableData();
};

// 单个删除
const deleteById = async () => {
  //打印id
  //获取row的ID
  // const id = row.product.id
  //打印id
  ElMessageBox.confirm("此操作将永久删除该产品, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    //获取row的ID
    const id = auditFormData.value.product.id;
    //打印id
    const res = await deleteProduct({id: id});
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: "删除成功!",
      });
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--;
      }
      getAllTableData();
    }
    if (res.code === -1) {
      forceDeleteDialogVisible.value = true;
      singleForceDeleteID.value = auditFormData.value.product.id;
    }
  });
};
const handleUploadSuccess = (response, file, fileList) => {
  if (response.code === 0) {
    isConfirmButtonDisabled.value = false;
    isUploaded.value = true;
    ElMessage({
      message: '文件上传成功',
      type: 'success'
    });
    // // 处理上传成功后的逻辑，例如更新文件列表或显示上传结果
    editForm.value.fileName = response.data.file.Filename;
    editForm.value.size = response.data.file.Size;
    console.log(isConfirmButtonDisabled.value);
  } else {
    ElMessage({
      message: '文件上传失败',
      type: 'error'
    });
  }
};
const handleChanged= (file) => {
  //上传开始和结束都会调用
  //如果上传成功，则直接返回，失败会继续禁用确定按钮
  if (isUploaded.value) {
    return ;
  }
  isConfirmButtonDisabled.value = true;
  console.log(isConfirmButtonDisabled.value);
}
const IsCreated = ref(true);
//同步产品
const syncById = async (id) =>  {
  //打印id
  const res = await syncElement({id: id})
  if (res.code === 0) {
    ElMessage({
      message: '同步成功',
      type: 'success'
    })
  }
  getAllTableData();
};
//新增和编辑的接口
const confirmProduct = async () => {
  if (activeName.value === 'model') {
    editForm.value.productType = 5;
  } else if (activeName.value === 'dataset') {
    editForm.value.productType = 6;
  }
  if (form.value) {
    await form.value.validate(); // 验证表单
    const param = {
      description: editForm.value.product.description,
      id: editForm.value.product.id,
      endTime: editForm.value.product.endTime,
      headImage: "string",
      name: editForm.value.product.name,
      price: Number(editForm.value.product.price),
      source: editForm.value.product.source,
      specJson: JSON.stringify(editForm.value.otherLabel),
      startTime: editForm.value.product.startTime,
      subType: editForm.value.productType,
      type: 2,
      typeId: editForm.value.secondType,
      user: editForm.value.product.user,
      userEmail: editForm.value.product.userEmail,
      userPhone: editForm.value.product.userPhone,
      validTime: editForm.value.product.validTime,
      isDownload: editForm.value.product.isDownload,
      files: editForm.value.files,
      url: editForm.value.url,
      size: editForm.value.size,
      openType:editForm.value.product.openType,
      fileName: editForm.value.fileName,
    };
    let res = null;
    if (IsCreated.value) res = await createProduct(param);
    else {
      param.auditStatus=0;
      param.status=0;
      res = await updateProduct(param);
    }
    if (res.code === 0) {
      ElMessage({
        type: "success",
        message: IsCreated.value ? "添加成功" : "修改成功",
        showClose: true,
      });
      productId.value = res.data.productId;
      visible.value = true;
    }
    getAllTableData();
    closeDialog();
  }
};

//审核结果的接口
const confirmAudit = async () => {
  try {
    const req = {
      ...auditFormData.value,
      status:formAudit.value.status,
    };
    req.auditStatus=formAudit.value.status;
    const response = await updateProduct(req);
    if (response.code === 0) {
      ElMessage.success("审核成功");
      closeDialog();
      await getAllTableData();
    } else {
      ElMessage.error("审核失败: " + (response.msg || "未知错误"));
    }
  } catch (error) {
    ElMessage.error("审核时发生错误: " + error.message);
  }
  //通知审核结果
  const param = {
    id: auditFormData.value.id,
    status: formAudit.value.auditStatus,
  };
  await notifyCloudProductStatus(param);
};


//初始化时获取数据
onMounted(() => {
  getAllTableData();
  getProductConfig();
});
// 添加计算文件大小的方法
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
</script>

<style lang="scss" scoped>
.warning {
  color: #dc143c;
}

.tabs {
  margin-bottom: 30px;
}

.form-section {
  border: 1px solid #dcdfe6;
  /* 边框颜色 */
  border-radius: 4px;
  /* 圆角 */
  padding: 16px;
  /* 内边距 */
  margin-bottom: 20px;
  /* 下外边距 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  /* 阴影 */
  background-color: #fff;
  /* 背景色 */
}

.button-group {
  display: flex;
  gap: 10px;
  /* Space between buttons */
}

.el-button.active {
  background-color: #409eff;
  /* Active button color */
  color: white;
  /* Text color when active */
}

.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-details {
  color: #77b2f3;
  cursor: pointer;
}
.active {
  color:  rgba(22, 93, 255, 1)!important;
  background-color:rgba(255, 255, 255, 1)!important;
  margin: 2px;
}

.blue {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background:rgba(22, 93, 255, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
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

.orange {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(255, 125, 0, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}

.grey {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background:rgba(134, 144, 156, 1);
  border-radius: 50%;
  order: 0;
  flex-shrink: 0;
  margin-right: 8px;
}

</style>
