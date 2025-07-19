<script setup>
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/pinia/modules/user';
import {ElMessage} from "element-plus";
import {useRouterStore} from "@/pinia/modules/router";
import router from '@/router'

function getHashParams() {
  const hash = window.location.hash
  if (!hash) return {}
  const hashPart = hash.replace(/^#\/?/, '')
  const [path, queryString] = hashPart.split('?')
  return queryString ? new URLSearchParams(queryString) : {}
}

function handleLoginAndRedirect() {
  try {
    const token = getHashParams().get('token');
    if (!token) throw new Error()
    userStore.setToken(token)
    window.location.href = window.location.pathname
  } catch (e) {
    ElMessage.error({
      message: '请求未携带正确的token',
      type: 'error',
    });
    router.push({ name: 'Login' });
  }
}

const userStore = useUserStore()
const routerStore = useRouterStore()

onMounted(() => {
  handleLoginAndRedirect()
})
</script>
