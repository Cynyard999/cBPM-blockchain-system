<template>
  <el-table
      :data="deliveryArrangements"
      :default-sort="{ prop: 'createTime', deliveryArrangement: 'descending' }"
      v-loading="loading"
      ref="tableRef"
      style='width: 100%'
      height="600px"
  >
    <el-table-column prop="createTime" label="创建时间"/>
    <el-table-column prop="assetName" label="名称"/>
    <el-table-column prop="quantity" label="数量"/>
    <el-table-column prop="fee" label="运费"/>
    <el-table-column prop="startPlace" label="发货地"/>
    <el-table-column prop="endPlace" label="收货地"/>
    <el-table-column prop="note" label="备注" overflow/>
    <el-table-column prop="ownerOrg" label="所属组织"/>
    <el-table-column prop="handler" label="处理组织"/>
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
        <el-button type="text" size="small" @click="getDeliveryArrangementForm(scope.row)">
          修改状态
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-dialog center width="500px" v-model="deliveryArrangementFormVisible" title="Change deliveryArrangement Status Confirm">
    <el-form label-position="right" label-width="135px">
      <el-form-item label="创建时间: ">
        {{selectedDeliveryArrangement.createTime}}
      </el-form-item>
      <el-form-item label="商品名称: ">
        {{selectedDeliveryArrangement.assetName}}
      </el-form-item>
      <el-form-item label="数量: ">
        {{selectedDeliveryArrangement.quantity}}
      </el-form-item>
      <el-form-item label="运费: ">
        {{selectedDeliveryArrangement.fee}}
      </el-form-item>
      <el-form-item label="发货地: ">
        {{selectedDeliveryArrangement.startPlace}}
      </el-form-item>
      <el-form-item label="收货地: ">
        {{selectedDeliveryArrangement.endPlace}}
      </el-form-item>
      <el-form-item label="备注: ">
        {{selectedDeliveryArrangement.note}}
      </el-form-item>
      <el-form-item v-show="noteDeliveryDetailVisible" label="DeliveryDetailNote:">
        <el-input placeholder="note for DeliveryDetail" v-model=this.noteDeliveryDetail />
      </el-form-item>
      <el-select @change="showNoteDeliveryDetail" style="margin-left: 135px" v-model="selectedDeliveryArrangement.newStatus" placeholder="Select">
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
            <el-button @click="deliveryArrangementFormVisible = false">取消</el-button>
            <el-button type="primary" @click="changeDeliveryArrangementStatus()">确定</el-button>
          </span>
    </template>
  </el-dialog>
</template>

<script>
import {request} from "../../api/axios";
import {ElMessage, ElNotification} from 'element-plus';

export default {
  name: "DeliveryArrangements",
  methods: {
    showNoteDeliveryDetail(){
      if(this.selectedDeliveryArrangement.newStatus===1){
        this.noteDeliveryDetailVisible=true;
      }else{
        this.noteDeliveryDetailVisible=false;
      }
    },
    getDeliveryArrangementForm(deliveryArrangementProxy) {
      this.selectedDeliveryArrangement = JSON.parse(JSON.stringify(deliveryArrangementProxy));
      this.selectedDeliveryArrangement.newStatus = this.selectedDeliveryArrangement.status;
      this.deliveryArrangementFormVisible = true;
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
    // get index of asset in deliveryArrangements(sort function causes the returning index is not correct)
    getDeliveryArrangementIndex(tradeID) {
      let index = -1;
      this.deliveryArrangements.forEach((deliveryArrangement, i) => {
        if (tradeID === asset["tradeID"]) {
          index = i;
        }
      });
      return index;
    },
    checkDisabled(status) {
      if (status - 1 > this.selectedDeliveryArrangement.status) {
        return true;
      }
      if (status < this.selectedDeliveryArrangement.status) {
        return true;
      }
      if (this.user.orgType === 'middleman') {
        return !(status === 3 || status === 2);
      } else {
        return !(status === 1 || status === 2 || status === 0);
      }
    },
    changeDeliveryArrangementStatus() {
      this.deliveryArrangementFormVisible = false;
      if (this.selectedDeliveryArrangement.status === this.selectedDeliveryArrangement.newStatus) {
        ElMessage({
          message: '修改成功',
          type: 'success',
        });
        this.getDeliveryArrangements();

      } else {
        let body = {
          channelName: "micchannel",
          contractName: "micchaincode",
          function: "",
          args: [this.selectedDeliveryArrangement.tradeID]
        };
        //deliveryOrder的修改status的body
        let bodyDeliveryOeder = {
          channelName: "scchannel",
          contractName: "scchaincode",
          function: "",
          args: [this.selectedDeliveryArrangement.tradeID]
        };

        if (this.selectedDeliveryArrangement.newStatus === 1) {
          body.function = "HandleDeliveryArrangement";
          bodyDeliveryOeder.function="HandleDeliveryOrder";
        }
        if (this.selectedDeliveryArrangement.newStatus === 2) {
          body.function = "FinishDeliveryArrangement";
          bodyDeliveryOeder.function="FinishDeliveryOrder";
        }
        if (this.selectedDeliveryArrangement.newStatus === 3) {
          body.function = "ConfirmFinishDeliveryArrangement";
        }
        let that = this;
        that.loading = true;
        //handle的时候顺便创造deliveryDetail
        if (this.selectedDeliveryArrangement.newStatus === 1) {
          this.createDeliveryDetail();
        }
        //同时修改deliveryOrder的status
        if(this.selectedDeliveryArrangement.newStatus !== 3) {
          request('/work/invoke', bodyDeliveryOeder, "POST").then(response => {
            ElMessage({
              message: '修改DeliveryStatus成功',
              type: 'success',
            });
          }).catch(error => {
            that.loading = false;
          });
        }
        request('/work/invoke', body, "POST").then(response => {
          ElMessage({
            message: '修改成功',
            type: 'success',
          });
          that.loading = false;
          that.noteDeliveryDetailVisible=false;
          that.getDeliveryArrangements();
        }).catch(error => {
          that.loading = false;
        });


      }
    },
    getDeliveryArrangements() {
      let body = {
        channelName: "micchannel",
        contractName: "micchaincode",
        function: "GetAllDeliveryArrangements",
        args: []
      };
      let that = this;
      this.loading = true;
      request('/work/query', body, "POST").then(response => {
        that.deliveryArrangements = response.data.result;

        that.loading = false;
      }).catch(error => {
        that.loading = false;
      });
    },
    getUser() {
      this.user = JSON.parse(window.localStorage.getItem("user"));
    },

    createDeliveryDetail(){
      let body={
        channelName: "cmachannel",
        contractName: "cmachaincode",
        function:  "CreateDeliveryDetail",
        transient:{
          detail:{
            TradeId: this.selectedDeliveryArrangement.tradeID,
            AssetName: this.selectedDeliveryArrangement.assetName,
            StartPlace: this.selectedDeliveryArrangement.startPlace,
            EndPlace: this.selectedDeliveryArrangement.endPlace,
            Contact: "beijing",
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
      deliveryArrangements: [],
      loading: true,
      user: {},
      deliveryArrangementFormVisible: false,
      noteDeliveryDetail: "",
      noteDeliveryDetailVisible: false,
      selectedDeliveryArrangement: {},
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
    this.getDeliveryArrangements();
  }
}
</script>

<style>

</style>