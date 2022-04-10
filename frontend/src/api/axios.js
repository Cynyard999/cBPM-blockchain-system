import axios from 'axios';
import {ShowMessage} from "./status";   // 引入状态码文件
import {ElMessage} from 'element-plus'
import 'element-plus/dist/index.css'

// 设置接口超时时间
axios.defaults.timeout = 60000;

//http request 拦截器
axios.interceptors.request.use(
    config => {
        // 配置请求头
        config.headers['Content-Type'] = 'application/json;charset=UTF-8';
        let token = window.localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = token;
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

//http response 拦截器
axios.interceptors.response.use(
    response => {
        // if (response.status === 401) {
        //     window.localStorage.removeItem('token');
        //     window.localStorage.removeItem('userInfo');
        //     router.replace('/home').then(r => console.log("switch to home"));
        //     return
        // }
        let token = response.headers.authorization;
        if (token) {
            console.log(response.data.result);
            window.localStorage.setItem('token', token);
            window.localStorage.setItem('user', JSON.stringify(response.data.result));
        }
        return response;
    },
    error => {
        const {response} = error;
        if (response) {// 传入响应码，匹配响应码对应信息
            ElMessage.warning(ShowMessage(response.status) + ": " + response.data.message);
            return Promise.reject(response);
        } else {
            ElMessage.warning('网络连接异常,请稍后再试!');
        }
    }
);

// 封装 GET POST 请求并导出
export function request(url = '', params = {}, type = 'POST') {
//设置 url params type 的默认值
    return new Promise((resolve, reject) => {
        let promise;
        if (type.toUpperCase() === 'GET') {
            promise = axios({
                url,
                params
            })
        } else if (type.toUpperCase() === 'POST') {
            promise = axios({
                method: 'POST',
                url,
                data: params
            })
        }
        //处理返回
        promise.then(res => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}

