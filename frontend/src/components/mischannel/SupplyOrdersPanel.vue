<template>
    <el-table
            :data="supplyOrders"
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
        <el-table-column prop="shippingAddress" label="发货地"/>
        <el-table-column prop="note" label="备注" overflow/>
        <el-table-column prop="ownerOrg" label="所属组织"/>
        <el-table-column prop="handlerOrg" label="处理组织"/>
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
                <el-button type="text" size="small" @click="getSupplyOrderForm(scope.row)">
                    修改状态
                </el-button>
            </template>
        </el-table-column>
    </el-table>

    <el-dialog center width="500px" v-model="supplyOrderFormVisible" title="Change Order Status Confirm">
        <el-form label-position="right" label-width="100px">
            <el-form-item label="创建时间: ">
                {{selectedSupplyOrder.createTime}}
            </el-form-item>
            <el-form-item label="商品名称: ">
                {{selectedSupplyOrder.assetName}}
            </el-form-item>
            <el-form-item label="数量: ">
                {{selectedSupplyOrder.quantity}}
            </el-form-item>
            <el-form-item label="总价: ">
                {{selectedSupplyOrder.totalPrice}}
            </el-form-item>
            <el-form-item label="发货地: ">
                {{selectedSupplyOrder.shippingAddress}}
            </el-form-item>
            <el-form-item label="备注: ">
                {{selectedSupplyOrder.note}}
            </el-form-item>
            <el-form-item v-show="noteDeliveryOrderVisible" label="DeliveryNote:">
                <el-input placeholder="note for deliveryOrder" v-model=this.noteForDeliveryOrder></el-input>
            </el-form-item>
            <el-select @change="showNoteDelivery" style="margin-left: 100px" v-model="selectedSupplyOrder.newStatus"
                       placeholder="Select">
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
            <el-button @click="supplyOrderFormVisible = false">取消</el-button>
            <el-button type="primary" @click="changeSupplyOrderStatus()">确定</el-button>
          </span>
        </template>
    </el-dialog>
</template>

<script>
    import {request} from "../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "SupplyOrders",
        methods: {
            showNoteDelivery() {
                if (this.selectedSupplyOrder.newStatus === 2) {
                    this.noteDeliveryOrderVisible = true;
                } else {
                    this.noteDeliveryOrderVisible = false;
                }
            },
            getSupplyOrderForm(orderProxy) {
                this.selectedSupplyOrder = JSON.parse(JSON.stringify(orderProxy));
                this.selectedSupplyOrder.newStatus = this.selectedSupplyOrder.status;
                this.supplyOrderFormVisible = true;
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
            getSupplyOrdersIndex(tradeID) {
                let index = -1;
                this.orders.forEach((order, i) => {
                    if (tradeID === asset["tradeID"]) {
                        index = i;
                    }
                });
                return index;
            },
            checkDisabled(status) {
                if (status - 1 > this.selectedSupplyOrder.status) {
                    return true;
                }
                if (status < this.selectedSupplyOrder.status) {
                    return true;
                }
                if (this.user.orgType === 'middleman') {
                    return !(status === 3);
                } else {
                    return !(status === 1 || status === 2 || status === 0);
                }
            },
            changeSupplyOrderStatus() {
                this.supplyOrderFormVisible = false;
                if (this.selectedSupplyOrder.status === this.selectedSupplyOrder.newStatus) {
                    ElMessage({
                        message: '修改成功',
                        type: 'success',
                    });
                    this.getOrders();

                } else {
                    let body = {
                        channelName: "mischannel",
                        contractName: "mischaincode",
                        function: "",
                        args: [this.selectedSupplyOrder.tradeID]
                    };
                    if (this.selectedSupplyOrder.newStatus === 1) {
                        body.function = "HandleSupplyOrder";
                    }
                    if (this.selectedSupplyOrder.newStatus === 2) {
                        body.function = "FinishSupplyOrder";
                    }
                    if (this.selectedSupplyOrder.newStatus === 3) {
                        body.function = "ConfirmFinishSupplyOrder";
                    }
                    let that = this;
                    that.loading = true;
                    //完成supplyerOrder后一起创建deleveryOrder
                    if (this.selectedSupplyOrder.newStatus === 2) {
                        this.createDeliveryOrder();
                    }
                    request('/work/invoke', body, "POST").then(response => {
                        ElMessage({
                            message: '修改成功',
                            type: 'success',
                        });
                        that.loading = false;
                        that.noteDeliveryOrderVisible = false;
                        that.getSupplyOrders();
                    }).catch(error => {
                        that.loading = false;
                    });
                }
            },
            getSupplyOrders() {
                let body = {
                    channelName: "cbpmchannel",
                    contractName: "cbpmchaincode",
                    function: "GetAllSupplyOrders",
                    args: []
                };
                let that = this;
                this.loading = true;
                request('/work/query', body, "POST").then(response => {
                    that.supplyOrders = response.data.result;
                    if (that.supplyOrders === null) {
                        that.supplyOrders = [];
                    }
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            getUser() {
                this.user = JSON.parse(window.localStorage.getItem("user"));
            },
            createDeliveryOrder() {
                let body = {
                    channelName: "scchannel",
                    contractName: "scchaincode",
                    function: "CreateDeliveryOrder",
                    transient: {
                        order: {
                            tradeID: this.selectedSupplyOrder.tradeID,
                            assetName: this.selectedSupplyOrder.assetName,
                            note: this.noteForDeliveryOrder,
                        }

                    }
                };
                request('/work/invoke', body, "POST").then(response => {
                    ElMessage({
                        message: '创建DeliveryOrder成功',
                        type: 'success',
                    });
                    console.log('创建DeliveryOrder成功')
                }).catch(error => {
                    console.log('创建DeliveryOrder失败')
                });
            }
        },
        data() {
            return {
                supplyOrders: [],
                loading: true,
                user: {},
                supplyOrderFormVisible: false,
                noteDeliveryOrderVisible: false,
                noteForDeliveryOrder: "",
                selectedSupplyOrder: {},
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
            this.getSupplyOrders();
        }
    }
</script>

<style>

</style>