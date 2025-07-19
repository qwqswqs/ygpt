

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
        <el-table-column label="序号">
          <template #default="scope">
            <span>{{ scope.row.ID }}</span>
          </template>
        </el-table-column>
        <el-table-column label="软件名称" width="120px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="类别">
          <template #default="scope">
            <span>{{ scope.row.type }}</span>
          </template>
        </el-table-column>
        <el-table-column label="版本">
          <template #default="scope">
            <span>{{ scope.row.version }}</span>
          </template>
        </el-table-column>
        <el-table-column label="url" width="200px">
          <template #default="scope">
            <span style="display:block;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">{{ scope.row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小">
          <template #default="scope">
            <span>{{ scope.row.size }}</span>
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
            :current-page="searchInfo.page"
            :page-size="searchInfo.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total,prev, pager, next,  sizes, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="'添加缓存信息'" width="40%">
      <el-form
          ref="addInfo"
          :model="addInfo"
          label-width="120px"
          :rules="IsCreated?rules:''"
      >
        <el-form-item label="工具名称" prop="name">
          <el-input v-model="addInfo.name" />
        </el-form-item>
        <el-form-item label="所在地址" prop="address">
          <el-input v-model="addInfo.address" />
        </el-form-item>
        <el-form-item label="作用说明" prop="description">
          <el-input v-model="addInfo.description" />
        </el-form-item>
        <el-form-item label="所在服务器ID" prop="working">
          <el-input v-model.number="addInfo.working" />
        </el-form-item>
        <el-form-item label="进程号" prop="taskID">
          <el-input v-model="addInfo.taskID" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary" @click="submitDialog">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  addSysSoftware,
  updateSysSoftware,
  deleteSysSoftware,
  getSysSoftwareList,
} from "@/api/yunguan/system/software";

export default {
  name: "ComplexTable",
  data() {
    return {
      searchInfo: {
        page: 1,
        pageSize: 10,
      },
      IsCreated: true,
      total: 0,
      list: [],
      dialogVisible: false,
      addInfo: {
        name:'',
        address: '',
        description:'',
        working:'',
        taskID:'',
        id:null,
      },

      rules: {
        name: [{ required: true, message: "请输入名称", trigger: "blur" }],
        address: [{ required: true, message: "请输入地址", trigger: "blur" }],
      },
    };
  },

  created() {
    this.GetSysSoftwareList();
  },
  methods: {
    handleCurrentChange(val) {
      this.searchInfo.page = val
      this.GetSysSoftwareList()
    },
    handleSizeChange(val) {
      this.searchInfo.pageSize = val
      this.GetSysSoftwareList()
    },
    OpenDialog() {
      this.dialogVisible = true;
    },
    CloseDialog() {
      this.dialogVisible = false;
    },
    handleEdit(val) {
      this.addInfo=val
      this.IsCreated = false
      this.dialogVisible = true
    },
    // 获取缓存列表
    GetSysSoftwareList() {
      getSysSoftwareList(this.searchInfo).then((response) => {
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
      }).then(() => {
            deleteSysSoftware(val.ID).then((response) => {
              if (response.code == 0) {
                this.$message({
                  message: "删除成功",
                  type: "success",
                });
                this.GetSysSoftwareList();
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
    AddSysSoftware() {
      this.$refs.addInfo.validate((valid) => {
        if (valid) {
          addSysSoftware(this.addInfo).then((response) => {
            if (response.code == 0) {
              this.$message({
                message: "添加成功",
                type: "success",
              });
              this.dialogVisible = false;
              this.GetSysSoftwareList();
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
    UpdateSysSoftware() {
      updateSysSoftware(this.addInfo).then(response => {
        if (response.code == 0) {
          this.$message({
            message: '编辑成功',
            type: 'success'
          });
          this.dialogVisible = false
          this.GetSysSoftwareList();
        } else {
          this.$message({
            message: '添加失败',
            type: 'error'
          });
        }
      })
    },
    submitDialog() {
      if (this.IsCreated)this.AddSysSoftware()
      else this.UpdateSysSoftware()
    },
  },
};
</script>
