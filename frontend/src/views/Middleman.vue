<template>
    <div class="header">
        <h1 style="color: white;text-align: center">中间商管理</h1>
    </div>
  <LeftSideBar/>
</template>

<script>
import LeftSideBar from "../components/MiddlemanLeftSideBar.vue";
    export default {
        name: "Middleman.vue",
        components:{
          LeftSideBar
      }
    }
</script>

<style scoped>
.header{
  width: 100%;
  height: 50px;
  background-color: #343a40;
}
</style>