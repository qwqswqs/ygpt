<template>
  <div>
    <div class="gva-table-box">
      <el-tabs v-model="activeName" class="tabs" @click="handleTabsClick">
        <el-tab-pane label="全部" name="all"></el-tab-pane>
        <el-tab-pane label="裸金属" name="baremetal"></el-tab-pane>
        <el-tab-pane label="云主机" name="virtual"></el-tab-pane>
        <el-tab-pane label="容器" name="container"></el-tab-pane>
<!--        <el-tab-pane label="存储" name="storage"></el-tab-pane>-->
      </el-tabs>
      <div class="gva-btn-list" style="margin-top: 10px">
        <el-form ref="searchForm" :inline="true" :model="searchInfo">
          <el-form-item>
            <el-input v-model="searchInfo.name" placeholder="请输入产品名称"/>
          </el-form-item>
          <el-form-item>
            <el-button  type="primary" @click="handleSearch">查询</el-button>
            <el-button  type="primary" @click="OpenDialog">添加</el-button>
            <el-button  :disabled="isSyncedNum==0" type="primary" @click="handleSyncedAll">立即同步
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格展示 -->
      <el-table :data="dataList" row-key="ID">
        <el-table-column label="名称">
          <template #default="scope">
            <el-button type="text" @click="showDetail(scope.row)">{{ scope.row.name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ productType[scope.row.type - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明" >
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
            <span>{{ statusType[scope.row.status - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="同步状态">
          <template #default="scope">
            <span>{{ scope.row.isSynced == 1 ? '未同步' : '已同步' }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column label="支付方式">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.payWay == 1 ? '微信' : '支付宝' }}</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" :disabled="scope.row.isSynced==2" @click="handleSyncedProduct(scope.row)">同步
            </el-button>
            <!--            <el-button type="text" @click="handleStatusChange(scope.row)">-->
            <!--              {{ scope.row.status != 2 ? '上架' : '下架' }}-->
            <!--            </el-button>-->
            <el-button type="text" @click="DeleteProductInfo(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
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
      <el-dialog v-model="dialogVisible" :title="IsCreate? '创建算力产品':'编辑算力产品'" width="1000px">
        <el-form ref="dataInfo" :model="dataInfo" :rules="IsCreate?dataRules:''" label-width="120px">

          <!--          <el-form-item>-->
          <!--            <el-checkbox v-model="dataInfo.isEnv" :true-value="1" :false-value="0" label="是否允许自定义环境"/>-->
          <!--            <el-checkbox v-model="dataInfo.isDemand" :true-value="1" :false-value="0" label="是否允许自定义需求"/>-->
          <!--            <el-checkbox v-model="dataInfo.isAi" :true-value="1" :false-value="0" label="是否允许自定义AI"/>-->
          <!--          </el-form-item>-->
          <el-tabs v-model="activeConfigName" class="tabs" @tabChange="handleConfigTabsClick">
            <!-- tab-pane 标准产品 -->
            <el-tab-pane :disabled="!IsCreate" label="标准产品" name="standard">
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="产品名称" prop="name">
                    <el-input v-model="dataInfo.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="产品类别" prop="type">
                    <el-select :disabled="!IsCreate" v-model="dataInfo.type" placeholder="请选择" style="width: 100%"
                               @change="handleChangeType">
                      <el-option v-for="(item,index) in productType"
                                 :key="index"
                                 :label="item"
                                 :value="index+1"/>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item :label="'支持的计费方式'">
                <el-checkbox-group :min="1" v-model="dataInfo.billing">
                  <el-checkbox label="计时" value="hour"/>
                  <el-checkbox label="包日" value="day"/>
                  <el-checkbox label="包月" value="month"/>
                  <el-checkbox label="包年" value="year"/>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="产品描述" prop="description">
                <el-input v-model="dataInfo.description" type="textarea"></el-input>
              </el-form-item>
              <el-form-item v-if="osList.length>0" label="操作系统" prop="os">
                <el-select v-model="dataInfo.os" multiple>
                  <el-option
                      v-for="(item,index) in osList"
                      :label="item.name"
                      :value="item.name"
                  />
                </el-select>
              </el-form-item>
              <el-form-item v-if="envList.length>0" label="实例环境" prop="description">
                <el-select v-model="dataInfo.env" multiple>
                  <el-option
                      v-for="(item,index) in envList[0].value"
                      :label="item"
                      :value="item"
                  />
                </el-select>
              </el-form-item>
              <!--              <el-form-item v-if="envList.length>0" label="实例环境" prop="description">-->
              <!--                <el-select v-model="dataInfo.env" multiple>-->
              <!--                  <el-option-->
              <!--                      v-for="(item,index) in envList[0].value"-->
              <!--                      :label="item"-->
              <!--                      :value="item"-->
              <!--                  />-->
              <!--                </el-select>-->
              <!--              </el-form-item>-->

              <el-form-item label="捆绑销售" prop="binding">
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
                <el-form-item label="数据集" prop="dataset">
                  <el-radio-group v-model="selectAi.isAllDataset">
                    <el-radio :value="1">
                      全部数据
                    </el-radio>
                    <el-radio :value="0">
                      选择数据
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="selectAi.isAllDataset === 0" label="选择数据集">
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
                >
                  <el-input-number
                      v-model="selectAi.datasetDiscount"
                      autocomplete="off"
                      :min="0"
                      :max="100"
                      placeholder="0~100"
                  />
                </el-form-item>
              </el-row>
              <el-row v-if="selectAi.binding === 1">
                <el-form-item label="模型" prop="modelSet">
                  <el-radio-group v-model="selectAi.isAllModel">
                    <el-radio :value="1">
                      全部模型
                    </el-radio>
                    <el-radio :value="0">
                      选择模型
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="selectAi.isAllModel === 0" label="选择模型">
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

              <el-descriptions
                  :column="2"
                  :size="'default'"
                  :title="'基本配置'"
                  border
                  v-if="showConfigList"
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
                  <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                    <el-input type="number" :min="0" v-model.number="price.base.hour"/>
                  </el-form-item>
                  <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                    <el-input type="number" :min="0" v-model.number="price.base.day"/>
                  </el-form-item>
                  <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                    <el-input type="number" :min="0" v-model.number="price.base.month"/>
                  </el-form-item>
                  <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                    <el-input type="number" :min="0" v-model.number="price.base.year"/>
                  </el-form-item>
                </el-descriptions-item>
              </el-descriptions>
              <el-descriptions
                  :column="1"
                  :size="'large'"
                  :title="'可选配置'"
                  border
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
                        <span>{{ scope.row.value }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="规格价格">
                      <template #default="scope">
                        <span v-if="dataInfo.billing.includes('hour')">{{ '计时: ' + price.bandwidth.hour }}</span><br/>
                        <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.bandwidth.day }}</span><br/>
                        <span v-if="dataInfo.billing.includes('month')">{{'包月: ' + price.bandwidth.month }}</span><br/>
                        <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.bandwidth.year }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="定义价格">
                      <template #default="scope">
                        <el-button type="text" @click="OpenBandWidthPriceDialog">定义</el-button>
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
                        <span>{{ scope.row.value }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="规格价格">
                      <template #default="scope">
                        <span v-if="dataInfo.billing.includes('hour')">{{ '计时: ' + price.sysDisk.hour }}</span><br/>
                        <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.sysDisk.day }}</span><br/>
                        <span v-if="dataInfo.billing.includes('month')">{{ '包月: ' + price.sysDisk.month }}</span><br/>
                        <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.sysDisk.year }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="定义价格">
                      <template #default="scope">
                        <el-button type="text" @click="OpenDiskPriceDialog">定义</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
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
                        <span>{{ scope.row.value }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="">
                      <template #default="scope">
                        <span v-if="dataInfo.billing.includes('hour')">{{ '计时: ' + price.dataDisk.hour }}</span><br/>
                        <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.dataDisk.day }}</span><br/>
                        <span v-if="dataInfo.billing.includes('month')">{{ '包月: ' + price.dataDisk.month }}</span><br/>
                        <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.dataDisk.year }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="">
                      <template #default="scope">
                        <el-button type="text" @click="OpenDataPriceDialog">定义</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </el-descriptions-item>
              </el-descriptions>
            </el-tab-pane>
            <!-- tab-pane 定制产品 -->
            <el-tab-pane :disabled="!IsCreate" label="定制产品" name="define">
              <el-row style="margin-top:20px">
                <el-col :span="12">
                  <el-form-item label="产品名称" prop="name">
                    <el-input v-model="dataInfo.name" placeholder="请输入产品名称"/>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="产品类别" prop="type">
                    <el-select v-model="dataInfo.type" placeholder="请选择" style="width: 100%"
                               @change="handleChangeType">
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
<!--                      <el-option label="存储" :value="4"/>-->
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item :label="'支持的计费方式'">
                <el-checkbox-group :min="1" v-model="dataInfo.billing">
                  <el-checkbox label="计时" value="hour"/>
                  <el-checkbox label="包日" value="day"/>
                  <el-checkbox label="包月" value="month"/>
                  <el-checkbox label="包年" value="year"/>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="产品描述" prop="description">
                <el-input v-model="dataInfo.description" type="textarea"></el-input>
              </el-form-item>
              <el-form-item v-if="osList.length>0" label="操作系统" prop="description">
                <el-select v-model="dataInfo.os" multiple>
                  <el-option
                      v-for="(item,index) in osList"
                      :label="item.name"
                      :value="item.name"
                  />
                </el-select>
              </el-form-item>
              <el-form-item v-if="envList.length>0" label="实例环境" prop="description">
                <el-select v-model="dataInfo.env" multiple>
                  <el-option
                      v-for="(item,index) in envList[0].value"
                      :label="item"
                      :value="item"
                  />
                </el-select>
              </el-form-item>

              <el-form-item label="捆绑销售" prop="binding">
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
                <el-form-item label="数据集" prop="dataset">
                  <el-radio-group v-model="selectAi.isAllDataset">
                    <el-radio :value="1">
                      全部数据
                    </el-radio>
                    <el-radio :value="0">
                      选择数据
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="selectAi.isAllDataset === 0" label="选择数据集">
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
                <el-form-item label="模型" prop="modelSet">
                  <el-radio-group v-model="selectAi.isAllModel">
                    <el-radio :value="1">
                      全部模型
                    </el-radio>
                    <el-radio :value="0">
                      选择模型
                    </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="selectAi.isAllModel === 0" label="选择模型">
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

              <el-descriptions
                  :column="2"
                  :size="'large'"
                  :title="'基本配置'"
                  border
                  v-if="showConfigList"
                  class="margin-top"
              >
                <el-descriptions-item>
                  <el-table
                      v-if="showConfigList"
                      :data="baseList"
                  >
                    <el-table-column label="名称" min-width="80px"align="left">
                      <template #default="scope">
                        <span>{{ scope.row.name }}</span>
                      </template>
                    </el-table-column>
                    <el-table-column label="规格内容"align="left">
                      <template #default="scope">
                        <el-select multiple v-model="scope.row.value">
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
                        <span>{{ scope.row.value }}</span>
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
                        <span>{{ scope.row.value }}</span>
                      </template>
                    </el-table-column>
                  </el-table>
                </el-descriptions-item>
              </el-descriptions>
            </el-tab-pane>
          </el-tabs>

        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="CloseDialog">取 消</el-button>
          <el-button type="warning" @click="handleConfirm">提 交</el-button>
        </div>
      </el-dialog>
      <!--      产品详情对话框-->
      <el-dialog v-model="detailVisible" :before-close="CloseDialog" :title="'产品详情'" width="800px">
        <!--        基本信息-->
        <el-descriptions
            :column="2"
            :size="'large'"
            :title="'基本信息'"
            border
            class="margin-top"
        >
          <el-descriptions-item>
            <template #label>
              <div class="cell-item">
                产品名称
              </div>
            </template>
            {{ detailData.name }}
          </el-descriptions-item>
          <el-descriptions-item>
            <template #label>
              <div class="cell-item">
                产品类别
              </div>
            </template>
            {{ productType[detailData.type - 1] }}
          </el-descriptions-item>
          <el-descriptions-item>
            <template #label>
              <div class="cell-item">
                产品说明
              </div>
            </template>
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
              <el-table-column label="规格内容" prop="value"/>
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
            <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
              <el-input type="number" :min="0" v-model.number="price.bandwidth.hour[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
              <el-input type="number" :min="0" v-model.number="price.bandwidth.day[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
              <el-input type="number" :min="0" v-model.number="price.bandwidth.month[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
              <el-input type="number" :min="0" v-model.number="price.bandwidth.year[index]"/>
            </el-form-item>
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
            <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
              <el-input type="number" :min="0" v-model.number="price.sysDisk.hour[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
              <el-input type="number" :min="0" v-model.number="price.sysDisk.day[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
              <el-input type="number" :min="0" v-model.number="price.sysDisk.month[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
              <el-input type="number" :min="0" v-model.number="price.sysDisk.year[index]"/>
            </el-form-item>
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
            <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
              <el-input type="number" :min="0" v-model.number="price.dataDisk.hour[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
              <el-input type="number" :min="0" v-model.number="price.dataDisk.day[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
              <el-input type="number" :min="0" v-model.number="price.dataDisk.month[index]"/>
            </el-form-item>
            <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
              <el-input type="number" :min="0" v-model.number="price.dataDisk.year[index]"/>
            </el-form-item>
          </div>
        </el-form>
        <div class="dialog-footer" style="margin: 0 auto;text-align: right;">
          <el-button @click="handleConfirmEditConfig">取 消</el-button>
          <el-button type="warning" @click="handleConfirmEditConfig">确 定</el-button>
        </div>
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
import {getAllConfigInfoList} from "@/api/yunguan/product/productConfig";
import _ from 'lodash';
import {GetDockerImage, GetSystemImage} from "@/api/yunguan/cloudpods/image/image";
import {GetBareMetalList} from "@/api/yunguan/cloudpods/baremetal/baremetal";
import {getAllDataSets, getAllModels} from "@/api/yunguan/element/element"; // 使用 lodash 库进行深拷贝
export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
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
      },
      baseList: [],
      baseAllList: [],
      bandList: [],
      DiskList: [],
      DataList: [],
      showConfigList: '',
      detailData: '',
      dialogVisible: false,
      detailVisible: false,
      bareList:[],
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
      },
      dataRules: {
        name: [
          {required: true, message: '请输入产品名称', trigger: 'blur'}
        ],
        type: [
          {required: true, message: '请选择产品类别', trigger: 'blur'}
        ]
      },
      IsCreate: true,
      activeName: 'all',
      activeConfigName: 'standard',
      price: {
        base: {
          hour: '',
          day: '',
          month: '',
          year: '',
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
      dockerImageList: [],
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
      }
      this.GetProductList()
    },
    handleConfigTabsClick() {
      switch (this.activeConfigName) {
        case "standard":
          this.dataInfo.isCustom = 0
          break;
        case "define":
          this.dataInfo.isCustom = 1
          break;
      }
      // window.alert(this.activeConfigName)
      // window.alert(this.dataInfo.isCustom)

      this.dataInfo = {}
      this.dataInfo.status = 1
      this.dataInfo.payWay = 1
      this.dataInfo.isSynced = 1
      // this.dataInfo.isCustom = 0
      this.dataInfo.isAi = 0
      this.dataInfo.isDemand = 0
      this.dataInfo.isEnv = 0
      this.dataInfo.billing = ["hour"]
      this.price.base = {}
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
      this.showConfigList = false

      this.selectAi = {
        selectDataset: [],
        selectModel: [],
      }

      this.GetListData()
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
      this.envList = _.cloneDeep(this.configList).filter(item => item.tag == "env")
      this.baseList = _.cloneDeep(this.configList).filter(item => item.tag === "cpu" || item.tag === "gpu" || item.tag === "memory"  || item.tag === "cpuType" || item.tag === "gpuType" || item.tag === "gpuNumber" || item.tag === "gpuMemory")
      this.baseAllList = _.cloneDeep(this.configList).filter(item => item.tag === "cpu" || item.tag === "gpu" || item.tag === "memory" || item.tag === "cpuType" || item.tag === "gpuType" || item.tag === "gpuNumber" || item.tag === "gpuMemory")
    },
    //获取镜像列表
    GetImageList() {
      GetSystemImage().then(res => {
        if (res.code == 0) {
          this.sysImageList = res.data.list ? JSON.parse(res.data.list) : []
          // console.log(this.sysImageList)
        }
      })
      GetDockerImage().then(res => {
        if (res.code == 0) {
          this.dockerImageList = res.data.list ? res.data.list : []
          console.log(this.dockerImageList)
        }
      })
    },
    //获取物理机列表
    GetBareMetalList(){
      GetBareMetalList({type:"physicalMachine"}).then(res=>{
        if(res.code==0){
          this.bareList=res.data.list
        }
      })
    },
    handleChangeType() {
      this.showConfigList = true
      this.GetListData()
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
        for (let i = 0; i < this.bandList[0].value.length; i++) {
          this.price.dataDisk.hour.push(0)
          this.price.dataDisk.day.push(0)
          this.price.dataDisk.month.push(0)
          this.price.dataDisk.year.push(0)
        }
      }
      if (this.dataInfo.type === 3) {
        this.dataInfo.os = []
        this.osList = this.dockerImageList
      } else {
        this.dataInfo.os = []
        this.osList = this.sysImageList
      }

      if (this.dataInfo.isCustom == 0 && this.dataInfo.type==1){
        this.GetBareMetalList()
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
      getAllModels().then((response) => {
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
      getAllDataSets().then((response) => {
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
      this.dataInfo = val
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
      this.dataInfo = {}
      this.dataInfo.status = 1
      this.dataInfo.payWay = 1
      this.dataInfo.isSynced = 1
      this.dataInfo.isCustom = this.activeConfigName === 'standard' ? 0 : 1
      this.dataInfo.isAi = 0
      this.dataInfo.isDemand = 0
      this.dataInfo.isEnv = 0
      this.dataInfo.billing = ["hour"]
      this.price.base = {}
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
      this.IsCreate = true
      this.showConfigList = false
      this.dialogVisible = true

      this.selectAi = {
        selectDataset: [],
        selectModel: [],
      }

      this.GetListData()
    },
    //取消
    CloseDialog() {
      this.dialogVisible = false
      this.detailVisible = false
    },
    //编辑
    handleEdit(val) {
      this.dataInfo = val
      if (val.isCustom == 0) this.activeConfigName = "standard"
      else this.activeConfigName = "define"
      this.showConfigList = true
      this.price = val.price
      this.IsCreate = false
      this.dialogVisible = true
      this.GetListData()
      let dataList = val.specs
      console.log(dataList)

      this.osList = this.dockerImageList

      this.baseList = dataList.filter(item => item.tag != "bandwidth" && item.tag != "sysDisk")
      this.bandList = dataList.filter(item => item.tag == "bandwidth")
      this.DiskList = dataList.filter(item => item.tag == "sysDisk")
      console.log(this.baseList)
      console.log(this.bandList)
      console.log(this.DiskList)
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
      this.dataInfo.specs = this.baseList.concat(this.bandList, this.DiskList, this.DataList)
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
      const dataset = this.selectAi.isAllDataset
          ? [{
            aiProductId: -2,
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
      this.dataInfo.specs = this.baseList.concat(this.bandList, this.DiskList, this.DataList)
      this.dataInfo.price = this.price
      this.dataInfo.isCustom = this.activeConfigName === 'standard' ? 0 : 1
      updateComputingInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '修改成功',
            type: 'success'
          })
          this.GetProductList(this.searchInfo)
          this.dialogVisible = false
        } else {
          this.$message({
            message: '修改失败',
            type: 'error'
          });
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
