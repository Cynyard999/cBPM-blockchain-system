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
                <li v-show="isAdmin" @click="registerFormVisible = true"><span style="cursor: pointer;">Register</span>
                </li>
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
                    <el-form :model="userInput" label-position="right" label-width="75px">
                        <el-form-item label="邮箱">
                            <el-input v-model="userInput.email" prefix-icon="Message"></el-input>
                        </el-form-item>
                        <el-form-item label="密码">
                            <el-input v-model="userInput.pwd" prefix-icon="Key" show-password></el-input>
                        </el-form-item>
                    </el-form>
                    <template #footer>
                          <span>
                            <el-button @click="loginOver">Cancel</el-button>
                            <el-button type="primary" @click="login">Confirm</el-button>
                          </span>
                    </template>
                </el-dialog>
                <el-dialog title="注册" v-model="registerFormVisible" center width="500px">
                    <el-form :model="adminInput" label-position="right" label-width="75px">
                        <el-form-item label="企业">
                            <el-select v-model="adminInput.orgType"
                                       placeholder="选择所属企业">
                                <el-option
                                        v-for="item in orgOptions"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="邮箱">
                            <el-input v-model="adminInput.email" prefix-icon="Message"></el-input>
                        </el-form-item>
                        <el-form-item label="用户名">
                            <el-input v-model="adminInput.userName" prefix-icon="User"></el-input>
                        </el-form-item>
                        <el-form-item label="密码">
                            <el-input v-model="adminInput.pwd" prefix-icon="Key" show-password></el-input>
                        </el-form-item>
                    </el-form>
                    <template #footer>
                          <span>
                            <el-button @click="registerOver">Cancel</el-button>
                            <el-button type="primary" @click="register">Confirm</el-button>
                          </span>
                    </template>
                </el-dialog>
            </ul>
        </header>
    </div>
</template>

<script>
    import {request} from "../api/axios";
    import router from "../router";
    import {ElMessage} from 'element-plus'
    import md5 from 'js-md5';
    import {ElNotification} from 'element-plus'

    export default {
        name: "Header",
        data() {
            return {
                isLogin: false,
                isAdmin: false,
                userName: "",
                userOrg: "",
                userInput: {
                    email: "",
                    pwd: ""
                },
                adminInput: {
                    email: "",
                    userName: "",
                    pwd: "",
                    orgType: "",
                },
                orgOptions: [
                    {
                        value: 'manufacturer',
                        label: 'manufacturer'
                    },
                    {
                        value: 'carrier',
                        label: 'carrier'
                    },
                    {
                        value: 'middleman',
                        label: 'middleman'
                    },
                    {
                        value: 'supplier',
                        label: 'supplier'
                    },
                ],
                loginFormVisible: false,
                registerFormVisible: false,
            };
        },
        mounted() {
            this.checkUserLogin();
        },
        methods: {
            register() {
                this.adminInput.pwd = md5(this.adminInput.email.charAt(0) + this.adminInput.pwd + this.adminInput.email.charAt(1));
                request('user/register', this.adminInput, 'POST').then(response => {
                    ElNotification({
                        title: '注册成功',
                        message: response.data.result.name,
                        type: 'success',
                        duration: 1000
                    });
                    this.registerOver();
                }).catch(error => {
                    ElNotification({
                        title: '注册失败',
                        message: error.data.message,
                        type: 'error',
                        duration: 1000
                    });
                });
            },
            registerOver() {
                this.registerFormVisible = false;
                this.adminInput = {
                    email: "",
                    userName: "",
                    pwd: "",
                    orgType: ""
                };
            },
            login() {
                let userEncryptedInput = {
                    "email": this.userInput.email,
                    "pwd": ""
                };
                userEncryptedInput.pwd = md5(this.userInput.email.charAt(0) + this.userInput.pwd + this.userInput.email.charAt(1));
                //5c690ccd483bab51498c83fd7fced363
                request('user/login', userEncryptedInput, 'POST').then(response => {
                    ElNotification({
                        title: '登录成功',
                        message: response.data.result.name,
                        type: 'success',
                        duration: 1000
                    });
                    this.loginOver();
                    this.checkUserLogin();
                }).catch(error => {
                    ElNotification({
                        title: '登录失败',
                        message: error.data.message,
                        type: 'error',
                        duration: 1000
                    });
                });
            },
            loginOver() {
                this.loginFormVisible = false;
                this.userInput = {
                    email: "",
                    pwd: ""
                };
            },
            logout() {
                this.userName = "";
                this.userOrg = "";
                this.isLogin = false;
                this.isAdmin = false;
                window.localStorage.removeItem("token");
                window.localStorage.removeItem("user");
                router.replace('/home').then(ElMessage({
                    message: '注销成功',
                    type: 'success',
                }));
            },
            checkUserLogin() {
                let token = window.localStorage.getItem("token");
                if (token) {
                    this.isLogin = true;
                    this.userName = JSON.parse(window.localStorage.getItem("user"))["name"];
                    this.userOrg = JSON.parse(window.localStorage.getItem("user"))["orgType"];
                    if (this.userOrg === "admin") {
                        this.isAdmin = true;
                    }
                } else {
                    this.isLogin = false;
                    this.isAdmin = false;
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
        content: " on Fabric";
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