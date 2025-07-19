<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">    <!-- 条件筛选栏 -->
    <el-tabs v-model="activeTab" @tab-click="(tab) => handleTabClick(tab)">
      <el-tab-pane v-for="item in tabList" :key="item.name" :label="item.label" :name="item.name" />
    </el-tabs>
    <div class="gva-search-box" >
      <div class="gva-btn-list">
        <el-form ref="searchForm" label-position="left" :model="searchInfo" :inline="true" style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="资源类型">
              <el-select v-model="searchInfo.subType" style="width: 240px" clearable placeholder="请选择">
                <el-option v-for="item in searchOptions.subType" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <!--        <el-form-item label="是否过期">-->
            <!--          <el-select v-model="searchInfo.isExpired" style="width: 200px" clearable placeholder="请选择">-->
            <!--            <el-option v-for="item in searchOptions.isExpired" :key="item.value" :label="item.label" :value="item.value" />-->
            <!--          </el-select>-->
            <!--        </el-form-item>-->
            <el-form-item label="状态">
              <el-select v-model="searchInfo.status" style="width: 240px" clearable placeholder="请选择">
                <el-option v-for="item in searchOptions.status" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <!--        <el-form-item label="来源">-->
            <!--          <el-select v-model="searchInfo.source" style="width: 200px" clearable placeholder="请选择">-->
            <!--            <el-option v-for="item in searchOptions.source" :key="item.value" :label="item.label" :value="item.value" />-->
            <!--          </el-select>-->
            <!--        </el-form-item>-->
            <el-form-item label="搜索">
              <el-input v-model="searchInfo.keyword" style="width: 240px" placeholder="请输入关键词" clearable />
            </el-form-item>
          </div>
          <el-form-item>
            <el-button class="button-gap" type="primary" @click="onSearch">
              查询
            </el-button>
            <el-button class="button-gap" type="info" @click="onReset">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <!-- 表格选择 -->

      <!-- 搜索条件 -->

    </div>

    <!-- 表格栏 -->
    <div class="gva-table-box">
      <!-- 表格操作按钮 -->
      <div style="margin-bottom: 16px">
        <el-button type="primary"  @click="addFunc">
          新增
        </el-button>
        <el-button type="grey" :disabled="!selectionList.length" @click="openDeleteDialog('batch')">
          删除
        </el-button>
      </div>
      <!-- 表格 -->
      <div v-if="tableData && tableData.length > 0">
        <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
          <el-table-column type="selection" width="55" />
          <!--          <el-table-column align="left" label="序号" type="index" min-width="100" />-->
          <el-table-column align="left" label="标题" min-width="200">
            <template #default="{ row }">
              <!--              点击跳转详情页-->
              <!--              <div class="text-overflow text-detail" @click="handleDetailClick(row)">-->
              <!--                {{ row.product.name }}-->
              <!--              </div>-->
              <div class="text-overflow">
                {{ row.product.name }}
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="描述" min-width="200">
            <template #default="{ row }">
              <el-tooltip
                  effect="light"
                  :content="row.product.description"
                  placement="top"
                  :open-delay="300"
                  :popper-options="{
        modifiers: [
          {
            name: 'computeStyle',
            options: {
              customStyles: {
                width: '200px',
                height: '102px',
                overflow: 'hidden',
                whiteSpace: 'normal',
                wordBreak: 'break-all'
              }
            }
          }
        ]
      }"
              >
                <template #content>
                  <div class="custom-tooltip-content">
                    {{ row.product.description }}
                  </div>
                </template>
                <div class="text-overflow">
                  {{ row.product.description }}
                </div>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column align="left" label="类型" min-width="100" sortable="custom">
            <template #default="{ row }">
              <div class="text-overflow">
                {{
                  searchOptions.subType.find(item => item.value === row.product.subtype)?.label || '-'
                }}
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="状态" min-width="100" sortable="custom">
            <template #default="{ row }">
              <div style="display: flex;align-items: center">
                <div :class="handleTagType(row.product.status)"> </div>
                <div v-if="row.product.status!=-1" class="text-overflow">
                  {{
                    searchOptions.status.find(item => item.value === row.product.status)?.label || '-'
                  }}
                </div>
                <el-button @click="dialogVisible=true;feedbackString=row.product.feedback" v-if="row.product.status==-1" type="text">{{searchOptions.status.find(item => item.value === row.product.status)?.label}}</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="操作" fixed="right" width="300">
            <template #default="{ row }">
              <el-button  class="text-blue-500"  type="text" link @click="editFunc(row)">
                编辑
              </el-button>
              <el-button  class="text-blue-500" v-if="row.product.status == 1" type="text" link @click="SyncSupplyDemand(row)">
                同步
              </el-button>
              <el-button class="text-blue-500" type="text" link @click="openDeleteDialog('single', row)">
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
            :current-page="searchInfo.pageNum"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 20, 30, 50, 100]"
            :total="searchInfo.total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 抽屉：增加/编辑 -->
    <el-drawer :show-close="false" v-model="drawerFormVisible" style="min-width: 500px;" size="1200px" :before-close="closeDrawer">

      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ drawerTitle }}</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="drawerFormVisible = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form label-position="left" ref="drawerFormRef" :model="drawerForm" :rules="rules">
        <el-collapse v-model="activeNames">
          <el-collapse-item isActive name="1">
            <template #title>
              <div :style="{ marginLeft: '10px' }">基本信息</div>
            </template>
            <div class="form-section">
              <el-form-item label="资源类型" label-width="100px" prop="subType">
                <div style="display: flex;background: rgba(242, 243, 245, 1);">
                  <el-button
                      type="key"
                      v-for="subtype in searchOptions.subType"
                      :key="subtype.value"
                      :class="{ active: drawerForm.subType === subtype.value }"
                      @click="drawerForm.subType = null;
                        drawerForm.subTypeId = null;
                        drawerForm.subType1Id = null;
                        drawerForm.subType2Id = null;
                        drawerForm.subType = subtype.value;
                        drawerForm.subTypeId = subtypeDesc[subtype.value].id;
                        drawerFormRef.clearValidate('subType');"
                  >
                    {{ subtype.label }}
                  </el-button>
                </div>
              </el-form-item>
              <el-form-item
                  v-if="drawerForm.subTypeId && subtypeMap[drawerForm.subTypeId] && subtypeMap[drawerForm.subTypeId].length > 0"
                  label-width="90px"
                  style="padding-left: 10px"
                  :label="subtypeDesc[drawerForm.subType].subLabel1"
                  prop="subType1Id"
              >
                <el-select
                    v-model="drawerForm.subType1Id"
                    style="width: 320px;"
                    placeholder="请选择"
                    clearable
                >
                  <el-option v-for="item in subtypeMap[drawerForm.subTypeId]" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
              </el-form-item>
              <el-form-item
                  v-if="drawerForm.subType1Id && subtypeMap[drawerForm.subType1Id] && subtypeMap[drawerForm.subType1Id].length > 0"
                  label-width="90px"
                  style="padding-left: 10px"
                  :label="subtypeDesc[drawerForm.subType].subLabel2"
                  prop="subType2Id"
              >
                <el-select
                    v-model="drawerForm.subType2Id"
                    style="width: 320px;"
                    placeholder="请选择"
                    clearable
                >
                  <el-option v-for="item in subtypeMap[drawerForm.subType1Id]" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
              </el-form-item>
              <el-form-item label="区域地址" label-width="100px" prop="cityId">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <!-- 单选按钮组 -->
                  <el-radio-group v-model="drawerForm.cityType" @change="handleIdChange">
                    <el-radio :label="0">不限</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <!-- 级联选择器始终显示 -->
                  <el-cascader
                      style="width: 320px;"
                      :options="options"
                      v-model="drawerForm.cityId"
                      @focus="drawerForm.cityType = 1"
                      @change="handleChange"
                      v-if="drawerForm.cityType === 1"
                  />
                </div>
              </el-form-item>
              <!-- 计费计价字段 -->
              <el-form-item label="计费计价" label-width="100px" prop="price">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <!-- 单选按钮组 -->
                  <el-radio-group v-model="drawerForm.priceType" @change="handlePriChange">
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <!-- 输入框始终显示 -->
                  <el-input
                      v-model="drawerForm.price"
                      style="width: 320px;"
                      placeholder="示例：10.00元/月/台"
                      clearable
                      @focus="drawerForm.priceType = 1"
                      @change="handleItemChange('price')"
                      v-if="drawerForm.priceType === 1"
                  />
                </div>
              </el-form-item>

              <!-- 供需数量字段 -->
              <el-form-item label="供需数量" label-width="100px" prop="number">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <!-- 单选按钮组 -->
                  <el-radio-group v-model="drawerForm.numberType" @change="handleNumChange">
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <!-- 输入框始终显示 -->
                  <el-input
                      v-model="drawerForm.number"
                      style="width: 320px;"
                      placeholder="示例：9台"
                      clearable
                      @focus="drawerForm.numberType = 1"
                      @change="handleItemChange('number')"
                      v-if="drawerForm.numberType === 1"
                  />
                </div>
              </el-form-item>
              <!-- 租赁时间字段 -->
              <el-form-item label="租赁时间" label-width="100px" prop="rentalTime">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <!-- 单选按钮组 -->
                  <el-radio-group v-model="drawerForm.rentalTimeType" @change="handleTimeChange">
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <!-- 日期选择器组 -->
                  <div v-if="drawerForm.rentalTimeType === 1" style="display: flex; gap: 15px;">
                    <el-date-picker
                        v-model="drawerForm.rentalStartTime"
                        placeholder="请选择租赁开始时间"
                        style="max-width: 200px;"
                        type="datetime"
                        @change="handleRentalStartTimeChange"
                        :picker-options="{
      disabledDate: (time) => time.getTime() < Date.now() - 8.64e7
    }"
                    />
                    <el-date-picker
                        v-model="drawerForm.rentalEndTime"
                        placeholder="请选择租赁结束时间"
                        style="max-width: 200px;"
                        type="datetime"
                        @change="handleRentalEndTimeChange"
                        :picker-options="rentalEndOptions"
                    />
                  </div>
                </div>
              </el-form-item>

              <!-- 信息有效期字段 -->
              <el-form-item label="信息有效期" label-width="100px" prop="time">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <!-- 单选按钮组 -->
                  <el-radio-group v-model="drawerForm.timeType" @change="handleTimeTypeChange">
                    <el-radio :label="0">默认</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <!-- 日期选择器组 -->
                  <div v-if="drawerForm.timeType === 1" style="display: flex; gap: 15px;">
                    <el-date-picker
                        v-model="drawerForm.startTime"
                        placeholder="请选择信息发布"
                        style="max-width: 200px;"
                        type="datetime"
                        @change="handleInfoStartTimeChange"
                        :picker-options="{
      disabledDate: (time) => time.getTime() < Date.now() - 8.64e7
    }"
                    />
                    <el-date-picker
                        v-model="drawerForm.endTime"
                        placeholder="请选择信息截止时间"
                        style="max-width: 200px;"
                        type="datetime"
                        @change="handleInfoEndTimeChange"
                        :picker-options="infoEndOptions"
                    />
                  </div>
                </div>
              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item name="2">
            <template #title>
              <div :style="{ marginLeft: '10px' }">详情描述</div>
            </template>
            <div class="form-section">
              <el-form-item label="标题" label-width="100px" prop="name">
                <el-input
                    v-model="drawerForm.name"
                    style="width: 320px;"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="描述" label-width="100px" prop="description">
                <el-input
                    v-model="drawerForm.description"
                    style="width: 1016px;"
                    type="textarea"
                    :rows="2"
                    maxlength="3000"
                    show-word-limit
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="规格" l label-width="90px" style="padding-left: 10px" >
                <div>
                  <el-select
                      v-model="currentKey"
                      style="width: 200px; margin-right: 10px;"
                      allow-create
                      filterable
                      clearable
                      placeholder="请选择"
                  >
                    <el-option
                        v-for="option in keys"
                        :key="option.name"
                        :label="option.name + '（单位:' + option.unit + '）'"
                        :value="option.name"
                        @change="unit = option.unit"
                    />
                  </el-select>
                  <el-input
                      v-model="currentValue"
                      style="width: 200px; margin-right: 10px;"
                      clearable
                      placeholder="请输入"
                  />
                  <el-button type="primary" @click="addKeyValue">
                    添加
                  </el-button>
                  <el-table :data="keyValuePairs" style="margin-top: 20px;">
                    <el-table-column label="规格名称" prop="key" />
                    <el-table-column label="规格值" prop="value" />
                    <el-table-column label="操作">
                      <template #default="{ row, $index }">
                        <el-button type="primary" @click="removeKeyValue($index)">
                          删除
                        </el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>

              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item name="3">
            <template #title>
              <div :style="{ marginLeft: '10px' }">联系方式</div>
            </template>
            <div class="form-section">
              <el-form-item label="联系人" label-width="100px" prop="contact">
                <el-input
                    v-model="drawerForm.contact"
                    style="width: 320px"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="电话" label-width="100px" prop="phone">
                <el-input
                    v-model="drawerForm.phone"
                    style="width: 320px"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="邮箱" label-width="90px" style="padding-left: 10px" prop="email">
                <el-input
                    v-model="drawerForm.email"
                    style="width: 320px"
                    placeholder="请输入（可选）"
                    clearable
                />
              </el-form-item>
            </div>
          </el-collapse-item>
        </el-collapse>



      </el-form>
      <div class="flex justify-end mt-4">
        <el-button type="info" @click="closeDrawer">
          取 消
        </el-button>
        <el-button type="primary" @click="confirmForm">
          确 定
        </el-button>
      </div>
    </el-drawer>

    <!--    暂时不用dialog(审暂时核使用drawer)-->
    <el-dialog v-model="dialogVisible" :show-close="true" style="min-width: 500px" title="审核意见" width="40%" center>
      <span>{{feedbackString!=''?feedbackString:'无'}}</span>
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
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import {
  createProduct,
  deleteProduct, deleteProductByIds,
  getProductConfigList,
  getProductList,
  getSupplyOrDemandCityList, SyncSupplyDemandInfo, updateProduct
} from "@/api/yunguan/product/productSupply"
import {getCategoryList} from "@/api/yunguan/system/type";
import { provinceAndCityData ,  CodeToText, TextToCode } from 'element-china-area-data'

/*********
 * 条件筛选
 *********/
const isActive = ref(true);

// 通过 computed 生成展开项数组
const activeNames = computed(() => {
  return isActive.value ? ['1', '2', '3'] : [];
});
// 搜索表单
const searchForm = ref(null)
// 已激活标签
const activeTab = ref(3)
// 标签列表
const tabList = ref([
  // { label: '全部', name: 0 },
  { label: '供给', name: 3 },
  { label: '需求', name: 4 },
])
const feedbackString =ref('');
// 条件选项
const searchOptions = ref({
  subType: [
    { label: '算力', value: 7 },
    { label: '硬件', value: 8 },
    { label: '机柜', value: 9 },
    { label: '带宽', value: 10 },
  ],
  isExpired: [
    { label: '未过期', value: 0 },
    { label: '已过期', value: 1 },
  ],
  status: [
    { label: '已同步', value: -2 },
    { label: '未通过', value: -1 },
    { label: '未审核', value: 0 },
    { label: '未同步', value: 1 },
  ],
  source: [
    { label: '平台', value: 1 },
    { label: '第三方', value: 2 },
    { label: '运营商', value: 3 },
    { label: '智算中心', value: 4 }
  ]
})

const handleTagType = (status) => {
  switch (status) {
    case -2:
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

// 搜索条件
const searchInfo = ref({
  page: 1,
  pageSize: 10,
  total: 0,
  sortBy: null,
  isDesc: null,
  keyword: null,
  type: activeTab.value,
})

// 点击筛选标签
const handleTabClick = (tab) => {
  resetSearchInfo()
  activeTab.value = tab.paneName
  searchInfo.value.type = tab.paneName
  getTableData()
}

// 搜索按钮
const onSearch = () => {
  getTableData()
}

// 重置按钮
const onReset = () => {
  resetSearchInfo()
  getTableData()
}

// 重置搜索条件
const resetSearchInfo = () => {
  searchInfo.value = {
    page: 1,
    pageSize: 10,
    total: 0,
    sortBy: null,
    isDesc: null,
    keyword: null,
    type: activeTab.value,
  }
}

// 搜索获取表格数据
const getTableData = async () => {
  const response = await getProductList(searchInfo.value)
  if (response.code === 0) {
    tableData.value = response.data.list
    searchInfo.value.pageSize = response.data.pageSize
    searchInfo.value.total = response.data.total
  }
  // console.log('searchInfo: ', searchInfo.value)
  // console.log('tableData: ', tableData.value)
  // console.log('total: ' + response.data.total)
  // console.log('page: ' + response.data.page)
  // console.log('pageSize: ' + response.data.pageSize)
}

/********
 * 城市选择器
 ********/
const options= provinceAndCityData ;
const selectedOptions=ref([])
const handleChange= (value)=> {
  btnStatus.value['regionId'] = 0
  const cityName=getCodeToText(null, value)
  console.log(cityName)
  console.log(cities)
  const data = cities.value.filter(c=>c.label===cityName)
  console.log(data)
  drawerForm.value.regionId=data[0].value;
  console.log(drawerForm.value.regionId)
}
const getCodeToText= (codeStr, codeArray)=> {
  if (null === codeStr && null === codeArray) {
    return null;
  } else if (null === codeArray) {
    codeArray = codeStr.split(",");
  }
  let area = "";
  if(CodeToText[codeArray[0]]=='北京市'||CodeToText[codeArray[0]]=='重庆市'||CodeToText[codeArray[0]]=='上海市'||CodeToText[codeArray[0]]=='天津市'||CodeToText[codeArray[0]]=='香港特别行政区'||CodeToText[codeArray[0]]=='澳门特别行政区'){
    area=CodeToText[codeArray[0]];
  }else{
    area=CodeToText[codeArray[1]];
  }
  console.log(area)
  return area;
}


/********
 * 表格栏
 ********/
// 表格数据
const tableData = ref([])
// 多选
const selectionList = ref([])

// 多选改变
const handleSelectionChange = (val) => {
  selectionList.value = val
}

// 分页
const handleCurrentChange = (val) => {
  searchInfo.value.page = val
  getTableData()
}
const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val
  getTableData()
}

// 排序
const sortChange = ({ prop, order }) => {
  if (prop) {
    searchInfo.value.orderBy = prop
    searchInfo.value.isDesc = order === 'descending'
  }
  getTableData()
}

//同步
const SyncSupplyDemand=async(row)=>{
  const resp = await SyncSupplyDemandInfo({ id: row.product.id })
  if (resp.code === 0) {
    ElMessage({
      type: 'success',
      message: '同步成功'
    })
    await getTableData()
  } else {
    ElMessage({
      showClose: true,
      message: '同步失败',
      type: 'error',
    })
  }
}

// 打开删除确认弹窗
const handleDelete = (row) => {
  currentProductId.value = row.product.id;
  deleteVisible.value = true;
};

// 关闭弹窗
const handleDialogClose = () => {
  deleteVisible.value = false;
  currentProductId.value = null;
};

// 打开删除弹窗（统一入口）
const openDeleteDialog = (type, row = null) => {
  deleteType.value = type;

  if (type === 'single' && row) {
    currentProductId.value = row.product.id;
  } else if (type === 'batch') {
    selectedIds.value = selectionList.value.map(item => item.product.id);
  }

  deleteVisible.value = true;
};

// 确认删除
const confirmDelete = async () => {
  try {
    let resp;

    if (deleteType.value === 'single') {
      // 调用单删接口
      resp = await deleteProduct({ id: currentProductId.value });
      ElMessage({
        type: 'success',
        message: '删除成功'
      });
    } else if (deleteType.value === 'batch') {
      // 调用批量删除接口
      resp = await deleteProductByIds({ ids: selectedIds.value });
      ElMessage({
        type: 'success',
        message: `批量删除 ${selectedIds.value.length} 条记录成功`
      });
    }

    if (resp.code === 0) {
      await getTableData(); // 刷新表格数据
      emits('confirm'); // 通知父组件删除成功
    } else {
      ElMessage({
        showClose: true,
        message: '删除失败',
        type: 'error',
      });
    }
  } catch (error) {
    ElMessage({
      showClose: true,
      message: '删除出错，请重试',
      type: 'error',
    });
    console.error('删除失败:', error);
  } finally {
    resetDeleteState();
    getTableData();
  }
};

// 重置删除状态
const resetDeleteState = () => {
  deleteVisible.value = false;
  deleteType.value = '';
  currentProductId.value = null;
  selectedIds.value = [];
};
// 保持原有方法不变
const handleItemChange = (item) => {
  btnStatus.value[item] = 0;
};

// 租赁时间配置
const rentalEndOptions = computed(() => ({
  disabledDate: (time) => {
    if (!drawerForm.value.rentalStartTime) return false;
    return time.getTime() < drawerForm.value.rentalStartTime.getTime();
  }
}));

// 信息有效期配置
const infoEndOptions = computed(() => ({
  disabledDate: (time) => {
    if (!drawerForm.value.startTime) return false;
    return time.getTime() < drawerForm.value.startTime.getTime();
  }
}));

// 租赁开始时间变更处理
const handleRentalStartTimeChange = () => {
  if (
      drawerForm.value.rentalStartTime &&
      drawerForm.value.rentalEndTime &&
      drawerForm.value.rentalEndTime < drawerForm.value.rentalStartTime
  ) {
    drawerForm.value.rentalEndTime = new Date(drawerForm.value.rentalStartTime);
    ElMessage.warning('租赁结束时间已自动调整为开始时间');
  }
  handleItemChange('rentalTime');
};

// 租赁结束时间变更处理
const handleRentalEndTimeChange = () => {
  if (
      drawerForm.value.rentalStartTime &&
      drawerForm.value.rentalEndTime &&
      drawerForm.value.rentalEndTime < drawerForm.value.rentalStartTime
  ) {
    drawerForm.value.rentalEndTime = new Date(drawerForm.value.rentalStartTime);
    ElMessage.warning('租赁结束时间不能早于开始时间，已自动调整');
  }
  handleItemChange('rentalTime');
};

// 信息开始时间变更处理
const handleInfoStartTimeChange = () => {
  if (
      drawerForm.value.startTime &&
      drawerForm.value.endTime &&
      drawerForm.value.endTime < drawerForm.value.startTime
  ) {
    drawerForm.value.endTime = new Date(drawerForm.value.startTime);
    ElMessage.warning('信息截止时间已自动调整为发布时间');
  }
  handleItemChange('time');
};

// 信息结束时间变更处理
const handleInfoEndTimeChange = () => {
  if (
      drawerForm.value.startTime &&
      drawerForm.value.endTime &&
      drawerForm.value.endTime < drawerForm.value.startTime
  ) {
    drawerForm.value.endTime = new Date(drawerForm.value.startTime);
    ElMessage.warning('信息截止时间不能早于发布时间，已自动调整');
  }
  handleItemChange('time');
};
// 批量删除
const onDelete = async () => {
  ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = selectionList.value.map(item => item.product.id)
    const resp = await deleteProductByIds({ids: ids})
    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      await getTableData()
    } else {
      ElMessage({
        showClose: true,
        message: '删除失败',
        type: 'error',
      })
    }
  })
}

/******
 * 抽屉
 ******/
const drawerTitle = ref('')
const drawerFormVisible = ref(false)
const drawerFormRef = ref(null)
const drawerForm = ref({})

const cities = ref([])

const getCities = async () => {
  const response = await getSupplyOrDemandCityList()
  if (response.code === 0) {
    response.data.list.forEach((item) => {
      if (item.cities.length > 0) {
        item.cities.forEach((city) => {
          cities.value.push({
            label: city.name,
            value: city.id,
          })
        })
      }
      else {
        cities.value.push({
          label: item.province.name,
          value: item.province.id,
        })
      }
    })
  }
}

const validateForm = (item, msg) => {
  return (rule, value, callback) => {
    if (item === 'rentalTime') {
      if (!btnStatus.value[item] && !(drawerForm.value.rentalStartTime && drawerForm.value.rentalEndTime)) {
        callback(new Error(msg))
      }
    } else if (item === 'time') {
      if (!btnStatus.value[item] && !(drawerForm.value.startTime && drawerForm.value.endTime)) {
        callback(new Error(msg))
      }
    } else  if (item === 'status') {
      if (drawerTitle.value === '审核' && !drawerForm.value.status) {
        callback(new Error(msg))
      }
    } else {
      if (!btnStatus.value[item] && !drawerForm.value[item]) {
        callback(new Error(msg));
      }
    }
    callback()
  }
}

const rules = ref({
  subType: [{ required: true, message: "请选择类型", trigger: "change" }],
  source: [{ required: true, message: "请选择来源", trigger: "change" }],
  regionId: [{ required: true, validator: validateForm('regionId', '请选择区域地址'), trigger: "change" }],
  cityId: [{ required: true, validator: validateForm('cityId', '请选择区域地址'), trigger: "change" }],
  price: [{ required: true, validator: validateForm('price', '请输入价格'), trigger: "blur" }],
  number: [{ required: true, validator: validateForm('number', '请输入数量'), trigger: "blur" }],
  rentalTime: [{ required: true, validator: validateForm('rentalTime', '请选择租赁时间'), trigger: "blur" }],
  time: [{ required: true, validator: validateForm('time', '请选择信息有效期'), trigger: "blur" }],
  name: [{ required: true, message: "请输入标题", trigger: "blur" }],
  description: [{ required: true, message: "请输入描述", trigger: "blur" }],
  contact: [{ required: true, message: "请输入联系人", trigger: "blur" }],
  phone: [
    { required: true, message: '请输入电话号码', trigger: 'blur' },
    {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入有效的手机号码',
      trigger: 'blur'
    }
  ],
  email: [
    {
      pattern: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
      message: '请输入有效的邮箱地址',
      trigger: 'blur'
    }
  ],
  status: [{ required: true, validator: validateForm('status', '请选择审核按钮'), trigger: "blur" }],
})

const currentKey = ref('')
const currentValue = ref('')
const keyValuePairs = ref([])
const keys = ref([])
const unit = computed(() => {
  let u = keys.value.find(item => item.name === currentKey.value)?.unit
  return u ? u : ''
})

const addKeyValue = () => {
  if (currentKey.value && currentValue.value) {
    keyValuePairs.value.push({key: currentKey.value, value: currentValue.value + ' ' + unit.value})
    currentKey.value = ''
    currentValue.value = ''
    unit.value = ''
  }
}

const removeKeyValue = (index) => {
  keyValuePairs.value.splice(index, 1)
}


const btnStatus = ref({})
const handleTimeTypeChange = (value) => {
  if (value === 0) {
    // 选择"默认"时触发原函数
    handleBtnClick('time');

  } else {
    // 选择"自定义"时的处理（可选）
    // 例如清空默认状态
  }
};

const handleTimeChange = (value) => {
  if (value === 0) {
    // 选择"默认"时触发原函数
    handleBtnClick('rentalTime');

  } else {
    // 选择"自定义"时的处理（可选）
    // 例如清空默认状态
  }
};
const handleNumChange = (value) => {
  if (value === 0) {
    // 选择"默认"时触发原函数
    handleBtnClick('number');

  } else {
    // 选择"自定义"时的处理（可选）
    // 例如清空默认状态
  }
};
const handlePriChange = (value) => {
  if (value === 0) {
    // 选择"默认"时触发原函数
    handleBtnClick('price');

  } else {
    // 选择"自定义"时的处理（可选）
    // 例如清空默认状态
  }
};
const handleIdChange = (value) => {
  if (value === 0) {
    // 选择"默认"时触发原函数
    handleBtnClick('regionId');

  } else {
    // 选择"自定义"时的处理（可选）
    // 例如清空默认状态
  }
};
const handleBtnClick = (btn) => {
  if (btn === 'rentalTime') {
    drawerForm.value.rentalStartTime = null;
    drawerForm.value.rentalEndTime = null;
  } else if (btn === 'time') {
    drawerForm.value.startTime = null;
    drawerForm.value.endTime = null;
  } else if (btn === 'regionId') {
    drawerForm.value.regionId = 0;
    drawerForm.value.cityId = [];
  } else {
    drawerForm.value[btn] = null;
  }
  btnStatus.value[btn] = 1
  drawerFormRef.value.validateField(btn)
}


/*
 * 子类别字典
 *
 * {
 *   1: [{id: 2, name: 'xxx'}, ...],
 *   ...
 * }
 */
const subtypeMap = ref({})
// 子类别描述
const subtypeDesc = ref({
  7: { name: '算力', value: 7, subLabel1: '计算类型', subLabel2: '卡类型' },
  8: { name: '硬件', value: 8, subLabel1: '品牌', subLabel2: '型号' },
  9: { name: '机柜', value: 9, subLabel1: '供电方式', subLabel2: null },
  10: { name: '带宽', value: 10, subLabel1: '运营商', subLabel2: '路由类型' }
})

const getSubtype = async () => {
  const response = await getProductConfigList({ type: 2 })
  if (response.code === 0) {
    response.data.list?.forEach(item => {
      if (item.root) {
        for (let key in subtypeDesc.value) {
          if (item.root.name === subtypeDesc.value[key].name) {
            subtypeDesc.value[key].id = item.root.id
          }
        }
      }
      if (item.root && item.second) {
        subtypeMap.value[item.root.id] = item.second.map(item_ => ({id: item_.id, name: item_.name}))
      }
      if (item.second && item.third) {
        item.second.forEach(secondItem => {
          subtypeMap.value[secondItem.id] = item.third.filter(item_ => item_.parent_id === secondItem.id).map(item_ => ({id: item_.id, name: item_.name}))
        })
      }
    })
    // console.log(response.data.list)
    // console.log(JSON.stringify(subtypeDesc.value))
    // console.log(JSON.stringify(subtypeMap.value))
  }
}

const getKeys = () => {
    getCategoryList({ type: 17 }).then((response) => {
    if (response.code === 0) {
      keys.value = response.data.list
    }
  })
}

const addFunc = function () {
  resetDrawerForm()
  openDrawer('add')
}

const editFunc = function (row) {
  drawerForm.value = JSON.parse(row.product.specJson)
  drawerForm.value.id = row.product.id
  keyValuePairs.value = drawerForm.value.keyValuePairs
  btnStatus.value = drawerForm.value.btnStatus
  openDrawer('edit')
}

const auditFunc = function (row) {
  drawerForm.value = JSON.parse(row.product.specJson)
  drawerForm.value.id = row.product.id
  keyValuePairs.value = drawerForm.value.keyValuePairs
  btnStatus.value = drawerForm.value.btnStatus
  openDrawer('audit')
}

const openDrawer = (key) => {
  switch (key) {
    case 'add':
      drawerTitle.value = '新增'
      drawerFormVisible.value = true
      break
    case 'edit':
      drawerTitle.value = '编辑'
      drawerFormVisible.value = true
      break
    case 'audit':
      drawerTitle.value = '审核'
      drawerFormVisible.value = true
      break
    default:
      break
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  resetDrawerForm()
}

const resetDrawerForm = () => {
  drawerForm.value = {}
  keyValuePairs.value = []
  btnStatus.value = {}
  if (drawerFormRef.value) {
    drawerFormRef.value.clearValidate()
  }
}

const handleUpOrDown = async (row) => {
  row.product.specJson = JSON.parse(row.product.specJson)
  row.product.specJson.status = row.product.status
  row.product.specJson = JSON.stringify(row.product.specJson)
  let data = {
    id: row.product.id,
    status: row.product.status,
    specJson: row.product.specJson
  }
  let auditResp = await updateProduct(data)
  if (auditResp.code === 0) {
    ElMessage({
      type: 'success',
      message: row.product.status === 2 ? '上架成功' : '下架成功',
      showClose: true
    })
    await getTableData()
  }
  //通知节点产品的审核状态
  const param = {
    id: row.product.id,
    status: row.product.status,
  };
  // await notifyCloudProductStatus(param);
}

// 新增/编辑确认按钮
const confirmForm = () => {
  drawerFormRef.value.validate(async valid => {
    if (valid) {
      // 封装表单信息
      drawerForm.value.keyValuePairs = keyValuePairs.value
      drawerForm.value.btnStatus = btnStatus.value
      drawerForm.value.specJson = JSON.stringify(keyValuePairs.value)
      // 封装参数信息
      let data = {
        id: drawerForm.value.id ? drawerForm.value.id : null,
        type: activeTab.value,
        subType: drawerForm.value.subType,
        typeId: drawerForm.value.subType2Id ? drawerForm.value.subType2Id : drawerForm.value.subType1Id,
        source: drawerForm.value.source,
        regionId: drawerForm.value.regionId,
        startTime: drawerForm.value.startTime,
        endTime: drawerForm.value.endTime,
        name: drawerForm.value.name,
        description: drawerForm.value.description,
        user: drawerForm.value.contact,
        userPhone: drawerForm.value.phone,
        userEmail: drawerForm.value.email,
        specJson: drawerForm.value.specJson,
        status: drawerForm.value.status
      };
      const par = {
        id: drawerForm.value.id,
        status: drawerForm.value.status
      }
      switch (drawerTitle.value) {
        case '新增':
          data.status = 0
          drawerForm.value.status = 0
          data.specJson = JSON.stringify(drawerForm.value)
          let addResp = await createProduct(data)
          if (addResp.code === 0) {
            ElMessage({
              type: 'success',
              message: '提交新增成功',
              showClose: true
            })
            await getTableData()
            closeDrawer()
          }
          break
        case '编辑':
          data.status = 0
          drawerForm.value.status = 0
          data.specJson = JSON.stringify(drawerForm.value)
          let editResp = await updateProduct(data)
          if (editResp.code === 0) {
            ElMessage({
              type: 'success',
              message: '提交更新成功',
              showClose: true
            })
            await getTableData()
            closeDrawer()
          }
          break
        case '审核':
          //通知审核结果
          // notifyCloudProductStatus(par)
          let auditResp = await updateProduct(data)
          if (auditResp.code === 0) {
            ElMessage({
              type: 'success',
              message: '提交审核成功',
              showClose: true
            })
            await getTableData()
            closeDrawer()
          }
          break
      }
    } else {
      ElMessage({
        type: 'warning',
        message: '请填写完整表单',
        showClose: true
      })
    }
  })
}

/********
 * 对话框(暂时不用)
 * ******
 */
const dialogTitle = ref('')
const dialogVisible = ref(false)
const deleteVisible =ref(false)
const deleteType = ref(''); // 删除类型：single(单删) / batch(批量删)
const currentProductId = ref(null); // 单删时的产品ID
const selectedIds = ref([]); // 批量删除时的ID数组
// const auditFunc = function (row) {
//   openDialog('audit')
// }

const openDialog = (key) => {
  switch (key) {
    case 'audit':
      dialogTitle.value = '审核'
      dialogVisible.value = true
      break
    case 'detail':
      dialogTitle.value = '详情'
      dialogVisible.value = true
      break
    default:
      break
  }
}

// 审核确认按钮
const confirmAudit = async () => {}

/********
 * 初始化
 ********/
onMounted(() => {
  getSubtype()
  getCities()
  getKeys()
  getTableData()
})
</script>


<style scoped lang="scss">
//.warning {
//  color: #dc143c;
//}

// 设置筛选标签的下边距
.tabs {
  margin-bottom: 30px;
}

// 设置按钮间隔且保证挤压时上下对齐
.button-gap {
  margin-left: 0;

}

// 文本溢出表格后截断
.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

// 文本详情
.text-detail {
  color: #77b2f3;
  cursor: pointer;
}

// 表单边框
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

// 按钮被激活
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
.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 自定义提示框样式 */
.custom-tooltip-content {
  max-width: 200px;
  max-height: 102px;
  overflow: hidden;
  white-space: normal;
  word-break: break-all;
}
</style>
