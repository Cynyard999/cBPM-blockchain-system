<template>
    <el-row class="tac">
        <el-col class="column" :span="4">
            <div class="name">
                <h5>供货商后台管理</h5>
            </div>
            <el-menu
                    default-active="1"
                    class="el-menu-vertical"
            >
                <el-sub-menu index="1">
                    <template #title>
                        <el-icon>
                            <rank/>
                        </el-icon>
                        <span>与中间商</span>
                    </template>
                    <el-menu-item index="1-1">货物信息</el-menu-item>
                    <el-menu-item index="1-2">进货订单</el-menu-item>
                </el-sub-menu>

                <el-sub-menu index="2">
                    <template #title>
                        <el-icon>
                            <sort/>
                        </el-icon>
                        <span>与运输商</span>
                    </template>
                    <el-menu-item index="2-1">运输信息</el-menu-item>
                </el-sub-menu>
            </el-menu>
        </el-col>

    </el-row>
</template>

<script>
    export default {
        name: "SupplierSideBar",
        methods: {
            clickMenuItem(menuItem) {
                console.log(menuItem.index)
            }
        }
    }

</script>
<style>

    .name {
        margin-left: 35px;
        font-size: large;
    }
</style>