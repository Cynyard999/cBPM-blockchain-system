<template>
    <div class="header">
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

    .header {
        background-color: #ffffff;
        width: 100%;
        height: 10%;
        /*position: fixed; 这样可以navbar一直在最上面*/
        position: sticky;
        box-shadow: 0 1.5px 2px 0 rgba(0, 0, 0, 0.175);
        z-index: 10;
    }

    header {
        top: 0;
        left: 0;
        justify-content: space-between;
        align-items: center;
        display: flex;
        padding: 20px 40px;
    }

    header .web-name {
        position: relative;
        font-size: 2em;
        text-decoration: none;
        color: #253237;
        text-transform: uppercase;
        cursor: default;
        font-weight: 400;
        user-select: none;
        letter-spacing: 5px;
    }

    header .web-name:after {
        content: " on Fabirc";
        height: 100%;
        text-decoration: none;
        color: #5C6B73;
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
        text-decoration: none;
        font-weight: 300;
        margin: 0 15px;
        color: #253237;
    }

    header ul li a:hover {
        color: #5C6B73;
    }

</style>