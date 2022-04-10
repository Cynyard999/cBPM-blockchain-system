<template>
    <Header/>
    <LeftSideBar/>
</template>

<script>
    import LeftSideBar from "../components/SupplierLeftSideBar.vue";
    import Header from "../components/Header.vue"
    import {request} from "../api/axios";

    export default {
        name: "Supplier",
        components: {
            LeftSideBar,
            Header
        },
        methods: {},
        mounted() {
            // request('/work/query', {}, 'POST').then(response => {
            //     console.log(response)
            // }).catch(error => {
            //     console.log(error)
            // });
        }
    }
</script>

<style scoped>

</style>