<template>
  <div style="padding: 24px;margin-top: 16px;background-color: white">
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-form
            ref="searchForm"
            :inline="true"
            :model="searchInfo"
            style="display: flex;justify-content: space-between"
        >
          <el-form-item>
            <el-form-item label="规则名称">
              <el-input
                  v-model="searchInfo.name"
                  placeholder="请输入"
              />
            </el-form-item>
            <el-form-item label="适用产品类别" >
              <el-select
                  style="width: 240px;"
                  v-model="searchInfo.type"
                  placeholder="请选择"
                  clearable
              >
                <el-option v-for="(item,index) in specType"
                           :label="item"
                           :key="index"
                           :value="index+1"/>
              </el-select>
            </el-form-item>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                @click="GetConfigInfoList()"
            >
              查询
            </el-button>
            <el-button
                type="info"
                @click="searchInfo.name='';searchInfo.type=nil;GetConfigInfoList()"
            >重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

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

      <el-table
          v-loading="false"
          :cell-style="{ 'text-align': 'left' }"
          :data="configList"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%; font-size: 15px"
      >
<!--        <el-table-column label="规格类别">-->
<!--          <template #default="scope">-->
<!--            <span>{{ scope.row.tag }}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
        <el-table-column label="规格名称" >
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="适用产品类别" >
          <template #default="scope">
            <span>{{ this.GetTypeName(scope.row.type) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="规格内容">
          <template #default="scope">
            <div style="display: flex; flex-wrap: wrap;gap: 4px">
              <el-tag style="" v-for="(item, key) in scope.row.value" :key="key">
                {{ item }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right">
          <template #default="scope">
            <el-button type="text" link @click="handleEdit(scope.row)">修改</el-button>
            <el-button type="text" link @click="openDeleteDialog(scope.row.ID)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="searchInfo.pageIndex"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer v-model="dialogVisible"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :show-close="false"
               size="520px">
      <div class="flex justify-between items-center" style="margin-bottom: 24px">
        <span class="text-lg" style="font-size: 24px;font-weight: 500;color: rgba(29, 33, 41, 1);">{{ IsCreated ? '新增' : '编辑' }}云产品规格</span>
        <el-icon
            class="cursor-pointer hover:text-gray-500 transition-colors"
            @click="CloseDialog"
            style="width: 16px; height: 16px; font-size: 16px"
        >
          <Close />
        </el-icon>
      </div>
      <el-form
          ref="dataInfo"
          :model="dataInfo"
          label-width="92px"
          :rules="IsCreated?rules:''"
      >
<!--        <el-form-item label="规格类别" prop="tag">-->
<!--          <el-select @change="ChangeTag" v-model="dataInfo.tag" placeholder="请选择规格类别">-->
<!--            <el-option  v-for="(item,index) in tagList"-->
<!--                       :label="item"-->
<!--                       :key="index"-->
<!--                       :value="item"/>-->
<!--          </el-select>-->
<!--        </el-form-item>-->
        <el-form-item label="适用产品" prop="type">
          <el-select style="width: 320px" v-model="dataInfo.type" multiple clearable placeholder="请选择产品类别">
            <el-option v-for="(item,index) in specType"
                       :label="item"
                       :key="index"
                       :value="index+1"/>
          </el-select>
        </el-form-item>
<!--        <el-form-item label="规格名称" prop="name">-->
<!--          <el-input v-model="dataInfo.name" placeholder="请输入规格名称"></el-input>-->
<!--        </el-form-item>-->
        <el-form-item label="规格名称" prop="name">
          <el-select style="width: 320px" v-model="dataInfo.name" placeholder="请选择规格名称">
            <el-option v-for="item in tagNameMap"
                       :label="item.name"
                       :value="item.tag"
                       @click="handleNameClick(item)"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="showValue" label="规格内容" prop="type">
          <el-select style="width: 320px" v-model="dataInfo.value" clearable multiple placeholder="请输入">

            <!-- 新增按钮作为特殊选项 -->
            <el-option disabled label="新增" value="__ADD_NEW__">
              <div>
                <div class="option-content">
                  <div v-if="isAdding" @click="AddValue" type="text"
                       style="width: 100%;background-color:lightgreen; height:40px;line-height: 40px;text-align: center;color: #4A9FF9">新增
                  </div>
                  <template v-else>
                    <el-input
                        v-model="valueName.value"
                        placeholder="请输入规格的值"
                    >
                    </el-input>
                    <el-select v-model="valueUnit" placeholder="单位" style="width: 10%">
                      <el-option v-for="(item,index) in valueUnitList[dataInfo.tag]" :key="index" :label="item" :value="item" />
                    </el-select>
                    <el-button type="primary" @click="ConfirmValue">
                      确定
                    </el-button>
                    <el-button @click="CancelValue">取消</el-button>
                  </template>
                </div>
              </div>
            </el-option>
            <el-option v-for="(item,index) in valueList"
                       :key="index"
                       :value="item.value">
              <template #default>
                <div class="option-content">
                  <span class="span-style">{{ item.value }}</span>
                  <el-button :disabled="isSelected(item,valueList)" @click.stop="DeleteValue(item)" type="text" class="delete-button"> 删除 </el-button>
                </div>
              </template>
            </el-option>

          </el-select>
        </el-form-item>

      </el-form>
      <div class="flex" style="justify-content: right;margin-top: 24px">
        <el-button type="info" @click="CloseDialog">取 消</el-button>

        <el-button type="primary" @click="submitDialog">确 定</el-button>
      </div>
    </el-drawer>

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

<script>
import {
  addConfigInfo,
  deleteConfigInfo,
  getConfigInfoList,
  updateConfigInfo,
} from "@/api/yunguan/product/productConfig";
import {addCategoryInfo, deleteCategoryInfo, getCategoryInfoList} from "@/api/yunguan/product/productCategory";

export default {
  name: "ComplexTable",
  data() {
    return {
      isAdding: true,
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      IsCreated: true,
      total: 0,
      configList: [],
      dialogVisible: false,
      valueName:{
        key:'',
        value:'',
      },
      showValue:false,
      dataInfo: {
        id: 0,
        name: '',
        tag: '',
        type: [],
        value: [],
      },
      tagList: ["CPU", "GPU", "Memory", "BandWith", "Disk", "OS"],
      tagNameMap: [
        {tag: 'cpu', name: 'CPU'},
        {tag: 'memory', name: '内存'},
        {tag: 'cpuType', name: 'CPU型号'},
        {tag: 'gpuType', name: 'GPU型号'},
        {tag: 'gpuNumber', name: 'GPU数量'},
        {tag: 'gpuMemory', name: 'GPU显存'},
        {tag: 'sysDisk', name: '系统盘'},
        {tag: 'dataDisk', name: '数据盘'},
        {tag: 'bandwidth', name: '带宽'},
        // {tag: 'os', name: '操作系统'},
        // {tag: 'env', name: '开发环境'}
        {tag: 'ip', name: 'IP'},
      ],
      valueUnit : "",
      valueUnitList:{
        'cpu': ['核'],
        'memory': ['GB'],
        'cpuType': [],
        'gpuType': [],
        'gpuNumber': ['卡'],
        'gpuMemory': ['GB'],
        'sysDisk': ['GB','TB'],
        'dataDisk': ['GB','TB'],
        'bandwidth': ['M'],
        'ip': ['个'],
      },
      valueList:[],
      cateList:[],
      specType: ["裸金属", "云主机", "容器"],
      rules: {
        tag: [{required: true, message: "请选择标签", trigger: "blur"}],
        name: [{required: true, message: "请输入名称", trigger: "blur"}],
        type: [{required: true, message: "请输入类别", trigger: "blur"}],
        default: [{required: true, message: "请选择默认选中", trigger: "blur"}],
      },
      deleteVisible: false,
      currentProductId: null

    };
  },

  mounted() {
    this.GetConfigInfoList();
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
        resp = await deleteConfigInfo(this.currentProductId);

        // 严格判断接口成功状态
        if (resp.code === 0) {
          // 显示对应成功消息
          this.$message({
            type: 'success',
            message: '删除成功'
          });
          this.deleteVisible = false;


          // 执行成功后的操作
          await this.GetConfigInfoList(); // 刷新表格
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
    isSelected(item,list) {
      // 检查该项是否被选中
      return this.dataInfo.value.includes(item.value)||list.length==1;
    },
    ChangeTag(){
      this.valueName.key=this.dataInfo.tag
      this.valueList = this.cateList.filter(data=>data.key == this.dataInfo.tag)
      this.showValue=true;
    },
    handleNameClick(item) {
      // console.log('item: ' + JSON.stringify(item))
      this.dataInfo.value = []
      this.dataInfo.name = item.name
      this.dataInfo.tag = item.tag
      this.valueName.key = this.dataInfo.tag
      this.valueList = this.cateList.filter(data => data.key === this.dataInfo.tag)
      this.showValue=true;
      this.valueUnit  = this.valueUnitList[item.tag][0]
    },
    AddValue() {
      this.isAdding = false;
    },
    ConfirmValue(){
      // 判断是否是有效的值
      // CpuType 和 GpuType 不需要判断
      if (this.valueName.key!=='cpuType'&&this.valueName.key!=='gpuType'){
        const value = this.valueName.value;
        // 先判断是否只包含数字
        if (/^\d+$/.test(value)) {
          const numValue = parseInt(value, 10);

          // 再判断是否在 0 ~ 9999 范围内
          if (numValue >= 0 && numValue <= 9999) {
            // 符合条件：仅数字，并且在 0~9999 之间
          } else {
            this.$message({
              message: '请输入 0 到 9999 之间的数字',
              type: 'warning'
            });
            return ;
          }
        } else {
          this.$message({
            message: '请输入有效的数字（仅限数字）',
            type: 'warning'
          });
          return ;
        }
      }
      // 追加单位
      if (this.valueUnit !== '' && this.valueUnit) {
        this.valueName.value = this.valueName.value + this.valueUnit
      }
      let v = this.valueList.find(data => data.value === this.valueName.value);
      if (v !== undefined) {
        this.$message({
          message: '该值已存在',
          type: 'warning',
        })
        return
      }
      addCategoryInfo(this.valueName).then(res => {
        if (res.code == 0) {
          this.valueName.value = ''
          this.isAdding = true
          this.$message({
            message: '创建成功',
            type: 'success'
          })
          this.GetCategoryList()
        }
      })
    },
    CancelValue(){
      this.isAdding=true
      this.valueName={
        key:'',
        value:'',
      }
    },
    DeleteValue(val){
      this.$confirm("此操作将永久删除该规格, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        deleteCategoryInfo(val).then(res=>{
          if (res.code == 0){
            this.$message({
              message:'成功删除',
              type:'success'
            })
            this.GetCategoryList()
            this.valueList=this.cateList.filter(data=>data.key == this.dataInfo.tag)
          }
        })
      }).catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除",
        });
      });
    },
    GetTypeName(val) {
      let name = "";
      if (val.length == 4) return "全部";
      if (val.length == 1) return this.specType[val[0]-1]
      for (let i = 0; i < val.length - 1; i++) {
        console.log(name)
        name += this.specType[val[i] - 1] + ","
      }
      name += this.specType[val[val.length - 1] - 1]
      return name
    },
    OpenDialog() {
      this.IsCreated = true
      this.showValue = false
      this.dialogVisible = true;
      this.dataInfo = {}
      this.dataInfo.value = []
      this.GetCategoryList()
    },
    CloseDialog() {
      this.dialogVisible = false;
      this.GetConfigInfoList()
    },
    handleEdit(val) {
      this.dataInfo = val
      this.showValue = true
      this.valueName.key = val.tag
      this.IsCreated = false
      this.dialogVisible = true
      this.GetCategoryList()
    },
    // 获取缓存列表
    GetConfigInfoList() {
      getConfigInfoList(this.searchInfo).then((response) => {
        if (response.code == 0) {
          this.configList = response.data.list;
          this.total = response.data.total;
        } else {
          this.$message({
            message: "获取失败",
            type: "error",
          });
        }
      });
    },
    // 删除缓存
    handleDelete(val) {
      this.$confirm("此操作将永久删除该规格, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        deleteConfigInfo(val.ID).then((response) => {
          if (response.code == 0) {
            this.$message({
              message: "删除成功",
              type: "success",
            });
            this.GetConfigInfoList();
          } else {
            this.$message({
              message: "删除失败",
              type: "error",
            });
          }
        });
      })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消删除",
            });
          });
    },
    // 添加缓存
    AddConfigInfo() {
      this.$refs.dataInfo.validate((valid) => {
        if (valid) {
          let tempData =this.configList.filter(item=>item.tag==this.dataInfo.tag)
          if (tempData.length>0&&tempData[0].type.some(item => this.dataInfo.type.includes(item))){
            this.$message({
              message:'该配置已存在!',
              type:'error'
            })
            return
          }
          addConfigInfo(this.dataInfo).then((response) => {
            if (response.code == 0) {
              this.$message({
                message: "添加成功",
                type: "success",
              });
              this.dialogVisible = false;
              this.GetConfigInfoList();
            }
          });
        }
      });
    },
    // 修改缓存
    UpdateConfigInfo() {
      // let tempData =this.configList.filter(item=>item.tag==this.dataInfo.tag)
      // if (tempData.length>0&&tempData[0].type.some(item => this.dataInfo.type.includes(item))){
      //   this.$message({
      //     message:'该配置已存在!',
      //     type:'error'
      //   })
      //   return
      // }
      updateConfigInfo(this.dataInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.dialogVisible = false
          this.GetConfigInfoList();
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    submitDialog() {
      if (this.IsCreated) this.AddConfigInfo()
      else this.UpdateConfigInfo()
    },
    //分页查询
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetConfigInfoList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetConfigInfoList()
    },

    //获取字典列表
    GetCategoryList(){
      getCategoryInfoList().then(res=>{
        if (res.code == 0){
          // todo 当cateList为空时，由于AddValue按钮所在select下，所以不会显示AddValue按钮，无法添加新的值
          this.cateList = res.data.list
          this.valueList=this.cateList.filter(data=>data.key == this.dataInfo.tag)
        }
      })
    }
  },
};
</script>
<style scoped>
/* 确保选项内容水平排列 */
.option-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

/* 设置左侧文本样式 */
.span-style {
  float: left;
  margin-right: 8px; /* 右边距 */
}

/* 设置删除按钮样式 */
.delete-button {
  float: right;
}

/* 当鼠标悬停时改变背景色 */
.el-select-dropdown__item:hover .option-content {
  background-color: #f5f7fa;
}


</style>
