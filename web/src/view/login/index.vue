<template>
  <div
    id="userLayout"
    class="w-full h-full relative"
  >
    <div
      class="rounded-lg flex items-center justify-evenly w-full h-full md:w-screen md:h-screen md:bg-[#194bfb] bg-white"
    >
      <div class="md:w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div class="oblique h-[130%] w-3/5 bg-white dark:bg-slate-900 transform -rotate-12 absolute -ml-52" />
        <!-- 分割斜块 -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">

              <img
                class="w-24"
                :src="$GIN_VUE_ADMIN.appLogo"
                alt
              >
            </div>
            <div class="mb-9">
              <p class="text-center text-4xl font-bold">{{ $GIN_VUE_ADMIN.appName }}</p>
<!--              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">A management platform using Golang and Vue</p>-->
            </div>
            <el-form
              ref="loginForm"
              :model="loginFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              @keyup.enter="submitForm"
              label-position="left"
            >
              <el-form-item
                prop="username"
                class="mb-6"
              >
                <el-input
                  v-model="loginFormData.username"
                  size="large"
                  placeholder="请输入用户名"
                  suffix-icon="user"
                />
              </el-form-item>
              <el-form-item
                prop="password"
                class="mb-6"
              >
                <el-input
                  v-model="loginFormData.password"
                  show-password
                  size="large"
                  type="password"
                  placeholder="请输入密码"
                />
              </el-form-item>
              <el-form-item
                v-if="loginFormData.openCaptcha"
                prop="captcha"
                class="mb-6"
              >
                <div class="flex w-full justify-between">
                  <el-input
                    v-model="loginFormData.captcha"
                    placeholder="请输入验证码"
                    size="large"
                    class="flex-1 mr-5"
                  />
                  <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
                    <img
                      v-if="picPath"
                      class="w-full h-full"
                      :src="picPath"
                      alt="请输入验证码"
                      @click="loginVerify()"
                    >
                  </div>
                </div>
              </el-form-item>
              <el-form-item class="mb-6">
                <span style="margin: 0 auto;margin-top: 0;font-size: 12px;color: #8c939d" v-if="isShowText!=''">{{isShowText}}</span>
                <el-button
                  class="shadow shadow-active h-11 w-full"
                  type="primary"
                  size="large"
                  :disabled="disableLogin"
                  @click="submitForm"
                >登 录</el-button>
              </el-form-item>
<!--              <el-form-item class="mb-6">-->
<!--                <el-button-->
<!--                  class="shadow shadow-active h-11 w-full"-->
<!--                  type="primary"-->
<!--                  size="large"-->
<!--                  @click="checkInit"-->
<!--                >前往初始化</el-button>-->

<!--              </el-form-item>-->
            </el-form>
          </div>
          <el-button @click="toggleLicenseConfig" v-if="clickTime>10">前往License配置</el-button>
        </div>
      </div>
      <div class="hidden md:block w-1/2 h-full float-right bg-[#194bfb]"><img
        class="h-full"
        src="@/assets/login_background.svg"
        alt="banner"
      >
      </div>
    </div>

    <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto  w-full z-20">
      <div class="links items-center justify-center gap-2 hidden md:flex">
        <a
          href="https://www.gin-vue-admin.com/"
          target="_blank"
        >
          <img
            src="@/assets/docs.png"
            class="w-8 h-8"
            alt="文档"
          >
        </a>
        <a
          href="https://support.qq.com/product/371961"
          target="_blank"
        >
          <img
            src="@/assets/kefu.png"
            class="w-8 h-8"
            alt="客服"
          >
        </a>
        <a
          href="https://ygpt"
          target="_blank"
        >
          <img
            src="@/assets/github.png"
            class="w-8 h-8"
            alt="github"
          >
        </a>
        <a
          href="https://space.bilibili.com/322210472"
          target="_blank"
        >
          <img
            src="@/assets/video.png"
            class="w-8 h-8"
            alt="视频站"
          >
        </a>
      </div>
    </BottomInfo>
  </div>
</template>

<script setup>
import { captcha } from '@/api/system/user'
import { checkDB } from '@/api/system/initdb'
import BottomInfo from '@/components/bottomInfo/bottomInfo.vue'
import {onMounted, reactive, ref} from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import {queryLicense} from "@/api/yunguan/config/licenseConfig";

defineOptions({
  name: "Login",
})

onMounted(()=>{
  loginVerify();
  GetLicenseConfig();
})
const router = useRouter()
// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}
const clickTime=ref(0)
// 获取验证码
const loginVerify = async() => {
  clickTime.value++;
  const ele = await captcha()
  rules.captcha.push({
    max: ele.data.captchaLength,
    min: ele.data.captchaLength,
    message: `请输入${ele.data.captchaLength}位验证码`,
    trigger: 'blur',
  })
  picPath.value = ele.data.picPath
  loginFormData.captchaId = ele.data.captchaId
  loginFormData.openCaptcha = ele.data.openCaptcha
}
// 登录相关操作
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: 'admin',
  password: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const login = async() => {
  return await userStore.LoginIn(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async(v) => {
    if (!v) {
      // 未通过前端静态验证
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      await loginVerify()
      return false
    }

    // 通过验证，请求登陆
    const flag = await login()

    // 登陆失败，刷新验证码
    if (!flag) {
      await loginVerify()
      return false
    }

    // 登陆成功
    return true
  })
}
const licenseConfig=ref({})
const GetLicenseConfig=()=> {
  queryLicense().then(res => {
    if (res.code == 0) {
      licenseConfig.value = res.data;
      CheckLicenseTime();
    }
  })
}
const isShowText=ref('');
const disableLogin=ref(false);
const CheckLicenseTime=()=>{
  const unixTimestamp = Math.floor(Date.now() / 1000);
  const aMonthTime=2592000;
  if (unixTimestamp<licenseConfig.value.startTime){
    isShowText.value="当前License将从"+formatTimestamp(licenseConfig.value.startTime)+"起开始生效";
    disableLogin.value=true;
    return;
  }
  if (unixTimestamp>licenseConfig.value.endTime){
    isShowText.value="当前License在"+formatTimestamp(licenseConfig.value.endTime)+"已过期";
    disableLogin.value=true;
    return;
  }
  if ((unixTimestamp+aMonthTime)>=licenseConfig.value.endTime){
    isShowText.value="当前License将在"+formatTimestamp(licenseConfig.value.endTime)+"过期";
    disableLogin.value=false;
    return;
  }
  isShowText.value="";
  return;
}

const toggleLicenseConfig = async() =>{
  userStore.SetLicenseConfig()
  await router.push({name: 'LicenseConfig'})
}
// 跳转初始化
const checkInit = async() => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      await router.push({name: 'Init'})
    } else {
      ElMessage({
        type: 'info',
        message: '已配置数据库信息，无法初始化',
      })
    }
  }
}

const formatTimestamp=(timestamp)=> {
  var tempTime = timestamp * 1000;
  const date = new Date(tempTime);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从 0 开始，需加 1
  const day = String(date.getDate()).padStart(2, '0');

  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
</script>
