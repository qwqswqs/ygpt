<template>
  <div class="container">
    <el-row class="row-dashboard">
      <el-col class="col-dashboard" :span="23"><div>驾驶舱视图</div></el-col>
    </el-row>
    <el-row class="row-fields">
      <el-col class="col-fields" :span="3"><div style="margin-left: 10px;">公司名称</div></el-col>
      <el-col class="col-fields" :span="3"><div style="margin-right: 60px;">资源详情</div></el-col>
      <el-col class="col-fields" :span="2"><div style="margin-right: 30px;">数量</div></el-col>
      <el-col class="col-fields" :span="2"><div style="margin-right: 18px;">使用时间</div></el-col>
      <el-col class="col-fields" :span="13">
        <div>数据统计</div>
        <el-select v-model="month" placeholder="请选择月份" style="width: 240px; position: relative; left: 150px;">
          <el-option-group
            v-for="group in options"
            :key="group.label"
            :label="group.label"
          >
            <el-option
              v-for="item in group.options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-option-group>
        </el-select>
      </el-col>
    </el-row>
    <!-- 2025.09 -->
    <el-row v-show="month === 9" class="row-data">
      <el-col :span="3"><div>万国数据</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>1台</div></el-col>
      <el-col :span="2"><div>2024.09</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">32核/128核</div>
            <div>
              <div style="width: 25px;" class="rectangle1"></div>
              <div style="width: 75px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">430G/512G</div>
            <div>
              <div style="width: 84px;" class="rectangle1"></div>
              <div style="width: 16px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">2.53T/4T</div>
            <div>
              <div style="width: 63px;" class="rectangle1"></div>
              <div style="width: 37px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">3G/10G</div>
            <div>
              <div style="width: 30px;" class="rectangle1"></div>
              <div style="width: 70px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 9" class="row-data">
      <el-col :span="3"><div>润迅数据</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>3台</div></el-col>
      <el-col :span="2"><div>2024.09</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">288核/384核</div>
            <div>
              <div style="width: 75px;" class="rectangle1"></div>
              <div style="width: 25px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">1140G/1536G</div>
            <div>
              <div style="width: 74px;" class="rectangle1"></div>
              <div style="width: 26px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">10.92T/12T</div>
            <div>
              <div style="width: 91px;" class="rectangle1"></div>
              <div style="width: 9px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 2024.10 -->
    <el-row v-show="month === 10" class="row-data">
      <el-col :span="3"><div>映潮</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>1台</div></el-col>
      <el-col :span="2"><div>2024.10</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">64核/128核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">402G/512G</div>
            <div>
              <div style="width: 79px;" class="rectangle1"></div>
              <div style="width: 21px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">3.11T/4T</div>
            <div>
              <div style="width: 78px;" class="rectangle1"></div>
              <div style="width: 22px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 10" class="row-data">
      <el-col :span="3"><div>奥飞数据</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>1台</div></el-col>
      <el-col :span="2"><div>2024.10</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">64核/128核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">433G/512G</div>
            <div>
              <div style="width: 79px;" class="rectangle1"></div>
              <div style="width: 21px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">3.19/4T</div>
            <div>
              <div style="width: 80px;" class="rectangle1"></div>
              <div style="width: 20px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">2G/10G</div>
            <div>
              <div style="width: 20px;" class="rectangle1"></div>
              <div style="width: 80px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 10" class="row-data">
      <el-col :span="3"><div>猿届云</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>150台</div></el-col>
      <el-col :span="2"><div>2024.10</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">9600核/19200核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">36000G/76800G</div>
            <div>
              <div style="width: 47px;" class="rectangle1"></div>
              <div style="width: 53;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">319.5T/600T</div>
            <div>
              <div style="width: 53px;" class="rectangle1"></div>
              <div style="width: 47px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">8G/10G</div>
            <div>
              <div style="width: 80px;" class="rectangle1"></div>
              <div style="width: 20px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 10" class="row-data">
      <el-col :span="3"><div>四川创阔智能科技有限公司</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>100台</div></el-col>
      <el-col :span="2"><div>2024.10</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">6400核/12800核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">28800G/51200G</div>
            <div>
              <div style="width: 56px;" class="rectangle1"></div>
              <div style="width: 44px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">128T/400T</div>
            <div>
              <div style="width: 32px;" class="rectangle1"></div>
              <div style="width: 68px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">7G/10G</div>
            <div>
              <div style="width: 70px;" class="rectangle1"></div>
              <div style="width: 30px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 2024.11 -->
    <el-row v-show="month === 11" class="row-data">
      <el-col :span="3"><div>四川创阔智能科技有限公司</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>100台</div></el-col>
      <el-col :span="2"><div>2024.11</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">9600核/12800核</div>
            <div>
              <div style="width: 75px;" class="rectangle1"></div>
              <div style="width: 25px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">25600G/51200G</div>
            <div>
              <div class="rectangle1"></div>
              <div class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">321T/400T</div>
            <div>
              <div style="width: 80px;" class="rectangle1"></div>
              <div style="width: 20px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">2G/10G</div>
            <div>
              <div style="width: 20px;" class="rectangle1"></div>
              <div style="width: 80px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 11" class="row-data">
      <el-col :span="3"><div>猿届云</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>150台</div></el-col>
      <el-col :span="2"><div>2024.11</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">9600核/19200核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">4G/10G</div>
            <div>
              <div style="width: 41px;" class="rectangle1"></div>
              <div style="width: 59px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">582T/600T</div>
            <div>
              <div style="width: 97px;" class="rectangle1"></div>
              <div style="width: 3px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 2024.12 -->
    <el-row v-show="month === 12" class="row-data">
      <el-col :span="3"><div>猿届云</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>150台</div></el-col>
      <el-col :span="2"><div>2024.12</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">14400核/19200核</div>
            <div>
              <div style="width: 75px;" class="rectangle1"></div>
              <div style="width: 25px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">42000G/76800G</div>
            <div>
              <div style="width: 55px;" class="rectangle1"></div>
              <div style="width: 45px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">423T/600T</div>
            <div>
              <div style="width: 70px;" class="rectangle1"></div>
              <div style="width: 30px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 12" class="row-data">
      <el-col :span="3"><div>四川创阔智能科技有限公司</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>100台</div></el-col>
      <el-col :span="2"><div>2024.12</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">6400核/12800核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">25600G/51200G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">183T/400T</div>
            <div>
              <div style="width: 46px;" class="rectangle1"></div>
              <div style="width: 54px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 2025.01 -->
    <el-row v-show="month === 1" class="row-data">
      <el-col :span="3"><div>猿届云</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>30台</div></el-col>
      <el-col :span="2"><div>2025.01</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">1920核/38400核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">7680G/15360G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">30.9T/120T</div>
            <div>
              <div style="width: 26px;" class="rectangle1"></div>
              <div style="width: 74px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 1" class="row-data">
      <el-col :span="3"><div>四川创阔智能科技有限公司</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>20台</div></el-col>
      <el-col :span="2"><div>2025.01</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">640核/2560核</div>
            <div>
              <div style="width: 25px;" class="rectangle1"></div>
              <div style="width: 75px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5120G/10240G</div>
            <div>
              <div class="rectangle1"></div>
              <div class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">46.2T/80T</div>
            <div>
              <div style="width: 58px;" class="rectangle1"></div>
              <div style="width: 42px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 2025.02 -->
    <el-row v-show="month === 2" class="row-data">
      <el-col :span="3"><div>猿届云</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>10台</div></el-col>
      <el-col :span="2"><div>2025.2</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">640核/1280核</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">1800G/2560G</div>
            <div>
              <div style="width: 70px;" class="rectangle1"></div>
              <div style="width: 30px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">34.1T/40T</div>
            <div>
              <div style="width:85px;" class="rectangle1"></div>
              <div style="width: 15px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">2G/10G</div>
            <div>
              <div style="width: 20px;" class="rectangle1"></div>
              <div style="width: 80px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row v-show="month === 2" class="row-data">
      <el-col :span="3"><div>思腾合力</div></el-col>
      <el-col :span="3"><div>裸金属4090*8卡</div></el-col>
      <el-col :span="2"><div>2台</div></el-col>
      <el-col :span="2"><div>2025.02</div></el-col>
      <el-col :span="4">
        <el-divider style="margin-right: 20px; font-size: 80px;" direction="vertical" />
        <div class="flex-data">
          <span class="data-key">vCPU</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">64核/256核</div>
            <div>
              <div style="width: 25px;" class="rectangle1"></div>
              <div style="width: 75px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">内存</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">520G/1024G</div>
            <div>
              <div style="width: 51px;" class="rectangle1"></div>
              <div style="width: 49px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">存储</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">7.36T/8T</div>
            <div>
              <div style="width: 92px;" class="rectangle1"></div>
              <div style="width: 8px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col class="col-left" :span="3">
        <div class="flex-data">
          <span class="data-key">带宽</span>
          <div>
            <div class="data-value1">已使用 / 总配额</div>
            <div class="data-value2">5G/10G</div>
            <div>
              <div style="width: 50px;" class="rectangle1"></div>
              <div style="width: 50px;" class="rectangle2"></div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const options = ref([
  {
    label: '月份',
    options: [
      { label: '2024年09月', value: 9 },
      { label: '2024年10月', value: 10 },
      { label: '2024年11月', value: 11 },
      { label: '2024年12月', value: 12 },
      { label: '2025年01月', value: 1 },
      { label: '2025年02月', value: 2 }
    ]
  }
])

const month = ref(9)

</script>

<style scoped>
/* 总容器 */
.container {
  width: 80vw;
  height: 80vh;
  //background-color: #ffffff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-self: center;
}

/* 行：通用 */
.el-row {
  width: 80vw;
  //min-width: 1500px;
  display: flex;
  justify-content: center;
}

/* 列：通用 */
.el-col {
  background-color: rgba(183, 210, 238, 0.15);
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 2px;
}

/* 行：驾驶舱 */
.row-dashboard {
  width: calc(80vw + 30px);
  height: 60px;
  margin: 20px 0;
  font-size: 30px;
  font-weight: 900;
  color: rgb(81, 96, 110);
}

/* 列：驾驶舱 */
.col-dashboard {
  background: linear-gradient(-45deg, rgba(90, 146, 219, 0.3), rgba(4, 59, 141, 0.4), rgba(130, 186, 238, 0.3), rgba(11, 104, 211, 0.3));
  background-size: 400% 400%;
  animation: GradientShift 20s ease infinite;
}
@keyframes GradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 行：字段 */
.row-fields {
  width: calc(80vw + 30px);
  //min-width: 1440px;
  height: 60px;
  margin-bottom: 5px;
  font-size: 22px;
  font-weight: 600;
  color: rgb(81, 96, 110);
}

/* 列：字段 */
.col-fields {
  background-color: rgba(158, 181, 206, 0.2);
}

/* 行：数据 */
.row-data {
  height: 100px;
  margin: 5px 0;
}

/* 列：左分隔 */
.col-left {
  border-radius: 8px;
  margin-left: 10px;
}

/* 弹性布局：数据 */
.flex-data {
  display: flex;
  align-items: center;
}

/* 数据：key */
.data-key {
  color: rgb(132, 133, 136);
  margin-right: 15px;
  font-size: 18px;
  font-weight: 600;
}

/* 数据：value1 */
.data-value1 {
  font-size: 14px;
  color: rgb(132, 133, 136);
}

/* 数据：value2 */
.data-value2 {
  color: rgb(68, 79, 107);
  font-weight: 600;
  margin-top: 2px;
  margin-bottom: 5px;
}

/* 横向矩形 */
.rectangle1, .rectangle2 {
  width: 50px;
  height: 10px;
  display: inline-block;
}
.rectangle1 {
  background-color: rgb(59, 123, 179);
  animation: colorChange1 4s infinite alternate;
}
.rectangle2 {
  background-color: rgb(176, 209, 247);
  animation: colorChange2 6s infinite alternate;
}
@keyframes colorChange1 {
  0% { background-color: rgba(59, 123, 179); }
  100% { background-color: rgb(86, 139, 209); }
}
@keyframes colorChange2 {
  0% { background-color: rgb(176, 209, 247); }
  100% { background-color: rgba(228, 240, 252, 0.986); }
}
</style>
