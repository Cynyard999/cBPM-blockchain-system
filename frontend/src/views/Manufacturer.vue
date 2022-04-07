<template>
    <div>
        <h1>Manufacturer</h1>
    </div>
</template>

<script>
    import {request} from "../api/axios";

    export default {
        name: 'Home',
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