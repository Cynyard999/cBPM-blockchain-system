<template>
    <el-table
            :data="assets"
            :default-sort="{ prop: '单价', order: 'descending' }"
            v-loading="loading"
            style='width: 100%'
            height="600px"
    >
        <!--        <el-table-column prop="assetID" label="ID"/>-->
        <el-table-column prop="assetName" label="名称"/>
        <el-table-column prop="assetPrice" label="单价" sortable/>
        <el-table-column prop="shippingAddress" label="发货地"/>
        <el-table-column prop="publicDescription" label="描述"/>
        <el-table-column prop="ownerOrg" label="所属组织"/>
        <el-table-column v-if="user.orgType === 'middleman'" label="Operations" width="120">
            <template #default="scope">
                <el-popconfirm title="确定要删除这个商品吗?" @confirm="deleteAsset(scope.row)">
                    <template #reference>
                        <el-button type="text" size="small">
                            删除
                        </el-button>
                    </template>
                </el-popconfirm>
              <el-button @click="beforeUpdateAsset(scope.row)" type="text" size="small">
                修改
              </el-button>
            </template>
        </el-table-column>
        <el-table-column v-if="user.orgType === 'manufacturer'" label="Operations" width="120">
            <template #default="scope">
                <el-button type="text" size="small" @click="getOrderForm(scope.row)">
                    订购
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-dialog center width="500px" v-model="orderFormVisible" title="Create Order Confirm">
        <el-form label-position="right" label-width="60px">
            <el-form-item label="名称: ">
                <el-input disabled v-model="selectedAsset.assetName"/>
            </el-form-item>
            <el-form-item label="单价: ">
                <el-input disabled v-model="selectedAsset.assetPrice"/>
            </el-form-item>
            <el-form-item label="描述: ">
                <el-input disabled v-model="selectedAsset.publicDescription"/>
            </el-form-item>
            <el-form-item label="发货地: ">
                <el-input disabled v-model="selectedAsset.shippingAddress"/>
            </el-form-item>
            <el-form-item label="数量: ">
                <el-input type="number" v-model="order.quantity"/>
            </el-form-item>
            <el-form-item label="总价: ">
                {{ (order.quantity * selectedAsset.assetPrice).toFixed(2) }}
            </el-form-item>
            <el-form-item label="收货地: ">
                <el-input autosize="true" v-model="order.receivingAddress"/>
            </el-form-item>
            <el-form-item label="备注: ">
                <el-input v-model="order.note"/>
            </el-form-item>
        </el-form>
        <template #footer>
          <span>
            <el-button @click="orderFormVisible = false">取消</el-button>
            <el-button type="primary" @click="createOrder()">创建</el-button>
          </span>
        </template>
    </el-dialog>
  <!--updateAsset的form-->
  <el-dialog center width="500px" v-model="updateAssetFormVisiable" title="update asset">
    <el-form label-position="right" label-width="60px">
      <el-form-item label="名称: ">
        <el-input disabled v-model="updateAssetData.assetName"></el-input>
      </el-form-item>
      <el-form-item label="单价: ">
        <el-input v-model=updateAssetData.assetPrice></el-input>
      </el-form-item>
      <el-form-item label="发货地: ">
        <el-input disabled v-model=updateAssetData.shippingAddress></el-input>
      </el-form-item>
      <el-form-item label="描述: ">
        <el-input v-model=updateAssetData.publicDescription></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
          <span>
            <el-button @click="updateAssetFormVisiable = false">取消</el-button>
            <el-button type="primary" @click="updateAsset()">确定</el-button>
          </span>
    </template>
  </el-dialog>
</template>

<script>
    import {request} from "../../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "assets",
        methods: {
            deleteAsset(assetProxy) {
                let index = this.getAssetIndex(assetProxy.assetID);
                let body = {
                    function: "DeleteAsset",
                    args: [assetProxy.assetID]
                };
                let that = this;
                that.loading = true;
                request('/work/invoke', body, "POST").then(response => {
                    that.assets.splice(index, 1);
                    that.loading = false;
                    ElMessage({
                        message: '删除成功',
                        type: 'success',
                    });
                }).catch(error => {
                    this.loading = false;
                })
            },

            getOrderForm(assetProxy) {
                this.selectedAsset = JSON.parse(JSON.stringify(assetProxy));
                this.order.assetID = this.selectedAsset.assetID;
                this.orderFormVisible = true;
            },

            createOrder() {
                this.order.quantity = parseInt(this.order.quantity);
                let body = {
                    function: "CreateOrder",
                    transient: {
                        order: this.order
                    }
                };
                let that = this;
                that.orderFormVisible = false;
                that.loading = true;
                request('/work/invoke', body, "POST").then(response => {
                    ElMessage({
                        message: '创建成功',
                        type: 'success',
                    });
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            // get index of asset in assets(sort function causes the returning index is not correct)
            getAssetIndex(assetID) {
                let index = -1;
                this.assets.forEach((asset, i) => {
                    if (assetID === asset["assetID"]) {
                        index = i;
                    }
                });
                return index;
            },
            getAssets() {
                let body = {
                    function: "GetAllAssets",
                    args: []
                };
                let that = this;
                request('/work/query', body, "POST").then(response => {
                    that.assets = response.data.result;
                    if (that.assets === null) {
                        that.assets = [];
                    }
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            getUser() {
                this.user = JSON.parse(window.localStorage.getItem("user"));
            },
          beforeUpdateAsset(assetProxy) {
            this.selectedAsset = JSON.parse(JSON.stringify(assetProxy));
            this.updateAssetFormVisiable = true;
            this.updateAssetData.assetID = this.selectedAsset.assetID;
            this.updateAssetData.assetName = this.selectedAsset.assetName;
            this.updateAssetData.assetPrice = this.selectedAsset.assetPrice;
            this.updateAssetData.shippingAddress = this.selectedAsset.shippingAddress;
            this.updateAssetData.publicDescription = this.selectedAsset.publicDescription;
          },
          updateAsset() {
            let body = {
              function: "UpdateAsset",
              args: [
                this.updateAssetData.assetID,
                this.updateAssetData.assetName,
                this.updateAssetData.publicDescription,
                this.updateAssetData.assetPrice * 1
                ]
            };
            if (this.updateAssetData.assetName.length === 0) {
              ElMessage({
                message: 'assetName不能为空',
                type: 'warning',
              });
              this.updateAssetFormVisiable = false;
              return;
            }
            if (this.updateAssetData.assetPrice * 1 < 0) {
              ElMessage({
                message: 'assetPrice必须>=0',
                type: 'warning',
              });
              this.updateAssetFormVisiable = false;
              return;
            }
            if (this.updateAssetData.shippingAddress.length === 0) {
              ElMessage({
                message: '发货地不能为空',
                type: 'warning',
              });
              this.updateAssetFormVisiable = false;
              return;
            }
            let that=this;
            that.updateAssetFormVisiable = false;
            that.loading=true;
            request('/work/invoke', body, "POST").then(response => {
              ElMessage({
                message: '修改成功',
                type: 'success',
              });
              that.loading=false;
              that.getAssets();
            }).catch(error => {
              this.loading = false;
            })
          },
        },
        data() {
            return {
                assets: [],
                loading: true,
                user: {},
                orderFormVisible: false,
                updateAssetFormVisiable: false,
                updateAssetData: {
                assetID: "",
                assetName: "",
                assetPrice: "",
                shippingAddress: "",
                publicDescription: ""
                  },
                selectedAsset: {},
                order: {
                    assetID: "",
                    quantity: 0,
                    receivingAddress: "",
                    note: "",
                }
            }
        },
        mounted() {
            this.getUser();
            this.getAssets();
        }
    }
</script>

<style>

</style>