<template>
  <div class="app-container">
    <el-card>
      <div id="remote-header" slot="header" :style="{color:isConnected?'#67C23A':'#E6A23C'}">
        {{ isConnected ?'成功连接算力平台':'待连接算力平台' }}
      </div>
      <el-form>
        <el-form-item label="算力平台地址">
          <el-input v-model="remoteAddr" :disabled="isConnected" placeholder="请输入云管平台地址" @input="checkInput">
            <template slot="prepend">Http://</template>
          </el-input>
          <i v-show="!isOkAddr" style="color: #F56C6C;">请输入合法的地址</i>
        </el-form-item>
        <el-form-item label="云管平台id">
          <el-input v-model="platformId" :disabled="isConnected" placeholder="请输入云管平台的id" @input="checkInput" />
          <i v-show="!isOkId" style="color: #F56C6C;">请输入合法的id</i>
        </el-form-item>
        <el-form-item label="远程连接授权码">
          <el-input v-model="authCode" type="password" :disabled="isConnected" placeholder="请输入远程连接的授权码" @input="checkInput" />
          <i v-show="!isOkAuthCode" style="color: #F56C6C;">请输入正确的远程授权码</i>
        </el-form-item>
        <el-form-item>
          <el-button v-show="!isConnected" :disabled="!isOKSubmit" type="primary" @click="isConnected=true">连接</el-button>
          <el-button v-show="isConnected" type="primary" @click="isConnected=false">断开</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'ComplexTable',
  data() {
    return {
      isConnected: false,
      remoteAddr: '',
      platformId: '',
      authCode: '',
      isOkAddr: true,
      isOkId: true,
      isOkAuthCode: true,
      isOKSubmit: false
    }
  },
  methods: {
    checkInput() {
      let isOk = true
      const urlPattern = /^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\*\+,;=.]+$/
      const idPattern = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{16}$/
      const authCodePattern = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{16}$/
      if (!urlPattern.test(this.remoteAddr)) {
        this.isOkAddr = false
        isOk = false
      } else {
        this.isOkAddr = true
      }
      if (!idPattern.test(this.platformId)) {
        this.isOkId = false
        isOk = false
      } else {
        this.isOkId = true
      }
      if (!authCodePattern.test(this.authCode)) {
        this.isOkAuthCode = false
        isOk = false
      } else {
        this.isOkAuthCode = true
      }
      this.isOKSubmit = isOk
      return isOk
    }
  }
}
</script>

<style scoped>
  #remote-header{
    font-size: 30px;
    font-family:Arial, Helvetica, sans-serif;
  }
</style>
