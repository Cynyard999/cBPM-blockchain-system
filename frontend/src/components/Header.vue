<template>
        <div class="top">
            <header>
                <div class="web-name">CBPM</div>
                <ul>
                    <li><a href="/home">Home</a></li>
                    <li><a href="/manufacturer">Manufacturer</a></li>
                    <li><a href="/carrier">Carrier</a></li>
                    <li><a href="/supplier">Supplier</a></li>
                    <li><a href="/middleman">Middleman</a></li>
                </ul>
            </header>
        </div>
</template>

<script>
    export default {
        name: "Header"
    }
</script>

<style scoped>

    .top {
        background-color: #ffffff;
        width: 100%;
        height: 98px;
        /*position: fixed; 这样可以navbar一直在最上面*/
        position: sticky;
        box-shadow: 0 1.5px 2px 0 rgba(0, 0, 0, 0.175);
        z-index: 10;
    }

    header {
        position: sticky;
        top: 0;
        left: 0;
        /*width: 100%;*/
        justify-content: space-between;
        align-items: center;
        display: flex;
        padding: 20px 40px;
        /*z-index: 10;*/
    }

    header .web-name {
        position: relative;
        font-size: 2em;
        text-decoration: none;
        color: #323845;
        text-transform: uppercase;
        cursor: default;
        transition: 0.6s;
        font-weight: 500;
        user-select: none;
        letter-spacing: 5px;
    }

    header .web-name:after {
        content: " on Fabirc";
        height: 100%;
        text-decoration: none;
        color: #5f5f5f;
        text-transform: uppercase;
        font-weight: 300;
        letter-spacing: 3px;
        font: italic 0.6em Georgia, serif;
    }

    header ul {
        position: relative;
        justify-content: center;
        display: flex;
        align-items: center;
    }


    header ul li {
        position: relative;
        list-style-type: none;
        user-select: none;
    }

    header ul li a {
        position: relative;
        text-decoration: none;
        font-weight: 300;

        margin: 0 15px;
        color: #323845;
        letter-spacing: 1px;
    }

    header ul li a:hover {
        color: #e74a3b;
        font-weight: 400;
    }
</style>