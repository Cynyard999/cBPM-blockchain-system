<template>
  <Header/>
  <LeftSideBar/>
</template>

<script>
    import {request} from "../api/axios";
    import LeftSideBar from "../components/ManufacturerLeftSideBar.vue";
    import Header from "../components/Header.vue"
    export default {
        name: 'Home',
      components:{
        LeftSideBar,
        Header
      },
        methods: {
            test() {
                let params = {
                    orgType: "supplier",
                    channelName: "mischannel",
                    contractName: "mischaincode",
                    function: "GetAsset",
                    args: ["8bf609c9-de72-4882-8638-e3b52e2f2fa6"]
                };
                request("/query",params,"POST").then(res=>{
                    console.log(res)
                }).catch(err=>{
                    console.log(err)
                })
            },
        },
        mounted() {
            // this.test()
        }
    }
</script>

<style scoped>

</style>