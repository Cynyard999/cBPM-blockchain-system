<template>
  <Header/>
  <LeftSideBar/>
</template>

<script>
import LeftSideBar from "../components/MiddlemanLeftSideBar.vue";
import Header from "../components/Header.vue"
    export default {
        name: "Middleman.vue",
        components:{
          LeftSideBar,
          Header
      }
    }
</script>

<style scoped>

</style>