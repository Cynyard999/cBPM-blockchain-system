<template>
    <el-table
            :data="deliveryOrders"
            :default-sort="{ prop: 'createTime', order: 'descending' }"
            v-loading="loading"
            ref="tableRef"
            style='width: 100%'
            height="600px"
    >
        <el-table-column prop="createTime" label="创建时间"/>
        <el-table-column prop="assetName" label="名称"/>
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
                <el-button v-if="this.user.orgType==='carrier'" type="text" size="small"
                           @click="getDeliveryOrderForm(scope.row)">
                    修改状态
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-dialog center width="600px" v-model="deliveryOrderFormVisible" title="Change Order Status Confirm">
        <el-form label-position="right" label-width="170px">
            <el-form-item label="创建时间: ">
                {{selectedDeliveryOrder.createTime}}
            </el-form-item>
            <el-form-item label="商品名称: ">
                {{selectedDeliveryOrder.assetName}}
            </el-form-item>
            <el-form-item label="备注: ">
                {{selectedDeliveryOrder.note}}
            </el-form-item>
            <el-form-item v-show="noteDeliveryDetailVisible" label="DeliveryDetailNote:">
                <el-input placeholder="note for DeliveryDetail" v-model=this.noteDeliveryDetail></el-input>
            </el-form-item>
            <el-form-item v-show="noteDeliveryDetailVisible" label="contactForDeliveryDetail:">
                <el-input placeholder="contact for DeliveryDetail" v-model=this.contact></el-input>
            </el-form-item>
            <el-select @change="showNoteDeliveryDetail" style="margin-left: 170px"
                       v-model="selectedDeliveryOrder.newStatus" placeholder="Select">
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
            <el-button @click="deliveryOrderFormVisible = false">取消</el-button>
            <el-button type="primary" @click="changeDeliveryOrderStatus()">确定</el-button>
          </span>
        </template>
    </el-dialog>
</template>

<script>
    import {request} from "../../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "DeliveryOrders",
        methods: {
            showNoteDeliveryDetail() {
                if (this.selectedDeliveryOrder.newStatus === 1) {
                    // 无法创建暂时永远设置为false
                    this.noteDeliveryDetailVisible = false;
                } else {
                    this.noteDeliveryDetailVisible = false;
                }
            },
            getDeliveryOrderForm(orderProxy) {
                this.selectedDeliveryOrder = JSON.parse(JSON.stringify(orderProxy));
                this.selectedDeliveryOrder.newStatus = this.selectedDeliveryOrder.status;
                this.deliveryOrderFormVisible = true;
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
            getDeliveryOrdersIndex(tradeID) {
                let index = -1;
                this.orders.forEach((order, i) => {
                    if (tradeID === asset["tradeID"]) {
                        index = i;
                    }
                });
                return index;
            },
            checkDisabled(status) {
                if (status - 1 > this.selectedDeliveryOrder.status) {
                    return true;
                }
                if (status < this.selectedDeliveryOrder.status) {
                    return true;
                }
                if (this.user.orgType === 'supplier') {
                    return true;
                } else {
                    return !(status === 1 || status === 2 || status === 0);
                }
            },
            changeDeliveryOrderStatus() {
                this.deliveryOrderFormVisible = false;
                if (this.selectedDeliveryOrder.status === this.selectedDeliveryOrder.newStatus) {
                    ElMessage({
                        message: '修改成功',
                        type: 'success',
                    });
                    this.getDeliveryOrders();

                } else {
                    let body = {
                        function: "",
                        args: [this.selectedDeliveryOrder.tradeID]
                    };
                    //DeliveryDetail修改status的body
                    let bodyDeliveryDetail = {
                        function: "",
                        args: [this.selectedDeliveryOrder.tradeID]
                    };

                    if (this.selectedDeliveryOrder.newStatus === 1) {
                        body.function = "HandleDeliveryOrder";
                        bodyDeliveryDetail.function = "HandleDeliveryArrangement";
                        // 无法通过deliveryOrder这边创建deliveryDetail
                        // this.createDeliveryDetail();
                    }
                    if (this.selectedDeliveryOrder.newStatus === 2) {
                        body.function = "FinishDeliveryOrder";
                        bodyDeliveryDetail.function = "FinishDeliveryArrangement";
                    }
                    if (this.selectedDeliveryOrder.newStatus === 3) {
                        body.function = "ConfirmFinishDeliveryOrder";
                        bodyDeliveryDetail.function = "ConfirmDeliveryArrangement";
                    }
                    let that = this;
                    that.loading = true;

                    //先修改DeliveryArrangement的status
                    request('/work/invoke', bodyDeliveryDetail, "POST").then(response => {
                        ElMessage({
                            message: '修改DeliveryArrangement成功',
                            type: 'success',
                        });
                    }).catch(error => {
                        that.loading = false;
                    });

                    request('/work/invoke', body, "POST").then(response => {
                        ElMessage({
                            message: '修改成功',
                            type: 'success',
                        });
                        that.loading = false;
                        that.getDeliveryOrders();
                    }).catch(error => {
                        that.loading = false;
                    });
                }
            },
            getDeliveryOrders() {
                let body = {
                    function: "GetAllDeliveryOrders",
                    args: []
                };
                let that = this;
                this.loading = true;
                request('/work/query', body, "POST").then(response => {
                    that.deliveryOrders = response.data.result;
                    if (that.deliveryOrders === null) {
                        that.deliveryOrders = [];
                    }
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            getUser() {
                this.user = JSON.parse(window.localStorage.getItem("user"));
            },
            createDeliveryDetail() {
                let body = {
                    function: "CreateDeliveryDetail",
                    transient: {
                        detail: {
                            TradeId: this.selectedDeliveryOrder.tradeID,
                            AssetName: this.selectedDeliveryOrder.assetName,
                            // 这两个属性deliveryorder并没有
                            StartPlace: this.selectedDeliveryOrder.startPlace,
                            EndPlace: this.selectedDeliveryOrder.endPlace,
                            Contact: this.contact,
                            note: this.noteDeliveryDetail,
                        }
                    }
                };
                request('/work/invoke', body, "POST").then(response => {
                    ElMessage({
                        message: '创建DeliveryDetail成功',
                        type: 'success',
                    });
                    console.log('创建DeliveryDetail成功')
                }).catch(error => {
                    console.log('创建DeliveryDetail失败')
                });
            },
        },
        data() {
            return {
                deliveryOrders: [],
                loading: true,
                user: {},
                deliveryOrderFormVisible: false,
                noteDeliveryDetail: "",
                contact: "",
                noteDeliveryDetailVisible: false,
                selectedDeliveryOrder: {},
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
                ]
            }
        },
        mounted() {
            this.getUser();
            this.getDeliveryOrders();
        }
    }
</script>

<style>

</style>