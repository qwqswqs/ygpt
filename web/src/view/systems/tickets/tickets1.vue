<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
        >
          <el-form-item>
            <el-select v-model="searchInfo.ticketType" style="width: 200px" placeholder="请选择工单类别">
              <el-option
                  v-for="item in [{ label: '全部', value: 0}, { label: '资源故障', value: 1}, { label: '资源分配', value: 2}]"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select v-model="searchInfo.status" style="width: 200px" placeholder="请选择状态">
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
            >
              查询
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

      <!-- 表格展示 -->
      <el-table
          :data="ticketList"
          row-key="ID"
      >
        <el-table-column label="工单创建时间" min-width="160px">
          <template #default="scope">
            <span>{{ FormatDateTime(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="工单类别">
          <template #default="scope">
            <span>{{ ticketType[scope.row.ticketType - 1] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" >
          <template #default="scope">
            <el-tag :type="handleTagType(scope.row.status)">
              {{ handleStatus[scope.row.status - 1] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <span class="text-overflow">{{ GetRenterName(scope.row.renterID) }}</span>
          </template>
        </el-table-column>
<!--        <el-table-column label="描述信息" min-width="200">-->
<!--          <template #default="scope">-->
<!--            <span class="text-overflow">{{ scope.row.description }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--        <el-table-column label="分配人">-->
<!--          <template #default="scope">-->
<!--            <span>{{ GetUserName(scope.row.createPerson) }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column label="处理人">
          <template #default="scope">
            <span class="text-overflow">{{ GetUserName(scope.row.handlePerson) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="100px" fixed="right">
          <template #default="scope">
<!--            <el-button-->
<!--                class="text-blue-500"-->
<!--                @click="handleDetail(scope.row)"-->
<!--                link-->
<!--            >-->
<!--              详情-->
<!--            </el-button>-->
            <el-button
                class="text-blue-500"
                :disabled="scope.row.status === 3"
                @click="handleEdit(scope.row)"
                link
            >
              处理
            </el-button>
            <el-button
                class="text-red-500"
                link
                @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
<!--            &lt;!&ndash; 上传发票、合同按钮 &ndash;&gt;-->
<!--            <el-button-->
<!--                type="text"-->
<!--                v-if="scope.row.ticketType === 2"-->
<!--                @click="showUploadInvoiceDrawer(scope.row)"-->
<!--            >-->
<!--              发票-->
<!--            </el-button>-->
<!--            <el-button-->
<!--                type="text"-->
<!--                v-if="scope.row.ticketType === 2"-->
<!--                @click="showUploadContractDrawer(scope.row)"-->
<!--            >-->
<!--              合同-->
<!--            </el-button>-->
<!--            &lt;!&ndash; 报价按钮 &ndash;&gt;-->
<!--            <el-button-->
<!--                type="text"-->
<!--                v-if="scope.row.isQuotation === 1"-->
<!--                @click="showQuoteDialog(scope.row)"-->
<!--            >-->
<!--              报价-->
<!--            </el-button>-->
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
      <el-drawer
          v-model="openDialog"
          style="min-width: 600px;"
          size="60%"
          :close-on-click-modal="false"
          :close-on-press-escape="false"
          :show-close="true"
          :title="IsCreated? '新增工单' : '处理工单'"
      >
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
            <span>{{ ticketType[dataInfo.ticketType-1] }}</span>
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
                  orderInfo.billingType === 1  ? '卡时' :
                  orderInfo.billingType === 2  ? '日' :
                  orderInfo.billingType === 4  ? '周' :
                  orderInfo.billingType === 8  ? '月' :
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
                      <el-input v-model="dataInfo.contractPartyA" style="width: 200px" placeholder="请输入" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="乙方名称">
                      <el-input v-model="dataInfo.contractPartyB" style="width: 200px" placeholder="请输入" />
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
                      <el-input v-model="dataInfo.invoicePartyA" style="width: 200px" placeholder="请输入" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="卖方名称">
                      <el-input v-model="dataInfo.invoicePartyB" style="width: 200px" placeholder="请输入" />
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
        <div class="flex justify-end mt-4">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="submitDialog">提 交</el-button>
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
        <el-form label-width="100px">
          <el-form-item label="文件上传" prop="file">
            <input ref="file" type="file" @change="getFileData"/>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="invoiceRemark" placeholder="请输入备注信息" />
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
        <el-form-item label="文件上传" prop="file">
          <input ref="file" type="file" @change="getFileData"/>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="invoiceRemark" placeholder="请输入备注信息" />
        </el-form-item>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="uploadContractDrawerVisible = false">取消</el-button>
            <el-button type="primary" @click="submitUploadContract">确定</el-button>
          </span>
        </template>
      </el-drawer>
    </div>
  </div>
</template>

<script>
import {
  addSysTicket,
  getSysTicketList,
  updateSysTicket,
  deleteSysTicket,
  postQuotation,
  uploadInvoice,
  uploadContract
} from "@/api/yunguan/system/ticket";
import { getAllRenterList, getAllUserList } from "@/api/yunguan/run/renter";
import { getBaseUrl } from "@/utils/format";
import {ElMessage} from "element-plus";

export default {
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
        status: 0,
        type:0,
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
        // submitPerson: [
        //   { required: true, message: '请选择硬件类型', trigger: 'blur' }
        // ],
        // handlePerson: [
        //   { required: true, message: '请选择硬件类型', trigger: 'blur' }
        // ]
        quotation: [{ required: true, message: '请输入报价', trigger: 'blur' }],
        isContract: [{ required: true, message: '请选择是否上传合同', trigger: 'change' }],
        isInvoice: [{ required: true, message: '请选择是否上传发票', trigger: 'change' }],
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
      file:null,
      // contract_file:null,
      invoiceRemark: '', // 发票备注信息
    };
  },
  created() {
    this.GetTicketList();
    this.GetUserList();
    this.GetRenterList();
  },
  methods: {
    getBaseUrl,
    getUploadBaseUrl() {
      return getBaseUrl();
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
      this.$confirm('此操作将永久删除该分区, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSysTicket(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            });
            this.GetTicketList(this.searchInfo);
          } else {
            this.$message({
              message: '删除失败',
              type: 'error'
            });
          }
        });
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
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
      return tenant? tenant.userName : '-';
    },
    //获取分区名称
    GetRenterName(id) {
      const tenant = this.renterList.find(t => t.ID === id);
      return tenant? tenant.userName : null;
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
      postQuotation({ id, quotation })
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
    showUploadInvoiceDrawer(row) {
      this.uploadInvoiceDrawerVisible = true;
      this.currentUploadInvoiceId = row.ID;
      this.selectedFile = null;
      this.invoiceRemark = '';
    },

    showUploadContractDrawer(row) {
      this.uploadContractDrawerVisible = true;
      this.currentUploadInvoiceId = row.ID;
      this.selectedFile = null;
      this.invoiceRemark = '';
    },
    // 处理文件选择

    handlefileUploadSuccess(response){
      if (response.code === 0) {
        this.selectedFile = response.data.file.url;  // 保存文件的URL
      }
    },
    // 提交上传发票
    submitUploadInvoice() {
      const formData = new FormData();
      formData.append('id', this.currentUploadInvoiceId);
      formData.append('file', this.file);
      formData.append('remark', this.invoiceRemark);

      uploadInvoice(formData)
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

    submitUploadContract() {
      const formData = new FormData();
      formData.append('id', this.currentUploadInvoiceId);
      formData.append('contract_file', this.file);
      formData.append('remark', this.invoiceRemark);

      uploadContract(formData)
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
    handleTagType(status) {
      switch (status) {
        case 1:
          return 'danger'
        case 2:
          return 'danger'
        case 3:
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
    }
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
</style>
