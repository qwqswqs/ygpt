<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <!-- 搜索和重置功能 -->
        <el-input
            v-model="searchQuery"
            class="filter-item"
            placeholder="请输入模型名称"
            style="width: 200px;"
        />

        <el-button class="filter-item"  type="primary" @click="handleSearch">
          搜索
        </el-button>
        <el-button class="filter-item" style="margin-left: 10px;" type="primary" @click="OpenDialog">
          添加
        </el-button>
      </div>

    <el-table
      v-loading="false"
      :cell-style="{'text-align':'left'}"
      :data="list"
      :header-cell-style="{ 'background-color': 'rgba(242, 243, 245, 1)','color': 'rgba(29, 33, 41, 1)','text-align':'left'}"
      style="width: 100%;font-size: 15px;"
    >
      <el-table-column label="序号">
        <template #default="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="名称">
        <template #default="scope">
          <span>{{ scope.row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="类型">
        <template #default="scope">
          <span>{{ scope.row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用途">
        <template #default="scope">
          <span>{{ scope.row.purpose }}</span>
        </template>
      </el-table-column>
      <el-table-column label="使用次数">
        <template #default="scope">
          <span>{{ scope.row.useTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="发布状态">
        <template #default="scope">
          <span>{{ scope.row.share }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
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
    <el-dialog
      :title="'创建数据集'"
      v-model="dialogVisible"
      width="40%"
    >
      <el-form
        ref="modelData"
        :model="data"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="data.name" placeholder="请输入数据集名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="data.type" placeholder="请选择数据集类型" style="width: 100%">
            <el-option label="文本" />
            <el-option label="图片" />
          </el-select>
        </el-form-item>
        <el-form-item label="用途" prop="desc">
          <el-input v-model="data.desc" placeholder="请描述数据集用途" />
        </el-form-item>
        <el-form-item label="是否发布" prop="share">
          <el-radio-group v-model="data.share">
            <el-radio :label="'是'" />
            <el-radio :label="'否'" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="文件上传">
          <el-upload
            class="upload-dragger"
            action="https://jsonplaceholder.typicode.com/posts/"
            multiple
            :limit="3"
          >
            <i class="el-icon-upload"></i>
            <span class="el-upload__text">将文件拖到此处，或<em>点击上传</em></span>
          </el-upload>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer" style="margin: 0 auto">
        <el-button type="primary">确 定</el-button>
        <el-button @click="CloseDialog">取 消</el-button>
      </div>
    </el-dialog>
  </div>
  </div>
</template>

<script>
export default {
  name: 'ComplexTable',
  data() {
    return {
      searchInfo: {
        pageIndex:1,
        pageSize:10,
      },
      total: 0,
      list: [
        { id: 1, name: '语言大模型', type: '文本', purpose: '各个国家的语言', useTime: 124, share: '已发布' },
        { id: 2, name: '语音大模型', type: '语音', purpose: '男女老少的语言', useTime: 34, share: '已发布' },
        { id: 3, name: '图片大模型', type: '图片', purpose: '各种小猫的图片', useTime: 245, share: '已发布' }
      ],
      dialogVisible: false,
      data: {
        name: '',
        type: '',
        desc: '',
        price: '',
        share: '是'
      },
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择类型', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    OpenDialog() {
      this.dialogVisible = true
    },
    CloseDialog() {
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.filter-container {
  margin-bottom: 20px;
}

.filter-item {
  margin-right: 10px;
}
</style>

