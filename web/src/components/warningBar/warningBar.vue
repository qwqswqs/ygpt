<template>
  <div
      class="flex items-center rounded-sm dark:bg-amber-700 dark:text-gray-200"
      style="background-color: #E8F3FF;color: #1D2129;margin-bottom: 16px;padding: 9px 16px;display: flex;align-items: center;"
      :class="href&&'cursor-pointer'"
      @click="open"
  >
    <img :src="warn" style="margin-right: 8px;width: 16px;height: 19px">
    <span style="font-size: 14px;font-weight: 400;">
      {{ title }}
    </span>
  </div>
</template>

<script setup>
import warn from '@/assets/warning.png';
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  href: {
    type: String,
    default: ''
  }
})

const open = () => {
  if (props.href) {
    window.open(props.href)
  }
}
</script>