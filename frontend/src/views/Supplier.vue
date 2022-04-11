<template>
    <Header/>
    <el-container>
        <el-aside>
            <el-row>
                <el-col>
                    <div class="name">
                        <h5>供货商后台管理</h5>
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
                            <el-menu-item index="/supplier/assets">货物信息</el-menu-item>
                            <el-menu-item index="/supplier/supply-orders">进货订单</el-menu-item>
                        </el-sub-menu>

                        <el-sub-menu index="2">
                            <template #title>
                                <el-icon>
                                    <sort/>
                                </el-icon>
                                <span>与运输商</span>
                            </template>
                            <el-menu-item index="/supplier/delivery-orders">运输订单</el-menu-item>
                        </el-sub-menu>
                    </el-menu>
                </el-col>
            </el-row>
        </el-aside>
        <el-main>
            <router-view name="supplierSubpage"></router-view>
        </el-main>
    </el-container>
</template>

<script>
    import Header from "../components/Header.vue"
    import {request} from "../api/axios";

    export default {
        name: "Supplier",
        components: {
            Header
        },
        methods: {},
        mounted() {
            // request('/work/query', {}, 'POST').then(response => {
            //     console.log(response)
            // }).catch(error => {
            //     console.log(error)
            // });
        }
    }
</script>

<style>
    .name {
        margin-left: 35px;
        font-size: large;
    }

</style>