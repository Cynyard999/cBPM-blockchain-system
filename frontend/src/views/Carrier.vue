<template>
    <Header/>
    <el-container>
        <el-aside>
            <el-row>
                <el-col>
                    <div class="name">
                        <h5>运输商后台管理</h5>
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
                            <el-menu-item index="/carrier/delivery-arrangements">运输安排</el-menu-item>
                        </el-sub-menu>

                        <el-sub-menu>
                            <template #title>
                                <el-icon>
                                    <goods/>
                                </el-icon>
                                <span>与生产商</span>
                            </template>
                            <el-menu-item index="/carrier/delivery-details">运输信息</el-menu-item>
                        </el-sub-menu>

                        <el-sub-menu index="2">
                            <template #title>
                                <el-icon>
                                    <sort/>
                                </el-icon>
                                <span>与供货商</span>
                            </template>
                            <el-menu-item index="/carrier/delivery-orders">运输订单</el-menu-item>
                        </el-sub-menu>
                    </el-menu>
                </el-col>
            </el-row>
        </el-aside>
        <el-main>
            <router-view name="carrierSubpage"></router-view>
        </el-main>
    </el-container>
</template>

<script>
    import Header from "../components/Header.vue"

    export default {
        name: "Carrier",
        components: {
            Header
        },
    }
</script>


<style>
    .name {
        margin-left: 35px;
        font-size: large;
    }

</style>