<template>
  <component
    :is="menuComponent"
    v-if="!routerInfo.hidden"
    :router-info="routerInfo"
    class="no-select"
  >
    <template v-if="routerInfo.children&&routerInfo.children.length">
      <AsideComponent
        v-for="item in routerInfo.children"
        :key="item.name"
        :router-info="item"
      />
    </template>
  </component>
</template>


<script setup>
import MenuItem from './menuItem.vue'
import AsyncSubmenu from './asyncSubmenu.vue'
import { computed } from 'vue'

defineOptions({
  name: 'AsideComponent',
})

const props = defineProps({
  routerInfo: {
    type: Object,
    default: () => null,
  },
  mode :{
    type : String,
    default: "vertical"
  }
})


const menuComponent = computed(() => {
  if (props.routerInfo.children && props.routerInfo.children.filter(item => !item.hidden).length) {
    return AsyncSubmenu
  } else {
    return MenuItem
  }
})

</script>

<style scoped>
.no-select {
  user-select: none; /* 禁止文本选择 */
  -webkit-user-select: none; /* Safari */
  -moz-user-select: none; /* Firefox */
  -ms-user-select: none; /* IE 10+ */
}

/* 如果需要更广泛的禁用选择，可以考虑全局样式 */
</style>