<template>
    <Header></Header>
    <div class="content">
        <div class="text">选择你的企业角色</div>
        <div class="container">
            <a href="/manufacturer" class="link">
                <el-card shadow="hover">
                    <h3 class="org-name">
                        Manufacturer
                    </h3>
                    <img src="manufacturer-background.png" id="manufacturer-img" alt="manufacturer">
                </el-card>
            </a>
            <a href="/middleman" class="link">
                <el-card shadow="hover">
                    <h3 class="org-name">
                        Middleman
                    </h3>
                    <img src="middleman-background.png" id="middleman-img" alt="middleman">
                </el-card>
            </a>
            <a href="/supplier" class="link">
                <el-card shadow="hover">
                    <h3 class="org-name">
                        Supplier
                    </h3>
                    <img src="supplier-background.png" id="supplier-img" alt="supplier">
                </el-card>
            </a>
            <a href="/carrier" class="link">
                <el-card shadow="hover">
                    <h3 class="org-name">
                        Carrier
                    </h3>
                    <img src="carrier-background.png" id="carrier-img" alt="carrier">
                </el-card>
            </a>
        </div>
    </div>
</template>

<script>
    import Header from "../components/Header.vue"

    export default {
        name: "Home",
        components: {Header},
    }
</script>

<style scoped>
    .content {
    }

    #manufacturer-img {
        opacity: 60%;
        position: relative;
    }

    #carrier-img {
        opacity: 60%;
        position: relative;
    }

    #middleman-img {
        opacity: 60%;
        position: relative;
    }

    #supplier-img {
        opacity: 60%;
        position: relative;
    }

    .org-name {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
        user-select: none;
        font-size: 30px;
        font-weight: normal;
    }

    .text {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
        user-select: none;
        text-align: center;
        margin-top: 4%;
        margin-left: auto;
        margin-right: auto;
        font-size: 40px;
        line-height: 1.07143;
        position: relative;
    }

    .link {
        text-decoration: none;
        border: none;
    }

    .container {
        display: grid;
        place-items: center;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 1vh;
        margin-top: 5%;
    }

    .el-card {
        height: 388px;
        width: 388px;
        border-radius: 1rem;
        text-align: center;
    }

</style>