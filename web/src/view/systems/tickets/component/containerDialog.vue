<template>

  <!-- 无状态创建表单 -->
  <el-drawer v-model="showAddDeployment" title="创建" size="520px">
    <el-form label-width="100px"
             ref="deploymentForm"
             label-position="left"
             :model="deploymentAddInfo"
             :rules="deploymentAddRule"
    >
      <el-form-item label="名称" prop="name">
        <el-input
            :disabled="isBindRenter"
            v-model="deploymentAddInfo.name"
            type="en-US"
            placeholder="请输入名称"
        ></el-input>
      </el-form-item>
      <el-form-item label="root密码" prop="rootPasswd">
        <el-input
            v-model="deploymentAddInfo.rootPasswd"

            type="password"
            placeholder="请输入密码"
        ></el-input>
      </el-form-item>
      <el-form-item label="集群" prop="clusterId">
        <el-select @change="handleSelectCluster" v-model="deploymentAddInfo.clusterId">
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
          <el-select v-model="selectImageRepoId" placeholder="镜像仓库选择">
            <el-option
                v-for="repo in imageRepos"
                :key="repo.id"
                :value="repo.id"
                :label="repo.name"
            >
            </el-option>
          </el-select>
          <!-- 镜像名称 -->
          <el-select v-model="selectImage" placeholder="镜像选择">
            <el-option
                v-for="image in images"
                :key="image.name"
                :label="image.name"
                :value="image.name"
            >
            </el-option>
          </el-select>
          <!-- 镜像标签 -->
          <el-select v-model="selectTag" placeholder="标签选择">
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
        <el-input-number disabled style="width: 100%" controls-position="right" :precision="0" type="number" :min="1"
                         v-model="deploymentAddInfo.cpu"
        ></el-input-number>
      </el-form-item>
      <el-form-item label="内存（GB）" prop="memory">
        <el-input-number disabled style="width: 100%" controls-position="right" :precision="0" type="number" :min="1"
                         v-model="deploymentAddInfo.memory"
        ></el-input-number>
      </el-form-item>
<!--      <el-form-item label="绑定租户">-->
<!--        <el-switch-->
<!--            v-model="isBindRenter"-->
<!--            :active-value="true"-->
<!--            :inactive-value="false"-->
<!--            @change="deploymentAddInfo.name=isBindRenter?deploymentAddInfo.name:'';deploymentAddInfo.renterID=''"-->
<!--            active-text="是"-->
<!--            class="ml-2"-->
<!--            inactive-text="否"-->
<!--            inline-prompt-->
<!--            style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"-->
<!--        />-->
<!--      </el-form-item>-->
      <el-form-item prop="renterID" v-if="isBindRenter" label="租户">
        <el-select disabled @change="generateUniqueName(deploymentAddInfo.renterID);" placeholder="请选择租户"
                   v-model="deploymentAddInfo.renterID">
          <el-option v-for="item in renterList"
                     :label="item.username"
                     :value="item.renterID"/>
        </el-select>
      </el-form-item>
      <el-form-item label="GPU">
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
        <el-form-item label="GPU数量" prop="gpu">
          <el-input type="number" :min="0"
                    v-model.number="deploymentAddInfo.gpu"
          ></el-input>
        </el-form-item>
        <el-form-item label="GPU厂商" prop="vendors">
          <el-select v-model="deploymentAddInfo.vendors"
                     placeholder="请选择GPU厂商">
            <el-option label="NVIDIA" :value="'nvidia'"/>
          </el-select>
        </el-form-item>
        <el-form-item label="GPU型号" prop="gpuModel">
          <el-select v-model="deploymentAddInfo.gpuModel">
            <el-option label="GeForce RTX 4090" value="4090"></el-option>
          </el-select>
        </el-form-item>
      </div>

      <el-form-item v-if="deploymentAddInfo.clusterId!=''" label="数据卷">
<!--        <div style="margin-bottom: 2px;width: 100%">-->
<!--          <el-radio-group v-model="createEvs" @change="containerEvsList=[]">-->
<!--            <el-radio :value="false">选择已有</el-radio>-->
<!--            <el-radio :value="true">新创建</el-radio>-->
<!--          </el-radio-group>-->
<!--        </div>-->
        <div v-for="(item, index) in containerEvsList">
          <div v-if="createEvs">
            <div class="selectsContainer" style="margin-top: 5px">
              <el-form-item label="Nas存储">
                <el-select style="width: 100%" v-model="item.containerNas" placeholder="选择容器Nas存储">
                  <el-option
                      v-for="tag in conNasList"
                      :label="tag.name"
                      :value="tag.name"
                  >
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="容量(GB)">
                <el-input-number disabled style="width: 100%" controls-position="right" :precision="0" type="number" :min="1"
                                 v-model="item.size"></el-input-number>
              </el-form-item>
            </div>
            <div class="selectsContainer">
              <el-form-item label="挂载点" style="flex: 2;border: 1px" prop="mountPath">
                <el-input @change="deploymentAddInfo.mountPath=item.mountPath;" style="width: 100%" v-model="item.mountPath"
                          placeholder="例如：/mnt"></el-input>
              </el-form-item>
<!--              <el-button type="text" @click="containerEvsList.splice(index, 1);" >移除</el-button>-->
            </div>
          </div>
          <div v-if="!createEvs" class="flex-container">

            <el-form-item label="云硬盘" style="width: 100%;">
              <el-select style="width: 100%;" v-model="item.containerEvsId" placeholder="云硬盘选择">
                <el-option
                    v-for="tag in conEvsList"
                    :label="tag.name"
                    :value="tag.id"
                >
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item style="width: 100%;" label="挂载点" prop="mountPath">
              <el-input @change="deploymentAddInfo.mountPath=item.mountPath;" style="width: 100%" v-model="item.mountPath"
                        placeholder="例如：/mnt"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="text" @click="containerEvsList.splice(index, 1);" >移除
              </el-button>
            </el-form-item>


          </div>
        </div>
<!--        <div style="margin-top: 5px;width: 100%">-->
<!--          <el-button @click="AddContainerEvs" type="primary" >添加数据卷</el-button>-->
<!--        </div>-->
      </el-form-item>
      <el-form-item label="操作">
        <el-button type="primary" @click="submitForm">提交</el-button>
        <el-button type="primary" @click="cancel">取消</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script setup>
import {getBaseK8SList} from "@/api/cloudpods/baseRes/k8s";
import {getImagesByRepository,} from "@/api/cloudpods/image/image";
import {getImageReposList} from "@/api/cloudpods/image/imageRepos";
import {
  addDeployment,
} from "@/api/cloudpods/k8s/deployment";
import {ElMessage, ElMessageBox} from "element-plus";
import { getRenterList} from "@/api/yunguan/run/renter";
import {getContainerNasList} from "@/api/cloudpods/storage/storage";
import {getContainerEvsList} from "@/api/cloudpods/storage/evs";
import {ref, watch} from "vue";

defineOptions({
  name: 'containerDialog'
})

const props = defineProps({
  containerData: {
    type: Object,
    default: function() {
      return {}
    }
  },
})
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
const renterList = ref([]);
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

const portMappings = ref([]);

// 表单验证规则名字，集群，镜像地址，cpu核数，内存必填
const deploymentAddRule = ref({
  name: [{required: true, message: "请输入名称", trigger: "blur"},
    {
      pattern: /^[A-Za-z][A-Za-z0-9-]{1,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
      message: '名称必须以字母开头，只能包含英文字符和数字，且长度应在2到12个字符之间',
      trigger: 'blur'
    }],
  rootPasswd: [
    {required: true, message: "请输入密码", trigger: "blur"},
    {
      pattern: /^[a-zA-Z0-9]{6,11}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
      message: '密码只能包含英文字符和数字，长度应在6到12个字符之间',
      trigger: 'blur'
    }
  ],
  mountPath: [
    {required: true,
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

const showAddDeployment = ref(false);
const useGPU = ref(false);
const createEvs = ref(false);
const clusterList = ref([]);
const imageRepos = ref([]);
const images = ref([]);
const tags = ref([]);
const selectImageRepoId = ref("");
const selectImage = ref("");
const selectTag = ref("");

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
const emit = defineEmits(['dialogSubmit'])
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
        containerEvsList.value[i].size = containerEvsList.value[i].size + 'G';
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
              portMappings.value = [
                {
                  port: 22,
                  targetPort: 22,
                  nodePort: 30000,
                  protocol: "TCP",
                },
              ];
              deploymentAddInfo.value = {};
              if (isBindRenter.value) {
                emit('dialogSubmit', JSON.parse(res.data).id,3);
              }
            }
          })
          .catch((err) => {
            ElMessage.error("添加失败");
          });
    }
  })

};


const handleAdd = () => {
  //请求镜像仓库列表
  getImageReposList({
    pageIndex: 1,
    pageSize: 10,
  }).then((resp) => {
    imageRepos.value = JSON.parse(resp.data).data;
  });
  getBaseK8SList({
    pageIndex: 1,
    pageSize: 10,
  }).then((resp) => {
    if (resp.code==0){
      const tempList = resp.data.data;
      clusterList.value = tempList.filter(item => item.status == "running");
    }
  })
  GetAllRenterList();
  deploymentForm.value=null;
  deploymentAddInfo.value={};
  deploymentAddInfo.value.cpu=props.containerData.cpu;
  deploymentAddInfo.value.memory=props.containerData.mem;
  containerEvsList.value=[];
  containerEvsList.value.push({
    containerNas: '',
    size: props.containerData.sysDisk+props.containerData.dataDisk,
    containerEvsId: '',
    mountPath: '',
  });
  generateUniqueName(props.containerData.renterID);
  addRenterResInfo.value.renterID=props.containerData.renterID;
  deploymentAddInfo.value.renterID=props.containerData.renterID;
  isBindRenter.value=true;
  useGPU.value=false;
  createEvs.value=true;
  showAddDeployment.value = true;
};
defineExpose({
  handleAdd,
});

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

/* 对话框样式优化 */
:deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background-color: #f5f7fa;
  margin: 0;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.el-dialog__body) {
  padding: 24px;
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

/* Drawer样式优化 */
:deep(.el-drawer) {
  border-radius: 12px 0 0 12px;
}

:deep(.el-drawer__header) {
  margin-bottom: 0;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
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
  padding: 0 16px;
}

:deep(.el-tabs__item) {
  height: 48px;
  line-height: 48px;
  font-size: 15px;
  transition: all 0.3s ease;
}

/* 描述列表样式优化 */
:deep(.el-descriptions) {
  padding: 16px;
  border-radius: 8px;
  background-color: #f8f9fa;
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