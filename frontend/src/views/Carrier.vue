<template>
    <Header/>
    <LeftSideBar/>
</template>

<script>
    import LeftSideBar from "../components/CarrierLeftSideBar.vue";
    import Header from "../components/Header.vue"

    export default {
        name: "Carrier",
        components: {
            LeftSideBar,
            Header
        },
    }
</script>


<style scoped>


</style>