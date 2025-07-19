<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div style="background-color: #fff;font-family: 'Noto Sans SC';">
      <el-tabs v-model="activeName" @tabChange="handleTabChange">

        <el-tab-pane label="标准产品" :name="'标准'"></el-tab-pane>
        <el-tab-pane label="定制产品" :name="'定制'"></el-tab-pane>
        <!--      <el-tab-pane label="投诉建议" :name="'投诉'"></el-tab-pane>-->
        <!--      <el-tab-pane label="合同发票" :name="'合同'"></el-tab-pane>-->
      </el-tabs>
      <div class="gva-search-box" style="margin-top:16px">
        <div class="gva-btn-list" style="">
          <el-form
              ref="searchForm"
              style="display: flex;justify-content: space-between"
              :inline="true"
              :model="searchInfo"
          >
            <el-form-item label="状态">
              <el-select v-model="searchInfo.status"
                         style="width: 200px;background-color: rgba(242, 243, 245, 1) !important;" placeholder="请选择">
                <el-option
                    v-for="item in [{ label: '全部', value: 0}, { label: '未处理', value: 2}, { label: '已处理', value: 3}]"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
              </el-select>
            </el-form-item>
            <!--          <el-form-item>-->
            <!--            <el-input-->
            <!--                v-model="searchInfo.name"-->
            <!--                placeholder="请输入名称关键字"-->
            <!--            />-->
            <!--          </el-form-item>-->
            <el-form-item>
              <el-button
                  type="primary"
                  @click="handleSearch"
                  style="background-color:rgba(22, 93, 255, 1);"
              >
                查询
              </el-button>
              <el-button class="button-gap" type="info"
                         @click="onReset">
                重置
              </el-button>
              <!--            <el-button-->
              <!--                -->
              <!--                type="primary"-->
              <!--                @click="addHardwareOpenDialog"-->
              <!--            >-->
              <!--              新增-->
              <!--            </el-button>-->
            </el-form-item>
          </el-form>
        </div>
      </div>
      <div class="gva-table-box">

        <!-- 表格展示 -->
        <el-table
            :data="ticketList"
            row-key="ID"
            @sort-change="handleSortChange"
            :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)'}"
        >
          <el-table-column label="工单创建时间" min-width="120px" sortable="custom">
            <template #default="scope">
              <span>{{ FormatDateTime(scope.row.CreatedAt) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态">
            <template #default="scope">
              <!--            <el-tag :type="handleTagType(scope.row.status)">-->
              <!--              -->
              <!--            </el-tag>-->
              <div style="display: flex;align-items: center">
                <div :class="handleTagType(scope.row.status)"></div>
                <div>{{ handleStatus[scope.row.status] }}</div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="所属租户">
            <template #default="scope">
              <span class="text-overflow">{{ GetRenterName(scope.row.renterID) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="处理人">
            <template #default="scope">
              <span class="text-overflow">{{ GetUserName(scope.row.handlePerson) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="100px" fixed="right">
            <template #default="scope">
              <el-button
                  class="text-blue-500"
                  @click="OpenHandleDrawer(scope.row)"
                  link
              >
                {{ scope.row.status == 1 ? '处理' : '详情' }}
              </el-button>
              <el-button
                  class="text-blue-500"
                  link
                  @click="HandleDeleteDialog(scope.row,'delete')"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
              :current-page="searchInfo.page"
              :page-size="searchInfo.pageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="total"
              layout="total,  prev, pager, next, sizes,jumper"
              @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
          />
        </div>
      </div>

      <el-drawer
          v-model="openDialog"
          style="min-width: 600px;"
          size="520px"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="false"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg"
                style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ IsCreated ? '新增工单' : '处理工单' }}</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="closeDialog"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <el-form
            ref="dataInfo"
            :model="dataInfo"
            label-width="120"
            label-position="left"
            :rules="rules"
        >
          <el-form-item label="创建时间">
            <span>{{ FormatDateTime(dataInfo.CreatedAt) }}</span>
          </el-form-item>
          <el-form-item label="所属租户">
            <span>{{ GetRenterName(dataInfo.renterID) }}</span>
          </el-form-item>
          <el-form-item label="工单类别">
            <span>{{ ticketType[dataInfo.ticketType - 1] }}</span>
          </el-form-item>
          <el-form-item label="工单描述">
            <el-descriptions style="width: 90%" :label-width="80" border>
              <el-descriptions-item :span="3" label="订单号">
                {{ orderInfo.code }}
              </el-descriptions-item>
              <el-descriptions-item label="产品名称">
                {{ orderInfo.productName }}
              </el-descriptions-item>
              <el-descriptions-item label="购买数量">
                {{ orderInfo.productNum }}
              </el-descriptions-item>
              <el-descriptions-item label="购买时长">
                {{ orderInfo.duration }}
                【{{
                  orderInfo.billingType === 1 ? '卡时' :
                      orderInfo.billingType === 2 ? '日' :
                          orderInfo.billingType === 4 ? '周' :
                              orderInfo.billingType === 8 ? '月' :
                                  orderInfo.billingType === 32 ? '年' : '?'
                }}】
              </el-descriptions-item>
              <el-descriptions-item
                  :rowspan="1"
                  label="已选规格"
              >
                <el-descriptions size="small" border>
                  <el-descriptions-item
                      v-for="(item, index) in orderInfo.specJson"
                      :key="index"
                      :label="item.name"
                      :span="orderInfo.specJson.length"
                      label-align="right"
                      align="left"
                  >
                    {{ item.value[item.valueIndex] }}
                  </el-descriptions-item>
                </el-descriptions>
              </el-descriptions-item>
            </el-descriptions>
          </el-form-item>
          <el-form-item v-if="dataInfo.isQuotation" label="订单报价" prop="quotation">
            <el-input-number
                style="width: 270px;"
                controls-position="right"
                v-model="dataInfo.quotation"
                :min="0"
                :precision="2"
                placeholder="请输入"
            >
              <template #prefix>
                <span>￥</span>
              </template>
            </el-input-number>
          </el-form-item>
          <el-form-item v-if="dataInfo.isThird" label="是否上传合同" prop="isContract">
            <div class="flex flex-col">
              <el-radio-group v-model="dataInfo.isContract">
                <el-radio :label="0">否</el-radio>
                <el-radio :label="1">是</el-radio>
              </el-radio-group>
              <el-form v-if="dataInfo.isContract" class="flex flex-col gap-2" label-width="70">
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="合同名称">
                      <el-input v-model="dataInfo.contractName" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="合同编号">
                      <el-input v-model="dataInfo.contractCode" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="甲方名称">
                      <el-input v-model="dataInfo.contractPartyA" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="乙方名称">
                      <el-input v-model="dataInfo.contractPartyB" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="签署时间">
                      <el-date-picker
                          v-model="dataInfo.contractStartTime"
                          style="width: 200px"
                          placeholder="请选择"
                          type="date"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="终止时间">
                      <el-date-picker
                          v-model="dataInfo.contractEndTime"
                          style="width: 200px"
                          placeholder="请选择"
                          type="date"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="总金额">
                      <el-input-number
                          v-model="dataInfo.contractAmount"
                          style="width: 200px;"
                          controls-position="right"
                          :min="0"
                          :precision="2"
                          placeholder="请输入"
                      >
                        <template #suffix>
                          <span>￥</span>
                        </template>
                      </el-input-number>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item label="合同附件">
                  <el-upload
                      style="width: 200px;"
                      :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
                      :on-error="uploadError"
                      :on-success="uploadContractSuccess"
                      :show-file-list="true"
                      :limit="1"
                  >
                    <el-button type="primary">上 传</el-button>
                  </el-upload>
                </el-form-item>
              </el-form>
            </div>
          </el-form-item>
          <el-form-item v-if="dataInfo.isThird" label="是否上传发票" prop="isInvoice">
            <div class="flex flex-col">
              <el-radio-group v-model="dataInfo.isInvoice">
                <el-radio :label="0">否</el-radio>
                <el-radio :label="1">是</el-radio>
              </el-radio-group>
              <el-form v-if="dataInfo.isInvoice" class="flex flex-col gap-2 right-0" label-width="70">
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="发票名称">
                      <el-input v-model="dataInfo.invoiceName" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="发票编号">
                      <el-input v-model="dataInfo.invoiceCode" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="买方名称">
                      <el-input v-model="dataInfo.invoicePartyA" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="卖方名称">
                      <el-input v-model="dataInfo.invoicePartyB" style="width: 200px" placeholder="请输入"/>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="税率">
                      <el-input-number
                          v-model="dataInfo.invoiceTaxRate"
                          style="width: 200px;"
                          controls-position="right"
                          :min="0"
                          :max="100"
                          :precision="2"
                          placeholder="请输入">
                        <template #suffix>
                          <span>%</span>
                        </template>
                      </el-input-number>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="税额">
                      <el-input-number
                          v-model="dataInfo.invoiceTaxAmount"
                          style="width: 200px;"
                          controls-position="right"
                          :min="0"
                          :max="100"
                          :precision="2"
                          placeholder="请输入"
                      >
                        <template #suffix>
                          <span>￥</span>
                        </template>
                      </el-input-number>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="总金额">
                      <el-input-number
                          v-model="dataInfo.invoiceAmount"
                          style="width: 200px;"
                          controls-position="right"
                          :min="0"
                          :precision="2"
                          placeholder="请输入"
                      >
                        <template #suffix>
                          <span>￥</span>
                        </template>
                      </el-input-number>
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-row :gutter="50">
                  <el-col :span="12">
                    <el-form-item label="开具时间">
                      <el-date-picker
                          v-model="dataInfo.invoiceStartTime"
                          style="width: 200px"
                          placeholder="请选择"
                          type="date"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item label="发票附件">
                  <el-upload
                      style="width: 200px;"
                      :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
                      :on-error="uploadError"
                      :on-success="uploadInvoiceSuccess"
                      :show-file-list="true"
                      :limit="1"
                  >
                    <el-button type="primary">上 传</el-button>
                  </el-upload>
                </el-form-item>
              </el-form>
            </div>
          </el-form-item>
        </el-form>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="closeDialog">取消</el-button>
          <el-button
              type="primary"
              @click="submitDialog"
          >确定
          </el-button>
        </div>
      </el-drawer>
      <!-- 报价弹窗 -->
      <el-dialog
          v-model="quoteDialogVisible"
          title="报价"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="true"
          width="30%"
      >
        <el-form label-width="100px">
          <el-form-item label="报价金额">
            <el-input
                v-model.number="quotation"
                placeholder="请输入报价金额"
                type="number"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="quoteDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitQuote">确定</el-button>
          </span>
        </template>
      </el-dialog>
      <!-- 上传发票抽屉 -->
      <el-drawer
          v-model="uploadInvoiceDrawerVisible"
          title="上传发票"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="true"
          size="60%"
      >
        <el-form :rules="rules" label-width="100px">
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="invoiceRemark" placeholder="请输入备注信息"/>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="uploadInvoiceDrawerVisible = false">取消</el-button>
            <el-button type="primary" @click="submitUploadInvoice">确定</el-button>
          </span>
        </template>
      </el-drawer>
      <!-- 上传合同抽屉 -->
      <el-drawer
          v-model="uploadContractDrawerVisible"
          title="上传合同"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="true"
          size="60%"
      >
        <el-form :rules="rules" label-width="100px">
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="invoiceRemark" placeholder="请输入备注信息"/>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="uploadContractDrawerVisible = false">取消</el-button>
            <el-button type="primary" @click="submitUploadContract">确定</el-button>
          </span>
        </template>
      </el-drawer>

      <!--    标准产品/定制产品处理对话框-->
      <el-drawer
          v-model="productDialog"
          size="520px"
          :close-on-click-modal="true"
          :close-on-press-escape="false"
          :show-close="false"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
          <span class="text-lg"
                style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ ticketInfo.status == 1 ? '处理工单' : '查看详情' }}</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="productDialog=false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <div v-if="resOrderList.length==0">
          <h3>该工单未生成租户资源相关信息</h3>
        </div>
        <!--      <div v-if="resOrderList.length!=0">-->
        <!--      </div>-->
        <div style="margin-top: 10px" v-for="(item,index) in resOrderList">
          <el-descriptions class="no-select" :column="1" label-align="left" label-width="100px">
            <el-descriptions-item :span="2" label="产品名称" v-if="item.description!=''">
              {{ JSON.parse(item.description).productName }}
            </el-descriptions-item>
            <el-descriptions-item :span="2" label="所分配资源">
              <span>{{ GetResName(item.resourceID, item.type) }}</span>
              <span v-if="GetResName(item.resourceID, item.type)!='未分配资源'|| ticketInfo.status==2 "
                    style="color:rgba(188, 191, 196, 1); margin-left: 8px"
              >分配
              </span>
              <span v-else
                    @click="callBareMethod(item)"
                    style="color:rgba(22, 93, 255, 1); margin-left: 8px"
              >分配
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="订单号" v-if="item.description!=''">
              {{ JSON.parse(item.description).code }}
            </el-descriptions-item>
            <el-descriptions-item label="所属租户">
              {{ GetRenterName(item.renterID) }}
            </el-descriptions-item>
            <el-descriptions-item label="产品类型">
              {{ productType[item.type] }}
            </el-descriptions-item>
            <el-descriptions-item label="购买时长"
                                  v-if="item.description!=''">
              {{ JSON.parse(item.description).duration }}
              【{{
                JSON.parse(item.description).billingType === 1 ? '卡时' :
                    JSON.parse(item.description).billingType === 2 ? '日' :
                        JSON.parse(item.description).billingType === 4 ? '周' :
                            JSON.parse(item.description).billingType === 8 ? '月' :
                                JSON.parse(item.description).billingType === 32 ? '年' : '?'
              }}】
            </el-descriptions-item>
            <el-descriptions-item
                v-if="item.description!=''"
                v-for="(temp, index1) in JSON.parse(JSON.parse(item.description).specJson)"
                :key="index1"
                :label="temp.name"
                label-align="left"
                label-width="100px"
                align="left"
            >
              {{ temp.value[temp.valueIndex != undefined ? temp.valueIndex : 0] }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
        <el-form-item style="margin-top: 5px" v-if="ticketInfo.isQuotation&&resOrderList.length!=0" label="订单报价"
                      prop="quotation">
          <el-input-number
              style="width: 270px;"
              controls-position="right"
              v-model="ticketInfo.quotation"
              :min="0"
              :disabled="ticketInfo.status==2"
              :precision="2"
              placeholder="请输入"
          >
            <template #prefix>
              <span>￥</span>
            </template>
          </el-input-number>
        </el-form-item>
        <el-form-item v-if="ticketInfo.isThird&&resOrderList.length!=0" label="是否上传合同" prop="isContract">
          <div class="flex flex-col">
            <el-radio-group :disabled="ticketInfo.status==2" v-model="ticketInfo.isContract">
              <el-radio :label="0">否</el-radio>
              <el-radio :label="1">是</el-radio>
            </el-radio-group>
            <div v-if="ticketInfo.isContract">
              <el-descriptions v-if="ticketInfo.contractJson!=''" :column="1" :title="'合同信息'" border>
                <el-descriptions-item label="文件名">{{ JSON.parse(ticketInfo.contractJson).fileName }}
                </el-descriptions-item>
                <el-descriptions-item label="备注">{{
                    JSON.parse(ticketInfo.contractJson).remark
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="上传时间">
                  {{ FormatDateTime(JSON.parse(ticketInfo.contractJson).uploadAt * 1000) }}
                </el-descriptions-item>
              </el-descriptions>
              <el-button v-if="ticketInfo.status==1" @click="showUploadContractDrawer">上传合同</el-button>
            </div>
          </div>
        </el-form-item>
        <el-form-item v-if="ticketInfo.isThird&&resOrderList.length!=0" label="是否上传发票" prop="isInvoice">
          <div class="flex flex-col">
            <el-radio-group :disabled="ticketInfo.status==2" v-model="ticketInfo.isInvoice">
              <el-radio :label="0">否</el-radio>
              <el-radio :label="1">是</el-radio>
            </el-radio-group>
            <div v-if="ticketInfo.isInvoice">
              <el-descriptions v-if="ticketInfo.invoiceJson!=''" :column="1" :title="'发票信息'" border>
                <el-descriptions-item label="文件名">{{ JSON.parse(ticketInfo.invoiceJson).fileName }}
                </el-descriptions-item>
                <el-descriptions-item label="备注">{{
                    JSON.parse(ticketInfo.invoiceJson).remark
                  }}
                </el-descriptions-item>
                <el-descriptions-item label="上传时间">
                  {{ FormatDateTime(JSON.parse(ticketInfo.invoiceJson).uploadAt * 1000) }}
                </el-descriptions-item>
              </el-descriptions>
              <el-button v-if="ticketInfo.status==1" @click="showUploadInvoiceDrawer">上传发票</el-button>
            </div>
          </div>
        </el-form-item>
        <div class="flex" style="justify-content: right;margin-top: 24px">
          <el-button type="info" @click="productDialog=false">{{ ticketInfo.status == 1 ? '取消' : '关闭 ' }}
          </el-button>
          <el-button
              v-if="ticketInfo.status==1"
              type="primary"
              @click="HandleFinish"
          >处理完成
          </el-button>
        </div>
      </el-drawer>
      <!--    分配资源对话框-->
      <el-dialog
          v-model="assignVisible"
          title="分配"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="true"
          width="30%"
      >
        <el-form label-width="100px" :model="assignData" :rules="rules" ref="assignData">
          <el-form-item label="资源">
            <el-select v-model="assignData.resourceID" placeholder="请选择资源">
              <el-option v-for="(item,index) in noAssignList"
                         :label="item.name"
                         :value="item.resourceID"/>
            </el-select>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="assignVisible = false">取消</el-button>
            <el-button type="primary" @click="ConfirmAssign">确定</el-button>
          </span>
        </template>
      </el-dialog>
      <baremetal @dialogSubmit="handleBareSubmitted" :renter-i-d="baremetalRenterID" ref="bareRef"/>
      <HostDialog @dialogSubmit="handleBareSubmitted" :hostData="hostData" ref="hostRef"/>
      <ConDialog @dialogSubmit="handleBareSubmitted" :container-data="conData" ref="conRef"/>
      <!--    删除弹框-->
      <el-dialog
          v-model="deleteVisible"
          :show-close="false"
          width="520px"
          @close="deleteVisible=false"
      >
        <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg"
              style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);display: flex;align-items: center;">
          <el-icon :style="{ color: 'rgba(255, 125, 0, 1)', fontSize: '1.25em',marginRight: '8px' }">
            <WarningFilled/>
          </el-icon>删除</span>
          <el-icon
              class="cursor-pointer hover:text-gray-500 transition-colors"
              @click="deleteVisible = false"
              style="width: 16px; height: 16px; font-size: 16px"
          >
            <Close/>
          </el-icon>
        </div>
        <span>该操作将永久删除该工单，是否继续?</span>
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
  addSysTicket,
  deleteSysTicket,
  getAllRes,
  getNoAssignRes,
  getSysTicketList,
  postQuotation,
  updateSysTicket,
  uploadContract,
  uploadInvoice
} from "@/api/yunguan/system/ticket";
import {getAllRenterList, getAllUserList, getRenterResByTicketId, updateRenterRes} from "@/api/yunguan/run/renter";
import {getBaseUrl} from "@/utils/format";
import {ElMessage} from "element-plus";
import _ from "lodash";
import Baremetal from './component/baremetalDialog.vue';
import HostDialog from './component/hostDialog.vue';
import ConDialog from './component/containerDialog.vue';

export default {
  components: {
    Baremetal,
    HostDialog,
    ConDialog,
  },
  data() {
    return {
      deleteId: '',
      deleteVisible: false,
      assignVisible: false,
      assignData: {
        id: '',
        resourceID: '',
        sshUser: '',
        sshPwd: '',
        sshHost: '',
        privateIp: '',
      },
      activeName: '标准',
      searchInfo: {
        page: 1,
        pageSize: 10,
        status: 0,
        ticketType: 1,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      dataInfo: {
        id: 0,
        type: '',
        ticketType: '',
        renterID: '',
        description: '',
        status: null,
        handlePerson: '',
        isQuotation: 0, // 新增字段，是否需要报价
      },
      orderInfo: null,
      handleStatus: ["", "未处理", "已处理"],
      rules: {
        quotation: [{required: true, message: '请输入报价', trigger: 'blur'}],
        isContract: [{required: true, message: '请选择是否上传合同', trigger: 'change'}],
        isInvoice: [{required: true, message: '请选择是否上传发票', trigger: 'change'}],
        resourceID: [{required: true, message: '请选择资源', trigger: 'blur'}],
        file: [{required: true, message: '请上传文件', trigger: 'blur'}],
      },
      ticketList: [],
      userList: [],
      renterList: [],
      ticketType: ["资源故障", "资源分配", "资源扩充"],
      quoteDialogVisible: false, // 报价弹窗显示状态
      quotation: 0, // 报价金额
      currentQuoteTicketId: null, // 当前报价工单的 ID
      uploadInvoiceDrawerVisible: false, // 上传发票抽屉显示状态
      uploadContractDrawerVisible: false,
      currentUploadInvoiceId: null, // 当前上传发票工单的 ID
      file: null,
      // contract_file:null,
      invoiceRemark: '', // 发票备注信息
      productDialog: false,
      resOrderList: [],
      productType: ["", "裸金属", "云主机", "容器"],
      allResBareList: [],
      allResHostList: [],
      allResConnList: [],
      noAssignBareList: [],
      noAssignHostList: [],
      noAssignConnList: [],
      noAssignList: [],
      modRenterRes: {
        id: '',
        resourceId: '',
      },
      ticketId: '',
      ticketInfo: {},
      baremetalRenterID: '',
      hostData: {
        cpu: 1,
        mem: 1,
        renterID: '',
        sysDisk: 10,
        dataDisk: 0,
        timeNum: 1,
        timeType: 1,
      },
      conData: {
        cpu: 1,
        mem: 1,
        renterID: '',
        sysDisk: 10,
        dataDisk: 0,
        timeNum: 1,
        timeType: 1,
      },
    }
  },
  created() {
    this.GetTicketList();
    this.GetUserList();
    this.GetRenterList();
    // this.GetAllResList()
    // this.GetNoAssignList()
  },
  watch: {
    resOrderList(newVal) {
      console.log("对话框数据已更新:", newVal);
    },
  },
  methods: {
    confirmDelete() {
      switch (this.deleteType) {
        case 'delete':
          this.handleDelete(this.deleteId)
          break;
      }
      this.deleteVisible = false;
    },
    HandleDeleteDialog(val, type) {
      this.deleteId = val.ID;
      this.deleteType=type;
      this.deleteVisible = true;
    },
    handleSortChange({column, prop, order}) {
      switch (order) {
        case "ascending":
          this.searchInfo.timeDesc = "asc";
          break;
        case "descending":
          this.searchInfo.timeDesc = "desc";
          break;
        default:
          this.searchInfo.timeDesc = "";
          break;
      }
      this.GetTicketList();
    },
    //解析json
    getSelectedValue(specJson, tag) {
      let item = specJson.find(item => item.tag === tag); // 找到指定 tag 的对象
      item = item ? item.value[item.valueIndex != undefined ? item.valueIndex : 0] : null; // 返回被选中的值
      if (item) {
        item = item.replace(/[^0-9]/g, '');
      }
      return parseInt(item);
    },
    //接收裸金属对话框数据
    handleBareSubmitted(value, type) {
      this.assignData.resourceID = value;
      this.GetNoAssignList(type);
      this.GetAllResList();
    },
    callBareMethod(item) {
      this.assignData.id = item.ID;
      switch (item.type) {
        case 1:
          this.baremetalRenterID = item.renterID;
          this.$refs.bareRef.OpenCreateDrawer();
          break;
        case 2:
          const temp = JSON.parse(JSON.parse(item.description).specJson);
          this.hostData.cpu = this.getSelectedValue(temp, "cpu");
          this.hostData.mem = this.getSelectedValue(temp, "memory");
          this.hostData.dataDisk = this.getSelectedValue(temp, "dataDisk");
          this.hostData.sysDisk = this.getSelectedValue(temp, "sysDisk");
          this.hostData.timeNum = JSON.parse(item.description).duration
          this.hostData.timeType = JSON.parse(item.description).billingType
          this.hostData.renterID = item.renterID;
          this.$refs.hostRef.OpenCreateDrawer();
          break;
        case 3:
          const tempData = JSON.parse(JSON.parse(item.description).specJson);
          this.conData.cpu = this.getSelectedValue(tempData, "cpu");
          this.conData.mem = this.getSelectedValue(tempData, "memory");
          this.conData.sysDisk = this.getSelectedValue(tempData, "sysDisk");
          this.conData.dataDisk = this.getSelectedValue(tempData, "dataDisk");
          this.conData.renterID = item.renterID;
          this.conData.timeNum = JSON.parse(item.description).duration
          this.conData.timeType = JSON.parse(item.description).billingType
          this.$refs.conRef.handleAdd();
          break;
      }
    },
    //处理完成
    HandleFinish() {
      for (let i = 0; i < this.resOrderList.length; i++) {
        if (this.resOrderList[i].resourceID == '') {
          this.$message({
            type: 'warning',
            message: '资源还未分配'
          })
          return;
        }
      }
      if (this.ticketInfo.isContract == 1) {
        if (this.ticketInfo.contractJson == '') {
          this.$message({
            type: 'warning',
            message: '请上传合同'
          })
          return;
        }
      }
      if (this.ticketInfo.isInvoice == 1) {
        if (this.ticketInfo.invoiceJson == '') {
          this.$message({
            type: 'warning',
            message: '请上传发票'
          })
          return;
        }
      }
      let data = {};
      data.id = this.ticketInfo.ID;
      data.status = 2;
      data.quotation = this.ticketInfo.quotation;
      updateSysTicket(data).then((res) => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '处理完成'
          })
          this.GetTicketList()
          this.productDialog = false;
        }
      })
    },
    //分配确定
    ConfirmAssign() {
      for (let i = 0; i < this.resOrderList.length; i++) {
        if (this.resOrderList[i].ID == this.assignData.id) {
          this.resOrderList[i].resourceID = this.assignData.resourceID;
          break;
        }
      }
      const temp = this.noAssignList.find(item => item.resourceID == this.assignData.resourceID)
      if (temp) {
        this.assignData.privateIp = temp.privateIp;
        this.assignData.sshHost = temp.privateIp;
        this.assignData.sshUser = temp.username;
        this.assignData.sshPwd = temp.password;
      }
      updateRenterRes(this.assignData).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '处理完成'
          })
          this.GetRenterResByOrderId();
        }
      })
    },
    //获取资源名称
    GetResName(id, type) {
      let res;
      switch (type) {
        case 1:
          res = this.allResBareList.find(t => t.resourceID === id)
          break;
        case 2:
          res = this.allResHostList.find(t => t.resourceID === id)
          break;
        case 3:
          res = this.allResConnList.find(t => t.resourceID === id)
          break;
      }
      return res ? res.name : '未分配资源';
    },
    //获取所有资源列表
    GetAllResList() {
      getAllRes({type: 1}).then(res => {
        if (res.code == 0) {
          this.allResBareList = res.data.list;
        }
      })
      getAllRes({type: 2}).then(res => {
        if (res.code == 0) {
          this.allResHostList = res.data.list;
        }
      })
      getAllRes({type: 3}).then(res => {
        if (res.code == 0) {
          this.allResConnList = res.data.list;
        }
      })
    },
    //获取未被分配资源列表
    GetNoAssignList(type) {
      switch (type) {
        case 1:
          getNoAssignRes({type: 1}).then(res => {
            if (res.code == 0) {
              this.noAssignBareList = res.data.list;
              this.noAssignList = _.cloneDeep(this.noAssignBareList);
              this.ConfirmAssign()
            }
          })
          break;
        case 2:
          getNoAssignRes({type: 2}).then(res => {
            if (res.code == 0) {
              this.noAssignHostList = res.data.list;
              this.noAssignList = _.cloneDeep(this.noAssignHostList);
              this.ConfirmAssign()
            }
          })
          break;
        case 3:
          getNoAssignRes({type: 3}).then(res => {
            if (res.code == 0) {
              this.noAssignConnList = res.data.list;
              this.noAssignList = _.cloneDeep(this.noAssignConnList);
              this.ConfirmAssign()
            }
          })
          break;
      }
    },
    //打开处理对话框
    OpenHandleDrawer(val) {
      this.ticketId = val.ID;
      this.GetRenterResByOrderId();
      this.ticketInfo = _.cloneDeep(val);
      this.productDialog = true;
      this.ticketInfo.quotation = 0;
      if (this.ticketInfo.contractJson != '') {
        this.ticketInfo.isContract = 1;
      } else this.ticketInfo.isContract = 0;
      if (this.ticketInfo.invoiceJson != '') {
        this.ticketInfo.isInvoice = 1;
      } else this.ticketInfo.isInvoice = 0;
    },
    //通过订单获取资源列表
    GetRenterResByOrderId() {
      getRenterResByTicketId({ticketId: this.ticketId}).then(res => {
        if (res.code == 0) {
          this.resOrderList = res.data;
        }
      })
    },
    //tab切换
    handleTabChange() {
      switch (this.activeName) {
        case "标准":
          this.searchInfo.ticketType = 1;
          break;
        case "定制":
          this.searchInfo.ticketType = 2;
          break;
        case "投诉":
          this.searchInfo.ticketType = 3;
          break;
        case "合同":
          this.searchInfo.ticketType = 4;
          break;
      }
      this.GetTicketList();
    },
    getFileData() {
      const inputFile = this.$refs.file.files[0];
      this.file = inputFile
    },
    getFilecontractData() {
      const inputFile = this.$refs.file.files[0];
      this.contract_file = inputFile
    },
    FormatDateTime(dateString) {
      if (!dateString) return ''
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
        // 如果日期是不合理的，返回空字符串
        return '';
      }
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hour = String(date.getHours()).padStart(2, '0');
      const minute = String(date.getMinutes()).padStart(2, '0');
      const second = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
    },
    handleDelete(val) {
      deleteSysTicket(val).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '删除成功',
            type: 'success'
          });
          this.GetTicketList(this.searchInfo);
        }
      });
    },
    handleEdit(val) {
      // this.dataInfo.id = val.ID;
      // this.dataInfo.name = val.name;
      // this.dataInfo.handlePerson = val.handlePerson;
      // this.dataInfo.description = val.description;
      // this.dataInfo.ticketType = val.ticketType;
      // this.dataInfo.isQuotation = val.isQuotation;
      // this.dataInfo.type = 1;

      // this.dataInfo = val
      // 深拷贝，保证dataInfo修改后与绑定值一致
      this.dataInfo = JSON.parse(JSON.stringify(val))
      this.orderInfo = JSON.parse(val.description)
      this.orderInfo.specJson = JSON.parse(this.orderInfo.specJson)
      this.orderInfo.storage = JSON.parse(this.orderInfo.storage)
      this.orderInfo.specJson.push({name: '存储类型', value: [this.orderInfo.storage[0].name], valueIndex: 0})
      this.orderInfo.specJson.push({name: '操作系统', value: [this.orderInfo.os], valueIndex: 0})

      this.IsCreated = false;
      this.openDialog = true;
    },
    handleDetail(val) {
      this.dataInfo = val
      this.IsCreated = false;
      this.openDialog = true;
    },
    handleSearch() {
      this.GetTicketList();
    },
    onReset() {
      this.searchInfo.status = 0;
      this.GetTicketList();
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val;
      this.GetTicketList();
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val;
      this.GetTicketList();
    },
    addHardwareOpenDialog() {
      this.openDialog = true;
      this.IsCreated = true;
      this.dataInfo = {};
      this.dataInfo.status = 1;
      this.dataInfo.type = 1;
      this.dataInfo.isQuotation = 0;
    },
    closeDialog() {
      this.openDialog = false
      this.dataInfo = {}
    },
    UpdateTicket() {
      this.dataInfo.status = 3
      this.dataInfo.contractJson = {
        name: this.dataInfo.contractName,
        code: this.dataInfo.contractCode,
        signedTime: this.dataInfo.contractStartTime,
        startTime: this.dataInfo.contractStartTime,
        endTime: this.dataInfo.contractEndTime,
        file: this.dataInfo.contractFile,
        amount: this.dataInfo.contractAmount,
        customerName: this.dataInfo.contractPartyA,
        supplierName: this.dataInfo.contractPartyB,
      }
      this.dataInfo.invoiceJson = {
        name: this.dataInfo.invoiceName,
        code: this.dataInfo.invoiceCode,
        signedTime: this.dataInfo.invoiceStartTime,
        startTime: this.dataInfo.invoiceStartTime,
        // endTime: this.dataInfo.invoiceEndTime,
        file: this.dataInfo.invoiceFile,
        amount: this.dataInfo.invoiceAmount,
        customerName: this.dataInfo.invoicePartyA,
        supplierName: this.dataInfo.invoicePartyB,
        taxRate: this.dataInfo.invoiceTaxRate,
        taxAmount: this.dataInfo.invoiceTaxAmount,
      }
      updateSysTicket(this.dataInfo).then(response => {
        if (response.code === 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false;
          this.GetTicketList();
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      });
    },
    // 新增分区
    AddHardware() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addSysTicket(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.openDialog = false;
              this.GetTicketList();
            } else {
              this.$message({
                message: '添加失败',
                type: 'error'
              });
            }
          });
        }
      });
    },
    // 提交
    submitDialog() {
      this.$refs.dataInfo.validate(async valid => {
        if (valid) {
          if (this.IsCreated) {
            await this.AddHardware()
          } else {
            await this.UpdateTicket()
          }
          this.GetTicketList()
        } else {
          this.$message({
            message: '请填写完整表单',
            type: 'error'
          })
        }
      })
    },
    // 获取工单列表
    GetTicketList() {
      getSysTicketList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.ticketList = response.data.list;
          this.total = response.data.total;
          // if (this.ticketInfo.ID != undefined) {
          //   this.ticketInfo = this.ticketList.find((item) => item.ID == this.ticketInfo.ID);
          //   if (this.ticketInfo.contractJson != '') {
          //     this.ticketInfo.isContract = 1;
          //   }
          //   if (this.ticketInfo.invoiceJson != '') {
          //     this.ticketInfo.isInvoice = 1;
          //   }
          // }
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      });
    },
    //获取分区名称
    GetUserName(id) {
      if (id === 0) {
        return '系统'
      }
      const tenant = this.userList.find(t => t.ID === id);
      return tenant ? tenant.userName : '-';
    },
    //获取分区名称
    GetRenterName(id) {
      const tenant = this.renterList.find(t => t.ID === id);
      return tenant ? tenant.userName : null;
    },
    // 获取用户列表
    GetUserList() {
      // todo 后期一定需要优化，获取全部用户并遍历处理不是规范写法
      getAllUserList().then(response => {
        if (response.code == 0) {
          this.userList = response.data.list;
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      });
    },
    // 获取用户列表
    GetRenterList() {
      getAllRenterList().then(response => {
        if (response.code == 0) {
          this.renterList = response.data.list;
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      });
    },
    // 上传合同的方法，需要根据实际情况实现
    uploadContract(row) {
      console.log('上传合同', row);
    },
    // 显示报价弹窗
    showQuoteDialog(row) {
      this.quoteDialogVisible = true;
      this.currentQuoteTicketId = row.ID;
      this.quotation = 0; // 清空报价金额
    },
    // 提交报价
    submitQuote() {
      if (this.quotation <= 0) {
        this.$message({
          message: '请输入有效的报价金额',
          type: 'warning'
        });
        return;
      }
      const id = this.currentQuoteTicketId;
      const quotation = this.quotation;
      postQuotation({id, quotation})
          .then(response => {
            if (response.code === 0) {
              this.$message({
                message: '报价提交成功',
                type: 'success'
              });
              this.quoteDialogVisible = false;
              this.GetTicketList();
            } else {
              this.$message({
                message: '报价提交失败：' + response.message,
                type: 'error'
              });
            }
          })
          .catch(error => {
            this.$message({
              message: '报价提交出错：' + error.message,
              type: 'error'
            });
          });
    },
    // 显示上传发票抽屉
    showUploadInvoiceDrawer() {
      this.uploadInvoiceDrawerVisible = true;
      this.currentUploadInvoiceId = this.ticketInfo.ID;
      this.selectedFile = null;
      this.invoiceRemark = '';
    },

    showUploadContractDrawer() {
      this.uploadContractDrawerVisible = true;
      this.currentUploadInvoiceId = this.ticketInfo.ID;
      this.selectedFile = null;
      this.invoiceRemark = '';
    },
    // 处理文件选择

    handlefileUploadSuccess(response) {
      if (response.code === 0) {
        this.selectedFile = response.data.file.url;  // 保存文件的URL
      }
    },
    // 提交上传发票
    submitUploadInvoice() {
      const formData = new FormData();
      if (this.file == null) {
        this.$message({
          type: 'warning',
          message: '请上传文件',
        })
        return;
      }
      formData.append('id', this.currentUploadInvoiceId);
      formData.append('file', this.file);
      formData.append('remark', this.invoiceRemark);
      this.file = null;

      uploadInvoice(formData)
          .then(response => {
            if (response.code === 0) {
              this.$message({
                message: '发票上传成功',
                type: 'success'
              });
              this.uploadInvoiceDrawerVisible = false;
              this.GetTicketList();
            } else {
              this.$message({
                message: '发票上传失败：' + response.message,
                type: 'error'
              });
            }
          })
          .catch(error => {
            this.$message({
              message: '发票上传出错：' + error.message,
              type: 'error'
            });
          });
    },

    submitUploadContract() {
      const formData = new FormData();
      if (this.file == null) {
        this.$message({
          type: 'warning',
          message: '请上传文件',
        })
        return;
      }
      formData.append('id', this.currentUploadInvoiceId);
      formData.append('contract_file', this.file);
      formData.append('remark', this.invoiceRemark);
      this.file = null;
      uploadContract(formData)
          .then(response => {
            if (response.code === 0) {
              this.$message({
                message: '合同上传成功',
                type: 'success'
              });
              this.uploadContractDrawerVisible = false;
              this.GetTicketList();
            } else {
              this.$message({
                message: '合同上传失败：' + response.message,
                type: 'error'
              });
            }
          })
          .catch(error => {
            this.$message({
              message: '合同上传出错：' + error.message,
              type: 'error'
            });
          });
    },
    handleTagType(status) {
      switch (status) {
        case 1:
          return 'danger'
        case 2:
          return 'success'
        default:
          return 'info'
      }
    },
    handleQuotationChange() {
      this.dataInfo.quotation = this.dataInfo.quotation.toFixed(2)
    },
    uploadContractSuccess(response) {
      this.dataInfo.contractFile = response.data?.file?.url
    },
    uploadInvoiceSuccess(response) {
      this.dataInfo.invoiceFile = response.data?.file?.url
    },
    uploadError() {
      ElMessage({
        type: 'error',
        message: '附件上传失败'
      })
    },
    getBaseUrl, getUploadBaseUrl() {
      return getBaseUrl();
    },
  }
};
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
