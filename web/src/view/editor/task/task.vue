<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <el-tabs v-model="activeName" @tab-change="handleTabsClick">
      <el-tab-pane :label="'裸金属'" name="baremetal">
        <div>
          <div class="gva-btn-list">
            <el-form ref="searchForm" :inline="true" :model="searchInfo">
              <el-form-item>
                <el-select placeholder="请选择实例" v-model="getTaskReq.privateIp" @change="HandleChangeInstance">
                  <el-option :key="item.privateIp" v-for="(item) in instanceList" :label="'裸金属:' + item.name + ' IP地址:' + item.privateIp"
                             :value="item.privateIp" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button  type="primary" :disabled="getTaskReq.privateIp == ''"
                            @click="addRenterTaskOpenDialog">
                  新增
                </el-button>
              </el-form-item>
            </el-form>
          </div>
          <!-- 表格展示 -->
          <el-table :data="taskList" row-key="ID" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
            <el-table-column label="任务名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.taskName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="框架名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.frameName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务类型">
              <template #default="scope">
                <span>{{ scope.row.mode == "train" ? "训练" : "推理" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="模型文件路径">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.modelFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="数据集文件路径" width="140">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.dataSetFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).message }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="进度">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).progress }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="输出">
              <template #default="scope">
                <!--            <span v-if="scope.row && scope.row.resultMsg">-->
                <!--            {{ parseResultMsg(scope.row.resultMsg).result }}-->
                <!--            </span>-->
                <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;"
                      :title="'结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth'"
                      v-if="parseResultMsg(scope.row.resultMsg).message === '任务结束'">
                    {{ '结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth' }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50 && parseResultMsg(scope.row.resultMsg).status !== 4"
                           @click="TaskManage(scope.row.taskID, 1)">暂停
                </el-button>
                <el-button type="text"
                           v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status === 4"
                           @click="TaskManage(scope.row.taskID, 2)">恢复
                </el-button>
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50"
                           @click="TaskManage(scope.row.taskID, 3)">终止
                </el-button>
                <span type="text"
                      v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status >= 50">已保存退出</span>
                <span type="text"
                      v-else-if="!scope.row && scope.row.resultMsg || !parseResultMsg(scope.row.resultMsg).status">等待响应</span>
                <span type="text" v-else-if="parseResultMsg(scope.row.resultMsg).status === -10">异常退出</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination :current-page="getTaskReq.page" :page-size="getTaskReq.pageSize"
                           :page-sizes="[10, 30, 50, 100]" :total="total" layout="total,prev, pager, next,  sizes, jumper"
                           @current-change="handleCurrentChange" @size-change="handleSizeChange" />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'云主机'" name="virtual">
        <div>
          <div class="gva-btn-list">
            <el-form ref="searchForm" :inline="true" :model="searchInfo">
              <el-form-item>
                <el-select placeholder="请选择实例" v-model="getTaskReq.privateIp" @change="HandleChangeInstance">
                  <el-option :key="item.privateIp" v-for="(item) in instanceList" :label="'云主机名称:' + item.name + ' IP地址:' + item.privateIp"
                             :value="item.privateIp" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button  type="primary" :disabled="getTaskReq.privateIp == ''"
                            @click="addRenterTaskOpenDialog">
                  新增
                </el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 表格展示 -->
          <el-table :data="taskList"  row-key="ID" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
            <el-table-column label="任务名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.taskName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="框架名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.frameName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务类型">
              <template #default="scope">
                <span>{{ scope.row.mode == "train" ? "训练" : "推理" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="模型文件路径">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.modelFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="数据集文件路径" width="140">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.dataSetFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).message }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="进度">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).progress }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="输出">
              <template #default="scope">
                <!--            <span v-if="scope.row && scope.row.resultMsg">-->
                <!--            {{ parseResultMsg(scope.row.resultMsg).result }}-->
                <!--            </span>-->
                <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;"
                      :title="'结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth'"
                      v-if="parseResultMsg(scope.row.resultMsg).message === '任务结束'">
                    {{ '结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth' }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50 && parseResultMsg(scope.row.resultMsg).status !== 4"
                           @click="TaskManage(scope.row.taskID, 1)">暂停
                </el-button>
                <el-button type="text"
                           v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status === 4"
                           @click="TaskManage(scope.row.taskID, 2)">恢复
                </el-button>
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50"
                           @click="TaskManage(scope.row.taskID, 3)">终止
                </el-button>
                <span type="text"
                      v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status >= 50">已保存退出</span>
                <span type="text"
                      v-else-if="!scope.row && scope.row.resultMsg || !parseResultMsg(scope.row.resultMsg).status">等待响应</span>
                <span type="text" v-else-if="parseResultMsg(scope.row.resultMsg).status === -10">异常退出</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination :current-page="getTaskReq.page" :page-size="getTaskReq.pageSize"
                           :page-sizes="[10, 30, 50, 100]" :total="total" layout="total,prev, pager, next,  sizes, jumper"
                           @current-change="handleCurrentChange" @size-change="handleSizeChange" />
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="'容器'" name="container">
        <div >
          <div class="gva-btn-list">
            <el-form ref="searchForm" :inline="true" :model="searchInfo">
              <el-form-item>
                <el-select placeholder="请选择实例" v-model="getTaskReq.privateIp" @change="HandleChangeInstance">
                  <el-option :key="item.privateIp" v-for="(item) in instanceList" :label="'容器:' + item.name + ' IP地址:' + item.privateIp"
                             :value="item.privateIp" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button  type="primary" :disabled="getTaskReq.privateIp == ''"
                            @click="addRenterTaskOpenDialog">
                  新增
                </el-button>
              </el-form-item>
            </el-form>
          </div>
          <!-- 表格展示 -->
          <el-table :data="taskList" row-key="ID" :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
            <el-table-column label="任务名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.taskName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="框架名称">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.frameName
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="任务类型">
              <template #default="scope">
                <span>{{ scope.row.mode == "train" ? "训练" : "推理" }}</span>
              </template>
            </el-table-column>
            <el-table-column label="模型文件路径">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.modelFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="数据集文件路径" width="140">
              <template #default="scope">
                  <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{
                      scope.row.dataSetFilePath
                    }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).message }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="进度">
              <template #default="scope">
                  <span v-if="scope.row && scope.row.resultMsg">
                    {{ parseResultMsg(scope.row.resultMsg).progress }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="输出">
              <template #default="scope">
                <!--            <span v-if="scope.row && scope.row.resultMsg">-->
                <!--            {{ parseResultMsg(scope.row.resultMsg).result }}-->
                <!--            </span>-->
                <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;"
                      :title="'结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth'"
                      v-if="parseResultMsg(scope.row.resultMsg).message === '任务结束'">
                    {{ '结果保存到:/usr/local/agent/' + parseResultMsg(scope.row.resultMsg).taskID + '.pth' }}
                  </span>
                <span v-else>
                    {{ "——" }}
                  </span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50 && parseResultMsg(scope.row.resultMsg).status !== 4"
                           @click="TaskManage(scope.row.taskID, 1)">暂停
                </el-button>
                <el-button type="text"
                           v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status === 4"
                           @click="TaskManage(scope.row.taskID, 2)">恢复
                </el-button>
                <el-button type="text"
                           v-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status < 50"
                           @click="TaskManage(scope.row.taskID, 3)">终止
                </el-button>
                <span type="text"
                      v-else-if="scope.row && scope.row.resultMsg && parseResultMsg(scope.row.resultMsg).status >= 50">已保存退出</span>
                <span type="text"
                      v-else-if="!scope.row && scope.row.resultMsg || !parseResultMsg(scope.row.resultMsg).status">等待响应</span>
                <span type="text" v-else-if="parseResultMsg(scope.row.resultMsg).status === -10">异常退出</span>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination :current-page="getTaskReq.page" :page-size="getTaskReq.pageSize"
                           :page-sizes="[10, 30, 50, 100]" :total="total" layout="total,prev, pager, next,  sizes, jumper"
                           @current-change="handleCurrentChange" @size-change="handleSizeChange" />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    <el-drawer v-model="openDialog" :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false"
      size="60%">

      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ IsCreated ? '新增' : '编辑' }}任务</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="openDialog = false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form ref="dataInfo" :model="dataInfo" :rules="IsCreated ? rules : ''" label-width="100px">
        <el-form-item label="任务名称" prop="taskName">
          <el-input v-model="dataInfo.taskName" />
        </el-form-item>
        <el-form-item label="框架名称" prop="frameName">
          <el-input v-model="dataInfo.frameName" />
        </el-form-item>
        <el-form-item label="任务类型" prop="mode" style="padding-left: 10px" label-width="90px">
          <el-select v-model="dataInfo.mode">
            <el-option label="训练" :value="'train'" />
            <el-option label="推理" :value="'inference'" />
          </el-select>
        </el-form-item>
        <el-form-item label="模型文件" prop="modelFileUUID" style="padding-left: 10px" label-width="90px">
          <el-select v-model="dataInfo.modelFileUUID" placeholder="请选择模型">
            <el-option :key="item.uuid"  :disabled="item.status!==100" v-for="(item) in modelList" :value="item.uuid" >
              <span style="float: left">{{ item.productName }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                <el-tag :type="item.status === 100 ? 'success' : 'info'">{{ item.status === 100 ? '可用' : '下载未完成' }}</el-tag>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="数据集文件" prop="dataSetFileUUID" style="padding-left: 10px" label-width="90px">
          <el-select v-model="dataInfo.dataSetFileUUID" placeholder="请选择数据集">
            <el-option :key="item.uuid" :disabled="item.status!==100" v-for="(item) in datasetList" :value="item.uuid">
              <span style="float: left">{{ item.productName }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                <el-tag :type="item.status === 100 ? 'success' : 'info'">{{ item.status === 100 ? '可用' : '下载未完成' }}</el-tag>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="flex justify-end mt-4">
        <div>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="submitDialog">确 定</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>
<script>
import { getAllShareElementList, GetDownloadElement, } from "@/api/yunguan/element/element";
import SelectImage from "@/components/selectImage/selectImage.vue";
import { deleteRenterTask, getRenterTaskList, updateRenterTask } from "@/api/yunguan/run/task";
import { getRenterResList } from "@/api/yunguan/resource/resInfo";
import { addRenterTask, listTaskProgress, manageRenterTask } from "@/api/yunguan/cloudpods/instance/instance";

export default {
  name: "NodeModel",
  components: { SelectImage },
  data() {
    return {
      activeName: 'baremetal',
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      total: 0,
      IsCreated: true, // true新增 false编辑
      openDialog: false,
      taskType: ["推理", "训练", "其他"],
      dataInfo: {
        id: 0,
        taskName: '',
        frameName: '',
        privateIp: '',
        mode: '',
        modelFileUUID: '',
        dataSetFileUUID: '',
      },
      rules: {
        taskName: [
          { required: true, message: '请输入任务名称', trigger: 'blur' }
        ],
        frameName: [
          { required: true, message: '请输入任务类型', trigger: 'blur' }
        ],
      },
      modelType: 1, //1节点模型 2我的模型
      datasetType: 1, //1节点数据集 2我的数据集
      algoType: 1, //1节点算法 2我的算法
      taskList: [],
      nodeModelList: [], //节点模型列表
      nodeDatasetList: [], //节点数据集列表
      nodeAlgoList: [], //节点算法列表
      elementList: [],

      instanceList: [],//实例数组
      getTaskReq: {
        page: 1,
        pageSize: 10,
        privateIp: '',
      },
      datasetList: [],
      modelList: [],
      taskTimer: null, // 定时器
      instanceInfo: {
        page: 1,
        pageSize: 10,
        type: 1,
      }
    }
  },
  created() {
    // this.GetRenterTaskList()
    // this.GetShareElementList()
    this.GetRenterInfoList()
    // this.startTaskTimer() // 启动定时器
  },
  beforeDestroy() {
    // this.stopTaskTimer() // 停止定时器
  },
  methods: {
    //tab切换
    handleTabsClick() {
      switch (this.activeName) {
        case "baremetal":
          this.instanceInfo.type = 1
          break;
        case "virtual":
          this.instanceInfo.type = 2
          break;
        case "container":
          this.instanceInfo.type = 3
          break;
      }
      this.getTaskReq.privateIp = '' //清空选择的机器
      this.taskList = [] //清空展示的任务
      this.GetRenterInfoList()
    },
    parseResultMsg(jsonStr) {
      try {
        if (jsonStr === ''|| jsonStr === undefined) {
          return;
        }
        return JSON.parse(jsonStr);
      } catch (e) {
        console.error("JSON 解析错误:", e);
        return {};
      }
    },
    HandleChangeInstance() {
      this.GetRenterTask()
    },
    FormatDateTime(dateString) {
      const date = new Date(dateString);
      if (date.toISOString() === '0001-01-01T00:00:00.000Z') {
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
    handleDelete(val) {
      this.$confirm('此操作将永久删除该任务, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteRenterTask(val.ID).then(response => {
          if (response.code == 0) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.GetRenterTaskList(this.searchInfo)
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
    handleEdit(val) {
      this.GetShareElementList()
      this.dataInfo = val
      this.IsCreated = false
      this.openDialog = true
      this.searchInfo.pageSize = 10
    },
    handleSearch() {
    },
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetRenterTask()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetRenterTask()
    },
    addRenterTaskOpenDialog() {
      this.GetDownloadElementList()
      this.searchInfo.pageSize = 10
      this.dataInfo = {}
      this.dataInfo.privateIp = this.getTaskReq.privateIp
    },
    closeDialog() {
      this.openDialog = false
    },
    UpdateRenterTask() {
      updateRenterTask(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.openDialog = false
          this.GetRenterTaskList()
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    AddRenterTask() {
      this.$refs.dataInfo.validate(valid => {
        if (valid) {
          addRenterTask(this.dataInfo).then(response => {
            if (response.code == 0) {
              this.openDialog = false
              this.$message({
                message: '添加成功',
                type: 'success'
              });
              this.GetRenterTask()
            } else {
              this.$message({
                message: '添加失败',
                type: 'error'
              });
            }
          })
        }
      });
    },
    submitDialog() {
      if (this.IsCreated) this.AddRenterTask()
      else this.UpdateRenterTask()
    },
    GetShareElementList() {
      getAllShareElementList().then(response => {
        if (response.code == 0) {
          this.elementList = response.data.list
          this.nodeModelList = response.data.list.filter(item => item.type == 1)
          this.nodeAlgoList = response.data.list.filter(item => item.type == 2)
          this.nodeDatasetList = response.data.list.filter(item => item.type == 3)
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    GetRenterTaskList() {
      getRenterTaskList(this.searchInfo).then(response => {
        if (response.code == 0) {
          this.taskList = response.data.list
          this.total = response.data.total
        } else {
          this.$message({
            message: '获取失败',
            type: 'error'
          });
        }
      })
    },
    GetElementName(id) {
      const tenant = this.elementList.find(t => t.ID === id);
      return tenant ? tenant.name : '';
    },
    GetRenterInfoList() {
      getRenterResList(this.instanceInfo).then(response => {
        if (response.code == 0) {
          this.instanceList = response.data.list.filter(item => item.status.toLowerCase() === "running")
          if (this.instanceList.length > 0) {
            this.getTaskReq.privateIp=this.instanceList[0].privateIp;
            this.GetRenterTask()
          }
        }
      })
    },
    GetRenterTask() {
      this.getTaskReq.page = this.searchInfo.page
      this.getTaskReq.pageSize = this.searchInfo.pageSize
      listTaskProgress(this.getTaskReq).then(res => {
        if (res.code == 0) {
          this.taskList = res.data.list
          this.total = res.data.total
        }
      })
    },
    GetDownloadElementList() {
      let obj = {}
      obj.privateIp = this.getTaskReq.privateIp
      GetDownloadElement(obj).then(res => {
        if (res.code == 0) {
          this.openDialog = true
          this.IsCreated = true
          this.datasetList = res.data.filter(item => item.productType == 2 )
          this.modelList = res.data.filter(item => item.productType == 1)
        }
      })
    },
    TaskManage(id, method) {
      let ip;
      ip = this.getTaskReq.privateIp
      manageRenterTask(ip, id, method).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '操作成功',
            type: 'success'
          })
          this.GetRenterTaskList(this.searchInfo)
        } else {
          this.$message({
            message: '操作失败',
            type: 'error'
          });
        }
      })
    },
    startTaskTimer() {
      this.taskTimer = setInterval(() => {
        if (this.getTaskReq.privateIp) {
          this.GetRenterTask()
        }
      }, 5000); // 每5秒请求一次
    },
    stopTaskTimer() {
      if (this.taskTimer) {
        clearInterval(this.taskTimer)
        this.taskTimer = null
      }
    },
  },
}
</script>