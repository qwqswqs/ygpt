<template>
  <div>
    <el-row
      :gutter="15"
      class="py-1"
    >
      <el-col :span="12">
        <el-card
          v-if="state.os"
          class="card_item"
        >
          <template #header>
            <div>Runtime</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">os:</el-col>
              <el-col
                :span="12"
                v-text="state.os.goos"
              />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">cpu nums:</el-col>
              <el-col
                :span="12"
                v-text="state.os.numCpu"
              />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">compiler:</el-col>
              <el-col
                :span="12"
                v-text="state.os.compiler"
              />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">go version:</el-col>
              <el-col
                :span="12"
                v-text="state.os.goVersion"
              />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">goroutine nums:</el-col>
              <el-col
                :span="12"
                v-text="state.os.numGoroutine"
              />
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card
          v-if="state.disk"
          class="card_item"
          :body-style="{ height: '180px', 'overflow-y': 'scroll' }"
        >
          <template #header>
            <div>Disk</div>
          </template>
          <div>
            <el-row
                v-for="(item, index) in state.disk"
                :key="index"
                :gutter="10"
                style="margin-bottom: 2rem"
            >
              <el-col :span="12">

                <el-row :gutter="10">
                  <el-col :span="12">MountPoint</el-col>
                  <el-col
                      :span="12"
                      v-text="item.mountPoint"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col
                      :span="12"
                      v-text="item.totalMb"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col
                      :span="12"
                      v-text="item.usedMb"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col
                      :span="12"
                      v-text="item.totalGb"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col
                      :span="12"
                      v-text="item.usedGb"
                  />
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                    type="dashboard"
                    :percentage="item.usedPercent"
                    :color="colors"
                />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row
      :gutter="15"
      class="py-1"
    >
      <el-col :span="12">
        <el-card
          v-if="state.cpu"
          class="card_item"
          :body-style="{ height: '180px', 'overflow-y': 'scroll' }"
        >
          <template #header>
            <div>CPU</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">physical number of cores:</el-col>
              <el-col
                :span="12"
                v-text="state.cpu.cores"
              />
            </el-row>
            <el-row
              v-for="(item, index) in state.cpu.cpus"
              :key="index"
              :gutter="10"
            >
              <el-col :span="12">core {{ index }}:</el-col>
              <el-col
                :span="12"
              ><el-progress
                type="line"
                :percentage="+item.toFixed(0)"
                :color="colors"
              /></el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card
          v-if="state.ram"
          class="card_item"
        >
          <template #header>
            <div>Ram</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col
                    :span="12"
                    v-text="state.ram.totalMb"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col
                    :span="12"
                    v-text="state.ram.usedMb"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col
                    :span="12"
                    v-text="state.ram.totalMb / 1024"
                  />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col
                    :span="12"
                    v-text="(state.ram.usedMb / 1024).toFixed(2)"
                  />
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                  type="dashboard"
                  :percentage="state.ram.usedPercent"
                  :color="colors"
                />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { getSystemState } from '@/api/system/system'
import { onUnmounted, ref } from 'vue'

defineOptions({
  name: 'State',
})

const timer = ref(null)
const state = ref({})
const colors = ref([
  { color: '#5cb87a', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#f56c6c', percentage: 80 }
])

const reload = async() => {
  const { data } = await getSystemState()
  state.value = data.server
}

reload()
timer.value = setInterval(() => {
  reload()
}, 1000 * 10)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
})

</script>

<style>

.card_item {
  @apply h-80 text-xl p-6  bg-white text-slate-700 dark:text-slate-400  dark:bg-slate-800 rounded m-2;;
}
</style>
