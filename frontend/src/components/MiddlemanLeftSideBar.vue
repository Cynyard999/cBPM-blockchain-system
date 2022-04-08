<template>
  <el-row class="tac">
    <el-col :span="4">
      <div class="mb-2">
      <h5 >选 项</h5>
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
            <span>与生产商</span>
          </template>
          <el-menu-item-group title="商品相关">
            <el-menu-item index="1-1">增加商品</el-menu-item>
            <el-menu-item index="1-2">查看所有商品</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="订单相关">
            <el-menu-item index="1-3">创建订单</el-menu-item>
            <el-menu-item index="1-4">查看所有订单</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>

        <el-sub-menu index="2">
          <template #title>
            <el-icon><sort /></el-icon>
            <span>与供货商</span>
          </template>
          <el-menu-item-group title="商品相关">
            <el-menu-item index="2-1">增加商品</el-menu-item>
            <el-menu-item index="2-2">查看所有商品</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="订单相关">
            <el-menu-item index="2-3">创建订单</el-menu-item>
            <el-menu-item index="2-4">查看所有订单</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>

        <el-sub-menu index="3">
          <template #title>
            <el-icon><refresh /></el-icon>
            <span>与运输商</span>
          </template>
          <el-menu-item-group title="商品相关">
            <el-menu-item index="3-1">增加商品</el-menu-item>
            <el-menu-item index="3-2">查看所有商品</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="订单相关">
            <el-menu-item index="3-3">创建订单</el-menu-item>
            <el-menu-item index="3-4">查看所有订单</el-menu-item>
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
  margin-left: 60px;
  font-size: large;

}
</style>