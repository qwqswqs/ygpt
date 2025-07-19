<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
            class="filter-item"
            placeholder="请输入缓存名称"
            style="width: 200px"
        />

        <el-button
            v-waves
            class="filter-item"
            type="primary"
        >
          搜索
        </el-button>
        <el-button
            class="filter-item"
            style="margin-left: 10px"
            type="primary"
            @click="OpenDialog"
        >
          添加
        </el-button>
      </div>

      <el-table
          v-loading="false"
          :cell-style="{ 'text-align': 'left' }"
          :data="list"
          :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
          style="width: 100%; font-size: 15px"
      >
        <el-table-column label="编号">
          <template #default="scope">
            <span>{{ scope.row.num }}</span>
          </template>
        </el-table-column>
        <el-table-column label="名称">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ scope.row.type }}</span>
          </template>
        </el-table-column>
        <el-table-column label="说明">
          <template #default="scope">
            <span>{{ scope.row.description }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right">
          <template #default="scope">
            <el-button type="text" @click="handleEdit(scope.row)" >修改</el-button>
            <el-button type="text" @click="handleDelete(scope.row)">删除</el-button>
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
               size="40%">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ IsCreated ? '新增' : '编辑' }}系统类别</span>
          <div>
            <el-button @click="CloseDialog">取 消</el-button>
            <el-button type="primary" @click="submitDialog">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form
          ref="addInfo"
          :model="addInfo"
          label-width="60px"
          :rules="IsCreated?rules:''"
      >
        <el-form-item label="编号" prop="num">
          <el-input-number style="width:100%" v-model="addInfo.num" :min="1"  />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="addInfo.name" />
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input v-model="addInfo.description" />
        </el-form-item>
        <el-form-item label="类别" prop="type">
          <div v-for="(item,index) in addInfo.type" :key="index" style="width: 100%;">
            <el-input v-model="addInfo.type[index]" style="width: 100%" @input="updateType(index, $event)"></el-input>
          </div>
        </el-form-item>
        <el-button type="text" @click="AddNewType">添加新类别</el-button>
        <el-button v-if="addInfo.type.length != 0" type="text" @click="RemoveNewType">移除类别</el-button>
      </el-form>
    </el-drawer>
  </div>
</template>

<script>
import {
  addSysType,
  updateSysType,
  deleteSysType,
  getSysTypeList,
} from "@/api/yunguan/system/type";

export default {
  name: "ComplexTable",
  data() {
    return {
      searchInfo: {
        pageIndex: 1,
        pageSize: 10,
      },
      IsCreated: true,
      total: 0,
      list: [],
      dialogVisible: false,
      addInfo: {
        num:0,
        name:'',
        type:[""],
        description:'',
        id:null,
      },
      typeList:[""],
      rules: {
        name: [{ required: true, message: "请输入名称", trigger: "blur" }],
        num: [{ required: true, message: "请输入编号", trigger: "blur" }],
        type: [{ required: true, message: "请输入类别", trigger: "blur" }],
        description: [{ required: true, message: "请输入说明", trigger: "blur" }],
      },
    };
  },

  mounted() {
    this.GetSysTypeList();
  },
  methods: {

    AddNewType () {
      this.addInfo.type.push("")
    },
    RemoveNewType () {
      this.addInfo.type.pop()
    },
    updateType(index, event) {
      // 可选：用于处理输入事件，确保类型安全等
      const newValue = event.target.value;
      this.$set(this.addInfo.type, index, newValue); // 使用 Vue 的 $set 方法来保证响应性
    },
    OpenDialog() {
      this.dialogVisible = true;
    },
    CloseDialog() {
      this.dialogVisible = false;
    },
    handleEdit(val) {
      this.addInfo.name=val.name
      this.addInfo.description=val.description
      this.addInfo.id=val.ID
      this.addInfo.num=val.num
      this.addInfo.type=val.type
      this.IsCreated = false
      this.dialogVisible = true
    },
    // 获取缓存列表
    GetSysTypeList() {
      getSysTypeList(this.searchInfo).then((response) => {
        if (response.code == 0) {
          this.list = response.data.list;
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
      this.$confirm("此操作将永久删除该缓存, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
          .then(() => {
            deleteSysType(val.ID).then((response) => {
              if (response.code == 0) {
                this.$message({
                  message: "删除成功",
                  type: "success",
                });
                this.GetSysTypeList();
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
    AddSysType() {
      this.$refs.addInfo.validate((valid) => {
        if (valid) {
          addSysType(this.addInfo).then((response) => {
            if (response.code == 0) {
              this.$message({
                message: "添加成功",
                type: "success",
              });
              this.dialogVisible = false;
              this.GetSysTypeList();
            } else {
              this.$message({
                message: "添加失败",
                type: "error",
              });
            }
          });
        }
      });
    },

    // 修改缓存
    UpdateSysType() {
      updateSysType(this.addInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.dialogVisible = false
          this.GetSysTypeList();
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    submitDialog() {
      if (this.IsCreated)this.AddSysType()
      else this.UpdateSysType()
    },
  },
};
</script>
