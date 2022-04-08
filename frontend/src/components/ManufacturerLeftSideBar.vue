<template>
  <el-row class="tac" >
    <el-col class="column" :span="4">
      <div class="mb-2">
        <h5 >生产商后台管理</h5>
      </div>
      <el-menu
          default-active="1"
          class="el-menu-vertical"
          @open="handleOpen"
          @close="handleClose"
      >
        <el-sub-menu index="1">
          <template #title>
            <el-icon ><goods /></el-icon>
            <span>与中间商</span>
          </template>
          <el-menu-item-group  title="商品相关">
            <el-menu-item  index="1-1">查看所有商品</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="订单相关">
            <el-menu-item index="1-2">创建订单</el-menu-item>
            <el-menu-item index="1-3">查看所有订单</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>

        <el-sub-menu index="2">
          <template #title>
            <el-icon><sort /></el-icon>
            <span>与运输商</span>
          </template>
          <el-menu-item-group title="订单相关">
            <el-menu-item index="2-1">查看所有订单</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>

      </el-menu>
    </el-col>

  </el-row>
</template>

<script lang="ts" setup>
import {
  Document,
  Menu as IconMenu,
  Location,
  Refresh,
  Setting,
  Sort,
  Goods,
} from '@element-plus/icons-vue'
const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}

</script>
<style>

.mb-2{
  margin-left: 35px;
  font-size: large;

}
</style>