<template>
    <div class="network-container">
        <div class="network-circle">
            <div class="card manufacturer">
                <h3 class="org-name">Manufacturer</h3>
                <p class="org-link">
                    <a href="/manufacturer/middleman">联系中间商</a>
                </p>
                <p class="org-link">
                    <a href="/manufacturer/carrier">联系运输商</a>
                </p>
            </div>
            <div class="card middleman">
                <h3 class="org-name">Middleman</h3>
                <p class="org-link">
                    <a href="/supplier/carrier">联系运输商</a>
                </p>
                <p class="org-link">
                    <a href="/supplier/manufacturer">联系生产商</a>
                </p>
                <p class="org-link">
                    <a href="/supplier/supplier">联系供应商</a>
                </p>
            </div>
            <div class="card supplier">
                <h3 class="org-name">Supplier</h3>
                <p class="org-link">
                    <a href="/supplier/middleman">联系中间商</a>
                </p>
                <p class="org-link">
                    <a href="/supplier/carrier">联系运输商</a>
                </p>
            </div>
            <div class="card carrier">
                <h3 class="org-name">Carrier</h3>
                <p class="org-link">
                    <a href="/supplier/middleman">联系中间商</a>
                </p>
                <p class="org-link">
                    <a href="/supplier/supplier">联系供应商</a>
                </p>
                <p class="org-link">
                    <a href="/supplier/manufacturer">联系生产商</a>
                </p>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Network"
    }
</script>

<style scoped>
    .network-container {
        position: center;
    }

    .network-container .network-circle {
        width: 80%;
        height: 60%;
        border-radius: 50%;
        border: 10px dashed #5f5f5f;
        position: absolute;
        top: 50%;
        left: 50%;
        -webkit-transform: translate(-50%, -50%);
        -ms-transform: translate(-50%, -50%);
        transform: translate(-50%, -50%);
    }

    .network-circle .card.manufacturer {
        bottom: 80%;
        left: 50%;
        -webkit-transform: translateX(-50%);
        -ms-transform: translateX(-50%);
        transform: translateX(-50%);
    }

    .network-circle .card.middleman {
        top: 50%;
        right: 0;
        -webkit-transform: translateY(-50%);
        -ms-transform: translateY(-50%);
        transform: translateY(-50%);
    }

    .network-circle .card.supplier {
        top: 80%;
        left: 50%;
        -webkit-transform: translateX(-50%);
        -ms-transform: translateX(-50%);
        transform: translateX(-50%);
    }

    .network-circle .card.carrier {
        left: 0;
        top: 50%;
        -webkit-transform: translateY(-50%);
        -ms-transform: translateY(-50%);
        transform: translateY(-50%);
    }

    .network-circle .card {

        position: absolute;
    }

    .card {
        background-color: transparent;
        border: none;
        -webkit-box-sizing: border-box;
        box-sizing: border-box;
        margin-bottom: 20px;
        position: relative;
    }

    a {
        color: #3b6caa;
        text-decoration: none;
    }

    .org-name {
        font-size: 1.5625rem;
        line-height: 1.875rem;
    }

    p.org-link {
        padding-bottom: 12px;
    }
</style>