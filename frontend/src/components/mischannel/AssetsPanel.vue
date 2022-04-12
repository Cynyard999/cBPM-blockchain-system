<template>
  <el-button @click="this.addAssetFormVisiable=true" v-if="user.orgType === 'supplier'" type="primary" style="margin-left: 1040px" plain>创建新Asset</el-button>
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
                <el-button @click="beforeUpdateAsset(scope.row)" type="text" size="small">
                    修改
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
<!--createAsset的form-->
  <el-dialog center width="500px"  v-model="addAssetFormVisiable" title="add new asset">
    <el-form>
      <el-form-item label="名称: ">
        <el-input  v-model=newAsset.assetName />
      </el-form-item>
      <el-form-item label="单价: ">
        <el-input v-model=newAsset.assetPrice />
      </el-form-item>
      <el-form-item label="发货地: ">
        <el-input v-model=newAsset.shippingAddress />
        </el-form-item>
        <el-form-item label="描述: ">
          <el-input v-model=newAsset.publicDescription />
        </el-form-item>
    </el-form>
    <template #footer>
          <span>
            <el-button @click="addAssetFormVisiable = false">取消</el-button>
            <el-button type="primary" @click="createAsset()">确定</el-button>
          </span>
    </template>
  </el-dialog>
  <!--updateAsset的form-->
  <el-dialog center width="500px"  v-model="updateAssetFormVisiable" title="update asset">
    <el-form>
      <el-form-item label="名称: ">
        <el-input  v-model="updateAssetData.assetName" />
      </el-form-item>
      <el-form-item label="单价: ">
        <el-input v-model=updateAssetData.assetPrice />
      </el-form-item>
      <el-form-item label="发货地: ">
        <el-input v-model=updateAssetData.shippingAddress />
      </el-form-item>
      <el-form-item label="描述: ">
        <el-input v-model=updateAssetData.publicDescription />
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
    import {request} from "../../api/axios";
    import {ElMessage, ElNotification} from 'element-plus';

    export default {
        name: "assets",
        methods: {
            deleteAsset(assetProxy) {
                let index = this.getAssetIndex(assetProxy.assetID);
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
            beforeUpdateAsset(assetProxy){
              this.selectedAsset = JSON.parse(JSON.stringify(assetProxy));
              this.updateAssetFormVisiable=true;
              this.updateAssetData.assetID=this.selectedAsset.assetID;
              this.updateAssetData.assetName=this.selectedAsset.assetName;
              this.updateAssetData.assetPrice=this.selectedAsset.assetPrice;
              this.updateAssetData.shippingAddress=this.selectedAsset.shippingAddress;
              this.updateAssetData.publicDescription=this.selectedAsset.publicDescription;
            },

            updateAsset(){
              let body = {
                channelName: "mischannel",
                contractName: "mischaincode",
                function: "UpdateAsset",
                args: [
                  this.updateAssetData.assetID,
                  this.updateAssetData.assetName,
                  this.updateAssetData.assetPrice*1,
                  this.updateAssetData.shippingAddress,
                  this.updateAssetData.publicDescription]
              };
              if(this.updateAssetData.assetName.length===0){
                ElMessage({
                  message: 'assetName不能为空',
                  type: 'warning',
                });
                this.updateAssetFormVisiable=false;
                return;
              }
              if(this.updateAssetData.assetPrice*1<0){
                ElMessage({
                  message: 'assetPrice必须>=0',
                  type: 'warning',
                });
                this.updateAssetFormVisiable=false;
                return;
              }
              if(this.updateAssetData.shippingAddress.length===0){
                ElMessage({
                  message: '发货地不能为空',
                  type: 'warning',
                });
                this.updateAssetFormVisiable=false;
                return;
              }
              request('/work/invoke', body, "POST").then(response => {
                ElMessage({
                  message: '修改成功',
                  type: 'success',
                });
                this.updateAssetFormVisiable=false;
              }).catch(error=>{
                this.loading = false;
              })
            },
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


            // get index of asset in orders(sort function causes the returning index is not correct)
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
            },
          createAsset(){
            this.addAssetFormVisiable=true;
            let body={
              channelName: "mischannel",
              contractName: "mischaincode",
              function: "CreateAsset",
              transient:{
                asset:{
                  assetName: this.newAsset.assetName,
                  assetPrice: this.newAsset.assetPrice*1,
                  shippingAddress: this.newAsset.shippingAddress,
                  publicDescription: this.newAsset.publicDescription
                }
              }
            };
            if(this.newAsset.assetName.length===0){
              ElMessage({
                message: 'assetName不能为空',
                type: 'warning',
              });
              this.addAssetFormVisiable=false;
              return;
            }
            if(this.newAsset.assetPrice*1<0){
              ElMessage({
                message: 'assetPrice必须>=0',
                type: 'warning',
              });
              this.addAssetFormVisiable=false;
              return;
            }
            if(this.newAsset.shippingAddress.length===0){
              ElMessage({
                message: '发货地不能为空',
                type: 'warning',
              });
              this.addAssetFormVisiable=false;
              return;
            }
            request('/work/invoke', body, "POST").then(response => {
              ElMessage({
                message: '创建newAsset成功',
                type: 'success',
              });
              console.log('创建newAsset成功');
              this.addAssetFormVisiable=false;
              this.newAsset.assetPrice="";
              this.newAsset.assetName="";
              this.newAsset.shippingAddress="";
              this.newAsset.publicDescription="";
            }).catch(error => {
              console.log('创建newAsset失败')
            });
          }
        },
        data() {
            return {
                assets: [],
                loading: true,
                user: {},
                publishAssetFormVisible: false,
                selectedAsset: {},
                addAssetFormVisiable: false,
// 创建新asset的数据结构
                newAsset:{
                  assetName:"",
                  assetPrice:"",
                  shippingAddress:"",
                  publicDescription:""
              },
              updateAssetFormVisiable:false,
               updateAssetData:{
                  assetID:"",
                 assetName:"",
                 assetPrice:"",
                 shippingAddress:"",
                 publicDescription:""
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