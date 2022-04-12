<template>
  <el-table
      :data="deliveryDetails"
      :default-sort="{ prop: 'createTime', order: 'descending' }"
      v-loading="loading"
      ref="tableRef"
      style='width: 100%'
      height="600px"
  >
    <el-table-column prop="createTime" label="创建时间"/>
    <el-table-column prop="assetName" label="名称"/>
    <el-table-column prop="note" label="备注" overflow/>
<!--    <el-table-column prop="ownerOrg" label="所属组织"/>-->
    <el-table-column prop="contact" label="contact"/>
    <el-table-column prop="startPlace" label="发货地"/>
    <el-table-column prop="endPlace" label="收货地"/>
    <el-table-column prop="updateTime" label="修改时间"/>
    <el-table-column
        prop="tag"
        label="状态"
        width="100"
        :filters="[{ text: '未处理', value: 0 },{ text: '开始处理', value: 1 },{ text: '处理完成', value: 2 }]"
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
        <el-button v-if="this.user.orgType==='carrier'" type="text" size="small" @click="getDeliveryDetailForm(scope.row)">
          修改状态
        </el-button>
        <el-button v-if="this.user.orgType==='manufacturer'" type="text" size="small" @click="getOrderWithMiddleman(scope.row)">
          查看
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-dialog center width="500px" v-model="deliveryDetailFormVisible" title="Change Detail Status Confirm">
    <el-form label-position="right" label-width="187px">
      <el-form-item label="创建时间: ">
        {{selectedDeliveryDetail.createTime}}
      </el-form-item>
      <el-form-item label="商品名称: ">
        {{selectedDeliveryDetail.assetName}}
      </el-form-item>
      <el-form-item label="contact: ">
        {{selectedDeliveryDetail.contact}}
      </el-form-item>
      <el-form-item label="发货地: ">
        {{selectedDeliveryDetail.startPlace}}
      </el-form-item>
      <el-form-item label="收货地: ">
        {{selectedDeliveryDetail.endPlace}}
      </el-form-item>
      <el-form-item label="备注: ">
        {{selectedDeliveryDetail.note}}
      </el-form-item>
      <el-select style="margin-left: 120px" v-model="selectedDeliveryDetail.newStatus" placeholder="Select">
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
            <el-button @click="deliveryDetailFormVisible = false">取消</el-button>
            <el-button type="primary" @click="changeDeliveryDetailStatus()">确定</el-button>
          </span>
    </template>
  </el-dialog>

<el-dialog center width="500px" v-model="orderDetailFormVisible" title="show orderDetail">
  <el-form label-position="right" label-width="187px">
    <el-form-item label="创建时间: ">
      {{orderDetails.createTime}}
    </el-form-item>
    <el-form-item label="商品名称: ">
      {{orderDetails.assetName}}
    </el-form-item>
    <el-form-item label="数量: ">
      {{orderDetails.quantity}}
    </el-form-item>
    <el-form-item label="发货地: ">
      {{orderDetails.shippingAddress}}
    </el-form-item>
    <el-form-item label="收货地: ">
      {{orderDetails.receivingAddress}}
    </el-form-item>
    <el-form-item label="备注: ">
      {{orderDetails.note}}
    </el-form-item>
  </el-form>
  <template #footer>
          <span>
            <el-button @click="orderDetailFormVisible = false">完成</el-button>
          </span>
  </template>
</el-dialog>
</template>

<script>
import {request} from "../../api/axios";
import {ElMessage, ElNotification} from 'element-plus';

export default {
  name: "DeliveryDetails",
  methods: {
    //获取对应order的详细信息以方便生产商确认收到货
    getOrderWithMiddleman(orderProxy){
      this.selectedDeliveryDetail = JSON.parse(JSON.stringify(orderProxy));
      let body={
        channelName: "mamichannel",
        contractName: "mamichaincode",
        function: "GetOrder",
        args:[this.selectedDeliveryDetail.tradeID]
      };
      let that = this;
      this.loading=true;
      request('/work/query', body, "POST").then(response => {
        that.orderDetails = response.data.result;
        that.loading = false;
        that.orderDetailFormVisible=true;
      }).catch(error => {
        that.loading = false;
      });
    },
    getDeliveryDetailForm(orderProxy) {
      this.selectedDeliveryDetail = JSON.parse(JSON.stringify(orderProxy));
      this.selectedDeliveryDetail.newStatus = this.selectedDeliveryDetail.status;
      this.deliveryDetailFormVisible = true;
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
    // get index of asset in Details(sort function causes the returning index is not correct)
    getDeliveryDetailsIndex(tradeID) {
      let index = -1;
      this.details.forEach((order, i) => {
        if (tradeID === asset["tradeID"]) {
          index = i;
        }
      });
      return index;
    },
    checkDisabled(status) {
      if (status - 1 > this.selectedDeliveryDetail.status) {
        return true;
      }
      if (status < this.selectedDeliveryDetail.status) {
        return true;
      }
      if (this.user.orgType === 'manufacturer') {
        //生产商无法操作deliveryDetail
        return true;
      } else {
        return !(status === 1 || status === 2 || status === 0);
      }
    },
    changeDeliveryDetailStatus() {
      this.deliveryDetailFormVisible = false;
      if (this.selectedDeliveryDetail.status === this.selectedDeliveryDetail.newStatus) {
        ElMessage({
          message: '修改成功',
          type: 'success',
        });
        this.getOrders();

      } else {
        let body = {
          channelName: "cmachannel",
          contractName: "cmachaincode",
          function: "",
          args: [this.selectedDeliveryDetail.tradeID]
        };
        if (this.selectedDeliveryDetail.newStatus === 1) {
          body.function = "HandleDeliveryDetail";
        }
        if (this.selectedDeliveryDetail.newStatus === 2) {
          body.function = "FinishDeliveryDetail";
        }
        if (this.selectedDeliveryDetail.newStatus === 3) {
          body.function = "ConfirmFinishDeliveryDetail";
        }
        let that = this;
        that.loading = true;

        request('/work/invoke', body, "POST").then(response => {
          ElMessage({
            message: '修改成功',
            type: 'success',
          });
          that.loading = false;
          that.getDeliveryDetails();
        }).catch(error => {
          that.loading = false;
        });
      }
    },
    getDeliveryDetails() {
      let body = {
        channelName: "cmachannel",
        contractName: "cmachaincode",
        function: "GetAllDeliveryDetails",
        args: []
      };
      let that = this;
      this.loading=true;
      request('/work/query', body, "POST").then(response => {
        that.deliveryDetails = response.data.result;
        that.loading = false;
      }).catch(error => {
        that.loading = false;
      });
    },
    getUser() {
      this.user = JSON.parse(window.localStorage.getItem("user"));
    },

  },
  data() {
    return {
      deliveryDetails: [],
      orderDetails:[],
      loading: true,
      user: {},
      deliveryDetailFormVisible: false,
      orderDetailFormVisible: false,
      selectedDeliveryDetail: {},
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
        }
      ]
    }
  },
  mounted() {
    this.getUser();
    this.getDeliveryDetails();
  }
}
</script>

<style>

</style>