<template>
  <div style="display: flex; justify-content: center; align-items: center; height: 400px;min-height: 100vh; min-width: 100vw;">
    <div class="gva-form-box" v-if="showLicenseConfig" style="width: 40%; min-width: 200px; padding: 20px; box-sizing: border-box;">
      <h4>License配置</h4>
      <el-descriptions border label-width="100px" :column="1">
        <el-descriptions-item label="设备数量">{{ licenseConfig.deviceNum }}</el-descriptions-item>
        <el-descriptions-item label="租户数量">{{ licenseConfig.renterNum }}</el-descriptions-item>
        <el-descriptions-item label="管理员数量">{{ licenseConfig.adminNum }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ formatTimestamp(licenseConfig.startTime) }}</el-descriptions-item>
        <el-descriptions-item label="结束时间">{{ formatTimestamp(licenseConfig.endTime) }}</el-descriptions-item>
      </el-descriptions>
      <div style="margin-top: 10px;display: flex">
        <el-button style="flex: 1" @click="OpenModPwdDialog" type="primary">修改License密码</el-button>
        <el-button style="flex: 1" @click="OpenModConfigDialog" type="primary">修改License配置</el-button>
        <el-button style="flex: 1" @click="goBack" type="primary">返回</el-button>
      </div>
    </div>
  </div>
  <!--    修改密码-->
  <el-dialog v-model="pwdVisible" :width="300" :show-close="true" :title="'修改密码'">
    <el-form :model="modPwdInfo" :rules="rules" ref="modPwdInfo" label-width="90px"
             label-position="left" @keyup.enter="modPwdConfirm">
      <el-form-item label="新密码" prop="licensePwd">
        <el-input placeholder="请输入新密码" v-model="modPwdInfo.licensePwd"></el-input>
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPwd">
        <el-input placeholder="请再次输入密码" v-model="modPwdInfo.confirmPwd"></el-input>
      </el-form-item>
      <el-button style="width: 100%" type="primary" @click="modPwdConfirm">确认</el-button>
    </el-form>
  </el-dialog>
  <!--    登录-->
  <el-dialog v-model="loginLicenseVisible" :width="300" :before-close="closeDialog" :title="'License配置登录'">
    <el-form ref="licenseConfigInfo" :model="licenseConfigInfo" :rules="rules" label-width="70px"
             label-position="left"
             @keyup.enter="loginLicenseConfig">
      <el-form-item label="用户名" prop="licenseUser">
        <el-input placeholder="请输入用户名" v-model="licenseConfigInfo.licenseUser"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="licensePwd">
        <el-input placeholder="请输入密码" v-model="licenseConfigInfo.licensePwd"></el-input>
      </el-form-item>
      <el-button style="width: 100%" type="primary" @click="loginLicenseConfig">确认</el-button>
    </el-form>
  </el-dialog>
  <!--    修改配置-->
  <el-dialog v-model="configVisible" :show-close="true" :title="'修改配置'">
    <el-form :model="modLicenseInfo" label-width="100px" label-position="left" @keyup.enter="modConfigConfirm">
      <el-form-item label="开始时间">
        <el-date-picker
            style="width: 100%"
            v-model="modLicenseInfo.startTime"
            type="datetime"
            placeholder="选择开始时间"
            :disabled-date="disabledDate"
        />
      </el-form-item>
      <el-form-item label="结束时间">
        <el-date-picker
            style="width: 100%"
            v-model="modLicenseInfo.endTime"
            type="datetime"
            placeholder="选择结束时间"
            :disabled-date="disabledDate"
        />
      </el-form-item>
      <el-form-item label="设备数量">
        <el-input-number style="width: 100%" controls-position="right" :precision="0" :min="0"
                         v-model="modLicenseInfo.deviceNum"/>
      </el-form-item>
      <el-form-item label="租户数量">
        <el-input-number style="width: 100%" controls-position="right" :precision="0" :min="0"
                         v-model="modLicenseInfo.renterNum"/>
      </el-form-item>
      <el-form-item label="管理员数量">
        <el-input-number style="width: 100%" controls-position="right" :precision="0" :min="0"
                         v-model="modLicenseInfo.adminNum"/>
      </el-form-item>
      <el-button style="width: 100%" type="primary" @click="modConfigConfirm">确认</el-button>
    </el-form>
  </el-dialog>
</template>
<script>

import {loginLicense, queryLicense, updateLicense, updateLicensePwd} from "@/api/yunguan/config/licenseConfig";
import _ from "lodash";

export default {
  data() {
    return {
      licenseConfig: {},
      showLicenseConfig: false,
      loginLicenseVisible: false,
      pwdVisible: false,
      configVisible: false,
      rules: {
        licenseUser: [{required: true, message: "请输入用户名", trigger: "blur"}],
        licensePwd: [
          {required: true, message: "请输入密码", trigger: "blur"},
          {
            pattern: /^.{6,20}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '长度应在6到20个字符之间',
            trigger: 'blur'
          }],
        confirmPwd: [
          {required: true, message: "请输入密码", trigger: "blur"},
          {
            pattern: /^.{6,20}$/, // 正则表达式，匹配以字母开头且仅包含字母和数字的字符串，长度为2到12
            message: '长度应在6到20个字符之间',
            trigger: 'blur'
          }],
      },
      licenseConfigInfo: {
        licenseUser: '',
        licensePwd: '',
      },
      modLicenseInfo: {
        id: '',
        deviceNum: 0,
        renterNum: 0,
        adminNum: 0,
        startTime: 0,
        endTime: 0,
      },
      tempLicenseInfo:{},
      modPwdInfo: {
        id: '',
        licensePwd: '',
        confirmPwd: '',
      },
    }
  },
  created() {
    this.OpenLoginConfig();
  },
  methods: {
    goBack(){
      this.$router.push({name: "Login"});
    },
    disabledDate(time) {
      return time.getTime() < Date.now() - 8.64e7;
    },
    modPwdConfirm() {
      this.$refs.modPwdInfo.validate(valid => {
        if (valid) {
          if (this.modPwdInfo.licensePwd != this.modPwdInfo.confirmPwd) {
            this.$message({
              type: 'warning',
              message: '两次密码不一致',
            })
            return
          }
          updateLicensePwd(this.modPwdInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '修改成功，请重新登录',
              })
              this.pwdVisible = false;
              this.OpenLoginConfig();
            }
          })
        }
      })
    },
    modConfigConfirm() {
      this.tempLicenseInfo=_.cloneDeep(this.modLicenseInfo)
      this.tempLicenseInfo.startTime/=1000;
      this.tempLicenseInfo.endTime/=1000;
      updateLicense(this.tempLicenseInfo).then(res => {
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '修改成功',
          })
          this.configVisible = false;
          this.GetLicenseConfig();
        }
      })
    },
    closeDialog() {
      this.loginLicenseVisible = false;
      this.$router.push({name: "Login"});
    },
    OpenLoginConfig() {
      this.showLicenseConfig = false;
      this.licenseConfigInfo = {};
      this.loginLicenseVisible = true;
    },
    OpenModPwdDialog() {
      this.modPwdInfo = {};
      this.modPwdInfo.id = _.cloneDeep(this.licenseConfig.id);
      this.pwdVisible = true;
    },
    OpenModConfigDialog() {
      this.modLicenseInfo.id = _.cloneDeep(this.licenseConfig.id);
      this.modLicenseInfo.renterNum = _.cloneDeep(this.licenseConfig.renterNum);
      this.modLicenseInfo.deviceNum = _.cloneDeep(this.licenseConfig.deviceNum);
      this.modLicenseInfo.adminNum = _.cloneDeep(this.licenseConfig.adminNum);
      this.modLicenseInfo.endTime = _.cloneDeep(this.licenseConfig.endTime*1000);
      this.modLicenseInfo.startTime = _.cloneDeep(this.licenseConfig.startTime*1000);
      this.configVisible = true;
    },
    loginLicenseConfig() {
      this.$refs.licenseConfigInfo.validate(valid => {
        if (valid) {
          loginLicense(this.licenseConfigInfo).then(res => {
            if (res.code == 0) {
              this.$message({
                type: 'success',
                message: '登录成功'
              })
              this.loginLicenseVisible = false;
              this.showLicenseConfig = true;
              this.GetLicenseConfig();
            }
          })
        }
      })
    },
    formatTimestamp(timestamp) {
      var tempTime = timestamp * 1000;
      const date = new Date(tempTime);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从 0 开始，需加 1
      const day = String(date.getDate()).padStart(2, '0');

      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');

      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    GetLicenseConfig() {
      queryLicense().then(res => {
        if (res.code == 0) {
          this.licenseConfig = res.data;
        }
      })
    }
  }
}

</script>

<style scoped lang="scss">

</style>