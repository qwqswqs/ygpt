<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <!-- 条件筛选栏 -->
    <div class="gva-search-box">
      <!-- 表格选择 -->
      <el-tabs v-model="activeTab"  @tab-click="(tab) => handleTabClick(tab)">
        <el-tab-pane v-for="item in tabList" :key="item.name" :label="item.label" :name="item.name"/>
      </el-tabs>
      <!-- 搜索条件 -->
      <div class="gva-btn-list" style="margin-top: 16px">

        <el-form ref="searchForm" :model="searchInfo" :inline="true" style="display: flex;justify-content: space-between">
          <div>
            <el-form-item label="标题">
              <el-input v-model="searchInfo.keyword" style="width: 200px" placeholder="请输入关键词" clearable/>
            </el-form-item>
            <el-form-item label="资源类型">
              <el-select v-model="searchInfo.subType" style="width: 200px" clearable placeholder="请选择">
                <el-option v-for="item in searchOptions.subType" :key="item.value" :label="item.label" :value="item.value"/>
              </el-select>
            </el-form-item>
            <!--        <el-form-item label="是否过期">-->
            <!--          <el-select v-model="searchInfo.isExpired" style="width: 200px" clearable placeholder="请选择">-->
            <!--            <el-option v-for="item in searchOptions.isExpired" :key="item.value" :label="item.label" :value="item.value" />-->
            <!--          </el-select>-->
            <!--        </el-form-item>-->
<!--                    <el-form-item label="状态">-->
<!--                      <el-select v-model="searchInfo.status" style="width: 200px" clearable placeholder="请选择">-->
<!--                        <el-option v-for="item in searchOptions.status" :key="item.value" :label="item.label" :value="item.value" />-->
<!--                      </el-select>-->
<!--                    </el-form-item>-->
            <!--        <el-form-item label="来源">-->
            <!--          <el-select v-model="searchInfo.source" style="width: 200px" clearable placeholder="请选择">-->
            <!--            <el-option v-for="item in searchOptions.source" :key="item.value" :label="item.label" :value="item.value" />-->
            <!--          </el-select>-->
            <!--        </el-form-item>-->

          </div>
          <el-form-item>
            <el-button class="button-gap" type="primary" @click="onSearch">
              查询
            </el-button>
            <el-button class="button-gap" type="info"  @click="onReset">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 表格栏 -->
    <div class="gva-table-box">
      <!-- 表格 -->
      <div v-if="tableData && tableData.length > 0">
        <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
          <el-table-column type="selection" width="55"/>
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
                <div class="text-overflow">
                  {{
                    searchOptions.status.find(item => item.value === row.product.status)?.label || '-'
                  }}
                </div>

              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="操作" fixed="right" width="300">
            <template #default="{ row }">
              <!--              <el-button  type="primary" link @click="editFunc(row)">-->
              <!--                编辑-->
              <!--              </el-button>-->
              <el-button :disabled="row.product.status !== 0" class="text-blue-500"  type="text" link
                         @click="auditFunc(row)">
                审核
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
    <el-drawer :show-close="false" v-model="drawerFormVisible" style="min-width: 500px;" size="70%" :before-close="closeDrawer">
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
      <el-form ref="drawerFormRef" :model="drawerForm" :rules="rules">
        <el-collapse v-model="activeNames">
          <el-collapse-item  name="1" >
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
                      @click="drawerForm.subType = null; /* 其他逻辑 */"
                      :disabled="true"
                  >
                  {{ subtype.label }}
                  </el-button>
                </div>
              </el-form-item>

              <el-form-item
                  v-if="drawerForm.subTypeId && subtypeMap[drawerForm.subTypeId] && subtypeMap[drawerForm.subTypeId].length > 0"
                  label-width="100px"
                  :label="subtypeDesc[drawerForm.subType].subLabel1"
                  prop="subType1Id"
              >
                <el-select
                    v-model="drawerForm.subType1Id"
                    style="width: 320px;"
                    placeholder="请选择"
                    clearable
                    :disabled="true"
                >
                <el-option
                    v-for="item in subtypeMap[drawerForm.subTypeId]"
                    :key="item.id"
                    :label="item.name"
                    :value="item.id"
                    :disabled="true"
                />
                </el-select>
              </el-form-item>

              <el-form-item
                  v-if="drawerForm.subType1Id && subtypeMap[drawerForm.subType1Id] && subtypeMap[drawerForm.subType1Id].length > 0"
                  label-width="100px"
                  :label="subtypeDesc[drawerForm.subType].subLabel2"
                  prop="subType2Id"
              >
                <el-select
                    v-model="drawerForm.subType2Id"
                    style="width: 320px;"
                    placeholder="请选择"
                    clearable
                    :disabled="true"
                >
                <el-option
                    v-for="item in subtypeMap[drawerForm.subType1Id]"
                    :key="item.id"
                    :label="item.name"
                    :value="item.id"
                    :disabled="true"
                />
                </el-select>
              </el-form-item>

              <el-form-item label="区域地址" label-width="100px" prop="cityId">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <el-radio-group v-model="drawerForm.regionType" :disabled="true"> <!-- 禁用单选组 -->
                    <el-radio :label="0">不限</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <el-cascader
                      style="width: 320px;"
                      :options="options"
                      v-model="drawerForm.cityId"
                      :disabled="true"
                  />
                </div>
              </el-form-item>

              <!-- 计费计价字段 -->
              <el-form-item label="计费计价" label-width="100px" prop="price">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <el-radio-group v-model="drawerForm.priceType" :disabled="true"> <!-- 禁用单选组 -->
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <el-input
                      v-model="drawerForm.price"
                      style="width: 320px;"
                      placeholder="示例：10.00元/月/台"
                      clearable
                      :disabled="true"
                  />
                </div>
              </el-form-item>

              <!-- 供需数量字段 -->
              <el-form-item label="供需数量" label-width="100px" prop="number">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <el-radio-group v-model="drawerForm.numberType" :disabled="true"> <!-- 禁用单选组 -->
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <el-input
                      v-model="drawerForm.number"
                      style="width: 320px;"
                      placeholder="示例：9台"
                      clearable
                      :disabled="true"
                  />
                </div>
              </el-form-item>

              <!-- 租赁时间字段 -->
              <el-form-item label="租赁时间" label-width="100px" prop="rentalTime">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <el-radio-group v-model="drawerForm.rentalTimeType" :disabled="true"> <!-- 禁用单选组 -->
                    <el-radio :label="0">面议</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <div v-if="drawerForm.rentalTimeType === 1" style="display: flex; gap: 15px;">
                    <el-date-picker
                        v-model="drawerForm.rentalStartTime"
                        placeholder="请选择租赁开始时间"
                        style="max-width: 200px;"
                        type="datetime"
                        :disabled="true"
                    />
                    <el-date-picker
                        v-model="drawerForm.rentalEndTime"
                        placeholder="请选择租赁结束时间"
                        style="max-width: 200px;"
                        type="datetime"
                        :disabled="true"
                    />
                  </div>
                </div>
              </el-form-item>

              <!-- 信息有效期字段 -->
              <el-form-item label="信息有效期" label-width="100px" prop="time">
                <div style="display: flex; gap: 15px; align-items: center;">
                  <el-radio-group v-model="drawerForm.timeType" :disabled="true"> <!-- 禁用单选组 -->
                    <el-radio :label="0">默认</el-radio>
                    <el-radio :label="1">自定义</el-radio>
                  </el-radio-group>

                  <div v-if="drawerForm.timeType === 1" style="display: flex; gap: 15px;">
                    <el-date-picker
                        v-model="drawerForm.startTime"
                        placeholder="请选择信息发布"
                        style="max-width: 200px;"
                        type="datetime"
                        :disabled="true"
                    />
                    <el-date-picker
                        v-model="drawerForm.endTime"
                        placeholder="请选择信息截止时间"
                        style="max-width: 200px;"
                        type="datetime"
                        :disabled="true"
                    />
                  </div>
                </div>
              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item title="" name="2" >
            <template #title>
              <div :style="{ marginLeft: '10px' }">详情描述</div>
            </template>
            <div class="form-section">
              <el-form-item label="标题" label-width="100px" prop="name">
                <el-input
                    disabled
                    v-model="drawerForm.name"
                    style="width: 320px;"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="描述" label-width="100px" prop="description">
                <el-input
                    disabled
                    v-model="drawerForm.description"
                    style="width: 1016px;height: 56px;"
                    type="textarea"
                    :rows="2"
                    maxlength="3000"
                    show-word-limit
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="规格" label-width="100px">
                <el-table :data="keyValuePairs" style="width: 90%" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
                  <el-table-column label="规格名称" prop="key"/>
                  <el-table-column label="规格值" prop="value"/>
                </el-table>
              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item title="" name="3" >
            <template #title>
              <div :style="{ marginLeft: '10px' }">联系方式</div>
            </template>
            <div class="form-section">
              <el-form-item label="联系人" label-width="100px" prop="contact">
                <el-input
                    disabled
                    v-model="drawerForm.contact"
                    style="width: 320px"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="电话" label-width="100px" prop="phone">
                <el-input
                    disabled
                    v-model="drawerForm.phone"
                    style="width: 320px"
                    placeholder="请输入"
                    clearable
                />
              </el-form-item>
              <el-form-item label="邮箱" label-width="100px" prop="email">
                <el-input
                    disabled
                    v-model="drawerForm.email"
                    style="width: 320px"
                    placeholder="请输入（可选）"
                    clearable
                />
              </el-form-item>
            </div>
          </el-collapse-item>
          <el-collapse-item title="" name="4" >
            <template #title>
              <div :style="{ marginLeft: '10px' }">审核信息</div>
            </template>
            <div v-if="drawerTitle === '审核'" class="form-section">
              <el-form-item label="审核状态" label-width="100px" prop="status">
                <el-radio-group v-model="drawerForm.status">
                  <el-radio :value="1" >
                    通过
                  </el-radio>
                  <el-radio :value="-1" >
                    驳回
                  </el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="审核意见" label-width="90px" prop="feedback" style="padding-left: 10px;">
                <el-input type="textarea" v-model="drawerForm.feedback"/>
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
    <el-dialog v-model="dialogVisible" style="min-width: 500px" title="Warning" width="70%" center>
      <span>
        It should be noted that the content will not be aligned in center by
        default
      </span>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="dialogVisible = false">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'

import {
  createProduct,
  deleteProduct,
  deleteProductByIds,
  getProductConfigList,
  getProductList,
  getSupplyOrDemandCityList,
  updateProduct
} from "@/api/yunguan/product/productSupply"
import {getCategoryList} from "@/api/yunguan/system/type"
import { provinceAndCityData ,  CodeToText, TextToCode } from 'element-china-area-data'

/*********
 * 条件筛选
 *********/
const isActive = ref(true);

// 通过 computed 生成展开项数组
const activeNames = computed(() => {
  return isActive.value ? ['1', '2', '3','4'] : [];
});
// 搜索表单
const searchForm = ref(null)
// 已激活标签
const activeTab = ref(3)
// 标签列表
const tabList = ref([
  // { label: '全部', name: 0 },
  {label: '供给', name: 3},
  {label: '需求', name: 4},
])

// 条件选项
const searchOptions = ref({
  subType: [
    {label: '算力', value: 7},
    {label: '硬件', value: 8},
    {label: '机柜', value: 9},
    {label: '带宽', value: 10},
  ],
  isExpired: [
    {label: '未过期', value: 0},
    {label: '已过期', value: 1},
  ],
  status: [
    {label: '未通过', value: -1},
    {label: '未审核', value: 0},
    {label: '已下架', value: 1},
    {label: '已上架', value: 2}
  ],
  source: [
    {label: '平台', value: 1},
    {label: '第三方', value: 2},
    {label: '运营商', value: 3},
    {label: '智算中心', value: 4}
  ]
})

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

// 搜索条件
const searchInfo = ref({
  page: 1,
  pageSize: 10,
  total: 0,
  sortBy: null,
  isDesc: null,
  keyword: null,
  type: activeTab.value,
  status:0,
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
    status:0,
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
  const data = cities.value.find(c=>c.label===cityName)
  console.log(data)
  drawerForm.value.regionId=data.value;
}
const getCodeToText= (codeStr, codeArray)=> {
  if (null === codeStr && null === codeArray) {
    return null;
  } else if (null === codeArray) {
    codeArray = codeStr.split(",");
  }
  let area = "";
  if(CodeToText[codeArray[0]]=='北京市'||CodeToText[codeArray[0]]=='重庆市'||CodeToText[codeArray[0]]=='上海市'||CodeToText[codeArray[0]]=='天津市'){
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
const sortChange = ({prop, order}) => {
  if (prop) {
    searchInfo.value.orderBy = prop
    searchInfo.value.isDesc = order === 'descending'
  }
  getTableData()
}

// 单个删除
const deleteFunc = async (row) => {
  ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const resp = await deleteProduct({id: row.product.id})
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
const feedbackForm = ref({})

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
    } else if (item === 'status') {
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
  status: [{required: true, validator: validateForm('status', '请选择审核按钮'), trigger: "blur"}],
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

const handleBtnClick = (btn) => {
  if (btn === 'rentalTime') {
    drawerForm.value.rentalStartTime = null;
    drawerForm.value.rentalEndTime = null;
  } else if (btn === 'time') {
    drawerForm.value.startTime = null;
    drawerForm.value.endTime = null;
  } else {
    drawerForm.value[btn] = null;
  }
  btnStatus.value[btn] = 1
  drawerFormRef.value.validateField(btn)
}

const handleItemChange = (item) => {
  btnStatus.value[item] = 0
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
  7: {name: '算力', value: 7, subLabel1: '计算类型', subLabel2: '卡类型'},
  8: {name: '硬件', value: 8, subLabel1: '品牌', subLabel2: '型号'},
  9: {name: '机柜', value: 9, subLabel1: '供电方式', subLabel2: null},
  10: {name: '带宽', value: 10, subLabel1: '运营商', subLabel2: '路由类型'}
})

const getSubtype = async () => {
  const response = await getProductConfigList({type: 2})
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
          subtypeMap.value[secondItem.id] = item.third.filter(item_ => item_.parent_id === secondItem.id).map(item_ => ({
            id: item_.id,
            name: item_.name
          }))
        })
      }
    })
    // console.log(response.data.list)
    // console.log(JSON.stringify(subtypeDesc.value))
    // console.log(JSON.stringify(subtypeMap.value))
  }
}

const getKeys = () => {
  getCategoryList({type: 17}).then((response) => {
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
        status: drawerForm.value.status,
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
          data.feedback=drawerForm.value.feedback;
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
const confirmAudit = async () => {
}

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
