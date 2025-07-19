<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box" style="margin-top:16px">
      <div class="gva-btn-list">
        <el-form style="display: flex;justify-content: space-between">
          <el-form-item label="集群选择">
            <el-select
                @change="FlushDeploymentData"
                v-model="selectClusterId"
                style="width: 240px;"
            >
              <el-option
                  v-for="cluster in clusterList"
                  :key="cluster.id"
                  :value="cluster.id"
                  :label="cluster.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="FlushDeploymentData"
            >刷新
            </el-button
            >
            <!--            <el-button type="primary" @click="monitorData">监控</el-button>-->
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div style="margin: 16px 0">
      <el-button type="primary" @click="handleAdd">新增</el-button>
    </div>
    <div class="pod-management-container">
      <!-- 操作按钮区域 -->

      <!-- 无状态表格 -->
      <el-table :data="deploymentList"
                :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}">
        <el-table-column type="selection" width="55"/>
        <el-table-column
            min-width="100px"
            fixed="left"
            label="名称"
            prop="name"
        >
          <template #default="scope">
            <el-button type="text" @click="OpenDrawer(scope.row)">
              {{ scope.row.name }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="所属租户">
          <template #default="scope">
            <el-button type="text">
              {{ GetRenterInfo(scope.row.renterId) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <div v-if="scope.row.pods">
              <el-icon
                  style="vertical-align: middle"
                  class="statusIcon"
                  color="green"
                  v-if="!scope.row.pods[0].warnings"
              >
                <SuccessFilled/>
              </el-icon>
              <span
                  style="vertical-align: middle"
                  v-if="scope.row.pods[0].status === 'Running'"
              >运行中</span
              >
              <span
                  style="vertical-align: middle"
                  v-else-if="scope.row.pods[0].status === 'Pending'"
              >挂起</span
              >
              <span
                  style="vertical-align: middle"
                  v-else-if="scope.row.pods[0].status === 'Succeeded'"
              >成功</span
              >
              <span
                  style="vertical-align: middle"
                  v-else-if="scope.row.pods[0].status === 'Failed'"
              >失败</span
              >
              <span
                  style="vertical-align: middle"
                  v-else-if="scope.row.pods[0].status === 'ContainerCreating'"
              >创建中</span
              >
              <span
                  style="vertical-align: middle"
                  v-else-if="scope.row.pods[0].status === 'Unknown'"
              >未知</span
              >
              <el-tooltip
                  v-if="scope.row.pods[0].warnings"
                  class="box-item"
                  effect="dark"
                  :content="scope.row.pods[0].warnings[0].message"
                  placement="top-start"
              >
                <div>
                  <el-icon
                      style="vertical-align: middle"
                      class="statusIcon"
                      color="red"
                  >
                    <WarnTriangleFilled/>
                  </el-icon>
                  <el-text style="vertical-align: middle">状态异常</el-text>
                </div>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="重启次数">
          <template #default="scope">
                <span v-if="scope.row.pods">
                  <span>{{ scope.row.pods[0].restarts }}</span>
                </span>
          </template>
        </el-table-column>
        <el-table-column label="ip地址" width="180px">
          <template #default="scope">
                <span v-if="scope.row.pods">
                  <span>{{ scope.row.pods[0].podIP }}</span>
                </span>
          </template>
        </el-table-column>
        <el-table-column :show-overflow-tooltip='true' min-width="180px" label="镜像地址">
          <template #default="scope">
                <span>{{
                    scope.row.containerImages === undefined
                        ? "无"
                        : scope.row.containerImages[0].image
                  }}</span>
          </template>
        </el-table-column>
        <el-table-column width="100px" label="创建时间">
          <template #default="scope">
                <span>{{
                    formatTimeDifference(scope.row.creationTimestamp)
                  }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" min-width="100px">
          <template #default="scope">
            <el-button link type="text" @click="connect(scope.row)"
            >远程连接
            </el-button
            >
            <el-button link type="text" @click="removeDeployment(scope.row)"
            >删除
            </el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="searchInfo.page" :page-size="searchInfo.pageSize"
                       :page-sizes="[10, 30, 50, 100]" :total="searchInfo.total"
                       layout="total,  prev, pager, next, sizes,jumper"
                       @current-change="handleCurrentChange" @size-change="handleSizeChange"/>
      </div>
    </div>
    <el-dialog
        v-model="deleteVisible"
        :show-close="false"
        width="520px"
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
    <!-- Pod原始信息展示和修改 -->
    <el-dialog v-model="showRawDialog" class="rawDialog">
      <!--   -->
      {{ rawData }}
    </el-dialog>
    <!-- 无状态创建表单 -->
    <el-dialog v-model="showAddDeployment"
               :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false"
               size="1200px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">创建</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="showAddDeployment=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <el-form label-width="100px"
               ref="deploymentForm"
               label-position="left"
               :model="deploymentAddInfo"
               :rules="deploymentAddRule"
      >
        <el-form-item label="名称" prop="name">
          <el-input
              style="width: 320px"
              disabled
              v-model="deploymentAddInfo.name"
              type="en-US"
              placeholder="请输入名称"
          ></el-input>
        </el-form-item>
        <el-form-item label="root密码" prop="rootPasswd">
          <el-input
              style="width: 320px"
              v-model="deploymentAddInfo.rootPasswd"

              type="password"
              placeholder="请输入密码"
          ></el-input>
        </el-form-item>
        <el-form-item label="集群" prop="clusterId">
          <el-select style="width: 320px" @change="handleSelectCluster" v-model="deploymentAddInfo.clusterId">
            <el-option
                :key="cluster.id"
                :value="cluster.id"
                :label="cluster.name"
                v-for="cluster in clusterList"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="镜像地址" prop="imageUrl">
          <div class="selectsContainer">
            <!-- 镜像仓库选择 -->
            <el-select style="width: 200px" v-model="selectImageRepoId" placeholder="镜像仓库选择">
              <el-option
                  v-for="repo in imageRepos"
                  :key="repo.id"
                  :value="repo.id"
                  :label="repo.name"
              >
              </el-option>
            </el-select>
            <!-- 镜像名称 -->
            <el-select style="width: 200px" v-model="selectImage" placeholder="镜像选择">
              <el-option
                  v-for="image in images"
                  :key="image.name"
                  :label="image.name"
                  :value="image.name"
              >
              </el-option>
            </el-select>
            <!-- 镜像标签 -->
            <el-select style="width: 200px" v-model="selectTag" placeholder="标签选择">
              <el-option
                  v-for="tag in tags"
                  :key="tag"
                  :label="tag"
                  :value="tag"
              >
              </el-option>
            </el-select>
          </div>
        </el-form-item>
        <el-form-item label="cpu核数" prop="cpu">
          <el-input-number controls-position="right" :precision="0" type="number" :min="1"
                           v-model="deploymentAddInfo.cpu"
          ></el-input-number>
        </el-form-item>
        <el-form-item label="内存（GB）" prop="memory">
          <el-input-number controls-position="right" :precision="0" type="number" :min="1"
                           v-model="deploymentAddInfo.memory"
          ></el-input-number>
        </el-form-item>
        <el-form-item label="绑定租户" style="padding-left: 10px" label-width="90">
          <el-switch
              v-model="isBindRenter"
              :active-value="true"
              :inactive-value="false"
              @change="generateUniqueName(0);deploymentAddInfo.renterID='';"
              active-text="是"
              class="ml-2"
              inactive-text="否"
              inline-prompt
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item prop="renterID" v-if="isBindRenter" label="租户">
          <el-select style="width: 320px" @change="generateUniqueName(deploymentAddInfo.renterID);"
                     placeholder="请选择租户"
                     v-model="deploymentAddInfo.renterID">
            <el-option v-for="item in renterList"
                       :label="item.username"
                       :value="item.renterID"/>
          </el-select>
        </el-form-item>
        <el-form-item label="GPU" style="padding-left: 10px" label-width="90">
          <el-switch
              v-model="useGPU"
              :active-value="true"
              :inactive-value="false"
              active-text="启用"
              class="ml-2"
              inactive-text="禁用"
              inline-prompt
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <div v-if="useGPU">
          <el-form-item label="GPU数量" prop="gpu" style="padding-left: 10px" label-width="90">
            <el-input style="width: 320px" type="number" :min="0"
                      v-model.number="deploymentAddInfo.gpu"
            ></el-input>
          </el-form-item>
          <el-form-item label="GPU厂商" prop="vendors" style="padding-left: 10px" label-width="90">
            <el-select style="width: 320px" v-model="deploymentAddInfo.vendors"
                       placeholder="请选择GPU厂商">
              <el-option label="NVIDIA" :value="'nvidia'"/>
            </el-select>
          </el-form-item>
          <el-form-item label="GPU型号" prop="gpuModel" style="padding-left: 10px" label-width="90">
            <el-select v-model="deploymentAddInfo.gpuModel">
              <el-option label="GeForce RTX 4090" value="GeForce RTX 4090"></el-option>
            </el-select>
          </el-form-item>
        </div>

        <el-form-item v-if="deploymentAddInfo.clusterId!=''" label="数据卷" style="padding-left: 10px" label-width="90">
          <div style="display: flex;flex-direction: column;width: 100%">
            <div style="display: block;width: 100%">
              <el-radio-group v-model="createEvs" @change="containerEvsList=[]">
                <el-radio :value="false">选择已有</el-radio>
                <el-radio :value="true">新创建</el-radio>
              </el-radio-group>
            </div>
            <div style="display: flex;flex-direction: column;">
              <div v-if="createEvs" style="display: block;width: 100%">
                <el-table :data="containerEvsList">
                  <el-table-column label="Nas存储" width="150">
                    <template #default="scope">
                      <el-select style="width: 100%" v-model="scope.row.containerNas" placeholder="选择容器Nas存储">
                        <el-option
                            v-for="tag in conNasList"
                            :label="tag.name"
                            :value="tag.name"
                        >
                        </el-option>
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column label="容量(GB)" width="150">
                    <template #default="scope">
                      <el-input-number style="width: 100%" controls-position="right" :precision="0" type="number"
                                       :min="1"
                                       v-model="scope.row.size"></el-input-number>
                    </template>
                  </el-table-column>
                  <el-table-column label="挂载点" width="150">
                    <template #default="scope">
                      <el-input @change="deploymentAddInfo.mountPath=scope.row.mountPath;" style="width: 100%"
                                v-model="scope.row.mountPath"
                                placeholder="例如：/mnt"></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" width="100">
                    <template #default="scope">
                      <el-button type="text" @click="containerEvsList.splice(scope.$index, 1);">移除</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <div v-if="!createEvs" style="display: block;width: 100%">
                <el-table :data="containerEvsList">
                  <el-table-column label="云硬盘">
                    <template #default="scope">
                      <el-select style="width: 100%;" v-model="scope.row.containerEvsId" placeholder="云硬盘选择">
                        <el-option
                            v-for="tag in conEvsList"
                            :label="tag.name"
                            :value="tag.id"
                        >
                        </el-option>
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column label="挂载点">
                    <template #default="scope">
                      <el-input @change="deploymentAddInfo.mountPath=scope.row.mountPath;" style="width: 100%"
                                v-model="scope.row.mountPath"
                                placeholder="例如：/mnt"></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作">
                    <template #default="scope">
                      <el-button type="text" @click="containerEvsList.splice(scope.$index, 1);">移除</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <div style="display: flex;width: 100%;margin-top: -10px ">
                <el-button style="width: 100%" @click="AddContainerEvs">添加数据卷</el-button>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>

      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="cancel">取消</el-button>
        <el-button type="primary" @click="submitForm">提交</el-button>
      </div>
    </el-dialog>
    <!-- 无状态的容器组信息 -->
    <el-drawer v-model="DeploymentPodsShow" size="60%">
      <el-table :data="deploymentPods" stripe border max-height="500">
        <el-table-column
            fixed
            prop="name"
            label="名称"
            width="180"
            show-overflow-tooltip
        >
          <template #default="scope">
            <el-link type="primary">{{ scope.row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column
            prop="podIP"
            label="Ip地址"
            width="180"
        ></el-table-column>
        <el-table-column
            prop="namespace"
            label="命令空间"
            width="180"
        ></el-table-column>
        <el-table-column prop="status" label="状态" width="180">
          <template #default="scope">
            <el-icon
                style="vertical-align: middle"
                class="statusIcon"
                color="green"
                v-if="!scope.row.warnings"
            >
              <SuccessFilled/>
            </el-icon>
            <span
                style="vertical-align: middle"
                v-if="scope.row.status === 'Running'"
            >运行中</span
            >
            <span
                style="vertical-align: middle"
                v-else-if="scope.row.status === 'Pending'"
            >挂起</span
            >
            <span
                style="vertical-align: middle"
                v-else-if="scope.row.status === 'Succeeded'"
            >成功</span
            >
            <span
                style="vertical-align: middle"
                v-else-if="scope.row.status === 'Failed'"
            >失败</span
            >
            <span
                style="vertical-align: middle"
                v-else-if="scope.row.status === 'ContainerCreating'"
            >创建中</span
            >
            <span
                style="vertical-align: middle"
                v-else-if="scope.row.status === 'Unknown'"
            >未知</span
            >
            <el-tooltip
                v-if="scope.row.warnings"
                class="box-item"
                effect="dark"
                :content="scope.row.warnings[0].message"
                placement="top-start"
            >
              <div>
                <el-icon
                    style="vertical-align: middle"
                    class="statusIcon"
                    color="red"
                >
                  <WarnTriangleFilled/>
                </el-icon>
                <el-text style="vertical-align: middle">状态异常</el-text>
              </div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column
            prop="ready"
            label="就绪情况"
            width="180"
        ></el-table-column>
        <el-table-column prop="restarts" label="重启次数"></el-table-column>
        <el-table-column prop="creationTimestamp" label="创建时间">
          <template #default="scope">
            <span>{{ formatTimeDifference(scope.row.creationTimestamp) }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
    <!--云主机详情-->
    <el-drawer v-model="renterVisible" :before-close="CloseDrawer" :title="'详情'" size="50%">
      <div style="margin: 0 10px" v-if="resInfo.renterID!=undefined">
        <el-descriptions :column="2" border title="基本信息">
          <el-descriptions-item label="所属租户">{{ GetRenterInfo(resInfo.renterID) }}</el-descriptions-item>
        </el-descriptions>
        <div v-if="orderInfo.code!=undefined">
          <el-descriptions :column="2" border>
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
                                orderInfo.billingType === 32 ? '年' : '无限制'
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
                  {{ item.value[item.valueIndex != undefined ? item.valueIndex : 0] }}
                </el-descriptions-item>
              </el-descriptions>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-drawer>

    <!--容器详情-->
    <el-drawer v-model="visible" :close-on-click-modal="true" :close-on-press-escape="true" :show-close="false"
               size="1200px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 16px;font-weight: 500;color: rgba(29, 33, 41, 1);">操作日志</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="visible=false"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close/>
        </el-icon>
      </div>
      <div style="margin: 0 10px">
        <el-tabs v-model="activeName">
          <el-tab-pane name="操作日志" label="操作日志">
            <Logs :log-data="logData"/>
          </el-tab-pane>
          <el-tab-pane name="监控" label="监控">
            <div
                style="margin-top: 10px;margin-bottom: 10px;display: flex;justify-content: space-between;align-items: center;/* 垂直居中 */gap: 10px;">
              <div style="display: flex;align-items: center;">
                <span>取值间隔（分钟）</span>
                <el-select @change="GetMonitorData" style="width: 200px" v-model="monitorInfo.interval">
                  <el-option v-for="(item,index) in monitorConfigData[monitorInfo.duration]"
                             :label="convertTimeInterval(item)"
                             :value="item"/>
                </el-select>
              </div>
              <div style="display: flex;align-items: center;background:  #F2F3F5;padding:0 5px">
                <el-radio-group v-model="monitorInfo.duration" @change="handleRadioChange">
                  <el-radio-button type="key" :value="1">近1小时</el-radio-button>
                  <el-radio-button type="key" :value="6">近6小时</el-radio-button>
                  <el-radio-button type="key" :value="12">近12小时</el-radio-button>
                  <el-radio-button type="key" :value="24">近24小时</el-radio-button>
                  <!--                      <el-radio-button value="168h">近1周</el-radio-button>-->
                  <!--                      <el-radio-button value="720h">近1月</el-radio-button>-->
                </el-radio-group>
                <el-button type="key" @click="GetMonitorData" icon="refresh" value="refresh"></el-button>
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
              <GvaCard v-if="monitorData.cpuData.xData.length>0" custom-class="h-42" title="CPU使用率">
                <GvaChart :line-title="'CPU使用率'" :x-line-data="monitorData.cpuData.xData"
                          :y-line-data="monitorData.cpuData.yData"
                          :type="8"/>
              </GvaCard>
              <GvaCard v-if="monitorData.memData.xData.length>0" custom-class="h-42" title="内存使用率">
                <GvaChart :line-title="'内存使用率'" :x-line-data="monitorData.memData.xData"
                          :y-line-data="monitorData.memData.yData"
                          :type="8"/>
              </GvaCard>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-drawer>
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
      <span>此操作将永久删除该容器，是否继续?</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="info" @click="deleteVisible = false">取消</el-button>
          <el-button
              type="primary"
              @click="confirmDelete1"
          >
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {getBaseK8SList} from "@/api/cloudpods/baseRes/k8s";
import {getImagesByRepository,} from "@/api/cloudpods/image/image";
import {getImageReposList} from "@/api/cloudpods/image/imageRepos";
import {
  addDeployment,
  deleteDeployment,
  deploymentMonitor,
  getDeploymentPods,
  listDeployment,
  webConnectDeployment,
} from "@/api/cloudpods/k8s/deployment";
import {getNameSpaceList, getPodsList} from "@/api/cloudpods/k8s/pods";
import router from "@/router";
import {ElMessage, ElMessageBox} from "element-plus";
import {onMounted, onUnmounted, ref, watch} from "vue";
import {addRenterRes, getRenterList, queryRenterResInfo} from "@/api/yunguan/run/renter";
import Logs from "@/view/cloudpods/component/log.vue";
import {getContainerNasList} from "@/api/cloudpods/storage/storage";
import {getContainerEvsList} from "@/api/cloudpods/storage/evs";
import {GvaCard, GvaChart} from "@/view/safe/monitor/componenst";
import {graphic} from "echarts";
import {format} from "date-fns";
import _ from "lodash";
// const monitorData=()=>{
//   deploymentMonitor({
//     interval:1, //时间间隔
//     duration:1, //监控时长
//     k8sName:"k8s", //集群名称
//     name:"nodeqwe",//容器 containerImages[0].name;
//     monitorType:"cpu",//监控类型 cpu memory
//   }).then(res=>{
//     if (res.code==0){
//       ElMessage({
//         type:"success",
//         message:"监控成功",
//       })
//     }
//   })
// }
const monitorData = ref({
  cpuData: {
    xData: [],
    yData: [],
  },
  memData: {
    xData: [],
    yData: [],
  }
})

//radio按钮监控时间变化
const handleRadioChange = () => {
  monitorInfo.value.interval = monitorConfigData.value[monitorInfo.value.duration][0]
  GetMonitorData(monitorInfo.value)
};

const clearData = () => {
  monitorData.value.cpuData.xData = []
  monitorData.value.cpuData.yData = []
  monitorData.value.memData.yData = []
  monitorData.value.memData.xData = []
};
const chartData = ref({
  type: 'line',
  smooth: true,
  showSymbol: false,
  symbol: 'circle',
  symbolSize: 6,
  name: '',
  data: [],
  lineStyle: {
    normal: {
      shadowBlur: 3,
      shadowColor: 'rgba(0, 0, 0, .12)',
      shadowOffsetX: 4,
      shadowOffsetY: 4,
    },
  },
  emphasis: {
    focus: 'series',
    itemStyle: {
      borderColor: 'blue',
      borderWidth: 2,
    },
    areaStyle: {
      color: new graphic.LinearGradient(0, 0, 0, 1, [
        {offset: 0, color: 'rgba(70, 130, 180, 0.7)'},
        {offset: 1, color: 'rgba(70, 130, 180, 0)'}
      ]),
    }
  },
  areaStyle: {
    color: new graphic.LinearGradient(0, 0, 0, 1, [
      {offset: 0, color: 'rgba(70, 130, 180, 0.5)'},
      {offset: 1, color: 'rgba(70, 130, 180, 0)'}
    ]),
  },
});
//获取监控数据
const GetMonitorData = (obj) => {
  clearData();
  deploymentMonitor({
    interval: monitorInfo.value.interval, //时间间隔
    duration: monitorInfo.value.duration, //监控时长
    k8sName: detailInfo.value.k8sName, //集群名称
    name: detailInfo.value.name,//容器 containerImages[0].name;
    monitorType: "cpu",//监控类型 cpu memory
  }).then(res => {
    if (res.code == 0 ) {
      let tempChartData = _.cloneDeep(chartData.value)
      tempChartData.name ="cpu使用率";
      for (let j = 0; j < res.data.data.result[0].values.length; j++) {
        let memData = res.data.data.result[0].values[j];
        tempChartData.data.push(Number(memData[1])*100);
        const formattedDate = format(new Date(memData[0]*1000), 'MM-dd HH:mm');
        monitorData.value.cpuData.xData.push(formattedDate);
      }
      monitorData.value.cpuData.yData.push(tempChartData);
    }
  })
  deploymentMonitor({
    interval: monitorInfo.value.interval, //时间间隔
    duration: monitorInfo.value.duration, //监控时长
    k8sName: detailInfo.value.k8sName, //集群名称
    name: detailInfo.value.name,//容器 containerImages[0].name;
    monitorType: "memory",//监控类型 cpu memory
  }).then(res => {
    if (res.code == 0 ) {
      let tempChartData = _.cloneDeep(chartData.value)
      tempChartData.name ="内存使用率";
      for (let j = 0; j < res.data.data.result[0].values.length; j++) {
        let memData = res.data.data.result[0].values[j];
        tempChartData.data.push(Number(memData[1])*100);
        const formattedDate = format(new Date(memData[0]*1000), 'MM-dd HH:mm');
        monitorData.value.memData.xData.push(formattedDate);
      }
      monitorData.value.memData.yData.push(tempChartData);
    }
  })
};
const monitorConfigData = ref({
  [1]: [1, 5],
  [6]: [3, 5, 10, 15,],
  [12]: [5, 10, 20, 30],
  [24]: [10, 20, 30, 60],
});
//时间间隔转化
const convertTimeInterval = (interval) => {
  const timeMap = {
    1: '1分钟',
    3: '3分钟',
    5: '5分钟',
    10: '10分钟',
    15: '15分钟',
    20: '20分钟',
    30: '30分钟',
    60: '1小时',
    // 可以根据需要添加更多映射
  };
  return timeMap[interval] || interval; // 如果没有找到匹配项，则返回原值
};
const monitorInfo = ref({
  interval: 1,
  duration: 1,
  k8sName: '',
  name: '',
  monitorType: '',
})
const deleteVisible = ref(false);
const deleteId = ref('');
const HandleDeleteDialog = (val) => {
  deleteVisible.value = true;
  deleteId.value = val;
}
const confirmDelete = () => {
  removeDeployment(deleteId);
  deleteVisible.value = false;
}
const isUpdateRenterRes = ref(false);
const modRenterRes = ref({
  id: '',
  resourceID: '',
  startTime: '',
  endTime: '',
});
const isBindRenter = ref(false);
const addRenterResInfo = ref({
  renterID: '',
  type: 3,
  sshHost: '',
  sshUser: '',
  sshPwd: '',
  resourceID: '',
  privateIp: '',
  aiList: '[]',
  startTime: '',
})
const generateRandomLetter = () => {
  // 随机生成一个字母（小写字母）
  const letters = "abcdefghijklmnopqrstuvwxyz";
  const randomIndex = Math.floor(Math.random() * letters.length);
  return letters[randomIndex];
};

const timestampToLetters = (timestamp) => {
  // 将时间戳转换为字符串，然后通过 ASCII 码转换为字母
  let timestampStr = timestamp.toString();
  // 确保至少有8位数字，不足的部分在前面补零
  if (timestampStr.length < 8) {
    timestampStr = timestampStr.padStart(8, '0');
  }
  // 截取后8位数字
  if (timestampStr.length > 8) {
    timestampStr = timestampStr.slice(-8);
  }
  let result = "";
  for (let char of timestampStr) {
    // 将字符转换为整数
    const digit = parseInt(char, 10);
    // 加上97得到新的ASCII值
    const asciiVal = digit + 97;
    // 转换为对应的ASCII字符
    result += String.fromCharCode(asciiVal);
  }
  return result;
};

const generateUniqueName = (renterId) => {
  // 生成纳秒级时间戳并转换为字符串，取前 8 位以控制长度
  const timeStr = timestampToLetters(Date.now());

  // 将租户 ID 转换为字符串
  const renterIdStr = renterId.toString();

  // 组合时间戳、随机字符串和租户 ID
  let uniqueName = `${timeStr}-${renterIdStr}-`;

  // 确保名称不超过 20 个字符
  if (uniqueName.length > 20) {
    uniqueName = uniqueName.slice(0, 20);
  }

  // 随机生成一个字母作为结尾
  if (!/[a-zA-Z]$/.test(uniqueName)) {
    uniqueName = uniqueName.slice(0, -1) + generateRandomLetter();
  }
  deploymentAddInfo.value.name = uniqueName;
};
const conNasList = ref([]);
const conEvsList = ref([]);
const containerEvsList = ref([]);
const containerEvsInfo = ref({
  containerNas: '',
  size: 1,
  containerEvsId: '',
  mountPath: '',
})
const AddContainerEvs = () => {
  containerEvsList.value.push({
    containerNas: '',
    size: 1,
    containerEvsId: '',
    mountPath: '',
  })
}
const handleSelectCluster = () => {
  getContainerNasList({pageIndex: 1, pageSize: 10000, cluster: deploymentAddInfo.value.clusterId}).then(res => {
    if (res.code == 0) {
      conNasList.value = res.data.list;
    }
  })
  getContainerEvsList({pageIndex: 1, pageSize: 10000, cluster: deploymentAddInfo.value.clusterId}).then(res => {
    if (res.code == 0) {
      conEvsList.value = res.data.list;
    }
  })
}
const deploymentForm = ref(null)
const deploymentAddInfo = ref({
  name: ref(''),
  renterID: ref(''),
  clusterId: ref(''),
  cpu: ref(1),
  memory: ref(1),
  rootPasswd: ref(''),
  portMappings: ref([]),
  containerEvs: ref([]),
});

const searchInfo = ref({
  page: 1,
  pageSize: 10,
  total: 0,
})
const handleCurrentChange = (val) => {
  searchInfo.value.page = val;
  FlushDeploymentData();
}
const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val;
  FlushDeploymentData();
}

const activeName = ref('操作日志')
const visible = ref(false)
const logData = ref({
  id: '',
  isContainer: true,
  ownerKind: '',
  ownerName: '',
  namespace: '',
  clusterId: '',
})
const detailInfo=ref({
  k8sName:'',
  name:'',
})
const OpenDrawer = (val) => {
  logData.value.id = val.id;
  logData.value.ownerKind = val.kind;
  logData.value.ownerName = val.name;
  logData.value.namespace = val.namespace;
  logData.value.clusterId = val.clusterId;
  detailInfo.value.k8sName=val.cluster;
  detailInfo.value.name=val.containerImages[0].name;
  GetMonitorData();
  visible.value = true;
}
const DeploymentPodsShow = ref(false);

const deploymentPods = ref([]);

const resInfo = ref({});
const orderInfo = ref({});
const specDetail = ref([]);
const renterList = ref([]);
const renterVisible = ref(false);
//打开绑定租户对话框
const GetResourceInfo = (val) => {
  resInfo.value = {}
  orderInfo.value = {}
  specDetail.value = [];
  queryRenterResInfo({resourceID: val.id}).then(res => {
    if (res.code == 0) {
      resInfo.value = res.data
      if (res.data.description != "") {
        orderInfo.value = JSON.parse(res.data.description)
        orderInfo.value.specJson = JSON.parse(orderInfo.value.specJson)
        orderInfo.value.storage = JSON.parse(orderInfo.value.storage)
        orderInfo.value.specJson.push({name: '存储类型', value: [orderInfo.value.storage[0].name], valueIndex: 0})
        orderInfo.value.specJson.push({name: '操作系统', value: [orderInfo.value.os], valueIndex: 0})
      }
    }
    renterVisible.value = true;
  })
};
const CloseDrawer = () => {
  renterVisible.value = false;
  visible.value = false;
}


//获取所有租户列表
const GetAllRenterList = () => {
  let obj = {};
  obj.page = 1;
  obj.pageSize = 1000;
  getRenterList(obj).then(res => {
    if (res.code == 0) {
      renterList.value = res.data.list;
    }
  })
};

//获取根据资源ID获取租户信息
function GetRenterInfo(val) {
  console.log(val)
  // 遍历租户资源列表
  for (const renter of renterList.value) {
    let res = renter
    if (res.renterID == val) {
      return res.username;
    }
  }
  // 如果没有找到匹配的租户，返回“无”
  return "无";
};
const showDeploymentPods = (row) => {
  //请求对应Pod拥有的数据
  getDeploymentPods({
    scope: "system",
    show_fail_reason: true,
    owner_kind: "Deployment",
    owner_name: row.name,
    namespace: "ygpt",
    cluster: selectClusterId.value,
    details: true,
    limit: 20,
  }).then((resp) => {
    deploymentPods.value = resp.data.data;
  });
  DeploymentPodsShow.value = true;
};

// const portMappings = ref([
//   {
//     port: 22,
//     targetPort: 22,
//     nodePort: 30000,
//     protocol: "TCP",
//   },
// ]);
const portMappings = ref([]);

// 表单验证规则名字，集群，镜像地址，cpu核数，内存必填
const deploymentAddRule = ref({
  name: [{required: true, message: "请输入名称", trigger: "blur"},],
  rootPasswd: [
    {required: true, message: "请输入密码", trigger: "blur"},
    {
      pattern: /^[a-zA-Z0-9]{6,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
      message: '密码只能包含英文字符和数字，长度应在6到12个字符之间',
      trigger: 'blur'
    }
  ],
  mountPath: [
    {
      required: true,
      pattern: /^\/[^/]+/, // 必须以 / 开头，且后面至少有一个非 / 字符
      message: '挂载点必须以 / 开头，且不能仅为 /',
      trigger: 'blur',
    },
  ],
  clusterId: [{required: true, message: "请选择集群", trigger: "blur"}],
  renterID: [{required: true, message: "请选择租户", trigger: "blur"}],
  imageUrl: [{required: true, message: "请输入镜像地址", trigger: "blur"}],
  cpu: [{required: true, message: "请输入cpu核数", trigger: "blur"}],
  memory: [{required: true, message: "请输入内存", trigger: "blur"}],
  portMappings: [
    {required: true, message: "请添加端口映射", trigger: "blur"},
  ],
});
const currentProductId = ref(null);
const showAddDeployment = ref(false);
const useGPU = ref(false);
const createEvs = ref(false);
const activeTab = ref("deployment");
const rawData = ref({});
const showRawDialog = ref(false);
const search = ref("");
const clusterList = ref([]);
const nameSpaceList = ref([]);
const selectClusterId = ref("");
const selectNameSpace = ref("");
const scope = ref("system");
const podTotal = ref(0);
const podList = ref([]);
const imageRepos = ref([]);
const images = ref([]);
const tags = ref([]);
const selectImageRepoId = ref("");
const selectImage = ref("");
const selectTag = ref("");
const deploymentList = ref([]);

// 镜像仓库选择更改
watch(selectImageRepoId, (newVal, oldVal) => {
  //请求镜像仓库有的镜像
  getImagesByRepository({
    id: selectImageRepoId.value,
  }).then((resp) => {
    images.value = resp.data.list;
  });
});

//镜像选择更改
watch(selectImage, (newVal, oldVal) => {
  //更改标签
  images.value.forEach((item) => {
    if (newVal === item.name) {
      tags.value = item.tags;
    }
  });
});

const cancel = () => {
  ElMessageBox.confirm("取消将清除所有的输入信息", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    showAddDeployment.value = false;
    portMappings.value = [];
    deploymentAddInfo.value = {};
  });
};

const submitForm = () => {
  deploymentAddInfo.value.portMappings = portMappings.value;
  //求出实际的镜像地址
  imageRepos.value.forEach((item) => {
    if (item.id === selectImageRepoId.value) {
      deploymentAddInfo.value.imageUrl = item.url + "/";
    }
  });
  images.value.forEach((item) => {
    if (item.name === selectImage.value) {
      const parts = item.name.split("/");
      deploymentAddInfo.value.imageUrl += parts[1] + ":" + selectTag.value;
    }
    deploymentAddInfo.value.imageUrl = deploymentAddInfo.value.imageUrl.replace(
        "https://",
        ""
    );
  });
  deploymentAddInfo.value.containerEvs = containerEvsList.value;
  console.log("数组的值：")
  console.log(containerEvsList.value)
  console.log(containerEvsList[0])
  if (containerEvsList.value.length > 0) {
    if (createEvs.value) {
      const seen = new Map(); // 用于存储已见过的 name 和 age 组合
      for (let i = 0; i < containerEvsList.value.length; i++) {
        containerEvsList.value[i].containerEvsId = '';
        containerEvsList.value[i].size = containerEvsList.value[i].size + 'Gi';
      }
      for (const item of containerEvsList.value) {
        const key = item.mountPath; // 使用 name 和 age 拼接成唯一键
        if (seen.has(key)) {
          ElMessage.error("挂载点不能重复");
          return; // 如果发现重复，直接退出
        }
        seen.set(key, true); // 记录当前组合
      }
    } else {
      for (let i = 0; i < containerEvsList.value.length; i++) {
        containerEvsList.value[i].size = '';
        containerEvsList.value[i].containerNas = '';
      }
      const seen = new Map(); // 用于存储已见过的 name 和 age 组合

      for (const item of containerEvsList.value) {
        const key = item.mountPath; // 使用 name 和 age 拼接成唯一键
        if (seen.has(key)) {
          ElMessage.error("挂载点不能重复");
          return; // 如果发现重复，直接退出
        }
        seen.set(key, true); // 记录当前组合
      }
      console.log(containerEvsList)
    }
  }
  deploymentForm.value.validate(async valid => {
    if (valid) {
      addDeployment(deploymentAddInfo.value)
          .then(res => {
            if (res.code == 0) {
              ElMessage.success("添加成功");
              showAddDeployment.value = false;
              FlushDeploymentData();
              portMappings.value = [
                {
                  port: 22,
                  targetPort: 22,
                  nodePort: 30000,
                  protocol: "TCP",
                },
              ];
              if (isBindRenter.value) {
                createRenterRes(JSON.parse(res.data))
              }
            }
          })
          .catch((err) => {
            ElMessage.error("添加失败");
          });
    }
  })

};
const createRenterRes = (val) => {
  addRenterResInfo.value.resourceID = val.id;
  addRenterResInfo.value.renterID = deploymentAddInfo.value.renterID;
  const now = new Date();
  const cTime = new Date(now.getTime() + 8 * 60 * 60 * 1000);
  const isoString = cTime.toISOString().replace("Z", "+08:00");
  modRenterRes.value.startTime = isoString;
  addRenterRes(addRenterResInfo.value).then(res => {
    if (res.code == 0) {
      ElMessage.success("绑定租户成功");
      isUpdateRenterRes.value = true;
      modRenterRes.value.id = res.data.ID;
      modRenterRes.value.resourceID = val.id;
    }
  })
}
const connect = (row) => {
  webConnectDeployment({
    id: row.id,
  }).then((resp) => {
    if (resp.code === 0) {
      const data = JSON.parse(resp.data).connect_params;
      console.log(data);
      const routeData = router.resolve({
        name: "terminal",
        query: {
          id: data,
        },
      });
      window.open(routeData.href, "_blank");
    }
  });
};

const removeDeployment = () => {
  deleteVisible.value = true;
  currentProductId.value = null;
};

// 确认删除
const confirmDelete1 = async () => {
  try {
    let resp;

    resp = await deleteDeployment({id: currentProductId.value});


    if (resp.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      });
      await FlushDeploymentData(); // 刷新表格数据
      emits('confirm'); // 通知父组件删除成功
    } else {
      ElMessage({
        showClose: true,
        message: '删除失败',
        type: 'error',
      });
    }
  } catch (error) {
    console.error('删除失败:', error);
  } finally {
    deleteVisible.value = false;
    FlushDeploymentData();
  }
};

const resetForm = () => {
};
const addPortMappings = () => {
  console.log(portMappings.value);
  portMappings.value.push({
    port: 80,
    targetPort: 80,
    nodePort: 30000,
    protocol: "TCP",
  });
};

const removePortMappings = (index) => {
  //根据下标删除
  portMappings.value.splice(index, 1);
};

/**
 * 刷新页面数据
 */
const FlushData = () => {
  //请求集群信息
  getBaseK8SList({
    pageIndex: 1,
    pageSize: 10,
    status: "running"
  }).then((resp) => {
    const tempList = resp.data.data;
    clusterList.value = tempList.filter(item => item.status == "running");
    //请求集群信息
    if (clusterList.value.length > 0) {
      selectClusterId.value = clusterList.value[0].id;
    }
    //请求对应NamSpace信息
    getNameSpaceList({
      cluster: selectClusterId.value,
      limit: 20,
      scope: scope.value,
    }).then((resp) => {
      nameSpaceList.value = resp.data.list;
      // 请求Pod信息
      getPodsList({
        scope: scope.value,
        show_fail_reason: true,
        detail: true,
        cluster: selectClusterId.value,
        namespace: selectNameSpace.value,
        limit: 20,
      }).then((resp) => {
        podList.value = resp.data.list.map((item) => {
          return JSON.parse(item);
        });
        podTotal.value = resp.data.total;
      });
    });
  });
};

const handleAdd = () => {
  //请求镜像仓库列表
  getImageReposList({
    pageIndex: 1,
    pageSize: 10,
  }).then((resp) => {
    imageRepos.value = JSON.parse(resp.data).data;
  });
  deploymentForm.value = null;
  deploymentAddInfo.value = {};
  deploymentAddInfo.value.cpu = 1;
  deploymentAddInfo.value.memory = 1;
  isBindRenter.value = false;
  useGPU.value = false;
  createEvs.value = false;
  generateUniqueName(0);
  showAddDeployment.value = true;
};

const changeTab = (tabName) => {
  //根据不同的tab去更新不同的数据
  switch (tabName) {
    case "pod":
      FlushData();
      break;
    case "deployment":
      FlushDeploymentData();
      break;
  }
};

const FlushDeploymentData = () => {
  GetAllRenterList();
  getBaseK8SList({
    pageIndex: 1,
    pageSize: 100000,
    status: "running"
  }).then((resp) => {
    if (resp.code == 0) {
      clusterList.value = resp.data.data;
      //请求集群信息
      if (clusterList.value.length > 0) {
        if (selectClusterId.value == "") selectClusterId.value = clusterList.value[0].id;
        // 请求无状态的信息
        listDeployment({
          cluster: selectClusterId.value,
          pageIndex: searchInfo.value.page,
          pageSize: searchInfo.value.pageSize,
        }).then(res => {
          if (res.code == 0) {
            deploymentList.value = res.data.data;
            searchInfo.value.total = res.data.total;
            // 请求无状态的Pod并且聚合
            deploymentList.value.forEach((item) => {
              getDeploymentPods({
                scope: "system",
                show_fail_reason: true,
                owner_kind: "Deployment",
                owner_name: item.name,
                namespace: "ygpt",
                cluster: selectClusterId.value,
                details: true,
              }).then(resp => {
                if (resp.code == 0) {
                  item.pods = resp.data.data;
                }
              })
            });
          }
        });
      }
    }
  });
};

const handleOpenRaw = (row) => {
  rawData.value = row;
  showRawDialog.value = true;
};

onMounted(() => {
  FlushDeploymentData(); //请求数据
});
onUnmounted(() => {
})

function formatTimeDifference(dateString) {
  // 将日期字符串转换为时间戳
  const givenDate = new Date(dateString).getTime();
  // 获取当前时间戳
  const currentDate = Date.now();
  // 计算时间差（以毫秒为单位）
  const timeDifference = currentDate - givenDate;

  // 一天的毫秒数
  const oneDay = 24 * 60 * 60 * 1000;

  if (timeDifference < oneDay) {
    // 小于一天，计算小时和分钟
    const hours = Math.floor(timeDifference / (1000 * 60 * 60));
    const minutes = Math.floor(
        (timeDifference % (1000 * 60 * 60)) / (1000 * 60)
    );
    if (isNaN(hours) || isNaN(minutes)) {
      return `无`;
    }
    return `${hours}小时${minutes}分钟前`;
  } else {
    // 大于一天，计算天数
    const days = Math.floor(timeDifference / oneDay);
    if (isNaN(days)) {
      return `无`;
    }
    return `${days}天前`;
  }
}

// 格式化Pod就绪状态字符串
const formatPodReadyStatusString = (podsInfo) => {
  if (podsInfo === undefined || !podsInfo) {
    return `无`;
  } else {
    const statusMessages = [];
    for (const [status, count] of Object.entries(podsInfo)) {
      if (count > 0) {
        statusMessages.push(`${status}: ${count}`);
      }
    }
    return statusMessages.length > 0
        ? statusMessages.join(", ")
        : "所有 Pod 状态正常";
  }
};
</script>

<style scoped lang="scss">
/* 整体容器样式 */
.pod-management-container {
  .button-group {
    margin-top: 16px;
    margin-bottom: 16px;
    display: flex;
    gap: 8px;
  }
}

.container {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

/* 表格样式优化 */
:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
  margin: 16px 0;
  transition: all 0.3s ease;
}

:deep(.el-table th) {
  background-color: #f5f7fa;
  font-weight: 600;
  color: #606266;
}

:deep(.el-table td) {
  padding: 12px 0;
}

/* 按钮样式优化 */
:deep(.el-button) {
  border-radius: 6px;
  transition: all 0.3s ease;
  font-weight: 500;
}

:deep(.el-button:not(.is-disabled):hover) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 表单样式优化 */
:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-input), :deep(.el-select) {
  width: 100%;
}

/* 状态图标优化 */
.statusIcon {
  margin-right: 8px;
  font-size: 16px;
  vertical-align: middle;
  transition: all 0.3s ease;
}

/* 分页容器样式 */
.pagination-container {
  margin-top: 24px;
  padding: 16px 0;
  display: flex;
  justify-content: center;
  background-color: #fff;
  border-radius: 8px;
}

/* 选择器容器样式优化 */
.selectsContainer {
  display: flex;
  gap: 2px;
  flex-wrap: wrap;
}

.selectsContainer > * {
  flex: 1;
  min-width: 200px;
}


/* 链接样式优化 */
:deep(.el-link) {
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-link:hover) {
  text-decoration: none;
}

/* Tabs样式优化 */
:deep(.el-tabs__nav-wrap) {

}

:deep(.el-tabs__item) {
  height: 48px;
  line-height: 48px;
  font-size: 15px;
  transition: all 0.3s ease;
}


:deep(.el-descriptions__title) {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
}

.flex-container {
  display: flex;
  align-items: center; /* 垂直居中对齐 */
  gap: 8px; /* 子元素之间的间距 */
  margin-top: 2px;
}

.select-container {
  display: flex;
  flex: 1;
  align-items: center;
  border: 1px solid #dcdfe6; /* 边框 */
}

.label {
  padding: 0 15px; /* 内边距 */
  background-color: #f0f2f5; /* 背景颜色 */
  font-size: 14px; /* 字体大小 */
  color: #606266; /* 文字颜色 */
  width: 100px;
  border-right: 1px solid #dcdfe6; /* 右侧边框 */
}

.mountPathStyle {

}
</style>
<style>
/* 清除 el-form-item 的所有样式 */
.el-form-item.mountPathStyle {
  margin: 0 !important;
  padding: 0 !important;
  border: none !important;
  background: none !important;
  line-height: normal !important;
}

.el-form-item.mountPathStyle .el-form-item__label {
  margin: 0 !important;
  padding: 0 !important;
  font-size: inherit !important;
  color: inherit !important;
}

.el-form-item.mountPathStyle .el-form-item__error {
  display: none !important;
}
</style>