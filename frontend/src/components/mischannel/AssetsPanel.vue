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
        <el-table-column v-if="user.orgType === 'supplier'" label="Operations" width="120">
            <template #default="scope">
                <el-popconfirm title="确定要删除这个货品吗?" @confirm="deleteAsset(scope.row)">
                    <template #reference>
                        <el-button type="text" size="small">
                            删除
                        </el-button>
                    </template>
                </el-popconfirm>
                <el-button type="text" size="small">

                </el-button>
            </template>
        </el-table-column>
        <el-table-column v-if="user.orgType === 'middleman'" label="Operations" width="120">
            <template #default="scope">
                <el-button type="text" size="small" @click="getPublishAssetForm(scope.row)">
                    发布
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-dialog center width="500px" v-model="publishAssetFormVisible" title="Publish Asset Confirm">
        <el-form>
            <el-form-item label="名称: ">
                <el-input disabled v-model="selectedAsset.assetName"/>
            </el-form-item>
            <el-form-item label="单价: ">
                <el-input v-model="selectedAsset.assetPrice"/>
            </el-form-item>
            <el-form-item label="描述: ">
                <el-input v-model="selectedAsset.publicDescription"/>
            </el-form-item>
            <el-form-item label="发货地: ">
                <el-input disabled v-model="selectedAsset.shippingAddress"/>
            </el-form-item>
        </el-form>
        <template #footer>
          <span>
            <el-button @click="publishAssetFormVisible = false">取消</el-button>
            <el-button type="primary" @click="publishAsset()">发布</el-button>
          </span>
        </template>
    </el-dialog>
</template>

<script>
    import {request} from "../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "assets",
        methods: {
            deleteAsset(assetProxy) {
                let index = this.getAssetIndexIndex(assetProxy.assetID);
                let body = {
                    channelName: "mischannel",
                    contractName: "mischaincode",
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
                }).catch(error=>{
                    this.loading = false;
                })
            },
            // TODO updateAsset
            // TODO 创建完成后显示创建的东西
            // TODO 值校验 非空 数字
            // TODO 前后端自定义返回的状态，从而打印不同的消息

            getPublishAssetForm(assetProxy) {
                this.selectedAsset = JSON.parse(JSON.stringify(assetProxy));
                delete this.selectedAsset.ownerOrg;
                this.publishAssetFormVisible = true;
            },

            publishAsset() {
                this.selectedAsset.assetPrice = parseFloat(this.selectedAsset.assetPrice);
                let body = {
                    channelName: "mamichannel",
                    contractName: "mamichaincode",
                    function: "CreateAsset",
                    transient: {
                        asset: this.selectedAsset
                    }
                };
                let that = this;
                that.publishAssetFormVisible = false;
                that.loading = true;
                request('/work/invoke', body, "POST").then(response => {
                    ElMessage({
                        message: '发布成功',
                        type: 'success',
                    });
                    that.loading = false;
                }).catch(error => {
                    that.loading = false;
                });
            },
            // get index of asset in assets(sort function causes the returning index is not correct)
            getAssetIndexIndex(assetID) {
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
                    channelName: "mischannel",
                    contractName: "mischaincode",
                    function: "GetAllAssets",
                    args: []
                };
                let that = this;
                request('/work/query', body, "POST").then(response => {
                    that.assets = response.data.result;
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
                assets: [],
                loading: true,
                user: {},
                publishAssetFormVisible: false,
                selectedAsset: {},
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