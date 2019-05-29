import React from 'react'
import ReactDOM from 'react-dom'
import Routers from './routers'
import { Provider } from 'react-redux'
import store from './store'
import axios from 'axios'
import { baseUrl } from './api/homeApi'


axios.defaults.baseURL = baseUrl // 全局配置请求地址
// axios.defaults.headers.common['Authorization'] = AUTH_TOKEN
// axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'

// request拦截
axios.interceptors.request.use(function(config) {
    // 在发送请求之前执行的代码
    return config
}, function(error) {
    // 请求出错以后执行的代码
    return Promise.reject(error)
})

// response拦截
axios.interceptors.response.use(function(response) {
    // 接收到响应以后执行的代码
    // console.log('请求接收以后拦截了')
    // console.log(response)
    if(response.data.err_code === 10303) {
        window.location.href = baseUrl + '/saml/login?RelayState=/'
        return
    }
    return response
}, function(error) {
    // 接收出错以后执行的代码
    return Promise.reject(error)
})

const rootNode = document.getElementById('root')

ReactDOM.render(
    <Provider store={store}> 
        <Routers />
    </Provider>, rootNode)