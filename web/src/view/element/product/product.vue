<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">

      <div class="gva-search-box" style="margin-top: 16px">
        <div class="gva-btn-list">
          <el-form ref="searchForm" :model="searchInfo" :inline="true" label-width="70px"
                   style="position: relative">
            <el-form-item label="产品名称">
              <el-input style="width: 240px;" v-model="searchInfo.name" placeholder="请输入"/>
            </el-form-item>
            <el-form-item label-width="40px" label="类别">
              <el-select style="width: 240px;" v-model="searchInfo.type" placeholder="请选择" clearable>
                <el-option label="裸金属" :value="1"/>
                <el-option label="云主机" :value="2"/>
                <el-option label="容器" :value="3"/>
              </el-select>
            </el-form-item>
            <el-form-item label="产品类型">
              <el-select style="width: 240px;" v-model="searchInfo.isCustom" placeholder="请选择" clearable>
                <el-option label="标准产品" :value="0"/>
                <el-option label="定制产品" :value="1"/>
              </el-select>
            </el-form-item>
            <el-form-item label-width="40px" label="状态">
              <el-select style="width: 240px;" v-model="searchInfo.auditStatus" placeholder="请选择" clearable>
                <el-option label="通过" :value="1"/>
                <el-option label="待审核" :value="0"/>
                <el-option label="驳回" :value="-1"/>
              </el-select>
            </el-form-item>
            <el-form-item label="同步状态">
              <el-select style="width: 240px;" v-model="searchInfo.isSynced" placeholder="请选择" clearable>
                <el-option label="未同步" :value="1"/>
                <el-option label="已同步" :value="2"/>
              </el-select>
            </el-form-item>
            <el-form-item style="position: absolute;right: 0">
              <el-button type="primary" @click="handleSearch">查询</el-button>
              <el-button type="info" @click="onReset">重置</el-button>

            </el-form-item>

          </el-form>


        </div>
        <!-- 条件选择 -->

      </div>
      <div class="gva-table-box">
        <div style="margin-bottom: 16px;">
          <el-button
              type="primary"
              @click="OpenDialog"
          >
            添加
          </el-button>
        </div>

        <el-table :data="dataList" row-key="ID"
                @sort-change="handleSortChange"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
        <el-table-column label="名称">
          <template #default="scope">
            <el-button link type="text" @click="showDetail(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ productType[scope.row.type - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">
              {{ scope.row.description }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="产品类型">
          <template #default="scope">
            <span>{{ scope.row.isCustom ? '定制产品' : '标准产品' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <div style="display: flex;align-items: center">
              <div :class="handleTagType(scope.row.auditStatus)"></div>
              <span v-if="scope.row.auditStatus==-1">{{ '驳回' }}</span>
              <span v-if="scope.row.auditStatus==0">{{ '待审核' }}</span>
              <span v-if="scope.row.auditStatus==1">{{ '通过' }}</span>
              <!--              <el-button type="text" @click="auditVisible=true;auditOpinion=scope.row.auditOpinion" v-else>{{ scope.row.auditStatus==-1?'驳回':'通过' }}</el-button>-->
            </div>
          </template>
        </el-table-column>
        <el-table-column label="同步状态">
          <template #default="scope">
            <div style="display: flex;align-items: center">
              <div :class="handleTagIsSynced(scope.row.isSynced)"></div>
              <span>{{ scope.row.isSynced == 1 ? '未同步' : '已同步' }}</span>
            </div>

          </template>
        </el-table-column>
        <el-table-column label="创建时间" sortable="custom" >
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button link type="text" class="text-blue-500" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button
                link
                type="text"
                class="text-blue-500 sync-btn"
            :disabled="scope.row.isSynced === 2"
            @click="handleSyncedProduct(scope.row)"
            >
            同步
            </el-button>
            <!--            <el-button type="text" @click="handleStatusChange(scope.row)">-->
            <!--              {{ scope.row.status != 2 ? '上架' : '下架' }}-->
            <!--            </el-button>-->
            <el-button link type="text" class="text-blue-500" @click="openDeleteDialog(scope.row.ID)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      </div>


      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
      <!--      新增和编辑对话框-->
      <el-dialog v-model="dialogVisible" class="no-select" :show-close="false"
                 width="1000px">
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{IsCreate? '创建算力产品':'编辑算力产品' }}</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="dialogVisible=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <el-form ref="dataInfo" :model="dataInfo" :rules="dataRules" label-width="120px">
          <el-tabs v-model="activeConfigName" class="tabs no-select" @tabChange="handleConfigTabsClick">
            <!-- tab-pane 标准产品 -->
            <el-tab-pane :disabled="!IsCreate" label="标准产品" :name="'standard'">
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="产品名称" prop="name">
                    <el-input style="width: 300px;" v-model="dataInfo.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="产品类别" prop="type">
                    <el-select  :disabled="!IsCreate" v-model="dataInfo.type" placeholder="请选择" style="width: 300px"
                               @change="handleChangeType">
                      <el-option v-for="(item,index) in productType"
                                 :key="index"
                                 :label="item"
                                 :value="index+1"/>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item style="padding-left: 10px" label-width="110px" :label="'支持的计费方式'">
                <el-checkbox-group :min="1" v-model="dataInfo.billing">
                  <el-checkbox label="计时" value="hour"/>
                  <el-checkbox label="包日" value="day"/>
                  <el-checkbox label="包月" value="month"/>
                  <el-checkbox label="包年" value="year"/>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item style="padding-left: 10px" label-width="110px" label="产品描述" prop="description">
                <el-input v-model="dataInfo.description" type="textarea"></el-input>
              </el-form-item>
              <el-form-item v-if="envList.length>0" label="实例环境" prop="description">
                <el-select v-model="dataInfo.env" multiple clearable>
                  <el-option
                      v-for="(item,index) in envList[0].value"
                      :label="item"
                      :value="item"
                  />
                </el-select>
              </el-form-item>
              <div v-if="IsCreate">
                <el-form-item style="padding-left: 10px" label-width="110px" label="捆绑销售" prop="binding">
                  <el-radio-group v-model="selectAi.binding">
                    <el-radio :value="0">
                      不绑定&nbsp;&nbsp;&nbsp;&nbsp;
                    </el-radio>
                    <el-radio :value="1">
                      绑定
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="数据集" prop="dataset">
                    <el-radio-group v-model="selectAi.isAllDataset">
                      <el-radio :value="1">
                        全部数据
                      </el-radio>
                      <el-radio :value="0">
                        选择数据
                      </el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item style="padding-left: 10px" label-width="110px" v-if="selectAi.isAllDataset === 0" label="选择数据集">
                    <el-select
                        v-model="selectAi.selectDataset"
                        multiple
                        clearable
                        filterable
                        placeholder="请选择/输入关键词"
                        style="width: 260px"
                        collapse-tags
                    >
                      <el-option
                          v-for="item in aiList.datasetList"
                          :key="item.product.id"
                          :label="item.product.name"
                          :value="JSON.stringify(item)"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                      style="padding-left: 10px" label-width="110px"
                  >
                    <el-input-number
                        v-model.number="selectAi.datasetDiscount"
                        autocomplete="off"
                        :min="0"
                        :max="100"
                        placeholder="0~100"
                    />
                  </el-form-item>
                </el-row>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="模型" prop="modelSet">
                    <el-radio-group v-model="selectAi.isAllModel">
                      <el-radio :value="1">
                        全部模型
                      </el-radio>
                      <el-radio :value="0">
                        选择模型
                      </el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item style="padding-left: 10px" label-width="110px" v-if="selectAi.isAllModel === 0" label="选择模型">
                    <el-select
                        v-model="selectAi.selectModel"
                        multiple
                        clearable
                        filterable
                        placeholder="请选择/输入关键词"
                        style="width: 260px"
                        collapse-tags
                    >
                      <el-option
                          v-for="item in aiList.modelList"
                          :key="item.product.id"
                          :label="item.product.name"
                          :value="JSON.stringify(item)"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item
                      style="padding-left: 10px" label-width="110px"
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                  >
                    <el-input-number
                        v-model.number="selectAi.modelDiscount"
                        autocomplete="off"
                        :min="0"
                        :max="100"
                        placeholder="0~100"
                    />
                  </el-form-item>
                </el-row>
              </div>
              <div v-else>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item label="绑定数据集">
                    <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectDataset">
                      {{ item.productName }}
                    </el-tag>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                  >
                    {{ selectAi.datasetDiscount }}
                  </el-form-item>
                </el-row>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item label="绑定模型">
                    <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectModel">
                      {{ item.productName }}
                    </el-tag>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                  >
                    {{ selectAi.modelDiscount }}
                  </el-form-item>
                </el-row>
              </div>

              <!--              云主机和容器的规格配置-->
              <div v-if="dataInfo.type!=1">
                <el-form-item label="操作系统" prop="os">
                  <el-select v-model="dataInfo.os" multiple clearable>
                    <el-option
                        v-for="(item,index) in osList"
                        :label="item.name"
                        :value="item.name"
                    />
                  </el-select>
                </el-form-item>
                <el-descriptions
                    :column="2"
                    :size="'default'"
                    :title="'基本配置'"
                    border
                    label-width="30px"
                    v-if="showConfigList"
                    class="margin-top"
                >
                  <el-descriptions-item label="基础配置">
                    <el-checkbox v-model="supplyGPU" v-if="IsCreate">
                      <span>提供GPU</span>
                    </el-checkbox>
                    <el-table
                        v-if="showConfigList"
                        :data="baseList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value[0]">
                            <el-option
                                v-for="(item, index) in baseAllList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="定义价格">
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.hour"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.day"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.month"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.year"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                  </el-descriptions-item>
                </el-descriptions>
                <el-descriptions
                    :column="1"
                    :title="'可选配置'"
                    border
                    label-width="30px"
                    v-if="showConfigList"
                    style="margin-top: 10px"
                >
                  <el-descriptions-item label="带宽配置">
                    <el-table
                        v-if="showConfigList"
                        :data="bandList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempBandList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <!--                      <el-table-column label="规格价格">-->
                      <!--                        <template #default="scope">-->
                      <!--                          <span v-if="dataInfo.billing.includes('hour')">{{-->
                      <!--                              '计时: ' + price.bandwidth.hour-->
                      <!--                            }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.bandwidth.day }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('month')">{{-->
                      <!--                              '包月: ' + price.bandwidth.month-->
                      <!--                            }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.bandwidth.year }}</span>-->
                      <!--                        </template>-->
                      <!--                      </el-table-column>-->
                      <el-table-column label="定义价格">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.ip.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.ip.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.ip.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.ip.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:60px;"> 元/Mbps </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                    <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="IP配置">
                    <el-table
                        v-if="showConfigList"
                        :data="IpList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempIpList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <el-table-column label="定义价格">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:60px;"> 元/个 </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                    <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="磁盘配置">
                    <el-table
                        v-if="showConfigList"
                        :data="DiskList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempDiskList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <!--                      <el-table-column label="规格价格">-->
                      <!--                        <template #default="scope">-->
                      <!--                          <span v-if="dataInfo.billing.includes('hour')">{{ '计时: ' + price.sysDisk.hour }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.sysDisk.day }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('month')">{{'包月: ' + price.sysDisk.month }}</span><br/>-->
                      <!--                          <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.sysDisk.year }}</span>-->
                      <!--                        </template>-->
                      <!--                      </el-table-column>-->
                      <el-table-column label="定义价格">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.sysDisk.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.sysDisk.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>

                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.sysDisk.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>

                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.sysDisk.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>

                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:75px;"> 元/GB </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                          <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-checkbox v-model="supplyDataDisk" v-if="IsCreate">
                      <span>提供数据盘</span>
                    </el-checkbox>
                    <el-table
                        v-if="showConfigList"
                        :data="DataList"
                    >
                      <el-table-column label="">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempDataList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.dataDisk.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.dataDisk.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.dataDisk.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.dataDisk.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:75px;"> 元/GB </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                          <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
              <!--              裸金属的规格配置-->
              <div v-if="dataInfo.type==1">
                <el-form ref="bareData" :model="bareData" :rules="IsCreate?dataRules:''" label-width="120px">
                  <el-form-item v-if="IsCreate" label="所属物理机" prop="bareData">
                    <el-select v-model="bareData" @change="HandleSelectBare">
                      <el-option
                          v-for="(item,index) in bareList"
                          :disabled="item.description!=undefined"
                          :label="'名称: '+item.name + (item.description!=undefined?' 已分配':'')"
                          :value="item"
                      />
                    </el-select>
                  </el-form-item>
                </el-form>
                <el-form-item label="操作系统" prop="os">
                  <el-select v-model="dataInfo.os" multiple clearable>
                    <el-option
                        v-for="(item,index) in osList"
                        :label="item.name"
                        :value="item.name"
                    />
                  </el-select>
                </el-form-item>
                <el-descriptions
                    :column="2"
                    :size="'default'"
                    :title="'基本配置'"
                    border
                    label-width="30px"
                    v-if="isBareChoose"
                    class="margin-top"
                >
                  <el-descriptions-item label="基础配置">
                    <el-table
                        v-if="showConfigList"
                        :data="baseList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <span>{{ scope.row.value[0] }}</span>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="定义价格">
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.hour"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.day"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.month"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                        <el-input-number controls-position="right" :precision="4" :min="0"
                                         v-model.number="price.base.year"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                  </el-descriptions-item>
                </el-descriptions>
                <el-descriptions
                    :column="1"
                    :title="'可选配置'"
                    border
                    label-width="30px"
                    v-if="isBareChoose"
                    style="margin-top: 10px"
                >
                  <el-descriptions-item label="带宽配置">
                    <el-table
                        v-if="showConfigList"
                        :data="bandList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempBandList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <el-table-column label="定义价格">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:60px;"> 元/Mbps </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                    <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="IP配置">
                    <el-table
                        v-if="showConfigList"
                        :data="IpList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempIpList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                      <el-table-column label="定义价格">
                        <template #default="scope">
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('hour')">
                              <span style="margin-right: 5px;">计时: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.hour" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('day')">
                              <span style="margin-right: 5px;">计日: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.day" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('month')">
                              <span style="margin-right: 5px;">计月: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.month" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <div style="display: flex;align-items: center; gap: 10px">
                            <div v-if="dataInfo.billing.includes('year')">
                              <span style="margin-right: 5px;">计年: </span>
                              <el-input-number v-model.number="unitPrice.bandwidth.year" controls-position="right"
                                               :precision="4" :min="0" style="width: 100px;"/>
                            </div>
                          </div>
                          <span style="margin-right: 5px;position: absolute;top:35%;right:60px;"> 元/个 </span>
                          <!--                          ”定义“按钮直接对每个规格的价格进行定义-->
                          <!--                    <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>-->
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>
            <!-- tab-pane 定制产品 -->
            <el-tab-pane :disabled="!IsCreate" label="定制产品" :name="'define'">
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="产品名称" prop="name">
                    <el-input style="width: 300px;" v-model="dataInfo.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="产品类别" prop="type">
                    <el-select  v-model="dataInfo.type" placeholder="请选择" style="width: 300px"
                               @change="handleChangeType">
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
                      <!--                      <el-option label="存储" :value="4"/>-->
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item style="padding-left: 10px" label-width="110px" :label="'支持的计费方式'">
                <el-checkbox-group :min="1" v-model="dataInfo.billing">
                  <el-checkbox label="计时" value="hour"/>
                  <el-checkbox label="包日" value="day"/>
                  <el-checkbox label="包月" value="month"/>
                  <el-checkbox label="包年" value="year"/>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item style="padding-left: 10px" label-width="110px" label="产品描述" prop="description">
                <el-input v-model="dataInfo.description" type="textarea"></el-input>
              </el-form-item>
              <el-form-item v-if="envList.length>0" label="实例环境" prop="description">
                <el-select v-model="dataInfo.env" multiple clearable>
                  <el-option
                      v-for="(item,index) in envList[0].value"
                      :label="item"
                      :value="item"
                  />
                </el-select>
              </el-form-item>
              <div v-if="IsCreate">
                <el-form-item style="padding-left: 10px" label-width="110px" label="捆绑销售" prop="binding">
                  <el-radio-group v-model="selectAi.binding">
                    <el-radio :value="0">
                      不绑定&nbsp;&nbsp;&nbsp;&nbsp;
                    </el-radio>
                    <el-radio :value="1">
                      绑定
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="数据集" prop="dataset">
                    <el-radio-group v-model="selectAi.isAllDataset">
                      <el-radio :value="1">
                        全部数据
                      </el-radio>
                      <el-radio :value="0">
                        选择数据
                      </el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item style="padding-left: 10px" label-width="110px" v-if="selectAi.isAllDataset === 0" label="选择数据集">
                    <el-select
                        v-model="selectAi.selectDataset"
                        multiple
                        clearable
                        filterable
                        placeholder="请选择/输入关键词"
                        style="width: 260px"
                        collapse-tags
                    >
                      <el-option
                          v-for="item in aiList.datasetList"
                          :key="item.ID"
                          :label="item.name"
                          :value="JSON.stringify(item)"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                      style="padding-left: 10px" label-width="110px"
                  >
                    <el-input-number
                        v-model.number="selectAi.datasetDiscount"
                        autocomplete="off"
                        :min="0"
                        :max="100"
                        placeholder="0~100"
                    />
                  </el-form-item>
                </el-row>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="模型" prop="modelSet">
                    <el-radio-group v-model="selectAi.isAllModel">
                      <el-radio :value="1">
                        全部模型
                      </el-radio>
                      <el-radio :value="0">
                        选择模型
                      </el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item style="padding-left: 10px" label-width="110px" v-if="selectAi.isAllModel === 0" label="选择模型">
                    <el-select
                        v-model="selectAi.selectModel"
                        multiple
                        clearable
                        filterable
                        placeholder="请选择/输入关键词"
                        style="width: 260px"
                        collapse-tags
                    >
                      <el-option
                          v-for="item in aiList.modelList"
                          :key="item.ID"
                          :label="item.name"
                          :value="JSON.stringify(item)"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                      style="padding-left: 10px" label-width="110px"
                  >
                    <el-input-number
                        v-model.number="selectAi.modelDiscount"
                        autocomplete="off"
                        :min="0"
                        :max="100"
                        placeholder="0~100"
                    />
                  </el-form-item>
                </el-row>
              </div>
              <div v-else>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="绑定数据集">
                    <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectDataset">
                      {{ item.productName }}
                    </el-tag>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                      style="padding-left: 10px" label-width="110px"
                  >
                    {{ selectAi.datasetDiscount }}
                  </el-form-item>
                </el-row>
                <el-row v-if="selectAi.binding === 1">
                  <el-form-item style="padding-left: 10px" label-width="110px" label="绑定模型">
                    <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectModel">
                      {{ item.productName }}
                    </el-tag>
                  </el-form-item>
                  <el-form-item
                      v-if="selectAi.binding === 1"
                      label="折扣率(%)"
                      style="padding-left: 10px" label-width="110px"
                  >
                    {{ selectAi.modelDiscount }}
                  </el-form-item>
                </el-row>
              </div>
              <el-form-item label="操作系统" prop="os">
                <el-select v-model="dataInfo.os" multiple clearable>
                  <el-option
                      v-for="(item,index) in osList"
                      :label="item.name"
                      :value="item.name"
                  />
                </el-select>
              </el-form-item>
              <div v-if="activeConfigName=='define'">
                <el-descriptions
                    :column="2"
                    :size="'large'"
                    :title="'基本配置'"
                    border
                    label-width="30px"
                    v-if="showConfigList"
                    class="margin-top"
                >
                  <el-descriptions-item label="基本配置">
                    <el-checkbox v-model="supplyGPU" v-if="IsCreate">
                      <span>提供GPU</span>
                    </el-checkbox>
                    <el-table
                        v-if="showConfigList"
                        :data="baseList"
                    >
                      <el-table-column label="名称" min-width="80px" align="left">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容" align="left">
                        <template #default="scope">
                          <el-select multiple clearable v-model="scope.row.value">
                            <el-option
                                v-for="(item, index) in baseAllList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"
                            />
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                </el-descriptions>
                <el-descriptions
                    :column="1"
                    :size="'large'"
                    :title="'可选配置'"
                    border
                    label-width="30px"
                    v-if="showConfigList"
                    style="margin-top: 10px"
                >
                  <el-descriptions-item label="带宽配置">
                    <el-table
                        v-if="showConfigList"
                        :data="bandList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempBandList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="IP配置">
                    <el-table
                        v-if="showConfigList"
                        :data="IpList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempIpList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="磁盘配置">
                    <el-table
                        v-if="showConfigList"
                        :data="DiskList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempDiskList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-checkbox v-model="supplyDataDisk" v-if="IsCreate">
                      <span>提供数据盘</span>
                    </el-checkbox>
                    <el-table
                        v-if="showConfigList"
                        :data="DataList"
                    >
                      <el-table-column label="">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <el-select v-model="scope.row.value" multiple>
                            <el-option
                                v-for="(item, index) in tempDataList[scope.$index].value"
                                :key="index"
                                :label="item"
                                :value="item"/>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>
          </el-tabs>

        </el-form>
        <div class="dialog-footer" style="margin: 10px auto;text-align: right;">
          <el-button @click="CloseDialog">取 消</el-button>
          <el-button type="warning" @click="handleConfirm">提 交</el-button>
        </div>
      </el-dialog>
      <!--      产品详情对话框-->
      <el-dialog v-model="detailVisible" :before-close="CloseDialog" :show-close="false" width="800px">
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">产品详情</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="detailVisible=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <!--        基本信息-->
        <el-descriptions
            :column="1"
            :size="'large'"
            :title="'基本信息'"
            class="margin-top"
        >
          <el-descriptions-item label="产品名称">
            {{ detailData.name }}
          </el-descriptions-item>
          <el-descriptions-item label="产品类别">
            {{ productType[detailData.type - 1] }}
          </el-descriptions-item>
          <el-descriptions-item label="产品说明">
            {{ detailData.description }}
          </el-descriptions-item>
        </el-descriptions>
        <!--        配置信息-->
        <el-descriptions
            :column="1"
            :size="'large'"
            :title="'规格配置'"
            class="margin-top"
            style="margin-top: 10px"
        >
          <el-descriptions-item>
            <el-table :data="configList" :header-cell-style="{background:'#eef1f6',color:'#606266'}">
              <el-table-column label="规格名称" prop="name"/>
              <el-table-column label="规格内容" prop="value">
                <template #default="scope">
                  <el-tag style="margin: 1px" v-for="(item,index) in scope.row.value">
                    {{item}}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-descriptions-item>
        </el-descriptions>
      </el-dialog>
      <!--          带宽价格配置对话框-->
      <el-dialog v-model="BandWidthPriceDialog" :before-close="handleConfirmEditConfig" :title="'价格定义'"
                 width="400px">
        <el-form label-width="120px">
          <div v-for="(item,index) in bandList[0].value">
            <h3>{{ item }}价格</h3>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.bandwidth.hour[index]"/>
              </el-form-item>
              <span class="unit">元</span>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.bandwidth.day[index]"/>
              </el-form-item>
              <span class="unit">元</span>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.bandwidth.month[index]"/>
              </el-form-item>
              <span class="unit">元</span>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.bandwidth.year[index]"/>
              </el-form-item>
              <span class="unit">元</span>
            </div>
          </div>
        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="handleConfirmEditConfig">取 消</el-button>
          <el-button type="warning" @click="handleConfirmEditConfig">确 定</el-button>
        </div>
      </el-dialog>
      <!--          磁盘价格配置对话框-->
      <el-dialog v-model="DiskPriceDialog" :before-close="handleConfirmEditConfig" :title="'价格定义'" width="400px">
        <el-form ref="dataInfo" :model="price.Disk" label-width="120px">
          <div v-for="(item,index) in DiskList[0].value">
            <h3>{{ item }}价格</h3>

            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.sysDisk.hour[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.sysDisk.day[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.sysDisk.month[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.sysDisk.year[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
          </div>
        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="handleConfirmEditConfig">取 消</el-button>
          <el-button type="warning" @click="handleConfirmEditConfig">确 定</el-button>
        </div>
      </el-dialog>
      <el-dialog v-model="DataPriceDialog" :before-close="handleConfirmEditConfig" :title="'价格定义'" width="400px">
        <el-form ref="dataInfo" :model="price.dataDisk" label-width="120px">
          <div v-for="(item,index) in DataList[0].value">
            <h3>{{ item }}价格</h3>

            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.dataDisk.hour[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.dataDisk.day[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.dataDisk.month[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
            <div style="display: flex;align-items: center; gap: 10px">
              <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                <el-input-number controls-position="right" :precision="4" :min="0"
                                 v-model.number="price.dataDisk.year[index]"/>
                <span class="unit">元</span>
              </el-form-item>
            </div>
          </div>
        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="handleConfirmEditConfig">取 消</el-button>
          <el-button type="warning" @click="handleConfirmEditConfig">确 定</el-button>
        </div>
      </el-dialog>

      <!--      审核意见-->
      <el-dialog v-model="auditVisible" :show-close="true" style="min-width: 500px" title="审核意见" width="40%" center>
        <span>{{ auditOpinion != '' ? auditOpinion : '无' }}</span>
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
  </div>
</template>
<script>
import {
  addComputingInfo,
  deleteComputingInfo,
  getComputingInfoList,
  syncedAllComputingInfo,
  syncedComputingInfo,
  updateComputingInfo,
} from "@/api/yunguan/product/productComputing";
import {deleteConfigInfo, getAllConfigInfoList} from "@/api/yunguan/product/productConfig";
import _ from 'lodash';
import {GetDockerImage} from "@/api/yunguan/cloudpods/image/image";
import {getElementProductList} from "@/api/yunguan/element/element";
import {getBaseBareHostList, updateBareHostDescription} from "@/api/cloudpods/baseRes/host";
import {getImageList} from "@/api/cloudpods/image/image";
import {getBaremetalDeviceList} from "@/api/cloudpods/baseRes/device"; // 使用 lodash 库进行深拷贝
export default {
  data() {
    return {
      supplyGPU: true,
      supplyDataDisk: true,
      auditVisible: false,
      auditOpinion: '',
      searchInfo: {
        page: 1,
        pageSize: 10,
        name: '',
      },
      total: 0,
      isSyncedNum: 0,
      billWay: ["hour", "day", "month", "year"],
      productType: ["裸金属", "云主机", "容器"],
      dataList: [],
      allConfigList: [],
      configList: [],
      envList: [],
      osList: [],
      aiList: [],
      selectAi: {
        selectDataset: [],
        selectModel: [],
        modelDiscount: 0,
        datasetDiscount: 0,
      },
      baseListBeforeCutOff: [],
      baseList: [],
      baseAllList: [],
      bandList: [],
      tempBandList: [],
      DiskList: [],
      tempDiskList: [],
      DataList: [],
      tempDataList: [],
      DataAllList: [],
      IpList: [],
      tempIpList: [],
      showConfigList: false,
      detailData: '',
      dialogVisible: false,
      detailVisible: false,
      dataInfo: {
        id: '',
        name: '',
        type: '',
        description: '',
        billing: ["hour"],
        price: '',
        specs: [],
        status: 1,
        isSynced: 1,
        isEnv: '',
        isDemand: '',
        isAi: '',
        isCustom: '',
        env: [],
        os: [],
        resourceId: '',
      },
      dataRules: {
        os: [
          {required: true, message: '请选择操作系统', trigger: 'blur'}
        ],
        name: [
          {required: true, message: '请输入产品名称', trigger: 'blur'}
        ],
        type: [
          {required: true, message: '请选择产品类别', trigger: 'blur'}
        ],
        bareData: [
          {required: true, message: '请选择裸金属', trigger: 'blur'}
        ]
      },
      IsCreate: true,
      activeName: 'all',
      activeConfigName: 'standard',
      unitPrice: {
        bandwidth: {
          hour: 0,
          day: 0,
          month: 0,
          year: 0,
        },
        ip: {
          hour: 0,
          day: 0,
          month: 0,
          year: 0,
        },
        sysDisk: {
          hour: 0,
          day: 0,
          month: 0,
          year: 0,
        },
        dataDisk: {
          hour: 0,
          day: 0,
          month: 0,
          year: 0,
        },
      },
      price: {
        base: {
          hour: 0,
          day: 0,
          month: 0,
          year: 0,
        },
        bandwidth: {
          hour: [],
          day: [],
          month: [],
          year: [],
        },
        sysDisk: {
          hour: [],
          day: [],
          month: [],
          year: [],
        },
        dataDisk: {
          hour: [],
          day: [],
          month: [],
          year: [],
        },

      },
      statusType: ["已保存", "已上架", "已下架"],//上架状态,
      basePriceDialog: false,
      BandWidthPriceDialog: false,
      DiskPriceDialog: false,
      DataPriceDialog: false,

      sysImageList: [],
      bareImageList: [],
      dockerImageList: [],
      bareList: [],
      bareData: null,
      isBareChoose: false,
      deleteVisible: false,
      currentProductId: null
    }
  },
  watch: {
    supplyGPU(newVal) {
      this.filterBaseList(newVal);
    },
    supplyDataDisk(newVal) {
      this.filterDataList(newVal);
    },
    // 由单价计算每种规格的价格
    unitPrice: {
      handler(newVal) {
        // 处理带宽各规格价格,带宽规格默认只有一种单位（Mbps）
        if (this.bandList && this.bandList.length > 0) {
          if (this.bandList[0]?.value) {
            for (let index = 0; index < this.bandList[0].value.length; index++) {
              // 获取规格值
              let spec = this.bandList[0].value[index]
              // 字符串中的获取规格单位
              let value = spec.match(/\d+/)[0]
              // let unit = spec.match(/[a-zA-Z]+/)[0]
              this.price.bandwidth.hour[index] = parseFloat((value * newVal.bandwidth.hour).toFixed(4))
              this.price.bandwidth.day[index] = parseFloat((value * newVal.bandwidth.day).toFixed(4))
              this.price.bandwidth.month[index] = parseFloat((value * newVal.bandwidth.month).toFixed(4))
              this.price.bandwidth.year[index] = parseFloat((value * newVal.bandwidth.year).toFixed(4))
            }
          }
        }
        // 处理系统盘各规格价格
        if (this.DiskList && this.DiskList.length > 0) {
          if (this.DiskList[0]?.value) {
            for (let index = 0; index < this.DiskList[0].value.length; index++) {
              // 获取规格值
              let spec = this.DiskList[0].value[index]
              // 字符串中的获取规格单位
              let value = spec.match(/\d+/)[0]
              let unit = spec.match(/[a-zA-Z]+/)[0]
              if (unit === 'TB') {
                value = value * 1024
              }
              this.price.sysDisk.hour[index] = parseFloat((value * newVal.sysDisk.hour).toFixed(4))
              this.price.sysDisk.day[index] = parseFloat((value * newVal.sysDisk.day).toFixed(4))
              this.price.sysDisk.month[index] = parseFloat((value * newVal.sysDisk.month).toFixed(4))
              this.price.sysDisk.year[index] = parseFloat((value * newVal.sysDisk.year).toFixed(4))
            }
          }
        }
        //  处理数据盘各规格价格
        if (this.DataList && this.DataList.length > 0) {
          if (this.DataList[0]?.value) {
            for (let index = 0; index < this.DataList[0].value.length; index++) {
              // 获取规格值
              let spec = this.DataList[0].value[index]
              // 字符串中的获取规格单位
              let value = spec.match(/\d+/)[0]
              let unit = spec.match(/[a-zA-Z]+/)[0]
              if (unit === 'TB') {
                value = value * 1024
              }
              this.price.dataDisk.hour[index] = parseFloat((value * newVal.dataDisk.hour).toFixed(4))
              this.price.dataDisk.day[index] = parseFloat((value * newVal.dataDisk.day).toFixed(4))
              this.price.dataDisk.month[index] = parseFloat((value * newVal.dataDisk.month).toFixed(4))
              this.price.dataDisk.year[index] = parseFloat((value * newVal.dataDisk.year).toFixed(4))
            }
          }
        }
      },
      deep: true,
      immediate: true,
    }
  },
  created() {
    this.GetProductList()
    this.GetAllConfigInfoList()
    this.GetImageList()
    this.GetModelList()
    this.GetDatasetList()
  },
  methods: {
    handleDialogClose() {
      this.deleteVisible = false;
      this.currentProductId = null;
    },
    // 打开删除弹窗（统一入口）
    openDeleteDialog(ID) {
      this.currentProductId = ID;
      console.log(ID, 'id');
      this.deleteVisible = true;
    },
    async confirmDelete() {
      try {
        let resp;
        resp = await deleteComputingInfo(this.currentProductId);

        // 严格判断接口成功状态
        if (resp.code === 0) {
          // 显示对应成功消息
          this.$message({
            type: 'success',
            message: '删除成功'
          });
          this.deleteVisible = false;


          // 执行成功后的操作
          await this.GetProductList(); // 刷新表格
        } else {
          // 接口返回失败状态时显示错误
          this.$message({
            showClose: true,
            message: resp.msg || '删除失败',
            type: 'error'
          });
          return; // 失败时直接返回，不执行后续操作
        }
      } catch (error) {
        console.error('删除过程中出现异常:', error);
      }
    },

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
      this.GetProductList();
    },
    //获取裸金属GPU信息
    GetBaremetalGpuList() {
      getBaremetalDeviceList({hostId: this.dataInfo.resourceId}).then(res => {
        if (res.code == 0) {
          let tempData = res.data;
          if (res.data.num > 0) {
            this.baseList.push({
              name: "GPU型号",
              tag: "gpuType",
              value: [tempData.model]
            })
            this.baseList.push({
              name: "GPU数量",
              tag: "gpuNumber",
              value: [tempData.num + '卡']
            })
          }
        }
      })
    },
    handleTagType(status) {
      switch (status) {
        case -1:
        case 0:
          return 'danger'
        case 1:
          return 'success'
        case 2:
          return 'blue'
        default:
          return 'info'
      }
    },
    handleTagIsSynced(status) {
      switch (status) {
        case 1:
          return 'danger'
        case 2:
          return 'blue'
        default:
          return 'info'
      }
    },
    filterBaseList(newVal) {
      if (newVal) {
        this.baseList = this.baseListBeforeCutOff
      } else {
        this.baseListBeforeCutOff = this.baseList
        this.baseList = this.baseList.filter(item => this.shouldDisplayOption(item));

      }

    },
    filterDataList(newVal) {
      if (newVal) {
        this.DataList = this.DataAllList;
      } else {
        this.DataList = [];
      }
    },
    // 判断是否显示选项
    shouldDisplayOption(row) {
      let isGPU = row.name.includes('GPU')
      return this.supplyGPU || !isGPU;
    },
    //格式化容量
    FormatCapacitySize(val) {
      var G = 1024
      var T = 1024 * 1024
      if (val < T) {
        return (val / G).toFixed(1) + "G";
      } else if (val >= G && val <= T) {
        return (val / G).toFixed(1) + "G";
      } else if (val > T) {
        return (val / T).toFixed(1) + "T";
      }
    },
    //选择物理机
    HandleSelectBare() {
      this.isBareChoose = true;
      this.baseList = [];
      this.dataInfo.resourceId = this.bareData.id;
      this.GetBaremetalGpuList()
      this.baseList.push({
        name: "CPU",
        tag: "cpu",
        value: [this.bareData.cpu_count + "核"]
      })
      this.baseList.push({
        name: "内存",
        tag: "memory",
        value: [this.FormatCapacitySize(this.bareData.mem_size)]
      })
      this.baseList.push({
        name: "磁盘",
        tag: "disk",
        value: [this.FormatCapacitySize(this.bareData.storage_size)]
      })
    },
    //获取物理机列表
    GetBareMetalList() {
      let obj = {}
      getBaseBareHostList(obj).then(response => {
        if (response.code == 0) {
          this.bareList = response.data.data
        }
      })
    },
    OpenBandWidthPriceDialog() {
      this.BandWidthPriceDialog = true
    },
    OpenDiskPriceDialog() {
      this.DiskPriceDialog = true
    },
    OpenDataPriceDialog() {
      this.DataPriceDialog = true
    },

    //展示详情
    showDetail(val) {
      this.detailVisible = true
      this.detailData = val;
      this.configList = val.specs
    },
    handleTabsClick() {
      switch (this.activeName) {
        case "all":
          this.searchInfo.type = 0
          break;
        case "baremetal":
          this.searchInfo.type = 1
          break;
        case "virtual":
          this.searchInfo.type = 2
          break;
        case "container":
          this.searchInfo.type = 3
          break;
        case "storage":
          this.searchInfo.type = 4
          break;
      }
      this.selectAi = {
        selectDataset: [],
        selectModel: [],
        datasetDiscount: 0,
        modelDiscount: 0,
      }
      this.searchInfo.name = "";
      this.GetProductList()
    },
    handleConfigTabsClick() {
      this.supplyGPU = true
      this.supplyDataDisk = true
      switch (this.activeConfigName) {
        case "standard":
          this.dataInfo.isCustom = 0
          if (this.baseList.length > 0) {
            for (let i = 0; i < this.baseList.length; i++) {
              this.baseList[i].value = [this.baseList[i].value[0]]
            }
          }
          break;
        case "define":
          this.dataInfo.isCustom = 1
          this.bareData = ""
          this.GetListData()
          break;
      }
      this.dataInfo.status = 1
      this.dataInfo.payWay = 1
      this.dataInfo.isSynced = 1
      // this.dataInfo.isCustom = 0
      this.dataInfo.isAi = 0
      this.dataInfo.isDemand = 0
      this.dataInfo.isEnv = 0
      this.dataInfo.resourceId = "";
      this.dataInfo.billing = ["hour"]
      this.price.base = {}
      // this.price.bandwidth.hour = []
      // this.price.bandwidth.day = []
      // this.price.bandwidth.month = []
      // this.price.bandwidth.year = []
      // this.price.sysDisk.hour = []
      // this.price.sysDisk.day = []
      // this.price.sysDisk.month = []
      // this.price.sysDisk.year = []
      // this.price.dataDisk.hour = []
      // this.price.dataDisk.day = []
      // this.price.dataDisk.month = []
      // this.price.dataDisk.year = []
      this.specs = []
      // this.showConfigList = false

      this.selectAi = {
        selectDataset: [],
        selectModel: [],
        datasetDiscount: 0,
        modelDiscount: 0,
      }
      this.selectAi.binding = 0;
    },
    handleConfirmEditConfig() {
      this.DiskPriceDialog = false
      this.DataPriceDialog = false
      this.BandWidthPriceDialog = false
    },
    syncSelections(newSelections) {
      // 清空之前的选中项
      this.configData.price = [];
      const sortedSelections = newSelections.sort((a, b) => {
        return this.valueList.indexOf(a) - this.valueList.indexOf(b);
      });
      this.configData.value = sortedSelections;
      // 同步新的选中项
      sortedSelections.forEach((selectedItem, selectedIndex) => {
        const correspondingIndex = this.valueList.indexOf(selectedItem);
        if (correspondingIndex !== -1 && correspondingIndex < this.priceList.length) {
          this.configData.price.push(this.priceList[correspondingIndex]);
        }
      });
    },
    GetListData() {
      this.configList = _.cloneDeep(this.allConfigList).filter(item => item.type.includes(this.dataInfo.type))
      // this.osList = _.cloneDeep(this.configList).filter(item => item.tag == "os")
      this.bandList = _.cloneDeep(this.configList).filter(item => item.tag == "bandwidth")
      this.DiskList = _.cloneDeep(this.configList).filter(item => item.tag == "sysDisk")
      this.DataList = _.cloneDeep(this.configList).filter(item => item.tag == "dataDisk")
      this.DataAllList = _.cloneDeep(this.configList).filter(item => item.tag == "dataDisk")
      this.envList = _.cloneDeep(this.configList).filter(item => item.tag == "env")
      this.IpList = _.cloneDeep(this.configList).filter(item => item.tag == "ip")
      this.baseList = _.cloneDeep(this.configList).filter(item => item.tag === "cpu" || item.tag === "gpu" || item.tag === "memory" || item.tag === "cpuType" || item.tag === "gpuType" || item.tag === "gpuNumber" || item.tag === "gpuMemory")
      this.baseAllList = _.cloneDeep(this.configList).filter(item => item.tag === "cpu" || item.tag === "gpu" || item.tag === "memory" || item.tag === "cpuType" || item.tag === "gpuType" || item.tag === "gpuNumber" || item.tag === "gpuMemory")
      this.tempBandList = _.cloneDeep(this.bandList);
      this.tempDataList = _.cloneDeep(this.DataList);
      this.tempDiskList = _.cloneDeep(this.DiskList);
      this.tempIpList = _.cloneDeep(this.IpList);
    },
    //获取镜像列表
    GetImageList() {
      let obj = {}
      obj.pageSize = 1000;
      obj.pageIndex = 1;
      getImageList(obj).then(res => {
        if (res.code == 0) {
          this.sysImageList = res.data.data
        }
      })
      obj.description = "裸金属"
      obj.description = "裸金属"
      //获取裸金属镜像
      getImageList(obj).then(res => {
        if (res.code == 0) {
          this.bareImageList = res.data.data
        }
      })
      GetDockerImage().then(res => {
        if (res.code == 0) {
          this.dockerImageList = res.data.list ? res.data.list : []
        }
      })
    },
    handleChangeType() {
      this.isBareChoose = false
      this.showConfigList = true
      this.dataInfo.resourceId = "";
      this.GetListData()
      // if (this.dataInfo.type === 1) {
      //   this.DiskList = []
      //   this.bandList = []
      // }
      if (this.baseList.length > 0) {
        for (let i = 0; i < this.baseList.length; i++) {
          this.baseList[i].value = [this.baseList[i].value[0]]
        }
      }
      if (this.DiskList.length > 0) {
        this.price.sysDisk.hour = []
        this.price.sysDisk.day = []
        this.price.sysDisk.month = []
        this.price.sysDisk.year = []
        for (let i = 0; i < this.DiskList[0].value.length; i++) {
          this.price.sysDisk.hour.push(0)
          this.price.sysDisk.day.push(0)
          this.price.sysDisk.month.push(0)
          this.price.sysDisk.year.push(0)
        }
      }
      if (this.bandList.length > 0) {
        this.price.bandwidth.hour = []
        this.price.bandwidth.day = []
        this.price.bandwidth.month = []
        this.price.bandwidth.year = []
        for (let i = 0; i < this.bandList[0].value.length; i++) {
          this.price.bandwidth.hour.push(0)
          this.price.bandwidth.day.push(0)
          this.price.bandwidth.month.push(0)
          this.price.bandwidth.year.push(0)
        }
      }
      if (this.DataList.length > 0) {
        this.price.dataDisk.hour = []
        this.price.dataDisk.day = []
        this.price.dataDisk.month = []
        this.price.dataDisk.year = []
        for (let i = 0; i < this.DataList[0].value.length; i++) {
          this.price.dataDisk.hour.push(0)
          this.price.dataDisk.day.push(0)
          this.price.dataDisk.month.push(0)
          this.price.dataDisk.year.push(0)
        }
      }
      if (this.dataInfo.type === 3) {
        this.dataInfo.os = []
        this.osList = this.dockerImageList
      } else if (this.dataInfo.type === 1) {
        this.dataInfo.os = []
        this.osList = this.bareImageList
      } else {
        this.dataInfo.os = []
        this.osList = this.sysImageList
      }
    },
    //获取规格配置
    GetAllConfigInfoList() {
      getAllConfigInfoList().then((response) => {
        if (response.code == 0) {
          this.allConfigList = response.data.list;
        } else {
          this.$message({
            message: "获取失败",
            type: "error",
          });
        }
      });
    },
    // 获取所有已上架模型列表
    GetModelList() {
      getElementProductList({page: 1, pageSize: 10000, subType: 5, isSynced: 1}).then((response) => {
        if (response.code === 0) {
          this.aiList.modelList = response.data.list;
        } else {
          this.$message({
            message: "获取失败",
            type: "error",
          });
        }
      })
    },
    // 获取所有已上架数据集列表
    GetDatasetList() {
      getElementProductList({page: 1, pageSize: 10000, subType: 6, isSynced: 1}).then((response) => {
        if (response.code === 0) {
          this.aiList.datasetList = response.data.list;
        } else {
          this.$message({
            message: "获取失败",
            type: "error",
          });
        }
      })
    },
    FormatDateTime(dateString) {
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
    },
    handleSearch() {
      this.GetProductList()
    },
    onReset() {
      this.searchInfo.name = ''
      this.searchInfo.type = null
      this.searchInfo.isCustom = null
      this.searchInfo.auditStatus = null
      this.searchInfo.isSynced = null
      // 重置页码
      this.searchInfo.page = 1
      this.GetProductList()
    },
    //同步全部
    handleSyncedAll() {
      syncedAllComputingInfo().then(response => {
        if (response.code == 0) {
          this.$message({
            message: '发送同步请求成功，等待响应',
            type: 'success'
          })
          this.GetProductList(this.searchInfo)
        } else {
          this.$message({
            message: '发送同步请求失败，请检查网络连接',
            type: 'error'
          });
        }
      })
    },
    //同步单个产品
    handleSyncedProduct(val) {
      this.dataInfo = _.cloneDeep(val)
      this.dataInfo.isSynced = 2
      this.SyncedProductInfo()
    },
    //上下架产品
    handleStatusChange(val) {
      this.dataInfo = val
      if (val.status != 2) {
        this.dataInfo.status = 2
      } else {
        this.dataInfo.status = 3
      }
      this.UpdateProductInfo()
    },
    //新增
    OpenDialog() {
      this.GetListData()
      //是否提供GPU是可选的
      this.supplyGPU = true
      //是否提供数据盘是可选的
      this.supplyDataDisk = true
      this.bareData = null
      this.dataInfo.id = 0;
      this.dataInfo.name = "";
      this.dataInfo.type = "";
      this.dataInfo.description = "";
      this.dataInfo.resourceId = "";
      this.dataInfo.os = [];
      this.isBareChoose = false;
      this.dataInfo.status = 1
      this.dataInfo.payWay = 1
      this.dataInfo.isSynced = 1
      this.dataInfo.isCustom = this.activeConfigName === 'standard' ? 0 : 1
      this.dataInfo.isAi = 0
      this.dataInfo.isDemand = 0
      this.dataInfo.isEnv = 0
      this.dataInfo.billing = ["hour"]
      this.price.base.hour = 0
      this.price.base.day = 0
      this.price.base.month = 0
      this.price.base.year = 0
      this.price.bandwidth.hour = []
      this.price.bandwidth.day = []
      this.price.bandwidth.month = []
      this.price.bandwidth.year = []
      this.price.sysDisk.hour = []
      this.price.sysDisk.day = []
      this.price.sysDisk.month = []
      this.price.sysDisk.year = []

      this.price.dataDisk.hour = []
      this.price.dataDisk.day = []
      this.price.dataDisk.month = []
      this.price.dataDisk.year = []
      this.specs = []

      this.selectAi = {
        selectDataset: [],
        selectModel: [],
        datasetDiscount: 0,
        modelDiscount: 0,
      }
      this.selectAi.binding = 0;
      this.GetBareMetalList()
      this.dialogVisible = true
      this.IsCreate = true
      this.showConfigList = false
    },
    //取消
    CloseDialog() {
      this.dialogVisible = false
      this.detailVisible = false
    },
    //编辑
    handleEdit(val) {
      this.dataInfo = _.cloneDeep(val)
      if (val.isCustom == 0) this.activeConfigName = "standard"
      else this.activeConfigName = "define"
      this.price = val.price
      this.GetListData()
      let dataList = val.specs

      this.osList = this.dockerImageList

      this.baseList = dataList.filter(item => item.tag != "bandwidth" && item.tag != "sysDisk" && item.tag != "dataDisk" && item.tag != "ip")
      this.bandList = dataList.filter(item => item.tag == "bandwidth")
      this.DiskList = dataList.filter(item => item.tag == "sysDisk")
      this.DataList = dataList.filter(item => item.tag == "dataDisk")
      if (this.dataInfo.type == 1) {
        this.isBareChoose = true;
        this.baseList = dataList
      }
      if (val.bind.model.length > 0) {
        this.selectAi.selectModel = val.bind.model;
        this.selectAi.selectDataset = val.bind.dataset;
        this.selectAi.modelDiscount = val.bind.model[0].discount;
        this.selectAi.datasetDiscount = val.bind.dataset[0].discount;
        this.selectAi.binding = 1;
        if (val.bind.dataset[0].aiProductId == -1) {
          this.selectAi.isAllDataset = 1;
        } else this.selectAi.isAllDataset = 0;
        if (val.bind.model[0].aiProductId == -1) {
          this.selectAi.isAllModel = 1;
        } else this.selectAi.isAllModel = 0;
      }
      this.IsCreate = false
      this.dialogVisible = true
      this.showConfigList = true
    },
    handleConfirm() {
      if (this.IsCreate) this.AddProductInfo()
      else this.UpdateProductInfo()
    },
    //新增产品
    AddProductInfo() {
      const dataset = this.selectAi.isAllDataset
          ? [{
            aiProductId: -1,
            productName: "全部数据集",
            reduction: 0,
            discount: this.selectAi.datasetDiscount
          }]
          : this.selectAi.selectDataset.map(item => {
            let obj = JSON.parse(item)
            return {
              aiProductId: obj.ID,
              productName: obj.name,
              reduction: 0,
              discount: this.selectAi.datasetDiscount
            }
          })
      const model = this.selectAi.isAllModel
          ? [{
            aiProductId: -1,
            productName: "全部模型",
            reduction: 0,
            discount: this.selectAi.modelDiscount
          }]
          : this.selectAi.selectModel.map(item => {
            let obj = JSON.parse(item)
            return {
              aiProductId: obj.ID,
              productName: obj.name,
              reduction: 0,
              discount: this.selectAi.modelDiscount
            }
          })
      this.dataInfo.bind = {model, dataset}
      const filteredBaseList = this.baseList.filter(item => {
        // 检查 value 是否为非空数组，并且不包含 '无'
        return Array.isArray(item.value) && item.value.length > 0 && !item.value.includes('无');
      });
      this.dataInfo.specs = filteredBaseList.concat(this.bandList, this.DiskList, this.DataList, this.IpList)
      this.dataInfo.price = this.price
      this.dataInfo.isCustom = this.activeConfigName === 'standard' ? 0 : 1
      this.dataInfo.storage = [{name: "本地存储", tag: "block", value: []}, {name: "文件存储", tag: "nas", value: []}]
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addComputingInfo(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '创建成功',
                type: 'success'
              })
              this.GetProductList(this.searchInfo)
              this.dialogVisible = false
              this.supplyGPU = true
              if (this.dataInfo.type == 1) {
                let obj = {}
                obj.id = this.dataInfo.resourceId;
                obj.description = "已分配";
                updateBareHostDescription(obj).then(res => {
                  if (res.code == 0) {
                  }
                })
              }
            } else {
              this.$message({
                message: '创建失败',
                type: 'error'
              });
            }
          })
        }
      });
    },
    //删除产品
    DeleteProductInfo(val) {
      this.$confirm('此操作将永久删除该产品, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteComputingInfo(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetProductList(this.searchInfo)
            if (val.type == 1) {
              let obj = {}
              obj.id = val.resourceId;
              obj.description = "";
              updateBareHostDescription(obj).then(res => {
                if (res.code == 0) {
                }
              })
            }
          } else {
            this.$message({
              message: '删除失败',
              type: 'error'
            });
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    },
    //修改产品
    UpdateProductInfo() {
      const dataset = this.selectAi.selectDataset;
      const model = this.selectAi.selectModel;
      this.dataInfo.bind = {model, dataset}
      this.dataInfo.specs = this.baseList.concat(this.bandList, this.DiskList, this.DataList)
      this.dataInfo.price = this.price
      this.dataInfo.auditStatus = 0
      this.dataInfo.isCustom = this.activeConfigName === 'standard' ? 0 : 1
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          updateComputingInfo(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '修改成功',
                type: 'success'
              })
              this.GetProductList(this.searchInfo)
              this.dialogVisible = false
              this.supplyGPU = true
            } else {
              this.$message({
                message: '修改失败',
                type: 'error'
              });
            }
          })
        }
      })
    },
    //同步产品
    SyncedProductInfo() {
      syncedComputingInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '同步成功',
            type: 'success'
          })
          this.GetProductList(this.searchInfo)
        } else {
          this.$message({
            message: '同步失败',
            type: 'error'
          });
        }
      })
    },
    GetProductList() {
      getComputingInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.searchInfo.page = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.total = response.data.total
          this.isSyncedNum = response.data.isSynced
        }
      })
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetProductList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetProductList()
    }
  }
}
</script>
<style scoped>

.unit {
  position: absolute;
  right: 40px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none; /* 确保用户不能与单位文本交互 */
  color: #606266; /* 与输入框文本颜色一致 */
}

.blue {
  position: relative;
  z-index: 0;
  width: 8px;
  height: 8px;
  background-blend-mode: normal;
  mix-blend-mode: normal;
  display: flex;
  background: rgba(22, 93, 255, 1);
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


</style>
