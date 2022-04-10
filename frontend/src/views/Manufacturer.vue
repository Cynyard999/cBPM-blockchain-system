<template>
    <Header/>
    <el-container>
        <el-aside>
            <el-row>
                <el-col>
                    <div class="name">
                        <h5>生产商后台管理</h5>
                    </div>
                    <el-menu
                            :router="true"
                    >
                        <el-sub-menu index="1">
                            <template #title>
                                <el-icon>
                                    <rank/>
                                </el-icon>
                                <span>与中间商</span>
                            </template>
                            <el-menu-item index="/manufacturer/assets">供应商品</el-menu-item>
                            <el-menu-item index="/manufacturer/orders">商品订单</el-menu-item>
                        </el-sub-menu>
                        <el-sub-menu index="2">
                            <template #title>
                                <el-icon>
                                    <sort/>
                                </el-icon>
                                <span>与运输商</span>
                            </template>
                            <el-menu-item index="/manufacturer/delivery-details">运输信息</el-menu-item>
                        </el-sub-menu>
                    </el-menu>
                </el-col>
            </el-row>
        </el-aside>
        <el-main>
            <router-view name="manufacturerSubpage"></router-view>
        </el-main>
    </el-container>
</template>

<script>
    import {request} from "../api/axios";
    import Header from "../components/Header.vue"

    export default {
        name: 'Home',
        components: {
            Header
        },
        methods: {},
        mounted() {
        }
    }
</script>

<style scoped>

</style>