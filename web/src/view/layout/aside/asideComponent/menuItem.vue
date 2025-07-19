<template>
  <el-menu-item
    :index="routerInfo.name"
    class="no-select dark:text-slate-300 overflow-x-hidden"
    :style="{
      height : sideHeight,
    }"
  >
    <el-icon v-if="routerInfo.meta.icon">
      <component :is="routerInfo.meta.icon" />
    </el-icon>
    <template #title>
      {{ routerInfo.meta.title }}
    </template>
  </el-menu-item>
</template>


<script setup>
import { computed } from 'vue'
import { useAppStore } from '@/pinia'
import { storeToRefs } from 'pinia'
const appStore = useAppStore()
const { config } = storeToRefs(appStore)

defineOptions({
  name: 'MenuItem',
})

defineProps({
  routerInfo: {
    default: function() {
      return null
    },
    type: Object
  },
})

const sideHeight = computed(() => {
  return  config.value.layout_side_item_height + 'px'
})

</script>

<style lang="scss">

</style>
<style scoped>
.no-select {
  user-select: none; /* 禁止文本选择 */
  -webkit-user-select: none; /* Safari */
  -moz-user-select: none; /* Firefox */
  -ms-user-select: none; /* IE 10+ */
}

/* 如果需要更广泛的禁用选择，可以考虑全局样式 */
</style>