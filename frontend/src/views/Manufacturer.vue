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
        components: {
            LeftSideBar,
            Header
        },
        methods: {
        },
        mounted() {
        }
    }
</script>

<style scoped>

</style>