<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-table-box">
      <el-tabs v-model="activeName" class="tabs no-select" @click="handleTabsClick">
        <el-tab-pane label="云产品" name="host">
          <div class="gva-search-box" style="margin-top: 16px">
            <div class="gva-btn-list">
              <el-form ref="searchForm" :model="searchInfo" :inline="true" style="display: flex;justify-content: space-between">
                <el-form-item>
                  <el-form-item label="名称">
                    <el-input v-model="searchInfo.name" placeholder="请输入"/>
                  </el-form-item>
                  <el-form-item label="类别">
                    <el-select style="width: 240px;" v-model="searchInfo.type" placeholder="请输入" clearable>
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="产品类型">
                    <el-select style="width: 240px;" v-model="searchInfo.isCustom" placeholder="请输入" clearable>
                      <el-option label="标准产品" :value="0"/>
                      <el-option label="定制产品" :value="1"/>
                    </el-select>
                  </el-form-item>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleSearch">查询</el-button>
                  <el-button type="info" @click="searchInfo.name='';searchInfo.isCustom=null;searchInfo.type=null;handleSearch()">重置</el-button>
                </el-form-item>
              </el-form>
            </div>
            <!-- 条件选择 -->

          </div>
          <div style="margin: 10px;">
            <!-- 表格展示 -->
            <el-table :data="dataList" row-key="ID" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
              <el-table-column label="名称">
                <template #default="scope">
                  <el-button type="text" link @click="handleEdit(scope.row)">{{ scope.row.name }}</el-button>
                </template>
              </el-table-column>
              <el-table-column label="类别">
                <template #default="scope">
                  <span>{{ productType[scope.row.type - 1] }}</span>
                </template>
              </el-table-column>
              <el-table-column label="说明">
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
              <el-table-column label="操作" width="140" fixed="right">
                <template #default="scope">
                  <el-button link type="text" @click="OpenAuditDialog(scope.row,1)">审核</el-button>
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
        </el-tab-pane>
        <el-tab-pane label="AI产品" name="ai">
          <div class="gva-search-box" style="margin-top: 16px">
            <div class="gva-btn-list">
              <el-form ref="searchForm" :model="searchInfo" :inline="true" style="display: flex;justify-content: space-between">
                <el-form-item>
                  <el-form-item label="名称">
                    <el-input
                        v-model="searchAiInfo.name"
                        placeholder="请输入"
                    />
                  </el-form-item>
                  <el-form-item label="类型">
                    <el-select v-model="searchInfo.isCustom" placeholder="请输入" clearable>
                      <el-option label="标准产品" :value="0"/>
                      <el-option label="定制产品" :value="1"/>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="类别">
                    <el-select v-model="searchInfo.type" placeholder="请输入" clearable>
                      <el-option label="裸金属" :value="1"/>
                      <el-option label="云主机" :value="2"/>
                      <el-option label="容器" :value="3"/>
                    </el-select>
                  </el-form-item>

                </el-form-item>

                <el-form-item>
                  <el-button
                      type="primary"
                      @click="handleSearch"
                  >
                    查询
                  </el-button>
                  <el-button type="info" @click="handleSearch">重置</el-button>
                </el-form-item>
              </el-form>
            </div>
            <!-- 条件选择 -->

          </div>
          <div style="margin: 10px;">
            <!-- 表格展示 -->
            <el-table
                :data="allAiList"
                row-key="ID"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
            >
              <el-table-column align="left" class="text-overflow" label="标题" min-width="200" prop="product.name">
                <template #default="scope">
                  <div class="text-overflow text-details" @click="handleDetailsClick(scope.row)">
                    {{ scope.row.product.name }}
                  </div>
                </template>
              </el-table-column>

              <el-table-column align="left" label="类型" min-width="200" prop="product.description">
                <template #default="scope">
                  <span>{{ scope.row.product.subtype==5?'模型':'数据集' }}</span>
                </template>
              </el-table-column>
              <el-table-column align="left" label="描述" min-width="200" prop="product.description">
                <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.product.description }}</span>
                </template>
              </el-table-column>
              <el-table-column align="left" label="价格（元）" min-width="120" prop="product.price" />

              <el-table-column align="left" label="类别" min-width="200" prop="sysType.name" />
              <el-table-column :fixed="'right'" label="操作" min-width="150"align="left">
                <template #default="scope">
                  <el-button type="text" @click="OpenAuditDialog(scope.row,2)">审核</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="gva-pagination">
              <el-pagination
                  :current-page="searchAiInfo.pageIndex"
                  :page-size="searchAiInfo.pageSize"
                  :page-sizes="[10, 30, 50, 100]"
                  :total="searchAiInfo.total"
                  layout="total,prev, pager, next,  sizes, jumper"
                  @current-change="handleAiCurrentChange"
                  @size-change="handleAiSizeChange"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      <!--      云产品详情对话框-->
      <el-dialog v-model="dialogVisible" class="no-select" :title="'产品详情'"
                 width="1000px">
        <el-form disabled ref="dataInfo" :model="dataInfo" :rules="IsCreate?dataRules:''" label-width="120px">
          <el-tabs v-model="activeConfigName" class="tabs no-select">
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
                    <el-select :disabled="!IsCreate" v-model="dataInfo.type" placeholder="请选择" style="width: 300px;" >
                      <el-option v-for="(item,index) in productType"
                                 :key="index"
                                 :label="item"
                                 :value="index+1"/>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item :label="'计费方式'">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in dataInfo.billing">{{billType[item]}}</el-tag>
              </el-form-item>
              <el-form-item label="产品描述" prop="description">
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
              <el-row v-if="selectAi.binding === 1">
                <el-form-item label="绑定数据集">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectDataset">{{item.productName}}</el-tag>
                </el-form-item>
                <el-form-item
                    v-if="selectAi.binding === 1"
                    label="折扣率(%)"
                >
                  {{selectAi.datasetDiscount}}
                </el-form-item>
              </el-row>
              <el-row v-if="selectAi.binding === 1">
                <el-form-item label="绑定模型">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectModel">{{item.productName}}</el-tag>
                </el-form-item>
                <el-form-item
                    v-if="selectAi.binding === 1"
                    label="折扣率(%)"
                >
                  {{selectAi.modelDiscount}}
                </el-form-item>
              </el-row>

              <!--              云主机和容器的规格配置-->
              <div v-if="dataInfo.type!=1">
                <el-form-item label="操作系统" prop="os">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in dataInfo.os">{{item}}</el-tag>
                </el-form-item>
                <el-descriptions
                    :column="2"
                    :size="'default'"
                    :title="'基本配置'"
                    border
                    class="margin-top"
                >
                  <el-descriptions-item label-width="100px" label="基础配置">
                    <el-table
                        :data="baseList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="定义价格">
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('hour')" label="计时价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.hour"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.day"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.month"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.year"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                  </el-descriptions-item>
                </el-descriptions>
                <el-descriptions
                    :column="1"
                    :size="'large'"
                    :title="'可选配置'"
                    border
                    style="margin-top: 10px"
                >
                  <el-descriptions-item label-width="100px"   label="带宽配置">
                    <el-table
                        :data="bandList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格价格">
                        <template #default="scope">
                          <span v-if="dataInfo.billing.includes('hour')">{{
                              '计时: ' + price.bandwidth.hour
                            }}</span><br/>
                          <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.bandwidth.day }}</span><br/>
                          <span v-if="dataInfo.billing.includes('month')">{{
                              '包月: ' + price.bandwidth.month
                            }}</span><br/>
                          <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.bandwidth.year }}</span>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="磁盘配置">
                    <el-table
                        :data="DiskList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格价格">
                        <template #default="scope">
                          <span v-if="dataInfo.billing.includes('hour')">{{ '计时: ' + price.sysDisk.hour }}</span><br/>
                          <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.sysDisk.day }}</span><br/>
                          <span v-if="dataInfo.billing.includes('month')">{{
                              '包月: ' + price.sysDisk.month
                            }}</span><br/>
                          <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.sysDisk.year }}</span>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-table
                        :data="DataList"
                    >
                      <el-table-column label="">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <span v-if="dataInfo.billing.includes('hour')">{{
                              '计时: ' + price.dataDisk.hour
                            }}</span><br/>
                          <span v-if="dataInfo.billing.includes('day')">{{ '包日: ' + price.dataDisk.day }}</span><br/>
                          <span v-if="dataInfo.billing.includes('month')">{{
                              '包月: ' + price.dataDisk.month
                            }}</span><br/>
                          <span v-if="dataInfo.billing.includes('year')">{{ '包年: ' + price.dataDisk.year }}</span>
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
                          :label="'名称: '+item.name"
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
                    v-if="isBareChoose"
                    class="margin-top"
                >
                  <el-descriptions-item label-width="100px"  label="基础配置">
                    <el-table
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
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.hour"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('day')" label="包日价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.day"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('month')" label="包月价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.month"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                    <div style="display: flex;align-items: center; gap: 10px">
                      <el-form-item v-if="dataInfo.billing.includes('year')" label="包年价格">
                        <el-input-number controls-position="right" :precision="2" :min="0"
                                         v-model.number="price.base.year"/>
                        <span class="unit">元</span>
                      </el-form-item>
                    </div>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>
            <!-- tab-pane 定制产品 -->
            <el-tab-pane :disabled="!IsCreate" label="定制产品" :name="'define'">
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
              <el-form-item :label="'计费方式'">
                <el-tag style="margin-left: 5px" v-for="(item,index) in dataInfo.billing">{{billType[item]}}</el-tag>
              </el-form-item>
              <el-form-item label="产品描述" prop="description">
                <el-input v-model="dataInfo.description" type="textarea"></el-input>
              </el-form-item>
              <el-form-item label="操作系统" prop="os">
                <el-tag style="margin-left: 5px" v-for="(item,index) in dataInfo.os">{{item}}</el-tag>
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

              <el-row v-if="selectAi.binding === 1">
                <el-form-item label="绑定数据集">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectDataset">{{item.productName}}</el-tag>
                </el-form-item>
                <el-form-item
                    v-if="selectAi.binding === 1"
                    label="折扣率(%)"
                >
                  {{selectAi.datasetDiscount}}
                </el-form-item>
              </el-row>
              <el-row v-if="selectAi.binding === 1">
                <el-form-item label="绑定模型">
                  <el-tag style="margin-left: 5px" v-for="(item,index) in selectAi.selectModel">{{item.productName}}</el-tag>
                </el-form-item>
                <el-form-item
                    v-if="selectAi.binding === 1"
                    label="折扣率(%)"
                >
                  {{selectAi.modelDiscount}}
                </el-form-item>
              </el-row>

              <div v-if="activeConfigName=='define'">
                <el-descriptions
                    :column="2"
                    :size="'large'"
                    :title="'基本配置'"
                    border
                    class="margin-top"
                >
                  <el-descriptions-item label-width="100px" >
                    <el-table
                        :data="baseList"
                    >
                      <el-table-column label="名称" min-width="80px"align="left">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容"align="left">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
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
                    style="margin-top: 10px"
                >
                  <el-descriptions-item label-width="100px"  label="带宽配置">
                    <el-table
                        :data="bandList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                  <el-descriptions-item label="磁盘配置">
                    <el-table
                        :data="DiskList"
                    >
                      <el-table-column label="名称">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="规格内容">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-table
                        :data="DataList"
                    >
                      <el-table-column label="">
                        <template #default="scope">
                          <span>{{ scope.row.name }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column label="">
                        <template #default="scope">
                          <el-tag style="margin-left: 2px" v-for="(item,key) in scope.row.value">
                            {{item}}
                          </el-tag>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-tab-pane>
          </el-tabs>

        </el-form>
      </el-dialog>
      <!--      AI产品详情对话框-->
      <el-dialog v-model="aiVisible" class="no-select" :title="'产品详情'"
                 width="500px">
        <el-descriptions
            :size="'default'"
            :column="1"
            border
            class="margin-top"
        >
          <el-descriptions-item label-width="100px" label="名称">
            {{aiData.name}}
          </el-descriptions-item>
          <el-descriptions-item label-width="100px" label="类型">
            {{aiType[aiData.type-1]}}
          </el-descriptions-item>
          <el-descriptions-item label-width="100px" label="描述">
            {{aiData.description}}
          </el-descriptions-item>
          <el-descriptions-item label-width="100px" label="标签">
            <el-tag
                style="margin-left: 2px"
                v-for="(tag,index) in aiData.usage"
                :key="index"
            >
              <span> {{ tag }}</span>
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label-width="100px" label="开放类型">
            {{openList[aiData.openType]}}
          </el-descriptions-item>
          <el-descriptions-item label-width="100px" label="价格">
            {{aiData.price}}
          </el-descriptions-item>
        </el-descriptions>
      </el-dialog>
      <!--      审核对话框-->
      <el-dialog v-model="auditVisible" :show-close="false" class="no-select"
                 width="520px">
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">审核</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="auditVisible=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close />
          </el-icon>
        </div>
        <el-form ref="auditForm" :model="auditForm" :rules="dataRules" >
          <el-form-item label="审核状态"  prop="auditStatus">
            <el-radio-group v-model="auditForm.auditStatus">
              <el-radio :value="1" >
                通过
              </el-radio>
              <el-radio :value="-1" >
                驳回
              </el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item
              label="审核意见"
              prop="auditOpinion"
          >
            <el-input
                style="width: 320px;height: 56px"
                type="textarea"
                v-model="auditForm.auditOpinion"
                placeholder="请输入"
                :maxlength="200"
                show-word-limit
            />
          </el-form-item>
        </el-form>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="auditVisible=false">取消</el-button>
          <el-button
              type="primary"
              @click="ConfirmAudit"
          >确定</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>
<script>
import {getComputingInfoList, updateComputingInfo,} from "@/api/yunguan/product/productComputing";
import {getAllConfigInfoList} from "@/api/yunguan/product/productConfig";
import _ from 'lodash';
import {GetDockerImage} from "@/api/yunguan/cloudpods/image/image";
import {
  getAllDataSets,
  getAllModels,
  getElementProductList,
  getNodeElementList,
  updateElement, updateProduct
} from "@/api/yunguan/element/element";
import {getBaseBareHostList} from "@/api/cloudpods/baseRes/host";
import {getImageList} from "@/api/cloudpods/image/image"; // 使用 lodash 库进行深拷贝
export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        auditStatus:0,
        name: '',
        total: 0,
      },
      searchAiInfo: {
        pageIndex: 1,
        pageSize: 10,
        auditStatus:0,
        status:0,
        source: 1, //来源：1节点 2租户
        name: '',
        total: 0,
      },
      openList: ['请选择开放类别', "不开放", "节点内可用", "算力可见节点内可用", "算力可见算力外可用"],
      billWay: ["hour", "day", "month", "year"],
      productType: ["裸金属", "云主机", "容器"],
      aiType: ["模型", "数据集", "算法"],
      dataList: [],
      allConfigList: [],
      configList: [],
      envList: [],
      osList: [],
      aiList: [],
      selectAi: {
        selectDataset: [],
        selectModel: [],
        modelDiscount:0,
        datasetDiscount:0,
      },
      baseList: [],
      baseAllList: [],
      bandList: [],
      DiskList: [],
      DataList: [],
      detailData: '',
      dialogVisible: false,
      billType:{
        ['hour']:'计时',
        ['day']:'按日',
        ['month']:'按月',
        ['year']:'按年',
      },
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
        ],
        auditStatus: [
          {required: true, message: '请选择审核状态', trigger: 'blur'}
        ],
        auditOpinion: [
          {required: true, message: '请输入审核意见', trigger: 'blur'}
        ]
      },
      IsCreate: true,
      activeName: 'host',
      activeConfigName: 'standard',
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
      sysImageList: [],
      bareImageList: [],
      dockerImageList: [],
      bareList: [],
      bareData: null,
      isBareChoose: false,
      allAiList: [],
      auditVisible: false,
      aiVisible: false,
      auditForm:{
        auditType:1,//1是云产品 2是AI产品
        ID:0,
        auditStatus:1,//1通过 -1驳回
        status:1,
        auditOpinion:'',//审核意见
      },
      aiData:{},
    }
  },
  created() {
    this.GetProductList()
    this.GetAllConfigInfoList()
    this.GetImageList()
    this.GetModelList()
    this.GetDatasetList()
    this.GetAiList()
    this.GetBareMetalList()
  },
  methods: {
    //打开ai产品详情对话框
    OpenDetailDialog(val) {
      this.aiData=val;
      this.aiVisible=true;
    },
    //打开审核对话框
    OpenAuditDialog(val,type){
      this.auditForm.auditType = type;
      if (type ===1) {
        this.auditForm.ID = val.ID;
      }else if (type ===2){
        this.auditForm.ID = val.product.id;
        this.auditForm.id = val.product.id;
      }
      this.auditForm.auditOpinion = '';
      this.auditForm.auditStatus = 1;
      this.auditVisible=true;
    },
    //审核确定
    ConfirmAudit(){
      this.$refs.auditForm.validate(valid => {
        if (valid) {
          if (this.auditForm.auditType==1){
            updateComputingInfo(this.auditForm).then(res=>{
              if (res.code == 0){
                this.$message({
                  type:'success',
                  message:'审核完成',
                })
                this.GetProductList()
                this.auditVisible=false;
              }
            })
          }else{
            this.auditForm.status=this.auditForm.auditStatus;
            updateProduct(this.auditForm).then(res=>{
              if (res.code == 0){
                this.$message({
                  type:'success',
                  message:'审核完成',
                })
                this.GetAiList()
                this.auditVisible=false;
              }
            })
          }
        }
      })
    },
    //获取Ai列表
    GetAiList() {
      getElementProductList(this.searchAiInfo).then(response => {
        if (response.code == 0) {
          this.allAiList = response.data.list
          this.searchAiInfo.pageIndex = response.data.page
          this.searchAiInfo.pageSize = response.data.pageSize
          this.searchAiInfo.total = response.data.total
        }
      })
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
    handleTabsClick() {
      switch (this.activeName) {
        case "ai":
          this.GetAiList();
          break;
        case "host":
          this.GetProductList();
          break;
      }
    },
    GetListData() {
      this.configList = _.cloneDeep(this.allConfigList).filter(item => item.type.includes(this.dataInfo.type))
      console.log(this.allConfigList)
      console.log(this.configList)
      console.log(this.dataInfo.type)
      // this.osList = _.cloneDeep(this.configList).filter(item => item.tag == "os")
      this.bandList = _.cloneDeep(this.configList).filter(item => item.tag == "bandwidth")
      this.DiskList = _.cloneDeep(this.configList).filter(item => item.tag == "sysDisk")
      this.DataList = _.cloneDeep(this.configList).filter(item => item.tag == "dataDisk")
      this.baseList = _.cloneDeep(this.configList).filter(item => item.tag != "bandwidth" && item.tag != "sysDisk" && item.tag != "dataDisk")
      this.baseAllList = _.cloneDeep(this.configList).filter(item => item.tag != "bandwidth" && item.tag != "sysDisk" && item.tag != "dataDisk")
      console.log(this.baseAllList)
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
          this.bareImageList =res.data.data
        }
      })
      GetDockerImage().then(res => {
        if (res.code == 0) {
          this.dockerImageList = res.data.list ? res.data.list : []
          console.log(this.dockerImageList)
        }
      })
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
    //编辑
    handleEdit(val) {
      this.dataInfo = _.cloneDeep(val)
      if (val.isCustom == 0) this.activeConfigName = "standard"
      else this.activeConfigName = "define"
      this.price = val.price
      this.GetListData()
      let dataList = val.specs

      this.osList = this.dockerImageList

      this.baseList = dataList.filter(item => item.tag != "bandwidth" && item.tag != "sysDisk" && item.tag != "dataDisk")
      this.bandList = dataList.filter(item => item.tag == "bandwidth")
      this.DiskList = dataList.filter(item => item.tag == "sysDisk")
      this.DataList = dataList.filter(item => item.tag == "dataDisk")
      if (this.dataInfo.type == 1) {
        this.isBareChoose = true;
        this.baseList = dataList
      }
      this.IsCreate = false
      this.dialogVisible = true
      this.showConfigList = true
      if (val.bind.model.length > 0) {
        this.selectAi.selectModel=val.bind.model;
        this.selectAi.selectDataset=val.bind.dataset;
        this.selectAi.modelDiscount=val.bind.model[0].discount;
        this.selectAi.datasetDiscount=val.bind.dataset[0].discount;
        this.selectAi.binding=1;
        if (val.bind.dataset[0].aiProductId==-1){
          this.selectAi.isAllDataset=1;
        }else this.selectAi.isAllDataset=0;
        if (val.bind.model[0].aiProductId==-1){
          this.selectAi.isAllModel=1;
        }else this.selectAi.isAllModel=0;
      }
    },
    GetProductList() {
      getComputingInfoList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.dataList = response.data.list
          this.searchInfo.page = response.data.page
          this.searchInfo.pageSize = response.data.pageSize
          this.searchInfo.total = response.data.total
          this.isSyncedNum = response.data.isSynced
        }
      })
    },
    handleAiCurrentChange(val) {
      this.searchAiInfo.pageIndex = val
      this.GetAiList()
    },
    handleAiSizeChange(val) {
      this.searchAiInfo.pageSize = val
      this.GetAiList()
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetProductList()
    },
    handleDetailsClick(val){

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

</style>
