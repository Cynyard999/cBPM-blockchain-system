<template>
    <el-row class="tac">
        <el-col class="column" :span="4">
            <div class="name">
                <h5>生产商后台管理</h5>
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
                    <el-menu-item index="2-1" @click="clickMenuItem">供应商品</el-menu-item>
                    <el-menu-item index="2-2" @click="clickMenuItem">商品订单</el-menu-item>
                </el-sub-menu>

                <el-sub-menu index="2">
                    <template #title>
                        <el-icon>
                            <sort/>
                        </el-icon>
                        <span>与运输商</span>
                    </template>
                    <el-menu-item index="2-1">运输订单</el-menu-item>
                </el-sub-menu>

            </el-menu>
        </el-col>

    </el-row>
</template>

<script>
    const handleOpen = (key, keyPath) => {
        console.log(key, keyPath)
    };
    const handleClose = (key, keyPath) => {
        console.log(key, keyPath)
    }

</script>
<style>

    .name {
        margin-left: 35px;
        font-size: large;

    }
</style>