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
                <li v-show="!isLogin" @click="loginFormVisible = true"><span style="cursor: pointer;">Login</span></li>
                <li v-show="isLogin">
                    <el-tooltip placement="bottom">
                        <template #content>{{userOrg}}</template>
                        <span style="cursor: default;">
                            {{userName}}
                        </span>
                    </el-tooltip>
                </li>
                <li v-show="isLogin" @click="logout">
                    <el-icon>
                        <close/>
                    </el-icon>
                </li>
                <el-dialog title="登录" v-model="loginFormVisible" center width="500px">
                    <el-form :model="userInput">
                        <el-form-item label="Email">
                            <el-input v-model="userInput.email" prefix-icon="User"></el-input>
                        </el-form-item>
                        <el-form-item label="Password">
                            <el-input v-model="userInput.pwd" prefix-icon="Key" show-password></el-input>
                        </el-form-item>
                    </el-form>
                    <template #footer>
                          <span>
                            <el-button @click="loginFormVisible = false">Cancel</el-button>
                            <el-button type="primary" @click="login">Confirm</el-button>
                          </span>
                    </template>
                </el-dialog>
            </ul>
        </header>
    </div>
</template>

<script>
    import {request} from "../api/axios";
    export default {
        name: "Header",
        data() {
            return {
                isLogin: false,
                userName: "",
                userOrg: "",
                userInput: {
                    "email": "",
                    "pwd": ""
                },
                loginFormVisible: false,
            }
        },
        mounted() {
            this.checkUserLogin();
        },
        methods: {
            login() {
                request('user/login', this.userInput, 'POST').then(response => {
                    this.$notify({
                        title: '登录成功',
                        message: response.data.result.name,
                        type: 'success',
                        duration: 2000
                    });
                    this.loginFormVisible = false;
                    this.checkUserLogin();
                }).catch(error => {
                    this.$notify({
                        title: '登录失败',
                        message: error.data.result.message,
                        type: 'success',
                        duration: 2000
                    });
                });
            },
            logout() {
                this.userName = "";
                this.userOrg = "";
                this.isLogin = false;
                window.localStorage.removeItem("token");
                window.localStorage.removeItem("user");
            },

            test() {
            },
            checkUserLogin() {
                let token = window.localStorage.getItem("token");
                if (token) {
                    this.isLogin = true;
                    this.userName = JSON.parse(window.localStorage.getItem("user"))["name"]
                    this.userOrg = JSON.parse(window.localStorage.getItem("user"))["orgType"]
                } else {
                    this.isLogin = false;
                    this.userName = "";
                }
            }
        }
    }
</script>

<style>

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

    header ul li span {
        text-decoration: none;
        font-weight: 300;
        margin: 0 15px;
        color: #253237;
    }

    header ul li a:hover {
        color: #5C6B73;
    }

    header ul li span:hover {
        color: #5C6B73;
    }

    .el-dialog {
        border-radius: 1rem;
    }

    .el-icon {
        top: 2px;
    }

    .el-icon:hover {
        color: #5C6B73;
        cursor: pointer;
    }

</style>