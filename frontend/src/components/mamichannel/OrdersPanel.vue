<template>
    <el-table
            :data="orders"
            :default-sort="{ prop: 'createTime', order: 'descending' }"
            v-loading="loading"
            ref="tableRef"
            style='width: 100%'
            height="600px"
    >
        <el-table-column prop="createTime" label="创建时间"/>
        <el-table-column prop="assetName" label="名称"/>
        <el-table-column prop="assetPrice" label="单价" sortable/>
        <el-table-column prop="quantity" label="数量"/>
        <el-table-column prop="totalPrice" label="总价"/>
        <el-table-column prop="receivingAddress" label="发货地"/>
        <el-table-column prop="shippingAddress" label="发货地"/>
        <el-table-column prop="note" label="备注" overflow/>
        <el-table-column prop="ownerOrg" label="所属组织"/>
        <el-table-column prop="updateTime" label="修改时间"/>
        <el-table-column
                prop="tag"
                label="状态"
                width="100"
                :filters="[{ text: '未处理', value: 0 },{ text: '开始处理', value: 1 },{ text: '处理完成', value: 2 },{ text: '确认完成', value: 3 }]"
                :filter-method="filterStatus"
                filter-placement="bottom-end">
            <template #default="scope">
                <el-tag
                        :type="getStatusTag(scope.row)"
                        disable-transitions
                >{{ getStatus(scope.row.status) }}
                </el-tag
                >
            </template>
        </el-table-column>
        <el-table-column label="Operations" width="120">
            <template #default="scope">
                <el-button type="text" size="small" @click="getOrderForm(scope.row)">
                    修改状态
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-dialog center width="500px" v-model="orderFormVisible" title="Change Order Status Confirm">
        <el-form>
            <el-form-item label="创建时间: ">
                {{selectedOrder.createTime}}
            </el-form-item>
            <el-form-item label="商品名称: ">
                {{selectedOrder.assetName}}
            </el-form-item>
            <el-form-item label="数量: ">
                {{selectedOrder.quantity}}
            </el-form-item>
            <el-form-item label="总价: ">
                {{selectedOrder.totalPrice}}
            </el-form-item>
            <el-form-item label="收货地: ">
                {{selectedOrder.receivingAddress}}
            </el-form-item>
            <el-form-item label="备注: ">
                {{selectedOrder.note}}
            </el-form-item>
            <el-select v-model="selectedOrder.newStatus" placeholder="Select">
                <el-option
                        v-for="item in statusOptions"
                        :key="item"
                        :label="item.label"
                        :value="item.value"
                        :disabled="checkDisabled(item.value)"
                />
            </el-select>
        </el-form>
        <template #footer>
          <span>
            <el-button @click="orderFormVisible = false">取消</el-button>
            <el-button type="primary" @click="changeOrderStatus()">确定</el-button>
          </span>
        </template>
    </el-dialog>
</template>

<script>
    import {request} from "../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "orders",
        methods: {

            getOrderForm(orderProxy) {
                this.selectedOrder = JSON.parse(JSON.stringify(orderProxy));
                this.selectedOrder.newStatus = this.selectedOrder.status;
                this.orderFormVisible = true;
            },
            filterStatus(value, row) {
                return row.status === value
            },
            getStatusTag(row) {
                if (row.status === 0) {
                    return 'info';
                }
                if (row.status === 1) {
                    return 'warning';
                }
                if (row.status === 2) {
                    return ''
                }
                return 'success'
            },
            getStatus(status) {
                return this.statusOptions[status].label;
            },
            // get index of asset in orders(sort function causes the returning index is not correct)
            getOrdersIndex(tradeID) {
                let index = -1;
                this.orders.forEach((order, i) => {
                    if (tradeID === asset["tradeID"]) {
                        index = i;
                    }
                });
                return index;
            },
            checkDisabled(status) {
                if (status - 1 > this.selectedOrder.status) {
                    return true;
                }
                if (status < this.selectedOrder.status) {
                    return true;
                }
                if (this.user.orgType === 'manufacturer') {
                    return !(status === 3 || status === 2);
                } else {
                    return !(status === 1 || status === 2 || status === 0);
                }
            },
            changeOrderStatus() {
                this.orderFormVisible = false;
                if (this.selectedOrder.status === this.selectedOrder.newStatus) {
                    ElMessage({
                        message: '修改成功',
                        type: 'success',
                    });
                    this.getOrders();

                } else {
                    let body = {
                        channelName: "mamichannel",
                        contractName: "mamichaincode",
                        function: "",
                        args: [this.selectedOrder.tradeID]
                    };
                    if (this.selectedOrder.newStatus === 1) {
                        body.function = "HandleOrder";
                    }
                    if (this.selectedOrder.newStatus === 2) {
                        body.function = "FinishOrder";
                    }
                    if (this.selectedOrder.newStatus === 3) {
                        body.function = "ConfirmFinishOrder";
                    }
                    let that = this;
                    that.loading = true;
                    request('/work/invoke', body, "POST").then(response => {
                        ElMessage({
                            message: '修改成功',
                            type: 'success',
                        });
                        that.loading = false;
                        that.getOrders();
                    }).catch(error => {
                        that.loading = false;
                    });
                }
            },
            getOrders() {
                let body = {
                    channelName: "mamichannel",
                    contractName: "mamichaincode",
                    function: "GetAllOrders",
                    args: []
                };
                let that = this;
                this.loading = true;
                request('/work/query', body, "POST").then(response => {
                    that.orders = response.data.result;
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            getUser() {
                this.user = JSON.parse(window.localStorage.getItem("user"));
            }
        },
        data() {
            return {
                orders: [],
                loading: true,
                user: {},
                orderFormVisible: false,
                selectedOrder: {},
                statusOptions: [
                    {
                        value: 0,
                        label: '未处理'
                    },
                    {
                        value: 1,
                        label: '开始处理'
                    },
                    {
                        value: 2,
                        label: '处理完成'
                    },
                    {
                        value: 3,
                        label: '确认完成'
                    }
                ]
            }
        },
        mounted() {
            this.getUser();
            this.getOrders();
        }
    }
</script>

<style>

</style>