<template>
    <Header/>
    <el-container>
        <el-aside>
            <el-row>
                <el-col>
                    <div class="name">
                        <h5>中间商后台管理</h5>
                    </div>
                    <el-menu
                            :router="true"
                    >
                        <el-sub-menu index="1">
                            <template #title>
                                <el-icon>
                                    <goods/>
                                </el-icon>
                                <span>与供应商</span>
                            </template>
                            <el-menu-item index="/middleman/assets">货物信息</el-menu-item>
                            <el-menu-item index="/middleman/supply-orders">进货订单</el-menu-item>
                        </el-sub-menu>

                        <el-sub-menu index="2">
                            <template #title>
                                <el-icon>
                                    <sort/>
                                </el-icon>
                                <span>与生产商</span>
                            </template>
                            <el-menu-item index="/middleman/goods" >供应商品</el-menu-item>
                            <el-menu-item index="/middleman/orders" >商品订单</el-menu-item>
                        </el-sub-menu>

                        <el-sub-menu index="3">
                            <template #title>
                                <el-icon>
                                    <refresh/>
                                </el-icon>
                                <span>与运输商</span>
                            </template>
                            <el-menu-item index="/middleman/delivery-arrangements" >运输安排</el-menu-item>
                        </el-sub-menu>

                    </el-menu>
                </el-col>
            </el-row>
        </el-aside>
        <el-main>
            <router-view name="middlemanSubpage"></router-view>
        </el-main>
    </el-container>
</template>

<script>
    import Header from "../components/Header.vue"

    export default {
        name: "Middleman.vue",
        components: {
            Header
        }
    }
</script>

<style>
    .name {
        margin-left: 35px;
        font-size: large;
    }

</style>